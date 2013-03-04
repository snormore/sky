package skyd

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "os"
  "strings"
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
  if resp.StatusCode != 200 || content != string(body) {
    t.Fatalf("%v: Expected [%v] %q, got [%v] %q.", message, statusCode, content, resp.StatusCode, string(body))
  }
}

func sendTestHttpRequest(method string, url string, contentType string, body string) (*http.Response, error) {
  client := &http.Client{Transport: &http.Transport{DisableKeepAlives:true}}
  req, _ := http.NewRequest(method, url, strings.NewReader(body))
  req.Header.Add("Content-Type", contentType)
  return client.Do(req)
}


func runTestServer(f func()) {
  path, _ := ioutil.TempDir("", "")
  defer os.RemoveAll(path)
  server := NewServer(8585, path)
  go server.ListenAndServe()
  defer server.Shutdown()
  f()
}

func setupTestTable(name string) {
  _, _ = sendTestHttpRequest("POST", "http://localhost:8585/tables", "application/json", fmt.Sprintf(`{"name":"%v"}`, name))
}

func setupTestProperty(tableName string, name string, typ string, dataType string) {
  _, _ = sendTestHttpRequest("POST", fmt.Sprintf("http://localhost:8585/tables/%v/properties", tableName), "application/json", fmt.Sprintf(`{"name":"%v", "type":"%v", "dataType":"%v"}`, name, typ, dataType))
}
