package httpserver_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/bernata/kvstore/internal/httpserver"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Parallel()
	buffer := bytes.NewBuffer(nil)
	logger := zerolog.New(buffer)
	srv := startTestServer(httpserver.LoggerOption(&logger))

	response, err := http.Get(srv.BaseURL() + "/v1/keys/foo")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode)
	require.NotEmpty(t, buffer.String())
	require.Contains(t, buffer.String(), "{\"level\":\"info\",\"path\":\"/v1/keys/foo\",\"method\":\"GET\"")
	require.Contains(t, buffer.String(), "\"status_code\":404")

	require.NoError(t, srv.Shutdown(context.Background()))
}
