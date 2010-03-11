#!/bin/sh

set -ev

pwd
for src in ../../../src/iolaus-*.go; do
    echo grep exit.Exit $src
    grep '\.Exit' $src
done

