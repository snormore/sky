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
	case c == '(':
		goto yystate3
	case c == ')':
		goto yystate4
	case c == ',':
		goto yystate5
	case c == ';':
		goto yystate6
	case c == 'A' || c >= 'C' && c <= 'F' || c >= 'H' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z' || c == '~':
		goto yystate7
	case c == 'B':
		goto yystate8
	case c == 'G':
		goto yystate10
	case c == 'S':
		goto yystate15
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule9
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	}

yystate3:
	c = y.getc()
	goto yyrule7

yystate4:
	c = y.getc()
	goto yyrule8

yystate5:
	c = y.getc()
	goto yyrule6

yystate6:
	c = y.getc()
	goto yyrule5

yystate7:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate8:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'Y':
		goto yystate9
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate9:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate10:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'R':
		goto yystate11
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate11:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'O':
		goto yystate12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate12:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'U':
		goto yystate13
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate13:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'P':
		goto yystate14
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate14:
	c = y.getc()
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate15:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'E':
		goto yystate16
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate16:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'L':
		goto yystate17
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate17:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'E':
		goto yystate18
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate18:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'C':
		goto yystate19
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c == 'T':
		goto yystate20
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yystate20:
	c = y.getc()
	switch {
	default:
		goto yyrule1
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate7
	}

yyrule1: // "SELECT"
	{
		return y.token(yylval, TSELECT)
	}
yyrule2: // "GROUP"
	{
		return y.token(yylval, TGROUP)
	}
yyrule3: // "BY"
	{
		return y.token(yylval, TBY)
	}
yyrule4: // [a-zA-Z_~][a-zA-Z0-9_]*
	{
		return y.strtoken(yylval, TIDENT)
	}
yyrule5: // ";"
	{
		return y.token(yylval, TSEMICOLON)
	}
yyrule6: // ","
	{
		return y.token(yylval, TCOMMA)
	}
yyrule7: // "("
	{
		return y.token(yylval, TLPAREN)
	}
yyrule8: // ")"
	{
		return y.token(yylval, TRPAREN)
	}
yyrule9: // [ \t\n\r]+

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
