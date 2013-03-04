package skyd

import (
  "io/ioutil"
  "os"
  "testing"
)

// Ensure that we can ping the server.
func TestServerPing(t *testing.T) {
  path, _ := ioutil.TempDir("", "")
  defer os.RemoveAll(path)
  server := NewServer(8585, path)
  go server.ListenAndServe()
  defer server.Shutdown()

  resp, err := sendTestHttpRequest("GET", "http://localhost:8585/ping", "application/json", "")
  if err != nil {
    t.Fatalf("Unable to ping: %v", err)
  }
  assertResponse(t, resp, 200, `{"message":"ok"}` + "\n", "GET /ping failed.")
}

func BenchmarkPing(b *testing.B) {
  path, _ := ioutil.TempDir("", "")
  defer os.RemoveAll(path)
  server := NewServer(8585, path)
  go server.ListenAndServe()
  defer server.Shutdown()

  for i := 0; i < b.N; i++ {
    _, _ = sendTestHttpRequest("GET", "http://localhost:8585/ping", "application/json", "")
  }
}

func BenchmarkRawPing(b *testing.B) {
  path, _ := ioutil.TempDir("", "")
  defer os.RemoveAll(path)
  server := NewServer(8585, path)
  go server.ListenAndServe()
  defer server.Shutdown()

  for i := 0; i < b.N; i++ {
    _, _ = sendTestHttpRequest("GET", "http://localhost:8585/rawping", "application/json", "")
  }
}

