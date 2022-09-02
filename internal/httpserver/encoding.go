package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
)

const applicationJSON = "application/json; charset=utf-8"
const textPlain = "text/plain; charset=utf-8"

func noOpDecode(_ context.Context, _ *http.Request) (interface{}, error) { return nil, nil }

func encodeText(_ context.Context, w http.ResponseWriter, response interface{}) error {
	text := response.(string)
	w.Header().Set("Content-Type", textPlain)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(text))
	return nil
}

func encodeJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	jsonResp, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", applicationJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
	return nil
}

func encodeNoContent(_ context.Context, w http.ResponseWriter, _ interface{}) error {
	w.Header().Set("Content-Type", applicationJSON)
	w.WriteHeader(http.StatusNoContent)
	return nil
}
