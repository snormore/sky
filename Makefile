PREFIX=/usr/local
DESTDIR=
BINDIR=${PREFIX}/bin
DATADIR=${PREFIX}/share

SKYD_SRCS = $(wildcard skyd/*.go)

BINARIES = skyd
BUILDDIR = build

all: $(BINARIES)

$(BUILDDIR)/%:
	mkdir -p $(dir $@)
	cd $* && go build -o $(abspath $@)

$(BINARIES): %: $(BUILDDIR)/%

# Dependencies
$(BUILDDIR)/skyd: $(SKYD_SRCS)

clean:
	rm -fr $(BUILDDIR)

# Targets
.PHONY: install clean all
# Programs
.PHONY: $(BINARIES)

install: $(BINARIES)
	install -m 755 -d ${DESTDIR}${BINDIR}
	install -m 755 $(BUILDDIR)/skyd ${DESTDIR}${BINDIR}/skyd
	install -m 755 -d ${DESTDIR}${DATADIR}
