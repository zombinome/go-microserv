package main

import (
	"context"
	"fmt"

	"github.com/zombinome/go-microserv/microserv"
	httphosting "github.com/zombinome/go-microserv/microserv/http/hosting"
	"github.com/zombinome/go-microserv/microserv/logging"
)

var _, host1, transport1 = CreateService("Remote Service 1", 8081)
var _, host2, transport2 = CreateService("Remote Service 2", 8082)
var _, host3, transport3 = CreateService("Remote Service 3", 8083)
var remoteHosts = []microserv.ServiceHost{host1, host2, host3}

func main() {
	// Starting remote hosts
	for _, remoteHost := range remoteHosts {
		go func(rh microserv.ServiceHost) {
			rh.Start()
		}(remoteHost)
	}

	var host = initHostedAndService()

	// Setting up local service
	fmt.Println("System start")
	err := host.Start()
	if err != nil {
		fmt.Print(err)
	}
}

func initHostedAndService() microserv.ServiceHost {
	var config = httphosting.HttpHostConfiguration{
		Address: "127.0.0.1:8080",
		Logger:  logging.NewConsoleLogger("MainHost", microserv.LogLevelInfo),
	}

	var host = httphosting.NewHttpHost(config)

	var internalServiceBase, _ = host.AddHostedService("internal")
	var internalService = internalServiceBase.(httphosting.HttpHostedService)
	internalService.AddEndpoint("stop", nil, nil, func(ctx context.Context, reqContext microserv.EndpointRequestContext, request interface{}) (interface{}, error) {
		reqContext.Logger.Info("Stopping...")
		var err = host.Stop(ctx)
		for _, rh := range remoteHosts {
			rh.Stop(ctx)
		}
		reqContext.Logger.Info("Stopped!")
		return nil, err
	})

	var testService = TestService{"Hosted Service"}
	var hostedService, _ = host.AddHostedService(ApiServiceName)
	hostedService.AddEndpoint(ApiEndpointName, TestRequestType, TestResponseType, func(ctx context.Context, requestContex microserv.EndpointRequestContext, request interface{}) (interface{}, error) {
		var reqModel = request.(TestRequestModel)
		return testService.DoWork(reqModel)
	})

	hostedService.AddEndpoint("testSelector", TestRequestType, TestResponseType, func(ctx context.Context, requestContex microserv.EndpointRequestContext, request interface{}) (interface{}, error) {
		return requestContex.Client.Invoke(ctx, ApiServiceName, ApiEndpointName, requestContex, request)
	})

	host.AddRemoteServiceWithSelector(ApiServiceName, TransportSelector{}, &transport1, &transport2, &transport3)
	return host
}
