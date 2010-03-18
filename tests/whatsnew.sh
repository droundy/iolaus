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

echo allas > aaa
echo bbb > b
iolaus-whatsnew
iolaus-whatsnew | grep hello
iolaus-whatsnew | grep bbb

# FIXME: The following requires us to make wh accept files
iolaus-whatsnew aaa
#iolaus-whatsnew aaa | grep bbb && exit 1
#iolaus-whatsnew b | grep hello && exit 1
iolaus-whatsnew aaa | grep goodbye | grep -- -

cd ..
