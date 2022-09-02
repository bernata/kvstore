package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) WriteValue(_ context.Context, request *WriteValueRequest) (WriteValueResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return WriteValueResponse{}, err
	}

	reqURL.Path = path.Join(
		reqURL.Path,
		fmt.Sprintf(
			"/v1/keys/%s",
			url.PathEscape(request.Key),
		),
	)

	body, err := json.Marshal(request)
	if err != nil {
		return WriteValueResponse{}, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Body:   io.NopCloser(bytes.NewReader(body)),
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return WriteValueResponse{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode >= 300 {
		return WriteValueResponse{}, decodeErrorJSON(response)
	}

	if response.ContentLength == 0 {
		return WriteValueResponse{}, nil
	}

	var ret WriteValueResponse
	err = decodeJSON(response, &ret)
	if err != nil {
		return WriteValueResponse{}, err
	}

	return ret, nil
}
