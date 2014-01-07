package server

import (
	"fmt"
	"os"
	"time"
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

func bench(s string, t ...time.Time) time.Time {
	if len(t) == 0 {
		warnf("%s ====\n", s)
	} else {
		warnf("%s: %v\n", s, time.Now().Sub(t[0]))
	}
	return time.Now()
}

func warn(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func warnf(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, v...)
}
