package httpserver

import (
	"fmt"
	"net/http"
)

type KeyValueError struct {
	StatusCode int
	Message    string
	Cause      error
}

func (kve KeyValueError) Error() string {
	if kve.Cause == nil {
		return fmt.Sprintf("[%d]: %s", kve.StatusCode, kve.Message)
	}
	return fmt.Sprintf("[%d]: %s: %s", kve.StatusCode, kve.Message, kve.Cause)
}

func BadRequest(message string, cause error) KeyValueError {
	return KeyValueError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Cause:      cause,
	}
}
