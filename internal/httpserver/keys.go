package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/bernata/kvstore/apiclient"

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
		requestRetrieve := request.(apiclient.RetrieveValueRequest)

		key, err := url.PathUnescape(requestRetrieve.Key)
		if err != nil {
			return nil, err
		}

		value, ok := store.Get(key)
		if !ok {
			return nil, apiclient.NotFoundErrorf("key '%s' not found", key)
		}

		return apiclient.RetrieveValueResponse{
			Key:   key,
			Value: value,
		}, nil
	}
}

func writeKeyEndpoint(store KVStore) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestWrite := request.(apiclient.WriteValueRequest)
		key, err := url.PathUnescape(requestWrite.Key)
		if err != nil {
			return nil, err
		}
		store.Write(key, requestWrite.Value)
		return apiclient.WriteValueResponse{}, nil
	}
}

func deleteKeyEndpoint(store KVStore) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestDelete := request.(apiclient.DeleteKeyRequest)
		key, err := url.PathUnescape(requestDelete.Key)
		if err != nil {
			return nil, err
		}
		store.Delete(key)
		return apiclient.DeleteKeyResponse{}, nil
	}
}

func decodeRetrieveValueRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key, ok := vars["key"] //path parameter name
	if !ok {
		return nil, apiclient.BadRequest("error reading key from path", nil)
	}

	return apiclient.RetrieveValueRequest{
		Key: key,
	}, nil
}

func decodeWriteValueRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key, ok := vars["key"] //path parameter name
	if !ok {
		return nil, apiclient.BadRequest("error reading key from path", nil)
	}

	var request apiclient.WriteValueRequest
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, apiclient.BadRequest("error reading request body", err)
	}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return nil, apiclient.BadRequest("error parsing request body", err)
	}
	request.Key = key
	return request, nil
}

func decodeDeleteKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key, ok := vars["key"] //path parameter name
	if !ok {
		return nil, apiclient.BadRequest("error reading key from path", nil)
	}

	return apiclient.DeleteKeyRequest{
		Key: key,
	}, nil
}
