package ast

import (
	"fmt"
)

// TemporalLoop represents a statement that iterates over time.
type TemporalLoop struct {
	Iterator   *VarRef
	Step       int
	Duration   int
	Statements Statements
}

func (l *TemporalLoop) node()      {}
func (l *TemporalLoop) block()     {}
func (l *TemporalLoop) statement() {}

// NewTemporalLoop creates a new TemporalLoop instance.
func NewTemporalLoop() *TemporalLoop {
	return &TemporalLoop{}
}

func (l *TemporalLoop) String() string {
	str := "FOR"
	if l.Iterator != nil {
		str += " " + l.Iterator.String()
	}
	if l.Step > 0 {
		quantity, units := SecondsToTimeSpan(l.Step)
		str += " " + fmt.Sprintf("EVERY %d %s", quantity, units)
	}
	if l.Duration > 0 {
		quantity, units := SecondsToTimeSpan(l.Duration)
		str += " " + fmt.Sprintf("WITHIN %d %s", quantity, units)
	}

	str += "\n"
	str += lineStartRegex.ReplaceAllString(l.Statements.String(), "  ") + "\n"
	str += "END"

	return str
}
