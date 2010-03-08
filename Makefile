# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

all: Makefile binaries

Makefile: scripts/make.header $(wildcard *.go)
	cp -f scripts/make.header $@
	gotmake >> $@

test: tests/example

install: installbins installpkgs


include $(GOROOT)/src/Make.$(GOARCH)

binaries:  bin/iolaus-initialize bin/iolaus-record bin/iolaus-whatsnew bin/pdiff
packages: 

ifndef GOBIN
GOBIN=$(HOME)/bin
endif

# ugly hack to deal with whitespaces in $GOBIN
nullstring :=
space := $(nullstring) # a space at the end
bindir=$(subst $(space),\ ,$(GOBIN))
pkgdir=$(subst $(space),\ ,$(GOROOT)/pkg/$(GOOS)_$(GOARCH))

.PHONY: test binaries packages install installbins installpkgs
.SUFFIXES: .$(O) .go .got .gotgo

.go.$(O):
	cd `dirname "$<"`; $(GC) `basename "$<"`
.got.gotgo:
	gotgo "$<"

# looks like we require src/git/gotgo/slice.got as installed package...
src/git/gotgo/slice(string).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/git/gotgo/
	$< 'string' > "$@"
src/git/git.$(O): src/git/git.go src/git/gotgo/slice(string).$(O) src/util/debug.$(O) src/util/exit.$(O)

src/git/gotgo/slice(string).$(O): src/git/gotgo/slice(string).go

src/git/plumbing.$(O): src/git/plumbing.go src/git/git.$(O) src/git/gotgo/slice(string).$(O) src/util/debug.$(O) src/util/error.$(O) src/util/patience.$(O)

src/git/porcelain.$(O): src/git/porcelain.go src/git/git.$(O)

bin/iolaus-initialize: src/iolaus-initialize.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-initialize: bin/iolaus-initialize
	cp $< $@
src/iolaus-initialize.$(O): src/iolaus-initialize.go src/git/git.$(O) src/git/porcelain.$(O) src/util/error.$(O) src/util/help.$(O)

bin/iolaus-record: src/iolaus-record.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-record: bin/iolaus-record
	cp $< $@
src/iolaus-record.$(O): src/iolaus-record.go src/git/git.$(O) src/git/plumbing.$(O) src/util/cook.$(O) src/util/error.$(O) src/util/help.$(O) src/util/out.$(O)

bin/iolaus-whatsnew: src/iolaus-whatsnew.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-whatsnew: bin/iolaus-whatsnew
	cp $< $@
src/iolaus-whatsnew.$(O): src/iolaus-whatsnew.go src/git/git.$(O) src/git/plumbing.$(O) src/util/help.$(O) src/util/out.$(O)

bin/pdiff: src/pdiff.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/pdiff: bin/pdiff
	cp $< $@
src/pdiff.$(O): src/pdiff.go src/util/patience.$(O)

src/util/cook.$(O): src/util/cook.go src/util/exit.$(O)

src/util/debug.$(O): src/util/debug.go

src/util/error.$(O): src/util/error.go src/util/cook.$(O) src/util/exit.$(O)

# looks like we require src/util/gotgo/slice.got as installed package...
src/util/gotgo/slice(func()).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/util/gotgo/
	$< 'func()' > "$@"
src/util/exit.$(O): src/util/exit.go src/util/gotgo/slice(func()).$(O)

src/util/gotgo/slice(func()).$(O): src/util/gotgo/slice(func()).go

src/util/gotgo/slice(int).$(O): src/util/gotgo/slice(int).go

src/util/help.$(O): src/util/help.go

src/util/out.$(O): src/util/out.go src/util/cook.$(O)

# looks like we require src/util/gotgo/slice.got as installed package...
src/util/gotgo/slice(int).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/util/gotgo/
	$< 'int' > "$@"
src/util/patience.$(O): src/util/patience.go src/util/gotgo/slice(int).$(O)

installbins:  $(bindir)/iolaus-initialize $(bindir)/iolaus-record $(bindir)/iolaus-whatsnew $(bindir)/pdiff
installpkgs: 
