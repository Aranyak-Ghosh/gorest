package interfaces

import (
	"time"

	"github.com/Aranyak-Ghosh/gorest/types"
)

type HttpClient interface {
	SetBaseUrl(baseUrl string) error
	SetMaxRetry(maxRetry int)
	SetTimeout(timeout time.Duration)
	Del(endpoint string, query types.Query, headers map[string]string) *HttpResponse
	Get(endpoint string, query types.Query, headers map[string]string) *HttpResponse
	Head(endpoint string, headers map[string]string) *HttpResponse
	Post(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *HttpResponse
	Put(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *HttpResponse
	Patch(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *HttpResponse
}
