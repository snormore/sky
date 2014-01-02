package ast

import (
	"fmt"
	"os"
)

func warn(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

