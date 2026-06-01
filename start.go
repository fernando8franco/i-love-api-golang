package iloveapigolang

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type StartParams struct {
	Token  string
	Tool   string
	Region string
}

type StartResponse struct {
	Server           string `json:"server"`
	Task             string `json:"task"`
	RemainingCredits int    `json:"remaining_credits"`
}

func (c *Client) Start(ctx context.Context, params StartParams) (StartResponse, error) {
	url := fmt.Sprintf(startURL, params.Tool, params.Region)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return StartResponse{}, fmt.Errorf("error creating request:\n%v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.GetToken())

	res, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return StartResponse{}, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		}
		return StartResponse{}, fmt.Errorf("error sending request:\n%v", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return StartResponse{}, handleError(res)
	}

	var response StartResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return StartResponse{}, fmt.Errorf("error decoding response:\n%v", err)
	}

	return response, nil
}
