package httpserver_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/bernata/kvstore/internal/httpserver"
	"github.com/stretchr/testify/require"
)

func TestServerRecover(t *testing.T) {
	t.Parallel()
	srv := startPanicServer()

	response, err := http.Get(srv.BaseURL() + "/v1/keys/foo")
	require.NoError(t, err)
	requireResponse(t, response, http.StatusInternalServerError, "{\"message\":\"[500]: server terminated abnormally: unit test panic\",\"status_code\":500}")

	require.NoError(t, srv.Shutdown(context.Background()))
}

func startPanicServer() httpserver.Server {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	store := kvPanic{}
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

type kvPanic struct{}

func (k kvPanic) Get(_ string) (string, bool) {
	panic("unit test panic")
}

func (k kvPanic) Delete(_ string) {
	panic("unit test panic")
}

func (k kvPanic) Write(_, _ string) {
	panic("unit test panic")
}
