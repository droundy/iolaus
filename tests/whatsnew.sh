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

env
pwd
which iolaus-whatsnew
echo allas > aaa
echo bbb > b
iolaus-whatsnew
iolaus-whatsnew | grep hello
iolaus-whatsnew | grep bbb

iolaus-whatsnew aaa
iolaus-whatsnew aaa | grep bbb && exit 1
iolaus-whatsnew b | grep goodbye && exit 1
iolaus-whatsnew aaa | grep goodbye | grep -- -

cd ..

rm -rf repo
mkdir repo
cd repo
echo -n > foo
echo I have 1 bottle of beer on the wall >> foo
for i in `seq 2 49`; do
    echo I have $i bottles of beer on the wall >> foo
done
cat foo
iolaus-initialize
iolaus-whatsnew
iolaus-record -am foo1..100

# Now I'll change the file
grep -v 'have . bottles' foo > foo.new
grep -v 'have 3. bottles' foo.new > foo
sed -e 's/21/twenty one/' < foo > foo.new
grep -v 'have 22 bottles' foo.new > foo
rm -f foo.new

git diff
iolaus-whatsnew

# We should be missing the thirty-something bottlese
iolaus-whatsnew | grep 'foo 1 ' # First context starts at line 1
iolaus-whatsnew | grep -- ' I have 29 bottles'
iolaus-whatsnew | grep -- ' I have 40 bottles'
for i in `seq 30 39`; do
    echo Checking we are missing $i bottles
    iolaus-whatsnew | grep -- "+I have $i bottles" && exit 1
    iolaus-whatsnew | grep -- " I have $i bottles" && exit 1
    iolaus-whatsnew | grep -- "-I have $i bottles"
done

iolaus-whatsnew | grep 'foo 10 ' # Second context starts at line 18
iolaus-whatsnew | grep 'foo 18 ' # Third context starts at line 18
# We should be missing 21-22, but have twenty one
iolaus-whatsnew | grep -- ' I have 20 bottle'
iolaus-whatsnew | grep -- '-I have 22 bottle'
iolaus-whatsnew | grep -- '-I have 21 bottle'
iolaus-whatsnew | grep -- '+I have twenty one bottle'
grep 23 foo
iolaus-whatsnew | grep -- ' I have 23 bottle'

# We should be missing the 2-9 bottles
iolaus-whatsnew | grep -- ' I have 1 bottle'
iolaus-whatsnew | grep -- ' I have 10 bottles'
for i in `seq 2 9`; do
    echo Checking we are missing $i bottles
    iolaus-whatsnew | grep -- "+I have $i bottles" && exit 1
    iolaus-whatsnew | grep -- " I have $i bottles" && exit 1
    iolaus-whatsnew | grep -- "-I have $i bottles"
done

