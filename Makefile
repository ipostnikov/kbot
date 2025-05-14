VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS=linux

format:
	gofmt -s -w ./

lint:
	go vet ./...
	staticcheck ./...

test:
	go test -v


build:
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${shell dpkg --print-architecture} go build -v -o kbot -ldflags "-X="github.com/ipostnikov/kbot/cmd.appVersion=${VERSION}

.PHONY: clean

clean:
	rm -rf kbot

.PHONY: tools

tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest