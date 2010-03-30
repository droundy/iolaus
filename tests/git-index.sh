#!/bin/sh

set -ev

# This is a test that GIT_INDEX_FILE is undefined when tests are run.

mkdir test1
cd test1
iolaus-initialize

cat > .test <<EOF
#!/bin/sh

set -ev
echo I am running test...
if test "\$GIT_INDEX_FILE" == ""; then
  echo good
else
  exit 1
fi
EOF
chmod +x .test

iolaus-whatsnew --debug

ls -l .test
iolaus-record -am testtest --debug > out 2>&1
cat out

grep good out
