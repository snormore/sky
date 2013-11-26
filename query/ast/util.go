package ast

import (
	"fmt"
	"reflect"
	"regexp"
)

const (
	secondsPerYear   = 60 * 60 * 24 * 7 * 365
	secondsPerWeek   = 60 * 60 * 24 * 7
	secondsPerDay    = 60 * 60 * 24
	secondsPerHour   = 60 * 60
	secondsPerMinute = 60
)

var lineStartRegex = regexp.MustCompile(`(?m)^`)
var nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z]+`)

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

// Converts a human readable time representation to seconds.
func TimeSpanToSeconds(quantity int, units string) int {
	switch units {
	case "YEAR", "YEARS":
		return quantity * secondsPerYear
	case "WEEK", "WEEKS":
		return quantity * secondsPerWeek
	case "DAY", "DAYS":
		return quantity * secondsPerDay
	case "HOUR", "HOURS":
		return quantity * secondsPerHour
	case "MINUTE", "MINUTES":
		return quantity * secondsPerMinute
	case "SECOND", "SECONDS":
		return quantity
	}
	panic(fmt.Sprintf("Duration not supported: %s", units))
}

// Converts a number of seconds a human readable representation.
func SecondsToTimeSpan(seconds int) (int, string) {
	var quantity int
	var units string
	if seconds >= secondsPerYear && seconds%secondsPerYear == 0 {
		quantity, units = seconds/secondsPerYear, "YEAR"
	} else if seconds >= secondsPerWeek && seconds%secondsPerWeek == 0 {
		quantity, units = seconds/secondsPerWeek, "WEEK"
	} else if seconds >= secondsPerDay && seconds%secondsPerDay == 0 {
		quantity, units = seconds/secondsPerDay, "DAY"
	} else if seconds >= secondsPerHour && seconds%secondsPerHour == 0 {
		quantity, units = seconds/secondsPerHour, "HOUR"
	} else if seconds >= secondsPerMinute && seconds%secondsPerMinute == 0 {
		quantity, units = seconds/secondsPerMinute, "MINUTE"
	} else {
		quantity, units = seconds, "SECOND"
	}

	if quantity != 1 {
		units += "S"
	}

	return quantity, units
}
