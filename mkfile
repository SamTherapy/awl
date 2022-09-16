# SPDX-License-Identifier: BSD-3-Clause
# Plan 9 mkfile

GO = go
PROG = awl

GOFLAGS = -ldflags=-s -ldflags=-w -ldflags=-X=main.version=PLAN9 -trimpath

CGO_ENABLED = 0

all:V: $PROG

$PROG:
  $GO build $GOFLAGS -o $target .

install:V:
  $GO install $GOFLAGS  .
  # cp doc/$PROG.1 /sys/man/1/$PROG

test:V:
  $GO test -v -cover ./...

fmt:V:
  gofmt -w -s .

vet:V:
  $GO vet ./...

lint:V: fmt vet

clean:V:
  $GO clean

nuke:V: clean
