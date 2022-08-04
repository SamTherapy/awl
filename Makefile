# SPDX-License-Identifier: BSD-3-Clause
# BSD/POSIX makefile

include template.mk

$(PROG):
	$(GO) build -o $(PROG) $(GOFLAGS) .

## install: installs awl and the manpage, RUN AS ROOT
install: all
	install -m755 $(PROG) $(PREFIX)/$(BIN)
	install -m644 doc/$(PROG).1 $(MAN)/man1

.PHONY: install