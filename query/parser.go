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
}

const TSELECT = 57346
const TGROUP = 57347
const TBY = 57348
const TSEMICOLON = 57349
const TCOMMA = 57350
const TLPAREN = 57351
const TRPAREN = 57352
const TIDENT = 57353

var yyToknames = []string{
	"TSELECT",
	"TGROUP",
	"TBY",
	"TSEMICOLON",
	"TCOMMA",
	"TLPAREN",
	"TRPAREN",
	"TIDENT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:126

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

const yyNprod = 16
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 24

var yyAct = []int{

	9, 18, 19, 24, 21, 10, 22, 15, 14, 23,
	3, 13, 7, 6, 16, 17, 5, 11, 1, 20,
	12, 8, 2, 4,
}
var yyPact = []int{

	12, -1000, 6, 5, -1000, -6, 12, -1000, 3, -1000,
	-2, -1000, -1000, -6, 9, -9, -1000, -7, -1000, -4,
	1, -1000, -1000, -8, -1000,
}
var yyPgo = []int{

	0, 23, 10, 22, 0, 21, 20, 19, 18,
}
var yyR1 = []int{

	0, 8, 3, 3, 3, 2, 1, 5, 5, 5,
	4, 4, 6, 6, 7, 7,
}
var yyR2 = []int{

	0, 1, 0, 2, 3, 1, 3, 0, 1, 3,
	3, 4, 0, 3, 1, 3,
}
var yyChk = []int{

	-1000, -8, -3, -2, -1, 4, 7, 7, -5, -4,
	11, -2, -6, 8, 5, 9, -4, 6, 10, 11,
	-7, 11, 10, 8, 11,
}
var yyDef = []int{

	2, -2, 1, 0, 5, 7, 0, 3, 12, 8,
	0, 4, 6, 0, 0, 0, 9, 0, 10, 0,
	13, 14, 11, 0, 15,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
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
		//line parser.y:40
		{
			l := yylex.(*yylexer)
			l.query.Statements = yyS[yypt-0].statements
		}
	case 2:
		//line parser.y:48
		{
			yyVAL.statements = make(Statements, 0)
		}
	case 3:
		//line parser.y:52
		{
			yyVAL.statements = make(Statements, 0)
			yyVAL.statements = append(yyVAL.statements, yyS[yypt-1].statement)
		}
	case 4:
		//line parser.y:57
		{
			yyVAL.statements = append(yyS[yypt-2].statements, yyS[yypt-0].statement)
		}
	case 5:
		//line parser.y:63
		{
			yyVAL.statement = Statement(yyS[yypt-0].selection)
		}
	case 6:
		//line parser.y:68
		{
			l := yylex.(*yylexer)
			yyVAL.selection = NewSelection(l.query)
			yyVAL.selection.Fields = yyS[yypt-1].selection_fields
			yyVAL.selection.Dimensions = yyS[yypt-0].strs
		}
	case 7:
		//line parser.y:78
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
		}
	case 8:
		//line parser.y:82
		{
			yyVAL.selection_fields = make([]*SelectionField, 0)
			yyVAL.selection_fields = append(yyVAL.selection_fields, yyS[yypt-0].selection_field)
		}
	case 9:
		//line parser.y:87
		{
			yyVAL.selection_fields = append(yyS[yypt-2].selection_fields, yyS[yypt-0].selection_field)
		}
	case 10:
		//line parser.y:94
		{
			yyVAL.selection_field = NewSelectionField("", yyS[yypt-2].str)
		}
	case 11:
		//line parser.y:98
		{
			yyVAL.selection_field = NewSelectionField(yyS[yypt-1].str, yyS[yypt-3].str)
		}
	case 12:
		//line parser.y:105
		{
			yyVAL.strs = make([]string, 0)
		}
	case 13:
		//line parser.y:109
		{
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 14:
		//line parser.y:116
		{
			yyVAL.strs = make([]string, 0)
			yyVAL.strs = append(yyVAL.strs, yyS[yypt-0].str)
		}
	case 15:
		//line parser.y:121
		{
			yyVAL.strs = append(yyS[yypt-2].strs, yyS[yypt-0].str)
		}
	}
	goto yystack /* stack new state and value */
}
