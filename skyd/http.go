package skyd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Sends an JSON encoded object to a given URL and parses the JSON response
// into an object.
func rpc(host string, port uint, method string, path string, reqObj interface{}, respObj interface{}) error {
	url := fmt.Sprintf("http://%s:%d%s", host, port, path)

	// Encode the object.
	var b bytes.Buffer
	if reqObj != nil {
		json.NewEncoder(&b).Encode(reqObj)
	}

	// Setup the request.
	req, _ := http.NewRequest(method, url, &b)
	req.Header.Set("Content-Type", "application/json")

	// Send request over client.
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Return an error if we don't receive a 200.
	if resp.StatusCode != http.StatusOK {
		message := resp.Header.Get("Sky-Error-Message")
		if message == "" {
			return fmt.Errorf("skyd: HTTP error: %d (%s %s%s)", resp.StatusCode, method, url, path)
		} else {
			return errors.New(message)
		}
	}

	// Decode response into object.
	if respObj != nil {
		if err := json.NewDecoder(resp.Body).Decode(&respObj); err != nil && err != io.EOF {
			return fmt.Errorf("skyd: Deserialization error: %v", err)
		}
	}

	return nil
}
