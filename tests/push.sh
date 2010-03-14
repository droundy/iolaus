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
iolaus-push --all

cd ../repo
# FIXME: here is a bug or missing feature: it'd be nice if push
# updated the working directory like darcs does, or refused to run on
# non-bare repositories.
git reset --hard

grep bye foo

# FIXME: need to test iolaus-push --interactive
