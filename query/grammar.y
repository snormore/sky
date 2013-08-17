%{

package query

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
    expr Expression
    var_ref *VarRef
    integer_literal *IntegerLiteral
    boolean_literal *BooleanLiteral
    string_literal *StringLiteral
}

%token <token> TSTARTQUERY, TSTARTSTATEMENT, TSTARTEXPRESSION
%token <token> TSELECT, TGROUP, TBY, TINTO
%token <token> TWHEN, TWITHIN, TTHEN, TEND
%token <token> TSEMICOLON, TCOMMA, TLPAREN, TRPAREN, TRANGE
%token <token> TEQUALS, TNOTEQUALS, TLT, TLTE, TGT, TGTE
%token <token> TAND, TOR, TPLUS, TMINUS, TMUL, TDIV
%token <token> TTRUE, TFALSE
%token <str> TIDENT, TSTRING, TWITHINUNITS
%token <integer> TINT

%type <query> query
%type <selection> selection
%type <statement> statement
%type <statements> statements
%type <selection_field> selection_field
%type <selection_fields> selection_fields
%type <strs> selection_group_by, selection_dimensions
%type <str> selection_name

%type <condition> condition
%type <condition_within> condition_within

%type <expr> expr
%type <integer_literal> integer_literal
%type <boolean_literal> boolean_literal
%type <string_literal> string_literal
%type <var_ref> var_ref

%left TOR
%left TAND
%left TEQUALS TNOTEQUALS
%left TLT TLTE
%left TGT TGTE
%left TPLUS TMINUS
%left TMUL TDIV

%start start

%%

start :
    TSTARTQUERY query
    {
        l := yylex.(*yylexer)
        l.query = $2
    }
|   TSTARTSTATEMENT statement
    {
        l := yylex.(*yylexer)
        l.statement = $2
    }
|   TSTARTEXPRESSION expr
    {
        l := yylex.(*yylexer)
        l.expression = $2
    }
;

query :
    statements
    {
        l := yylex.(*yylexer)
        l.query.Statements = $1
        $$ = l.query
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
    TWHEN expr condition_within TTHEN statements TEND
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

expr :
    expr TEQUALS expr     { $$ = &BinaryExpression{op:OpEquals, lhs:$1, rhs:$3} }
|   expr TNOTEQUALS expr  { $$ = &BinaryExpression{op:OpNotEquals, lhs:$1, rhs:$3} }
|   expr TLT expr         { $$ = &BinaryExpression{op:OpLessThan, lhs:$1, rhs:$3} }
|   expr TLTE expr        { $$ = &BinaryExpression{op:OpLessThanOrEqualTo, lhs:$1, rhs:$3} }
|   expr TGT expr         { $$ = &BinaryExpression{op:OpGreaterThan, lhs:$1, rhs:$3} }
|   expr TGTE expr        { $$ = &BinaryExpression{op:OpGreaterThanOrEqualTo, lhs:$1, rhs:$3} }
|   expr TAND expr        { $$ = &BinaryExpression{op:OpAnd, lhs:$1, rhs:$3} }
|   expr TOR expr         { $$ = &BinaryExpression{op:OpOr, lhs:$1, rhs:$3} }
|   expr TPLUS expr       { $$ = &BinaryExpression{op:OpPlus, lhs:$1, rhs:$3} }
|   expr TMINUS expr      { $$ = &BinaryExpression{op:OpMinus, lhs:$1, rhs:$3} }
|   expr TMUL expr        { $$ = &BinaryExpression{op:OpMultiply, lhs:$1, rhs:$3} }
|   expr TDIV expr        { $$ = &BinaryExpression{op:OpDivide, lhs:$1, rhs:$3} }
|   var_ref               { $$ = Expression($1) }
|   integer_literal       { $$ = Expression($1) }
|   boolean_literal       { $$ = Expression($1) }
|   string_literal        { $$ = Expression($1) }
|   TLPAREN expr TRPAREN { $$ = $2 }
;

var_ref :
    TIDENT
    {
        $$ = &VarRef{value:$1}
    }
;

integer_literal :
    TINT
    {
        $$ = &IntegerLiteral{value:$1}
    }
;

boolean_literal :
    TTRUE
    {
        $$ = &BooleanLiteral{value:true}
    }
|   TFALSE
    {
        $$ = &BooleanLiteral{value:false}
    }
;

string_literal :
    TSTRING
    {
        $$ = &StringLiteral{value:$1}
    }
;

%%

type within struct {
    start int
    end int
    units string
}
