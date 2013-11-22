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

func (p *Parser) Parse(r io.Reader) (*Query, error) {
	l := newLexer(bufio.NewReader(r), TSTARTQUERY)
	yyParse(l)
	return l.query, l.err
}

func (p *Parser) ParseString(s string) (*Query, error) {
	return p.Parse(bytes.NewBufferString(s))
}
