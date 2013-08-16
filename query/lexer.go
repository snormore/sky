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
	case c == '(':
		goto yystate6
	case c == ')':
		goto yystate7
	case c == ',':
		goto yystate8
	case c == ';':
		goto yystate9
	case c == 'A' || c >= 'C' && c <= 'F' || c == 'H' || c >= 'J' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z' || c == '~':
		goto yystate10
	case c == 'B':
		goto yystate11
	case c == 'G':
		goto yystate13
	case c == 'I':
		goto yystate18
	case c == 'S':
		goto yystate22
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule11
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
	goto yyrule9

yystate7:
	c = y.getc()
	goto yyrule10

yystate8:
	c = y.getc()
	goto yyrule8

yystate9:
	c = y.getc()
	goto yyrule7

yystate10:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate11:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'Y':
		goto yystate12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate12:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate13:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'R':
		goto yystate14
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate14:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'O':
		goto yystate15
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate15:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'U':
		goto yystate16
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate16:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'P':
		goto yystate17
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate17:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate18:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'N':
		goto yystate19
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'T':
		goto yystate20
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate20:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'O':
		goto yystate21
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate21:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate22:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'E':
		goto yystate23
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate23:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'L':
		goto yystate24
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate24:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'E':
		goto yystate25
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate25:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'C':
		goto yystate26
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate26:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c == 'T':
		goto yystate27
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate27:
	c = y.getc()
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

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
yyrule6: // [a-zA-Z_~][a-zA-Z0-9_]*
	{
		return y.strtoken(yylval, TIDENT)
	}
yyrule7: // ";"
	{
		return y.token(yylval, TSEMICOLON)
	}
yyrule8: // ","
	{
		return y.token(yylval, TCOMMA)
	}
yyrule9: // "("
	{
		return y.token(yylval, TLPAREN)
	}
yyrule10: // ")"
	{
		return y.token(yylval, TRPAREN)
	}
yyrule11: // [ \t\n\r]+

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
