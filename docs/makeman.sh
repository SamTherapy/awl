#!/bin/sh

set -e

rm -f docs/awl.1.gz
scdoc <docs/awl.1.scd >docs/awl.1
gzip -9 -n docs/awl.1
