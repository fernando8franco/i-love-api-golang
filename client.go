package iloveapigolang

import "net/http"

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, apiKey, token string) *Client {
	return &Client{
		httpClient: httpClient,
	}
}
