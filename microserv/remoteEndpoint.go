package microserv

import (
	"context"
	"reflect"
)

type remoteEndpoint struct {
	name         string
	service      *remoteService
	requestType  reflect.Type
	responseType reflect.Type
}

func (ep remoteEndpoint) Name() string {
	return ep.name
}

func (ep remoteEndpoint) RequestType() reflect.Type {
	return ep.requestType
}

func (ep remoteEndpoint) ResponseType() reflect.Type {
	return ep.responseType
}

func (ep remoteEndpoint) Invoke(ctx context.Context, requestContext EndpointRequestContext, requestModel interface{}) (interface{}, error) {
	var host = *(ep.service.host)
	var hostedService, _ = host.HostedService(ep.service.name)
	var svc, tr, err = ep.service.transportSelector.SelectService(ctx, hostedService, ep.service.Transports(), requestContext)
	if err != nil {
		return nil, err
	}

	if svc != nil {
		var endpoint, _ = svc.Endpoint(ep.name)
		return endpoint.HandleRequest(ctx, requestContext, requestModel)
	}

	return tr.SendRequest(ctx, ep.service.name, ep.name, requestContext.CorrelationId, &requestContext.Logger, requestContext)
}
