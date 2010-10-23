# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

all: binaries

test: all
	./scripts/harness
clean:
	rm -f */*.$(O) */*/*.$(O) */*/*/*.$(O) bin/*

EXE=
ifeq ($(GOOS),windows)
	EXE=.exe
endif

install: installbins installman

web: doc/index.html doc/manual.html doc/install.html doc/TODO.html \
	$(subst src,doc,$(subst .go,.html,$(wildcard src/iolaus-*.go))) \
	doc/hydra.svg doc/iolaus.css

doc/index.html: README.md scripts/mkdown scripts/header.html scripts/footer.html
	./scripts/mkdown -o doc/index.html README.md

doc/manual.html: scripts/mkmanual scripts/header.html scripts/footer.html
	./scripts/mkmanual src/iolaus-*.go

doc/install.html: INSTALL.md scripts/mkdown scripts/header.html scripts/footer.html
	./scripts/mkdown -o $@ $<

doc/TODO.html: TODO.md scripts/mkdown scripts/header.html scripts/footer.html
	./scripts/mkdown -o $@ $<

doc/%.svg: scripts/%.svg
	cp -f $< $@

doc/%.css: scripts/%.css
	cp -f $< $@

man: $(subst src,doc/man/man1,$(subst .go,.1,$(wildcard src/iolaus-*.go)))
installman: $(subst src,doc/man/man1,$(subst .go,.1,$(wildcard src/iolaus-*.go)))
	echo cp -f $? /usr/share/man/man1/

doc/man/man1/iolaus-%.1: bin/iolaus-%
	@mkdir -p `dirname $@`
	$< --create-manpage > $@

doc/%.html: doc/man/man1/%.1
	cat scripts/header.html | sed -e "s/Iolaus/$*/" > $@
	groff -man -Thtml $< | tail -n +19 >> $@


include $(GOROOT)/src/Make.inc

binaries:  scripts/harness scripts/mkdown scripts/mkmanual scripts/pdiff \
	bin/iolaus-initialize$(EXE) \
	bin/iolaus-pull$(EXE) bin/iolaus-push$(EXE) \
	bin/iolaus-record$(EXE) bin/iolaus-stash$(EXE) \
	bin/iolaus-whatsnew$(EXE) bin/iolaus-changes$(EXE)

ifndef GOBIN
GOBIN=$(HOME)/bin
endif

# ugly hack to deal with whitespaces in $GOBIN
nullstring :=
space := $(nullstring) # a space at the end
bindir=$(subst $(space),\ ,$(GOBIN))
pkgdir=$(subst $(space),\ ,$(GOROOT)/pkg/$(GOOS)_$(GOARCH))
srcpkgdir=$(subst $(space),\ ,$(GOROOT)/src/pkg)

gooptpkg=$(pkgdir)/github.com/droundy/goopt.a

.PHONY: test binaries install installbins man installman clean
.SUFFIXES: .$(O) .go .got

.go.$(O):
	cd `dirname "$<"`; $(GC) `basename "$<"`

bin/iolaus-%$(EXE): src/iolaus-%.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
$(bindir)/%: bin/%
	cp $< $@

scripts/harness: scripts/harness.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/harness.$(O): src/util/error.$(O) src/util/exit.$(O)

scripts/mkdown: scripts/mkdown.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/mkdown.$(O): src/util/slice(string).$(O) \
	src/util/debug.$(O) src/util/error.$(O)

scripts/mkmanual: scripts/mkmanual.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/mkmanual.$(O): src/util/slice(string).$(O) \
	src/util/debug.$(O) src/util/error.$(O)

scripts/pdiff: scripts/pdiff.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
scripts/pdiff.$(O): src/util/patience.$(O)

src/git/color.$(O): src/git/git.$(O)

# Automatically install goopt if it isn't yet installed...
$(gooptpkg):
	@echo Installing goopt as $(gooptpkg)
	goinstall github.com/droundy/goopt
src/util/debug.$(O): $(gooptpkg)

src/util/machine.$(O): src/util/machine_$(GOOS).go
	cd `dirname "$<"`; $(GC) -o `basename "$@"` `basename "$<"`

src/git/git.$(O): src/util/slice(string).$(O) \
	src/util/debug.$(O) src/util/exit.$(O)

src/git/gotgo/slice(git.CommitHash).$(O): src/git/git.$(O)

src/git/plumbing.$(O): src/git/color.$(O) src/git/git.$(O) \
	src/git/gotgo/slice(git.CommitHash).$(O) \
	src/util/slice(string).$(O) src/util/debug.$(O) \
	src/util/error.$(O) src/util/patience.$(O)

src/git/porcelain.$(O): src/git/git.$(O)

src/gotgo/slice(git.Commitish).$(O): src/git/git.$(O)

src/iolaus/core.$(O): src/git/color.$(O) src/git/plumbing.$(O) \
	src/util/out.$(O) src/util/patience.$(O)

src/iolaus/gotgo/box(git.CommitHash,git.Commitish).$(O): src/git/git.$(O)

src/iolaus/prompt.$(O): src/git/color.$(O) src/iolaus/core.$(O) \
	src/util/error.$(O)

src/iolaus/promptcommit.$(O): src/iolaus/core.$(O) \
	src/iolaus/gotgo/box(git.CommitHash,git.Commitish).$(O) \
	src/util/error.$(O)

src/iolaus/test.$(O): src/git/plumbing.$(O) \
	src/iolaus/gotgo/box(git.CommitHash,git.Commitish).$(O) \
	src/util/out.$(O) src/util/machine.$(O) \

src/iolaus-initialize.$(O): src/git/git.$(O) src/git/porcelain.$(O) \
	src/util/error.$(O) src/util/help.$(O)

src/iolaus-changes.$(O): src/iolaus-changes.go \
		src/private-cs.go src/git/plumbing.$(O) \
		src/util/error.$(O) src/util/out.$(O) src/util/help.$(O)
	cd src; $(GC) iolaus-changes.go private-cs.go

src/iolaus-pull.$(O): src/iolaus/core.$(O) src/iolaus/promptcommit.$(O) \
	src/iolaus/test.$(O) src/util/help.$(O) src/util/out.$(O)

src/iolaus-push.$(O): src/iolaus/core.$(O) src/iolaus/promptcommit.$(O) \
	src/iolaus/test.$(O) src/util/help.$(O) src/util/out.$(O)

src/iolaus-record.$(O): src/gotgo/slice(git.Commitish).$(O) \
	src/iolaus/core.$(O) \
	src/iolaus/prompt.$(O) src/iolaus/test.$(O) \
	src/util/error.$(O) src/util/help.$(O) src/util/out.$(O)

src/iolaus-whatsnew.$(O): src/iolaus/core.$(O) src/iolaus/prompt.$(O) \
	src/util/error.$(O) src/util/help.$(O)

src/util/cook.$(O): src/util/exit.$(O)
src/util/error.$(O): src/util/cook.$(O)
src/util/out.$(O): src/util/cook.$(O)

ifneq ($(strip $(shell which gotgo)),)
src/private-cs.go: $(srcpkgdir)/gotgo/slice.got
	gotgo --package-name=main --prefix=cs -o "$@" "$<" ./git/git.CommitHash
src/util/slice(string).go: $(srcpkgdir)/gotgo/slice.got
	gotgo -o "$@" "$<" string
src/git/gotgo/slice(git.CommitHash).go: $(srcpkgdir)/gotgo/slice.got
	mkdir -p src/git/gotgo/
	gotgo -o "$@" "$<" "../git.CommitHash"
src/iolaus/gotgo/box(git.CommitHash,git.Commitish).go: $(srcpkgdir)/gotgo/box.got
	mkdir -p src/iolaus/gotgo/
	gotgo -o "$@" "$<" ../../git/git.CommitHash ../../git/git.Commitish
src/gotgo/slice(git.Commitish).go: $(srcpkgdir)/gotgo/slice.got
	mkdir -p src/gotgo/
	gotgo -o "$@" "$<" ../git/git.Commitish
src/util/slicePatienceElem.go: $(srcpkgdir)/gotgo/slice.got
	gotgo --package-name=patience --prefix pe -o "$@" "$<" PatienceElem
src/util/gotgo/slice(int).go: $(srcpkgdir)/gotgo/slice.got
	mkdir -p src/util/gotgo/
	gotgo $< 'int' > "$@"
src/util/sliceStringChunk.go: $(srcpkgdir)/gotgo/slice.got
	gotgo --package-name=patience --prefix sc -o "$@" "$<" StringChunk
src/util/slicePatienceElems.go: $(srcpkgdir)/gotgo/slice.got
	gotgo --package-name=patience --prefix pes -o "$@" "$<" '[]PatienceElem'
endif
src/util/patience.$(O): src/util/patience.go \
	src/util/slicePatienceElem.go \
	src/util/gotgo/slice(int).$(O) src/git/color.$(O) \
	src/util/slicePatienceElems.go \
	src/util/sliceStringChunk.go
	cd src/util && $(O)g -o patience.$(O) patience.go \
		slicePatienceElem.go slicePatienceElems.go sliceStringChunk.go

src/util/patienceTypes.$(O): src/util/patienceTypes.go src/git/color.$(O)

installbins:  $(bindir)/iolaus-initialize$(EXE) \
	$(bindir)/iolaus-pull$(EXE) $(bindir)/iolaus-push$(EXE) \
	$(bindir)/iolaus-record$(EXE) \
	$(bindir)/iolaus-whatsnew$(EXE)  $(bindir)/iolaus-changes$(EXE)
