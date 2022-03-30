package implementation

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/Aranyak-Ghosh/gorest/types"
	"github.com/asaskevich/govalidator"

	"github.com/Aranyak-Ghosh/gorest/interfaces"
)

type httpClient struct {
	client   *http.Client
	baseUrl  string
	maxRetry int
}

var _ interfaces.HttpClient = (*httpClient)(nil)

const (
	defaultMaxRetry = 3
)

// SetBaseUrl sets the base url for the http client
func (h *httpClient) SetBaseUrl(baseUrl string) error {
	err := validateUrl(baseUrl)
	if err != nil {
		return err
	} else {
		h.baseUrl = baseUrl
		return nil
	}
}

func (h *httpClient) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

// SetMaxRetry sets the max retry for the http client
func (h *httpClient) SetMaxRetry(maxRetry int) {
	h.maxRetry = maxRetry
}

// Delete makes a delete request to the http client
func (h *httpClient) Del(endpoint string, query types.Query, headers map[string]string) *interfaces.HttpResponse {
	return nil
}

// Get makes a get request to the http client
func (h *httpClient) Get(endpoint string, query types.Query, headers map[string]string) *interfaces.HttpResponse {
	return nil
}

// Head makes a head request
func (h *httpClient) Head(endpoint string, headers map[string]string) *interfaces.HttpResponse {
	return nil
}

// Post makes a post request to the http client
func (h *httpClient) Post(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *interfaces.HttpResponse {
	return nil
}

// Put makes a put request to the http client
func (h *httpClient) Put(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *interfaces.HttpResponse {
	return nil
}

// Patch makes a patch request to the http client
func (h *httpClient) Patch(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) *interfaces.HttpResponse {
	return nil
}

func isSuccessResponse(response *httpResponse) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}

func handleRequestBody(request *http.Request, body any, bodyType types.ContentType) error {

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", bodyType.Header())
	}

	switch bodyType {
	case types.JSON:
		dat, err := json.Marshal(body)
		if err != nil {
			return types.SerializeBodyError(err)
		} else {
			request.Body = ioutil.NopCloser(bytes.NewReader(dat))
		}
	case types.FormData, types.EncodedFormData:
		if reflect.TypeOf(body).Kind() != reflect.Map {
			return types.SerializeBodyError(fmt.Errorf("Body must be of type map with form content type"))
		} else {
			var mapData = map[string]any(body.(map[string]any))

			var data url.Values

			for key, value := range mapData {

				if value == nil {
					delete(mapData, key)
				}

				if reflect.TypeOf(value).Kind() == reflect.Slice {
					data[key] = make([]string, len(value.([]any)))
					for _, v := range value.([]any) {
						d, err := json.Marshal(v)
						if err == nil {
							data[key] = append(data[key], string(d))
						} else {
							return types.SerializeBodyError(fmt.Errorf("Failed to serialize key %s", key))
						}
					}
				} else {
					d, err := json.Marshal(value)
					if err != nil {
						return types.SerializeBodyError(fmt.Errorf("Failed to serialize key %s", key))
					} else {
						data.Add(key, string(d))
					}
				}
			}

			var serializedData = data.Encode()

			request.Body = ioutil.NopCloser(strings.NewReader(serializedData))
		}
	case types.XML:
		dat, err := xml.Marshal(body)
		if err != nil {
			return types.SerializeBodyError(err)
		} else {
			request.Body = ioutil.NopCloser(bytes.NewReader(dat))
		}
	case types.Binary, types.RAW:
		var b bytes.Buffer
		enc := gob.NewEncoder(&b)
		if err := enc.Encode(body); err != nil {
			return types.SerializeBodyError(err)
		} else {
			request.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
		}
	}

	return nil
}

func constructHeaders(req *http.Request, headers map[string]string) {
	for key, val := range headers {
		req.Header.Add(key, val)
	}
}

func handleResponse(response *http.Response, res *httpResponse) {
	res.StatusCode = response.StatusCode
	res.ResponseHeaders = map[string][]string(response.Header)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		res.Error = err
	} else {
		res.ResponseData = body
	}
}

func validateUrl(url string) error {
	if ok := govalidator.IsURL(url); !ok {
		return types.ValidationError("Url")
	} else {
		return nil
	}
}
