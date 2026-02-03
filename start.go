package iloveapigolang

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StartResponse struct {
	Server           string `json:"server"`
	Task             string `json:"task"`
	RemainingCredits int    `json:"remaining_credits"`
}

func (ac ApiCredentials) Start(tool, region string) (response StartResponse, err error) {
	url := fmt.Sprintf(startURL, tool, region)

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return StartResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+ac.AuthToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return StartResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return StartResponse{}, handleError(res)
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return StartResponse{}, err
	}

	return response, nil
}
