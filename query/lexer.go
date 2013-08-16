package query

import (
	"bufio"
	"fmt"
	"log"
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
	case c == '"':
		goto yystate3
	case c == '&':
		goto yystate6
	case c == '(':
		goto yystate8
	case c == ')':
		goto yystate9
	case c == ',':
		goto yystate10
	case c == ';':
		goto yystate11
	case c == '=':
		goto yystate12
	case c == 'A' || c == 'C' || c == 'D' || c == 'F' || c == 'H' || c >= 'J' && c <= 'R' || c == 'U' || c == 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z' || c == '~':
		goto yystate14
	case c == 'B':
		goto yystate15
	case c == 'E':
		goto yystate17
	case c == 'G':
		goto yystate20
	case c == 'I':
		goto yystate25
	case c == 'S':
		goto yystate29
	case c == 'T':
		goto yystate35
	case c == 'W':
		goto yystate39
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	case c == '|':
		goto yystate43
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	}

yystate3:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate4
	case c == '\\':
		goto yystate5
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate3
	}

yystate4:
	c = y.getc()
	goto yyrule1

yystate5:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate3
	}

yystate6:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '&':
		goto yystate7
	}

yystate7:
	c = y.getc()
	goto yyrule9

yystate8:
	c = y.getc()
	goto yyrule15

yystate9:
	c = y.getc()
	goto yyrule16

yystate10:
	c = y.getc()
	goto yyrule14

yystate11:
	c = y.getc()
	goto yyrule13

yystate12:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '=':
		goto yystate13
	}

yystate13:
	c = y.getc()
	goto yyrule11

yystate14:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate15:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'Y':
		goto yystate16
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate16:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate17:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'N':
		goto yystate18
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate18:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'D':
		goto yystate19
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate20:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'R':
		goto yystate21
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate21:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'O':
		goto yystate22
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate22:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'U':
		goto yystate23
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate23:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'P':
		goto yystate24
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate24:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate25:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'N':
		goto yystate26
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate26:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'T':
		goto yystate27
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate27:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'O':
		goto yystate28
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate28:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate29:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'E':
		goto yystate30
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate30:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'L':
		goto yystate31
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate31:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'E':
		goto yystate32
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate32:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'C':
		goto yystate33
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate33:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'T':
		goto yystate34
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate34:
	c = y.getc()
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate35:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'H':
		goto yystate36
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate36:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'E':
		goto yystate37
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate37:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'N':
		goto yystate38
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate38:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate39:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'H':
		goto yystate40
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate40:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'E':
		goto yystate41
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate41:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c == 'N':
		goto yystate42
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate42:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate14
	}

yystate43:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '|':
		goto yystate44
	}

yystate44:
	c = y.getc()
	goto yyrule10

yyrule1: // \"(\\.|[^\\"])*\"
	{
		return y.quotedstrtoken(yylval, TSTRING)
	}
yyrule2: // "SELECT"
	{
		return y.token(yylval, TSELECT)
	}
yyrule3: // "GROUP"
	{
		return y.token(yylval, TGROUP)
	}
yyrule4: // "BY"
	{
		return y.token(yylval, TBY)
	}
yyrule5: // "INTO"
	{
		return y.token(yylval, TINTO)
	}
yyrule6: // "WHEN"
	{
		return y.token(yylval, TWHEN)
	}
yyrule7: // "THEN"
	{
		return y.token(yylval, TTHEN)
	}
yyrule8: // "END"
	{
		return y.token(yylval, TEND)
	}
yyrule9: // "&&"
	{
		return y.token(yylval, TAND)
	}
yyrule10: // "||"
	{
		return y.token(yylval, TOR)
	}
yyrule11: // "=="
	{
		return y.token(yylval, TEQUALS)
	}
yyrule12: // [a-zA-Z_~][a-zA-Z0-9_]*
	{
		return y.strtoken(yylval, TIDENT)
	}
yyrule13: // ";"
	{
		return y.token(yylval, TSEMICOLON)
	}
yyrule14: // ","
	{
		return y.token(yylval, TCOMMA)
	}
yyrule15: // "("
	{
		return y.token(yylval, TLPAREN)
	}
yyrule16: // ")"
	{
		return y.token(yylval, TRPAREN)
	}
yyrule17: // [ \t\n\r]+

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
