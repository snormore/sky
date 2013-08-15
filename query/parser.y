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
    str string
    query *Query
    statement QueryStep
    statements QueryStepList
    selection *QuerySelection
}

%token <token> TSELECT
%token <token> TSEMICOLON
%token <str> TIDENT

%type <statement> select
%type <statement> statement
%type <statements> statements


%%

query :
    statements
    {
        l := yylex.(*yylexer)
        l.query.Steps = $1
    }
;

statements :
/* empty */
    {
        $$ = make(QueryStepList, 0)
    }
|   statement TSEMICOLON
    {
        $$ = make(QueryStepList, 0)
        $$ = append($$, $1)
    }
|   statements TSEMICOLON statement
    {
        $$ = append($1, $3)
    }
;

statement : select
;

select :
    TSELECT TIDENT
    {
        l := yylex.(*yylexer)
        s := NewQuerySelection(l.query)
        s.Name = $2
        $$ = s
    }
;

%%

type Parser struct {
}

func NewParser() *Parser {
    return &Parser{}
}

func (p *Parser) Parse(r io.Reader) *Query {
    l := newLexer(bufio.NewReader(r))
    yyParse(l)
    return l.query
}

func (p *Parser) ParseString(s string) *Query {
    return p.Parse(bytes.NewBufferString(s))
}

