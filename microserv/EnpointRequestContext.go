package microserv

type EndpointRequestContext struct {
	CorrelationId string
	Client        ServiceClient
	Logger        LoggerScope
}

func NewEndpointRequestContext(
	host ServiceHost,
	logger LoggerScope,
	serviceName string,
	actionName string,
	correlationId string) EndpointRequestContext {
	return EndpointRequestContext{
		Logger:        logger,
		Client:        host.ServiceClient(),
		CorrelationId: correlationId,
	}
}
