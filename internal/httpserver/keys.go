package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

func retrieveKeyHandler() http.Handler {
	return newHandler(
		retrieveKeyEndpoint(),
		decodeRetrieveValueRequest,
		encodeJSON,
	)
}

func writeKeyHandler() http.Handler {
	return newHandler(
		writeKeyEndpoint(),
		decodeWriteValueRequest,
		encodeNoContent,
	)
}

func deleteKeyHandler() http.Handler {
	return newHandler(
		deleteKeyEndpoint(),
		decodeDeleteValueRequest,
		encodeNoContent,
	)
}

func retrieveKeyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		requestRetrieve := request.(RetrieveValueRequest)
		return RetrieveValueResponse{
			Key:   requestRetrieve.Key,
			Value: "base64 me",
		}, nil
	}
}

func writeKeyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//requestWrite := request.(WriteValueRequest)
		return WriteValueResponse{}, nil
	}
}

func deleteKeyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//requestWrite := request.(DeleteKeyRequest)
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

	request := WriteValueRequest{
		Key: key,
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, BadRequest("error reading request body", err)
	}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return nil, BadRequest("error parsing request body", err)
	}

	return request, nil
}

func decodeDeleteValueRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WriteValueResponse struct {
}

type DeleteKeyRequest struct {
	Key string `json:"key"`
}

type DeleteKeyResponse struct {
}
