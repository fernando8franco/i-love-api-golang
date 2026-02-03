package iloveapigolang

import (
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
