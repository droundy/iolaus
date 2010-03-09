#!/bin/sh

set -ev

# verify that the harness is working properly

cd $(dirname $(which iolaus-initialize))
pwd
test -d ../src
test -f ../src/iolaus-initialize.go
