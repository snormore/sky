CFLAGS=`llvm-config-3.4 --cflags`
LDFLAGS="`llvm-config-3.4 --ldflags` -Wl,-L`llvm-config-3.4 --libdir` -lLLVM-`llvm-config-3.4 --version`"
TEST=.
PKG=./...

all: test

flags:
	@echo "CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS)"

grammar:
	${MAKE} -C query/parser

test: grammar
	CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go test -v -test.run=$(TEST) $(PKG)

bench: grammar
	CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go test -v -test.bench=. $(PKG)

.PHONY: test
