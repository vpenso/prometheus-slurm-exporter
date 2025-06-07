/* Copyright 2020 Joeri Hermans, Victor Penso, Matteo Dessalvi

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

type NPUsMetrics struct {
	alloc       float64
	idle        float64
	total       float64
	utilization float64
}

func NPUsGetMetrics() *NPUsMetrics {
	return ParseNPUsMetrics()
}

func ParseAllocatedNPUs() float64 {
	var num_npus = 0.0

	args := []string{"-a", "-X", "--format=AllocTRES", "--state=RUNNING", "--noheader", "--parsable2"}
	output := string(Execute("sacct", args))
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				line = strings.Trim(line, "\"")
				for _, resource := range strings.Split(line, ",") {
					if strings.HasPrefix(resource, "gres/npu=") {
						descriptor := strings.TrimPrefix(resource, "gres/npu=")
						if job_npus, err := strconv.ParseFloat(descriptor, 64); err == nil {
							num_npus += job_npus
						}
					}
				}
			}
		}
	}

	return num_npus
}

func ParseTotalNPUs() float64 {
	var num_npus = 0.0

	args := []string{"-h", "-o", "%n,%G"}
	output := string(Execute("sinfo", args))
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			if len(line) > 0 {
				line = strings.TrimSpace(line)
				fields := strings.Split(line, ",")
				if len(fields) >= 2 {
					gres := fields[1]
					// gres column format: comma-delimited list of resources
					for _, resource := range strings.Split(gres, ",") {
						if strings.HasPrefix(resource, "npu:") {
							// format: npu:<type>:N(S:<something>), e.g. npu:A100:2(S:0)
							parts := strings.Split(resource, ":")
							if len(parts) >= 3 {
								descriptor := strings.Split(parts[2], "(")[0]
								if node_npus, err := strconv.ParseFloat(descriptor, 64); err == nil {
									num_npus += node_npus
								}
							}
						}
					}
				}
			}
		}
	}

	return num_npus
}

func ParseNPUsMetrics() *NPUsMetrics {
	var nm NPUsMetrics
	total_npus := ParseTotalNPUs()
	allocated_npus := ParseAllocatedNPUs()
	nm.alloc = allocated_npus
	nm.idle = total_npus - allocated_npus
	nm.total = total_npus
	if total_npus > 0 {
		nm.utilization = allocated_npus / total_npus
	} else {
		nm.utilization = 0
	}
	return &nm
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

func NewNPUsCollector() *NPUsCollector {
	return &NPUsCollector{
		alloc:       prometheus.NewDesc("slurm_npus_alloc", "Allocated NPUs", nil, nil),
		idle:        prometheus.NewDesc("slurm_npus_idle", "Idle NPUs", nil, nil),
		total:       prometheus.NewDesc("slurm_npus_total", "Total NPUs", nil, nil),
		utilization: prometheus.NewDesc("slurm_npus_utilization", "Total NPU utilization", nil, nil),
	}
}

type NPUsCollector struct {
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
}

// Send all metric descriptions
func (cc *NPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.total
	ch <- cc.utilization
}
func (cc *NPUsCollector) Collect(ch chan<- prometheus.Metric) {
	nm := NPUsGetMetrics()
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, nm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, nm.idle)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, nm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, nm.utilization)
}
