#!/bin/sh

set -ev

# verify that the harness is working properly

HERE=`pwd`
cd ../../..
GRANDPARENT=`pwd`
cd "$HERE"
# The current directory
grep "$GRANDPARENT" $(which iolaus-initialize)

cd $(dirname $(which iolaus-initialize))
pwd
test -d ../src
test -f ../src/iolaus-initialize.go
