package skyd

import (
	"net/http"
	"io/ioutil"
	"testing"
)

func assertProperty(t *testing.T, property *Property, id int64, name string, typ string, dataType string) {
	if property.Id != id {
		t.Fatalf("Unexpected property id. Expected %v, got %v", id, property.Id)
	}
	if property.Name != name {
		t.Fatalf("Unexpected property name. Expected %v, got %v", name, property.Name)
	}
	if property.Type != typ {
		t.Fatalf("Unexpected property type. Expected %v, got %v", typ, property.Type)
	}
	if property.DataType != dataType {
		t.Fatalf("Unexpected property data type. Expected %v, got %v", dataType, property.DataType)
	}
}

func assertResponse(t *testing.T, resp *http.Response, statusCode int, content string, message string) {
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode != 200 || string(body) != "" {
		t.Fatalf("%v: Expected [%v] '%v', got [%v] '%v.", message, statusCode, content, resp.StatusCode, string(body))
  }
}
