%{

package query

import (
    "bufio"
    "bytes"
    "io"
)

%}

%union{
    value float64
}

%token  TSELECT

%%

select: TSELECT
;

%%

type Parser struct {
}

func NewParser() *Parser {
    return &Parser{}
}

func (p *Parser) Parse(r io.Reader) {
    yyParse(newLexer(bufio.NewReader(r)))
}

func (p *Parser) ParseString(s string) {
    p.Parse(bytes.NewBufferString(s))
}

