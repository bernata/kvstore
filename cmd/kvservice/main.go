package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/bernata/kvstore/internal/httpserver"
)

type ServiceCommandLine struct {
	Port int `help:"Port to listen on" short:"p" default:"8282"`
}

func main() {
	var commandLine ServiceCommandLine
	_ = kong.Parse(&commandLine)

	srv, err := newServer(commandLine.Port)
	if err != nil {
		panic(err)
	}

	err = srv.Listen()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func newServer(port int) (httpserver.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return httpserver.Server{}, err
	}

	return httpserver.New(listener)
}
