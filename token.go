package iloveapigolang

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GenerateToken(ctx context.Context) error {
	_, err, _ := c.sfGroup.Do("generate_token", func() (any, error) {
		data := struct {
			PublicKey string `json:"public_key"`
		}{
			PublicKey: c.apiKey,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error encoding request: %w", err)
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			authURL,
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
			}
			return nil, fmt.Errorf("error sending request: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return nil, handleError(res)
		}

		var response struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		c.mu.Lock()
		c.token = response.Token
		fmt.Println(c.token)
		c.mu.Unlock()

		return nil, nil
	})

	return err
}
