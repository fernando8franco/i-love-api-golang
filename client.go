package iloveapigolang

import (
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	token      string
}

func NewClient(httpClient *http.Client, apiKey string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}
