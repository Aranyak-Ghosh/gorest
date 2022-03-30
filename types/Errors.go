package types

type HttpClientError struct {
	ErrorCode    int
	ErrorMessage string
	ErrorDetails error
}

func (e *HttpClientError) Error() string {
	return e.ErrorMessage
}
