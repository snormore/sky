all: lexer.go grammar.go

test: all
	go test -v

lexer.go: lexer.l
	golex -t $< > $@
	go fmt

grammar.go: grammar.y
	go tool yacc -o $@ $<
	go fmt
