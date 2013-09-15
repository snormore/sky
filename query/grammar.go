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
	exit             *Exit
	debug            *Debug
	selection        *Selection
	selection_field  *SelectionField
	selection_fields []*SelectionField
	condition        *Condition
	condition_within *within
	event_loop       *EventLoop
	session_loop     *SessionLoop
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
const TEXIT = 57358
const TDEBUG = 57359
const TSELECT = 57360
const TGROUP = 57361
const TBY = 57362
const TINTO = 57363
const TWHEN = 57364
const TWITHIN = 57365
const TTHEN = 57366
const TEND = 57367
const TFOR = 57368
const TEACH = 57369
const TEVERY = 57370
const TIN = 57371
const TEVENT = 57372
const TSESSION = 57373
const TDELIMITED = 57374
const TSEMICOLON = 57375
const TCOMMA = 57376
const TLPAREN = 57377
const TRPAREN = 57378
const TRANGE = 57379
const TEQUALS = 57380
const TNOTEQUALS = 57381
const TLT = 57382
const TLTE = 57383
const TGT = 57384
const TGTE = 57385
const TAND = 57386
const TOR = 57387
const TPLUS = 57388
const TMINUS = 57389
const TMUL = 57390
const TDIV = 57391
const TASSIGN = 57392
const TTRUE = 57393
const TFALSE = 57394
const TAMPERSAND = 57395
const TIDENT = 57396
const TAMPIDENT = 57397
const TQUOTEDSTRING = 57398
const TWITHINUNITS = 57399
const TTIMEUNITS = 57400
const TINT = 57401

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
	"TEXIT",
	"TDEBUG",
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
	"TSESSION",
	"TDELIMITED",
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
	"TAMPERSAND",
	"TIDENT",
	"TAMPIDENT",
	"TQUOTEDSTRING",
	"TWITHINUNITS",
	"TTIMEUNITS",
	"TINT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line grammar.y:441
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

const yyNprod = 78
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 286

var yyAct = []int{

	8, 30, 44, 69, 45, 29, 99, 132, 37, 131,
	24, 124, 107, 102, 138, 134, 29, 125, 137, 35,
	36, 32, 33, 48, 35, 36, 34, 47, 114, 31,
	98, 120, 32, 33, 62, 46, 36, 34, 60, 61,
	31, 63, 64, 58, 59, 60, 61, 65, 139, 127,
	35, 36, 71, 68, 42, 78, 79, 80, 81, 82,
	83, 84, 85, 86, 87, 88, 89, 29, 67, 92,
	96, 52, 53, 54, 55, 103, 100, 58, 59, 60,
	61, 104, 128, 32, 33, 77, 35, 36, 34, 54,
	55, 31, 101, 58, 59, 60, 61, 74, 75, 116,
	106, 95, 119, 117, 122, 118, 123, 50, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 60, 61, 97,
	70, 91, 129, 50, 51, 52, 53, 54, 55, 135,
	136, 58, 59, 60, 61, 39, 1, 28, 93, 140,
	50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	60, 61, 90, 73, 50, 51, 52, 53, 54, 55,
	56, 57, 58, 59, 60, 61, 70, 27, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	109, 110, 111, 112, 113, 2, 4, 3, 5, 26,
	50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	60, 61, 50, 51, 52, 53, 54, 55, 56, 57,
	58, 59, 60, 61, 50, 51, 52, 53, 54, 55,
	56, 105, 58, 59, 60, 61, 18, 19, 20, 21,
	76, 17, 16, 22, 15, 25, 141, 23, 18, 19,
	20, 21, 40, 72, 14, 22, 126, 9, 133, 23,
	18, 19, 20, 21, 41, 108, 94, 22, 115, 49,
	130, 23, 18, 19, 20, 21, 66, 43, 13, 22,
	12, 11, 121, 23, 18, 19, 20, 21, 10, 38,
	7, 22, 6, 0, 0, 23,
}
var yyPact = []int{

	181, -1000, -1000, -1000, 259, 32, -1000, 122, 259, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -35, -1000,
	19, -19, 32, -4, 164, -1000, -1000, -1000, -1000, 32,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 259, -1000, -35,
	-1000, -8, 32, 34, -1000, 152, 17, 130, 67, 57,
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 116, 107, 32, 102, 80, -19, 99, -1000,
	-24, -30, 68, -46, -1000, 49, 77, -47, 31, 31,
	47, 47, -3, -3, 85, 176, -10, -10, -1000, -1000,
	-1000, 172, 164, -1000, -1000, -28, -1000, -35, -1000, 106,
	69, -1000, -6, 247, 84, -1000, -48, -41, 14, -1000,
	-1000, -1000, -1000, -1000, -1000, 48, -1000, -1000, 106, 235,
	-50, -1000, -52, 223, -43, -1000, -1000, -35, -35, -1000,
	-1000, -39, -44, -1000, -1000, 12, -1000, -1000, -1000, -1000,
	211, -1000,
}
var yyPgo = []int{

	0, 282, 280, 279, 278, 271, 270, 268, 242, 0,
	2, 267, 266, 258, 1, 256, 3, 255, 246, 244,
	243, 234, 232, 231, 230, 221, 4, 189, 167, 137,
	235, 136,
}
var yyR1 = []int{

	0, 31, 31, 31, 31, 1, 9, 9, 8, 8,
	8, 8, 8, 8, 8, 8, 2, 2, 3, 18,
	18, 17, 17, 17, 17, 17, 4, 5, 6, 7,
	11, 11, 11, 10, 10, 10, 16, 16, 12, 12,
	13, 13, 15, 15, 19, 20, 20, 21, 22, 23,
	24, 24, 25, 25, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 30, 14, 14, 27, 28, 28, 29,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 0, 2, 5, 0,
	3, 1, 1, 1, 1, 1, 4, 1, 4, 4,
	0, 1, 3, 2, 4, 5, 0, 2, 0, 3,
	1, 3, 0, 2, 6, 0, 5, 5, 9, 6,
	0, 3, 0, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 1, 1, 1, 1,
	3, 1, 1, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -31, 4, 6, 5, 7, -1, -2, -9, -8,
	-4, -5, -6, -7, -19, -21, -22, -23, 15, 16,
	17, 18, 22, 26, -26, -30, -27, -28, -29, 35,
	-14, 59, 51, 52, 56, 54, 55, -9, -3, 13,
	-8, -30, 35, -11, -10, -26, 54, -26, 27, -30,
	38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
	48, 49, -26, -14, 50, -26, -12, 34, 19, -16,
	14, 35, -20, 23, 30, 31, -24, 28, -26, -26,
	-26, -26, -26, -26, -26, -26, -26, -26, -26, -26,
	36, 14, -26, 36, -15, 21, -10, 20, 54, 36,
	-26, 24, 59, -9, 32, -25, 23, 59, -17, 8,
	9, 10, 11, 12, 56, -13, -14, -16, 36, -9,
	37, 25, 20, -9, 59, 58, -18, 35, 34, -16,
	25, 59, 59, 25, 58, -14, -14, 57, 58, 36,
	-9, 25,
}
var yyDef = []int{

	0, -2, 16, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 11, 12, 13, 14, 15, 0, 27,
	0, 30, 0, 0, 4, 66, 67, 68, 69, 0,
	71, 74, 75, 76, 77, 72, 73, 5, 17, 0,
	7, 0, 0, 38, 31, 36, 72, 45, 0, 50,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 42, 0, 0, 33,
	0, 0, 0, 0, 6, 0, 52, 0, 54, 55,
	56, 57, 58, 59, 60, 61, 62, 63, 64, 65,
	70, 0, 26, 28, 29, 0, 32, 0, 37, 36,
	0, 6, 0, 0, 0, 6, 0, 0, 19, 21,
	22, 23, 24, 25, 43, 39, 40, 34, 36, 0,
	0, 47, 0, 0, 0, 51, 18, 0, 0, 35,
	44, 0, 0, 49, 53, 0, 41, 46, 6, 20,
	0, 48,
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
	52, 53, 54, 55, 56, 57, 58, 59,
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
		//line grammar.y:97
		{
			l := yylex.(*yylexer)
			l.query = yyS[yypt-0].query
		}
	case 2:
		//line grammar.y:102
		{
			l := yylex.(*yylexer)
			l.statements = yyS[yypt-0].statements
		}
	case 3:
		//line grammar.y:107
		{
			l := yylex.(*yylexer)
			l.statement = yyS[yypt-0].statement
		}
	case 4:
		//line grammar.y:112
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 5:
		//line grammar.y:120
		{
			yyVAL.query = NewQuery()
			yyVAL.query.SetDeclaredVariables(yyS[yypt-1].variables)
			yyVAL.query.SetStatements(yyS[yypt-0].statements)
		}
	case 6:
		//line grammar.y:129
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 7:
		//line grammar.y:133
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 8:
		//line grammar.y:139
		{
			yyVAL.statement = Statement(yyS[yypt-0].assignment)
		}
	case 9:
		//line grammar.y:140
		{
			yyVAL.statement = Statement(yyS[yypt-0].exit)
		}
	case 10:
		//line grammar.y:141
		{
			yyVAL.statement = Statement(yyS[yypt-0].debug)
		}
	case 11:
		//line grammar.y:142
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 12:
		//line grammar.y:143
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 13:
		//line grammar.y:144
		{
			yyVAL.statement = Statement(yyS[yypt-0].event_loop)
		}
	case 14:
		//line grammar.y:145
		{
			yyVAL.statement = Statement(yyS[yypt-0].session_loop)
		}
	case 15:
		//line grammar.y:146
		{
			yyVAL.statement = Statement(yyS[yypt-0].temporal_loop)
		}
	case 16:
		//line grammar.y:151
		{
			yyVAL.variables = make([]*Variable, 0)
		}
	case 17:
		//line grammar.y:155
		{
			yyVAL.variables = append(yyS[yypt-1].variables, yyS[yypt-0].variable)
		}
	case 18:
		//line grammar.y:162
		{
			yyVAL.variable = NewVariable(yyS[yypt-3].str, yyS[yypt-1].str)
			yyVAL.variable.Association = yyS[yypt-0].str
		}
	case 19:
		//line grammar.y:170
		{
			yyVAL.str = ""
		}
	case 20:
		//line grammar.y:174
		{
			yyVAL.str = yyS[yypt-1].str
		}
	case 21:
		//line grammar.y:180
		{
			yyVAL.str = core.FactorDataType
		}
	case 22:
		//line grammar.y:181
		{
			yyVAL.str = core.StringDataType
		}
	case 23:
		//line grammar.y:182
		{
			yyVAL.str = core.IntegerDataType
		}
	case 24:
		//line grammar.y:183
		{
			yyVAL.str = core.FloatDataType
		}
	case 25:
		//line grammar.y:184
		{
			yyVAL.str = core.BooleanDataType
		}
	case 26:
		//line grammar.y:189
		{
			yyVAL.assignment = NewAssignment()
			yyVAL.assignment.SetTarget(yyS[yypt-2].var_ref)
			yyVAL.assignment.SetExpression(yyS[yypt-0].expr)
		}
	case 27:
		//line grammar.y:198
		{
			yyVAL.exit = NewExit()
		}
	case 28:
		//line grammar.y:205
		{
			yyVAL.debug = NewDebug()
			yyVAL.debug.SetExpression(yyS[yypt-1].expr)
		}
	case 29:
		//line grammar.y:213
		{
			yyVAL.selection = NewSelection()
			yyVAL.selection.SetFields(yyS[yypt-2].selection_fields)
			yyVAL.selection.Dimensions = yyS[yypt-1].strs
			yyVAL.selection.Name = yyS[yypt-0].str
		}
	case 30:
		//line grammar.y:223
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 31:
		//line grammar.y:227
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 32:
		//line grammar.y:232
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 33:
		//line grammar.y:239
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, "", yyS[yypt-1].expr)
		}
	case 34:
		//line grammar.y:243
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-3].str, nil)
		}
	case 35:
		//line grammar.y:247
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str, yyS[yypt-2].expr)
		}
	case 36:
		//line grammar.y:253
		{
			yyVAL.str = ""
		}
	case 37:
		//line grammar.y:254
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 38:
		//line grammar.y:259
		{
			yyVAL.strs = make([]string, 0)
		}
	case 39:
		//line grammar.y:263
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 40:
		//line grammar.y:270
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 41:
		//line grammar.y:275
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 42:
		//line grammar.y:282
		{
			yyVAL.str = ""
		}
	case 43:
		//line grammar.y:286
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 44:
		//line grammar.y:293
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 45:
		//line grammar.y:305
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 46:
		//line grammar.y:309
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
	case 47:
		//line grammar.y:324
		{
			yyVAL.event_loop = NewEventLoop()
			yyVAL.event_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 48:
		//line grammar.y:332
		{
			yyVAL.session_loop = NewSessionLoop()
			yyVAL.session_loop.IdleDuration = timeSpanToSeconds(yyS[yypt-3].integer, yyS[yypt-2].str)
			yyVAL.session_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 49:
		//line grammar.y:341
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
	case 50:
		//line grammar.y:358
		{
			yyVAL.integer = 0
		}
	case 51:
		//line grammar.y:362
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 52:
		//line grammar.y:369
		{
			yyVAL.integer = 0
		}
	case 53:
		//line grammar.y:373
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 54:
		//line grammar.y:379
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 55:
		//line grammar.y:380
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 56:
		//line grammar.y:381
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 57:
		//line grammar.y:382
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 58:
		//line grammar.y:383
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 59:
		//line grammar.y:384
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 60:
		//line grammar.y:385
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 61:
		//line grammar.y:386
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 62:
		//line grammar.y:387
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 63:
		//line grammar.y:388
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 64:
		//line grammar.y:389
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 65:
		//line grammar.y:390
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 66:
		//line grammar.y:391
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 67:
		//line grammar.y:392
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 68:
		//line grammar.y:393
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 69:
		//line grammar.y:394
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 70:
		//line grammar.y:395
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 71:
		//line grammar.y:400
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 72:
		//line grammar.y:407
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 73:
		//line grammar.y:411
		{
			yyVAL.str = yyS[yypt-0].str[1:]
		}
	case 74:
		//line grammar.y:418
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 75:
		//line grammar.y:425
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 76:
		//line grammar.y:429
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 77:
		//line grammar.y:436
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
