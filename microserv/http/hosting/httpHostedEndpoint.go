package hosting

import (
	"context"
	"reflect"

	"github.com/zombinome/go-microserv/microserv"
)

type endpointHandler = func(context.Context, microserv.EndpointRequestContext, interface{}) (interface{}, error)

type httpHostedEndpoint struct {
	name         string
	requestType  reflect.Type
	responseType reflect.Type
	handler      endpointHandler
	middlewares  []EndpointMiddleware
}

func (ep httpHostedEndpoint) Name() string {
	return ep.name
}

func (ep httpHostedEndpoint) RequestType() reflect.Type {
	return ep.requestType
}

func (ep httpHostedEndpoint) ResponseType() reflect.Type {
	return ep.responseType
}

func (ep httpHostedEndpoint) HandleRequest(ctx context.Context, requestContex microserv.EndpointRequestContext, requestModel interface{}) (interface{}, error) {
	return ep.handler(ctx, requestContex, requestModel)
}

func (ep httpHostedEndpoint) AddMiddleware(newMiddleware EndpointMiddleware) httpHostedEndpoint {
	ep.middlewares = append(ep.middlewares, newMiddleware)
	return ep
}

func NewHostedEndpoint(name string, handler endpointHandler, requestType reflect.Type, responseType reflect.Type) httpHostedEndpoint {
	return httpHostedEndpoint{name, requestType, responseType, handler, []EndpointMiddleware{}}
}
