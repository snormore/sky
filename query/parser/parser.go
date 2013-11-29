package parser

import (
	"bufio"
	"bytes"
	"io"

	"github.com/skydb/sky/query/ast"
)

type Parser struct {
}

// New creates a new Parser instance.
func New() *Parser {
	return &Parser{}
}

// Parse parses a SkyQL query string from a reader and returns the AST structure.
func Parse(r io.Reader) (*ast.Query, error) {
	return New().Parse(r)
}

// ParseString parses a SkyQL query string and returns the AST structure.
func ParseString(s string) (*ast.Query, error) {
	return New().Parse(bytes.NewBufferString(s))
}

// Parse parses a SkyQL query string from a reader and returns the AST structure.
func (p *Parser) Parse(r io.Reader) (*ast.Query, error) {
	l := newLexer(bufio.NewReader(r), TSTARTQUERY)
	yyParse(l)
	return l.query, l.err
}

// ParseString parses a SkyQL query string and returns the AST structure.
func (p *Parser) ParseString(s string) (*ast.Query, error) {
	return p.Parse(bytes.NewBufferString(s))
}
