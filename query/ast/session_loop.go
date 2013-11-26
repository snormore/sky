package ast

import (
	"fmt"
)

// SessionLoop represents a statement that iterates over individual sessions.
type SessionLoop struct {
	Statements   Statements
	IdleDuration int
}

func (l *SessionLoop) node() {}

// NewSessionLoop creates a new SessionLoop instance.
func NewSessionLoop() *SessionLoop {
	return &SessionLoop{}
}

func (l *SessionLoop) String() string {
	quantity, units := SecondsToTimeSpan(l.IdleDuration)
	str := fmt.Sprintf("FOR EACH SESSION DELIMITED BY %d %s\n", quantity, units)
	str += lineStartRegex.ReplaceAllString(l.Statements.String(), "  ") + "\n"
	str += "END"
	return str
}
