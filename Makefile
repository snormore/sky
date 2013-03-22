################################################################################
# Variables
################################################################################

VERSION=0.3.0
PREFIX=/usr/local
DESTDIR=
BINDIR=${PREFIX}/bin
DATADIR=${PREFIX}/share
INCLUDEDIR=${PREFIX}/include
LIBDIR=${PREFIX}/lib

SKYD_SRCS = $(wildcard skyd/*.go)

BINARIES = skyd
BUILDDIR = build

all: $(BINARIES)

$(BUILDDIR)/%:
	mkdir -p $(dir $@)
	cd $* && go build -o $(abspath $@)

$(BINARIES): %: $(BUILDDIR)/%

UNAME=$(shell uname)
ifeq ($(UNAME), Darwin)
SOSUFFIX=dylib
endif
ifeq ($(UNAME), Linux)
SOSUFFIX=so
endif

################################################################################
# Dependencies
################################################################################

leveldb:
	${MAKE} clean -C deps/leveldb-1.9.0
	${MAKE} -C deps/leveldb-1.9.0
	install -m 755 -d ${DESTDIR}${INCLUDEDIR}/leveldb
	install -m 755 deps/leveldb-1.9.0/include/leveldb/* ${DESTDIR}${INCLUDEDIR}/leveldb
	install -m 755 -d ${DESTDIR}${LIBDIR}
	install -m 755 deps/leveldb-1.9.0/libleveldb.a ${DESTDIR}${LIBDIR}
	install -m 755 deps/leveldb-1.9.0/libleveldb.${SOSUFFIX}.1.9 ${DESTDIR}${LIBDIR}
	ln -sf ${DESTDIR}${LIBDIR}/libleveldb.${SOSUFFIX}.1.9 ${DESTDIR}${LIBDIR}/libleveldb.${SOSUFFIX}.1
	ln -sf ${DESTDIR}${LIBDIR}/libleveldb.${SOSUFFIX}.1.9 ${DESTDIR}${LIBDIR}/libleveldb.${SOSUFFIX}

luajit:
	${MAKE} -C deps/LuaJIT-2.0.1 clean PREFIX=${PREFIX}
	${MAKE} -C deps/LuaJIT-2.0.1 PREFIX=${PREFIX}
	${MAKE} -C deps/LuaJIT-2.0.1 install PREFIX=${PREFIX}

csky:
	${MAKE} -C deps/csky clean install PREFIX=${PREFIX}

data:
	mkdir -p /var/lib/sky
	chmod 777 /var/lib/sky

deps: leveldb luajit csky data


################################################################################
# Build
################################################################################

$(BUILDDIR)/skyd: $(SKYD_SRCS)


################################################################################
# Clean
################################################################################

clean:
	rm -fr $(BUILDDIR)

.PHONY: install clean all csky leveldb luajit data
.PHONY: $(BINARIES)

install: $(BINARIES)
	install -m 755 -d ${DESTDIR}${BINDIR}
	install -m 755 $(BUILDDIR)/skyd ${DESTDIR}${BINDIR}/skyd
	install -m 755 -d ${DESTDIR}${DATADIR}
