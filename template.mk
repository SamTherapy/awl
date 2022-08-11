# SPDX-License-Identifier: BSD-3-Clause
# Template for the BSD/GNU makefiles

HASH ?= `git describe --always --dirty --broken || echo "UNKNOWN"`
VERSION ?= "git-$(HASH)"

SOURCES ?= $(shell find . -name "*.go" -type f ! -name '*_test*')
TEST_SOURCES ?= $(shell find . -name "*_test.go" -type f)

CGO_ENABLED ?= 0
GO ?= go
TEST ?= $(GO) test -v -race
COVER ?= $(GO) tool cover
GOFLAGS ?= -ldflags "-s -w -X=main.version=$(VERSION)" -trimpath
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


## test: run go test
test: $(TEST_SOURCES)
	$(GO) test -cover -coverprofile=coverage/coverage.out ./...

test-ci:
	$(TEST)

## fuzz: runs fuzz tests
fuzz:
	cd cli
	$(TEST) -fuzz=FuzzFlags -fuzztime 10000x
	$(TEST) -fuzz=FuzzDig -fuzztime 10000x
	$(TEST) -fuzz=FuzzParseArgs -fuzztime 10000x
	cd ..

coverage/coverage.out: test
	$(COVER) -func=coverage/coverage.out
	$(COVER) -html=coverage/coverage.out -o coverage/cover.html

## cover: generates test coverage, output as HTML
cover: coverage/coverage.out

fmt:
	gofmt -w -s .

vet:
	$(GO) vet ./...

## lint: lint awl, using fmt, vet and golangci-lint
lint: fmt vet
	golangci-lint run --fix

## clean: clean the build files
clean:
	$(GO) clean
	rm doc/$(PROG).1

## help: Prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: clean lint test fmt vet help