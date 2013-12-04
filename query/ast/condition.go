package ast

import (
	"fmt"
	"strings"
)

const (
	UnitSteps    = "steps"
	UnitSessions = "sessions"
	UnitSeconds  = "seconds"
)

// Condition represents a conditional statement made within a query.
type Condition struct {
	Expression       Expression
	WithinRangeStart int
	WithinRangeEnd   int
	WithinUnits      string
	Statements       Statements
}

func (c *Condition) node() {}
func (c *Condition) block() {}
func (c *Condition) statement() {}

// NewCondition returns a new Condition instance.
func NewCondition() *Condition {
	return &Condition{
		WithinRangeStart: 0,
		WithinRangeEnd:   0,
		WithinUnits:      UnitSteps,
	}
}

// Converts the condition to a string-based representation.
func (c *Condition) String() string {
	str := "WHEN"
	if str != "" {
		str += " " + c.Expression.String()
	}
	if c.WithinRangeStart != 0 || c.WithinRangeStart != 0 || c.WithinUnits != UnitSteps {
		str += fmt.Sprintf(" WITHIN %d .. %d %s", c.WithinRangeStart, c.WithinRangeEnd, strings.ToUpper(c.WithinUnits))
	}
	str += " THEN\n"
	str += lineStartRegex.ReplaceAllString(c.Statements.String(), "  ") + "\n"
	str += "END"
	return str
}
