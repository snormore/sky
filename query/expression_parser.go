package query

import (
	"bufio"
	"bytes"
	"io"
)

type ExpressionParser struct {
}

func NewExpressionParser() *ExpressionParser {
	return &ExpressionParser{}
}

func (p *ExpressionParser) Parse(query *Query, r io.Reader) Expression {
	l := newLexer(bufio.NewReader(r), TSTARTEXPRESSION)
	l.query = query
	yyParse(l)
	return l.expression
}

func (p *ExpressionParser) ParseString(query *Query, s string) Expression {
	return p.Parse(query, bytes.NewBufferString(s))
}
