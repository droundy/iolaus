#!/bin/sh

set -ev

mkdir repo
cd repo
echo hello > foo
iolaus-initialize
iolaus-record -am addfoo

cd ..
mkdir new
cd new
iolaus-initialize
iolaus-pull --dry-run ../repo
iolaus-pull --all ../repo
ls
git log
grep hello foo
