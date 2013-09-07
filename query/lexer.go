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
	lineidx    int
	charidx    int
	tlineidx   int
	tcharidx   int
	startToken int
	err        error
	query      *Query
	statement  Statement
	statements Statements
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
		y.charidx++

		// Reset line and character index at "\n"
		if y.current == 10 {
			y.lineidx++
			y.charidx = 0
		}
	}
	return y.current
}

func (y *yylexer) Error(e string) {
	y.err = fmt.Errorf("Unexpected '%s' at line %d, char %d, %s", y.buf, y.tlineidx+1, y.tcharidx+1, e)
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

	y.tlineidx, y.tcharidx = y.lineidx, y.charidx
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
	case c == 'C' || c >= 'J' && c <= 'L' || c >= 'N' && c <= 'R' || c == 'U' || c == 'V' || c == 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 's' || c >= 'u' && c <= 'z' || c == '~':
		goto yystate31
	case c == 'D':
		goto yystate41
	case c == 'E':
		goto yystate61
	case c == 'F':
		goto yystate76
	case c == 'G':
		goto yystate88
	case c == 'H':
		goto yystate93
	case c == 'I':
		goto yystate96
	case c == 'M':
		goto yystate104
	case c == 'S':
		goto yystate109
	case c == 'T':
		goto yystate132
	case c == 'W':
		goto yystate136
	case c == 'Y':
		goto yystate147
	case c == '\'':
		goto yystate10
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate2
	case c == 'f':
		goto yystate149
	case c == 't':
		goto yystate154
	case c == '|':
		goto yystate158
	case c >= '0' && c <= '9':
		goto yystate22
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule52
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
	goto yyrule34

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
	goto yyrule39

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
	goto yyrule48

yystate14:
	c = y.getc()
	goto yyrule49

yystate15:
	c = y.getc()
	goto yyrule43

yystate16:
	c = y.getc()
	goto yyrule41

yystate17:
	c = y.getc()
	goto yyrule47

yystate18:
	c = y.getc()
	goto yyrule42

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
	goto yyrule45

yystate21:
	c = y.getc()
	goto yyrule44

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
	goto yyrule46

yystate24:
	c = y.getc()
	switch {
	default:
		goto yyrule36
	case c == '=':
		goto yystate25
	}

yystate25:
	c = y.getc()
	goto yyrule35

yystate26:
	c = y.getc()
	switch {
	default:
		goto yyrule50
	case c == '=':
		goto yystate27
	}

yystate27:
	c = y.getc()
	goto yyrule33

yystate28:
	c = y.getc()
	switch {
	default:
		goto yyrule38
	case c == '=':
		goto yystate29
	}

yystate29:
	c = y.getc()
	goto yyrule37

yystate30:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'S':
		goto yystate32
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate31:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate32:
	c = y.getc()
	switch {
	default:
		goto yyrule12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate33:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate34
	case c == 'Y':
		goto yystate40
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate34:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate35
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate35:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'L':
		goto yystate36
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate36:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate37
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate37:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'A':
		goto yystate38
	case c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate38:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate39
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate39:
	c = y.getc()
	switch {
	default:
		goto yyrule28
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate40:
	c = y.getc()
	switch {
	default:
		goto yyrule10
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate41:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'A':
		goto yystate42
	case c == 'E':
		goto yystate45
	case c >= '0' && c <= '9' || c >= 'B' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate42:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'Y':
		goto yystate43
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate43:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c == 'S':
		goto yystate44
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate44:
	c = y.getc()
	switch {
	default:
		goto yyrule29
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate45:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'B':
		goto yystate46
	case c == 'C':
		goto yystate49
	case c == 'L':
		goto yystate54
	case c >= '0' && c <= '9' || c == 'A' || c >= 'D' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate46:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'U':
		goto yystate47
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate47:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'G':
		goto yystate48
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate48:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate49:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'L':
		goto yystate50
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'K' || c >= 'M' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate50:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'A':
		goto yystate51
	case c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate51:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'R':
		goto yystate52
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate52:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate53
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate53:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate54:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'I':
		goto yystate55
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate55:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'M':
		goto yystate56
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'L' || c >= 'N' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate56:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'I':
		goto yystate57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate57:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate58:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate59:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'D':
		goto yystate60
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate60:
	c = y.getc()
	switch {
	default:
		goto yyrule23
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate61:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'A':
		goto yystate62
	case c == 'N':
		goto yystate65
	case c == 'V':
		goto yystate67
	case c == 'X':
		goto yystate73
	case c >= '0' && c <= '9' || c >= 'B' && c <= 'M' || c >= 'O' && c <= 'U' || c == 'W' || c == 'Y' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate62:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'C':
		goto yystate63
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate63:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'H':
		goto yystate64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate64:
	c = y.getc()
	switch {
	default:
		goto yyrule18
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate65:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'D':
		goto yystate66
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate66:
	c = y.getc()
	switch {
	default:
		goto yyrule16
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate67:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate68
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate68:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate69
	case c == 'R':
		goto yystate71
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate69:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate70
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate70:
	c = y.getc()
	switch {
	default:
		goto yyrule21
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate71:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'Y':
		goto yystate72
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'X' || c == 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate72:
	c = y.getc()
	switch {
	default:
		goto yyrule19
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate73:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'I':
		goto yystate74
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate74:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate75:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate76:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'A':
		goto yystate77
	case c == 'L':
		goto yystate82
	case c == 'O':
		goto yystate86
	case c >= '0' && c <= '9' || c >= 'B' && c <= 'K' || c == 'M' || c == 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate77:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'C':
		goto yystate78
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate78:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate79
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate79:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate80
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate80:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'R':
		goto yystate81
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate81:
	c = y.getc()
	switch {
	default:
		goto yyrule24
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate82:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate83
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate83:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'A':
		goto yystate84
	case c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate84:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate85
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate85:
	c = y.getc()
	switch {
	default:
		goto yyrule27
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate86:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'R':
		goto yystate87
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate87:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate88:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'R':
		goto yystate89
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate89:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate90
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate90:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'U':
		goto yystate91
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate91:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'P':
		goto yystate92
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate92:
	c = y.getc()
	switch {
	default:
		goto yyrule9
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate93:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate94
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate94:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'U':
		goto yystate95
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate95:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'R':
		goto yystate43
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate96:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate97
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate97:
	c = y.getc()
	switch {
	default:
		goto yyrule20
	case c == 'T':
		goto yystate98
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate98:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate99
	case c == 'O':
		goto yystate103
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate99:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'G':
		goto yystate100
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate100:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate101
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate101:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'R':
		goto yystate102
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate102:
	c = y.getc()
	switch {
	default:
		goto yyrule26
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate103:
	c = y.getc()
	switch {
	default:
		goto yyrule11
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate104:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'I':
		goto yystate105
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate105:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate106
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate106:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'U':
		goto yystate107
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'T' || c >= 'V' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate107:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate108
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate108:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate43
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate109:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate110
	case c == 'T':
		goto yystate124
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate110:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'C':
		goto yystate111
	case c == 'L':
		goto yystate114
	case c == 'S':
		goto yystate118
	case c == 'T':
		goto yystate123
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'K' || c >= 'M' && c <= 'R' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate111:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate112
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate112:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate113
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate113:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'D':
		goto yystate43
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'C' || c >= 'E' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate114:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate115
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate115:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'C':
		goto yystate116
	case c >= '0' && c <= '9' || c == 'A' || c == 'B' || c >= 'D' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate116:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate117
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate117:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate118:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'S':
		goto yystate119
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate119:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'I':
		goto yystate120
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate120:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'O':
		goto yystate121
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'N' || c >= 'P' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate121:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate122
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate122:
	c = y.getc()
	switch {
	default:
		goto yyrule22
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate123:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate124:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate125
	case c == 'R':
		goto yystate128
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Q' || c >= 'S' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate125:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'P':
		goto yystate126
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'O' || c >= 'Q' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate126:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'S':
		goto yystate127
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'R' || c >= 'T' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate127:
	c = y.getc()
	switch {
	default:
		goto yyrule30
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate128:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'I':
		goto yystate129
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate129:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate130
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate130:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'G':
		goto yystate131
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'H' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate131:
	c = y.getc()
	switch {
	default:
		goto yyrule25
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate132:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'H':
		goto yystate133
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate133:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate134
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate134:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate135
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate135:
	c = y.getc()
	switch {
	default:
		goto yyrule15
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate136:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate137
	case c == 'H':
		goto yystate139
	case c == 'I':
		goto yystate142
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c == 'F' || c == 'G' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate137:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate138
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate138:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'K':
		goto yystate43
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'J' || c >= 'L' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate139:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate140
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate140:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate141
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate141:
	c = y.getc()
	switch {
	default:
		goto yyrule13
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate142:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'T':
		goto yystate143
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'S' || c >= 'U' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate143:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'H':
		goto yystate144
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'G' || c >= 'I' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate144:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'I':
		goto yystate145
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'H' || c >= 'J' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate145:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'N':
		goto yystate146
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'M' || c >= 'O' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate146:
	c = y.getc()
	switch {
	default:
		goto yyrule14
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate147:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'E':
		goto yystate148
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'D' || c >= 'F' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate148:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'A':
		goto yystate95
	case c >= '0' && c <= '9' || c >= 'B' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate149:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'a':
		goto yystate150
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate31
	}

yystate150:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'l':
		goto yystate151
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate31
	}

yystate151:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 's':
		goto yystate152
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate31
	}

yystate152:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'e':
		goto yystate153
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate31
	}

yystate153:
	c = y.getc()
	switch {
	default:
		goto yyrule32
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate154:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'r':
		goto yystate155
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate31
	}

yystate155:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'u':
		goto yystate156
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate31
	}

yystate156:
	c = y.getc()
	switch {
	default:
		goto yyrule51
	case c == 'e':
		goto yystate157
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate31
	}

yystate157:
	c = y.getc()
	switch {
	default:
		goto yyrule31
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate158:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '|':
		goto yystate159
	}

yystate159:
	c = y.getc()
	goto yyrule40

yyrule1: // \"(\\.|[^\\"])*\"
	{
		return y.quotedstrtoken(yylval, TQUOTEDSTRING)
	}
yyrule2: // \'(\\.|[^\\'])*\'
	{
		return y.quotedstrtoken(yylval, TQUOTEDSTRING)
	}
yyrule3: // [0-9]+
	{
		return y.inttoken(yylval, TINT)
	}
yyrule4: // "DECLARE"
	{
		return y.token(yylval, TDECLARE)
	}
yyrule5: // "SET"
	{
		return y.token(yylval, TSET)
	}
yyrule6: // "EXIT"
	{
		return y.token(yylval, TEXIT)
	}
yyrule7: // "DEBUG"
	{
		return y.token(yylval, TDEBUG)
	}
yyrule8: // "SELECT"
	{
		return y.token(yylval, TSELECT)
	}
yyrule9: // "GROUP"
	{
		return y.token(yylval, TGROUP)
	}
yyrule10: // "BY"
	{
		return y.token(yylval, TBY)
	}
yyrule11: // "INTO"
	{
		return y.token(yylval, TINTO)
	}
yyrule12: // "AS"
	{
		return y.token(yylval, TAS)
	}
yyrule13: // "WHEN"
	{
		return y.token(yylval, TWHEN)
	}
yyrule14: // "WITHIN"
	{
		return y.token(yylval, TWITHIN)
	}
yyrule15: // "THEN"
	{
		return y.token(yylval, TTHEN)
	}
yyrule16: // "END"
	{
		return y.token(yylval, TEND)
	}
yyrule17: // "FOR"
	{
		return y.token(yylval, TFOR)
	}
yyrule18: // "EACH"
	{
		return y.token(yylval, TEACH)
	}
yyrule19: // "EVERY"
	{
		return y.token(yylval, TEVERY)
	}
yyrule20: // "IN"
	{
		return y.token(yylval, TIN)
	}
yyrule21: // "EVENT"
	{
		return y.token(yylval, TEVENT)
	}
yyrule22: // "SESSION"
	{
		return y.token(yylval, TSESSION)
	}
yyrule23: // "DELIMITED"
	{
		return y.token(yylval, TDELIMITED)
	}
yyrule24: // "FACTOR"
	{
		return y.token(yylval, TFACTOR)
	}
yyrule25: // "STRING"
	{
		return y.token(yylval, TSTRING)
	}
yyrule26: // "INTEGER"
	{
		return y.token(yylval, TINTEGER)
	}
yyrule27: // "FLOAT"
	{
		return y.token(yylval, TFLOAT)
	}
yyrule28: // "BOOLEAN"
	{
		return y.token(yylval, TBOOLEAN)
	}
yyrule29: // (SECOND|SECONDS|MINUTE|MINUTES|HOUR|HOURS|DAY|DAYS|WEEK|WEEKS|YEAR|YEARS)
	{
		return y.strtoken(yylval, TTIMEUNITS)
	}
yyrule30: // (STEPS)
	{
		return y.strtoken(yylval, TWITHINUNITS)
	}
yyrule31: // "true"
	{
		return y.token(yylval, TTRUE)
	}
yyrule32: // "false"
	{
		return y.token(yylval, TFALSE)
	}
yyrule33: // "=="
	{
		return y.token(yylval, TEQUALS)
	}
yyrule34: // "!="
	{
		return y.token(yylval, TNOTEQUALS)
	}
yyrule35: // "<="
	{
		return y.token(yylval, TLTE)
	}
yyrule36: // "<"
	{
		return y.token(yylval, TLT)
	}
yyrule37: // ">="
	{
		return y.token(yylval, TGTE)
	}
yyrule38: // ">"
	{
		return y.token(yylval, TGT)
	}
yyrule39: // "&&"
	{
		return y.token(yylval, TAND)
	}
yyrule40: // "||"
	{
		return y.token(yylval, TOR)
	}
yyrule41: // "+"
	{
		return y.token(yylval, TPLUS)
	}
yyrule42: // "-"
	{
		return y.token(yylval, TMINUS)
	}
yyrule43: // "*"
	{
		return y.token(yylval, TMUL)
	}
yyrule44: // "/"
	{
		return y.token(yylval, TDIV)
	}
yyrule45: // ".."
	{
		return y.token(yylval, TRANGE)
	}
yyrule46: // ";"
	{
		return y.token(yylval, TSEMICOLON)
	}
yyrule47: // ","
	{
		return y.token(yylval, TCOMMA)
	}
yyrule48: // "("
	{
		return y.token(yylval, TLPAREN)
	}
yyrule49: // ")"
	{
		return y.token(yylval, TRPAREN)
	}
yyrule50: // "="
	{
		return y.token(yylval, TASSIGN)
	}
yyrule51: // [a-zA-Z_~][a-zA-Z0-9_]*
	{
		return y.strtoken(yylval, TIDENT)
	}
yyrule52: // [ \t\n\r]+

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
