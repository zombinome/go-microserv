package hosting

import (
	"reflect"
	"strings"

	"github.com/zombinome/go-microserv/microserv"
)

type httpHostedService struct {
	name        string
	middlewares []EndpointMiddleware
	endpoints   map[string]httpHostedEndpoint
}

func (svc httpHostedService) Name() string {
	return svc.name
}

func (svc httpHostedService) Endpoints() []microserv.HostedEndpoint {
	var result = make([]microserv.HostedEndpoint, len(svc.endpoints))
	var i = 0
	for _, ep := range svc.endpoints {
		result[i] = ep
		i++
	}

	return result
}

func (svc httpHostedService) Endpoint(name string) (microserv.HostedEndpoint, bool) {
	var key = strings.ToUpper(name)
	ep, isPresent := svc.endpoints[key]
	return ep, isPresent
}

func (svc httpHostedService) AddEndpoint(name string, requestType, responseType reflect.Type, handler microserv.EndpointHandler) (microserv.HostedEndpoint, error) {
	var key = strings.ToUpper(name)
	if _, isPresent := svc.endpoints[key]; isPresent {
		return nil, microserv.ErrEndpointAlreadyRegistered
	}

	var newEndpoint = httpHostedEndpoint{
		name:         name,
		requestType:  requestType,
		responseType: responseType,
		handler:      handler,
		middlewares:  []EndpointMiddleware{},
	}

	svc.endpoints[key] = newEndpoint

	return newEndpoint, nil
}

func (svc httpHostedService) AddMiddleware(newMiddleware EndpointMiddleware) HttpHostedService {
	svc.middlewares = append(svc.middlewares, newMiddleware)
	return svc
}
