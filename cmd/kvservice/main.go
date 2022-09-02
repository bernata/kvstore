package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/bernata/kvstore/internal/httpserver"
	"github.com/bernata/kvstore/internal/kv"
	"github.com/rs/zerolog"
)

type ServiceCommandLine struct {
	Port int `help:"Port to listen on" short:"p" default:"8282"`
}

func main() {
	var commandLine ServiceCommandLine
	_ = kong.Parse(&commandLine)

	logger := newLogger()
	store := kv.NewStore(100)
	srv, err := newServer(commandLine.Port, store, &logger)
	if err != nil {
		panic(err)
	}

	err = srv.Listen()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func newServer(port int, store *kv.Store, logger *zerolog.Logger) (httpserver.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return httpserver.Server{}, err
	}

	return httpserver.New(listener, store, httpserver.LoggerOption(logger))
}

func newLogger() zerolog.Logger {
	//Example: HOSTNAME=andrewb-host
	host := os.Getenv("HOSTNAME")
	return zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Str("host", host).
		Logger()
}
