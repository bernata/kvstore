package httpserver

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func pingHandler() http.Handler {
	return newHandler(
		pingEndpoint(),
		noOpDecode,
		encodeText,
	)
}

func pingEndpoint() endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return "OK", nil
	}
}
