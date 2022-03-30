package interfaces

import "github.com/Aranyak-Ghosh/gorest/types"

type HttpClient interface {
	SetBaseUrl(baseUrl string)
	SetMaxRetry(maxRetry int)
	Get(endpoint string, query types.Query, headers map[string]string) *HttpResponse
	Post(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *HttpResponse
	Put(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *HttpResponse
	Patch(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *HttpResponse
	Del(endpoint string, query types.Query, headers map[string]string) *HttpResponse
}
