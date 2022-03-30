package types

type Query map[string]string

type HttpResponse struct {
	ResponseData    []byte
	Error           error
	ResponseCode    int
	ResponseHeaders map[string][]string
}
