package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zombinome/go-microserv/microserv"
	"github.com/zombinome/go-microserv/microserv/http/hosting"
	httpTransport "github.com/zombinome/go-microserv/microserv/http/transport"
	"github.com/zombinome/go-microserv/microserv/logging"
)

type TestService struct {
	Tag string
}

const ApiServiceName = "testService"
const ApiEndpointName = "doWork"

func (svc TestService) DoWork(request TestRequestModel) (TestResponseModel, error) {
	var response = TestResponseModel{
		Service: svc.Tag,
	}

	if request.Zulu {
		response.Result = request.Bar
	} else {
		response.Result = request.Foo
	}

	return response, nil
}

func CreateService(tag string, port int) (TestService, microserv.ServiceHost, microserv.RemoteTransport) {
	var service = TestService{tag}
	var address = "127.0.0.1:" + fmt.Sprint(port)
	var hostConfig = hosting.HttpHostConfiguration{
		Address: address,
		Logger:  logging.NewConsoleLogger(tag, microserv.LogLevelInfo),
	}

	var svcTransport = fmt.Sprintf("http://127.0.0.1:%d/%s", port, ApiServiceName)
	var transport = httpTransport.NewHttpTransport(svcTransport, http.MethodPost, TestResponseType)
	transport.SetRetryPolicy(5, true, true, true)

	var host = hosting.NewHttpHost(hostConfig)
	var apiService, _ = host.AddHostedService(ApiServiceName)
	apiService.AddEndpoint(
		ApiEndpointName,
		TestRequestType,
		TestResponseType,
		func(ctx context.Context, requestContex microserv.EndpointRequestContext, request interface{}) (interface{}, error) {
			var requestModel = request.(TestRequestModel)
			return service.DoWork(requestModel)
		},
	)

	return service, host, transport
}
