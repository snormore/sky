package parser

import (
	"bufio"
	"bytes"
	"io"

	"github.com/skydb/sky/query/ast"
)

type StatementParser struct {
}

func NewStatementParser() *StatementParser {
	return &StatementParser{}
}

func (p *StatementParser) Parse(r io.Reader) (ast.Statement, error) {
	l := newLexer(bufio.NewReader(r), TSTARTSTATEMENT)
	yyParse(l)
	return l.statement, l.err
}

func (p *StatementParser) ParseString(s string) (ast.Statement, error) {
	return p.Parse(bytes.NewBufferString(s))
}
