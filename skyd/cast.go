package skyd

import (
	"fmt"
	"reflect"
)

// Casts a value to an int64.
func castInt64(value interface{}) (int64, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return int64(value.(int)), nil
	case reflect.Int8:
		return int64(value.(int8)), nil
	case reflect.Int16:
		return int64(value.(int16)), nil
	case reflect.Int32:
		return int64(value.(int32)), nil
	case reflect.Int64:
		return value.(int64), nil
	case reflect.Uint:
		return int64(value.(uint)), nil
	case reflect.Uint8:
		return int64(value.(uint8)), nil
	case reflect.Uint16:
		return int64(value.(uint16)), nil
	case reflect.Uint32:
		return int64(value.(uint32)), nil
	case reflect.Uint64:
		return int64(value.(uint64)), nil
	}
	return 0, fmt.Errorf("skyd.castInt64: Unable to cast '%v'", value)
}

// Casts a value to a uint64.
func castUint64(value interface{}) (uint64, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return uint64(value.(int)), nil
	case reflect.Int8:
		return uint64(value.(int8)), nil
	case reflect.Int16:
		return uint64(value.(int16)), nil
	case reflect.Int32:
		return uint64(value.(int32)), nil
	case reflect.Int64:
		return uint64(value.(int64)), nil
	case reflect.Uint:
		return uint64(value.(uint)), nil
	case reflect.Uint8:
		return uint64(value.(uint8)), nil
	case reflect.Uint16:
		return uint64(value.(uint16)), nil
	case reflect.Uint32:
		return uint64(value.(uint32)), nil
	case reflect.Uint64:
		return value.(uint64), nil
	}
	return 0, fmt.Errorf("skyd.castUint64: Unable to cast '%v'", value)
}
