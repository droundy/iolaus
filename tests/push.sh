#!/bin/sh

set -ev

mkdir repo
cd repo
echo hello > foo
iolaus-initialize
iolaus-record -am addfoo

cd ..
git clone repo new

cd new
echo bye > foo
iolaus-record -am modfoo
iolaus-push --dry-run > out
cat out
grep modfoo out
iolaus-push --all

cd ../repo
git reset --hard

# FIXME: here is a bug!
grep bye foo
