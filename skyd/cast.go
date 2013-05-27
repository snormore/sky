package skyd

import (
	"reflect"
)

// Normalizes a value. Int and Uint types are combined into int64 and Float types
// are combined into float64. All other types are left alone.
func normalize(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return value
}

// Casts to a uint64 if possible.
func castUint64(value interface{}) (uint64, bool) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return uint64(v.Float()), true
	}
	return 0, false
}
