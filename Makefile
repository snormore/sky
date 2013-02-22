################################################################################
# Variables
################################################################################

VERSION=0.2.3

CFLAGS=-g -Wall -Wextra -Wno-strict-overflow -std=gnu99 -D_GNU_SOURCE -D_FILE_OFFSET_BITS=64 -Ideps/leveldb-1.7.0/include -Ideps/LuaJIT-2.0.0/src -I/usr/local/include
CXXFLAGS=-g -Wall -Wextra -Wno-strict-overflow -std=gnu99 -D_GNU_SOURCE -D_FILE_OFFSET_BITS=64 -Ideps/leveldb-1.7.0/include -Ideps/LuaJIT-2.0.0/src -I/usr/local/include
LIBS=-lpthread -lzmq -ldl

SOURCES=$(wildcard src/**/*.c src/**/**/*.c src/*.c)
OBJECTS=$(patsubst %.c,%.o,${SOURCES})
BIN_SOURCES=src/skyd.c
BIN_OBJECTS=$(patsubst %.c,%.o,${BIN_SOURCES})
LIB_SOURCES=$(filter-out ${BIN_SOURCES},${SOURCES})
LIB_OBJECTS=$(filter-out ${BIN_OBJECTS},${OBJECTS})
TEST_SOURCES=$(wildcard tests/*_tests.c tests/**/*_tests.c)
TEST_OBJECTS=$(patsubst %.c,%,${TEST_SOURCES})

PACKAGE=pkg/sky-${VERSION}.tar.gz
PKGTMPDIR=pkg/tmp/sky-${VERSION}

UNAME=$(shell uname)
ifeq ($(UNAME), Darwin)
LUAJIT_FLAGS=-pagezero_size 10000 -image_base 100000000
CFLAGS+=-Wno-self-assign
endif
ifeq ($(UNAME), Linux)
CFLAGS+=-rdynamic -L/usr/local/lib
CXXFLAGS+=-rdynamic -L/usr/local/lib
endif

PREFIX?=/usr/local

.PHONY: test valgrind

################################################################################
# Main Targets
################################################################################

compile: bin/libleveldb.a bin/libluajit.a bin/libsky.a bin/skyd
all: compile test

################################################################################
# Dependencies
################################################################################

bin/libleveldb.a: bin
	${MAKE} -C deps/leveldb-1.7.0
	mv deps/leveldb-1.7.0/libleveldb.a bin/libleveldb.a

bin/libluajit.a: bin
	${MAKE} -C deps/LuaJIT-2.0.0
	mv deps/LuaJIT-2.0.0/src/libluajit.a bin/libluajit.a

################################################################################
# Installation
################################################################################

install:
	install -d $(DESTDIR)/$(PREFIX)/sky/
	install -d $(DESTDIR)/$(PREFIX)/sky/bin/
	install -d $(DESTDIR)/$(PREFIX)/sky/data/
	install bin/skyd $(DESTDIR)/$(PREFIX)/sky/bin/
	rm -f $(DESTDIR)/$(PREFIX)/bin/skyd
	ln -s $(DESTDIR)/$(PREFIX)/sky/bin/skyd $(DESTDIR)/$(PREFIX)/bin/skyd

################################################################################
# Package
################################################################################

package: cleaner
	rm -rf pkg
	mkdir -p ${PKGTMPDIR}
	cp Makefile LICENSE README.md ${PKGTMPDIR}
	cp -r deps ${PKGTMPDIR}/deps
	cp -r src ${PKGTMPDIR}/src
	cp -r tests ${PKGTMPDIR}/tests
	tar czvf ${PACKAGE} -C pkg/tmp .
	rm -rf pkg/tmp

################################################################################
# Binaries
################################################################################

bin/libsky.a: bin ${LIB_OBJECTS}
	rm -f bin/libsky.a
	ar rcs $@ ${LIB_OBJECTS}
	ranlib $@

bin/skyd: bin ${OBJECTS} bin/libsky.a bin/libleveldb.a bin/libluajit.a
	$(CC) $(CFLAGS) -Isrc -c -o $@.o src/skyd.c
	$(CXX) $(CXXFLAGS) -Isrc $(LUAJIT_FLAGS) -o $@ $@.o bin/libsky.a bin/libleveldb.a bin/libluajit.a $(LIBS)
	chmod 700 $@

bin:
	mkdir -p bin


################################################################################
# jsmn
################################################################################

src/jsmn/jsmn.o: src/jsmn/jsmn.c
	$(CC) $(CFLAGS) -Wno-sign-compare -Isrc -c -o $@ $<



################################################################################
# Memory Leak Check
################################################################################

valgrind: VALGRIND=valgrind
valgrind: CFLAGS+=-DVALGRIND
valgrind: clean all

################################################################################
# Release Build (No Debug)
################################################################################

release: CFLAGS+=-DNDEBUG
release: clean all


################################################################################
# Tests
################################################################################

test: $(TEST_OBJECTS) tmp
	@sh ./tests/runtests.sh $(VALGRIND)

$(TEST_OBJECTS): %: %.c bin/libsky.a bin/libleveldb.a bin/libluajit.a
	$(CC) $(CFLAGS) -Isrc -c -o $@.o $<
	$(CXX) $(CXXFLAGS) -Isrc $(LUAJIT_FLAGS) -o $@ $@.o bin/libsky.a bin/libleveldb.a bin/libluajit.a $(LIBS)


################################################################################
# Misc
################################################################################

info:
	echo $(CFLAGS)

tmp:
	mkdir -p tmp


################################################################################
# Clean
################################################################################

clean: 
	rm -rf bin ${OBJECTS} ${TEST_OBJECTS}
	rm -rf tests/*.dSYM tests/**/*.dSYM
	rm -rf  tests/*.o tests/**/*.o
	rm -rf tmp pkg

clean-leveldb:
	rm -f bin/libleveldb.a
	${MAKE} clean -C deps/leveldb-1.7.0

clean-luajit:
	rm -f bin/libluajit.a
	${MAKE} clean -C deps/LuaJIT-2.0.0

cleaner: clean clean-leveldb clean-luajit
