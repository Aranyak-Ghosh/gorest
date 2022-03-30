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
		ErrorMessage: "Failed to serialize supplied body with error",
	}
}

func ValidationError(attrName string) *HttpClientError {
	return &HttpClientError{
		ErrorCode:    2001,
		ErrorDetails: fmt.Errorf("Validation for attribute %s failed", attrName),
		ErrorMessage: "Failed to validate attribute",
	}
}

func UnsupportedMIMETypeError(additional string) *HttpClientError {
	return &HttpClientError{
		ErrorCode:    2002,
		ErrorDetails: fmt.Errorf("Unsupported mime type"),
		ErrorMessage: fmt.Sprintf("Specified Content type is not currently supported.%s", additional),
	}
}

func UnMarshallError(err error) *HttpClientError {
	return &HttpClientError{
		ErrorCode:    2003,
		ErrorDetails: err,
		ErrorMessage: "Failed to unmarshal response to given object",
	}
}
