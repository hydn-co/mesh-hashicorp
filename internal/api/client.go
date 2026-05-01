package api

import "net/http"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}
