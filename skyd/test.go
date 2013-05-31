package skyd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const TestPort = 8800
const TestElectionTimeout = 60 * time.Millisecond
const TestHeartbeatTimeout = 20 * time.Millisecond

func assertProperty(t *testing.T, property *Property, id int64, name string, transient bool, dataType string) {
	if property.Id != id {
		t.Fatalf("Unexpected property id. Expected %v, got %v", id, property.Id)
	}
	if property.Name != name {
		t.Fatalf("Unexpected property name. Expected %v, got %v", name, property.Name)
	}
	if property.Transient != transient {
		t.Fatalf("Unexpected property transiency. Expected %v, got %v", transient, property.Transient)
	}
	if property.DataType != dataType {
		t.Fatalf("Unexpected property data type. Expected %v, got %v", dataType, property.DataType)
	}
}

func assertResponse(t *testing.T, resp *http.Response, statusCode int, content string, message string) {
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 || content != string(body) {
		t.Fatalf("%v:\nexp:[%v] %s\ngot:[%v] %s.", message, statusCode, content, resp.StatusCode, string(body))
	}
}

func sendTestHttpRequest(method string, url string, contentType string, body string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Add("Content-Type", contentType)
	return client.Do(req)
}

func runTestServerAt(path string, f func(s *Server)) {
	server, err := NewServer(TestPort, path)
	if err != nil {
		panic(fmt.Sprintf("Unable to run test server: %s", err))
	}
	server.Silence()
	server.ListenAndServe(nil)
	defer server.Shutdown()
	f(server)
}

func runTestServer(f func(s *Server)) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	runTestServerAt(path, f)
}

func runTestServers(callbacks ...func(s *Server)) {
	var wg sync.WaitGroup
	servers := []*Server{}
	for i, _ := range callbacks {
		wg.Add(1)

		path, _ := ioutil.TempDir("", "")
		defer os.RemoveAll(path)

		server, err := NewServer(uint(TestPort+i), path)
		if err != nil {
			panic(fmt.Sprintf("Unable to run test server: %s", err))
		}
		server.ElectionTimeout = TestElectionTimeout
		server.HeartbeatTimeout = TestHeartbeatTimeout
		server.Silence()
		server.ListenAndServe(nil)
		defer server.Shutdown()

		servers = append(servers, server)
	}

	// Execute all the callbacks at the same time.
	for i, _f := range callbacks {
		f := _f
		go func() {
			defer wg.Done()
			f(servers[i])
		}()
	}

	// Wait until everything is done.
	wg.Wait()
}

func createTempTable(t *testing.T) *Table {
	path, err := ioutil.TempDir("", "")
	os.RemoveAll(path)

	table := NewTable("test", path)
	err = table.Create()
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}

	return table
}

func setupTestTable(name string) {
	resp, _ := sendTestHttpRequest("POST", "http://localhost:8800/tables", "application/json", fmt.Sprintf(`{"name":"%v"}`, name))
	resp.Body.Close()
}

func setupTestProperty(tableName string, name string, transient bool, dataType string) {
	resp, _ := sendTestHttpRequest("POST", fmt.Sprintf("http://localhost:8800/tables/%v/properties", tableName), "application/json", fmt.Sprintf(`{"name":"%v", "transient":%v, "dataType":"%v"}`, name, transient, dataType))
	resp.Body.Close()
}

func setupTestData(t *testing.T, tableName string, items [][]string) {
	for i, item := range items {
		resp, _ := sendTestHttpRequest("PUT", fmt.Sprintf("http://localhost:8800/tables/%s/objects/%s/events/%s", tableName, item[0], item[1]), "application/json", item[2])
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("setupTestData[%d]: Expected 200, got %v.", i, resp.StatusCode)
		}
	}
}

func _codegen(t *testing.T, tableName string, query string) {
	resp, _ := sendTestHttpRequest("POST", fmt.Sprintf("http://localhost:8800/tables/%s/query/codegen", tableName), "application/json", query)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func _dumpObject(t *testing.T, tableName string, objectId string) {
	resp, _ := sendTestHttpRequest("GET", fmt.Sprintf("http://localhost:8800/tables/%s/objects/%s/events", tableName, objectId), "application/json", "")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
