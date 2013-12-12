package test

import (
	"fmt"
	"os"
	"time"

	"github.com/skydb/sky/core"
)

func testevent(timestamp string, args ...interface{}) *core.Event {
	e := &core.Event{Timestamp: musttime(timestamp)}
	e.Data = make(map[int64]interface{})
	for i := 0; i < len(args); i += 2 {
		key := args[i].(int)
		e.Data[int64(key)] = args[i+1]
	}
	return e
}

func musttime(timestamp string) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		panic(err)
	}
	return t
}

func debugln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
