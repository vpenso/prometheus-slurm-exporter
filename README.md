# Prometheus Slurm Exporter

Prometheus collector and exporter for metrics extracted from the [Slurm](https://slurm.schedmd.com/overview.html) resource scheduling system.

## Exported Metrics

### State of the Nodes

* **Allocated**: nodes which has been allocated to one or more jobs.
* **Completing**: all jobs associated with these nodes are in the process of being completed.
* **Down**:  nodes which are unavailable for use.
* **Fail**: these nodes are expected to fail soon and are unavailable for use per system administrator request.
* **Error**: nodes which are currently in an error state and not capable of running any jobs.
* **Idle**: nodes not allocated to any jobs and thus available for use.
* **Maint**: nodes which are currently marked with the __maintenance__ flag.
* **Mixed**: nodes which have some of their CPUs ALLOCATED while others are IDLE.
* **Resv**: these nodes are in an advanced reservation and not generally available.

[Information extracted from the SLURM **sinfo** command](https://slurm.schedmd.com/sinfo.html)

### Status of the Jobs

* **PENDING**: Jobs awaiting for resource allocation.
* **RUNNING**: Jobs currently allocated.
* **SUSPENDED**: Job has an allocation but execution has been suspended and CPUs have been released for other jobs.
* **CANCELLED**: Jobs which were explicitly cancelled by the user or system administrator.
* **COMPLETING**: Jobs which are in the process of being completed.
* **COMPLETED**: Jobs have terminated all processes on all nodes with an exit code of zero.
* **CONFIGURING**: Jobs have been allocated resources, but are waiting for them to become ready for use.
* **FAILED**: Jobs terminated with a non-zero exit code or other failure condition.
* **TIMEOUT**: Jobs terminated upon reaching their time limit.
* **PREEMPTED**: Jobs terminated due to preemption.
* **NODE_FAIL**: Jobs terminated due to failure of one or more allocated nodes.

[Information extracted from the SLURM **squeue** command](https://slurm.schedmd.com/squeue.html)

### Scheduler Information

* **Server Thread count**: The number of current active ``slurmctld`` threads. 
* **Queue size**: The length of the scheduler queue.
* **Last cycle**: Time in microseconds for last scheduling cycle.
* **Mean cycle**: Mean of scheduling cycles since last reset.
* **Cycles per minute**: Counter of scheduling executions per minute.
* **(Backfill) Last cycle**: Time in microseconds of last backfilling cycle.
* **(Backfill) Mean cycle**: Mean of backfilling scheduling cycles in microseconds since last reset.
* **(Backfill) Depth mean**: Mean of processed jobs during backfilling scheduling cycles since last reset.

[Information extracted from the SLURM **sdiag** command](https://slurm.schedmd.com/sdiag.html)

## How to install the exporter

### Debian

Install the Prometheus [Go client library](https://github.com/prometheus/client_golang)

    >>> apt install -t jessie-backports golang-github-prometheus-client-golang-dev

Use the [Makefile](Makefile) to build and test the code.

### CentOS

Under CentOS not all the GOlang dependencies are available as packages.

In order to use the [Makefile](Makefile) provided with this repository you can proceed as follows:

1. Install the GOlang compiler plus GIT and make:
```bash
yum install git golang-bin make
```

2. Clone this repo and export the *GOPATH* environment variable:
```bash
git clone https://github.com/vpenso/prometheus-slurm-exporter.git
cd prometheus-slurm-exporter
export GOPATH=$(pwd):/usr/share/gocode
```

3. Install all the necessary GOlang dependencies:
```bash
go get github.com/prometheus/client_golang
go get github.com/prometheus/client_model
go get github.com/prometheus/common
go get github.com/prometheus/procfs
go get github.com/beorn7/perks/quantile
go get github.com/golang/protobuf/proto
go get github.com/matttproud/golang_protobuf_extensions/pbutil
go get github.com/sirupsen/logrus
go get gopkg.in/alecthomas/kingpin.v2
```

**NOTE**: all these packages will be saved under the ``src`` subdirectory of the Slurm exporter project.

In some cases, the ``go get`` command may report the following message:
```
package github.com/[...]/[...]: no buildable Go source files
```
This should be considered harmless and it will not affect the build process.

4. Build the executable:
```bash
make build
```

## Prometheus references

* [GOlang Package Documentation](https://godoc.org/github.com/prometheus/client_golang/prometheus)
* [Metric Types](https://prometheus.io/docs/concepts/metric_types/)
* [Writing Exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)
* [Available Exporters](https://prometheus.io/docs/instrumenting/exporters/)


## License

Copyright 2017 Victor Penso

This is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.


