%{
package query

import (
    "bufio"
    "fmt"
    "strconv"
)

type yylexer struct{
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
        src: src,
        startToken:startToken,
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
%}

%yyc c
%yyn c = y.getc()

%%
    y.buf = y.buf[:0]

\"(\\.|[^\\"])*\"          return y.quotedstrtoken(yylval, TSTRING);
\'(\\.|[^\\'])*\'          return y.quotedstrtoken(yylval, TSTRING);
[0-9]+                     return y.inttoken(yylval, TINT);

"SELECT"                   return y.token(yylval, TSELECT);
"GROUP"                    return y.token(yylval, TGROUP);
"BY"                       return y.token(yylval, TBY);
"INTO"                     return y.token(yylval, TINTO);
"AS"                       return y.token(yylval, TAS);
"WHEN"                     return y.token(yylval, TWHEN);
"WITHIN"                   return y.token(yylval, TWITHIN);
"THEN"                     return y.token(yylval, TTHEN);
"END"                      return y.token(yylval, TEND);
(STEPS|SESSIONS)           return y.strtoken(yylval, TWITHINUNITS);
"true"                     return y.token(yylval, TTRUE);
"false"                    return y.token(yylval, TFALSE);
"=="                       return y.token(yylval, TEQUALS);
"!="                       return y.token(yylval, TNOTEQUALS);
"<="                       return y.token(yylval,  TLTE);
"<"                        return y.token(yylval,  TLT);
">="                       return y.token(yylval,  TGTE);
">"                        return y.token(yylval,  TGT);
"&&"                       return y.token(yylval, TAND);
"||"                       return y.token(yylval, TOR);
"+"                        return y.token(yylval, TPLUS);
"-"                        return y.token(yylval, TMINUS);
"*"                        return y.token(yylval, TMUL);
"/"                        return y.token(yylval, TDIV);
".."                       return y.token(yylval, TRANGE);
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

// Saves the quoted string in the buffer and the token to the parser value
// and returns the token.
func (y *yylexer) quotedstrtoken(yylval *yySymType, tok int) int {
    str := string(y.buf)
    yylval.str = str[1:len(str)-1]
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
