%{
package query

import (
    "bufio"
    "fmt"
    "log"
)

type yylexer struct{
    src     *bufio.Reader
    buf     []byte
    empty   bool
    current byte
    index   int
    query   *Query
}

func newLexer(src *bufio.Reader) *yylexer {
    y := &yylexer{src: src, query:&Query{}}
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
%}

%yyc c
%yyn c = y.getc()

%%
    y.buf = y.buf[:0]

"SELECT"                   return y.token(yylval, TSELECT);
"GROUP"                    return y.token(yylval, TGROUP);
"BY"                       return y.token(yylval, TBY);
[a-zA-Z_~][a-zA-Z0-9_]*    return y.strtoken(yylval, TIDENT);
";"                        return y.token(yylval, TSEMICOLON);
","                        return y.token(yylval, TCOMMA);
"("                        return y.token(yylval, TLPAREN);
")"                        return y.token(yylval, TRPAREN);
[ \t\n\r]+

%%
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
