package query

import (
	"bufio"
	"bytes"
	"io"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader) *Query {
	l := newLexer(bufio.NewReader(r), TSTARTQUERY)
	yyParse(l)
	return l.query
}

func (p *Parser) ParseString(s string) *Query {
	return p.Parse(bytes.NewBufferString(s))
}
