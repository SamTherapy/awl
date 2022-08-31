# SPDX-License-Identifier: BSD-3-Clause
# Template for the BSD/GNU makefiles

HASH ?= `git describe --always --dirty --broken | sed 's/\([^-]*-g\)/r\1/;s/-/./g' || echo "UNKNOWN"`

SOURCES ?= $(shell find . -name "*.go" -type f ! -name '*_test*')
TEST_SOURCES ?= $(shell find . -name "*_test.go" -type f)

CGO_ENABLED ?= 0
GO ?= go
TEST ?= $(GO) test -race
COVER ?= $(GO) tool cover
GOFLAGS ?= -ldflags "-s -w -X=main.version=$(HASH)" -trimpath
DESTDIR :=

PREFIX ?= /usr/local
BIN ?= bin

SCDOC ?= scdoc
MAN ?= $(PREFIX)/share/man

PROG ?= awl

# hehe
all: $(PROG) doc/$(PROG).1

doc/$(PROG).1: doc/$(PROG).1.scd
	$(SCDOC) <doc/$(PROG).1.scd >doc/$(PROG).1

.PHONY: fmt
fmt:
	gofmt -w -s .

.PHONY: vet
vet:
	$(GO) vet ./...

## lint: lint awl, using fmt, vet and golangci-lint
.PHONY: lint
lint: fmt vet
	golangci-lint run --fix

## test: run go test
test: $(TEST_SOURCES)
	$(TEST) -cover -coverprofile=coverage/coverage.out ./...

test-ci:
	$(TEST) -v

## fuzz: runs fuzz tests
fuzz: $(TEST_SOURCES)
	cd cli
	$(TEST) -fuzz=FuzzFlags -fuzztime 10000x
	$(TEST) -fuzz=FuzzDig -fuzztime 10000x
	$(TEST) -fuzz=FuzzParseArgs -fuzztime 10000x
	cd ..

.PHONY: full_test
full_test: test fuzz

coverage/coverage.out: test
	$(COVER) -func=coverage/coverage.out
	$(COVER) -html=coverage/coverage.out -o coverage/cover.html

## cover: generates test coverage, output as HTML
cover: coverage/coverage.out

## clean: clean the build files
.PHONY: clean
clean:
	$(GO) clean
	rm doc/$(PROG).1

## help: Prints this help message
.PHONY: help
help:
	@echo "Usage: "
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'