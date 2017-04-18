# Prometheus Slurm Exporter

Prometheus collector and exporter for metric form the [Slurm](https://slurm.schedmd.com/overview.html) resource scheduling system.

Cf. Prometheus documentation:

* [Package Documentation](https://godoc.org/github.com/prometheus/client_golang/prometheus)
* [Metric Types](https://prometheus.io/docs/concepts/metric_types/)
* [Writing Exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)

Install the Prometheus [Go client library](https://github.com/prometheus/client_golang)

    >>> apt install -t jessie-backports golang-github-prometheus-client-golang-dev

Use the [Makefile](Makefile) to build and test the code.

## License

Copyright 2017 Victor Penso

This is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.


