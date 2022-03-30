package types

type Query map[string]string

type HttpResponse struct {
	ResponseData    []byte
	Error           error
	StatusCode      int
	ResponseHeaders map[string][]string
}

type ContentType int

const (
	JSON ContentType = iota
	XML
	RAW
	FormData
	EncodedFormData
	Binary
)

func (c *ContentType) Header() string {
	switch *c {
	case JSON:
		return "application/json"
	case XML:
		return "application/xml"
	case RAW:
		return "application/octet-stream"
	case FormData:
		return "application/x-www-form-urlencoded"
	case EncodedFormData:
		return "application/x-www-form-urlencoded"
	case Binary:
		return "application/octet-stream"
	default:
		return ""
	}
}

func (c *ContentType) FromHeader(header string) ContentType {
	switch header {
	case "application/json":
		return JSON
	case "application/xml":
		return XML
	case "application/octet-stream":
		return RAW
	case "application/x-www-form-urlencoded":
		return FormData
	default:
		return JSON
	}
}
