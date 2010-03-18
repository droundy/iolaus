#!/bin/sh

set -ev

mkdir temp
cd temp
iolaus-initialize

cat > foo <<EOF
A
B
C
D
E
F
G
H
I
J
K
L
M
N
O
P
EOF

iolaus-record -am 'addfoo'

cat > foo <<EOF
A
B
C
d
E
F
G
H
i
J
K
L
M
n
O
P
EOF

iolaus-whatsnew | grep d
iolaus-whatsnew | grep i
iolaus-whatsnew | grep n

echo znyny | iolaus-record --interactive

iolaus-whatsnew | grep d
iolaus-whatsnew | grep i && exit 1
iolaus-whatsnew | grep n

echo znyy | iolaus-record --interactive

iolaus-whatsnew | grep i && exit 1
iolaus-whatsnew | grep n && exit 1
iolaus-whatsnew | grep d
