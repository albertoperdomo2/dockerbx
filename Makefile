VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BINARY_NAME := dockerbx

.PHONY: build
build:
	go build -buildmode=pie -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(BINARY_NAME) cmd/dockerbx/main.go

.PHONY: install
install: build
	mv $(BINARY_NAME) /usr/local/bin/

.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)_darwin_arm64 $(BINARY_NAME)_darwin_amd64

.PHONY: release
release: clean
	GOOS=darwin GOARCH=arm64 go build -buildmode=pie -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(BINARY_NAME)_darwin_arm64 cmd/dockerbx/main.go
	GOOS=darwin GOARCH=amd64 go build -buildmode=pie -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(BINARY_NAME)_darwin_amd64 cmd/dockerbx/main.go
	shasum -a 256 $(BINARY_NAME)_darwin_arm64 $(BINARY_NAME)_darwin_amd64 > checksums.txt

.PHONY: all
all: clean build release
