package query

import (
	"bufio"
	"bytes"
	"io"
)

type StatementsParser struct {
}

func NewStatementsParser() *StatementsParser {
	return &StatementsParser{}
}

func (p *StatementsParser) Parse(r io.Reader) (Statements, error) {
	l := newLexer(bufio.NewReader(r), TSTARTSTATEMENTS)
	yyParse(l)
	return l.statements, l.err
}

func (p *StatementsParser) ParseString(s string) (Statements, error) {
	return p.Parse(bytes.NewBufferString(s))
}
