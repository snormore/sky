package parser

import (
	"bufio"
	"bytes"
	"io"

	"github.com/skydb/sky/query/ast"
)

type StatementsParser struct {
}

func NewStatementsParser() *StatementsParser {
	return &StatementsParser{}
}

func (p *StatementsParser) Parse(r io.Reader) (ast.Statements, error) {
	l := newLexer(bufio.NewReader(r), TSTARTSTATEMENTS)
	yyParse(l)
	return l.statements, l.err
}

func (p *StatementsParser) ParseString(s string) (ast.Statements, error) {
	return p.Parse(bytes.NewBufferString(s))
}
