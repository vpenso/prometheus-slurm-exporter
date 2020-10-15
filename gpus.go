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
	"github.com/prometheus/common/log"
	"io/ioutil"
	"os/exec"
	"strings"
)

type GPUsMetrics struct {
	alloc float64
	idle  float64
	other float64
	total float64
}

func GPUsGetMetrics() *GPUsMetrics {
	return ParseGPUsMetrics()
}

func ParseAllocatedGPUs() float64 {
	return 0.0 // TODO Implement
}

func ParseIdleGPUs() float64 {
	return 0.0 // TOOD Implement
}

func ParseOtherGPUs() float64 {
	return 0.0 // TODO Implement
}

func ParseTotalGPUs() float64 {
	args := []string{"-h", "-o \"%n %G\""}
	output := string(Execute("sinfo", args))
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			log.Infof("Line %s: ", line)
			descriptor := strings.Split(line, " ")[0]
			log.Infof("Descriptor %s: ", descriptor)
		}
	}

	return 0.0
}

func ParseGPUsMetrics() *GPUsMetrics {
	var gm GPUsMetrics
	gm.alloc = ParseAllocatedGPUs()
	gm.idle = ParseIdleGPUs()
	gm.other = ParseOtherGPUs()
	gm.total = ParseTotalGPUs()
	return &gm
}

// Execute the sinfo command and return its output
func Execute(command string, arguments []string) []byte {
	cmd := exec.Command(command, arguments...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

func NewGPUsCollector() *GPUsCollector {
	return &GPUsCollector{
		alloc: prometheus.NewDesc("slurm_gpus_alloc", "Allocated GPUs", nil, nil),
		idle:  prometheus.NewDesc("slurm_gpus_idle", "Idle GPUs", nil, nil),
		other: prometheus.NewDesc("slurm_gpus_other", "Mix GPUs", nil, nil),
		total: prometheus.NewDesc("slurm_gpus_total", "Total GPUs", nil, nil),
	}
}

type GPUsCollector struct {
	alloc *prometheus.Desc
	idle  *prometheus.Desc
	other *prometheus.Desc
	total *prometheus.Desc
}

// Send all metric descriptions
func (cc *GPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
}
func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) {
	cm := GPUsGetMetrics()
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
}
