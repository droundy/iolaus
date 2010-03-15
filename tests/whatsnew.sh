#!/bin/sh

set -ev

mkdir test1
cd test1
iolaus-initialize

echo goodbye > aaa
echo hello > b

echo yq | iolaus-whatsnew --interactive
echo yq | iolaus-whatsnew --interactive | grep "Added bbb" && exit 1
echo yq | iolaus-whatsnew --interactive | grep "Added aaa"

cd ..
