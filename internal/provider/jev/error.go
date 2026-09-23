package jev

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("Jev API returned HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("Jev API returned HTTP %d: %s", e.StatusCode, e.Message)
}
