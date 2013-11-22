package query

import (
	"bufio"
	"bytes"
	"io"
)

type StatementParser struct {
}

func NewStatementParser() *StatementParser {
	return &StatementParser{}
}

func (p *StatementParser) Parse(r io.Reader) (Statement, error) {
	l := newLexer(bufio.NewReader(r), TSTARTSTATEMENT)
	yyParse(l)
	return l.statement, l.err
}

func (p *StatementParser) ParseString(s string) (Statement, error) {
	return p.Parse(bytes.NewBufferString(s))
}
