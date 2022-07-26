# SPDX-License-Identifier: BSD-3-Clause

include template.mk

ifeq ($(OS),Windows_NT)
	EXE := $(PROG).exe
else
	EXE := $(PROG)
endif


$(PROG):
	$(GO) build -o $(EXE) $(GOFLAGS) .

## install: installs awl
install: all
ifeq ($(OS),Windows_NT)
	$(GO) install $(GOFLAGS) .
else
	install -m755 $(PROG) $(PREFIX)/$(BIN)
	install -m644 doc/$(PROG).1 $(MAN)/man1
endif

.PHONY: install