# SPDX-License-Identifier: BSD-3-Clause
# BSD/POSIX makefile

include template.mk

$(PROG): $(SOURCES)
	$(GO) build -o $(PROG) $(GOFLAGS) .

## install: installs awl and the manpage, RUN AS ROOT
install: all
	install -Dm755 $(PROG) $(DESTDIR)$(PREFIX)/$(BIN)/$(PROG)
	install -Dm644 doc/$(PROG).1 $(DESTDIR)$(MAN)/man1/$(PROG).1

.PHONY: install