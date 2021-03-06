package implementation

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Aranyak-Ghosh/gorest/interfaces"
	"github.com/Aranyak-Ghosh/gorest/types"
)

type httpResponse struct {
	responseData    []byte
	receivedError   error
	statusCode      int
	responseHeaders http.Header
	nativeResponse  *http.Response
}

var _ interfaces.HttpResponse = (*httpResponse)(nil)

func (h *httpResponse) Result(val any) error {
	if h.receivedError != nil {
		return h.receivedError
	}

	ct := h.responseHeaders.Get("Content-Type")

	var contentType types.ContentType

	if h.responseData == nil {
		h.responseData = h.RawData()
	}

	if ok := contentType.FromHeader(ct); !ok {
		return types.UnsupportedMIMETypeError("Access raw bytes by using RawData method")
	}

	switch contentType {
	case types.JSON:
		err := json.Unmarshal(h.responseData, val)
		if err != nil {
			return types.UnMarshallError(err)
		} else {
			return nil
		}
	case types.XML:
		err := xml.Unmarshal(h.responseData, val)
		if err != nil {
			return types.UnMarshallError(err)
		} else {
			return nil
		}
	default:
		return types.UnsupportedMIMETypeError("Access raw bytes by using RawData method")
	}
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

	if h.responseData == nil {
		defer h.nativeResponse.Body.Close()
		body, err := ioutil.ReadAll(h.nativeResponse.Body)

		if err != nil {
			h.receivedError = err
		} else {
			h.responseData = body
		}
	}
	return []byte(h.responseData)
}
func (h *httpResponse) Error() error {
	if h.receivedError != nil {
		return fmt.Errorf("%w", h.receivedError)
	} else {
		return nil
	}
}
func (h *httpResponse) RawResponse() *http.Response {
	return h.nativeResponse
}

func isJsonArray(data []byte) bool {
	var x []byte

	copy(x, data)

	x = bytes.TrimLeft(x, " \t\r\n")

	return len(x) > 0 && x[0] == '['
}
