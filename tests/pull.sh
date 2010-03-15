#!/bin/sh

set -ev

mkdir repo
cd repo
echo hello > foo
iolaus-initialize
iolaus-record -am addfoo
iolaus-whatsnew | grep foo && exit 1

cd ..
mkdir new
cd new
iolaus-initialize
iolaus-pull --dry-run ../repo
iolaus-pull --all ../repo
ls
git log
grep hello foo

# FIXME: need to test iolaus-pull --interactive
