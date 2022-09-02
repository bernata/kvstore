package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		_ = srv.Shutdown(context.Background())
		logger.Info().Str("api", "notify").Str("reason", "shutting_down").Msg("")
	}()

	err = srv.Listen()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
	logger.Info().Str("api", "main").Str("reason", "shutdown").Msg("")

	time.Sleep(time.Second * 5)
}

func newServer(port int, store *kv.Store, logger *zerolog.Logger) (httpserver.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return httpserver.Server{}, err
	}

	return httpserver.New(listener, store, httpserver.LoggerOption(logger))
}

func newLogger() zerolog.Logger {
	host, _ := os.Hostname()
	return zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Str("host", host).
		Logger()
}
