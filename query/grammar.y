%{

package query

import (
    "github.com/skydb/sky/core"
)

%}

%union{
    token int
    integer int
    str string
    strs []string
    query *Query
    variable *Variable
    variables []*Variable
    statement Statement
    statements Statements
    assignment *Assignment
    selection *Selection
    selection_field *SelectionField
    selection_fields []*SelectionField
    condition *Condition
    condition_within *within
    event_loop *EventLoop
    temporal_loop *TemporalLoop
    expr Expression
    var_ref *VarRef
    integer_literal *IntegerLiteral
    boolean_literal *BooleanLiteral
    string_literal *StringLiteral
}

%token <token> TSTARTQUERY, TSTARTSTATEMENT, TSTARTSTATEMENTS, TSTARTEXPRESSION
%token <token> TFACTOR, TSTRING, TINTEGER, TFLOAT, TBOOLEAN
%token <token> TDECLARE, TAS, TSET
%token <token> TSELECT, TGROUP, TBY, TINTO
%token <token> TWHEN, TWITHIN, TTHEN, TEND
%token <token> TFOR, TEACH, TEVERY, TIN, TEVENT
%token <token> TSEMICOLON, TCOMMA, TLPAREN, TRPAREN, TRANGE
%token <token> TEQUALS, TNOTEQUALS, TLT, TLTE, TGT, TGTE
%token <token> TAND, TOR, TPLUS, TMINUS, TMUL, TDIV, TASSIGN
%token <token> TTRUE, TFALSE
%token <str> TIDENT, TQUOTEDSTRING, TWITHINUNITS, TTIMEUNITS
%token <integer> TINT

%type <query> query
%type <variables> variables
%type <variable> variable
%type <assignment> assignment
%type <selection> selection
%type <statement> statement
%type <statements> statements
%type <selection_field> selection_field
%type <selection_fields> selection_fields
%type <strs> selection_group_by, selection_dimensions
%type <str> selection_name
%type <str> data_type

%type <condition> condition
%type <condition_within> condition_within

%type <event_loop> event_loop
%type <temporal_loop> temporal_loop
%type <integer> temporal_loop_step temporal_loop_duration

%type <expr> expr
%type <integer_literal> integer_literal
%type <boolean_literal> boolean_literal
%type <string_literal> string_literal
%type <var_ref> var_ref

%left TASSIGN
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
|   TSTARTSTATEMENTS statements
    {
        l := yylex.(*yylexer)
        l.statements = $2
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
    variables statements
    {
        $$ = &Query{}
        $$.SetDeclaredVariables($1)
        $$.SetStatements($2)
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
    assignment    { $$ = Statement($1) }
|   selection     { $$ = Statement($1) }
|   condition     { $$ = Statement($1) }
|   event_loop    { $$ = Statement($1) }
|   temporal_loop { $$ = Statement($1) }
;

variables :
    /* empty */
    {
        $$ = make([]*Variable, 0)
    }
|   variables variable
    {
        $$ = append($1, $2)
    }
;

variable :
    TDECLARE TIDENT TAS data_type
    {
        $$ = NewVariable($2, $4)
    }
;

data_type :
    TFACTOR  { $$ = core.FactorDataType }
|   TSTRING  { $$ = core.StringDataType }
|   TINTEGER { $$ = core.IntegerDataType }
|   TFLOAT   { $$ = core.FloatDataType }
|   TBOOLEAN { $$ = core.BooleanDataType }
;

assignment :
    TSET var_ref TASSIGN expr
    {
        $$ = NewAssignment()
        $$.SetTarget($2)
        $$.SetExpression($4)
    }
;

selection :
    TSELECT selection_fields selection_group_by selection_name
    {
        $$ = NewSelection()
        $$.SetFields($2)
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
    TIDENT TLPAREN TRPAREN TAS TIDENT
    {
        $$ = NewSelectionField($5, $1, nil)
    }
|   TIDENT TLPAREN expr TRPAREN TAS TIDENT
    {
        $$ = NewSelectionField($6, $1, $3)
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
|   TINTO TQUOTEDSTRING
    {
        $$ = $2
    }
;

condition :
    TWHEN expr condition_within TTHEN statements TEND
    {
        $$ = NewCondition()
        $$.SetExpression($2)
        $$.WithinRangeStart = $3.start
        $$.WithinRangeEnd = $3.end
        $$.WithinUnits = $3.units
        $$.SetStatements($5)
    }
;

condition_within :
    /* empty */
    {
        $$ = &within{start:0, end:0, units:UnitSteps}
    }
|   TWITHIN TINT TRANGE TINT TWITHINUNITS
    {
        $$ = &within{start:$2, end:$4}
        switch $5 {
        case "STEPS":
            $$.units = UnitSteps
        case "SESSIONS":
            $$.units = UnitSessions
        case "SECONDS":
            $$.units = UnitSeconds
        }
    }
;

event_loop :
    TFOR TEACH TEVENT statements TEND
    {
        $$ = NewEventLoop()
        $$.SetStatements($4)
    }
;

temporal_loop :
    TFOR var_ref temporal_loop_step temporal_loop_duration statements TEND
    {
        $$ = NewTemporalLoop()
        $$.SetRef($2)
        $$.step = $3
        $$.duration = $4
        $$.SetStatements($5)

        // Default steps to 1 of the unit of the duration.
        if $$.step == 0 && $$.duration > 0 {
            _, units := secondsToTimeSpan($4)
            $$.step = timeSpanToSeconds(1, units)
        }
    }
;

temporal_loop_step :
    /* empty */
    {
        $$ = 0
    }
|   TEVERY TINT TTIMEUNITS
    {
        $$ = timeSpanToSeconds($2, $3)
    }
;

temporal_loop_duration :
    /* empty */
    {
        $$ = 0
    }
|   TWITHIN TINT TTIMEUNITS
    {
        $$ = timeSpanToSeconds($2, $3)
    }
;

expr :
    expr TEQUALS expr     { $$ = NewBinaryExpression(OpEquals, $1, $3) }
|   expr TNOTEQUALS expr  { $$ = NewBinaryExpression(OpNotEquals, $1, $3) }
|   expr TLT expr         { $$ = NewBinaryExpression(OpLessThan, $1, $3) }
|   expr TLTE expr        { $$ = NewBinaryExpression(OpLessThanOrEqualTo, $1, $3) }
|   expr TGT expr         { $$ = NewBinaryExpression(OpGreaterThan, $1, $3) }
|   expr TGTE expr        { $$ = NewBinaryExpression(OpGreaterThanOrEqualTo, $1, $3) }
|   expr TAND expr        { $$ = NewBinaryExpression(OpAnd, $1, $3) }
|   expr TOR expr         { $$ = NewBinaryExpression(OpOr, $1, $3) }
|   expr TPLUS expr       { $$ = NewBinaryExpression(OpPlus, $1, $3) }
|   expr TMINUS expr      { $$ = NewBinaryExpression(OpMinus, $1, $3) }
|   expr TMUL expr        { $$ = NewBinaryExpression(OpMultiply, $1, $3) }
|   expr TDIV expr        { $$ = NewBinaryExpression(OpDivide, $1, $3) }
|   var_ref               { $$ = Expression($1) }
|   integer_literal       { $$ = Expression($1) }
|   boolean_literal       { $$ = Expression($1) }
|   string_literal        { $$ = Expression($1) }
|   TLPAREN expr TRPAREN  { $$ = $2 }
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
    TQUOTEDSTRING
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
