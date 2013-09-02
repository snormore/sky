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

//line grammar.y:408
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

const yyNprod = 71
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 244

var yyAct = []int{

	8, 23, 41, 28, 91, 122, 114, 98, 34, 94,
	28, 124, 115, 128, 105, 127, 126, 125, 119, 31,
	32, 29, 33, 43, 107, 30, 31, 32, 29, 33,
	58, 42, 30, 59, 46, 47, 48, 49, 50, 51,
	52, 61, 54, 55, 56, 57, 29, 60, 71, 72,
	73, 74, 75, 76, 77, 78, 79, 80, 81, 82,
	56, 57, 85, 111, 129, 117, 89, 92, 109, 95,
	46, 47, 48, 49, 50, 51, 52, 53, 54, 55,
	56, 57, 46, 47, 48, 49, 50, 51, 52, 53,
	54, 55, 56, 57, 110, 118, 86, 113, 46, 47,
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	83, 67, 46, 47, 48, 49, 50, 51, 52, 53,
	54, 55, 56, 57, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 46, 47, 48, 49,
	50, 51, 44, 68, 54, 55, 56, 57, 48, 49,
	50, 51, 65, 39, 54, 55, 56, 57, 50, 51,
	64, 70, 54, 55, 56, 57, 29, 54, 55, 56,
	57, 93, 24, 63, 17, 18, 19, 20, 97, 88,
	90, 21, 120, 108, 123, 22, 17, 18, 19, 20,
	38, 84, 36, 21, 1, 45, 121, 22, 17, 18,
	19, 20, 37, 27, 26, 21, 25, 9, 112, 22,
	17, 18, 19, 20, 96, 69, 16, 21, 15, 66,
	14, 22, 100, 101, 102, 103, 104, 2, 4, 3,
	5, 116, 99, 87, 106, 62, 40, 13, 12, 11,
	10, 35, 7, 6,
}
var yyPact = []int{

	223, -1000, -1000, -1000, 195, -23, -1000, 179, 195, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -5, -1000, 120,
	-20, -23, 115, 46, -1000, -1000, -1000, -1000, -23, -1000,
	-1000, -1000, -1000, -1000, 195, -1000, -18, -1000, -1, -23,
	141, -1000, 119, 88, 113, 133, -23, -23, -23, -23,
	-23, -23, -23, -23, -23, -23, -23, -23, 76, 177,
	-23, 62, 158, -20, 160, -30, 147, -46, -1000, 155,
	-48, 110, 110, 118, 118, 123, 123, 100, -2, 14,
	14, -1000, -1000, -1000, 214, 46, -1000, -1000, -38, -1000,
	-27, 169, 34, -1000, 28, 183, -1000, -49, -42, 32,
	-1000, -1000, -1000, -1000, -1000, -1000, 63, -1000, -33, 168,
	171, -50, -1000, 159, -43, -1000, -1000, -34, -35, -1000,
	-36, -1000, -40, -1000, -1000, 30, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 243, 242, 241, 240, 239, 238, 237, 202, 0,
	2, 236, 235, 234, 233, 232, 231, 220, 219, 218,
	216, 215, 214, 1, 206, 204, 203, 172, 194,
}
var yyR1 = []int{

	0, 28, 28, 28, 28, 1, 9, 9, 8, 8,
	8, 8, 8, 8, 8, 2, 2, 3, 16, 16,
	15, 15, 15, 15, 15, 4, 5, 6, 7, 11,
	11, 11, 10, 10, 12, 12, 13, 13, 14, 14,
	17, 18, 18, 19, 20, 21, 21, 22, 22, 23,
	23, 23, 23, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 27, 24, 25, 25,
	26,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 0, 2, 5, 0, 3,
	1, 1, 1, 1, 1, 4, 1, 4, 4, 0,
	1, 3, 5, 6, 0, 3, 1, 3, 0, 2,
	6, 0, 5, 5, 6, 0, 3, 0, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 1, 1, 1, 1, 3, 1, 1, 1, 1,
	1,
}
var yyChk = []int{

	-1000, -28, 4, 6, 5, 7, -1, -2, -9, -8,
	-4, -5, -6, -7, -17, -19, -20, 15, 16, 17,
	18, 22, 26, -23, -27, -24, -25, -26, 33, 51,
	55, 49, 50, 52, -9, -3, 13, -8, -27, 33,
	-11, -10, 51, -23, 27, -27, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, -23, 51,
	48, -23, -12, 32, 19, 33, -18, 23, 30, -21,
	28, -23, -23, -23, -23, -23, -23, -23, -23, -23,
	-23, -23, -23, 34, 14, -23, 34, -14, 21, -10,
	20, 34, -23, 24, 55, -9, -22, 23, 55, -15,
	8, 9, 10, 11, 12, 52, -13, 51, 14, 34,
	-9, 35, 25, -9, 55, 54, -16, 33, 32, 51,
	14, 25, 55, 25, 54, 51, 51, 51, 53, 34,
}
var yyDef = []int{

	0, -2, 15, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 11, 12, 13, 14, 0, 26, 0,
	29, 0, 0, 4, 61, 62, 63, 64, 0, 66,
	67, 68, 69, 70, 5, 16, 0, 7, 0, 0,
	34, 30, 0, 41, 0, 45, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 38, 0, 0, 0, 0, 0, 6, 47,
	0, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	58, 59, 60, 65, 0, 25, 27, 28, 0, 31,
	0, 0, 0, 6, 0, 0, 6, 0, 0, 18,
	20, 21, 22, 23, 24, 39, 35, 36, 0, 0,
	0, 0, 43, 0, 0, 46, 17, 0, 0, 32,
	0, 40, 0, 44, 48, 0, 37, 33, 42, 19,
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
			yyVAL.variable = NewVariable(yyS[yypt-3].str, yyS[yypt-1].str)
			yyVAL.variable.Association = yyS[yypt-0].str
		}
	case 18:
		//line grammar.y:166
		{
			yyVAL.str = ""
		}
	case 19:
		//line grammar.y:170
		{
			yyVAL.str = yyS[yypt-1].str
		}
	case 20:
		//line grammar.y:176
		{
			yyVAL.str = core.FactorDataType
		}
	case 21:
		//line grammar.y:177
		{
			yyVAL.str = core.StringDataType
		}
	case 22:
		//line grammar.y:178
		{
			yyVAL.str = core.IntegerDataType
		}
	case 23:
		//line grammar.y:179
		{
			yyVAL.str = core.FloatDataType
		}
	case 24:
		//line grammar.y:180
		{
			yyVAL.str = core.BooleanDataType
		}
	case 25:
		//line grammar.y:185
		{
			yyVAL.assignment = NewAssignment()
			yyVAL.assignment.SetTarget(yyS[yypt-2].var_ref)
			yyVAL.assignment.SetExpression(yyS[yypt-0].expr)
		}
	case 26:
		//line grammar.y:194
		{
			yyVAL.exit = NewExit()
		}
	case 27:
		//line grammar.y:201
		{
			yyVAL.debug = NewDebug()
			yyVAL.debug.SetExpression(yyS[yypt-1].expr)
		}
	case 28:
		//line grammar.y:209
		{
			yyVAL.selection = NewSelection()
			yyVAL.selection.SetFields(yyS[yypt-2].selection_fields)
			yyVAL.selection.Dimensions = yyS[yypt-1].strs
			yyVAL.selection.Name = yyS[yypt-0].str
		}
	case 29:
		//line grammar.y:219
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 30:
		//line grammar.y:223
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 31:
		//line grammar.y:228
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 32:
		//line grammar.y:235
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str, nil)
		}
	case 33:
		//line grammar.y:239
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-5].str, yyS[yypt-3].expr)
		}
	case 34:
		//line grammar.y:246
		{
			yyVAL.strs = make([]string, 0)
		}
	case 35:
		//line grammar.y:250
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 36:
		//line grammar.y:257
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 37:
		//line grammar.y:262
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 38:
		//line grammar.y:269
		{
			yyVAL.str = ""
		}
	case 39:
		//line grammar.y:273
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 40:
		//line grammar.y:280
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 41:
		//line grammar.y:292
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 42:
		//line grammar.y:296
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
	case 43:
		//line grammar.y:311
		{
			yyVAL.event_loop = NewEventLoop()
			yyVAL.event_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 44:
		//line grammar.y:319
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
	case 45:
		//line grammar.y:336
		{
			yyVAL.integer = 0
		}
	case 46:
		//line grammar.y:340
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 47:
		//line grammar.y:347
		{
			yyVAL.integer = 0
		}
	case 48:
		//line grammar.y:351
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 49:
		//line grammar.y:357
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 50:
		//line grammar.y:358
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 51:
		//line grammar.y:359
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 52:
		//line grammar.y:360
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 53:
		//line grammar.y:361
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 54:
		//line grammar.y:362
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 55:
		//line grammar.y:363
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 56:
		//line grammar.y:364
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 57:
		//line grammar.y:365
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 58:
		//line grammar.y:366
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 59:
		//line grammar.y:367
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 60:
		//line grammar.y:368
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 61:
		//line grammar.y:369
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 62:
		//line grammar.y:370
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 63:
		//line grammar.y:371
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 64:
		//line grammar.y:372
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 65:
		//line grammar.y:373
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 66:
		//line grammar.y:378
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 67:
		//line grammar.y:385
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 68:
		//line grammar.y:392
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 69:
		//line grammar.y:396
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 70:
		//line grammar.y:403
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
