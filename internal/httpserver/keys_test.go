package httpserver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/bernata/kvstore/internal/httpserver"
	"github.com/stretchr/testify/require"
)

func TestRetrieveKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	response, err := http.Get(srv.BaseURL() + "/v1/keys/foo")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusOK, "{\"key\":\"foo\",\"value\":\"base64 me\"}")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestRetrieveKeyWithSlash(t *testing.T) {
	t.Parallel()
	srv := startTestServer()
	response, err := http.Get(srv.BaseURL() + "/v1/keys/foo/bar/baz")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusOK, "{\"key\":\"foo/bar/baz\",\"value\":\"base64 me\"}")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestWriteKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	response, err := http.Post(srv.BaseURL()+"/v1/keys/foo", applicationJSON, reader(httpserver.WriteValueRequest{Key: "k1", Value: "v1"}))
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNoContent, "")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	response, err := http.Post(srv.BaseURL()+"/v1/keys/foo", applicationJSON, reader(httpserver.DeleteKeyRequest{Key: "k1"}))
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNoContent, "")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func reader(data any) io.Reader {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

const applicationJSON = "application/json; charset=utf-8"
