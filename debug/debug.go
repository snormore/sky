package sky

import (
	"fmt"
	"os"
)

// Prints to the standard error.
// Arguments are handled in the manner of fmt.Print.
func Warn(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
}

// Prints to the standard logger if debug mode is enabled.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

// Prints to the standard logger if debug mode is enabled.
// Arguments are handled in the manner of fmt.Println.
func Warnln(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

