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
    selection_field *QuerySelectionField
    selection_fields []*QuerySelectionField
}

%token <token> TSELECT
%token <token> TSEMICOLON, TCOMMA, TLPAREN, TRPAREN
%token <str> TIDENT

%type <selection> selection
%type <statement> statement
%type <statements> statements
%type <selection_field> selection_field
%type <selection_fields> selection_fields


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

statement :
    selection  { $$ = QueryStep($1) }
;

selection :
    TSELECT selection_fields
    {
        l := yylex.(*yylexer)
        $$ = NewQuerySelection(l.query)
        $$.Fields = $2
    }
;

selection_fields :
    /* empty */
    {
        $$ = make([]*QuerySelectionField, 0)
    }
|   selection_field
    {
        $$ = make([]*QuerySelectionField, 0)
        $$ = append($$, $1)
    }
|   selection_fields TCOMMA selection_field
    {
        $$ = append($1, $3)
    }
;

selection_field :
    TIDENT TLPAREN TRPAREN
    {
        $$ = NewQuerySelectionField("", $1)
    }
|   TIDENT TLPAREN TIDENT TRPAREN
    {
        $$ = NewQuerySelectionField($3, $1)
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

