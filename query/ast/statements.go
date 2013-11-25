package ast

import (
	"bytes"
	"fmt"
	"strings"
)

type Statements []Statement

func (s Statements) String() string {
	output := []string{}
	for _, statement := range s {
		output = append(output, statement.String())
	}
	return strings.Join(output, "\n")
}
