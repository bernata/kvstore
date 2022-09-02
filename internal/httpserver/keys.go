package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

func retrieveKeyHandler(store KVStore) http.Handler {
	return newHandler(
		retrieveKeyEndpoint(store),
		decodeRetrieveValueRequest,
		encodeJSON,
	)
}

func writeKeyHandler(store KVStore) http.Handler {
	return newHandler(
		writeKeyEndpoint(store),
		decodeWriteValueRequest,
		encodeNoContent,
	)
}

func deleteKeyHandler(store KVStore) http.Handler {
	return newHandler(
		deleteKeyEndpoint(store),
		decodeDeleteKeyRequest,
		encodeNoContent,
	)
}

func retrieveKeyEndpoint(store KVStore) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestRetrieve := request.(RetrieveValueRequest)

		value, ok := store.Get(requestRetrieve.Key)
		if !ok {
			return nil, NotFoundErrorf("key '%s' not found", requestRetrieve.Key)
		}

		return RetrieveValueResponse{
			Key:   requestRetrieve.Key,
			Value: value,
		}, nil
	}
}

func writeKeyEndpoint(store KVStore) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestWrite := request.(WriteValueRequest)
		store.Write(requestWrite.Key, requestWrite.Value)
		return WriteValueResponse{}, nil
	}
}

func deleteKeyEndpoint(store KVStore) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestWrite := request.(DeleteKeyRequest)
		store.Delete(requestWrite.Key)
		return DeleteKeyResponse{}, nil
	}
}

func decodeRetrieveValueRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key, ok := vars["key"] //path parameter name
	if !ok {
		return nil, BadRequest("error reading key from path", nil)
	}

	return RetrieveValueRequest{
		Key: key,
	}, nil
}

func decodeWriteValueRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key, ok := vars["key"] //path parameter name
	if !ok {
		return nil, BadRequest("error reading key from path", nil)
	}

	var request WriteValueRequest
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, BadRequest("error reading request body", err)
	}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return nil, BadRequest("error parsing request body", err)
	}
	request.Key = key
	return request, nil
}

func decodeDeleteKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key, ok := vars["key"] //path parameter name
	if !ok {
		return nil, BadRequest("error reading key from path", nil)
	}

	return DeleteKeyRequest{
		Key: key,
	}, nil
}

type RetrieveValueRequest struct {
	Key string
}

type RetrieveValueResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WriteValueRequest struct {
	Key   string `json:"-"`
	Value string `json:"value"`
}

type WriteValueResponse struct {
}

type DeleteKeyRequest struct {
	Key string `json:"key"`
}

type DeleteKeyResponse struct {
}
