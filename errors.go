package iloveapigolang

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	statusCode int
	message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.statusCode, e.message)
}

func (e *APIError) StatusCode() int {
	return e.statusCode
}

func handleError(res *http.Response) error {
	apiErr := &APIError{
		statusCode: res.StatusCode,
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		apiErr.message = "failed to read error response body"
		return apiErr
	}

	apiErr.message = string(bodyBytes)

	var errorRes struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(bodyBytes, &errorRes); err == nil && errorRes.Message != "" {
		apiErr.message = errorRes.Message
	}

	return apiErr
}
