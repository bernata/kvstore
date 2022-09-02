package apiclient

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPDoer
	baseURL    string
}

type ClientOption func(*Client)

func HTTPOption(httpDo HTTPDoer) ClientOption {
	return func(c *Client) {
		c.httpClient = httpDo
	}
}

func BaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: defaultClient(),
		baseURL:    "http://localhost:8282",
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func defaultClient() *http.Client {
	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	return &http.Client{
		Transport: httpTransport,
		Timeout:   time.Second * 60,
	}
}

func decodeErrorJSON(response *http.Response) error {
	var errDefault KeyValueError
	err := decodeJSON(response, &errDefault)
	if err == nil {
		err = &errDefault
	}
	return err
}

func decodeJSON(response *http.Response, data interface{}) error {
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBytes, data)
}
