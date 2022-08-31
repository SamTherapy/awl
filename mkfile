# SPDX-License-Identifier: BSD-3-Clause
# Plan 9 mkfile

GO = go
PROG = awl

CGO_ENABLED = 0

$PROG:
  $GO build -ldflags="-s -w -X=main.version=PLAN9" -o $PROG .

install:
  $GO install -ldflags="-s -w -X=main.version=PLAN9" .
  cp doc/$PROG.1 /sys/man/1/$PROG

test:
  $GO test -cover -coverprofile=coverage/coverage.out ./...

fmt:
  gofmt -w -s .

vet:
  $GO vet ./...

lint: fmt vet

clean:
  $GO clean

nuke: clean
