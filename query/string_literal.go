package query

type StringLiteral struct {
	value string
}

func (l *StringLiteral) String() string {
	// TODO: Escape string.
	return "\"" + l.value + "\""
}
