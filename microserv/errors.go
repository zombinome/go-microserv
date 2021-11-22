package microserv

import "errors"

var ErrServiceNotFound = errors.New("unknown service")
var ErrEndpointNotFound = errors.New("endpoint was not found on service")

var ErrServiceAlreadyRegistered = errors.New("service with the same name is already registered in host")
var ErrEndpointAlreadyRegistered = errors.New("endpoint with the same name already registered")
