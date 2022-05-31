/* Copyright 2017 Victor Penso, Matteo Dessalvi

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
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type QueueMetrics struct {
	pending     float64
	pending_dep float64
	running     float64
	suspended   float64
	cancelled   float64
	completing  float64
	completed   float64
	configuring float64
	failed      float64
	timeout     float64
	preempted   float64
	node_fail   float64
	out_of_memory float64
}

// Returns the scheduler metrics
func QueueGetMetrics() *QueueMetrics {
	return ParseQueueMetrics(QueueData())
}

func ParseQueueMetrics(input []byte) *QueueMetrics {
	var qm QueueMetrics
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		if strings.Contains(line, ",") {
			splitted := strings.Split(line, ",")
			state := splitted[1]
			switch state {
			case "PENDING":
				qm.pending++
				if len(splitted) > 2 && splitted[2] == "Dependency" {
					qm.pending_dep++
				}
			case "RUNNING":
				qm.running++
			case "SUSPENDED":
				qm.suspended++
			case "CANCELLED":
				qm.cancelled++
			case "COMPLETING":
				qm.completing++
			case "COMPLETED":
				qm.completed++
			case "CONFIGURING":
				qm.configuring++
			case "FAILED":
				qm.failed++
			case "TIMEOUT":
				qm.timeout++
			case "PREEMPTED":
				qm.preempted++
			case "NODE_FAIL":
				qm.node_fail++
                        case "OUT_OF_MEMORY":
                                qm.out_of_memory++
			}
		}
	}
	return &qm
}

// Execute the squeue command and return its output
func QueueData() []byte {
	cmd := exec.Command("squeue", "-a", "-r", "-h", "-o %A,%T,%r", "--states=all")
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
 * Slurm queue metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

func NewQueueCollector() *QueueCollector {
	return &QueueCollector{
		pending:     prometheus.NewDesc("slurm_queue_pending", "Pending jobs in queue", nil, nil),
		pending_dep: prometheus.NewDesc("slurm_queue_pending_dependency", "Pending jobs because of dependency in queue", nil, nil),
		running:     prometheus.NewDesc("slurm_queue_running", "Running jobs in the cluster", nil, nil),
		suspended:   prometheus.NewDesc("slurm_queue_suspended", "Suspended jobs in the cluster", nil, nil),
		cancelled:   prometheus.NewDesc("slurm_queue_cancelled", "Cancelled jobs in the cluster", nil, nil),
		completing:  prometheus.NewDesc("slurm_queue_completing", "Completing jobs in the cluster", nil, nil),
		completed:   prometheus.NewDesc("slurm_queue_completed", "Completed jobs in the cluster", nil, nil),
		configuring: prometheus.NewDesc("slurm_queue_configuring", "Configuring jobs in the cluster", nil, nil),
		failed:      prometheus.NewDesc("slurm_queue_failed", "Number of failed jobs", nil, nil),
		timeout:     prometheus.NewDesc("slurm_queue_timeout", "Jobs stopped by timeout", nil, nil),
		preempted:   prometheus.NewDesc("slurm_queue_preempted", "Number of preempted jobs", nil, nil),
		node_fail:   prometheus.NewDesc("slurm_queue_node_fail", "Number of jobs stopped due to node fail", nil, nil),
                out_of_memory: prometheus.NewDesc("slurm_queue_out_of_memory", "Number of jobs stopped by oomkiller", nil, nil),
	}
}

type QueueCollector struct {
	pending     *prometheus.Desc
	pending_dep *prometheus.Desc
	running     *prometheus.Desc
	suspended   *prometheus.Desc
	cancelled   *prometheus.Desc
	completing  *prometheus.Desc
	completed   *prometheus.Desc
	configuring *prometheus.Desc
	failed      *prometheus.Desc
	timeout     *prometheus.Desc
	preempted   *prometheus.Desc
	node_fail   *prometheus.Desc
        out_of_memory *prometheus.Desc
}

func (qc *QueueCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- qc.pending
	ch <- qc.pending_dep
	ch <- qc.running
	ch <- qc.suspended
	ch <- qc.cancelled
	ch <- qc.completing
	ch <- qc.completed
	ch <- qc.configuring
	ch <- qc.failed
	ch <- qc.timeout
	ch <- qc.preempted
	ch <- qc.node_fail
        ch <- qc.out_of_memory
}

func (qc *QueueCollector) Collect(ch chan<- prometheus.Metric) {
	qm := QueueGetMetrics()
	ch <- prometheus.MustNewConstMetric(qc.pending, prometheus.GaugeValue, qm.pending)
	ch <- prometheus.MustNewConstMetric(qc.pending_dep, prometheus.GaugeValue, qm.pending_dep)
	ch <- prometheus.MustNewConstMetric(qc.running, prometheus.GaugeValue, qm.running)
	ch <- prometheus.MustNewConstMetric(qc.suspended, prometheus.GaugeValue, qm.suspended)
	ch <- prometheus.MustNewConstMetric(qc.cancelled, prometheus.GaugeValue, qm.cancelled)
	ch <- prometheus.MustNewConstMetric(qc.completing, prometheus.GaugeValue, qm.completing)
	ch <- prometheus.MustNewConstMetric(qc.completed, prometheus.GaugeValue, qm.completed)
	ch <- prometheus.MustNewConstMetric(qc.configuring, prometheus.GaugeValue, qm.configuring)
	ch <- prometheus.MustNewConstMetric(qc.failed, prometheus.GaugeValue, qm.failed)
	ch <- prometheus.MustNewConstMetric(qc.timeout, prometheus.GaugeValue, qm.timeout)
	ch <- prometheus.MustNewConstMetric(qc.preempted, prometheus.GaugeValue, qm.preempted)
	ch <- prometheus.MustNewConstMetric(qc.node_fail, prometheus.GaugeValue, qm.node_fail)
        ch <- prometheus.MustNewConstMetric(qc.out_of_memory, prometheus.GaugeValue, qm.out_of_memory)
}
