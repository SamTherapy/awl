#!/bin/sh

set -e

rm -f doc/awl.1.gz
scdoc <doc/awl.1.scd >doc/awl.1
gzip -9 -n doc/awl.1
