package iloveapigolang

import (
	"fmt"
	"io"
	"net/http"
)

func (c *Client) Dowload(server, task string, file io.Writer) (io.ReadCloser, error) {
	dowloadUrl := fmt.Sprintf(dowloadURL, server, task)

	req, err := http.NewRequest("GET", dowloadUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request:\n%v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request:\n%v", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		defer res.Body.Close()
		return nil, handleError(res)
	}

	return res.Body, nil
}
