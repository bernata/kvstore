package httpserver_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadiness(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	response, err := http.Get(srv.BaseURL() + "/v1/readiness")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusOK, "OK")

	require.NoError(t, srv.Shutdown(context.Background()))
}
