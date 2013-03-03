package skyd

import (
	"io/ioutil"
	"os"
	"testing"
)

// Ensure that we can create a property through the server.
func TestServerCreateProperty(t *testing.T) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	server := NewServer(8585, path)
	go server.ListenAndServe()
	defer server.Shutdown()
  setupTestTable("foo")
  
  resp, err := sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/properties", "application/json", `{"name":"bar", "type":"object", "dataType":"string"}`)
  if err != nil {
		t.Fatalf("Unable to create property: %v", err)
  }
  assertResponse(t, resp, 200, `{"id":1,"name":"bar","type":"object","dataType":"string"}`+"\n", "POST /tables/:name/properties failed.")
}
