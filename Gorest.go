package gorest

import (
	"github.com/Aranyak-Ghosh/gorest/implementation"
	"github.com/Aranyak-Ghosh/gorest/interfaces"
)

type HttpClient interfaces.HttpClient
type HttpResponse interfaces.HttpResponse

func NewClient() HttpClient {
	return implementation.Client()
}
