package iloveapigolang

import "net/http"

type Client struct {
	httpClient *http.Client
	apiKey     string
	token      string
}

func NewClient(httpClient *http.Client, apiKey, token string) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
		token:      token,
	}
}

func (c *Client) UpdateAPIKey(newAPIKey string) {
	c.apiKey = newAPIKey
}

func (c *Client) UpdateToken(newToken string) {
	c.token = newToken
}
