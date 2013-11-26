package ast

// BooleanLiteral represents a hardcoded boolean value.
type BooleanLiteral struct {
	Value bool
}

// NewBooleanLiteral creates a new BooleanLiteral instance.
func NewBooleanLiteral(value bool) *BooleanLiteral {
	return &BooleanLiteral{Value: value}
}

func (l *BooleanLiteral) node()       {}
func (l *BooleanLiteral) expression() {}

func (l *BooleanLiteral) String() string {
	if l.Value {
		return "true"
	} else {
		return "false"
	}
}
