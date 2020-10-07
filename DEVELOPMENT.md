Setup the development environment on a node with access to the Slurm user
command-line interface, in particular with the `sinfo`, `squeue`, and `sdiag`
commands.

Install Go from source:

```bash
export VERSION=1.13 OS=linux ARCH=amd64
wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz
tar -xzvf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PWD/go/bin:$PATH
```

_Alternatively install Go from a package of your Linux distribution._

Use Git to clone the source code of the exporter, and download all Go dependency
modules:

```bash
# clone the source code
git clone https://github.com/vpenso/prometheus-slurm-exporter.git
cd prometheus-slurm-exporter
# download dependencies
export GOPATH=$PWD/go/modules
go mod download
```

### Build

Build the exporter:

```bash
go build -o bin/prometheus-slurm-exporter {main,cpus,nodes,queue,scheduler}.go
```

Start the exporter (foreground), and query all metrics:

```bash
bin/prometheus-slurm-exporter
...
# query all metrics (default port)
curl http://localhost:8080/metrics
```

### Tests

Run all tests included in `_test.go` files:

```bash
go test -v *.go
```

### Development

References:

* [GOlang Package Documentation](https://godoc.org/github.com/prometheus/client_golang/prometheus)
* [Metric Types](https://prometheus.io/docs/concepts/metric_types/)
* [Writing Exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)
* [Available Exporters](https://prometheus.io/docs/instrumenting/exporters/)


<https://prometheus.io/docs/instrumenting/writing_exporters>
