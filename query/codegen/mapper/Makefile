CFLAGS=`llvm-config-3.4 --cflags`
LDFLAGS="`llvm-config-3.4 --ldflags` -Wl,-L`llvm-config-3.4 --libdir` -lLLVM-`llvm-config-3.4 --version`"

test:
	CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go test -v
