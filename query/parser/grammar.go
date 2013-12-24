//line grammar.y:2
package parser

import __yyfmt__ "fmt"

//line grammar.y:3
import (
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

//line grammar.y:12
type yySymType struct {
	yys              int
	token            int
	integer          int
	boolean          bool
	str              string
	strs             []string
	query            *ast.Query
	var_decl         *ast.VarDecl
	var_decls        ast.VarDecls
	statement        ast.Statement
	statements       ast.Statements
	assignment       *ast.Assignment
	exit             *ast.Exit
	debug            *ast.Debug
	selection        *ast.Selection
	field            *ast.Field
	fields           ast.Fields
	condition        *ast.Condition
	condition_within *within
	event_loop       *ast.EventLoop
	session_loop     *ast.SessionLoop
	temporal_loop    *ast.TemporalLoop
	expr             ast.Expression
	var_ref          *ast.VarRef
	integer_literal  *ast.IntegerLiteral
	boolean_literal  *ast.BooleanLiteral
	string_literal   *ast.StringLiteral
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
const TDISTINCT = 57364
const TWHEN = 57365
const TWITHIN = 57366
const TTHEN = 57367
const TEND = 57368
const TFOR = 57369
const TEACH = 57370
const TEVERY = 57371
const TIN = 57372
const TEVENT = 57373
const TSESSION = 57374
const TDELIMITED = 57375
const TSEMICOLON = 57376
const TCOMMA = 57377
const TLPAREN = 57378
const TRPAREN = 57379
const TRANGE = 57380
const TEQUALS = 57381
const TNOTEQUALS = 57382
const TLT = 57383
const TLTE = 57384
const TGT = 57385
const TGTE = 57386
const TAND = 57387
const TOR = 57388
const TPLUS = 57389
const TMINUS = 57390
const TMUL = 57391
const TDIV = 57392
const TASSIGN = 57393
const TTRUE = 57394
const TFALSE = 57395
const TAMPERSAND = 57396
const TIDENT = 57397
const TAMPIDENT = 57398
const TQUOTEDSTRING = 57399
const TWITHINUNITS = 57400
const TTIMEUNITS = 57401
const TINT = 57402

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
	"TDISTINCT",
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

//line grammar.y:451
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

const yyNprod = 80
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 288

var yyAct = []int{

	8, 30, 69, 45, 44, 133, 132, 125, 37, 24,
	108, 29, 103, 140, 135, 126, 139, 35, 36, 98,
	115, 48, 64, 29, 121, 128, 47, 32, 33, 141,
	35, 36, 34, 62, 71, 31, 42, 60, 61, 32,
	33, 63, 46, 36, 34, 129, 65, 31, 35, 36,
	58, 59, 60, 61, 78, 79, 80, 81, 82, 83,
	84, 85, 86, 87, 88, 89, 105, 107, 92, 74,
	75, 77, 96, 102, 130, 104, 50, 51, 52, 53,
	54, 55, 56, 57, 58, 59, 60, 61, 52, 53,
	54, 55, 95, 123, 58, 59, 60, 61, 97, 117,
	70, 91, 118, 120, 119, 39, 93, 124, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	54, 55, 101, 40, 58, 59, 60, 61, 9, 1,
	136, 137, 28, 138, 27, 26, 106, 99, 76, 73,
	90, 142, 50, 51, 52, 53, 54, 55, 56, 57,
	58, 59, 60, 61, 50, 51, 52, 53, 54, 55,
	56, 57, 58, 59, 60, 61, 70, 50, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 60, 61, 50,
	51, 52, 53, 54, 55, 56, 17, 58, 59, 60,
	61, 50, 51, 52, 53, 54, 55, 56, 57, 58,
	59, 60, 61, 50, 51, 52, 53, 54, 55, 68,
	16, 58, 59, 60, 61, 18, 19, 20, 21, 2,
	4, 3, 5, 22, 15, 67, 143, 23, 18, 19,
	20, 21, 72, 14, 100, 127, 22, 109, 94, 134,
	23, 18, 19, 20, 21, 116, 66, 43, 13, 22,
	12, 11, 131, 23, 18, 19, 20, 21, 25, 10,
	38, 7, 22, 6, 0, 122, 23, 18, 19, 20,
	21, 0, 0, 0, 0, 22, 0, 41, 0, 23,
	0, 0, 49, 110, 111, 112, 113, 114,
}
var yyPact = []int{

	215, -1000, -1000, -1000, 252, -25, -1000, 92, 252, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -38, -1000,
	0, -13, -25, -7, 128, -1000, -1000, -1000, -1000, -25,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 252, -1000, -38,
	-1000, -29, -25, 190, -1000, 152, -2, 115, 38, 42,
	-25, -25, -25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, 103, 87, -25, 69, 71, -13, 78, -1000,
	-36, 100, 48, -48, -1000, 33, 43, -50, 47, 47,
	77, 77, 3, 3, 164, 140, -12, -12, -1000, -1000,
	-1000, 275, 128, -1000, -1000, -37, -1000, -38, -1000, 86,
	-25, -1000, -1000, -14, 239, 73, -1000, -53, -44, -11,
	-1000, -1000, -1000, -1000, -1000, -1000, 10, -1000, -1000, 37,
	226, -54, -1000, -55, 213, -45, -1000, -1000, -38, -38,
	86, -1000, -42, -46, -1000, -1000, -8, -1000, -1000, -1000,
	-1000, -1000, 200, -1000,
}
var yyPgo = []int{

	0, 263, 261, 260, 259, 251, 250, 248, 123, 0,
	4, 247, 246, 245, 1, 238, 2, 237, 235, 234,
	233, 232, 224, 210, 186, 138, 136, 3, 135, 134,
	132, 258, 129,
}
var yyR1 = []int{

	0, 32, 32, 32, 32, 1, 9, 9, 8, 8,
	8, 8, 8, 8, 8, 8, 2, 2, 3, 18,
	18, 17, 17, 17, 17, 17, 4, 5, 6, 7,
	11, 11, 11, 10, 10, 10, 19, 19, 16, 16,
	12, 12, 13, 13, 15, 15, 20, 21, 21, 22,
	23, 24, 25, 25, 26, 26, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 31, 14, 14, 28, 29, 29, 30,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 0, 2, 5, 0,
	3, 1, 1, 1, 1, 1, 4, 1, 4, 4,
	0, 1, 3, 2, 4, 6, 0, 1, 0, 2,
	0, 3, 1, 3, 0, 2, 6, 0, 5, 5,
	9, 6, 0, 3, 0, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 1, 1,
	1, 1, 3, 1, 1, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -32, 4, 6, 5, 7, -1, -2, -9, -8,
	-4, -5, -6, -7, -20, -22, -23, -24, 15, 16,
	17, 18, 23, 27, -27, -31, -28, -29, -30, 36,
	-14, 60, 52, 53, 57, 55, 56, -9, -3, 13,
	-8, -31, 36, -11, -10, -27, 55, -27, 28, -31,
	39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, 50, -27, -14, 51, -27, -12, 35, 19, -16,
	14, 36, -21, 24, 31, 32, -25, 29, -27, -27,
	-27, -27, -27, -27, -27, -27, -27, -27, -27, -27,
	37, 14, -27, 37, -15, 21, -10, 20, 55, 37,
	-19, 22, 25, 60, -9, 33, -26, 24, 60, -17,
	8, 9, 10, 11, 12, 57, -13, -14, -16, -27,
	-9, 38, 26, 20, -9, 60, 59, -18, 36, 35,
	37, 26, 60, 60, 26, 59, -14, -14, -16, 58,
	59, 37, -9, 26,
}
var yyDef = []int{

	0, -2, 16, 6, 0, 0, 1, 6, 2, 3,
	8, 9, 10, 11, 12, 13, 14, 15, 0, 27,
	0, 30, 0, 0, 4, 68, 69, 70, 71, 0,
	73, 76, 77, 78, 79, 74, 75, 5, 17, 0,
	7, 0, 0, 40, 31, 38, 74, 47, 0, 52,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 44, 0, 0, 33,
	0, 36, 0, 0, 6, 0, 54, 0, 56, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 66, 67,
	72, 0, 26, 28, 29, 0, 32, 0, 39, 38,
	0, 37, 6, 0, 0, 0, 6, 0, 0, 19,
	21, 22, 23, 24, 25, 45, 41, 42, 34, 0,
	0, 0, 49, 0, 0, 0, 53, 18, 0, 0,
	38, 46, 0, 0, 51, 55, 0, 43, 35, 48,
	6, 20, 0, 50,
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
	52, 53, 54, 55, 56, 57, 58, 59, 60,
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
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
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
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
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
		//line grammar.y:100
		{
			l := yylex.(*yylexer)
			l.query = yyS[yypt-0].query
		}
	case 2:
		//line grammar.y:105
		{
			l := yylex.(*yylexer)
			l.statements = yyS[yypt-0].statements
		}
	case 3:
		//line grammar.y:110
		{
			l := yylex.(*yylexer)
			l.statement = yyS[yypt-0].statement
		}
	case 4:
		//line grammar.y:115
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 5:
		//line grammar.y:123
		{
			yyVAL.query = ast.NewQuery()
			yyVAL.query.DeclaredVarDecls = yyS[yypt-1].var_decls
			yyVAL.query.Statements = yyS[yypt-0].statements
		}
	case 6:
		//line grammar.y:132
		{
			yyVAL.statements = make(ast.Statements, 0)
		}
	case 7:
		//line grammar.y:136
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 8:
		//line grammar.y:142
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].assignment)
		}
	case 9:
		//line grammar.y:143
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].exit)
		}
	case 10:
		//line grammar.y:144
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].debug)
		}
	case 11:
		//line grammar.y:145
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].selection)
		}
	case 12:
		//line grammar.y:146
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].condition)
		}
	case 13:
		//line grammar.y:147
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].event_loop)
		}
	case 14:
		//line grammar.y:148
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].session_loop)
		}
	case 15:
		//line grammar.y:149
		{
			yyVAL.statement = ast.Statement(yyS[yypt-0].temporal_loop)
		}
	case 16:
		//line grammar.y:154
		{
			yyVAL.var_decls = make(ast.VarDecls, 0)
		}
	case 17:
		//line grammar.y:158
		{
			yyVAL.var_decls = append(yyS[yypt-1].var_decls, yyS[yypt-0].var_decl)
		}
	case 18:
		//line grammar.y:165
		{
			yyVAL.var_decl = ast.NewVarDecl(0, yyS[yypt-3].str, yyS[yypt-1].str)
			yyVAL.var_decl.Association = yyS[yypt-0].str
		}
	case 19:
		//line grammar.y:173
		{
			yyVAL.str = ""
		}
	case 20:
		//line grammar.y:177
		{
			yyVAL.str = yyS[yypt-1].str
		}
	case 21:
		//line grammar.y:183
		{
			yyVAL.str = core.FactorDataType
		}
	case 22:
		//line grammar.y:184
		{
			yyVAL.str = core.StringDataType
		}
	case 23:
		//line grammar.y:185
		{
			yyVAL.str = core.IntegerDataType
		}
	case 24:
		//line grammar.y:186
		{
			yyVAL.str = core.FloatDataType
		}
	case 25:
		//line grammar.y:187
		{
			yyVAL.str = core.BooleanDataType
		}
	case 26:
		//line grammar.y:192
		{
			yyVAL.assignment = ast.NewAssignment()
			yyVAL.assignment.Target = yyS[yypt-2].var_ref
			yyVAL.assignment.Expression = yyS[yypt-0].expr
		}
	case 27:
		//line grammar.y:201
		{
			yyVAL.exit = ast.NewExit()
		}
	case 28:
		//line grammar.y:208
		{
			yyVAL.debug = ast.NewDebug()
			yyVAL.debug.Expression = yyS[yypt-1].expr
		}
	case 29:
		//line grammar.y:216
		{
			yyVAL.selection = ast.NewSelection()
			yyVAL.selection.Fields = yyS[yypt-2].fields
			yyVAL.selection.Dimensions = yyS[yypt-1].strs
			yyVAL.selection.Name = yyS[yypt-0].str
		}
	case 30:
		//line grammar.y:226
		{
			yyVAL.fields = make(ast.Fields, 0)
		}
	case 31:
		//line grammar.y:230
		{
			yyVAL.fields = make(ast.Fields, 0)
			yyVAL.fields = append(yyVAL.fields, yyS[yypt-0].field)
		}
	case 32:
		//line grammar.y:235
		{
			yyVAL.fields = append(yyS[yypt-2].fields, yyS[yypt-0].field)
		}
	case 33:
		//line grammar.y:242
		{
			yyVAL.field = ast.NewField(yyS[yypt-0].str, "", yyS[yypt-1].expr)
		}
	case 34:
		//line grammar.y:246
		{
			yyVAL.field = ast.NewField(yyS[yypt-0].str, yyS[yypt-3].str, nil)
		}
	case 35:
		//line grammar.y:250
		{
			yyVAL.field = ast.NewField(yyS[yypt-0].str, yyS[yypt-5].str, yyS[yypt-2].expr)
			yyVAL.field.Distinct = yyS[yypt-3].boolean
		}
	case 36:
		//line grammar.y:257
		{
			yyVAL.boolean = false
		}
	case 37:
		//line grammar.y:258
		{
			yyVAL.boolean = true
		}
	case 38:
		//line grammar.y:262
		{
			yyVAL.str = ""
		}
	case 39:
		//line grammar.y:263
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 40:
		//line grammar.y:268
		{
			yyVAL.strs = make([]string, 0)
		}
	case 41:
		//line grammar.y:272
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 42:
		//line grammar.y:279
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 43:
		//line grammar.y:284
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 44:
		//line grammar.y:291
		{
			yyVAL.str = ""
		}
	case 45:
		//line grammar.y:295
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 46:
		//line grammar.y:302
		{
			yyVAL.condition = ast.NewCondition()
			yyVAL.condition.Expression = yyS[yypt-4].expr
			yyVAL.condition.Start = yyS[yypt-3].condition_within.start
			yyVAL.condition.End = yyS[yypt-3].condition_within.end
			yyVAL.condition.UOM = yyS[yypt-3].condition_within.units
			yyVAL.condition.Statements = yyS[yypt-1].statements
		}
	case 47:
		//line grammar.y:314
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: ast.UnitSteps}
		}
	case 48:
		//line grammar.y:318
		{
			yyVAL.condition_within = &within{start: yyS[yypt-3].integer, end: yyS[yypt-1].integer}
			switch yyS[yypt-0].str {
			case "STEPS":
				yyVAL.condition_within.units = ast.UnitSteps
			case "SESSIONS":
				yyVAL.condition_within.units = ast.UnitSessions
			case "SECONDS":
				yyVAL.condition_within.units = ast.UnitSeconds
			}
		}
	case 49:
		//line grammar.y:333
		{
			yyVAL.event_loop = ast.NewEventLoop()
			yyVAL.event_loop.Statements = yyS[yypt-1].statements
		}
	case 50:
		//line grammar.y:341
		{
			yyVAL.session_loop = ast.NewSessionLoop()
			yyVAL.session_loop.IdleDuration = ast.TimeSpanToSeconds(yyS[yypt-3].integer, yyS[yypt-2].str)
			yyVAL.session_loop.Statements = yyS[yypt-1].statements
		}
	case 51:
		//line grammar.y:350
		{
			yyVAL.temporal_loop = ast.NewTemporalLoop()
			yyVAL.temporal_loop.Iterator = yyS[yypt-4].var_ref
			yyVAL.temporal_loop.Step = yyS[yypt-3].integer
			yyVAL.temporal_loop.Duration = yyS[yypt-2].integer
			yyVAL.temporal_loop.Statements = yyS[yypt-1].statements

			// Default steps to 1 of the unit of the duration.
			if yyVAL.temporal_loop.Step == 0 && yyVAL.temporal_loop.Duration > 0 {
				_, units := ast.SecondsToTimeSpan(yyS[yypt-2].integer)
				yyVAL.temporal_loop.Step = ast.TimeSpanToSeconds(1, units)
			}
		}
	case 52:
		//line grammar.y:367
		{
			yyVAL.integer = 0
		}
	case 53:
		//line grammar.y:371
		{
			yyVAL.integer = ast.TimeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 54:
		//line grammar.y:378
		{
			yyVAL.integer = 0
		}
	case 55:
		//line grammar.y:382
		{
			yyVAL.integer = ast.TimeSpanToSeconds(yyS[yypt-1].integer, yyS[yypt-0].str)
		}
	case 56:
		//line grammar.y:388
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 57:
		//line grammar.y:389
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 58:
		//line grammar.y:390
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 59:
		//line grammar.y:391
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 60:
		//line grammar.y:392
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 61:
		//line grammar.y:393
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 62:
		//line grammar.y:394
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 63:
		//line grammar.y:395
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 64:
		//line grammar.y:396
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 65:
		//line grammar.y:397
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 66:
		//line grammar.y:398
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 67:
		//line grammar.y:399
		{
			yyVAL.expr = ast.NewBinaryExpression(ast.OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 68:
		//line grammar.y:400
		{
			yyVAL.expr = ast.Expression(yyS[yypt-0].var_ref)
		}
	case 69:
		//line grammar.y:401
		{
			yyVAL.expr = ast.Expression(yyS[yypt-0].integer_literal)
		}
	case 70:
		//line grammar.y:402
		{
			yyVAL.expr = ast.Expression(yyS[yypt-0].boolean_literal)
		}
	case 71:
		//line grammar.y:403
		{
			yyVAL.expr = ast.Expression(yyS[yypt-0].string_literal)
		}
	case 72:
		//line grammar.y:404
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 73:
		//line grammar.y:409
		{
			yyVAL.var_ref = ast.NewVarRef()
			yyVAL.var_ref.Name = yyS[yypt-0].str
		}
	case 74:
		//line grammar.y:417
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 75:
		//line grammar.y:421
		{
			yyVAL.str = yyS[yypt-0].str[1:]
		}
	case 76:
		//line grammar.y:428
		{
			yyVAL.integer_literal = ast.NewIntegerLiteral(yyS[yypt-0].integer)
		}
	case 77:
		//line grammar.y:435
		{
			yyVAL.boolean_literal = ast.NewBooleanLiteral(true)
		}
	case 78:
		//line grammar.y:439
		{
			yyVAL.boolean_literal = ast.NewBooleanLiteral(false)
		}
	case 79:
		//line grammar.y:446
		{
			yyVAL.string_literal = ast.NewStringLiteral(yyS[yypt-0].str)
		}
	}
	goto yystack /* stack new state and value */
}
