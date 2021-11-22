package main

import (
	"context"

	"github.com/zombinome/go-microserv/microserv"
)

type TransportSelector struct {
}

func (t TransportSelector) SelectService(
	ctx context.Context,
	localService microserv.HostedService,
	remoteTransports []microserv.RemoteTransport,
	endpointReqCtx microserv.EndpointRequestContext,
) (microserv.HostedService, microserv.RemoteTransport, error) {
	var correlationId = endpointReqCtx.CorrelationId
	var divider = byte(len(remoteTransports) + 1)
	var num = correlationId[4] % divider
	if num == 0 {
		endpointReqCtx.Logger.InfoF("Selected service %s", localService.Name())
		return localService, nil, nil
	}

	var transport = remoteTransports[num-1]
	endpointReqCtx.Logger.InfoF("Selected remote service (transport %s)", transport.Id())
	return nil, transport, nil
}
