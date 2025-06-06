// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package certloader

import (
	"log/slog"
	"sync"

	"github.com/cilium/cilium/pkg/fswatcher"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/time"
)

const watcherEventCoalesceWindow = 100 * time.Millisecond

// Watcher is a set of TLS configuration files including CA files, and a
// certificate along with its private key. The files are watched for change and
// reloaded automatically.
type Watcher struct {
	*FileReloader
	log       *slog.Logger
	fswatcher *fswatcher.Watcher
	stop      chan struct{}
}

// NewWatcher returns a Watcher that watch over the given file
// paths. The given files are expected to already exists when this function is
// called. On success, the returned Watcher is ready to use.
func NewWatcher(log *slog.Logger, caFiles []string, certFile, privkeyFile string) (*Watcher, error) {
	r, err := NewFileReloaderReady(caFiles, certFile, privkeyFile)
	if err != nil {
		return nil, err
	}
	// An error here would be unexpected as we were able to create a
	// FileReloader having read the files, so the files should exist and be
	// "watchable".
	fswatcher, err := newFsWatcher(log, caFiles, certFile, privkeyFile)
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		FileReloader: r,
		log:          log,
		fswatcher:    fswatcher,
		stop:         make(chan struct{}),
	}

	w.Watch()
	return w, nil
}

// FutureWatcher returns a channel where exactly one Watcher will be sent once
// the given files are ready and loaded. This can be useful when the file paths
// are well-known, but the files themselves don't exist yet. Note that the
// requirement is that the file directories must exists.
func FutureWatcher(log *slog.Logger, caFiles []string, certFile, privkeyFile string) (<-chan *Watcher, error) {
	r, err := NewFileReloader(caFiles, certFile, privkeyFile)
	if err != nil {
		return nil, err
	}
	fswatcher, err := newFsWatcher(log, caFiles, certFile, privkeyFile)
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		FileReloader: r,
		log:          log,
		fswatcher:    fswatcher,
		stop:         make(chan struct{}),
	}

	res := make(chan *Watcher)
	go func(res chan<- *Watcher) {
		defer close(res)
		// Attempt a reload without having received any fs notification in case
		// all the files are already there. Note that the keypair and CA are
		// reloaded separately as a "partial update" is still useful: If the
		// FileReloader is "half-ready" (e.g. has loaded the keypair but failed
		// to load the CA), we only need a successfully handled CA related fs
		// notify event to become Ready (in other words, we don't need to
		// receive a fs event for the keypair in that case to become ready).
		_, keypairErr := w.ReloadKeypair()
		_, caErr := w.ReloadCA()
		ready := w.Watch()
		if keypairErr == nil && caErr == nil {
			log.Debug("TLS configuration ready")
			res <- w
			return
		}
		log.Debug("Waiting on fsnotify update to be ready")
		select {
		case <-ready:
			log.Debug("TLS configuration ready")
			res <- w
		case <-w.stop:
		}
	}(res)

	return res, nil
}

// Watch initialize the files watcher and update goroutine. It returns a ready
// channel that will be closed once an update made the underlying FileReloader
// ready.
func (w *Watcher) Watch() <-chan struct{} {
	// prepare the ready channel to be returned. We will close it exactly once.
	var once sync.Once
	ready := make(chan struct{})
	markReady := func() {
		once.Do(func() {
			close(ready)
		})
	}

	// build maps for the CA files and keypair files to help detecting what has
	// changed in order to reload only the appropriate certificates.
	keypairMap := make(map[string]struct{})
	caMap := make(map[string]struct{})
	if w.FileReloader.certFile != "" {
		keypairMap[w.FileReloader.certFile] = struct{}{}
	}
	if w.FileReloader.privkeyFile != "" {
		keypairMap[w.FileReloader.privkeyFile] = struct{}{}
	}
	for _, path := range w.FileReloader.caFiles {
		caMap[path] = struct{}{}
	}

	// used to coalesce fswatcher events that arrive within the same time window
	var keypairReload, caReload <-chan time.Time

	go func() {
		defer w.fswatcher.Close()
		for {
			select {
			case event := <-w.fswatcher.Events:
				path := event.Name
				w.log.Debug("Received fswatcher event",
					logfields.Path, path,
					logfields.Operation, event.Op,
				)

				_, keypairUpdated := keypairMap[path]
				_, caUpdated := caMap[path]

				if keypairUpdated {
					if keypairReload == nil {
						keypairReload = time.After(watcherEventCoalesceWindow)
					}
				} else if caUpdated {
					if caReload == nil {
						caReload = time.After(watcherEventCoalesceWindow)
					}
				} else {
					// fswatcher should never send events for unknown files
					w.log.Warn("Unknown file, ignoring.",
						logfields.Path, path,
						logfields.Operation, event.Op,
					)
					continue
				}
			case <-keypairReload:
				keypairReload = nil

				keypair, err := w.ReloadKeypair()
				if err != nil {
					w.log.Warn("Keypair update failed", logfields.Error, err)
					continue
				}
				id := keypairId(keypair)
				w.log.Info("Keypair updated", logfields.KeyPairSN, id)
				if w.Ready() {
					markReady()
				}
			case <-caReload:
				caReload = nil

				if _, err := w.ReloadCA(); err != nil {
					w.log.Warn("Certificate authority update failed", logfields.Error, err)
					continue
				}
				w.log.Info("Certificate authority updated")
				if w.Ready() {
					markReady()
				}
			case err := <-w.fswatcher.Errors:
				w.log.Warn("fswatcher error", logfields.Error, err)
			case <-w.stop:
				w.log.Info("Stopping fswatcher")
				return
			}
		}
	}()

	return ready
}

// Stop watching the files.
func (w *Watcher) Stop() {
	select {
	case <-w.stop:
	default:
		close(w.stop)
	}
}

// newFsWatcher returns a fswatcher.Watcher watching over the given files.
// The fswatcher.Watcher supports watching over files which do not exist yet.
// A create event will be emitted once the file is added.
func newFsWatcher(logger *slog.Logger, caFiles []string, certFile, privkeyFile string) (*fswatcher.Watcher, error) {
	trackFiles := []string{}

	if certFile != "" {
		trackFiles = append(trackFiles, certFile)
	}
	if privkeyFile != "" {
		trackFiles = append(trackFiles, privkeyFile)
	}
	for _, path := range caFiles {
		if path != "" {
			trackFiles = append(trackFiles, path)
		}
	}

	return fswatcher.New(logger, trackFiles)
}
