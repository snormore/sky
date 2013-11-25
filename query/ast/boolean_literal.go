package ast

// BooleanLiteral represents a hardcoded boolean value.
type BooleanLiteral struct {
	Value bool
}

func (l *BooleanLiteral) node() string {}

func (l *BooleanLiteral) String() string {
	if l.Value {
		return "true"
	} else {
		return "false"
	}
}
