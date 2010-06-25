#!/bin/sh

set -ev

# This is a test of whatsnew...

mkdir test1
cd test1
iolaus-initialize

cat > file <<EOF
line 1
line 2
line 3
line 4
line 5
line 6
EOF

iolaus-record -am 'add file'

cat > file <<EOF
LINE 1
LINE 2
line 3
line 4
line 5
line 6
EOF

iolaus-whatsnew
git diff

iolaus-whatsnew > out

grep ' file 1 ' out
grep '-line 1' out
grep '-line 2' out
grep '+LINE 1' out
grep '+LINE 2' out
grep ' LINE 1' out && exit 1
grep 'line 6' out && exit 1
grep ' line 3' out
grep ' line 4' out
grep ' line 5' out
