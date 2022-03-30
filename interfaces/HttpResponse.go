package interfaces

import "net/http"

type HttpResponse interface {
	Result(*any) error
	IsSuccessfulResponse() bool
	Status() int
	Headers() http.Header
	RawData() []byte
}