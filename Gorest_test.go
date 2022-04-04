package gorest_test

import (
	"reflect"
	"testing"

	"github.com/Aranyak-Ghosh/gorest"
	"github.com/Aranyak-Ghosh/gorest/types"
)

func TestGet(t *testing.T) {
	ep := "http://localhost:3000/"
	cl := gorest.NewClient()

	cl.SetBaseUrl(ep)

	resp := cl.Get("", nil, nil)

	if resp.Error() != nil || resp.Status() != 200 {
		t.Errorf("Expected status code 200, got %d", resp.Status())
	}
}

func TestJsonPost(t *testing.T) {
	ep := "http://localhost:3000/"
	cl := gorest.NewClient()

	cl.SetBaseUrl(ep)

	expected := map[string]string{
		"name": "Aranyak Ghosh",
	}
	var data map[string]string

	resp := cl.Post("", nil, nil, expected, types.JSON)

	if resp.Error() != nil {
		t.Errorf("Expected no error, got %v", resp.Error())
	} else if !resp.IsSuccessfulResponse() {
		t.Errorf("Expected success status, got %d", resp.Status())
	} else {
		resp.Result(&data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	}
}

func TestXMLPost(t *testing.T) {
	ep := "http://localhost:3000/"
	cl := gorest.NewClient()

	cl.SetBaseUrl(ep)

	type bodyModel struct {
		Name string `xml:"name"`
	}

	expected := bodyModel{
		Name: "Aranyak Ghosh",
	}
	var data bodyModel

	resp := cl.Post("", nil, nil, expected, types.XML)

	if resp.Error() != nil {
		t.Errorf("Expected no error, got %v", resp.Error())
	} else if !resp.IsSuccessfulResponse() {
		t.Errorf("Expected success status, got %d", resp.Status())
	} else {
		resp.Result(&data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	}
}

func TestJsonPatch(t *testing.T) {
	ep := "http://localhost:3000/"
	cl := gorest.NewClient()

	cl.SetBaseUrl(ep)

	type bodyModel struct {
		Name string
	}

	expected := bodyModel{
		Name: "Aranyak Ghosh",
	}
	var data bodyModel

	resp := cl.Patch("", nil, nil, expected, types.JSON)

	if resp.Error() != nil {
		t.Errorf("Expected no error, got %v", resp.Error())
	} else if !resp.IsSuccessfulResponse() {
		t.Errorf("Expected success status, got %d", resp.Status())
	} else {
		resp.Result(&data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	}
}
