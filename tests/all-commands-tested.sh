#!/bin/sh

RETVAL=0

for cmd in ../../../bin/*; do
    echo -n Checking if `basename $cmd` is tested... ' '
    if grep `basename $cmd` ../../*.sh > /dev/null; then
        echo yes.
    else
        echo NO!
        RETVAL=1
    fi
done

exit $RETVAL
