package hosting

import (
	"context"
	"net/http"

	"github.com/zombinome/go-microserv/microserv"
)

type HttpHostedService interface {
	microserv.HostedService

	AddMiddleware(newMiddleware EndpointMiddleware) HttpHostedService
}

type HttpHostedEndpoint interface {
	microserv.HostedEndpoint

	AddMiddleware(newMiddleware EndpointMiddleware) httpHostedEndpoint
}

type nextServiceMiddlewareInvoker = func(ctx context.Context) error
type HostMiddleware = func(
	ctx context.Context,
	rw http.ResponseWriter,
	httpRequest *http.Request,
	logger microserv.LoggerScope,
	next nextServiceMiddlewareInvoker,
) error

type nextEndpointMiddlewareInvoker = func(ctx context.Context, requestModel interface{}) (interface{}, error)
type EndpointMiddleware = func(
	ctx context.Context,
	requestContex microserv.EndpointRequestContext,
	requestModel interface{},
	next nextEndpointMiddlewareInvoker,
) (interface{}, error)

// func callServiceMiddlewareChain(
// 	ctx context.Context,
// 	host *httpHost,
// 	rw http.ResponseWriter,
// 	req *http.Request,
// 	service httpHostedService,
// 	endpointName string,
// 	logger microserv.LoggerScope,
// ) error {
// 	var length = len(host.middlewares)
// 	if length == 0 {
// 		return host.handleRequestEndpointStage(ctx, rw, req, service, endpointName, logger)
// 	}

// 	var i = 0
// 	var runner nextServiceMiddlewareInvoker = nil
// 	runner = func(ctx context.Context) error {
// 		if i < length {
// 			var middleware = service.Middleware[i]
// 			i++
// 			return middleware(ctx, rw, req, logger, runner)
// 		}

// 		return host.handleRequestEndpointStage(ctx, rw, req, service, endpointName, logger)
// 	}

// 	return runner(ctx)
// }

// func callEndpointMiddlewareChain(
// 	ctx context.Context,
// 	host *httpHost,
// 	service httpHostedService,
// 	endpoint httpHostedEndpoint,
// 	reqContext microserv.EndpointRequestContext,
// 	requestModel interface{},
// ) (interface{}, error) {
// 	var length = len(service.Middleware)
// 	if length == 0 {
// 		return endpoint.HandleRequest(ctx, reqContext, requestModel)
// 	}

// 	var i = 0
// 	var runner nextEndpointMiddlewareInvoker = nil
// 	runner = func(ctx context.Context, model interface{}) (interface{}, error) {
// 		if i < length {
// 			var middleware = endpoint.Middleware[i]
// 			i++

// 			return middleware(ctx, reqContext, model, runner)
// 		}

// 		return endpoint.HandleRequest(ctx, reqContext, model)
// 	}

// 	return runner(ctx, requestModel)
// }
