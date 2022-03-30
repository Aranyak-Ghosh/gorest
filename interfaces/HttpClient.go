package interfaces

import "github.com/Aranyak-Ghosh/gorest/types"

type HttpClient interface {
	SetBaseUrl(baseUrl string)
	SetMaxRetry(maxRetry int)
	Get(endpoint string, Query, headers map[string]string) *types.HttpResponse
	Post(endpoint string, Query, headers map[string]string, body interface{}) *types.HttpResponse
	Put(endpoint string, Query, headers map[string]string, body interface{}) *types.HttpResponse
	Patch(endpoint string, Query, headers map[string]string, body interface{}) *types.HttpResponse
	Del(endpoint string, Query, headers map[string]string) *types.HttpResponse
}
