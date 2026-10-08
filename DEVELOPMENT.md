# Development

Setup the development environment on a node with access to the Slurm user
command-line interface, in particular with the `sinfo`, `squeue`, `sacct`,
`sdiag`, and `sshare` commands. `sinfo`, `squeue` and `sacct` must support
`--json` (Slurm 23.02+); `sdiag` and `sshare` are still parsed as legacy
text output (see the comments atop `scheduler.go`/`sshare.go`).

Live Slurm access is only needed to exercise the `*GetMetrics` integration
paths and to run the exporter binary itself - the unit test suite (`make
test`) runs entirely against the JSON fixtures under `test_data/` and needs
no Slurm installation.

The optional AMD ROCm telemetry collector (`-rocm-acct`) additionally shells
out to `amd-smi` (or `rocm-smi`, via `-rocm-smi-cmd`) on the node it runs
on; neither binary is required to build or run the unit tests, only to
exercise that collector live.

## Install Go from source

```bash
export VERSION=1.21 OS=linux ARCH=amd64
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
