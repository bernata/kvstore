package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type KeyValueError struct {
	StatusCodeHTTP int    `json:"status_code"`
	Message        string `json:"message"`
	Cause          error  `json:"-"`
}

func (kve KeyValueError) Error() string {
	if kve.Cause == nil {
		return fmt.Sprintf("[%d]: %s", kve.StatusCodeHTTP, kve.Message)
	}
	return fmt.Sprintf("[%d]: %s: %s", kve.StatusCodeHTTP, kve.Message, kve.Cause)
}

func (kve KeyValueError) StatusCode() int {
	return kve.StatusCodeHTTP
}

func (kve KeyValueError) MarshalJSON() ([]byte, error) {
	err := struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	}{
		Message:    kve.Error(),
		StatusCode: kve.StatusCode(),
	}
	return json.Marshal(err)
}

func BadRequest(message string, cause error) KeyValueError {
	return KeyValueError{
		StatusCodeHTTP: http.StatusBadRequest,
		Message:        message,
		Cause:          cause,
	}
}

func NotFoundErrorf(message string, parameters ...any) KeyValueError {
	return KeyValueError{
		StatusCodeHTTP: http.StatusNotFound,
		Message:        fmt.Sprintf(message, parameters...),
	}
}

func InternalServerError(message string, cause error) KeyValueError {
	return KeyValueError{
		StatusCodeHTTP: http.StatusInternalServerError,
		Message:        message,
		Cause:          cause,
	}
}
