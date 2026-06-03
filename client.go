package iloveapigolang

import (
	"net/http"
	"sync"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	token      string
	mu         sync.RWMutex

	tokenInflight bool
	tokenDone     chan struct{}
}

func NewClient(httpClient *http.Client, apiKey, token string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
		token:      token,
	}
}

func (c *Client) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
}

func (c *Client) GetToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}
