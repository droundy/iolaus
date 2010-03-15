#!/bin/sh

set -ev

mkdir test1
cd test1
iolaus-initialize

date > aaa
echo hello b > b

iolaus-whatsnew

echo yn | iolaus-record --interactive --patch 'Hello world'

iolaus-whatsnew
iolaus-whatsnew | grep aaa && exit 1
iolaus-whatsnew | grep 'Added b'

cd ..
mkdir test2
cd test2
iolaus-initialize

date > aaa
echo hello b > b

iolaus-whatsnew

echo ny | iolaus-record --interactive --patch 'Hello world'

iolaus-whatsnew
iolaus-whatsnew | grep 'Added b' && exit 1
iolaus-whatsnew | grep aaa

cd ..

mkdir test-interactive-cancel
cd test-interactive-cancel
iolaus-initialize
date > aaa
echo hello b > b
iolaus-whatsnew | grep 'Added aaa'
iolaus-whatsnew | grep 'Added b'
echo yq | iolaus-record --interactive --patch 'Hello world'
# Verify that a cancelled record doesn't affect the index!
iolaus-whatsnew | grep 'Added aaa'
iolaus-whatsnew | grep 'Added b'
cd ..

mkdir test-interactive-view
cd test-interactive-view
iolaus-initialize
date > aaa
iolaus-whatsnew | grep 'Added aaa'
echo vq | iolaus-record --interactive --patch 'Hello world' > out
cat out
grep 'Added aaa' out
cd ..
