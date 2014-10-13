BIN=cbot

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOCMD) get ./... 
GOFMT=gofmt -w
 
default: build

build:
#	GOARCH=amd64 GOOS=linux $(GOBUILD) -a -o bin/linux-amd64/$(BIN)
	GOARCH=386 GOOS=linux $(GOBUILD) -a -o bin/linux-386/$(BIN)
#	GOARCH=amd64 GOOS=darwin $(GOBUILD) -a -o bin/darwin-amd64/$(BIN)

format:
	$(GOFMT) ./**/*.go

clean:
	$(GOCLEAN)

test:
	$(GODEP) && $(GOTEST) -v ./...
