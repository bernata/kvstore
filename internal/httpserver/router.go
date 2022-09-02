package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func router() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/v1/ping", pingHandler()).Methods(http.MethodGet)
	r.Handle("/v1/readiness", readinessHandler()).Methods(http.MethodGet)
	r.Handle("/v1/keys/{key:.*}", retrieveKeyHandler()).Methods(http.MethodGet)
	r.Handle("/v1/keys/{key:.*}", writeKeyHandler()).Methods(http.MethodPost)
	r.Handle("/v1/keys/{key:.*}", deleteKeyHandler()).Methods(http.MethodDelete)

	return r
}
