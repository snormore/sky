package parser

import (
	"bufio"
	"bytes"
	"io"

	"github.com/skydb/sky/query/ast"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader) (*ast.Query, error) {
	l := newLexer(bufio.NewReader(r), TSTARTQUERY)
	yyParse(l)
	return l.query, l.err
}

func (p *Parser) ParseString(s string) (*ast.Query, error) {
	return p.Parse(bytes.NewBufferString(s))
}
