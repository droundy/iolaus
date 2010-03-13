#!/bin/sh

set -ev

iolaus-initialize

date > .test

iolaus-whatsnew
iolaus-whatsnew | grep 'Added .test'

iolaus-record --all --patch 'Hello world'

chmod +x .test

iolaus-record --debug --all --patch 'Failing test' && exit 1

true
