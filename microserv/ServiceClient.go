package microserv

import (
	"context"
)

type ServiceClient interface {
	Invoke(ctx context.Context, service string, endpoint string, context EndpointRequestContext, request interface{}) (interface{}, error)
}

type apiServiceClient struct {
	host *ServiceHost
}

func (client apiServiceClient) Invoke(ctx context.Context, service string, endpoint string, reqContext EndpointRequestContext, request interface{}) (interface{}, error) {
	var hostedService, hasHostedSvc = (*client.host).HostedService(service)
	var remoteService, hasRemoteSvc = (*client.host).RemoteService(service)

	// If no services at all, return error
	if !hasHostedSvc && !hasRemoteSvc {
		return nil, ErrServiceNotFound
	}

	// Hosted service is present, but remote service is not
	if !hasRemoteSvc {
		if ep, isPresent := hostedService.Endpoint(endpoint); isPresent {
			return ep.HandleRequest(ctx, reqContext, request)
		}

		return nil, ErrEndpointNotFound
	}

	// Remote service is present
	var transports []RemoteTransport = remoteService.Transports()

	svc, tr, err := remoteService.TransportSelector().SelectService(ctx, hostedService, transports, reqContext)
	if err != nil {
		return nil, err
	}

	if svc != nil {
		if ep, isPresent := svc.Endpoint(endpoint); isPresent {
			return ep.HandleRequest(ctx, reqContext, request)
		}

		return nil, ErrEndpointNotFound
	}

	return tr.SendRequest(ctx, service, endpoint, reqContext.CorrelationId, &reqContext.Logger, request)
}

func NewApiClient(host *ServiceHost) ServiceClient {
	return apiServiceClient{host}
}
