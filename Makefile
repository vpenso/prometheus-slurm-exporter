PROJECT_NAME = prometheus-slurm-exporter
GOPATH=$(shell pwd):/usr/share/gocode
GOFILES=main.go nodes.go queue.go scheduler.go
GOBIN=bin/$(PROJECT_NAME)

build:
	mkdir -p $(shell pwd)/bin
	@echo "Build $(GOFILES) to $(GOBIN)"
	@GOPATH=$(GOPATH) go build -o $(GOBIN) $(GOFILES)

test:
	@GOPATH=$(GOPATH) go test -v *.go

run:
	@GOPATH=$(GOPATH) go run $(GOFILES)
