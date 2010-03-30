#!/bin/sh

set -ev

# This is a test of whatsnew...

mkdir test1
cd test1
iolaus-initialize

echo goodbye > aaa
echo hello > b

echo yq | iolaus-whatsnew --interactive
echo yq | iolaus-whatsnew --interactive | grep "Added bbb" && exit 1
echo yq | iolaus-whatsnew --interactive | grep "Added aaa"

iolaus-record -am addfiles

env
pwd
which iolaus-whatsnew
echo allas > aaa
echo bbb > b
iolaus-whatsnew
iolaus-whatsnew | grep hello
iolaus-whatsnew | grep bbb

iolaus-whatsnew aaa
iolaus-whatsnew aaa | grep bbb && exit 1
iolaus-whatsnew b | grep goodbye && exit 1
iolaus-whatsnew aaa | grep goodbye | grep -- -

cd ..
