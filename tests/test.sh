#!/bin/sh

set -ev

iolaus-initialize

date > .test

iolaus-whatsnew
iolaus-whatsnew | grep 'Added .test'

iolaus-record --all --patch 'Hello world'

chmod +x .test

iolaus-record --all --patch 'Failing test' && exit 1

iolaus-record --test --all --patch 'Failing test' && exit 1

iolaus-record --no-test --all --patch 'Failing test'

cat > .test <<EOF
#!/bin/sh
true
EOF
chmod +x .test

iolaus-record -am 'passing test'

