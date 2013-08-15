%{

package query

import (
    "bufio"
    "bytes"
    "io"
)

%}

%union{
    token int
    buf []byte
}

%token TSELECT
%token TIDENT
%token TSEMICOLON

%%

select: TSELECT TIDENT TSEMICOLON
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

