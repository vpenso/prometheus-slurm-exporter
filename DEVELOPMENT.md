## Development

Setup the development environment on a node with access to the Slurm user
commnad-line interface, in particular with the `sinfo` and `squeue` commands.

Install Go from source:

```bash
export VERSION=1.13 OS=linux ARCH=amd64
wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz
tar -xzvf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PWD/go/bin:$PATH
```

_Alternatively install Go from a package of your Linux distribution._

Use Git to clone the source code the exporter, and download all Go dependency
libraries:

```bash
# clone the source code
git clone https://github.com/vpenso/prometheus-slurm-exporter.git
cd prometheus-slurm-exporter
# download dependencies
export GOPATH=$PWD/go/modules
go mod download
```

### Build

```bash
# build the exporter
go build -o bin/prometheus-slurm-exporter {main,cpus,nodes,queue,scheduler}.go
# start the exporter (foreground)
bin/prometheus-slurm-exporter
...
# query all metrics (default port)
curl http://localhost:8080/metrics
```

Run all tests included in `_test.go` files:

```bash
go test -v *.go
```
