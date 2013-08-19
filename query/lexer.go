package query

import (
	"bufio"
	"fmt"
	"strconv"
)

type yylexer struct {
	src        *bufio.Reader
	buf        []byte
	empty      bool
	current    byte
	index      int
	startToken int
	err        error
	query      *Query
	statement  Statement
	expression Expression
}

func newLexer(src *bufio.Reader, startToken int) *yylexer {
	y := &yylexer{
		src:        src,
		startToken: startToken,
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

func (y *yylexer) Error(e string) {
	y.err = fmt.Errorf("Unexpected '%c' at character %d, %s", y.current, y.index, e)
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
		goto yystate13
	case c == ')':
		goto yystate14
	case c == '*':
		goto yystate15
	case c == '+':
		goto yystate16
	case c == ',':
		goto yystate17
	case c == '-':
		goto yystate18
	case c == '.':
		goto yystate19
	case c == '/':
		goto yystate21
	case c == ';':
		goto yystate23
	case c == '<':
		goto yystate24
	case c == '=':
		goto yystate26
	case c == '>':
		goto yystate28
	case c == 'A':
		goto yystate30
	case c == 'B':
		goto yystate33
	case c == 'C' || c == 'D' || c == 'F' || c == 'H' || c >= 'J' && c <= 'R' || c == 'U' || c == 'V' || c >= 'X' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 's' || c >= 'u' && c <= 'z' || c == '~':
		goto yystate31
	case c == 'E':
		goto yystate35
	case c == 'G':
		goto yystate38
	case c == 'I':
		goto yystate43
	case c == 'S':
		goto yystate47
	case c == 'T':
		goto yystate61
	case c == 'W':
		goto yystate65
	case c == '\'':
		goto yystate10
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	case c == 'f':
		goto yystate74
	case c == 't':
		goto yystate79
	case c == '|':
		goto yystate83
	case c >= '0' && c <= '9':
		goto yystate22
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule34
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
	goto yyrule17

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
	goto yyrule22

yystate10:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '\'':
		goto yystate11
	case c == '\\':
		goto yystate12
	case c >= '\x01' && c <= '&' || c >= '(' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate10
	}

yystate11:
	c = y.getc()
	goto yyrule2

yystate12:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate10
	}

yystate13:
	c = y.getc()
	goto yyrule32

yystate14:
	c = y.getc()
	goto yyrule33

yystate15:
	c = y.getc()
	goto yyrule26

yystate16:
	c = y.getc()
	goto yyrule24

yystate17:
	c = y.getc()
	goto yyrule31

yystate18:
	c = y.getc()
	goto yyrule25

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate20
	}

yystate20:
	c = y.getc()
	goto yyrule28

yystate21:
	c = y.getc()
	goto yyrule27

yystate22:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9':
		goto yystate22
	}

yystate23:
	c = y.getc()
	goto yyrule30

yystate24:
	c = y.getc()
	switch {
	default:
		goto yyrule19
	case c == '=':
		goto yystate25
	}

yystate25:
	c = y.getc()
	goto yyrule18

yystate26:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '=':
		goto yystate27
	}

yystate27:
	c = y.getc()
	goto yyrule16

yystate28:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c == '=':
		goto yystate29
	}

yystate29:
	c = y.getc()
	goto yyrule20

yystate30:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'S':
		goto yystate32
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate31:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate32:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate33:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'Y':
		goto yystate34
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate34:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate35:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'N':
		goto yystate36
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate36:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'D':
		goto yystate37
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate37:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate38:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'R':
		goto yystate39
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate39:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'O':
		goto yystate40
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate40:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'U':
		goto yystate41
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate41:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'P':
		goto yystate42
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate42:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate43:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'N':
		goto yystate44
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate44:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'T':
		goto yystate45
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate45:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'O':
		goto yystate46
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate46:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate47:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'E':
		goto yystate48
	case c == 'T':
		goto yystate59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate48:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'L':
		goto yystate49
	case c == 'S':
		goto yystate53
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate49:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'E':
		goto yystate50
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate50:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'C':
		goto yystate51
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate51:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'T':
		goto yystate52
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate52:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate53:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'S':
		goto yystate54
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate54:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'I':
		goto yystate55
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate55:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'O':
		goto yystate56
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate56:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'N':
		goto yystate57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate57:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'S':
		goto yystate58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate58:
	c = y.getc()
	switch {
	default:
		goto yyrule13
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate59:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'E':
		goto yystate60
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate60:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'P':
		goto yystate57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate61:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'H':
		goto yystate62
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate62:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'E':
		goto yystate63
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate63:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'N':
		goto yystate64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate64:
	c = y.getc()
	switch {
	default:
		goto yyrule11
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate65:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'H':
		goto yystate66
	case c == 'I':
		goto yystate69
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate66:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'E':
		goto yystate67
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate67:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'N':
		goto yystate68
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate68:
	c = y.getc()
	switch {
	default:
		goto yyrule9
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate69:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'T':
		goto yystate70
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate70:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'H':
		goto yystate71
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate71:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'I':
		goto yystate72
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate72:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'N':
		goto yystate73
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate73:
	c = y.getc()
	switch {
	default:
		goto yyrule10
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate74:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'a':
		goto yystate75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate31
	}

yystate75:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'l':
		goto yystate76
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate31
	}

yystate76:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 's':
		goto yystate77
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate31
	}

yystate77:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'e':
		goto yystate78
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate31
	}

yystate78:
	c = y.getc()
	switch {
	default:
		goto yyrule15
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate79:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'r':
		goto yystate80
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate31
	}

yystate80:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'u':
		goto yystate81
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate31
	}

yystate81:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'e':
		goto yystate82
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate31
	}

yystate82:
	c = y.getc()
	switch {
	default:
		goto yyrule14
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate83:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '|':
		goto yystate84
	}

yystate84:
	c = y.getc()
	goto yyrule23

yyrule1: // \"(\\.|[^\\"])*\"
	{
		return y.quotedstrtoken(yylval, TSTRING)
	}
yyrule2: // \'(\\.|[^\\'])*\'
	{
		return y.quotedstrtoken(yylval, TSTRING)
	}
yyrule3: // [0-9]+
	{
		return y.inttoken(yylval, TINT)
	}
yyrule4: // "SELECT"
	{
		return y.token(yylval, TSELECT)
	}
yyrule5: // "GROUP"
	{
		return y.token(yylval, TGROUP)
	}
yyrule6: // "BY"
	{
		return y.token(yylval, TBY)
	}
yyrule7: // "INTO"
	{
		return y.token(yylval, TINTO)
	}
yyrule8: // "AS"
	{
		return y.token(yylval, TAS)
	}
yyrule9: // "WHEN"
	{
		return y.token(yylval, TWHEN)
	}
yyrule10: // "WITHIN"
	{
		return y.token(yylval, TWITHIN)
	}
yyrule11: // "THEN"
	{
		return y.token(yylval, TTHEN)
	}
yyrule12: // "END"
	{
		return y.token(yylval, TEND)
	}
yyrule13: // (STEPS|SESSIONS)
	{
		return y.strtoken(yylval, TWITHINUNITS)
	}
yyrule14: // "true"
	{
		return y.token(yylval, TTRUE)
	}
yyrule15: // "false"
	{
		return y.token(yylval, TFALSE)
	}
yyrule16: // "=="
	{
		return y.token(yylval, TEQUALS)
	}
yyrule17: // "!="
	{
		return y.token(yylval, TNOTEQUALS)
	}
yyrule18: // "<="
	{
		return y.token(yylval, TLTE)
	}
yyrule19: // "<"
	{
		return y.token(yylval, TLT)
	}
yyrule20: // ">="
	{
		return y.token(yylval, TGTE)
	}
yyrule21: // ">"
	{
		return y.token(yylval, TGT)
	}
yyrule22: // "&&"
	{
		return y.token(yylval, TAND)
	}
yyrule23: // "||"
	{
		return y.token(yylval, TOR)
	}
yyrule24: // "+"
	{
		return y.token(yylval, TPLUS)
	}
yyrule25: // "-"
	{
		return y.token(yylval, TMINUS)
	}
yyrule26: // "*"
	{
		return y.token(yylval, TMUL)
	}
yyrule27: // "/"
	{
		return y.token(yylval, TDIV)
	}
yyrule28: // ".."
	{
		return y.token(yylval, TRANGE)
	}
yyrule29: // [a-zA-Z_~][a-zA-Z0-9_]*
	{
		return y.strtoken(yylval, TIDENT)
	}
yyrule30: // ";"
	{
		return y.token(yylval, TSEMICOLON)
	}
yyrule31: // ","
	{
		return y.token(yylval, TCOMMA)
	}
yyrule32: // "("
	{
		return y.token(yylval, TLPAREN)
	}
yyrule33: // ")"
	{
		return y.token(yylval, TRPAREN)
	}
yyrule34: // [ \t\n\r]+

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
