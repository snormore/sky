package query

type BooleanLiteral struct {
	value bool
}

func (l *BooleanLiteral) String() string {
	if l.value {
		return "true"
	} else {
		return "false"
	}
}
