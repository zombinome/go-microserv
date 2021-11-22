package microserv

import "context"

type DefaultTransportSelector struct {
}

func (ts DefaultTransportSelector) SelectService(
	ctx context.Context,
	localService HostedService,
	remoteTransports []RemoteTransport,
	endpointReqCtx EndpointRequestContext,
) (HostedService, RemoteTransport, error) {
	if localService != nil {
		return localService, nil, nil
	}

	if len(remoteTransports) > 0 {
		return nil, remoteTransports[0], nil
	}

	return nil, nil, ErrServiceNotFound
}
