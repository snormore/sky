package query

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
)

type yylexer struct {
	src     *bufio.Reader
	buf     []byte
	empty   bool
	current byte
	index   int
	query   *Query
}

func newLexer(src *bufio.Reader) *yylexer {
	y := &yylexer{src: src, query: &Query{}}
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
	log.Fatal(fmt.Sprintf("[%d] %s", y.index, e))
}

func (y *yylexer) Lex(yylval *yySymType) int {
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
	case c == ',':
		goto yystate12
	case c == '.':
		goto yystate13
	case c == ';':
		goto yystate16
	case c == '<':
		goto yystate17
	case c == '=':
		goto yystate19
	case c == '>':
		goto yystate21
	case c == 'A' || c == 'C' || c == 'D' || c == 'F' || c == 'H' || c >= 'J' && c <= 'R' || c == 'U' || c == 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z' || c == '~':
		goto yystate23
	case c == 'B':
		goto yystate24
	case c == 'E':
		goto yystate26
	case c == 'G':
		goto yystate29
	case c == 'I':
		goto yystate34
	case c == 'S':
		goto yystate38
	case c == 'T':
		goto yystate52
	case c == 'W':
		goto yystate56
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	case c == '|':
		goto yystate65
	case c >= '0' && c <= '9':
		goto yystate15
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule26
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
	goto yyrule24

yystate11:
	c = y.getc()
	goto yyrule25

yystate12:
	c = y.getc()
	goto yyrule23

yystate13:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate14
	}

yystate14:
	c = y.getc()
	goto yyrule20

yystate15:
	c = y.getc()
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9':
		goto yystate15
	}

yystate16:
	c = y.getc()
	goto yyrule22

yystate17:
	c = y.getc()
	switch {
	default:
		goto yyrule15
	case c == '=':
		goto yystate18
	}

yystate18:
	c = y.getc()
	goto yyrule14

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '=':
		goto yystate20
	}

yystate20:
	c = y.getc()
	goto yyrule12

yystate21:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == '=':
		goto yystate22
	}

yystate22:
	c = y.getc()
	goto yyrule16

yystate23:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate24:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'Y':
		goto yystate25
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate25:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate26:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'N':
		goto yystate27
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate27:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'D':
		goto yystate28
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate28:
	c = y.getc()
	switch {
	default:
		goto yyrule10
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate29:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'R':
		goto yystate30
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate30:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'O':
		goto yystate31
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate31:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'U':
		goto yystate32
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate32:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'P':
		goto yystate33
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate33:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate34:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'N':
		goto yystate35
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate35:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'T':
		goto yystate36
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate36:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'O':
		goto yystate37
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate37:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate38:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'E':
		goto yystate39
	case c == 'T':
		goto yystate50
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate39:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'L':
		goto yystate40
	case c == 'S':
		goto yystate44
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate40:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'E':
		goto yystate41
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate41:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'C':
		goto yystate42
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate42:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'T':
		goto yystate43
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate43:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate44:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'S':
		goto yystate45
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate45:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'I':
		goto yystate46
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate46:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'O':
		goto yystate47
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate47:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'N':
		goto yystate48
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate48:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'S':
		goto yystate49
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate49:
	c = y.getc()
	switch {
	default:
		goto yyrule11
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate50:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'E':
		goto yystate51
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate51:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'P':
		goto yystate48
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate52:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'H':
		goto yystate53
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate53:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'E':
		goto yystate54
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate54:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'N':
		goto yystate55
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate55:
	c = y.getc()
	switch {
	default:
		goto yyrule9
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate56:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'H':
		goto yystate57
	case c == 'I':
		goto yystate60
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate57:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'E':
		goto yystate58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate58:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'N':
		goto yystate59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate59:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate60:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'T':
		goto yystate61
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate61:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'H':
		goto yystate62
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate62:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'I':
		goto yystate63
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate63:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == 'N':
		goto yystate64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate64:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate23
	}

yystate65:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '|':
		goto yystate66
	}

yystate66:
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
yyrule20: // ".."
	{
		return y.token(yylval, TRANGE)
	}
yyrule21: // [a-zA-Z_~][a-zA-Z0-9_]*
	{
		return y.strtoken(yylval, TIDENT)
	}
yyrule22: // ";"
	{
		return y.token(yylval, TSEMICOLON)
	}
yyrule23: // ","
	{
		return y.token(yylval, TCOMMA)
	}
yyrule24: // "("
	{
		return y.token(yylval, TLPAREN)
	}
yyrule25: // ")"
	{
		return y.token(yylval, TRPAREN)
	}
yyrule26: // [ \t\n\r]+

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
