package iloveapigolang

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ApiCredentials struct {
	APIKey    string
	AuthToken string
}

func (ac ApiCredentials) GetToken() (token string, err error) {
	data := struct {
		PublicKey string `json:"public_key"`
	}{
		PublicKey: ac.APIKey,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		authURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", handleError(res)
	}

	response := struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Token, nil
}
