# Development

Setup the development environment on a node with access to the Slurm user
command-line interface, in particular with the `sinfo`, `squeue`, and `sdiag`
commands.

## Install Go from source

```bash
export VERSION=1.15 OS=linux ARCH=amd64
wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz
tar -xzvf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PWD/go/bin:$PATH
```

_Alternatively install Go using the packaging system of your Linux distribution._

## Clone this repository and build

Use Git to clone the source code of the exporter, run all the tests and build the binary:

```bash
# clone the source code
git clone https://github.com/vpenso/prometheus-slurm-exporter.git
cd prometheus-slurm-exporter
make
```

To just run the tests:

```bash
make test
```

Start the exporter (foreground), and query all metrics:

```bash
./bin/prometheus-slurm-exporter
```

If you wish to run the exporter on a different port, or the default port (8080) is already in use, run with the following argument:

```bash
./bin/prometheus-slurm-exporter --listen-address="0.0.0.0:<port>"
...

# query all metrics (default port)
curl http://localhost:8080/metrics
```

## References

* [GOlang Package Documentation](https://godoc.org/github.com/prometheus/client_golang/prometheus)
* [Metric Types](https://prometheus.io/docs/concepts/metric_types/)
* [Writing Exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)
* [Available Exporters](https://prometheus.io/docs/instrumenting/exporters/)
