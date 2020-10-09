PROJECT_NAME = prometheus-slurm-exporter
ifndef GOPATH
	GOPATH=$(shell pwd):/usr/share/gocode
endif
GOFILES=accounts.go cpus.go main.go nodes.go queue.go scheduler.go users.go
GOBIN=bin/$(PROJECT_NAME)

build:
	mkdir -p $(shell pwd)/bin
	@echo "Build $(GOFILES) to $(GOBIN)"
	@GOPATH=$(GOPATH) go build -o $(GOBIN) $(GOFILES)

test:
	@GOPATH=$(GOPATH) go test -v *.go

run:
	@GOPATH=$(GOPATH) go run $(GOFILES)

clean:
	if [ -f ${GOBIN} ] ; then rm -f ${GOBIN} ; fi
