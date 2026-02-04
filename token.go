package iloveapigolang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GenerateToken() (string, error) {
	data := struct {
		PublicKey string `json:"public_key"`
	}{
		PublicKey: c.apiKey,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error encoding request:\n%v", err)
	}

	req, err := http.NewRequest(
		"POST",
		authURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("error creating request:\n%v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request:\n%v", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", handleError(res)
	}

	response := struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response:\n%v", err)
	}

	return response.Token, nil
}
