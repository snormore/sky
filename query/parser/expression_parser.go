package parser

import (
	"bufio"
	"bytes"
	"io"

	"github.com/skydb/sky/query/ast"
)

type ExpressionParser struct {
}

func NewExpressionParser() *ExpressionParser {
	return &ExpressionParser{}
}

func (p *ExpressionParser) Parse(r io.Reader) (ast.Expression, error) {
	l := newLexer(bufio.NewReader(r), TSTARTEXPRESSION)
	yyParse(l)
	return l.expression, l.err
}

func (p *ExpressionParser) ParseString(s string) (ast.Expression, error) {
	return p.Parse(bytes.NewBufferString(s))
}
