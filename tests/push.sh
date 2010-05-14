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
# verify that repo doesn't have foo yet
grep bye ../repo/foo && exit 1
iolaus-push --dry-run > out
cat out
grep modfoo out
# verify that repo still doesn't have foo, since we only pushed --dry-run
grep bye ../repo/foo && exit 1
# we can't push to a non-bare repository by default
iolaus-push --all && exit 1

cd ../repo
git config --bool core.bare true
cd ../new
iolaus-push --all

cd ../repo
# To check that the repository has been changed, we'll just make it
# unbare.  This is a stupid way to check this...
git config --bool core.bare false
git reset --hard
grep bye foo

# FIXME: need to test iolaus-push --interactive
