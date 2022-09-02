package httpserver

import (
	"context"
	"net"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Server struct {
	server   *http.Server
	listener net.Listener
}

func New(listener net.Listener) (Server, error) {
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: router(),
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

func newHandler(
	e endpoint.Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc,
) *httptransport.Server {
	return httptransport.NewServer(e, dec, enc)
}
