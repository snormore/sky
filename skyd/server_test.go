package skyd

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
	"testing"
)

// Ensure that we can create a new table through the server.
func TestServerCreateTable(t *testing.T) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	server := NewServer(8585, path)
	go server.ListenAndServe()
	defer server.Shutdown()

  resp, err := http.Post("http://localhost:8585/tables", "application/json", strings.NewReader(`{"name":"foo"}`))
  if err != nil {
		t.Fatalf("Unable to create table: %v", err)
  }
  assertResponse(t, resp, 200, "", "POST /tables failed.")
	if _, err := os.Stat(fmt.Sprintf("%v/foo", path)); os.IsNotExist(err) {
		t.Fatalf("POST /tables did not create table.")
	}
}

func assertResponse(t *testing.T, resp *http.Response, statusCode int, content string, message string) {
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode != 200 || string(body) != "" {
		t.Fatalf("%v: Expected [%v] '%v', got [%v] '%v.", message, statusCode, content, resp.StatusCode, string(body))
  }
}