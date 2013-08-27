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
	event_loop       *EventLoop
	temporal_loop    *TemporalLoop
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
const TFOR = 57366
const TEACH = 57367
const TEVERY = 57368
const TIN = 57369
const TEVENT = 57370
const TSEMICOLON = 57371
const TCOMMA = 57372
const TLPAREN = 57373
const TRPAREN = 57374
const TRANGE = 57375
const TEQUALS = 57376
const TNOTEQUALS = 57377
const TLT = 57378
const TLTE = 57379
const TGT = 57380
const TGTE = 57381
const TAND = 57382
const TOR = 57383
const TPLUS = 57384
const TMINUS = 57385
const TMUL = 57386
const TDIV = 57387
const TASSIGN = 57388
const TTRUE = 57389
const TFALSE = 57390
const TIDENT = 57391
const TQUOTEDSTRING = 57392
const TWITHINUNITS = 57393
const TTIMEUNITS = 57394
const TINT = 57395

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
	"TFOR",
	"TEACH",
	"TEVERY",
	"TIN",
	"TEVENT",
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
	"TTIMEUNITS",
	"TINT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line grammar.y:375
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

const yyNprod = 65
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 215

var yyAct = []int{

	8, 19, 36, 113, 107, 24, 84, 91, 30, 87,
	115, 108, 24, 118, 98, 117, 39, 116, 110, 38,
	100, 27, 28, 25, 29, 37, 53, 26, 27, 28,
	25, 29, 45, 46, 26, 54, 49, 50, 51, 52,
	25, 25, 55, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 51, 52, 79, 104, 59,
	82, 85, 102, 88, 41, 42, 43, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 43, 44, 45, 46,
	109, 62, 49, 50, 51, 52, 64, 103, 86, 77,
	106, 41, 42, 43, 44, 45, 46, 47, 48, 49,
	50, 51, 52, 61, 41, 42, 43, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 41, 42, 43, 44,
	45, 46, 47, 48, 49, 50, 51, 52, 41, 42,
	43, 44, 45, 46, 47, 90, 49, 50, 51, 52,
	41, 42, 43, 44, 45, 46, 58, 81, 49, 50,
	51, 52, 49, 50, 51, 52, 15, 16, 83, 57,
	111, 17, 15, 16, 114, 18, 20, 17, 15, 16,
	112, 18, 101, 17, 15, 16, 105, 18, 78, 17,
	32, 1, 34, 18, 23, 40, 93, 94, 95, 96,
	97, 2, 4, 3, 5, 33, 22, 21, 89, 63,
	9, 14, 13, 60, 12, 92, 80, 99, 56, 35,
	11, 10, 31, 7, 6,
}
var yyPact = []int{

	187, -1000, -1000, -1000, 159, -19, -1000, 167, 159, -1000,
	-1000, -1000, -1000, -1000, -1000, -8, -24, -19, -9, 70,
	-1000, -1000, -1000, -1000, -19, -1000, -1000, -1000, -1000, -1000,
	159, -1000, -14, -1000, -4, 129, -1000, 28, 82, 53,
	60, -19, -19, -19, -19, -19, -19, -19, -19, -19,
	-19, -19, -19, 57, 164, -19, 128, -24, 140, -26,
	66, -44, -1000, 114, -46, 40, 40, -6, -6, 110,
	110, 106, 94, 11, 11, -1000, -1000, -1000, 178, 70,
	-1000, -36, -1000, -29, 158, 30, -1000, 25, 153, -1000,
	-49, -41, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 50,
	-1000, -31, 146, 147, -50, -1000, 141, -42, -1000, -32,
	-1000, -34, -1000, -38, -1000, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 214, 213, 212, 211, 210, 195, 0, 2, 209,
	208, 207, 206, 205, 204, 203, 202, 201, 199, 198,
	1, 197, 196, 184, 166, 181,
}
var yyR1 = []int{

	0, 25, 25, 25, 25, 1, 7, 7, 6, 6,
	6, 6, 6, 2, 2, 3, 13, 13, 13, 13,
	13, 4, 5, 9, 9, 9, 8, 8, 10, 10,
	11, 11, 12, 12, 14, 15, 15, 16, 17, 18,
	18, 19, 19, 20, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20, 20, 20,
	24, 21, 22, 22, 23,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 1, 1, 0, 2, 4, 1, 1, 1, 1,
	1, 4, 4, 0, 1, 3, 5, 6, 0, 3,
	1, 3, 0, 2, 6, 0, 5, 5, 6, 0,
	3, 0, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 1, 1, 1, 1, 3,
	1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -25, 4, 6, 5, 7, -1, -2, -7, -6,
	-4, -5, -14, -16, -17, 15, 16, 20, 24, -20,
	-24, -21, -22, -23, 31, 49, 53, 47, 48, 50,
	-7, -3, 13, -6, -24, -9, -8, 49, -20, 25,
	-24, 34, 35, 36, 37, 38, 39, 40, 41, 42,
	43, 44, 45, -20, 49, 46, -10, 30, 17, 31,
	-15, 21, 28, -18, 26, -20, -20, -20, -20, -20,
	-20, -20, -20, -20, -20, -20, -20, 32, 14, -20,
	-12, 19, -8, 18, 32, -20, 22, 53, -7, -19,
	21, 53, -13, 8, 9, 10, 11, 12, 50, -11,
	49, 14, 32, -7, 33, 23, -7, 53, 52, 30,
	49, 14, 23, 53, 23, 52, 49, 49, 51,
}
var yyDef = []int{

	0, -2, 13, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 11, 12, 0, 23, 0, 0, 4,
	55, 56, 57, 58, 0, 60, 61, 62, 63, 64,
	5, 14, 0, 7, 0, 28, 24, 0, 35, 0,
	39, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 32, 0, 0, 0,
	0, 0, 6, 41, 0, 43, 44, 45, 46, 47,
	48, 49, 50, 51, 52, 53, 54, 59, 0, 21,
	22, 0, 25, 0, 0, 0, 6, 0, 0, 6,
	0, 0, 15, 16, 17, 18, 19, 20, 33, 29,
	30, 0, 0, 0, 0, 37, 0, 0, 40, 0,
	26, 0, 34, 0, 38, 42, 31, 27, 36,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53,
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
		//line grammar.y:90
		{
			l := yylex.(*yylexer)
			l.query = yyS[yypt-0].query
		}
	case 2:
		//line grammar.y:95
		{
			l := yylex.(*yylexer)
			l.statements = yyS[yypt-0].statements
		}
	case 3:
		//line grammar.y:100
		{
			l := yylex.(*yylexer)
			l.statement = yyS[yypt-0].statement
		}
	case 4:
		//line grammar.y:105
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 5:
		//line grammar.y:113
		{
			yyVAL.query = &Query{}
			yyVAL.query.SetDeclaredVariables(yyS[yypt-1].variables)
			yyVAL.query.SetStatements(yyS[yypt-0].statements)
		}
	case 6:
		//line grammar.y:122
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 7:
		//line grammar.y:126
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 8:
		//line grammar.y:132
		{
			yyVAL.statement = Statement(yyS[yypt-0].assignment)
		}
	case 9:
		//line grammar.y:133
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 10:
		//line grammar.y:134
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 11:
		//line grammar.y:135
		{
			yyVAL.statement = Statement(yyS[yypt-0].event_loop)
		}
	case 12:
		//line grammar.y:136
		{
			yyVAL.statement = Statement(yyS[yypt-0].temporal_loop)
		}
	case 13:
		//line grammar.y:141
		{
			yyVAL.variables = make([]*Variable, 0)
		}
	case 14:
		//line grammar.y:145
		{
			yyVAL.variables = append(yyS[yypt-1].variables, yyS[yypt-0].variable)
		}
	case 15:
		//line grammar.y:152
		{
			yyVAL.variable = NewVariable(yyS[yypt-2].str, yyS[yypt-0].str)
		}
	case 16:
		//line grammar.y:158
		{
			yyVAL.str = core.FactorDataType
		}
	case 17:
		//line grammar.y:159
		{
			yyVAL.str = core.StringDataType
		}
	case 18:
		//line grammar.y:160
		{
			yyVAL.str = core.IntegerDataType
		}
	case 19:
		//line grammar.y:161
		{
			yyVAL.str = core.FloatDataType
		}
	case 20:
		//line grammar.y:162
		{
			yyVAL.str = core.BooleanDataType
		}
	case 21:
		//line grammar.y:167
		{
			yyVAL.assignment = NewAssignment()
			yyVAL.assignment.SetTarget(yyS[yypt-2].var_ref)
			yyVAL.assignment.SetExpression(yyS[yypt-0].expr)
		}
	case 22:
		//line grammar.y:176
		{
			yyVAL.selection = NewSelection()
			yyVAL.selection.SetFields(yyS[yypt-2].selection_fields)
			yyVAL.selection.Dimensions = yyS[yypt-1].strs
			yyVAL.selection.Name = yyS[yypt-0].str
		}
	case 23:
		//line grammar.y:186
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 24:
		//line grammar.y:190
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 25:
		//line grammar.y:195
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 26:
		//line grammar.y:202
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str, nil)
		}
	case 27:
		//line grammar.y:206
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-5].str, yyS[yypt-3].expr)
		}
	case 28:
		//line grammar.y:213
		{
			yyVAL.strs = make([]string, 0)
		}
	case 29:
		//line grammar.y:217
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 30:
		//line grammar.y:224
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 31:
		//line grammar.y:229
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 32:
		//line grammar.y:236
		{
			yyVAL.str = ""
		}
	case 33:
		//line grammar.y:240
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 34:
		//line grammar.y:247
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 35:
		//line grammar.y:259
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 36:
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
	case 37:
		//line grammar.y:278
		{
			yyVAL.event_loop = NewEventLoop()
			yyVAL.event_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 38:
		//line grammar.y:286
		{
			yyVAL.temporal_loop = NewTemporalLoop()
			yyVAL.temporal_loop.SetRef(yyS[yypt-4].var_ref)
			yyVAL.temporal_loop.step = yyS[yypt-3].integer
			yyVAL.temporal_loop.duration = yyS[yypt-2].integer
			yyVAL.temporal_loop.SetStatements(yyS[yypt-1].statements)

			// Default steps to 1 of the unit of the duration.
			if yyVAL.temporal_loop.step == 0 && yyVAL.temporal_loop.duration > 0 {
				_, units := secondsToTimeSpan(yyS[yypt-2].integer)
				yyVAL.temporal_loop.step = timeSpanToSeconds(1, units)
			}
		}
	case 39:
		//line grammar.y:303
		{
			yyVAL.integer = 0
		}
	case 40:
		//line grammar.y:307
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 41:
		//line grammar.y:314
		{
			yyVAL.integer = 0
		}
	case 42:
		//line grammar.y:318
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 43:
		//line grammar.y:324
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 44:
		//line grammar.y:325
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 45:
		//line grammar.y:326
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 46:
		//line grammar.y:327
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 47:
		//line grammar.y:328
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 48:
		//line grammar.y:329
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 49:
		//line grammar.y:330
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 50:
		//line grammar.y:331
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 51:
		//line grammar.y:332
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 52:
		//line grammar.y:333
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 53:
		//line grammar.y:334
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 54:
		//line grammar.y:335
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 55:
		//line grammar.y:336
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 56:
		//line grammar.y:337
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 57:
		//line grammar.y:338
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 58:
		//line grammar.y:339
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 59:
		//line grammar.y:340
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 60:
		//line grammar.y:345
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 61:
		//line grammar.y:352
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 62:
		//line grammar.y:359
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 63:
		//line grammar.y:363
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 64:
		//line grammar.y:370
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
