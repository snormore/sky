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
	str              string
	strs             []string
	query            *Query
	statement        Statement
	statements       Statements
	selection        *Selection
	selection_field  *SelectionField
	selection_fields []*SelectionField
	condition        *Condition
}

const TSELECT = 57346
const TGROUP = 57347
const TBY = 57348
const TINTO = 57349
const TWHEN = 57350
const TTHEN = 57351
const TEND = 57352
const TSEMICOLON = 57353
const TCOMMA = 57354
const TLPAREN = 57355
const TRPAREN = 57356
const TEQUALS = 57357
const TAND = 57358
const TOR = 57359
const TIDENT = 57360
const TSTRING = 57361

var yyToknames = []string{
	"TSELECT",
	"TGROUP",
	"TBY",
	"TINTO",
	"TWHEN",
	"TTHEN",
	"TEND",
	"TSEMICOLON",
	"TCOMMA",
	"TLPAREN",
	"TRPAREN",
	"TEQUALS",
	"TAND",
	"TOR",
	"TIDENT",
	"TSTRING",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:179

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

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 23
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 39

var yyAct = []int{

	12, 31, 29, 25, 2, 37, 33, 26, 13, 10,
	18, 9, 20, 34, 16, 17, 36, 19, 6, 30,
	28, 15, 7, 27, 35, 6, 22, 23, 24, 7,
	1, 11, 5, 21, 32, 14, 8, 3, 4,
}
var yyPact = []int{

	-1000, -1000, 21, -1000, -1000, -1000, -9, -10, 9, -1000,
	2, 1, -1000, -3, 19, -9, 22, -11, -1000, -10,
	-17, 8, -18, -1000, -12, -1000, -1, 14, -1000, -1000,
	-1000, -1000, 4, -1000, -1000, -1000, -13, -1000,
}
var yyPgo = []int{

	0, 38, 37, 4, 11, 36, 35, 34, 33, 32,
	31, 0, 30,
}
var yyR1 = []int{

	0, 12, 3, 3, 2, 2, 1, 5, 5, 5,
	4, 4, 6, 6, 7, 7, 8, 8, 9, 10,
	10, 10, 11,
}
var yyR2 = []int{

	0, 1, 0, 2, 1, 1, 5, 0, 1, 3,
	3, 4, 0, 3, 1, 3, 0, 2, 5, 0,
	1, 3, 3,
}
var yyChk = []int{

	-1000, -12, -3, -2, -1, -9, 4, 8, -5, -4,
	18, -10, -11, 18, -6, 12, 5, 13, 9, 16,
	15, -8, 7, -4, 6, 14, 18, -3, -11, 19,
	11, 19, -7, 18, 14, 10, 12, 18,
}
var yyDef = []int{

	2, -2, 1, 3, 4, 5, 7, 19, 12, 8,
	0, 0, 20, 0, 16, 0, 0, 0, 2, 0,
	0, 0, 0, 9, 0, 10, 0, 0, 21, 22,
	6, 17, 13, 14, 11, 18, 0, 15,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19,
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
		//line parser.y:47
		{
			l := yylex.(*yylexer)
			l.query.Statements = yyS[yypt-0].statements
		}
	case 2:
		//line parser.y:55
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 3:
		//line parser.y:59
		{
			yyVAL.statements = append(yyS[yypt-1].statements, yyS[yypt-0].statement)
		}
	case 4:
		//line parser.y:66
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 5:
		//line parser.y:70
		{
			yyVAL.statement = Statement(yyS[yypt-0].condition)
		}
	case 6:
		//line parser.y:77
		{
			l := yylex.(*yylexer)
			yyVAL.selection = NewSelection(l.query)
			yyVAL.selection.Fields = yyS[yypt-3].selection_fields
			yyVAL.selection.Dimensions = yyS[yypt-2].strs
			yyVAL.selection.Name = yyS[yypt-1].str
		}
	case 7:
		//line parser.y:88
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 8:
		//line parser.y:92
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 9:
		//line parser.y:97
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 10:
		//line parser.y:104
		{
			yyVAL.selection_field = NewSelectionField("", yyS[yypt-2].str)
		}
	case 11:
		//line parser.y:108
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-1].str, yyS[yypt-3].str)
		}
	case 12:
		//line parser.y:115
		{
			yyVAL.strs = make([]string, 0)
		}
	case 13:
		//line parser.y:119
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 14:
		//line parser.y:126
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 15:
		//line parser.y:131
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	case 16:
		//line parser.y:138
		{
			yyVAL.str = ""
		}
	case 17:
		//line parser.y:142
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 18:
		//line parser.y:149
		{
			l := yylex.(*yylexer)
			yyVAL.condition = NewCondition(l.query)
			yyVAL.condition.Expression = yyS[yypt-3].str
			yyVAL.condition.Statements = yyS[yypt-1].statements
		}
	case 19:
		//line parser.y:159
		{
			yyVAL.str = ""
		}
	case 20:
		//line parser.y:163
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 21:
		//line parser.y:167
		{
			yyVAL.str = yyS[yypt-2].str + " && " + yyS[yypt-0].str
		}
	case 22:
		//line parser.y:174
		{
			yyVAL.str = yyS[yypt-2].str + " == \"" + yyS[yypt-0].str + "\""
		}
	}
	goto yystack /* stack new state and value */
}
