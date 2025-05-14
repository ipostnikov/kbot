
# app name and registry
# APP := $(shell basename $(shell git remote get-url origin))
APP := kbot
VERSION := $(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
REGISTRY := ipostnikov
IMAGE_TAG := $(REGISTRY)/$(APP):$(VERSION)-$(TARGETARCH)

# variables for build
GOOS ?= linux
GOARCH ?= amd64

format:
	@echo "Formatting Go code..."
	@gofmt -s -w ./

get:
	@go get

install_tools:
	@go install honnef.co/go/tools/cmd/staticcheck@latest

# run vet and staticcheck
vet:
	@go vet ./...
	@staticcheck ./...

# run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# do build
build: format get
	@echo "Building version..."
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APP) -ldflags "-X=github.com/ipostnikov/kbot/cmd.appVersion=$(VERSION)"

# do build specifying os and arch
linux:
	@echo "Building for Linux amd64..."
	$(MAKE) build GOOS=linux GOARCH=amd64

windows:
	@echo "Building for Windows amd64..."
	$(MAKE) build GOOS=windows GOARCH=amd64

macos:
	@echo "Building for Macos amd64..."
	$(MAKE) build GOOS=darwin GOARCH=amd64

arm64:
	@echo "Building for ARM64..."
	@$(MAKE) build GOOS=linux GOARCH=arm64

# build docker image
image:
	@echo "Building Docker image..."
	@docker build . -t $(IMAGE_TAG)

# push docker image

push:
	@echo "Pushing Docker image..."
	@docker push $(IMAGE_TAG)

# clean artifacts

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(APP)
	@docker rmi $(IMAGE_TAG) || true


