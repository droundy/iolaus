#!/bin/sh

set -ev

mkdir test1
cd test1
iolaus-initialize

date > aaa

iolaus-record -am 'Hello world'

rm aaa
iolaus-whatsnew | grep aaa

iolaus-record -am 'remove aaa' --debug

iolaus-whatsnew | grep aaa && exit 1
iolaus-whatsnew
