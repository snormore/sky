//line parser.y:2
package query

import __yyfmt__ "fmt"

//line parser.y:3
import (
	"bufio"
	"bytes"
	"io"
)

//line parser.y:13
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
}

const TSELECT = 57346
const TGROUP = 57347
const TBY = 57348
const TINTO = 57349
const TWHEN = 57350
const TWITHIN = 57351
const TTHEN = 57352
const TEND = 57353
const TSEMICOLON = 57354
const TCOMMA = 57355
const TLPAREN = 57356
const TRPAREN = 57357
const TRANGE = 57358
const TEQUALS = 57359
const TAND = 57360
const TOR = 57361
const TIDENT = 57362
const TSTRING = 57363
const TWITHINUNITS = 57364
const TINT = 57365

var yyToknames = []string{
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
	"TAND",
	"TOR",
	"TIDENT",
	"TSTRING",
	"TWITHINUNITS",
	"TINT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:197

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader) *Query {
	l := newLexer(bufio.NewReader(r))
	yyParse(l)
	return l.query
}

func (p *Parser) ParseString(s string) *Query {
	return p.Parse(bytes.NewBufferString(s))
}

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

const yyNprod = 25
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 45

var yyAct = []int{

	2, 41, 30, 43, 33, 31, 26, 42, 35, 13,
	10, 27, 12, 20, 21, 38, 36, 9, 16, 17,
	6, 39, 19, 32, 7, 28, 15, 40, 23, 37,
	6, 25, 29, 24, 7, 1, 18, 11, 5, 22,
	34, 14, 8, 3, 4,
}
var yyPact = []int{

	-1000, -1000, 26, -1000, -1000, -1000, -10, -11, 13, -1000,
	5, 4, -1000, -3, 21, -10, 25, -9, 15, -11,
	-21, -16, 11, -17, -1000, -12, -1000, 1, -1000, -1000,
	-1, -1000, -1000, -1000, 8, -1000, -1000, 16, -22, -13,
	-1000, -19, -1000, -1000,
}
var yyPgo = []int{

	0, 44, 43, 0, 17, 42, 41, 40, 39, 38,
	37, 12, 36, 35,
}
var yyR1 = []int{

	0, 13, 3, 3, 2, 2, 1, 5, 5, 5,
	4, 4, 6, 6, 7, 7, 8, 8, 9, 10,
	10, 10, 11, 12, 12,
}
var yyR2 = []int{

	0, 1, 0, 2, 1, 1, 5, 0, 1, 3,
	3, 4, 0, 3, 1, 3, 0, 2, 6, 0,
	1, 3, 3, 0, 5,
}
var yyChk = []int{

	-1000, -13, -3, -2, -1, -9, 4, 8, -5, -4,
	20, -10, -11, 20, -6, 13, 5, 14, -12, 18,
	9, 17, -8, 7, -4, 6, 15, 20, 10, -11,
	23, 21, 12, 21, -7, 20, 15, -3, 16, 13,
	11, 23, 20, 22,
}
var yyDef = []int{

	2, -2, 1, 3, 4, 5, 7, 19, 12, 8,
	0, 23, 20, 0, 16, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 9, 0, 10, 0, 2, 21,
	0, 22, 6, 17, 13, 14, 11, 0, 0, 0,
	18, 0, 15, 24,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23,
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
		//line parser.y:51
		{
			l := yylex.(*yylexer)
			l.query.Statements = yyS[yypt-0].statements
		}
	case 2:
		//line parser.y:59
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 3:
		//line parser.y:63
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 4:
		//line parser.y:70
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 5:
		//line parser.y:74
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 6:
		//line parser.y:81
		{
			l := yylex.(*yylexer)
			yyVAL.selection = NewSelection(l.query)
			yyVAL.selection.Fields = yyS[yypt-3].selection_fields
			yyVAL.selection.Dimensions = yyS[yypt-2].strs
			yyVAL.selection.Name = yyS[yypt-1].str
		}
	case 7:
		//line parser.y:92
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 8:
		//line parser.y:96
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 9:
		//line parser.y:101
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 10:
		//line parser.y:108
		{
			yyVAL.selection_field = NewSelectionField("", yyS[yypt-2].str)
		}
	case 11:
		//line parser.y:112
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-1].str, yyS[yypt-3].str)
		}
	case 12:
		//line parser.y:119
		{
			yyVAL.strs = make([]string, 0)
		}
	case 13:
		//line parser.y:123
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 14:
		//line parser.y:130
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 15:
		//line parser.y:135
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 16:
		//line parser.y:142
		{
			yyVAL.str = ""
		}
	case 17:
		//line parser.y:146
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 18:
		//line parser.y:153
		{
			l := yylex.(*yylexer)
			yyVAL.condition = NewCondition(l.query)
			yyVAL.condition.Expression = yyS[yypt-4].str
			yyVAL.condition.WithinRangeStart = yyS[yypt-3].condition_within.start
			yyVAL.condition.WithinRangeEnd = yyS[yypt-3].condition_within.end
			yyVAL.condition.WithinUnits = yyS[yypt-3].condition_within.units
			yyVAL.condition.Statements = yyS[yypt-1].statements
		}
	case 19:
		//line parser.y:166
		{
			yyVAL.str = ""
		}
	case 20:
		//line parser.y:170
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 21:
		//line parser.y:174
		{
			yyVAL.str = yyS[yypt-2].str + " && " + yyS[yypt-0].str
		}
	case 22:
		//line parser.y:181
		{
			yyVAL.str = yyS[yypt-2].str + " == \"" + yyS[yypt-0].str + "\""
		}
	case 23:
		//line parser.y:188
		{
			yyVAL.condition_within = &within{start: 0, end: 0, units: "steps"}
		}
	case 24:
		//line parser.y:192
		{
			yyVAL.condition_within = &within{start: yyS[yypt-3].integer, end: yyS[yypt-1].integer, units: yyS[yypt-0].str}
		}
	}
	goto yystack /* stack new state and value */
}
