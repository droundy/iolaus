#!/bin/sh

set -ev

mkdir repo
cd repo
echo hello > foo
iolaus-initialize
iolaus-record -am addfoo
iolaus-whatsnew --all | grep foo && exit 1

date > bar
iolaus-record -am 'addbar'

iolaus-changes
iolaus-changes | grep addbar
iolaus-changes | grep addfoo
