PROJECT_NAME = prometheus-slurm-exporter
ROOTDIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))
SRCDIR := $(ROOTDIR)
GOPATH := $(SRCDIR)/go/modules
GOBIN := bin/$(PROJECT_NAME)
GOFILES := $(wildcard $(SRCIR)/*.go)
GOFLAGS = -v
GOTESTFLAGS =

all: $(GOBIN)

download: go/modules/pkg/mod

.PHONY: build
build: $(GOBIN)

$(GOBIN): $(GOFILES)
	mkdir -p bin
	@echo "Building $(GOBIN)"
	go build $(GOFLAGS) -o $(GOBIN)

go/modules/pkg/mod: go.mod
	go mod download

.PHONY: test
test: go/modules/pkg/mod $(GOFILES)
	go test -v

run: $(GOBIN)
	$(GOBIN)

.PHONY: clean
clean:
	go clean -modcache
	rm -fr bin/ go/
