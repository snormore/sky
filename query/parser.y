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
    integer int
    str string
    strs []string
    query *Query
    statement Statement
    statements Statements
    selection *Selection
    selection_field *SelectionField
    selection_fields []*SelectionField
    condition *Condition
    condition_within *within
}

%token <token> TSELECT, TGROUP, TBY, TINTO
%token <token> TWHEN, TWITHIN, TTHEN, TEND
%token <token> TSEMICOLON, TCOMMA, TLPAREN, TRPAREN, TRANGE
%token <token> TEQUALS, TAND, TOR
%token <str> TIDENT, TSTRING, TWITHINUNITS
%token <integer> TINT

%type <selection> selection
%type <statement> statement
%type <statements> statements
%type <selection_field> selection_field
%type <selection_fields> selection_fields
%type <strs> selection_group_by, selection_dimensions
%type <str> selection_name

%type <condition> condition
%type <str> conditionals conditional
%type <condition_within> condition_within

%%

query :
    statements
    {
        l := yylex.(*yylexer)
        l.query.Statements = $1
    }
;

statements :
    /* empty */
    {
        $$ = make(Statements, 0)
    }
|   statements statement
    {
        $$ = append($1, $2)
    }
;

statement :
    selection
    {
        $$ = Statement($1)
    }
|   condition
    {
        $$ = Statement($1)
    }
;

selection :
    TSELECT selection_fields selection_group_by selection_name TSEMICOLON
    {
        l := yylex.(*yylexer)
        $$ = NewSelection(l.query)
        $$.Fields = $2
        $$.Dimensions = $3
        $$.Name = $4
    }
;

selection_fields :
    /* empty */
    {
        $$ = make([]*SelectionField, 0)
    }
|   selection_field
    {
        $$ = make([]*SelectionField, 0)
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
        $$ = NewSelectionField("", $1)
    }
|   TIDENT TLPAREN TIDENT TRPAREN
    {
        $$ = NewSelectionField($3, $1)
    }
;

selection_group_by :
    /* empty */
    {
        $$ = make([]string, 0)
    }
|   TGROUP TBY selection_dimensions
    {
        $$ = $3
    }
;

selection_dimensions :
    TIDENT
    {
        $$ = make([]string, 0)
        $$ = append($$, $1)
    }
|   selection_dimensions TCOMMA TIDENT
    {
        $$ = append($1, $3)
    }
;

selection_name :
    /* empty */
    {
        $$ = ""
    }
|   TINTO TSTRING
    {
        $$ = $2
    }
;

condition :
    TWHEN conditionals condition_within TTHEN statements TEND
    {
        l := yylex.(*yylexer)
        $$ = NewCondition(l.query)
        $$.Expression = $2
        $$.WithinRangeStart = $3.start
        $$.WithinRangeEnd = $3.end
        $$.WithinUnits = $3.units
        $$.Statements = $5
    }
;

conditionals :
    /* empty */
    {
        $$ = ""
    }
|   conditional
    {
        $$ = $1
    }
|   conditionals TAND conditional
    {
        $$ = $1 + " && " + $3
    }
;

conditional :
    TIDENT TEQUALS TSTRING
    {
        $$ = $1 + " == \"" + $3 + "\""
    }
;

condition_within :
    /* empty */
    {
        $$ = &within{start:0, end:0, units:"steps"}
    }
|   TWITHIN TINT TRANGE TINT TWITHINUNITS
    {
        $$ = &within{start:$2, end:$4, units:$5}
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

type within struct {
    start int
    end int
    units string
}
