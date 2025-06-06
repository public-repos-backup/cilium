//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ExpressRouteConnectionsClient contains the methods for the ExpressRouteConnections group.
// Don't use this type directly, use NewExpressRouteConnectionsClient() instead.
type ExpressRouteConnectionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewExpressRouteConnectionsClient creates a new instance of ExpressRouteConnectionsClient with the specified values.
//   - subscriptionID - The subscription credentials which uniquely identify the Microsoft Azure subscription. The subscription
//     ID forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewExpressRouteConnectionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ExpressRouteConnectionsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ExpressRouteConnectionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates a connection between an ExpressRoute gateway and an ExpressRoute circuit.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - expressRouteGatewayName - The name of the ExpressRoute gateway.
//   - connectionName - The name of the connection subresource.
//   - putExpressRouteConnectionParameters - Parameters required in an ExpressRouteConnection PUT operation.
//   - options - ExpressRouteConnectionsClientBeginCreateOrUpdateOptions contains the optional parameters for the ExpressRouteConnectionsClient.BeginCreateOrUpdate
//     method.
func (client *ExpressRouteConnectionsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, putExpressRouteConnectionParameters ExpressRouteConnection, options *ExpressRouteConnectionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ExpressRouteConnectionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, expressRouteGatewayName, connectionName, putExpressRouteConnectionParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ExpressRouteConnectionsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ExpressRouteConnectionsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates a connection between an ExpressRoute gateway and an ExpressRoute circuit.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
func (client *ExpressRouteConnectionsClient) createOrUpdate(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, putExpressRouteConnectionParameters ExpressRouteConnection, options *ExpressRouteConnectionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ExpressRouteConnectionsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, connectionName, putExpressRouteConnectionParameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ExpressRouteConnectionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, putExpressRouteConnectionParameters ExpressRouteConnection, options *ExpressRouteConnectionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}/expressRouteConnections/{connectionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	if connectionName == "" {
		return nil, errors.New("parameter connectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{connectionName}", url.PathEscape(connectionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, putExpressRouteConnectionParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes a connection to a ExpressRoute circuit.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - expressRouteGatewayName - The name of the ExpressRoute gateway.
//   - connectionName - The name of the connection subresource.
//   - options - ExpressRouteConnectionsClientBeginDeleteOptions contains the optional parameters for the ExpressRouteConnectionsClient.BeginDelete
//     method.
func (client *ExpressRouteConnectionsClient) BeginDelete(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, options *ExpressRouteConnectionsClientBeginDeleteOptions) (*runtime.Poller[ExpressRouteConnectionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, expressRouteGatewayName, connectionName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ExpressRouteConnectionsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ExpressRouteConnectionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes a connection to a ExpressRoute circuit.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
func (client *ExpressRouteConnectionsClient) deleteOperation(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, options *ExpressRouteConnectionsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ExpressRouteConnectionsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, connectionName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ExpressRouteConnectionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, options *ExpressRouteConnectionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}/expressRouteConnections/{connectionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	if connectionName == "" {
		return nil, errors.New("parameter connectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{connectionName}", url.PathEscape(connectionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified ExpressRouteConnection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - expressRouteGatewayName - The name of the ExpressRoute gateway.
//   - connectionName - The name of the ExpressRoute connection.
//   - options - ExpressRouteConnectionsClientGetOptions contains the optional parameters for the ExpressRouteConnectionsClient.Get
//     method.
func (client *ExpressRouteConnectionsClient) Get(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, options *ExpressRouteConnectionsClientGetOptions) (ExpressRouteConnectionsClientGetResponse, error) {
	var err error
	const operationName = "ExpressRouteConnectionsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, connectionName, options)
	if err != nil {
		return ExpressRouteConnectionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExpressRouteConnectionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExpressRouteConnectionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ExpressRouteConnectionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, connectionName string, options *ExpressRouteConnectionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}/expressRouteConnections/{connectionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	if connectionName == "" {
		return nil, errors.New("parameter connectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{connectionName}", url.PathEscape(connectionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ExpressRouteConnectionsClient) getHandleResponse(resp *http.Response) (ExpressRouteConnectionsClientGetResponse, error) {
	result := ExpressRouteConnectionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExpressRouteConnection); err != nil {
		return ExpressRouteConnectionsClientGetResponse{}, err
	}
	return result, nil
}

// List - Lists ExpressRouteConnections.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - expressRouteGatewayName - The name of the ExpressRoute gateway.
//   - options - ExpressRouteConnectionsClientListOptions contains the optional parameters for the ExpressRouteConnectionsClient.List
//     method.
func (client *ExpressRouteConnectionsClient) List(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, options *ExpressRouteConnectionsClientListOptions) (ExpressRouteConnectionsClientListResponse, error) {
	var err error
	const operationName = "ExpressRouteConnectionsClient.List"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, options)
	if err != nil {
		return ExpressRouteConnectionsClientListResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExpressRouteConnectionsClientListResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExpressRouteConnectionsClientListResponse{}, err
	}
	resp, err := client.listHandleResponse(httpResp)
	return resp, err
}

// listCreateRequest creates the List request.
func (client *ExpressRouteConnectionsClient) listCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, options *ExpressRouteConnectionsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}/expressRouteConnections"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ExpressRouteConnectionsClient) listHandleResponse(resp *http.Response) (ExpressRouteConnectionsClientListResponse, error) {
	result := ExpressRouteConnectionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExpressRouteConnectionList); err != nil {
		return ExpressRouteConnectionsClientListResponse{}, err
	}
	return result, nil
}
