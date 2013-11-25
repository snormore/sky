package ast

import (
	"bytes"
	"fmt"
	"github.com/skydb/sky/core"
)

// TemporalLoop represents a statement that iterates over time.
type TemporalLoop struct {
	Ref        *VarRef
	Step       int
	Duration   int
	Statements Statements
}

func (l *TemporalLoop) node() string {}

// NewTemporalLoop creates a new TemporalLoop instance.
func NewTemporalLoop() *TemporalLoop {
	return &TemporalLoop{}
}

func (l *TemporalLoop) String() string {
	str := "FOR"
	if l.ref != nil {
		str += " " + l.ref.String()
	}
	if l.step > 0 {
		quantity, units := secondsToTimeSpan(l.step)
		str += " " + fmt.Sprintf("EVERY %d %s", quantity, units)
	}
	if l.duration > 0 {
		quantity, units := secondsToTimeSpan(l.duration)
		str += " " + fmt.Sprintf("WITHIN %d %s", quantity, units)
	}

	str += "\n"
	str += lineStartRegex.ReplaceAllString(l.statements.String(), "  ") + "\n"
	str += "END"

	return str
}
