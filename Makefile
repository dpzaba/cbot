BIN=cbot

# Go parameters
GOCMD=/usr/local/go/bin/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOCMD) get ./...
GOFMT=gofmt -w

default: build

build:
	$(GOBUILD) -a -o bin/$(BIN)

format:
	$(GOFMT) ./**/*.go

clean:
	$(GOCLEAN)

test:
	$(GODEP) && $(GOTEST) -v ./...
