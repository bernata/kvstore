package httpserver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/bernata/kvstore/apiclient"

	"github.com/stretchr/testify/require"
)

func TestRetrieveKeyNotFound(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	response, err := http.Get(srv.BaseURL() + "/v1/keys/foo")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNotFound, "{\"message\":\"[404]: key 'foo' not found\",\"status_code\":404}")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestRetrieveKeyWithSlash(t *testing.T) {
	t.Parallel()
	srv := startTestServer()
	response, err := http.Get(srv.BaseURL() + "/v1/keys/foo/bar/baz")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNotFound, "{\"message\":\"[404]: key 'foo/bar/baz' not found\",\"status_code\":404}")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestWriteKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	// NOTE: k1 is ignored because the key is determined from the path
	response, err := http.Post(srv.BaseURL()+"/v1/keys/foo/bar", applicationJSON, reader(apiclient.WriteValueRequest{Key: "k1", Value: "v1"}))
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNoContent, "")

	response, err = http.Get(srv.BaseURL() + "/v1/keys/foo/bar")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusOK, "{\"key\":\"foo/bar\",\"value\":\"v1\"}")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestWriteDeleteGetKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	// NOTE: k1 is ignored because the key is determined from the path
	response, err := http.Post(srv.BaseURL()+"/v1/keys/foo", applicationJSON, reader(apiclient.WriteValueRequest{Value: "v1"}))
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNoContent, "")

	req, err := http.NewRequest(http.MethodDelete, srv.BaseURL()+"/v1/keys/foo", nil)
	require.NoError(t, err)
	response, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNoContent, "")

	response, err = http.Get(srv.BaseURL() + "/v1/keys/foo")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusNotFound, "{\"message\":\"[404]: key 'foo' not found\",\"status_code\":404}")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	req, err := http.NewRequest(http.MethodDelete, srv.BaseURL()+"/v1/keys/foo", nil)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(req)
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
