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
	op  int
	lhs Expression
	rhs Expression
}

func (l *BinaryExpression) String() string {
	var str string

	switch l.lhs.(type) {
	case *BinaryExpression:
		str += "(" + l.lhs.String() + ")"
	default:
		str += l.lhs.String()
	}

	switch l.op {
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
	default:
		str += " <missing> "
	}

	switch l.rhs.(type) {
	case *BinaryExpression:
		str += "(" + l.rhs.String() + ")"
	default:
		str += l.rhs.String()
	}

	return str
}
