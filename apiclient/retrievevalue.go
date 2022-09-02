package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) RetrieveValue(_ context.Context, request *RetrieveValueRequest) (RetrieveValueResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return RetrieveValueResponse{}, err
	}

	reqURL.Path = path.Join(
		reqURL.Path,
		fmt.Sprintf(
			"/v1/keys/%s",
			url.PathEscape(request.Key),
		),
	)

	req := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return RetrieveValueResponse{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode >= 300 {
		return RetrieveValueResponse{}, decodeErrorJSON(response)
	}

	var ret RetrieveValueResponse
	err = decodeJSON(response, &ret)
	if err != nil {
		return RetrieveValueResponse{}, err
	}

	return ret, nil
}
