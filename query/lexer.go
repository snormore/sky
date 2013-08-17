package query

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
)

type yylexer struct {
	src        *bufio.Reader
	buf        []byte
	empty      bool
	current    byte
	index      int
	startToken int
	query      *Query
	statement  Statement
	expression Expression
}

func newLexer(src *bufio.Reader, startToken int) *yylexer {
	y := &yylexer{
		src:        src,
		startToken: startToken,
		query:      &Query{},
	}
	y.current, _ = src.ReadByte()
	return y
}

func (y *yylexer) getc() byte {
	var err error
	if y.current != 0 {
		y.buf = append(y.buf, y.current)
	}

	if y.current, err = y.src.ReadByte(); err == nil {
		y.index++
	}
	return y.current
}

func (y yylexer) Error(e string) {
	log.Fatal(fmt.Sprintf("[%d] Unexpected: '%c', %s", y.index, y.current, e))
}

func (y *yylexer) Lex(yylval *yySymType) int {
	if y.startToken != 0 {
		token := y.startToken
		y.startToken = 0
		return token
	}
	c := y.current
	if y.empty {
		c, y.empty = y.getc(), false
	}

yystate0:

	y.buf = y.buf[:0]

	goto yystart1

	goto yystate1 // silence unused label error
yystate1:
	c = y.getc()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '!':
		goto yystate3
	case c == '"':
		goto yystate5
	case c == '&':
		goto yystate8
	case c == '(':
		goto yystate10
	case c == ')':
		goto yystate11
	case c == '*':
		goto yystate12
	case c == '+':
		goto yystate13
	case c == ',':
		goto yystate14
	case c == '-':
		goto yystate15
	case c == '.':
		goto yystate16
	case c == '/':
		goto yystate18
	case c == ';':
		goto yystate20
	case c == '<':
		goto yystate21
	case c == '=':
		goto yystate23
	case c == '>':
		goto yystate25
	case c == 'A' || c == 'C' || c == 'D' || c == 'F' || c == 'H' || c >= 'J' && c <= 'R' || c == 'U' || c == 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z' || c == '~':
		goto yystate27
	case c == 'B':
		goto yystate28
	case c == 'E':
		goto yystate30
	case c == 'G':
		goto yystate33
	case c == 'I':
		goto yystate38
	case c == 'S':
		goto yystate42
	case c == 'T':
		goto yystate56
	case c == 'W':
		goto yystate60
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	case c == '|':
		goto yystate69
	case c >= '0' && c <= '9':
		goto yystate19
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule30
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	}

yystate3:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '=':
		goto yystate4
	}

yystate4:
	c = y.getc()
	goto yyrule13

yystate5:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate6
	case c == '\\':
		goto yystate7
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate5
	}

yystate6:
	c = y.getc()
	goto yyrule1

yystate7:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate5
	}

yystate8:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '&':
		goto yystate9
	}

yystate9:
	c = y.getc()
	goto yyrule18

yystate10:
	c = y.getc()
	goto yyrule28

yystate11:
	c = y.getc()
	goto yyrule29

yystate12:
	c = y.getc()
	goto yyrule22

yystate13:
	c = y.getc()
	goto yyrule20

yystate14:
	c = y.getc()
	goto yyrule27

yystate15:
	c = y.getc()
	goto yyrule21

yystate16:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate17
	}

yystate17:
	c = y.getc()
	goto yyrule24

yystate18:
	c = y.getc()
	goto yyrule23

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9':
		goto yystate19
	}

yystate20:
	c = y.getc()
	goto yyrule26

yystate21:
	c = y.getc()
	switch {
	default:
		goto yyrule15
	case c == '=':
		goto yystate22
	}

yystate22:
	c = y.getc()
	goto yyrule14

yystate23:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '=':
		goto yystate24
	}

yystate24:
	c = y.getc()
	goto yyrule12

yystate25:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == '=':
		goto yystate26
	}

yystate26:
	c = y.getc()
	goto yyrule16

yystate27:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate28:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'Y':
		goto yystate29
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate29:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate30:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'N':
		goto yystate31
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate31:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'D':
		goto yystate32
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate32:
	c = y.getc()
	switch {
	default:
		goto yyrule10
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate33:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'R':
		goto yystate34
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate34:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'O':
		goto yystate35
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate35:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'U':
		goto yystate36
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate36:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'P':
		goto yystate37
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate37:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate38:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'N':
		goto yystate39
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate39:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'T':
		goto yystate40
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate40:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'O':
		goto yystate41
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate41:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate42:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'E':
		goto yystate43
	case c == 'T':
		goto yystate54
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate43:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'L':
		goto yystate44
	case c == 'S':
		goto yystate48
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate44:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'E':
		goto yystate45
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate45:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'C':
		goto yystate46
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate46:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'T':
		goto yystate47
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate47:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate48:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'S':
		goto yystate49
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate49:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'I':
		goto yystate50
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate50:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'O':
		goto yystate51
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate51:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'N':
		goto yystate52
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate52:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'S':
		goto yystate53
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate53:
	c = y.getc()
	switch {
	default:
		goto yyrule11
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate54:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'E':
		goto yystate55
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate55:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'P':
		goto yystate52
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate56:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'H':
		goto yystate57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate57:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'E':
		goto yystate58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate58:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'N':
		goto yystate59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate59:
	c = y.getc()
	switch {
	default:
		goto yyrule9
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate60:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'H':
		goto yystate61
	case c == 'I':
		goto yystate64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate61:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'E':
		goto yystate62
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate62:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'N':
		goto yystate63
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate63:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate64:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'T':
		goto yystate65
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate65:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'H':
		goto yystate66
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate66:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'I':
		goto yystate67
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate67:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c == 'N':
		goto yystate68
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate68:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate27
	}

yystate69:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '|':
		goto yystate70
	}

yystate70:
	c = y.getc()
	goto yyrule19

yyrule1: // \"(\\.|[^\\"])*\"
	{
		return y.quotedstrtoken(yylval, TSTRING)
	}
yyrule2: // [0-9]+
	{
		return y.inttoken(yylval, TINT)
	}
yyrule3: // "SELECT"
	{
		return y.token(yylval, TSELECT)
	}
yyrule4: // "GROUP"
	{
		return y.token(yylval, TGROUP)
	}
yyrule5: // "BY"
	{
		return y.token(yylval, TBY)
	}
yyrule6: // "INTO"
	{
		return y.token(yylval, TINTO)
	}
yyrule7: // "WHEN"
	{
		return y.token(yylval, TWHEN)
	}
yyrule8: // "WITHIN"
	{
		return y.token(yylval, TWITHIN)
	}
yyrule9: // "THEN"
	{
		return y.token(yylval, TTHEN)
	}
yyrule10: // "END"
	{
		return y.token(yylval, TEND)
	}
yyrule11: // (STEPS|SESSIONS)
	{
		return y.strtoken(yylval, TWITHINUNITS)
	}
yyrule12: // "=="
	{
		return y.token(yylval, TEQUALS)
	}
yyrule13: // "!="
	{
		return y.token(yylval, TNOTEQUALS)
	}
yyrule14: // "<="
	{
		return y.token(yylval, TLTE)
	}
yyrule15: // "<"
	{
		return y.token(yylval, TLT)
	}
yyrule16: // ">="
	{
		return y.token(yylval, TGTE)
	}
yyrule17: // ">"
	{
		return y.token(yylval, TGT)
	}
yyrule18: // "&&"
	{
		return y.token(yylval, TAND)
	}
yyrule19: // "||"
	{
		return y.token(yylval, TOR)
	}
yyrule20: // "+"
	{
		return y.token(yylval, TPLUS)
	}
yyrule21: // "-"
	{
		return y.token(yylval, TMINUS)
	}
yyrule22: // "*"
	{
		return y.token(yylval, TMUL)
	}
yyrule23: // "/"
	{
		return y.token(yylval, TDIV)
	}
yyrule24: // ".."
	{
		return y.token(yylval, TRANGE)
	}
yyrule25: // [a-zA-Z_~][a-zA-Z0-9_]*
	{
		return y.strtoken(yylval, TIDENT)
	}
yyrule26: // ";"
	{
		return y.token(yylval, TSEMICOLON)
	}
yyrule27: // ","
	{
		return y.token(yylval, TCOMMA)
	}
yyrule28: // "("
	{
		return y.token(yylval, TLPAREN)
	}
yyrule29: // ")"
	{
		return y.token(yylval, TRPAREN)
	}
yyrule30: // [ \t\n\r]+

	goto yystate0
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	y.empty = true
	return int(c)
}

// Saves the token to the parser value and returns the token.
func (y *yylexer) token(yylval *yySymType, tok int) int {
	yylval.token = tok
	return tok
}

// Saves the string in the buffer and the token to the parser value
// and returns the token.
func (y *yylexer) strtoken(yylval *yySymType, tok int) int {
	yylval.str = string(y.buf)
	return y.token(yylval, tok)
}

// Saves the quoted string in the buffer and the token to the parser value
// and returns the token.
func (y *yylexer) quotedstrtoken(yylval *yySymType, tok int) int {
	str := string(y.buf)
	yylval.str = str[1 : len(str)-1]
	return y.token(yylval, tok)
}

// Saves the integer in the buffer and the token to the parser value
// and returns the token.
func (y *yylexer) inttoken(yylval *yySymType, tok int) int {
	var err error
	if yylval.integer, err = strconv.Atoi(string(y.buf)); err != nil {
		panic("strconv failed: " + string(y.buf))
	}
	return y.token(yylval, tok)
}
