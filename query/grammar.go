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
const TQUOTEDSTRING = 57397
const TWITHINUNITS = 57398
const TTIMEUNITS = 57399
const TINT = 57400

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
	"TQUOTEDSTRING",
	"TWITHINUNITS",
	"TTIMEUNITS",
	"TINT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line grammar.y:432
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

const yyNprod = 75
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 266

var yyAct = []int{

	8, 30, 24, 44, 29, 96, 130, 129, 37, 121,
	104, 99, 137, 132, 122, 136, 47, 29, 36, 35,
	32, 33, 36, 35, 34, 46, 111, 31, 135, 126,
	45, 62, 61, 32, 33, 36, 35, 34, 59, 60,
	31, 63, 36, 35, 64, 65, 57, 58, 59, 60,
	117, 138, 76, 77, 78, 79, 80, 81, 82, 83,
	84, 85, 86, 87, 124, 69, 42, 90, 125, 101,
	75, 94, 97, 100, 115, 68, 49, 50, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 60, 53, 54,
	67, 98, 57, 58, 59, 60, 103, 113, 93, 116,
	72, 73, 91, 120, 49, 50, 51, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 119, 95, 127, 114,
	49, 50, 51, 52, 53, 54, 133, 134, 57, 58,
	59, 60, 89, 39, 1, 28, 71, 88, 139, 49,
	50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	60, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	58, 59, 60, 49, 50, 51, 52, 53, 54, 55,
	56, 57, 58, 59, 60, 49, 50, 51, 52, 53,
	54, 55, 25, 57, 58, 59, 60, 51, 52, 53,
	54, 27, 26, 57, 58, 59, 60, 18, 19, 20,
	21, 41, 102, 74, 22, 17, 48, 140, 23, 18,
	19, 20, 21, 40, 16, 15, 22, 70, 9, 131,
	23, 18, 19, 20, 21, 14, 123, 105, 22, 92,
	112, 128, 23, 18, 19, 20, 21, 66, 43, 13,
	22, 12, 11, 118, 23, 18, 19, 20, 21, 10,
	38, 7, 22, 6, 0, 0, 23, 106, 107, 108,
	109, 110, 2, 4, 3, 5,
}
var yyPact = []int{

	258, -1000, -1000, -1000, 230, -18, -1000, 120, 230, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -35, -1000,
	31, -24, -18, -11, 125, -1000, -1000, -1000, -1000, -18,
	-1000, -1000, -1000, -1000, -1000, -1000, -23, 230, -1000, -35,
	-1000, -6, -18, 56, -1000, 30, 113, 70, 42, -18,
	-18, -18, -18, -18, -18, -18, -18, -18, -18, -18,
	-18, 101, -1000, 118, -18, 66, 77, -24, 97, -31,
	67, -47, -1000, 37, 73, -48, 147, 147, 46, 46,
	0, 0, 82, 137, -10, -10, -1000, -1000, -1000, 249,
	125, -1000, -1000, -29, -1000, -35, 105, 38, -1000, 13,
	218, 96, -1000, -49, -43, 29, -1000, -1000, -1000, -1000,
	-1000, -1000, 34, -1000, -25, 104, 206, -51, -1000, -52,
	194, -44, -1000, -1000, -35, -35, -1000, -26, -1000, -41,
	-45, -1000, -1000, 15, -1000, -1000, -1000, -1000, -1000, 182,
	-1000,
}
var yyPgo = []int{

	0, 253, 251, 250, 249, 242, 241, 239, 213, 0,
	3, 238, 237, 230, 1, 229, 227, 226, 225, 217,
	215, 214, 205, 203, 202, 2, 192, 191, 135, 182,
	134,
}
var yyR1 = []int{

	0, 30, 30, 30, 30, 1, 9, 9, 8, 8,
	8, 8, 8, 8, 8, 8, 2, 2, 3, 17,
	17, 16, 16, 16, 16, 16, 4, 5, 6, 7,
	11, 11, 11, 10, 10, 12, 12, 13, 13, 15,
	15, 18, 19, 19, 20, 21, 22, 23, 23, 24,
	24, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 29, 14,
	14, 26, 27, 27, 28,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 0, 2, 5, 0,
	3, 1, 1, 1, 1, 1, 4, 1, 4, 4,
	0, 1, 3, 5, 6, 0, 3, 1, 3, 0,
	2, 6, 0, 5, 5, 9, 6, 0, 3, 0,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 1, 1, 1, 1, 3, 1, 1,
	2, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -30, 4, 6, 5, 7, -1, -2, -9, -8,
	-4, -5, -6, -7, -18, -20, -21, -22, 15, 16,
	17, 18, 22, 26, -25, -29, -26, -27, -28, 35,
	-14, 58, 51, 52, 55, 54, 53, -9, -3, 13,
	-8, -29, 35, -11, -10, 54, -25, 27, -29, 38,
	39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, -25, 54, -14, 50, -25, -12, 34, 19, 35,
	-19, 23, 30, 31, -23, 28, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25, 36, 14,
	-25, 36, -15, 21, -10, 20, 36, -25, 24, 58,
	-9, 32, -24, 23, 58, -16, 8, 9, 10, 11,
	12, 55, -13, -14, 14, 36, -9, 37, 25, 20,
	-9, 58, 57, -17, 35, 34, 54, 14, 25, 58,
	58, 25, 57, -14, -14, 54, 56, 57, 36, -9,
	25,
}
var yyDef = []int{

	0, -2, 16, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 11, 12, 13, 14, 15, 0, 27,
	0, 30, 0, 0, 4, 63, 64, 65, 66, 0,
	68, 71, 72, 73, 74, 69, 0, 5, 17, 0,
	7, 0, 0, 35, 31, 0, 42, 0, 47, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 70, 0, 0, 0, 39, 0, 0, 0,
	0, 0, 6, 0, 49, 0, 51, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 61, 62, 67, 0,
	26, 28, 29, 0, 32, 0, 0, 0, 6, 0,
	0, 0, 6, 0, 0, 19, 21, 22, 23, 24,
	25, 40, 36, 37, 0, 0, 0, 0, 44, 0,
	0, 0, 48, 18, 0, 0, 33, 0, 41, 0,
	0, 46, 50, 0, 38, 34, 43, 6, 20, 0,
	45,
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
	52, 53, 54, 55, 56, 57, 58,
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
			yyVAL.query = &Query{}
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
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str, nil)
		}
	case 34:
		//line grammar.y:243
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-5].str, yyS[yypt-3].expr)
		}
	case 35:
		//line grammar.y:250
		{
			yyVAL.strs = make([]string, 0)
		}
	case 36:
		//line grammar.y:254
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 37:
		//line grammar.y:261
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 38:
		//line grammar.y:266
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 39:
		//line grammar.y:273
		{
			yyVAL.str = ""
		}
	case 40:
		//line grammar.y:277
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 41:
		//line grammar.y:284
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 42:
		//line grammar.y:296
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 43:
		//line grammar.y:300
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
	case 44:
		//line grammar.y:315
		{
			yyVAL.event_loop = NewEventLoop()
			yyVAL.event_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 45:
		//line grammar.y:323
		{
			yyVAL.session_loop = NewSessionLoop()
			yyVAL.session_loop.IdleDuration = timeSpanToSeconds(yyS[yypt-3].integer, yyS[yypt-2].str)
			yyVAL.session_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 46:
		//line grammar.y:332
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
	case 47:
		//line grammar.y:349
		{
			yyVAL.integer = 0
		}
	case 48:
		//line grammar.y:353
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 49:
		//line grammar.y:360
		{
			yyVAL.integer = 0
		}
	case 50:
		//line grammar.y:364
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 51:
		//line grammar.y:370
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 52:
		//line grammar.y:371
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 53:
		//line grammar.y:372
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 54:
		//line grammar.y:373
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 55:
		//line grammar.y:374
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 56:
		//line grammar.y:375
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 57:
		//line grammar.y:376
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 58:
		//line grammar.y:377
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 59:
		//line grammar.y:378
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 60:
		//line grammar.y:379
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 61:
		//line grammar.y:380
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 62:
		//line grammar.y:381
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 63:
		//line grammar.y:382
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 64:
		//line grammar.y:383
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 65:
		//line grammar.y:384
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 66:
		//line grammar.y:385
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 67:
		//line grammar.y:386
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 68:
		//line grammar.y:391
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 69:
		//line grammar.y:398
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 70:
		//line grammar.y:402
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 71:
		//line grammar.y:409
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 72:
		//line grammar.y:416
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 73:
		//line grammar.y:420
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 74:
		//line grammar.y:427
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
