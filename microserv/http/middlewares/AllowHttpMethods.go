package middlewares

import (
	"context"
	"net/http"

	"github.com/zombinome/go-microserv/microserv"
	"github.com/zombinome/go-microserv/microserv/http/hosting"
)

type nextEndpointMiddlewareInvoker = func(ctx context.Context, requestModel interface{}) (interface{}, error)

func AllowHttpMethods(args ...string) hosting.EndpointMiddleware {
	return func(
		ctx context.Context,
		requestContex microserv.EndpointRequestContext,
		requestModel interface{},
		next nextEndpointMiddlewareInvoker,
	) (interface{}, error) {
		httpRequestPtr := ctx.Value(microserv.ContextKeys.HttpRequest)
		httpRequest := httpRequestPtr.(*http.Request)

		if !containsString(args, httpRequest.Method) {
			return nil, microserv.ErrEndpointNotFound
		}

		return next(ctx, requestModel)
	}
}

func containsString(array []string, value string) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}

	return false
}
