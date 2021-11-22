package microserv

type ContextKey string

type contextKeysStruct struct {
	CorrelationId   ContextKey
	EndpointContext ContextKey
	HttpRequest     ContextKey
}

var ContextKeys = contextKeysStruct{
	CorrelationId:   "CorrelationId",
	EndpointContext: "EndpointRequestContext",
	HttpRequest:     "HttpRequest",
}
