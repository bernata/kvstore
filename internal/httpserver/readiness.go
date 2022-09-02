package httpserver

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func readinessHandler() http.Handler {
	return newHandler(
		readinessEndpoint(),
		noOpDecode,
		encodeText,
	)
}

func readinessEndpoint() endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return "OK", nil
	}
}
