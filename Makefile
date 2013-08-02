PACKAGES=core factors query query/engine server skyd/config
PKGPATHS=$(patsubst %,github.com/snormore/sky/%,$(PACKAGES))

test:
	go test -v $(PKGPATHS)

fmt:
	go fmt $(PKGPATHS)

