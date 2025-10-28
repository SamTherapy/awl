#!/bin/sh

set -e

rm -f docs/awl.1.gz
scdoc <docs/awl.1.scd >docs/awl.1
gzip -9kn docs/awl.1
gzip -9kn README.md
gzip -9kn docs/CONTRIBUTING.md
