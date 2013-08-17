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

func (p *StatementParser) Parse(query *Query, r io.Reader) Statement {
	l := newLexer(bufio.NewReader(r), TSTARTSTATEMENT)
	l.query = query
	yyParse(l)
	return l.statement
}

func (p *StatementParser) ParseString(query *Query, s string) Statement {
	return p.Parse(query, bytes.NewBufferString(s))
}
