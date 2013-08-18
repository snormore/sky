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
const TSTARTEXPRESSION = 57348
const TSELECT = 57349
const TGROUP = 57350
const TBY = 57351
const TINTO = 57352
const TWHEN = 57353
const TWITHIN = 57354
const TTHEN = 57355
const TEND = 57356
const TSEMICOLON = 57357
const TCOMMA = 57358
const TLPAREN = 57359
const TRPAREN = 57360
const TRANGE = 57361
const TEQUALS = 57362
const TNOTEQUALS = 57363
const TLT = 57364
const TLTE = 57365
const TGT = 57366
const TGTE = 57367
const TAND = 57368
const TOR = 57369
const TPLUS = 57370
const TMINUS = 57371
const TMUL = 57372
const TDIV = 57373
const TTRUE = 57374
const TFALSE = 57375
const TIDENT = 57376
const TSTRING = 57377
const TWITHINUNITS = 57378
const TINT = 57379

var yyToknames = []string{
	"TSTARTQUERY",
	"TSTARTSTATEMENT",
	"TSTARTEXPRESSION",
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

//line grammar.y:264

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

const yyNprod = 46
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 146

var yyAct = []int{

	6, 12, 77, 67, 79, 17, 69, 28, 29, 30,
	31, 32, 33, 27, 25, 36, 37, 38, 39, 40,
	20, 21, 18, 22, 78, 19, 71, 64, 26, 74,
	47, 48, 49, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 65, 28, 29, 30, 31, 32, 33,
	34, 35, 36, 37, 38, 39, 72, 62, 32, 33,
	46, 44, 36, 37, 38, 39, 75, 73, 28, 29,
	30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	28, 29, 30, 31, 32, 33, 34, 35, 36, 37,
	38, 39, 28, 29, 30, 31, 32, 33, 34, 68,
	36, 37, 38, 39, 30, 31, 32, 33, 38, 39,
	36, 37, 38, 39, 36, 37, 38, 39, 43, 66,
	10, 10, 61, 1, 11, 11, 42, 76, 63, 2,
	3, 4, 23, 13, 16, 15, 7, 14, 45, 9,
	60, 70, 41, 24, 8, 5,
}
var yyPact = []int{

	125, -1000, -1000, 114, -12, -1000, 114, -1000, -1000, -1000,
	-6, -12, 60, -1000, -1000, -1000, -1000, -12, -1000, -1000,
	-1000, -1000, -1000, -1000, 110, -1000, 44, 48, -12, -12,
	-12, -12, -12, -12, -12, -12, -12, -12, -12, -12,
	24, 112, -6, 119, 9, 106, -34, 82, 82, 34,
	34, 86, 86, -13, 72, 78, 78, -1000, -1000, -1000,
	84, -29, -1000, -8, -1000, 38, -1000, 10, -1000, -1000,
	50, -1000, -1000, 113, -35, -10, -1000, -32, -1000, -1000,
}
var yyPgo = []int{

	0, 145, 144, 132, 0, 14, 143, 142, 141, 140,
	139, 138, 1, 137, 135, 134, 133, 123,
}
var yyR1 = []int{

	0, 17, 17, 17, 1, 4, 4, 3, 3, 2,
	6, 6, 6, 5, 5, 7, 7, 8, 8, 9,
	9, 10, 11, 11, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 16, 13, 14, 14, 15,
}
var yyR2 = []int{

	0, 2, 2, 2, 1, 0, 2, 1, 1, 5,
	0, 1, 3, 3, 4, 0, 3, 1, 3, 0,
	2, 6, 0, 5, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 1, 1, 1, 1,
	3, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -17, 4, 5, 6, -1, -4, -3, -2, -10,
	7, 11, -12, -16, -13, -14, -15, 17, 34, 37,
	32, 33, 35, -3, -6, -5, 34, -12, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	-12, -7, 16, 8, 17, -11, 12, -12, -12, -12,
	-12, -12, -12, -12, -12, -12, -12, -12, -12, 18,
	-9, 10, -5, 9, 18, 34, 13, 37, 15, 35,
	-8, 34, 18, -4, 19, 16, 14, 37, 34, 36,
}
var yyDef = []int{

	0, -2, 5, 0, 0, 1, 4, 2, 7, 8,
	10, 0, 3, 36, 37, 38, 39, 0, 41, 42,
	43, 44, 45, 6, 15, 11, 0, 22, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 19, 0, 0, 0, 0, 0, 24, 25, 26,
	27, 28, 29, 30, 31, 32, 33, 34, 35, 40,
	0, 0, 12, 0, 13, 0, 5, 0, 9, 20,
	16, 17, 14, 0, 0, 0, 21, 0, 18, 23,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37,
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
			l.statement = yyS[yypt-0].statement
		}
	case 3:
		//line grammar.y:79
		{
			l := yylex.(*yylexer)
			l.expression = yyS[yypt-0].expr
		}
	case 4:
		//line grammar.y:87
		{
			l := yylex.(*yylexer)
			l.query.Statements = yyS[yypt-0].statements
			yyVAL.query = l.query
		}
	case 5:
		//line grammar.y:96
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 6:
		//line grammar.y:100
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 7:
		//line grammar.y:107
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 8:
		//line grammar.y:111
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 9:
		//line grammar.y:118
		{
			l := yylex.(*yylexer)
			yyVAL.selection = NewSelection(l.query)
			yyVAL.selection.Fields = yyS[yypt-3].selection_fields
			yyVAL.selection.Dimensions = yyS[yypt-2].strs
			yyVAL.selection.Name = yyS[yypt-1].str
		}
	case 10:
		//line grammar.y:129
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 11:
		//line grammar.y:133
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 12:
		//line grammar.y:138
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 13:
		//line grammar.y:145
		{
			yyVAL.selection_field = NewSelectionField("", yyS[yypt-2].str)
		}
	case 14:
		//line grammar.y:149
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-1].str, yyS[yypt-3].str)
		}
	case 15:
		//line grammar.y:156
		{
			yyVAL.strs = make([]string, 0)
		}
	case 16:
		//line grammar.y:160
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 17:
		//line grammar.y:167
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 18:
		//line grammar.y:172
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 19:
		//line grammar.y:179
		{
			yyVAL.str = ""
		}
	case 20:
		//line grammar.y:183
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 21:
		//line grammar.y:190
		{
			l := yylex.(*yylexer)
			yyVAL.condition = NewCondition(l.query)
			yyVAL.condition.Expression = yyS[yypt-4].expr
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.Statements = yyS[yypt-1].statements
		}
	case 22:
		//line grammar.y:203
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: "steps"}
		}
	case 23:
		//line grammar.y:207
		{
			yyVAL.condition_within = &within{start: yyS[yypt-3].integer, end: yyS[yypt-1].integer, units: yyS[yypt-0].str}
		}
	case 24:
		//line grammar.y:213
		{
			yyVAL.expr = &BinaryExpression{op: OpEquals, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 25:
		//line grammar.y:214
		{
			yyVAL.expr = &BinaryExpression{op: OpNotEquals, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 26:
		//line grammar.y:215
		{
			yyVAL.expr = &BinaryExpression{op: OpLessThan, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 27:
		//line grammar.y:216
		{
			yyVAL.expr = &BinaryExpression{op: OpLessThanOrEqualTo, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 28:
		//line grammar.y:217
		{
			yyVAL.expr = &BinaryExpression{op: OpGreaterThan, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 29:
		//line grammar.y:218
		{
			yyVAL.expr = &BinaryExpression{op: OpGreaterThanOrEqualTo, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 30:
		//line grammar.y:219
		{
			yyVAL.expr = &BinaryExpression{op: OpAnd, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 31:
		//line grammar.y:220
		{
			yyVAL.expr = &BinaryExpression{op: OpOr, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 32:
		//line grammar.y:221
		{
			yyVAL.expr = &BinaryExpression{op: OpPlus, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 33:
		//line grammar.y:222
		{
			yyVAL.expr = &BinaryExpression{op: OpMinus, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 34:
		//line grammar.y:223
		{
			yyVAL.expr = &BinaryExpression{op: OpMultiply, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 35:
		//line grammar.y:224
		{
			yyVAL.expr = &BinaryExpression{op: OpDivide, lhs: yyS[yypt-2].expr, rhs: yyS[yypt-0].expr}
		}
	case 36:
		//line grammar.y:225
		{
			yyVAL.expr = Expression(yyS[yypt-0].var_ref)
		}
	case 37:
		//line grammar.y:226
		{
			yyVAL.expr = Expression(yyS[yypt-0].integer_literal)
		}
	case 38:
		//line grammar.y:227
		{
			yyVAL.expr = Expression(yyS[yypt-0].boolean_literal)
		}
	case 39:
		//line grammar.y:228
		{
			yyVAL.expr = Expression(yyS[yypt-0].string_literal)
		}
	case 40:
		//line grammar.y:229
		{
			yyVAL.expr = yyS[yypt-1].expr
		}
	case 41:
		//line grammar.y:234
		{
			yyVAL.var_ref = &VarRef{value: yyS[yypt-0].str}
		}
	case 42:
		//line grammar.y:241
		{
			yyVAL.integer_literal = &IntegerLiteral{value: yyS[yypt-0].integer}
		}
	case 43:
		//line grammar.y:248
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: true}
		}
	case 44:
		//line grammar.y:252
		{
			yyVAL.boolean_literal = &BooleanLiteral{value: false}
		}
	case 45:
		//line grammar.y:259
		{
			yyVAL.string_literal = &StringLiteral{value: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
