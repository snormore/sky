all: lexer.go parser.go

test: all
	go test -v

lexer.go: lexer.l
	golex -t $< | gofmt > $@

parser.go: parser.y
	go tool yacc -o $@ $<
	go fmt
