VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: build
build:
	go build -buildmode=pie -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o dockerbx cmd/dockerbx/main.go

.PHONY: install
install: build
	mv dockerbx /usr/local/bin/

.PHONY: clean
clean:
	rm -f dockerbx
