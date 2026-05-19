package iloveapigolang

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type DowloadParams struct {
	Server string
	Task   string
}

func (c *Client) Download(ctx context.Context, params DowloadParams) (io.ReadCloser, error) {
	dowloadUrl := fmt.Sprintf(dowloadURL, params.Server, params.Task)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		dowloadUrl,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request:\n%v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.getToken())

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request:\n%v", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		res.Body.Close()
		return nil, handleError(res)
	}

	return res.Body, nil
}
