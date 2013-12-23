package ast

const (
	OpEquals = iota
	OpNotEquals
	OpGreaterThan
	OpGreaterThanOrEqualTo
	OpLessThan
	OpLessThanOrEqualTo
	OpPlus
	OpMinus
	OpMultiply
	OpDivide
	OpAnd
	OpOr
)

// opstr is a lookup from operations to their string representation.
var opstr = map[int]string{
	OpEquals:               "==",
	OpNotEquals:            "!=",
	OpLessThan:             "<",
	OpLessThanOrEqualTo:    "<=",
	OpGreaterThan:          ">",
	OpGreaterThanOrEqualTo: ">=",
	OpAnd:      "&&",
	OpOr:       "||",
	OpPlus:     "+",
	OpMinus:    "-",
	OpMultiply: "*",
	OpDivide:   "/",
}

// BinaryExpression represents an operation that is performed between
// two Expressions.
type BinaryExpression struct {
	Op  int
	LHS Expression
	RHS Expression
}

func (e *BinaryExpression) node()       {}
func (e *BinaryExpression) expression() {}

// Creates a new binary expression.
func NewBinaryExpression(op int, lhs Expression, rhs Expression) *BinaryExpression {
	return &BinaryExpression{
		Op:  op,
		LHS: lhs,
		RHS: rhs,
	}
}

// Returns a string representation of the expression operator.
func (e *BinaryExpression) OpString() string {
	return opstr[e.Op]
}

// Returns a string representation of the expression.
func (e *BinaryExpression) String() string {
	var str string

	switch e.LHS.(type) {
	case *BinaryExpression:
		str += "(" + e.LHS.String() + ")"
	default:
		str += e.LHS.String()
	}

	str += " " + e.OpString() + " "

	switch e.RHS.(type) {
	case *BinaryExpression:
		str += "(" + e.RHS.String() + ")"
	default:
		str += e.RHS.String()
	}

	return str
}
