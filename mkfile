# SPDX-License-Identifier: BSD-3-Clause
# Plan 9 mkfile

GO = go
PROG = awl
LDFLAGS = '-s -w'
GOFLAGS = -ldflags=$LDFLAGS

CGO_ENABLED = 0

$PROG: 
  $GO build $GOFLAGS -o $PROG '-buildvcs=false' .

install: $PROG
  $GO install $GOFLAGS .
  cp doc/$PROG.1 /sys/man/1/$PROG

test:
  $GO test -v -cover -coverprofile=coverage/coverage.out ./...

fmt:
  gofmt -w -s .

vet:
  $GO vet ./...

lint: fmt vet

clean:
  $GO clean
