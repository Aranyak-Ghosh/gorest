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
	maxRetry *int
}

const blankString = ""

const (
	GET   = "GET"
	POST  = "POST"
	PUT   = "PUT"
	PATCH = "PATCH"
	DEL   = "DELETE"
	HEAD  = "HEAD"
)

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
	h.maxRetry = new(int)
	*h.maxRetry = maxRetry
}

// Delete makes a delete request to the http client
func (h *httpClient) Del(endpoint string, query types.Query, headers map[string]string) interfaces.HttpResponse {
	var res = &httpResponse{}

	if url, err := h.constructAndValidateUrl(endpoint, query); err != nil {
		res.Error = err
	} else {
		if req, e := http.NewRequest(DEL, url, nil); e != nil {
			res.Error = e
		} else {
			constructHeaders(req, headers)
			h.exec(req, res)
		}
	}
	return res
}

// Get makes a get request to the http client
func (h *httpClient) Get(endpoint string, query types.Query, headers map[string]string) interfaces.HttpResponse {

	var res = &httpResponse{}

	if url, err := h.constructAndValidateUrl(endpoint, query); err != nil {
		res.Error = err
	} else {
		if req, e := http.NewRequest(GET, url, nil); e != nil {
			res.Error = e
		} else {
			constructHeaders(req, headers)
			h.exec(req, res)
		}
	}
	return res
}

// Head makes a head request
func (h *httpClient) Head(endpoint string, headers map[string]string) interfaces.HttpResponse {
	var res = &httpResponse{}

	if url, err := h.constructAndValidateUrl(endpoint, nil); err != nil {
		res.Error = err
	} else {
		if req, e := http.NewRequest(HEAD, url, nil); e != nil {
			res.Error = e
		} else {
			constructHeaders(req, headers)
			h.exec(req, res)
		}
	}
	return res
}

// Post makes a post request to the http client
func (h *httpClient) Post(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) interfaces.HttpResponse {

	var res = &httpResponse{}

	if url, err := h.constructAndValidateUrl(endpoint, query); err != nil {
		res.Error = err
	} else {
		if req, e := http.NewRequest(POST, url, nil); e != nil {
			res.Error = e
		} else {
			constructHeaders(req, headers)
			e = handleRequestBody(req, body, bodyType)
			if e != nil {
				res.Error = e
			} else {
				h.exec(req, res)
			}
		}
	}
	return res
}

// Put makes a put request to the http client
func (h *httpClient) Put(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) interfaces.HttpResponse {
	var res = &httpResponse{}

	if url, err := h.constructAndValidateUrl(endpoint, query); err != nil {
		res.Error = err
	} else {
		if req, e := http.NewRequest(PUT, url, nil); e != nil {
			res.Error = e
		} else {
			constructHeaders(req, headers)
			e = handleRequestBody(req, body, bodyType)
			if e != nil {
				res.Error = e
			} else {
				h.exec(req, res)
			}

			h.exec(req, res)
		}
	}
	return res
}

// Patch makes a patch request to the http client
func (h *httpClient) Patch(endpoint string, query types.Query, headers map[string]string, body any, bodyType types.ContentType) interfaces.HttpResponse {
	var res = &httpResponse{}

	if url, err := h.constructAndValidateUrl(endpoint, query); err != nil {
		res.Error = err
	} else {
		if req, e := http.NewRequest(PUT, url, nil); e != nil {
			res.Error = e
		} else {
			constructHeaders(req, headers)
			e = handleRequestBody(req, body, bodyType)
			if e != nil {
				res.Error = e
			} else {
				h.exec(req, res)
			}

			h.exec(req, res)
		}
	}
	return res
}

func isSuccessResponse(response *httpResponse) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}

func handleRequestBody(request *http.Request, body any, bodyType types.ContentType) error {

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", bodyType.Header())
	} else {
		if ok := bodyType.FromHeader(request.Header.Get("Content-Type")); !ok {
			return fmt.Errorf("Invalid Content-Type header set")
		} else {
			request.Header.Set("Content-Type", bodyType.Header())
		}
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
	if ok := (!isNullOrEmpty(url) && govalidator.IsURL(url)); !ok {
		return types.ValidationError("Url")
	} else {
		return nil
	}
}

func isNullOrEmpty(str string) bool {
	return (str == blankString)
}

func (h *httpClient) constructAndValidateUrl(endpoint string, query types.Query) (string, error) {
	var url string
	if isNullOrEmpty(h.baseUrl) {
		url = endpoint
	} else {
		url = fmt.Sprintf("%s%s", url, endpoint)
	}

	if query != nil {
		q, err := query.UrlEncode()
		if err != nil {
			return "", err
		}
		url = fmt.Sprintf("%s?%s", url, q)
	}
	return url, validateUrl(url)
}

func (h *httpClient) exec(req *http.Request, res *httpResponse) {
	var isSuccess = false
	var attempt = 0

	if h.maxRetry == nil {
		h.SetMaxRetry(defaultMaxRetry)
	}

	for (attempt <= *h.maxRetry) && !isSuccess {
		response, e := h.client.Do(req)
		if e != nil {
			res.Error = e
			attempt++
		} else {
			isSuccess = true
			handleResponse(response, res)
		}
	}
}
