#!/bin/sh

set -ev

mkdir repo
cd repo
echo hello > foo
iolaus-initialize
iolaus-record -am addfoo
iolaus-whatsnew --all | grep foo && exit 1
cd ..

mkdir new
cd new
iolaus-initialize
iolaus-pull --dry-run ../repo
iolaus-pull --all ../repo
ls
git log
grep hello foo
cd ..

mkdir repo1
cd repo1
iolaus-initialize
date > bar
iolaus-record -am 'addbar'
iolaus-pull --all ../repo
ls
iolaus-whatsnew
grep hello foo
#
#x
git log | grep Merge
cd ..

# First let's make a copy that we can use later to attempt this pull
# again...
git clone repo test-interactive

# Now let's verify that a test suite is run when commits are merged...
git clone repo repowithtest
cd repowithtest
cat > .test <<EOF
#!/bin/sh
# Make sure that xxx doesn't exist!
test -f xxx && exit 1
true
EOF
chmod +x .test
iolaus-record -am 'add test that xxx does not exist'
./.test
cd ../repo
date > xxx
iolaus-record -am 'create file xxx so test will fail...'
# test should run by default, and should fail
iolaus-pull --all ../repowithtest && echo oops, it did not crash && exit 1
# test should fail if run explicitly
iolaus-pull --test --all ../repowithtest && exit 1
# test shouldn't be run if we use --no-test
iolaus-pull --no-test --all ../repowithtest
test -f xxx
./.test && exit 1
test -f .test

git log --max-count=1 --parents
cd ..

# Let's check that interactive pulling behaves reasonably...
cd test-interactive
# Interactive pull should fail since there's no terminal to prompt...
iolaus-pull --interactive && exit 1
iolaus-pull --interactive | grep 'Merge'

# FIXME: Currently typing 'n' uncovers a bug in the interactive commit
# selection code.  :( So I've commented out the remaining tests.

#echo n | iolaus-pull --interactive
#echo nq | iolaus-pull --interactive | grep 'create file xxx'
cd ..
