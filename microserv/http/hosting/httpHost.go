package hosting

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zombinome/go-microserv/microserv"
	"github.com/zombinome/go-microserv/microserv/serialization"
)

type httpHost struct {
	server         *http.Server
	hostedServices map[string]httpHostedService
	remoteServices map[string]microserv.RemoteService
	logger         microserv.Logger
	client         microserv.ServiceClient
	middlewares    []HostMiddleware
}

func (host httpHost) AddHostedService(name string) (microserv.HostedService, error) {
	var key = strings.ToUpper(name)
	if _, isPreset := host.hostedServices[key]; isPreset {
		return nil, microserv.ErrServiceAlreadyRegistered
	}

	var newService = httpHostedService{
		name,
		[]EndpointMiddleware{},
		make(map[string]httpHostedEndpoint),
	}
	host.hostedServices[key] = newService

	return newService, nil
}

func (host httpHost) AddRemoteService(name string, transport *microserv.RemoteTransport) (microserv.RemoteService, error) {
	var key = strings.ToUpper(name)
	if _, isPreset := host.hostedServices[key]; isPreset {
		return nil, microserv.ErrServiceAlreadyRegistered
	}

	var newService = microserv.NewRemoteService(&host, name, transport)
	host.remoteServices[key] = newService
	return newService, nil
}

func (host httpHost) AddRemoteServiceWithSelector(
	name string,
	selector microserv.TransportSelector,
	transports ...*microserv.RemoteTransport,
) (microserv.RemoteService, error) {
	var key = strings.ToUpper(name)
	if _, isPreset := host.remoteServices[key]; isPreset {
		return nil, microserv.ErrServiceAlreadyRegistered
	}

	var newService = microserv.NewRemoteServiceWithSelector(&host, name, selector, transports...)
	host.remoteServices[key] = newService
	return newService, nil
}

func (host httpHost) HostedServices() []microserv.HostedService {
	var result = make([]microserv.HostedService, len(host.hostedServices))
	var i = 0
	for _, service := range host.hostedServices {
		result[i] = service
		i++
	}

	return result
}

func (host httpHost) RemoteSerices() []microserv.RemoteService {
	var result = make([]microserv.RemoteService, len(host.remoteServices))
	var i = 0
	for _, service := range host.remoteServices {
		result[i] = service
		i++
	}

	return result
}

func (host httpHost) HostedService(name string) (microserv.HostedService, bool) {
	var key = strings.ToUpper(name)
	service, isPresent := host.hostedServices[key]
	return service, isPresent
}

func (host httpHost) RemoteService(name string) (microserv.RemoteService, bool) {
	var key = strings.ToUpper(name)
	service, isPresent := host.remoteServices[key]
	return service, isPresent
}

func (host httpHost) ServiceClient() microserv.ServiceClient {
	return host.client
}

func (host httpHost) Start() error {
	return host.server.ListenAndServe()
}

func (host httpHost) Stop(ctx context.Context) error {
	return host.server.Shutdown(ctx)
}

func (host *httpHost) AddMiddleware(newMiddleware HostMiddleware) {
	host.middlewares = append(host.middlewares, newMiddleware)
}

func NewHttpHost(config HttpHostConfiguration) *httpHost {
	var host = httpHost{
		hostedServices: make(map[string]httpHostedService),
		remoteServices: make(map[string]microserv.RemoteService),
		logger:         config.Logger,
		middlewares:    []HostMiddleware{},
	}

	host.server = &http.Server{
		Addr: config.Address,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
			host.handleHttpRequest(rw, request)
		}),
	}

	var apiHost microserv.ServiceHost = host
	host.client = microserv.NewApiClient(&apiHost)

	return &host
}

func (host *httpHost) handleHttpRequest(rw http.ResponseWriter, request *http.Request) {
	var reqCtx = request.Context()
	reqCtx = context.WithValue(reqCtx, microserv.ContextKeys.HttpRequest, request)
	// Generating correlation ID which can be used to track single client request
	// across different microservices
	var correlationId = getOrCreateCorrelationId(request)
	reqCtx = context.WithValue(reqCtx, microserv.ContextKeys.CorrelationId, correlationId)

	// Creating logger scope for this request
	// Setting request handling time measurement
	var logger = host.logger.CreateScope("EP", correlationId)
	logger.BeginMeasure("Begin handling request " + request.RequestURI)
	defer logger.EndMeasure("Finished handling request " + request.RequestURI)

	// Parsing service and endpoint name and checking that they are belong to the known endpoint
	var tokens = strings.Split(strings.TrimPrefix(request.URL.Path, "/"), "/")
	if len(tokens) != 2 {
		http.NotFound(rw, request)
		return
	}

	var serviceName = strings.ToUpper(tokens[0])
	service, isPresent := host.hostedServices[serviceName]
	if !isPresent {
		http.NotFound(rw, request)
		return
	}

	var endpointName = strings.ToUpper(tokens[1])
	endpoint, isPresent := service.endpoints[endpointName]
	if !isPresent {
		http.NotFound(rw, request)
		return
	}

	var unhandledError = host.callHostMiddlewares(reqCtx, rw, request, logger, &service, &endpoint)
	if unhandledError != nil {
		logger.Error(unhandledError.Error())
		http.Error(rw, unhandledError.Error(), http.StatusInternalServerError)
	}
}

func (host *httpHost) callHostMiddlewares(
	ctx context.Context,
	rw http.ResponseWriter,
	request *http.Request,
	logger microserv.LoggerScope,
	service *httpHostedService,
	endpoint *httpHostedEndpoint,
) error {
	var length = len(host.middlewares)
	if length == 0 {
		return host.handleServiceRequest(ctx, rw, request, logger, service, endpoint)
	}

	var i = 0
	var runner nextServiceMiddlewareInvoker = nil
	runner = func(ctx context.Context) error {
		if i < length {
			var middleware = host.middlewares[i]
			i++
			return middleware(ctx, rw, request, logger, runner)
		}

		return host.handleServiceRequest(ctx, rw, request, logger, service, endpoint)
	}

	return runner(ctx)
}

func (host *httpHost) handleServiceRequest(
	ctx context.Context,
	rw http.ResponseWriter,
	request *http.Request,
	logger microserv.LoggerScope,
	service *httpHostedService,
	endpoint *httpHostedEndpoint,
) error {
	var correlationId = ctx.Value(microserv.ContextKeys.CorrelationId).(string)
	var reqContext = microserv.NewEndpointRequestContext(*host, logger, service.name, endpoint.name, correlationId)
	var requestModel, err = serialization.ReadRequest(endpoint.requestType, request)
	if _, isDict := requestModel.(map[string]interface{}); isDict {
		logger.Error("Serialization problem")
	}
	if err != nil {
		return err
	}

	var response interface{}
	var length = len(service.middlewares) + len(endpoint.middlewares)
	if length == 0 {
		response, err = endpoint.HandleRequest(ctx, reqContext, requestModel)
	} else {
		var allMiddlewares []EndpointMiddleware = []EndpointMiddleware{}
		allMiddlewares = append(append(allMiddlewares[:0], service.middlewares...), service.middlewares...)

		var i = 0
		var runner nextEndpointMiddlewareInvoker = nil
		runner = func(ctx context.Context, reqModel interface{}) (interface{}, error) {
			if i < length {
				var middleware = allMiddlewares[i]
				i++
				return middleware(ctx, reqContext, reqModel, runner)
			}

			return endpoint.HandleRequest(ctx, reqContext, requestModel)
		}

		response, err = runner(ctx, requestModel)
	}

	if err != nil {
		return err
	}

	if response != nil {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(rw).Encode(response)
		if err != nil {
			return err
		}
	}

	return nil
}
