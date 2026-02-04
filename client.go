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

func (c *Client) SetAPIKey(newAPIKey string) {
	c.apiKey = newAPIKey
}

func (c *Client) SetToken(newToken string) {
	c.token = newToken
}

func (c *Client) GetToken() string {
	return c.token
}
