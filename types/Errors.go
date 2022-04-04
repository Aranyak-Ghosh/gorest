package types

import "fmt"

type HttpClientError struct {
	ErrorCode    ErrorCode
	ErrorMessage string
	ErrorDetails error
}

type ErrorCode int

const (
	SerializeBodyErrorCode ErrorCode = 2000 + iota
	ValidationErrorCode
	UnsupportedMIMETypeErrorCode
	UnMarshallErrorCode
)

func (e *HttpClientError) Error() string {
	return fmt.Errorf("%s: %w", e.ErrorMessage, e.ErrorDetails).Error()
}

func SerializeBodyError(err error) *HttpClientError {
	if err == nil {
		return nil
	} else {
		return &HttpClientError{
			ErrorCode:    SerializeBodyErrorCode,
			ErrorDetails: err,
			ErrorMessage: "Failed to serialize supplied body with error",
		}
	}
}

func ValidationError(attrName string) *HttpClientError {
	return &HttpClientError{
		ErrorCode:    ValidationErrorCode,
		ErrorDetails: fmt.Errorf("Validation for attribute %s failed", attrName),
		ErrorMessage: "Failed to validate attribute",
	}
}

func UnsupportedMIMETypeError(additional string) *HttpClientError {
	return &HttpClientError{
		ErrorCode:    UnsupportedMIMETypeErrorCode,
		ErrorDetails: fmt.Errorf("Unsupported mime type"),
		ErrorMessage: fmt.Sprintf("Specified Content type is not currently supported.%s", additional),
	}
}

func UnMarshallError(err error) *HttpClientError {
	if err != nil {
		return &HttpClientError{
			ErrorCode:    UnMarshallErrorCode,
			ErrorDetails: err,
			ErrorMessage: "Failed to unmarshal response to given object",
		}
	} else {
		return nil
	}
}
