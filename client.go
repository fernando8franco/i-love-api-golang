package iloveapigolang

import (
	"net/http"
	"sync"
	"time"
)

type Client struct {
	httpClient *http.Client
	token      string
	mu         sync.RWMutex
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) getToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

func (c *Client) UpdateToken(newToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = newToken
}
