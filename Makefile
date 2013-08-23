PACKAGES=core factors query server skyd/config
PKGPATHS=$(patsubst %,github.com/skydb/sky/%,$(PACKAGES))

all: fmt test

test:
	${MAKE} -C core test
	${MAKE} -C factors test
	${MAKE} -C query test
	${MAKE} -C server test

fmt:
	go fmt $(PKGPATHS)

