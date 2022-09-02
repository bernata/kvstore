package httpserver

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func router(store KVStore, serverOptions *serverOptions) *mux.Router {
	r := mux.NewRouter()
	r.Use(loggerContext(serverOptions.logger))

	r.Handle("/v1/ping", pingHandler()).Methods(http.MethodGet)
	r.Handle("/v1/readiness", readinessHandler()).Methods(http.MethodGet)

	subrouter := r.PathPrefix("/v1").Subrouter()
	subrouter.Use(requestLogger)
	subrouter.Handle("/keys/{key:.*}", retrieveKeyHandler(store)).Methods(http.MethodGet)
	subrouter.Handle("/keys/{key:.*}", writeKeyHandler(store)).Methods(http.MethodPost)
	subrouter.Handle("/keys/{key:.*}", deleteKeyHandler(store)).Methods(http.MethodDelete)

	return r
}

func loggerContext(logger *zerolog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := logger.WithContext(request.Context())
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())
		ctx := logger.WithContext(r.Context())
		wrapped := &responseWriter{ResponseWriter: w}
		defer func(begin time.Time) {
			elapsed := time.Since(begin)
			logger.Info().
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Str("remote_addr", r.RemoteAddr).
				Int64("content_length", r.ContentLength).
				Int("status_code", wrapped.statusCode).
				Int64("elapsed_ms", elapsed.Milliseconds()).
				Msg("")
		}(time.Now())

		next.ServeHTTP(wrapped, r.WithContext(ctx))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
