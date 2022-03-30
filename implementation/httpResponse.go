package implementation

import (
	"fmt"
	"net/http"

	"github.com/Aranyak-Ghosh/gorest/interfaces"
)

type httpResponse struct {
	responseData    []byte
	receivedError   error
	statusCode      int
	responseHeaders http.Header
}

var _ interfaces.HttpResponse = (*httpResponse)(nil)

func (h *httpResponse) Result(val *any) error {
	if h.receivedError != nil {
		return h.receivedError
	}

	return nil
}

func (h *httpResponse) IsSuccessfulResponse() bool {
	return h.statusCode >= 200 && h.statusCode < 300
}

func (h *httpResponse) Status() int {
	return h.statusCode
}
func (h *httpResponse) Headers() http.Header {
	return http.Header(h.responseHeaders)
}
func (h *httpResponse) RawData() []byte {
	return []byte(h.responseData)
}
func (h *httpResponse) Error() error {
	return fmt.Errorf("%w", h.receivedError)
}
