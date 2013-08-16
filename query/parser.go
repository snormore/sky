//line parser.y:2
package query

import __yyfmt__ "fmt"

//line parser.y:3
import (
	"bufio"
	"bytes"
	"io"
)

//line parser.y:13
type yySymType struct {
	yys              int
	token            int
	integer          int
	str              string
	strs             []string
	query            *Query
	statement        Statement
	statements       Statements
	selection        *Selection
	selection_field  *SelectionField
	selection_fields []*SelectionField
	condition        *Condition
	condition_within *within
	expr             Expression
	var_ref          *VarRef
	integer_literal  *IntegerLiteral
	string_literal   *StringLiteral
}

const TSELECT = 57346
const TGROUP = 57347
const TBY = 57348
const TINTO = 57349
const TWHEN = 57350
const TWITHIN = 57351
const TTHEN = 57352
const TEND = 57353
const TSEMICOLON = 57354
const TCOMMA = 57355
const TLPAREN = 57356
const TRPAREN = 57357
const TRANGE = 57358
const TEQUALS = 57359
const TNOTEQUALS = 57360
const TLT = 57361
const TLTE = 57362
const TGT = 57363
const TGTE = 57364
const TAND = 57365
const TOR = 57366
const TIDENT = 57367
const TSTRING = 57368
const TWITHINUNITS = 57369
const TINT = 57370

var yyToknames = []string{
	"TSELECT",
	"TGROUP",
	"TBY",
	"TINTO",
	"TWHEN",
	"TWITHIN",
	"TTHEN",
	"TEND",
	"TSEMICOLON",
	"TCOMMA",
	"TLPAREN",
	"TRPAREN",
	"TRANGE",
	"TEQUALS",
	"TNOTEQUALS",
	"TLT",
	"TLTE",
	"TGT",
	"TGTE",
	"TAND",
	"TOR",
	"TIDENT",
	"TSTRING",
	"TWITHINUNITS",
	"TINT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:226

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
	end   int
	units string
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 36
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 91

var yyAct = []int{

	2, 11, 50, 60, 24, 25, 26, 27, 28, 29,
	30, 31, 49, 62, 52, 38, 61, 33, 24, 25,
	26, 27, 28, 29, 15, 39, 41, 42, 43, 44,
	45, 46, 47, 48, 32, 16, 18, 54, 17, 10,
	57, 56, 24, 25, 26, 27, 28, 29, 30, 31,
	24, 25, 26, 27, 28, 29, 30, 26, 27, 28,
	29, 28, 29, 9, 55, 22, 21, 58, 6, 51,
	40, 37, 7, 6, 20, 59, 1, 7, 35, 12,
	14, 13, 23, 5, 36, 34, 53, 19, 8, 3,
	4,
}
var yyPact = []int{

	-1000, -1000, 69, -1000, -1000, -1000, 14, 10, 61, -1000,
	51, 25, -1000, -1000, -1000, 10, -1000, -1000, -1000, 71,
	14, 65, 0, 60, 10, 10, 10, 10, 10, 10,
	10, 10, -16, -13, 57, -12, -1000, 12, -1000, 49,
	-1000, 38, 38, 40, 40, -1000, -1000, 1, 33, 24,
	-1000, -1000, -1000, 54, -1000, -1000, 64, -25, -9, -1000,
	-14, -1000, -1000,
}
var yyPgo = []int{

	0, 90, 89, 0, 63, 88, 87, 86, 85, 83,
	82, 1, 81, 80, 79, 76,
}
var yyR1 = []int{

	0, 15, 3, 3, 2, 2, 1, 5, 5, 5,
	4, 4, 6, 6, 7, 7, 8, 8, 9, 10,
	10, 11, 11, 11, 11, 11, 11, 11, 11, 11,
	11, 11, 11, 14, 12, 13,
}
var yyR2 = []int{

	0, 1, 0, 2, 1, 1, 5, 0, 1, 3,
	3, 4, 0, 3, 1, 3, 0, 2, 6, 0,
	5, 3, 3, 3, 3, 3, 3, 3, 3, 1,
	1, 1, 3, 1, 1, 1,
}
var yyChk = []int{

	-1000, -15, -3, -2, -1, -9, 4, 8, -5, -4,
	25, -11, -14, -12, -13, 14, 25, 28, 26, -6,
	13, 5, 14, -10, 17, 18, 19, 20, 21, 22,
	23, 24, 9, -11, -8, 7, -4, 6, 15, 25,
	10, -11, -11, -11, -11, -11, -11, -11, -11, 28,
	15, 12, 26, -7, 25, 15, -3, 16, 13, 11,
	28, 25, 27,
}
var yyDef = []int{

	2, -2, 1, 3, 4, 5, 7, 0, 12, 8,
	0, 19, 29, 30, 31, 0, 33, 34, 35, 16,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 9, 0, 10, 0,
	2, 21, 22, 23, 24, 25, 26, 27, 28, 0,
	32, 6, 17, 13, 14, 11, 0, 0, 0, 18,
	0, 15, 20,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line parser.y:66
		{
			l := yylex.(*yylexer)
			l.query.Statements = yyS[yypt-0].statements
		}
	case 2:
		//line parser.y:74
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 3:
		//line parser.y:78
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 4:
		//line parser.y:85
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 5:
		//line parser.y:89
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 6:
		//line parser.y:96
		{
			l := yylex.(*yylexer)
			yyVAL.selection = NewSelection(l.query)
			yyVAL.selection.Fields = yyS[yypt-3].selection_fields
			yyVAL.selection.Dimensions = yyS[yypt-2].strs
			yyVAL.selection.Name = yyS[yypt-1].str
		}
	case 7:
		//line parser.y:107
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 8:
		//line parser.y:111
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 9:
		//line parser.y:116
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 10:
		//line parser.y:123
		{
			yyVAL.selection_field = NewSelectionField("", yyS[yypt-2].str)
		}
	case 11:
		//line parser.y:127
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-1].str, yyS[yypt-3].str)
		}
	case 12:
		//line parser.y:134
		{
			yyVAL.strs = make([]string, 0)
		}
	case 13:
		//line parser.y:138
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 14:
		//line parser.y:145
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 15:
		//line parser.y:150
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 16:
		//line parser.y:157
		{
			yyVAL.str = ""
		}
	case 17:
		//line parser.y:161
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 18:
		//line parser.y:168
		{
			l := yylex.(*yylexer)
			yyVAL.condition = NewCondition(l.query)
			yyVAL.condition.Expression = yyS[yypt-4].expr.String()
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.Statements = yyS[yypt-1].statements
		}
	case 19:
		//line parser.y:181
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: "steps"}
		}
	case 20:
		//line parser.y:185
		{
			yyVAL.condition_within = &within{start: yyS[yypt-3].integer, end: yyS[yypt-1].integer, units: yyS[yypt-0].str}
		}
	case 21:
		//line parser.y:191
		{
			yyVAL.expr = &BinaryExpression{op: OpEquals, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 22:
		//line parser.y:192
		{
			yyVAL.expr = &BinaryExpression{op: OpNotEquals, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 23:
		//line parser.y:193
		{
			yyVAL.expr = &BinaryExpression{op: OpLessThan, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 24:
		//line parser.y:194
		{
			yyVAL.expr = &BinaryExpression{op: OpLessThanOrEqualTo, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 25:
		//line parser.y:195
		{
			yyVAL.expr = &BinaryExpression{op: OpGreaterThan, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 26:
		//line parser.y:196
		{
			yyVAL.expr = &BinaryExpression{op: OpGreaterThanOrEqualTo, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 27:
		//line parser.y:197
		{
			yyVAL.expr = &BinaryExpression{op: OpAnd, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 28:
		//line parser.y:198
		{
			yyVAL.expr = &BinaryExpression{op: OpOr, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 29:
		//line parser.y:199
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 30:
		//line parser.y:200
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 31:
		//line parser.y:201
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 32:
		//line parser.y:202
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 33:
		//line parser.y:207
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 34:
		//line parser.y:214
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 35:
		//line parser.y:221
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
