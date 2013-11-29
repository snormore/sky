package test

import (
	"testing"

	"github.com/skydb/sky/core"
	"github.com/stretchr/testify/assert"
)

func TestSelectCount(t *testing.T) {
	query := `SELECT count()`
	result, err := runMapper(query, map[string][]*core.Event{
		"foo": []*core.Event{
			testevent("2000-01-01T00:00:00Z", 1, 10),
			testevent("2000-01-01T00:00:02Z", 1, 20),
		}})

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
