package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) DeleteKey(_ context.Context, request *DeleteKeyRequest) (DeleteKeyResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return DeleteKeyResponse{}, err
	}

	reqURL.Path = path.Join(
		reqURL.Path,
		fmt.Sprintf(
			"/v1/keys/%s",
			url.PathEscape(request.Key),
		),
	)

	req := &http.Request{
		Method: http.MethodDelete,
		URL:    reqURL,
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return DeleteKeyResponse{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode >= 300 {
		return DeleteKeyResponse{}, decodeErrorJSON(response)
	}

	if response.ContentLength == 0 {
		return DeleteKeyResponse{}, nil
	}

	var ret DeleteKeyResponse
	err = decodeJSON(response, &ret)
	if err != nil {
		return DeleteKeyResponse{}, err
	}

	return ret, nil
}
