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
const TSELECT = 57359
const TGROUP = 57360
const TBY = 57361
const TINTO = 57362
const TWHEN = 57363
const TWITHIN = 57364
const TTHEN = 57365
const TEND = 57366
const TFOR = 57367
const TEACH = 57368
const TEVERY = 57369
const TIN = 57370
const TEVENT = 57371
const TSEMICOLON = 57372
const TCOMMA = 57373
const TLPAREN = 57374
const TRPAREN = 57375
const TRANGE = 57376
const TEQUALS = 57377
const TNOTEQUALS = 57378
const TLT = 57379
const TLTE = 57380
const TGT = 57381
const TGTE = 57382
const TAND = 57383
const TOR = 57384
const TPLUS = 57385
const TMINUS = 57386
const TMUL = 57387
const TDIV = 57388
const TASSIGN = 57389
const TTRUE = 57390
const TFALSE = 57391
const TIDENT = 57392
const TQUOTEDSTRING = 57393
const TWITHINUNITS = 57394
const TTIMEUNITS = 57395
const TINT = 57396

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

//line grammar.y:385
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

const yyNprod = 67
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 220

var yyAct = []int{

	8, 21, 38, 115, 109, 93, 89, 117, 32, 110,
	120, 100, 119, 26, 86, 118, 45, 46, 47, 48,
	26, 40, 51, 52, 53, 54, 112, 102, 55, 29,
	30, 27, 31, 41, 39, 28, 29, 30, 27, 31,
	56, 27, 28, 53, 54, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 27, 57, 81,
	106, 61, 84, 87, 104, 90, 43, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 53, 54, 47, 48,
	111, 64, 51, 52, 53, 54, 66, 88, 92, 105,
	83, 79, 108, 43, 44, 45, 46, 47, 48, 49,
	50, 51, 52, 53, 54, 63, 43, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 53, 54, 43, 44,
	45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
	43, 44, 45, 46, 47, 48, 49, 85, 51, 52,
	53, 54, 43, 44, 45, 46, 47, 48, 60, 113,
	51, 52, 53, 54, 51, 52, 53, 54, 16, 17,
	18, 59, 103, 80, 19, 34, 22, 116, 20, 16,
	17, 18, 35, 1, 25, 19, 24, 9, 114, 20,
	16, 17, 18, 36, 23, 91, 19, 42, 65, 107,
	20, 16, 17, 18, 15, 14, 62, 19, 13, 94,
	82, 20, 95, 96, 97, 98, 99, 2, 4, 3,
	5, 101, 58, 37, 12, 11, 10, 33, 7, 6,
}
var yyPact = []int{

	203, -1000, -1000, -1000, 176, -12, -1000, 152, 176, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -9, -1000, -16, -12,
	7, 71, -1000, -1000, -1000, -1000, -12, -1000, -1000, -1000,
	-1000, -1000, 176, -1000, -10, -1000, 11, 130, -1000, 29,
	83, 52, 59, -12, -12, -12, -12, -12, -12, -12,
	-12, -12, -12, -12, -12, 58, 149, -12, 70, -16,
	118, -19, 64, -48, -1000, 66, -49, -21, -21, 39,
	39, 111, 111, 107, 95, -2, -2, -1000, -1000, -1000,
	194, 71, -1000, -40, -1000, -23, 148, 31, -1000, 26,
	165, -1000, -50, -44, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 49, -1000, -24, 135, 154, -51, -1000, 143, -46,
	-1000, -35, -1000, -38, -1000, -42, -1000, -1000, -1000, -1000,
	-1000,
}
var yyPgo = []int{

	0, 219, 218, 217, 216, 215, 214, 172, 0, 2,
	213, 212, 211, 200, 199, 198, 196, 195, 194, 188,
	185, 1, 184, 176, 174, 166, 173,
}
var yyR1 = []int{

	0, 26, 26, 26, 26, 1, 8, 8, 7, 7,
	7, 7, 7, 7, 2, 2, 3, 14, 14, 14,
	14, 14, 4, 5, 6, 10, 10, 10, 9, 9,
	11, 11, 12, 12, 13, 13, 15, 16, 16, 17,
	18, 19, 19, 20, 20, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	21, 21, 25, 22, 23, 23, 24,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 0, 2, 4, 1, 1, 1,
	1, 1, 4, 1, 4, 0, 1, 3, 5, 6,
	0, 3, 1, 3, 0, 2, 6, 0, 5, 5,
	6, 0, 3, 0, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 1, 1, 1,
	1, 3, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -26, 4, 6, 5, 7, -1, -2, -8, -7,
	-4, -5, -6, -15, -17, -18, 15, 16, 17, 21,
	25, -21, -25, -22, -23, -24, 32, 50, 54, 48,
	49, 51, -8, -3, 13, -7, -25, -10, -9, 50,
	-21, 26, -25, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, -21, 50, 47, -11, 31,
	18, 32, -16, 22, 29, -19, 27, -21, -21, -21,
	-21, -21, -21, -21, -21, -21, -21, -21, -21, 33,
	14, -21, -13, 20, -9, 19, 33, -21, 23, 54,
	-8, -20, 22, 54, -14, 8, 9, 10, 11, 12,
	51, -12, 50, 14, 33, -8, 34, 24, -8, 54,
	53, 31, 50, 14, 24, 54, 24, 53, 50, 50,
	52,
}
var yyDef = []int{

	0, -2, 14, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 11, 12, 13, 0, 23, 25, 0,
	0, 4, 57, 58, 59, 60, 0, 62, 63, 64,
	65, 66, 5, 15, 0, 7, 0, 30, 26, 0,
	37, 0, 41, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 34, 0,
	0, 0, 0, 0, 6, 43, 0, 45, 46, 47,
	48, 49, 50, 51, 52, 53, 54, 55, 56, 61,
	0, 22, 24, 0, 27, 0, 0, 0, 6, 0,
	0, 6, 0, 0, 16, 17, 18, 19, 20, 21,
	35, 31, 32, 0, 0, 0, 0, 39, 0, 0,
	42, 0, 28, 0, 36, 0, 40, 44, 33, 29,
	38,
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
	52, 53, 54,
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
		//line grammar.y:92
		{
			l := yylex.(*yylexer)
			l.query = yyS[yypt-0].query
		}
	case 2:
		//line grammar.y:97
		{
			l := yylex.(*yylexer)
			l.statements = yyS[yypt-0].statements
		}
	case 3:
		//line grammar.y:102
		{
			l := yylex.(*yylexer)
			l.statement = yyS[yypt-0].statement
		}
	case 4:
		//line grammar.y:107
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 5:
		//line grammar.y:115
		{
			yyVAL.query = &Query{}
			yyVAL.query.SetDeclaredVariables(yyS[yypt-1].variables)
			yyVAL.query.SetStatements(yyS[yypt-0].statements)
		}
	case 6:
		//line grammar.y:124
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 7:
		//line grammar.y:128
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 8:
		//line grammar.y:134
		{
			yyVAL.statement = Statement(yyS[yypt-0].assignment)
		}
	case 9:
		//line grammar.y:135
		{
			yyVAL.statement = Statement(yyS[yypt-0].exit)
		}
	case 10:
		//line grammar.y:136
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 11:
		//line grammar.y:137
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 12:
		//line grammar.y:138
		{
			yyVAL.statement = Statement(yyS[yypt-0].event_loop)
		}
	case 13:
		//line grammar.y:139
		{
			yyVAL.statement = Statement(yyS[yypt-0].temporal_loop)
		}
	case 14:
		//line grammar.y:144
		{
			yyVAL.variables = make([]*Variable, 0)
		}
	case 15:
		//line grammar.y:148
		{
			yyVAL.variables = append(yyS[yypt-1].variables, yyS[yypt-0].variable)
		}
	case 16:
		//line grammar.y:155
		{
			yyVAL.variable = NewVariable(yyS[yypt-2].str, yyS[yypt-0].str)
		}
	case 17:
		//line grammar.y:161
		{
			yyVAL.str = core.FactorDataType
		}
	case 18:
		//line grammar.y:162
		{
			yyVAL.str = core.StringDataType
		}
	case 19:
		//line grammar.y:163
		{
			yyVAL.str = core.IntegerDataType
		}
	case 20:
		//line grammar.y:164
		{
			yyVAL.str = core.FloatDataType
		}
	case 21:
		//line grammar.y:165
		{
			yyVAL.str = core.BooleanDataType
		}
	case 22:
		//line grammar.y:170
		{
			yyVAL.assignment = NewAssignment()
			yyVAL.assignment.SetTarget(yyS[yypt-2].var_ref)
			yyVAL.assignment.SetExpression(yyS[yypt-0].expr)
		}
	case 23:
		//line grammar.y:179
		{
			yyVAL.exit = NewExit()
		}
	case 24:
		//line grammar.y:186
		{
			yyVAL.selection = NewSelection()
			yyVAL.selection.SetFields(yyS[yypt-2].selection_fields)
			yyVAL.selection.Dimensions = yyS[yypt-1].strs
			yyVAL.selection.Name = yyS[yypt-0].str
		}
	case 25:
		//line grammar.y:196
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 26:
		//line grammar.y:200
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 27:
		//line grammar.y:205
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 28:
		//line grammar.y:212
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str, nil)
		}
	case 29:
		//line grammar.y:216
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-5].str, yyS[yypt-3].expr)
		}
	case 30:
		//line grammar.y:223
		{
			yyVAL.strs = make([]string, 0)
		}
	case 31:
		//line grammar.y:227
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 32:
		//line grammar.y:234
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 33:
		//line grammar.y:239
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 34:
		//line grammar.y:246
		{
			yyVAL.str = ""
		}
	case 35:
		//line grammar.y:250
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 36:
		//line grammar.y:257
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 37:
		//line grammar.y:269
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 38:
		//line grammar.y:273
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
	case 39:
		//line grammar.y:288
		{
			yyVAL.event_loop = NewEventLoop()
			yyVAL.event_loop.SetStatements(yyS[yypt-1].statements)
		}
	case 40:
		//line grammar.y:296
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
	case 41:
		//line grammar.y:313
		{
			yyVAL.integer = 0
		}
	case 42:
		//line grammar.y:317
		{
			yyVAL.integer = timeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
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
		//line grammar.y:334
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 46:
		//line grammar.y:335
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 47:
		//line grammar.y:336
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 48:
		//line grammar.y:337
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 49:
		//line grammar.y:338
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 50:
		//line grammar.y:339
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 51:
		//line grammar.y:340
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 52:
		//line grammar.y:341
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 53:
		//line grammar.y:342
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 54:
		//line grammar.y:343
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 55:
		//line grammar.y:344
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 56:
		//line grammar.y:345
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 57:
		//line grammar.y:346
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 58:
		//line grammar.y:347
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 59:
		//line grammar.y:348
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 60:
		//line grammar.y:349
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 61:
		//line grammar.y:350
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 62:
		//line grammar.y:355
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 63:
		//line grammar.y:362
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 64:
		//line grammar.y:369
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 65:
		//line grammar.y:373
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 66:
		//line grammar.y:380
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
