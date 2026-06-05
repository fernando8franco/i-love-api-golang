package iloveapigolang

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GenerateToken(ctx context.Context) error {
	c.mu.Lock()

	if c.tokenInflight {
		ch := c.tokenDone
		c.mu.Unlock()
		select {
		case <-ch:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	c.tokenInflight = true
	c.tokenDone = make(chan struct{})
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.tokenInflight = false
		close(c.tokenDone)
		c.mu.Unlock()
	}()

	data := struct {
		PublicKey string `json:"public_key"`
	}{
		PublicKey: c.apiKey,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		authURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		}
		return fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return handleError(res)
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	c.mu.Lock()
	c.token = response.Token
	c.mu.Unlock()

	return nil
}
