package httpserver_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/bernata/kvstore/apiclient"
	"github.com/stretchr/testify/require"
)

func TestClientRetrieveKeyNotFound(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	client := apiclient.NewClient(apiclient.BaseURL(srv.BaseURL()))
	_, err := client.RetrieveValue(context.Background(), &apiclient.RetrieveValueRequest{Key: "k1"})
	requireNotFound(t, err)

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestClientRetrieveKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	client := apiclient.NewClient(apiclient.BaseURL(srv.BaseURL()))
	_, err := client.WriteValue(context.Background(), &apiclient.WriteValueRequest{Key: "k1", Value: "v1"})
	require.NoError(t, err)

	response, err := client.RetrieveValue(context.Background(), &apiclient.RetrieveValueRequest{Key: "k1"})
	require.NoError(t, err)
	require.Equal(t, apiclient.RetrieveValueResponse{Key: "k1", Value: "v1"}, response)

	require.NoError(t, srv.Shutdown(context.Background()))
}

func TestClientDeleteKey(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	client := apiclient.NewClient(apiclient.BaseURL(srv.BaseURL()))
	_, err := client.WriteValue(context.Background(), &apiclient.WriteValueRequest{Key: "k1", Value: "v1"})
	require.NoError(t, err)
	_, err = client.DeleteKey(context.Background(), &apiclient.DeleteKeyRequest{Key: "k1"})
	require.NoError(t, err)

	_, err = client.RetrieveValue(context.Background(), &apiclient.RetrieveValueRequest{Key: "k1"})
	requireNotFound(t, err)

	require.NoError(t, srv.Shutdown(context.Background()))
}

func requireNotFound(t *testing.T, err error) {
	var kvError *apiclient.KeyValueError
	require.True(t, errors.As(err, &kvError))
	require.Equal(t, http.StatusNotFound, kvError.StatusCodeHTTP)
}
