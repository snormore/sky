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

func (p *ExpressionParser) Parse(r io.Reader) (Expression, error) {
	l := newLexer(bufio.NewReader(r), TSTARTEXPRESSION)
	yyParse(l)
	return l.expression, l.err
}

func (p *ExpressionParser) ParseString(s string) (Expression, error) {
	return p.Parse(bytes.NewBufferString(s))
}
