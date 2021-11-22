package microserv

import (
	"context"
	"reflect"
)

type EndpointHandler = func(ctx context.Context, requestContex EndpointRequestContext, request interface{}) (interface{}, error)

type ServiceHost interface {
	AddHostedService(name string) (HostedService, error)

	AddRemoteService(name string, transport *RemoteTransport) (RemoteService, error)

	AddRemoteServiceWithSelector(name string, selector TransportSelector, transports ...*RemoteTransport) (RemoteService, error)

	HostedServices() []HostedService

	RemoteSerices() []RemoteService

	HostedService(name string) (HostedService, bool)

	RemoteService(name string) (RemoteService, bool)

	ServiceClient() ServiceClient

	Start() error

	Stop(ctx context.Context) error
}

type HostedService interface {
	Name() string

	Endpoints() []HostedEndpoint

	Endpoint(name string) (HostedEndpoint, bool)

	AddEndpoint(name string, requestType, responseType reflect.Type, handler EndpointHandler) (HostedEndpoint, error)
}

type HostedEndpoint interface {
	Name() string

	RequestType() reflect.Type

	ResponseType() reflect.Type

	HandleRequest(ctx context.Context, requestContext EndpointRequestContext, requestModel interface{}) (interface{}, error)
}

type RemoteService interface {
	Name() string

	Transport() RemoteTransport

	Transports() []RemoteTransport

	TransportSelector() TransportSelector

	AddEndpoint(name string, requestType reflect.Type, responseType reflect.Type) (RemoteEndpoint, error)

	Endpoints() []RemoteEndpoint

	Endpoint(name string) (RemoteEndpoint, bool)
}

type RemoteEndpoint interface {
	Name() string

	RequestType() reflect.Type

	ResponseType() reflect.Type

	Invoke(ctx context.Context, requestContext EndpointRequestContext, requestModel interface{}) (interface{}, error)
}

type RemoteTransport interface {
	Id() string

	// Sends request to external service
	SendRequest(
		ctx context.Context,
		serviceName string,
		endpointName string,
		correlationId string,
		logger *LoggerScope,
		request interface{},
	) (interface{}, error)
}

type TransportSelector interface {
	// Sends request to appropriate service
	// Returns service result
	SelectService(
		ctx context.Context,
		localService HostedService,
		remoteTransports []RemoteTransport,
		endpointReqCtx EndpointRequestContext,
	) (HostedService, RemoteTransport, error)
}
