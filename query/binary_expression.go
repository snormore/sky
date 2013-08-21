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

const (
	queryMode = "query"
	luaMode   = "lua"
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

// Returns a Lua representation of the expression.
func (e *BinaryExpression) Codegen() (string, error) {
	var str string

	// Codegen the left-hand side.
	expr, err := e.lhs.Codegen()
	if err != nil {
		return "", err
	}

	switch e.lhs.(type) {
	case *BinaryExpression:
		str += "(" + expr + ")"
	default:
		str += expr
	}

	str += " " + e.stringifyOperator(luaMode) + " "

	// Codegen the right-hand side.
	expr, err = e.rhs.Codegen()
	if err != nil {
		return "", err
	}

	switch e.rhs.(type) {
	case *BinaryExpression:
		str += "(" + expr + ")"
	default:
		str += expr
	}

	return str, nil
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

	str += " " + e.stringifyOperator(queryMode) + " "

	switch e.rhs.(type) {
	case *BinaryExpression:
		str += "(" + e.rhs.String() + ")"
	default:
		str += e.rhs.String()
	}

	return str
}

// Converts the binary expression's operator to a query or Lua string.
func (e *BinaryExpression) stringifyOperator(mode string) string {
	switch e.op {
	case OpEquals:
		return "=="
	case OpNotEquals:
		if mode == luaMode {
			return "~="
		} else {
			return "!="
		}
	case OpLessThan:
		return "<"
	case OpLessThanOrEqualTo:
		return "<="
	case OpGreaterThan:
		return ">"
	case OpGreaterThanOrEqualTo:
		return ">="
	case OpAnd:
		if mode == luaMode {
			return "and"
		} else {
			return "&&"
		}
	case OpOr:
		if mode == luaMode {
			return "or"
		} else {
			return "||"
		}
	case OpPlus:
		return "+"
	case OpMinus:
		return "-"
	case OpMultiply:
		return "*"
	case OpDivide:
		return "/"
	}

	return "<missing>"
}
