package hosting

import (
	"net/http"

	"github.com/zombinome/go-microserv/microserv"
)

func isEmpty(val string) bool {
	return len(val) == 0
}

func getOrCreateCorrelationId(request *http.Request) string {
	var correlationId = request.Header.Get(CorrelationIdHeader)
	if !isEmpty(correlationId) {
		return correlationId
	}

	return microserv.GenerateCorrelationId()
}
