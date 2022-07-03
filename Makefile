GO:=go
GOFLAGS:= -ldflags '-s -w'
PREFIX:=/usr/local
BINPATH=$(PREFIX)/bin

# hehe
all: awl

awl: .
	$(GO) build -o awl $(GOFLAGS) .

test:
	$(GO) test ./...

fmt:
	$(GO) fmt

vet:
	$(GO) vet

lint: fmt vet

install: awl
	install awl $(BINPATH)  || echo "You probably need to run `sudo make install`"

clean:
	$(GO) clean