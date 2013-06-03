package skyd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Converts untyped map to a map[string]interface{} if passed a map.
func ConvertToStringKeys(value interface{}) interface{} {
	if m, ok := value.(map[interface{}]interface{}); ok {
		ret := make(map[string]interface{})
		for k, v := range m {
			ret[fmt.Sprintf("%v", k)] = ConvertToStringKeys(v)
		}
		return ret
	}

	return value
}

// Parse response to map and/or error.
func parseResponse(resp *http.Response) (map[string]interface{}, error) {
	// Parse the response into a map.
	ret := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil && err != io.EOF {
		return nil, err
	}

	// Extract the error if one is returned.
	var retErr error
	if errObj, ok := ret["error"].(map[string]interface{}); ok {
		errorMessage, _ := errObj["message"].(string)
		retErr = errors.New(errorMessage)
		delete(ret, "error")
	}
	
	return ret, retErr
}

// Writes to standard error.
func warn(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", v...)
}
