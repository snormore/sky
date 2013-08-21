//line grammar.y:2
package query

import __yyfmt__ "fmt"

//line grammar.y:3
//line grammar.y:7
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
	boolean_literal  *BooleanLiteral
	string_literal   *StringLiteral
}

const TSTARTQUERY = 57346
const TSTARTSTATEMENT = 57347
const TSTARTSTATEMENTS = 57348
const TSTARTEXPRESSION = 57349
const TSELECT = 57350
const TGROUP = 57351
const TBY = 57352
const TINTO = 57353
const TAS = 57354
const TWHEN = 57355
const TWITHIN = 57356
const TTHEN = 57357
const TEND = 57358
const TSEMICOLON = 57359
const TCOMMA = 57360
const TLPAREN = 57361
const TRPAREN = 57362
const TRANGE = 57363
const TEQUALS = 57364
const TNOTEQUALS = 57365
const TLT = 57366
const TLTE = 57367
const TGT = 57368
const TGTE = 57369
const TAND = 57370
const TOR = 57371
const TPLUS = 57372
const TMINUS = 57373
const TMUL = 57374
const TDIV = 57375
const TTRUE = 57376
const TFALSE = 57377
const TIDENT = 57378
const TSTRING = 57379
const TWITHINUNITS = 57380
const TINT = 57381

var yyToknames = []string{
	"TSTARTQUERY",
	"TSTARTSTATEMENT",
	"TSTARTSTATEMENTS",
	"TSTARTEXPRESSION",
	"TSELECT",
	"TGROUP",
	"TBY",
	"TINTO",
	"TAS",
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
	"TTRUE",
	"TFALSE",
	"TIDENT",
	"TSTRING",
	"TWITHINUNITS",
	"TINT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line grammar.y:274
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

const yyNprod = 47
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 152

var yyAct = []int{

	7, 14, 19, 82, 8, 69, 85, 71, 84, 66,
	83, 79, 73, 27, 28, 29, 77, 22, 23, 20,
	24, 42, 21, 34, 35, 67, 75, 38, 39, 40,
	41, 46, 49, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 60, 61, 70, 30, 31, 32, 33,
	34, 35, 36, 37, 38, 39, 40, 41, 64, 40,
	41, 78, 48, 38, 39, 40, 41, 68, 80, 76,
	30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 30, 31, 32, 33, 34, 35, 36, 37,
	38, 39, 40, 41, 30, 31, 32, 33, 34, 35,
	36, 74, 38, 39, 40, 41, 30, 31, 32, 33,
	34, 35, 63, 45, 38, 39, 40, 41, 32, 33,
	34, 35, 44, 12, 38, 39, 40, 41, 13, 12,
	65, 81, 1, 15, 13, 2, 4, 3, 5, 25,
	18, 17, 16, 47, 9, 11, 62, 72, 43, 26,
	10, 6,
}
var yyPact = []int{

	131, -1000, -1000, -1000, 121, -17, -1000, 121, 121, -1000,
	-1000, -1000, -22, -17, 60, -1000, -1000, -1000, -1000, -17,
	-1000, -1000, -1000, -1000, -1000, -1000, 104, -1000, 12, 48,
	-17, -17, -17, -17, -17, -17, -17, -17, -17, -17,
	-17, -17, 24, 101, -22, 120, -11, 52, -34, 94,
	94, -3, -3, 33, 33, 84, 72, 27, 27, -1000,
	-1000, -1000, 28, -30, -1000, -24, 89, 6, -1000, -5,
	-1000, -1000, 43, -1000, -25, 56, 115, -36, -26, -1000,
	-28, -1000, -32, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 151, 150, 139, 0, 13, 149, 148, 147, 146,
	145, 143, 1, 142, 141, 140, 133, 132,
}
var yyR1 = []int{

	0, 17, 17, 17, 17, 1, 4, 4, 3, 3,
	2, 6, 6, 6, 5, 5, 7, 7, 8, 8,
	9, 9, 10, 11, 11, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 16, 13, 14, 14, 15,
}
var yyR2 = []int{

	0, 2, 2, 2, 2, 1, 0, 2, 1, 1,
	5, 0, 1, 3, 5, 6, 0, 3, 1, 3,
	0, 2, 6, 0, 5, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 1, 1, 1,
	1, 3, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -17, 4, 6, 5, 7, -1, -4, -4, -3,
	-2, -10, 8, 13, -12, -16, -13, -14, -15, 19,
	36, 39, 34, 35, 37, -3, -6, -5, 36, -12,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, -12, -7, 18, 9, 19, -11, 14, -12,
	-12, -12, -12, -12, -12, -12, -12, -12, -12, -12,
	-12, 20, -9, 11, -5, 10, 20, 36, 15, 39,
	17, 37, -8, 36, 12, 20, -4, 21, 18, 36,
	12, 16, 39, 36, 36, 38,
}
var yyDef = []int{

	0, -2, 6, 6, 0, 0, 1, 5, 2, 3,
	8, 9, 11, 0, 4, 37, 38, 39, 40, 0,
	42, 43, 44, 45, 46, 7, 16, 12, 0, 23,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 20, 0, 0, 0, 0, 0, 25,
	26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
	36, 41, 0, 0, 13, 0, 0, 0, 6, 0,
	10, 21, 17, 18, 0, 0, 0, 0, 0, 14,
	0, 22, 0, 19, 15, 24,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39,
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
		//line grammar.y:69
		{
			l := yylex.(*yylexer)
			l.query = yyS[yypt-0].query
		}
	case 2:
		//line grammar.y:74
		{
			l := yylex.(*yylexer)
			l.statements = yyS[yypt-0].statements
		}
	case 3:
		//line grammar.y:79
		{
			l := yylex.(*yylexer)
			l.statement = yyS[yypt-0].statement
		}
	case 4:
		//line grammar.y:84
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 5:
		//line grammar.y:92
		{
			yyVAL.query = &Query{}
			yyVAL.query.SetStatements(yyS[yypt-0].statements)
		}
	case 6:
		//line grammar.y:100
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 7:
		//line grammar.y:104
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 8:
		//line grammar.y:111
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 9:
		//line grammar.y:115
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 10:
		//line grammar.y:122
		{
			yyVAL.selection = NewSelection()
			yyVAL.selection.SetFields(yyS[yypt-3].selection_fields)
			yyVAL.selection.Dimensions = yyS[yypt-2].strs
			yyVAL.selection.Name = yyS[yypt-1].str
		}
	case 11:
		//line grammar.y:132
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 12:
		//line grammar.y:136
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 13:
		//line grammar.y:141
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 14:
		//line grammar.y:148
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-4].str+"()")
		}
	case 15:
		//line grammar.y:152
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-0].str, yyS[yypt-5].str+"("+yyS[yypt-3].str+")")
		}
	case 16:
		//line grammar.y:159
		{
			yyVAL.strs = make([]string, 0)
		}
	case 17:
		//line grammar.y:163
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 18:
		//line grammar.y:170
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 19:
		//line grammar.y:175
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 20:
		//line grammar.y:182
		{
			yyVAL.str = ""
		}
	case 21:
		//line grammar.y:186
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 22:
		//line grammar.y:193
		{
			yyVAL.condition = NewCondition()
			yyVAL.condition.SetExpression(yyS[yypt-4].expr)
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.SetStatements(yyS[yypt-1].statements)
		}
	case 23:
		//line grammar.y:205
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: UnitSteps}
		}
	case 24:
		//line grammar.y:209
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
	case 25:
		//line grammar.y:223
		{
			yyVAL.expr = NewBinaryExpression(OpEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 26:
		//line grammar.y:224
		{
			yyVAL.expr = NewBinaryExpression(OpNotEquals, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 27:
		//line grammar.y:225
		{
			yyVAL.expr = NewBinaryExpression(OpLessThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 28:
		//line grammar.y:226
		{
			yyVAL.expr = NewBinaryExpression(OpLessThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 29:
		//line grammar.y:227
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThan, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 30:
		//line grammar.y:228
		{
			yyVAL.expr = NewBinaryExpression(OpGreaterThanOrEqualTo, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 31:
		//line grammar.y:229
		{
			yyVAL.expr = NewBinaryExpression(OpAnd, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 32:
		//line grammar.y:230
		{
			yyVAL.expr = NewBinaryExpression(OpOr, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 33:
		//line grammar.y:231
		{
			yyVAL.expr = NewBinaryExpression(OpPlus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 34:
		//line grammar.y:232
		{
			yyVAL.expr = NewBinaryExpression(OpMinus, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 35:
		//line grammar.y:233
		{
			yyVAL.expr = NewBinaryExpression(OpMultiply, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 36:
		//line grammar.y:234
		{
			yyVAL.expr = NewBinaryExpression(OpDivide, yyS[yypt-2].expr, yyS[yypt-0].expr)
		}
	case 37:
		//line grammar.y:235
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 38:
		//line grammar.y:236
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 39:
		//line grammar.y:237
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 40:
		//line grammar.y:238
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 41:
		//line grammar.y:239
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 42:
		//line grammar.y:244
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 43:
		//line grammar.y:251
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 44:
		//line grammar.y:258
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 45:
		//line grammar.y:262
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 46:
		//line grammar.y:269
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
