package gorest

import (
	"github.com/Aranyak-Ghosh/gorest/implementation"
	"github.com/Aranyak-Ghosh/gorest/interfaces"
)

func NewClient() interfaces.HttpClient {
	return implementation.Client()
}
