package query

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

type BinaryExpression struct {
	queryElementImpl
	op  int
	lhs Expression
	rhs Expression
}

// Creates a new binary expression.
func NewBinaryExpression(op int, lhs Expression, rhs Expression) *BinaryExpression {
	e := &BinaryExpression{op: op}
	e.SetLhs(lhs)
	e.SetRhs(rhs)
	return e
}

// Returns the left-hand side expression.
func (e *BinaryExpression) Lhs() Expression {
	return e.lhs
}

// Sets the left-hand side expression.
func (e *BinaryExpression) SetLhs(expression Expression) {
	if e.lhs != nil {
		e.lhs.SetParent(nil)
	}
	e.lhs = expression
	if e.lhs != nil {
		e.lhs.SetParent(e)
	}
}

// Returns the right-hand side expression.
func (e *BinaryExpression) Rhs() Expression {
	return e.rhs
}

// Sets the right-hand side expression.
func (e *BinaryExpression) SetRhs(expression Expression) {
	if e.rhs != nil {
		e.rhs.SetParent(nil)
	}
	e.rhs = expression
	if e.rhs != nil {
		e.rhs.SetParent(e)
	}
}

// Returns a string representation of the expression.
func (e *BinaryExpression) String() string {
	var str string

	switch e.lhs.(type) {
	case *BinaryExpression:
		str += "(" + e.lhs.String() + ")"
	default:
		str += e.lhs.String()
	}

	switch e.op {
	case OpEquals:
		str += " == "
	case OpNotEquals:
		str += " != "
	case OpLessThan:
		str += " < "
	case OpLessThanOrEqualTo:
		str += " <= "
	case OpGreaterThan:
		str += " > "
	case OpGreaterThanOrEqualTo:
		str += " >= "
	case OpAnd:
		str += " && "
	case OpOr:
		str += " || "
	case OpPlus:
		str += " + "
	case OpMinus:
		str += " - "
	case OpMultiply:
		str += " * "
	case OpDivide:
		str += " / "
	default:
		str += " <missing> "
	}

	switch e.rhs.(type) {
	case *BinaryExpression:
		str += "(" + e.rhs.String() + ")"
	default:
		str += e.rhs.String()
	}

	return str
}
