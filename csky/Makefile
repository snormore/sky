################################################################################
# Variables
################################################################################

MAJOR_VERSION=0
MINOR_VERSION=3
PATCH_VERSION=0

PREFIX?=/usr/local
CFLAGS=-g -Wall -Wextra -std=c99 -Iinclude
LDFLAGS=-lleveldb
SOURCES=$(wildcard src/*.c)
OBJECTS=$(patsubst %.c,%.o,${SOURCES})
TEST_SOURCES=$(wildcard tests/*_tests.c)
TEST_OBJECTS=$(patsubst %.c,%,${TEST_SOURCES})

UNAME=$(shell uname)
ifeq ($(UNAME), Darwin)
LDFLAGS+=-pagezero_size 10000 -image_base 100000000
endif
ifeq ($(UNAME), Linux)
LDFLAGS+=-rdynamic
endif


################################################################################
# Targets
################################################################################

default: build

build: libcsky.a

libcsky.a: ${OBJECTS}
	rm -f bin/libsky.a
	ar rcs libcsky.a ${OBJECTS}
	ranlib libcsky.a

install:
	install -d $(DESTDIR)/$(PREFIX)/include/sky
	install include/sky/* $(DESTDIR)/$(PREFIX)/include/sky
	install -d $(DESTDIR)/$(PREFIX)/lib
	install libcsky* $(DESTDIR)/$(PREFIX)/lib

clean: 
	rm -rf bin ${OBJECTS} ${TEST_OBJECTS}
	rm -rf tests/*.dSYM tests/**/*.dSYM
	rm -rf tests/*.o tests/**/*.o


################################################################################
# Tests
################################################################################

check: $(TEST_OBJECTS)
	@sh ./tests/runtests.sh

$(TEST_OBJECTS): %: %.c libcsky.a
	$(CC) $(CFLAGS) -Itests -o $@ $< libcsky.a
