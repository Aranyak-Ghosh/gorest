package types

import (
	"fmt"
	"net/url"
	"reflect"
)

type Query map[string]any

func (q *Query) UrlEncode() string {

	var data url.Values

	for key, value := range *q {
		if value == nil {
			delete(*q, key)
		}
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			data[key] = make([]string, len(value.([]any)))
			for _, v := range value.([]any) {
				data[key] = append(data[key], fmt.Sprintf("%v", v))
			}
		} else {
			data.Add(key, fmt.Sprintf("%v", value))
		}
	}
	return data.Encode()
}

func (q *Query) Add(key string, value any) {
	(*q)[key] = value
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
