#!/bin/sh

set -ev

../../../scripts/pdiff bad-old bad-new

../../../scripts/pdiff bad-old bad-new
