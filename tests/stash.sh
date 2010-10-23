#!/bin/sh

set -ev

mkdir test1
cd test1
iolaus-initialize

echo hello b > b

iolaus-whatsnew

iolaus-record -am 'Hello world'

date > a
test -f a
iolaus-stash --all
# Actually --all is inherently broken, since it means to *keep* all
# the changes, which is to say, to stash none of the changes.
ls -lh
# FIXME: I should make iolaus-stash remove created files!
#test -f a && exit 1

echo goodbye > b
grep goodbye b
echo nnn | iolaus-stash --interactive
cat b
grep goodbye b && exit 1
grep hello b

echo foo > b
echo nnn | iolaus-stash -m 'stash name'
cat b
grep foo b && exit 1

git stash list
git stash list | grep 'stash name'

true
