package httpserver_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	response, err := http.Get(srv.BaseURL() + "/v1/ping")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusOK, "OK")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func requireResponse(t *testing.T, response *http.Response, statusCode int, body string) {
	require.Equal(t, statusCode, response.StatusCode)
	data, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, body, string(data))
}
