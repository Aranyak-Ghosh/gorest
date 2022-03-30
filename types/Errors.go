package types

import "fmt"

type HttpClientError struct {
	ErrorCode    int
	ErrorMessage string
	ErrorDetails error
}

func (e *HttpClientError) Error() string {
	return fmt.Errorf("%s: %w", e.ErrorMessage, e.ErrorDetails).Error()
}

func SerializeBodyError(err error) *HttpClientError {
	return &HttpClientError{
		ErrorCode:    2000,
		ErrorDetails: err,
		ErrorMessage: "failed to serialize supplied body with error",
	}
}
