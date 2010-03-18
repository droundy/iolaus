#!/bin/sh

set -ev

mkdir repo
cd repo
echo hello > foo
iolaus-initialize
iolaus-record -am addfoo
iolaus-whatsnew --all | grep foo && exit 1
cd ..

mkdir new
cd new
iolaus-initialize
iolaus-pull --dry-run ../repo
iolaus-pull --all ../repo
ls
git log
grep hello foo
cd ..

mkdir repo1
cd repo1
iolaus-initialize
date > bar
iolaus-record -am 'addbar'
git status && true
iolaus-pull --debug --all ../repo
ls
iolaus-whatsnew
grep hello foo
#
#x
git log | grep Merge
cd ..

# FIXME: need to test iolaus-pull --interactive
