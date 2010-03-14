#!/bin/sh

RETVAL=0

iolaus-initialize

for cmdx in ../../../bin/*; do
    cmd=`basename $cmdx`
    for flag in `$cmd --list-options`; do
        if test $flag == --list-options; then
            echo -n
        elif test $flag == --debug; then
            echo -n
        elif test $flag == --help; then
            echo -n
        else
            echo -n Checking if $cmd $flag is tested... ' '
            if grep $cmd ../../*.sh | grep -- $flag > /dev/null; then
                echo yes.
            else
                echo NO!
                RETVAL=1
            fi
        fi
    done
done

exit $RETVAL
