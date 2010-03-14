#!/bin/sh

set -ev

mkdir test1
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
iolaus-initialize

date > aaa
echo hello b > b

iolaus-whatsnew

echo ny | iolaus-record --interactive --patch 'Hello world'

iolaus-whatsnew
iolaus-whatsnew | grep 'Added b' && exit 1
iolaus-whatsnew | grep aaa
