package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"

	"github.com/zombinome/go-microserv/microserv"
	"github.com/zombinome/go-microserv/microserv/http/hosting"
	"github.com/zombinome/go-microserv/microserv/serialization"
)

type httpTransport struct {
	baseUrl      string
	method       string
	contentType  string
	client       *http.Client
	responseType reflect.Type
	retryPolicy  *httpTransportConfig
}

type httpTransportConfig struct {
	RetryCount    uint
	RetryOnPost   bool
	RetryOnPut    bool
	RetryOnDelete bool
}

var DefaultConfig = httpTransportConfig{
	RetryCount:  5,
	RetryOnPost: true,
	RetryOnPut:  true,
}

func (transport httpTransport) Id() string {
	return transport.baseUrl
}

func (t httpTransport) SendRequest(
	ctx context.Context,
	serviceName string,
	endpointName string,
	correlationId string,
	logger *microserv.LoggerScope,
	request interface{},
) (interface{}, error) {

	httpRequest, err := t.toHttpRequest(ctx, endpointName, request, correlationId)
	if err != nil {
		return nil, err
	}

	var httpResponse *http.Response
	var reqNumber uint = 0
	for ; ; reqNumber++ {
		httpResponse, err = http.DefaultClient.Do(httpRequest)
		if err == nil || !t.shouldRetry(err, reqNumber) {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if t.responseType == nil {
		return nil, nil
	}

	// Deserializing response here
	return serialization.DeserializeFromJsonReader(httpResponse.Body, t.responseType)
}

// Sets retry policy for current trasport
// retryCount - count of retries, use 0 to disable retry policy
func (t *httpTransport) SetRetryPolicy(retryCount uint, retryOnPost bool, retryOnPut bool, retryOnDelete bool) *httpTransport {
	t.retryPolicy.RetryCount = retryCount
	t.retryPolicy.RetryOnPost = retryOnPost
	t.retryPolicy.RetryOnPut = retryOnPut
	t.retryPolicy.RetryOnDelete = retryOnDelete
	return t
}

func NewHttpTransport(baseUrl string, method string, responseType reflect.Type) httpTransport {
	var baseUrlNorm = baseUrl
	var lastChar = baseUrl[len(baseUrl)-1]
	if lastChar != '/' {
		baseUrlNorm = baseUrl + "/"
	}

	var client = http.DefaultClient
	return httpTransport{
		baseUrlNorm,
		method,
		"application/json; charset=utf-8",
		client,
		responseType,
		&httpTransportConfig{0, false, true, false},
	}
}

func (t *httpTransport) toHttpRequest(ctx context.Context, endpointName string, model interface{}, correlationId string) (*http.Request, error) {
	query, err := serialization.WriteToQueryParameters(model)
	if err != nil {
		return nil, err
	}

	var fullAddress = t.baseUrl + endpointName + "?" + query.Encode()
	body, err := serialization.SerializeToJsonReader(model)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, t.method, fullAddress, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", t.contentType)
	request.Header.Set(hosting.CorrelationIdHeader, correlationId)

	request.Close = true

	return request, nil
}

func (t *httpTransport) shouldRetry(err error, tryNumber uint) bool {
	var rp = t.retryPolicy
	if rp.RetryCount == 0 || tryNumber >= rp.RetryCount || !errors.Is(err, io.EOF) {
		return false
	}

	if t.method == http.MethodPost {
		return rp.RetryOnPost
	}

	if t.method == http.MethodPut {
		return rp.RetryOnPut
	}

	if t.method == http.MethodDelete {
		return rp.RetryOnDelete
	}

	return true
}
