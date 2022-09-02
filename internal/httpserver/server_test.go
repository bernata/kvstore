package httpserver_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/bernata/kvstore/internal/kv"

	"github.com/bernata/kvstore/internal/httpserver"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Parallel()
	srv := startTestServer()

	response, err := http.Get(srv.BaseURL())
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	require.NoError(t, srv.Shutdown(context.Background()))
}

func startTestServer() httpserver.Server {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	store := kv.NewStore(10)

	srv, err := httpserver.New(listener, store)
	if err != nil {
		panic(err)
	}

	go func() {
		err := srv.Listen()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return srv
}
