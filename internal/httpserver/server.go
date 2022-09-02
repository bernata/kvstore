package httpserver

import (
	"context"
	"net"
	"net/http"
)

type Server struct {
	server   *http.Server
	listener net.Listener
}

func New(listener net.Listener) (Server, error) {
	server := &http.Server{
		Addr: listener.Addr().String(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("done"))
		}),
	}

	return Server{server: server, listener: listener}, nil
}

func (s Server) BaseURL() string {
	//goland:noinspection HttpUrlsUsage
	return "http://" + s.listener.Addr().String()
}

func (s Server) Listen() error {
	return s.server.Serve(s.listener)
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
