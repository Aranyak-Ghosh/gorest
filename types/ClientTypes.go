package types

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

type Query map[string]any

func (q *Query) UrlEncode() (string, error) {

	var data url.Values

	for key, value := range *q {

		if value == nil {
			delete(*q, key)
		}

		if reflect.TypeOf(value).Kind() == reflect.Slice {
			data[key] = make([]string, len(value.([]any)))
			for _, v := range value.([]any) {
				d, err := json.Marshal(v)
				if err == nil {
					data[key] = append(data[key], string(d))
				} else {
					return "", SerializeBodyError(fmt.Errorf("Failed to serialize key %s", key))
				}
			}
		} else {
			d, err := json.Marshal(value)
			if err != nil {
				return "", SerializeBodyError(fmt.Errorf("Failed to serialize key %s", key))
			} else {
				data.Add(key, string(d))
			}
		}
	}
	return data.Encode(), nil
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

func (c *ContentType) FromHeader(header string) bool {
	switch header {
	case "application/json":
		*c = JSON
		return true
	case "application/xml":
		*c = XML
		return true
	case "application/octet-stream":
		*c = RAW
		return true
	case "application/x-www-form-urlencoded":
		*c = FormData
		return true
	default:
		//TODO: Current implementation defaults to handling other content-type as octet-stream. If an error needs to be thrown, default return can be changed to false. Alternatively, more enums and cases can be added and a PR can be created to handle other content-types
		*c = RAW
		return true
	}
}
