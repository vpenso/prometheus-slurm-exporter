PROJECT_NAME = prometheus-slurm-exporter
SHELL := $(shell which bash) -eu -o pipefail

GOPATH := $(shell pwd)/go/modules
GOBIN := bin/$(PROJECT_NAME)
GOFILES := $(shell ls *.go)

.PHONY: build
build: test $(GOBIN)

$(GOBIN): go/modules/pkg/mod $(GOFILES)
	mkdir -p bin
	@echo "Building $(GOBIN)"
	go build -v -o $(GOBIN)

go/modules/pkg/mod: go.mod
	go mod download

.PHONY: test
test: go/modules/pkg/mod $(GOFILES)
	go test -v

run: $(GOBIN)
	$(GOBIN)

clean:
	go clean -modcache
	rm -fr bin/ go/
