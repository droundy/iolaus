# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.


all: Makefile  bin/iolaus-initialize bin/iolaus-record bin/iolaus-whatsnew bin/mkmake bin/pdiff

include $(GOROOT)/src/Make.$(GOARCH)

.PHONY: test
.SUFFIXES: .$(O) .go .got .gotgo

Makefile: bin/mkmake
	./bin/mkmake > $@
.go.$(O):
	cd `dirname "$<"`; $(GC) `basename "$<"`
.got.gotgo:
	gotgo "$<"

src/git/git.$(O): src/git/git.go src/util/debug.$(O) src/util/exit.$(O)

src/git/plumbing.$(O): src/git/plumbing.go src/util/patience.$(O) src/util/debug.$(O) src/util/error.$(O) src/git/git.$(O)

src/git/porcelain.$(O): src/git/porcelain.go src/git/git.$(O)

bin/iolaus-initialize: src/iolaus-initialize.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
src/iolaus-initialize.$(O): src/iolaus-initialize.go src/util/error.$(O) src/util/help.$(O) src/git/porcelain.$(O) src/git/git.$(O)

bin/iolaus-record: src/iolaus-record.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
src/iolaus-record.$(O): src/iolaus-record.go src/util/error.$(O) src/util/out.$(O) src/util/help.$(O) src/git/git.$(O) src/util/cook.$(O) src/git/plumbing.$(O)

bin/iolaus-whatsnew: src/iolaus-whatsnew.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
src/iolaus-whatsnew.$(O): src/iolaus-whatsnew.go src/util/out.$(O) src/util/help.$(O) src/git/git.$(O) src/git/plumbing.$(O)

bin/mkmake: src/mkmake.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
src/mkmake.$(O): src/mkmake.go

bin/pdiff: src/pdiff.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
src/pdiff.$(O): src/pdiff.go src/util/patience.$(O)

src/util/cook.$(O): src/util/cook.go src/util/exit.$(O)

src/util/debug.$(O): src/util/debug.go

src/util/error.$(O): src/util/error.go src/util/exit.$(O) src/util/cook.$(O)

src/util/exit.$(O): src/util/exit.go

src/util/help.$(O): src/util/help.go

src/util/out.$(O): src/util/out.go src/util/cook.$(O)

src/util/patience.$(O): src/util/patience.go

