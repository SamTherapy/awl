# SPDX-License-Identifier: BSD-3-Clause

include template.mk

ifeq ($(OS),Windows_NT)
	EXE := $(PROG).exe
else
	EXE := $(PROG)
endif


$(PROG):
	$(GO) build -o $(EXE) $(GOFLAGS) .


ifeq ($(OS),Windows_NT)
## install: installs awl
install:
	$(GO) install $(GOFLAGS) .
else
install: all
	install -m755 $(PROG) $(PREFIX)/$(BIN)
	install -m644 doc/$(PROG).1 $(MAN)/man1
endif

.PHONY: install