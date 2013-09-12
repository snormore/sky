%{
package query

import (
    "bufio"
    "fmt"
    "strconv"
    "strings"
)

type yylexer struct{
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
%}

%yyc c
%yyn c = y.getc()

%%
    y.tlineidx, y.tcharidx = y.lineidx, y.charidx
    y.buf = y.buf[:0]

\"(\\.|[^\\"])*\"          return y.quotedstrtoken(yylval, TQUOTEDSTRING, "\"");
\'(\\.|[^\\'])*\'          return y.quotedstrtoken(yylval, TQUOTEDSTRING, "'");
[0-9]+                     return y.inttoken(yylval, TINT);

"DECLARE"                  return y.token(yylval, TDECLARE);
"SET"                      return y.token(yylval, TSET);
"EXIT"                     return y.token(yylval, TEXIT);
"DEBUG"                    return y.token(yylval, TDEBUG);
"SELECT"                   return y.token(yylval, TSELECT);
"GROUP"                    return y.token(yylval, TGROUP);
"BY"                       return y.token(yylval, TBY);
"INTO"                     return y.token(yylval, TINTO);
"AS"                       return y.token(yylval, TAS);
"WHEN"                     return y.token(yylval, TWHEN);
"WITHIN"                   return y.token(yylval, TWITHIN);
"THEN"                     return y.token(yylval, TTHEN);
"END"                      return y.token(yylval, TEND);
"FOR"                      return y.token(yylval, TFOR);
"EACH"                     return y.token(yylval, TEACH);
"EVERY"                    return y.token(yylval, TEVERY);
"IN"                       return y.token(yylval, TIN);
"EVENT"                    return y.token(yylval, TEVENT);
"SESSION"                  return y.token(yylval, TSESSION);
"DELIMITED"                return y.token(yylval, TDELIMITED);
"FACTOR"                   return y.token(yylval, TFACTOR);
"STRING"                   return y.token(yylval, TSTRING);
"INTEGER"                  return y.token(yylval, TINTEGER);
"FLOAT"                    return y.token(yylval, TFLOAT);
"BOOLEAN"                  return y.token(yylval, TBOOLEAN);
(SECOND|SECONDS|MINUTE|MINUTES|HOUR|HOURS|DAY|DAYS|WEEK|WEEKS|YEAR|YEARS)  return y.strtoken(yylval, TTIMEUNITS);
(STEPS)                    return y.strtoken(yylval, TWITHINUNITS);
"true"                     return y.token(yylval, TTRUE);
"false"                    return y.token(yylval, TFALSE);
"=="                       return y.token(yylval, TEQUALS);
"!="                       return y.token(yylval, TNOTEQUALS);
"<="                       return y.token(yylval, TLTE);
"<"                        return y.token(yylval, TLT);
">="                       return y.token(yylval, TGTE);
">"                        return y.token(yylval, TGT);
"&&"                       return y.token(yylval, TAND);
"||"                       return y.token(yylval, TOR);
"+"                        return y.token(yylval, TPLUS);
"-"                        return y.token(yylval, TMINUS);
"*"                        return y.token(yylval, TMUL);
"/"                        return y.token(yylval, TDIV);
".."                       return y.token(yylval, TRANGE);
";"                        return y.token(yylval, TSEMICOLON);
","                        return y.token(yylval, TCOMMA);
"("                        return y.token(yylval, TLPAREN);
")"                        return y.token(yylval, TRPAREN);
"="                        return y.token(yylval, TASSIGN);

\@\@?[a-zA-Z_~][a-zA-Z0-9_]*  return y.strtoken(yylval, TAMPIDENT);
[a-zA-Z_~][a-zA-Z0-9_]*    return y.strtoken(yylval, TIDENT);
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
func (y *yylexer) quotedstrtoken(yylval *yySymType, tok int, quote string) int {
    str := string(y.buf)
    str = str[1:len(str)-1]
    str = strings.Replace(str, "\\" + quote, quote, -1)
    yylval.str = str
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
