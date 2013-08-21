//line grammar.y:2
package query

import __yyfmt__ "fmt"

//line grammar.y:3
import (
	"github.com/skydb/sky/core"
)

//line grammar.y:11
type yySymType struct {
	yys              int
	token            int
	integer          int
	str              string
	strs             []string
	query            *Query
	variable         *Variable
	variables        []*Variable
	statement        Statement
	statements       Statements
	assignment       *Assignment
	selection        *Selection
	selection_field  *SelectionField
	selection_fields []*SelectionField
	condition        *Condition
	condition_within *within
	expr             Expression
	var_ref          *VarRef
	integer_literal  *IntegerLiteral
	boolean_literal  *BooleanLiteral
	string_literal   *StringLiteral
}

const TSTARTQUERY = 57346
const TSTARTSTATEMENT = 57347
const TSTARTSTATEMENTS = 57348
const TSTARTEXPRESSION = 57349
const TFACTOR = 57350
const TSTRING = 57351
const TINTEGER = 57352
const TFLOAT = 57353
const TBOOLEAN = 57354
const TDECLARE = 57355
const TAS = 57356
const TSET = 57357
const TSELECT = 57358
const TGROUP = 57359
const TBY = 57360
const TINTO = 57361
const TWHEN = 57362
const TWITHIN = 57363
const TTHEN = 57364
const TEND = 57365
const TSEMICOLON = 57366
const TCOMMA = 57367
const TLPAREN = 57368
const TRPAREN = 57369
const TRANGE = 57370
const TEQUALS = 57371
const TNOTEQUALS = 57372
const TLT = 57373
const TLTE = 57374
const TGT = 57375
const TGTE = 57376
const TAND = 57377
const TOR = 57378
const TPLUS = 57379
const TMINUS = 57380
const TMUL = 57381
const TDIV = 57382
const TASSIGN = 57383
const TTRUE = 57384
const TFALSE = 57385
const TIDENT = 57386
const TQUOTEDSTRING = 57387
const TWITHINUNITS = 57388
const TINT = 57389

var yyToknames = []string{
	"TSTARTQUERY",
	"TSTARTSTATEMENT",
	"TSTARTSTATEMENTS",
	"TSTARTEXPRESSION",
	"TFACTOR",
	"TSTRING",
	"TINTEGER",
	"TFLOAT",
	"TBOOLEAN",
	"TDECLARE",
	"TAS",
	"TSET",
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
	"TPLUS",
	"TMINUS",
	"TMUL",
	"TDIV",
	"TASSIGN",
	"TTRUE",
	"TFALSE",
	"TIDENT",
	"TQUOTEDSTRING",
	"TWITHINUNITS",
	"TINT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line grammar.y:328
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

const yyNprod = 57
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 170

var yyAct = []int{

	8, 16, 33, 21, 98, 79, 101, 87, 27, 40,
	41, 100, 99, 44, 45, 46, 47, 35, 56, 24,
	25, 22, 26, 48, 23, 95, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 57, 58,
	59, 60, 61, 62, 63, 64, 65, 66, 67, 68,
	89, 34, 71, 49, 69, 74, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 38, 39,
	40, 41, 22, 50, 44, 45, 46, 47, 93, 92,
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45,
	46, 47, 36, 37, 38, 39, 40, 41, 42, 76,
	44, 45, 46, 47, 36, 37, 38, 39, 40, 41,
	46, 47, 44, 45, 46, 47, 77, 44, 45, 46,
	47, 91, 13, 14, 54, 53, 94, 15, 86, 73,
	97, 78, 75, 52, 13, 14, 96, 90, 70, 15,
	81, 82, 83, 84, 85, 17, 29, 2, 4, 3,
	5, 30, 1, 20, 19, 18, 9, 55, 12, 31,
	80, 72, 88, 51, 32, 11, 10, 28, 7, 6,
}
var yyPact = []int{

	143, -1000, -1000, -1000, 119, -23, -1000, 133, 119, -1000,
	-1000, -1000, -1000, 28, 7, -23, 51, -1000, -1000, -1000,
	-1000, -23, -1000, -1000, -1000, -1000, -1000, 119, -1000, 9,
	-1000, 32, 108, -1000, 98, -3, -23, -23, -23, -23,
	-23, -23, -23, -23, -23, -23, -23, -23, 27, 124,
	-23, 110, 7, 114, 72, 109, -42, 37, 37, -24,
	-24, 80, 80, 75, 63, 71, 71, -1000, -1000, -1000,
	132, 51, 104, -38, -1000, 6, 123, 94, -1000, 50,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 101, -1000,
	-19, 122, 107, -43, -32, -1000, -33, -1000, -40, -1000,
	-1000, -1000,
}
var yyPgo = []int{

	0, 169, 168, 167, 166, 165, 151, 0, 2, 164,
	163, 162, 161, 160, 158, 157, 1, 155, 154, 153,
	145, 152,
}
var yyR1 = []int{

	0, 21, 21, 21, 21, 1, 7, 7, 6, 6,
	6, 2, 2, 3, 13, 13, 13, 13, 13, 4,
	5, 9, 9, 9, 8, 8, 10, 10, 11, 11,
	12, 12, 14, 15, 15, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 20, 17, 18, 18, 19,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 0, 2, 4, 1, 1, 1, 1, 1, 4,
	5, 0, 1, 3, 5, 6, 0, 3, 1, 3,
	0, 2, 6, 0, 5, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 1, 1, 1,
	1, 3, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -21, 4, 6, 5, 7, -1, -2, -7, -6,
	-4, -5, -14, 15, 16, 20, -16, -20, -17, -18,
	-19, 26, 44, 47, 42, 43, 45, -7, -3, 13,
	-6, -20, -9, -8, 44, -16, 29, 30, 31, 32,
	33, 34, 35, 36, 37, 38, 39, 40, -16, 44,
	41, -10, 25, 17, 26, -15, 21, -16, -16, -16,
	-16, -16, -16, -16, -16, -16, -16, -16, -16, 27,
	14, -16, -12, 19, -8, 18, 27, 44, 22, 47,
	-13, 8, 9, 10, 11, 12, 24, 45, -11, 44,
	14, 27, -7, 28, 25, 44, 14, 23, 47, 44,
	44, 46,
}
var yyDef = []int{

	0, -2, 11, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 0, 21, 0, 4, 47, 48, 49,
	50, 0, 52, 53, 54, 55, 56, 5, 12, 0,
	7, 0, 26, 22, 0, 33, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 30, 0, 0, 0, 0, 0, 35, 36, 37,
	38, 39, 40, 41, 42, 43, 44, 45, 46, 51,
	0, 19, 0, 0, 23, 0, 0, 0, 6, 0,
	13, 14, 15, 16, 17, 18, 20, 31, 27, 28,
	0, 0, 0, 0, 0, 24, 0, 32, 0, 29,
	25, 34,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47,
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
		//line grammar.y:83
		{
			l := yylex.(*yylexer)
			l.query = yyS[yypt-0].query
		}
	case 2:
		//line grammar.y:88
		{
			l := yylex.(*yylexer)
			l.statements = yyS[yypt-0].statements
		}
	case 3:
		//line grammar.y:93
		{
			l := yylex.(*yylexer)
			l.statement = yyS[yypt-0].statement
		}
	case 4:
		//line grammar.y:98
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 5:
		//line grammar.y:106
		{
			yyVAL.query = &Query{}
			yyVAL.query.SetDeclaredVariables(yyS[yypt-1].variables)
			yyVAL.query.SetStatements(yyS[yypt-0].statements)
		}
	case 6:
		//line grammar.y:115
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 7:
		//line grammar.y:119
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 8:
		//line grammar.y:126
		{
			yyVAL.statement = Statement(yyS[yypt-0].assignment)
		}
	case 9:
		//line grammar.y:130
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 10:
		//line grammar.y:134
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 11:
		//line grammar.y:141
		{
			yyVAL.variables = make([]*Variable, 0)
		}
	case 12:
		//line grammar.y:145
		{
			yyVAL.variables = append(yyS[yypt-1].variables, yyS[yypt-0].variable)
		}
	case 13:
		//line grammar.y:152
		{
			yyVAL.variable = NewVariable(yyS[yypt-2].str, yyS[yypt-0].str)
		}
	case 14:
		//line grammar.y:158
		{
			yyVAL.str = core.FactorDataType
		}
	case 15:
		//line grammar.y:159
		{
			yyVAL.str = core.StringDataType
		}
	case 16:
		//line grammar.y:160
		{
			yyVAL.str = core.IntegerDataType
		}
	case 17:
		//line grammar.y:161
		{
			yyVAL.str = core.FloatDataType
		}
	case 18:
		//line grammar.y:162
		{
			yyVAL.str = core.BooleanDataType
		}
	case 19:
		//line grammar.y:167
		{
			yyVAL.assignment = NewAssignment()
			yyVAL.assignment.SetTarget(yyS[yypt-2].var_ref)
			yyVAL.assignment.SetExpression(yyS[yypt-0].expr)
		}
	case 20:
		//line grammar.y:176
		{
			yyVAL.selection = NewSelection()
			yyVAL.selection.SetFields(yyS[yypt-3].selection_fields)
			yyVAL.selection.Dimensions = yyS[yypt-2].strs
			yyVAL.selection.Name = yyS[yypt-1].str
		}
	case 21:
		//line grammar.y:186
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 22:
		//line grammar.y:190
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 23:
		//line grammar.y:195
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 24:
		//line grammar.y:202
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str+"()")
		}
	case 25:
		//line grammar.y:206
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-5].str+"("+yyS[yypt-3].str+")")
		}
	case 26:
		//line grammar.y:213
		{
			yyVAL.strs = make([]string, 0)
		}
	case 27:
		//line grammar.y:217
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 28:
		//line grammar.y:224
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 29:
		//line grammar.y:229
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 30:
		//line grammar.y:236
		{
			yyVAL.str = ""
		}
	case 31:
		//line grammar.y:240
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 32:
		//line grammar.y:247
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 33:
		//line grammar.y:259
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 34:
		//line grammar.y:263
		{
			yyVAL.condition_within = &within{start: yyS[yypt-3].integer, end: yyS[yypt-1].integer}
			switch yyS[yypt-0].str {
			case "STEPS":
				yyVAL.condition_within.units = UnitSteps
			case "SESSIONS":
				yyVAL.condition_within.units = UnitSessions
			case "SECONDS":
				yyVAL.condition_within.units = UnitSeconds
			}
		}
	case 35:
		//line grammar.y:277
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 36:
		//line grammar.y:278
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 37:
		//line grammar.y:279
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 38:
		//line grammar.y:280
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 39:
		//line grammar.y:281
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 40:
		//line grammar.y:282
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 41:
		//line grammar.y:283
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 42:
		//line grammar.y:284
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 43:
		//line grammar.y:285
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 44:
		//line grammar.y:286
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 45:
		//line grammar.y:287
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 46:
		//line grammar.y:288
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 47:
		//line grammar.y:289
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 48:
		//line grammar.y:290
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 49:
		//line grammar.y:291
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 50:
		//line grammar.y:292
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 51:
		//line grammar.y:293
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 52:
		//line grammar.y:298
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 53:
		//line grammar.y:305
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 54:
		//line grammar.y:312
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 55:
		//line grammar.y:316
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 56:
		//line grammar.y:323
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
