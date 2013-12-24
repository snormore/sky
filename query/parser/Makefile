all: lexer.go grammar.go

lexer.go: lexer.l
	golex -t $< > $@
	go fmt $@

grammar.go: grammar.y
	go tool yacc -o $@ $<
	go fmt $@
