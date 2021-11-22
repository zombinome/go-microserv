package microserv

import (
	"reflect"
	"strings"
)

type remoteService struct {
	name              string
	endpoints         map[string]remoteEndpoint
	transportSelector TransportSelector
	transports        map[string]RemoteTransport
	host              *ServiceHost
}

// type remoteServiceInternal struct {
// 	transportSelector TransportSelector
// 	transports        []RemoteTransport
// }

func (svc remoteService) Name() string {
	return svc.name
}

func (svc remoteService) Transport() RemoteTransport {
	for _, transport := range svc.transports {
		return transport
	}

	return nil
}

func (svc remoteService) Transports() []RemoteTransport {
	var result = make([]RemoteTransport, len(svc.transports))
	var i = 0
	for _, transport := range svc.transports {
		result[i] = transport
		i++
	}

	return result
}

func (svc remoteService) TransportSelector() TransportSelector {
	return svc.transportSelector
}

func (svc remoteService) AddEndpoint(name string, requestType reflect.Type, responseType reflect.Type) (RemoteEndpoint, error) {
	var key = strings.ToUpper(name)
	if _, isPresent := svc.endpoints[key]; isPresent {
		return nil, ErrEndpointAlreadyRegistered
	}

	var newEndpoint = remoteEndpoint{name, &svc, requestType, responseType}
	svc.endpoints[key] = newEndpoint

	return newEndpoint, nil
}

func (svc remoteService) Endpoints() []RemoteEndpoint {
	var result = make([]RemoteEndpoint, len(svc.endpoints))
	var i = 0
	for _, value := range svc.endpoints {
		result[i] = value
		i++
	}

	return result
}

func (svc remoteService) Endpoint(name string) (RemoteEndpoint, bool) {
	var key = strings.ToUpper(name)
	ep, isPresent := svc.endpoints[key]
	return ep, isPresent
}

func NewRemoteService(host ServiceHost, name string, transport *RemoteTransport) RemoteService {
	return NewRemoteServiceWithSelector(host, name, DefaultTransportSelector{}, transport)
}

func NewRemoteServiceWithSelector(host ServiceHost, name string, selector TransportSelector, transports ...*RemoteTransport) RemoteService {
	allTransports := make(map[string]RemoteTransport, len(transports))
	for _, transport := range transports {
		allTransports[(*transport).Id()] = *transport
	}

	return remoteService{
		name:              name,
		endpoints:         make(map[string]remoteEndpoint),
		transportSelector: selector,
		transports:        allTransports,
		host:              &host,
	}
}
