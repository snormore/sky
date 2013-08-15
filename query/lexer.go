package query

import (
	"bufio"
	"log"
)

type yylexer struct {
	src     *bufio.Reader
	buf     []byte
	empty   bool
	current byte
}

func newLexer(src *bufio.Reader) (y *yylexer) {
	y = &yylexer{src: src}
	if b, err := src.ReadByte(); err == nil {
		y.current = b
	}
	return
}

func (y *yylexer) getc() byte {
	if y.current != 0 {
		y.buf = append(y.buf, y.current)
	}
	y.current = 0
	if b, err := y.src.ReadByte(); err == nil {
		y.current = b
	}
	return y.current
}

func (y yylexer) Error(e string) {
	log.Fatal(e)
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
	case c == 'S':
		goto yystate3
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule2
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	}

yystate3:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 'E':
		goto yystate4
	}

yystate4:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 'L':
		goto yystate5
	}

yystate5:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 'E':
		goto yystate6
	}

yystate6:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 'C':
		goto yystate7
	}

yystate7:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 'T':
		goto yystate8
	}

yystate8:
	c = y.getc()
	goto yyrule1

yyrule1: // "SELECT"
	{
		return TSELECT
	}
yyrule2: // [ \t\n\r]+

	goto yystate0
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	y.empty = true
	return int(c)
}
