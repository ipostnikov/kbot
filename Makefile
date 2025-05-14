APP=$(shell basename $(shell git remote get-url origin))
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
REGISTRY=ipostnikov
TARGETOS=linux
TARGETARCH=arm64

format:
	gofmt -s -w ./

get:
	go get

lint:
	go vet ./...
	staticcheck ./...

test:
	go test -v


build: format get
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${shell dpkg --print-architecture} go build -v -o kbot -ldflags "-X="github.com/ipostnikov/kbot/cmd.appVersion=${VERSION}

image:
	docker build . -t ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}

push:
	docker build . -t ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}

.PHONY: clean

clean:
	rm -rf kbot

.PHONY: tools

tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest