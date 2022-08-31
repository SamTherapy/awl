# SPDX-License-Identifier: BSD-3-Clause

include template.mk

ifeq ($(OS),Windows_NT)
	EXE := $(PROG).exe
else
	EXE := $(PROG)

endif

$(PROG): $(SOURCES)
	$(GO) build -o $(EXE) $(GOFLAGS) .

## install: installs awl
ifeq ($(OS),Windows_NT)
install:
	$(GO) install $(GOFLAGS) .
else
install: all
	install -Dm755 $(PROG) $(DESTDIR)$(PREFIX)/$(BIN)/$(PROG)
	install -Dm644 doc/$(PROG).1 $(DESTDIR)$(MAN)/man1/$(PROG).1
endif

.PHONY: install