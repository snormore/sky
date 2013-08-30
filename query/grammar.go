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
const TSEMICOLON = 57373
const TCOMMA = 57374
const TLPAREN = 57375
const TRPAREN = 57376
const TRANGE = 57377
const TEQUALS = 57378
const TNOTEQUALS = 57379
const TLT = 57380
const TLTE = 57381
const TGT = 57382
const TGTE = 57383
const TAND = 57384
const TOR = 57385
const TPLUS = 57386
const TMINUS = 57387
const TMUL = 57388
const TDIV = 57389
const TASSIGN = 57390
const TTRUE = 57391
const TFALSE = 57392
const TIDENT = 57393
const TQUOTEDSTRING = 57394
const TWITHINUNITS = 57395
const TTIMEUNITS = 57396
const TINT = 57397

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

//line grammar.y:396
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

const yyNprod = 69
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 240

var yyAct = []int{

	8, 23, 41, 28, 91, 120, 114, 98, 34, 94,
	28, 122, 115, 125, 105, 124, 123, 117, 107, 31,
	32, 29, 33, 43, 42, 30, 31, 32, 29, 33,
	58, 59, 30, 29, 46, 47, 48, 49, 50, 51,
	52, 61, 54, 55, 56, 57, 56, 57, 71, 72,
	73, 74, 75, 76, 77, 78, 79, 80, 81, 82,
	60, 111, 85, 65, 39, 116, 89, 92, 109, 95,
	46, 47, 48, 49, 50, 51, 52, 53, 54, 55,
	56, 57, 46, 47, 48, 49, 50, 51, 52, 53,
	54, 55, 56, 57, 110, 68, 86, 113, 46, 47,
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	83, 67, 46, 47, 48, 49, 50, 51, 52, 53,
	54, 55, 56, 57, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 46, 47, 48, 49,
	50, 51, 44, 93, 54, 55, 56, 57, 48, 49,
	50, 51, 70, 97, 54, 55, 56, 57, 50, 51,
	64, 88, 54, 55, 56, 57, 29, 54, 55, 56,
	57, 90, 24, 63, 17, 18, 19, 20, 118, 108,
	84, 21, 36, 1, 121, 22, 17, 18, 19, 20,
	38, 27, 26, 21, 25, 45, 119, 22, 17, 18,
	19, 20, 37, 96, 69, 21, 16, 9, 112, 22,
	17, 18, 19, 20, 15, 66, 14, 21, 99, 87,
	106, 22, 100, 101, 102, 103, 104, 2, 4, 3,
	5, 62, 40, 13, 12, 11, 10, 35, 7, 6,
}
var yyPact = []int{

	223, -1000, -1000, -1000, 195, -23, -1000, 169, 195, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -18, -1000, 31,
	-27, -23, 115, 46, -1000, -1000, -1000, -1000, -23, -1000,
	-1000, -1000, -1000, -1000, 195, -1000, -20, -1000, 12, -23,
	141, -1000, 30, 88, 65, 124, -23, -23, -23, -23,
	-23, -23, -23, -23, -23, -23, -23, -23, 76, 166,
	-23, 62, 140, -27, 151, -30, 119, -46, -1000, 130,
	-48, 110, 110, 118, 118, 123, 123, 100, -2, 0,
	0, -1000, -1000, -1000, 214, 46, -1000, -1000, -38, -1000,
	-33, 165, 34, -1000, 26, 183, -1000, -49, -42, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 33, -1000, -34, 164,
	171, -50, -1000, 159, -43, -1000, -35, -1000, -36, -1000,
	-40, -1000, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 239, 238, 237, 236, 235, 234, 233, 202, 0,
	2, 232, 231, 220, 219, 218, 216, 215, 214, 206,
	204, 203, 1, 194, 192, 191, 172, 183,
}
var yyR1 = []int{

	0, 27, 27, 27, 27, 1, 9, 9, 8, 8,
	8, 8, 8, 8, 8, 2, 2, 3, 15, 15,
	15, 15, 15, 4, 5, 6, 7, 11, 11, 11,
	10, 10, 12, 12, 13, 13, 14, 14, 16, 17,
	17, 18, 19, 20, 20, 21, 21, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 26, 23, 24, 24, 25,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 0, 2, 4, 1, 1,
	1, 1, 1, 4, 1, 4, 4, 0, 1, 3,
	5, 6, 0, 3, 1, 3, 0, 2, 6, 0,
	5, 5, 6, 0, 3, 0, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 1,
	1, 1, 1, 3, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -27, 4, 6, 5, 7, -1, -2, -9, -8,
	-4, -5, -6, -7, -16, -18, -19, 15, 16, 17,
	18, 22, 26, -22, -26, -23, -24, -25, 33, 51,
	55, 49, 50, 52, -9, -3, 13, -8, -26, 33,
	-11, -10, 51, -22, 27, -26, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, -22, 51,
	48, -22, -12, 32, 19, 33, -17, 23, 30, -20,
	28, -22, -22, -22, -22, -22, -22, -22, -22, -22,
	-22, -22, -22, 34, 14, -22, 34, -14, 21, -10,
	20, 34, -22, 24, 55, -9, -21, 23, 55, -15,
	8, 9, 10, 11, 12, 52, -13, 51, 14, 34,
	-9, 35, 25, -9, 55, 54, 32, 51, 14, 25,
	55, 25, 54, 51, 51, 53,
}
var yyDef = []int{

	0, -2, 15, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 11, 12, 13, 14, 0, 24, 0,
	27, 0, 0, 4, 59, 60, 61, 62, 0, 64,
	65, 66, 67, 68, 5, 16, 0, 7, 0, 0,
	32, 28, 0, 39, 0, 43, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 36, 0, 0, 0, 0, 0, 6, 45,
	0, 47, 48, 49, 50, 51, 52, 53, 54, 55,
	56, 57, 58, 63, 0, 23, 25, 26, 0, 29,
	0, 0, 0, 6, 0, 0, 6, 0, 0, 17,
	18, 19, 20, 21, 22, 37, 33, 34, 0, 0,
	0, 0, 41, 0, 0, 44, 0, 30, 0, 38,
	0, 42, 46, 35, 31, 40,
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
	52, 53, 54, 55,
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
		//line grammar.y:94
		{
			l := yylex.(*yylexer)
			l.query = yyS[yypt-0].query
		}
	case 2:
		//line grammar.y:99
		{
			l := yylex.(*yylexer)
			l.statements = yyS[yypt-0].statements
		}
	case 3:
		//line grammar.y:104
		{
			l := yylex.(*yylexer)
			l.statement = yyS[yypt-0].statement
		}
	case 4:
		//line grammar.y:109
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 5:
		//line grammar.y:117
		{
			yyVAL.query = &Query{}
			yyVAL.query.SetDeclaredVariables(yyS[yypt-1].variables)
			yyVAL.query.SetStatements(yyS[yypt-0].statements)
		}
	case 6:
		//line grammar.y:126
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 7:
		//line grammar.y:130
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 8:
		//line grammar.y:136
		{
			yyVAL.statement = Statement(yyS[yypt-0].assignment)
		}
	case 9:
		//line grammar.y:137
		{
			yyVAL.statement = Statement(yyS[yypt-0].exit)
		}
	case 10:
		//line grammar.y:138
		{
			yyVAL.statement = Statement(yyS[yypt-0].debug)
		}
	case 11:
		//line grammar.y:139
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 12:
		//line grammar.y:140
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 13:
		//line grammar.y:141
		{
			yyVAL.statement = Statement(yyS[yypt-0].event_loop)
		}
	case 14:
		//line grammar.y:142
		{
			yyVAL.statement = Statement(yyS[yypt-0].temporal_loop)
		}
	case 15:
		//line grammar.y:147
		{
			yyVAL.variables = make([]*Variable, 0)
		}
	case 16:
		//line grammar.y:151
		{
			yyVAL.variables = append(yyS[yypt-1].variables, yyS[yypt-0].variable)
		}
	case 17:
		//line grammar.y:158
		{
			yyVAL.variable = NewVariable(yyS[yypt-2].str, yyS[yypt-0].str)
		}
	case 18:
		//line grammar.y:164
		{
			yyVAL.str = core.FactorDataType
		}
	case 19:
		//line grammar.y:165
		{
			yyVAL.str = core.StringDataType
		}
	case 20:
		//line grammar.y:166
		{
			yyVAL.str = core.IntegerDataType
		}
	case 21:
		//line grammar.y:167
		{
			yyVAL.str = core.FloatDataType
		}
	case 22:
		//line grammar.y:168
		{
			yyVAL.str = core.BooleanDataType
		}
	case 23:
		//line grammar.y:173
		{
			yyVAL.assignment = NewAssignment()
			yyVAL.assignment.SetTarget(yyS[yypt-2].var_ref)
			yyVAL.assignment.SetExpression(yyS[yypt-0].expr)
		}
	case 24:
		//line grammar.y:182
		{
			yyVAL.exit = NewExit()
		}
	case 25:
		//line grammar.y:189
		{
			yyVAL.debug = NewDebug()
			yyVAL.debug.SetExpression(yyS[yypt-1].expr)
		}
	case 26:
		//line grammar.y:197
		{
			yyVAL.selection = NewSelection()
			yyVAL.selection.SetFields(yyS[yypt-2].selection_fields)
			yyVAL.selection.Dimensions = yyS[yypt-1].strs
			yyVAL.selection.Name = yyS[yypt-0].str
		}
	case 27:
		//line grammar.y:207
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 28:
		//line grammar.y:211
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 29:
		//line grammar.y:216
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 30:
		//line grammar.y:223
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str, nil)
		}
	case 31:
		//line grammar.y:227
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-5].str, yyS[yypt-3].expr)
		}
	case 32:
		//line grammar.y:234
		{
			yyVAL.strs = make([]string, 0)
		}
	case 33:
		//line grammar.y:238
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 34:
		//line grammar.y:245
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 35:
		//line grammar.y:250
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 36:
		//line grammar.y:257
		{
			yyVAL.str = ""
		}
	case 37:
		//line grammar.y:261
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 38:
		//line grammar.y:268
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 39:
		//line grammar.y:280
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 40:
		//line grammar.y:284
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
	case 41:
		//line grammar.y:299
		{
			yyVAL.event_loop = NewEventLoop()
			yyVAL.event_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 42:
		//line grammar.y:307
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
	case 43:
		//line grammar.y:324
		{
			yyVAL.integer = 0
		}
	case 44:
		//line grammar.y:328
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 45:
		//line grammar.y:335
		{
			yyVAL.integer = 0
		}
	case 46:
		//line grammar.y:339
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 47:
		//line grammar.y:345
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 48:
		//line grammar.y:346
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 49:
		//line grammar.y:347
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 50:
		//line grammar.y:348
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 51:
		//line grammar.y:349
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 52:
		//line grammar.y:350
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 53:
		//line grammar.y:351
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 54:
		//line grammar.y:352
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 55:
		//line grammar.y:353
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 56:
		//line grammar.y:354
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 57:
		//line grammar.y:355
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 58:
		//line grammar.y:356
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 59:
		//line grammar.y:357
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 60:
		//line grammar.y:358
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 61:
		//line grammar.y:359
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 62:
		//line grammar.y:360
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 63:
		//line grammar.y:361
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 64:
		//line grammar.y:366
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 65:
		//line grammar.y:373
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 66:
		//line grammar.y:380
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 67:
		//line grammar.y:384
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 68:
		//line grammar.y:391
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
