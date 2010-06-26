#!/bin/sh

set -ev

mkdir repo
cd repo
echo hello > foo
iolaus-initialize
iolaus-record -am addfoo

cd ..
git clone repo new

cd new
echo bye > foo
iolaus-record -am modfoo
# verify that repo doesn't have foo yet
grep bye ../repo/foo && exit 1
iolaus-push --dry-run > out
cat out
grep modfoo out
# verify that repo still doesn't have foo, since we only pushed --dry-run
grep bye ../repo/foo && exit 1
# we can't push to a non-bare repository by default
iolaus-push --all && exit 1

cd ../repo
git config --bool core.bare true
cd ../new
iolaus-push --all

cd ../repo
# To check that the repository has been changed, we'll just make it
# unbare.  This is a stupid way to check this...
git config --bool core.bare false
git reset --hard
grep bye foo
git config --bool core.bare true
cd ..

# FIXME: need to test iolaus-push --interactive


# check that a test suite works as expected
git clone repo repowithtest
cd new
date > xxx
iolaus-record -am 'add xxx file so test will fail'
git push --all
cd ../repowithtest
cat > .test <<EOF
#!/bin/sh
# Make sure that xxx doesn't exist!
test -f xxx && exit 1
true
EOF
chmod +x .test
iolaus-record -am 'add test that xxx does not exist'
./.test
# Test should fail and prevent pushing
iolaus-push --all && exit 1
# Push should also fail if we ask for --test
iolaus-push --test --all && exit 1
# But without testing, it should work just fine
iolaus-push --no-test --all
cd ../new
# Verify that the push actually worked as expected Note that the pull
# here won't perform a test, since it's a fast-forward pull.
iolaus-pull --all
./.test && exit 1
test -f .test
cd ..
