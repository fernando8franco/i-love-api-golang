package iloveapigolang

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	statusCode int
	message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error: status %d\nresponse body: %s", e.statusCode, e.message)
}

func (e *APIError) IsUnauthorized() bool {
	return e.statusCode == http.StatusUnauthorized
}

func handleError(res *http.Response) error {
	var errorRes struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
		return err
	}

	return &APIError{
		statusCode: res.StatusCode,
		message:    errorRes.Message,
	}
}
