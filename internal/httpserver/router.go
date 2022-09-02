package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func router(store KVStore) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/v1/ping", pingHandler()).Methods(http.MethodGet)
	r.Handle("/v1/readiness", readinessHandler()).Methods(http.MethodGet)
	r.Handle("/v1/keys/{key:.*}", retrieveKeyHandler(store)).Methods(http.MethodGet)
	r.Handle("/v1/keys/{key:.*}", writeKeyHandler(store)).Methods(http.MethodPost)
	r.Handle("/v1/keys/{key:.*}", deleteKeyHandler(store)).Methods(http.MethodDelete)

	return r
}
