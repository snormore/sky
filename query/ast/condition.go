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
	Expression Expression
	Start      int
	End        int
	UOM        string
	Statements Statements
}

func (c *Condition) node()      {}
func (c *Condition) block()     {}
func (c *Condition) statement() {}

// NewCondition returns a new Condition instance.
func NewCondition() *Condition {
	return &Condition{
		Start: 0,
		End:   0,
		UOM:   UnitSteps,
	}
}

// Converts the condition to a string-based representation.
func (c *Condition) String() string {
	str := "WHEN"
	if c.Expression != nil {
		str += " " + c.Expression.String()
	}
	if c.Start != 0 || c.End != 0 || c.UOM != UnitSteps {
		str += fmt.Sprintf(" WITHIN %d .. %d %s", c.Start, c.End, strings.ToUpper(c.UOM))
	}
	str += " THEN\n"
	str += lineStartRegex.ReplaceAllString(c.Statements.String(), "  ") + "\n"
	str += "END"
	return str
}
