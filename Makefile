# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

ifneq ($(strip $(shell which gotmake)),)
all: Makefile binaries web

Makefile: scripts/make.header scripts/mkmake $(wildcard */*/.go) $(wildcard */*.go)
	./scripts/mkmake
else
all: binaries web
endif

test: all
	./scripts/harness

install: installbins installpkgs

web: doc/index.html doc/manual.html \
	$(subst src,doc,$(subst .go,.html,$(wildcard src/*.go))) \
	doc/hydra.svg doc/iolaus.css

doc/index.html: README.md scripts/mkdown scripts/header.html scripts/footer.html
	./scripts/mkdown -o doc/index.html README.md

doc/manual.html: scripts/mkmanual scripts/header.html scripts/footer.html
	./scripts/mkmanual src/*.go

doc/%.svg: scripts/%.svg
	cp -f $< $@

doc/%.css: scripts/%.css
	cp -f $< $@

EXTRAPHONY=man installman

man: $(subst src,doc/man/man1,$(subst .go,.1,$(wildcard src/*.go)))
installman: $(subst src,doc/man/man1,$(subst .go,.1,$(wildcard src/*.go)))
	echo cp -f $? /usr/share/man/man1/

doc/man/man1/%.1: bin/%
	@mkdir -p `dirname $@`
	$< --create-manpage > $@

doc/%.html: doc/man/man1/%.1
	cat scripts/header.html | sed -e "s/Iolaus/$*/" > $@
	groff -man -Thtml $< | tail -n +19 >> $@


include $(GOROOT)/src/Make.$(GOARCH)

binaries:  scripts/harness scripts/mkdown scripts/mkmanual scripts/pdiff bin/iolaus-initialize bin/iolaus-pull bin/iolaus-push bin/iolaus-record bin/iolaus-whatsnew
packages: 

ifndef GOBIN
GOBIN=$(HOME)/bin
endif

# ugly hack to deal with whitespaces in $GOBIN
nullstring :=
space := $(nullstring) # a space at the end
bindir=$(subst $(space),\ ,$(GOBIN))
pkgdir=$(subst $(space),\ ,$(GOROOT)/pkg/$(GOOS)_$(GOARCH))

.PHONY: test binaries packages install installbins installpkgs $(EXTRAPHONY)
.SUFFIXES: .$(O) .go .got .gotgo $(EXTRASUFFIXES)

.go.$(O):
	cd `dirname "$<"`; $(GC) `basename "$<"`
.got.gotgo:
	gotgo "$<"

scripts/gotgo/slice(string).$(O): scripts/gotgo/slice(string).go

scripts/harness: scripts/harness.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/harness.$(O): scripts/harness.go src/util/error.$(O) src/util/exit.$(O)

ifneq ($(strip $(shell which gotgo)),)
# looks like we require scripts/gotgo/slice.got as installed package...
scripts/gotgo/slice(string).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p scripts/gotgo/
	$< 'string' > "$@"
endif
scripts/mkdown: scripts/mkdown.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/mkdown.$(O): scripts/mkdown.go scripts/gotgo/slice(string).$(O) src/util/debug.$(O) src/util/error.$(O)

scripts/mkmanual: scripts/mkmanual.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/mkmanual.$(O): scripts/mkmanual.go scripts/gotgo/slice(string).$(O) src/util/debug.$(O) src/util/error.$(O)

scripts/pdiff: scripts/pdiff.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/pdiff.$(O): scripts/pdiff.go src/util/patience.$(O)

src/git/color.$(O): src/git/color.go src/git/git.$(O)

ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/git/gotgo/slice.got as installed package...
src/git/gotgo/slice(string).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/git/gotgo/
	$< 'string' > "$@"
endif
src/git/git.$(O): src/git/git.go src/git/gotgo/slice(string).$(O) src/util/debug.$(O) src/util/exit.$(O)

src/git/gotgo/slice(git.CommitHash).$(O): src/git/gotgo/slice(git.CommitHash).go src/git/git.$(O)

src/git/gotgo/slice(string).$(O): src/git/gotgo/slice(string).go

ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/git/gotgo/slice.got as installed package...
src/git/gotgo/slice(git.CommitHash).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/git/gotgo/
	$< --import 'import git "../git"' 'git.CommitHash' > "$@"
endif
src/git/plumbing.$(O): src/git/plumbing.go src/git/color.$(O) src/git/git.$(O) src/git/gotgo/slice(git.CommitHash).$(O) src/git/gotgo/slice(string).$(O) src/util/debug.$(O) src/util/error.$(O) src/util/patience.$(O)

src/git/porcelain.$(O): src/git/porcelain.go src/git/git.$(O)

src/gotgo/slice(git.Commitish).$(O): src/gotgo/slice(git.Commitish).go src/git/git.$(O)

src/iolaus/gotgo/box(git.CommitHash,git.Commitish).$(O): src/iolaus/gotgo/box(git.CommitHash,git.Commitish).go src/git/git.$(O)

ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/iolaus/gotgo/box.got as installed package...
src/iolaus/gotgo/box(git.CommitHash,git.Commitish).go: $(pkgdir)/./gotgo/box.gotgo
	mkdir -p src/iolaus/gotgo/
	$< --import 'import git "../../git/git"' 'git.CommitHash' 'git.Commitish' > "$@"
endif
src/iolaus/test.$(O): src/iolaus/test.go src/git/git.$(O) src/git/plumbing.$(O) src/iolaus/gotgo/box(git.CommitHash,git.Commitish).$(O) src/util/out.$(O)

bin/iolaus-initialize: src/iolaus-initialize.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-initialize: bin/iolaus-initialize
	cp $< $@
src/iolaus-initialize.$(O): src/iolaus-initialize.go src/git/git.$(O) src/git/porcelain.$(O) src/util/error.$(O) src/util/help.$(O)

ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/gotgo/slice.got as installed package...
src/gotgo/slice(git.Commitish).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/gotgo/
	$< --import 'import git "../git/git"' 'git.Commitish' > "$@"
endif
bin/iolaus-pull: src/iolaus-pull.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-pull: bin/iolaus-pull
	cp $< $@
src/iolaus-pull.$(O): src/iolaus-pull.go src/git/git.$(O) src/git/plumbing.$(O) src/gotgo/slice(git.Commitish).$(O) src/util/error.$(O) src/util/exit.$(O) src/util/help.$(O) src/util/out.$(O)

bin/iolaus-push: src/iolaus-push.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-push: bin/iolaus-push
	cp $< $@
src/iolaus-push.$(O): src/iolaus-push.go src/git/git.$(O) src/git/plumbing.$(O) src/gotgo/slice(git.Commitish).$(O) src/util/error.$(O) src/util/exit.$(O) src/util/help.$(O) src/util/out.$(O)

bin/iolaus-record: src/iolaus-record.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-record: bin/iolaus-record
	cp $< $@
src/iolaus-record.$(O): src/iolaus-record.go src/git/git.$(O) src/git/plumbing.$(O) src/gotgo/slice(git.Commitish).$(O) src/iolaus/test.$(O) src/util/cook.$(O) src/util/error.$(O) src/util/help.$(O) src/util/out.$(O)

bin/iolaus-whatsnew: src/iolaus-whatsnew.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/iolaus-whatsnew: bin/iolaus-whatsnew
	cp $< $@
src/iolaus-whatsnew.$(O): src/iolaus-whatsnew.go src/git/color.$(O) src/git/git.$(O) src/git/plumbing.$(O) src/util/exit.$(O) src/util/help.$(O) src/util/out.$(O)

src/util/cook.$(O): src/util/cook.go src/util/exit.$(O)

src/util/debug.$(O): src/util/debug.go

src/util/error.$(O): src/util/error.go src/util/cook.$(O) src/util/exit.$(O)

ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/util/gotgo/slice.got as installed package...
src/util/gotgo/slice(func()).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/util/gotgo/
	$< 'func()' > "$@"
endif
src/util/exit.$(O): src/util/exit.go src/util/gotgo/slice(func()).$(O)

src/util/gotgo/slice([]pt.PatienceElem).$(O): src/util/gotgo/slice([]pt.PatienceElem).go src/util/patienceTypes.$(O)

src/util/gotgo/slice(func()).$(O): src/util/gotgo/slice(func()).go

src/util/gotgo/slice(int).$(O): src/util/gotgo/slice(int).go

src/util/gotgo/slice(pt.PatienceElem).$(O): src/util/gotgo/slice(pt.PatienceElem).go src/util/patienceTypes.$(O)

src/util/gotgo/slice(pt.StringChunk).$(O): src/util/gotgo/slice(pt.StringChunk).go src/util/patienceTypes.$(O)

src/util/help.$(O): src/util/help.go

src/util/out.$(O): src/util/out.go src/util/cook.$(O)

ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/util/gotgo/slice.got as installed package...
src/util/gotgo/slice([]pt.PatienceElem).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/util/gotgo/
	$< --import 'import pt "../patienceTypes"' '[]pt.PatienceElem' > "$@"
endif
ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/util/gotgo/slice.got as installed package...
src/util/gotgo/slice(int).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/util/gotgo/
	$< 'int' > "$@"
endif
ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/util/gotgo/slice.got as installed package...
src/util/gotgo/slice(pt.StringChunk).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/util/gotgo/
	$< --import 'import pt "../patienceTypes"' 'pt.StringChunk' > "$@"
endif
ifneq ($(strip $(shell which gotgo)),)
# looks like we require src/util/gotgo/slice.got as installed package...
src/util/gotgo/slice(pt.PatienceElem).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p src/util/gotgo/
	$< --import 'import pt "../patienceTypes"' 'pt.PatienceElem' > "$@"
endif
src/util/patience.$(O): src/util/patience.go src/util/gotgo/slice([]pt.PatienceElem).$(O) src/util/gotgo/slice(int).$(O) src/util/gotgo/slice(pt.PatienceElem).$(O) src/util/gotgo/slice(pt.StringChunk).$(O) src/util/patienceTypes.$(O)

src/util/patienceTypes.$(O): src/util/patienceTypes.go src/git/color.$(O)

installbins:  $(bindir)/iolaus-initialize $(bindir)/iolaus-pull $(bindir)/iolaus-push $(bindir)/iolaus-record $(bindir)/iolaus-whatsnew
installpkgs: 
