/* Copyright 2021 Victor Penso

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
)

type PartitionMetrics struct {
	allocated float64
	idle      float64
	other     float64
	pending   float64
	running   float64
	total     float64
}

func ParsePartitionsMetrics() map[string]*PartitionMetrics {
	partitions := make(map[string]*PartitionMetrics)

	sinfo, err := SinfoData()
	if err != nil {
		log.Errorf("partitions: %v", err)
		return partitions
	}
	for _, n := range sinfo.Nodes {
		for _, p := range n.Partitions {
			if _, ok := partitions[p]; !ok {
				partitions[p] = &PartitionMetrics{}
			}
			partitions[p].idle += float64(n.IdleCPUs)
			partitions[p].other += float64(n.OtherCPUs())
			partitions[p].total += float64(n.CPUs)
		}
	}

	squeue, err := SqueueData()
	if err != nil {
		log.Errorf("partitions: %v", err)
		return partitions
	}
	// Allocated CPUs are derived from running jobs rather than the nodes'
	// alloc_cpus field: a node shared across several partitions would
	// otherwise have its allocated CPUs counted into each of them.
	for _, j := range squeue.Jobs {
		p, ok := partitions[j.Partition]
		if !ok {
			continue
		}
		switch j.State() {
		case "PENDING":
			p.pending++
		case "RUNNING":
			p.running++
			p.allocated += float64(j.CPUs)
		}
	}

	return partitions
}

type PartitionsCollector struct {
	allocated *prometheus.Desc
	idle      *prometheus.Desc
	other     *prometheus.Desc
	pending   *prometheus.Desc
	running   *prometheus.Desc
	total     *prometheus.Desc
}

func NewPartitionsCollector() *PartitionsCollector {
	labels := []string{"partition"}
	return &PartitionsCollector{
		allocated: prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUs for partition", labels, nil),
		idle:      prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUs for partition", labels, nil),
		other:     prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUs for partition", labels, nil),
		pending:   prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", labels, nil),
		running:   prometheus.NewDesc("slurm_partition_jobs_running", "Running jobs for partition", labels, nil),
		total:     prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUs for partition", labels, nil),
	}
}

func (pc *PartitionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pc.allocated
	ch <- pc.idle
	ch <- pc.other
	ch <- pc.pending
	ch <- pc.running
	ch <- pc.total
}

func (pc *PartitionsCollector) Collect(ch chan<- prometheus.Metric) {
	pm := ParsePartitionsMetrics()
	for p := range pm {
		if pm[p].allocated > 0 {
			ch <- prometheus.MustNewConstMetric(pc.allocated, prometheus.GaugeValue, pm[p].allocated, p)
		}
		if pm[p].idle > 0 {
			ch <- prometheus.MustNewConstMetric(pc.idle, prometheus.GaugeValue, pm[p].idle, p)
		}
		if pm[p].other > 0 {
			ch <- prometheus.MustNewConstMetric(pc.other, prometheus.GaugeValue, pm[p].other, p)
		}
		if pm[p].pending > 0 {
			ch <- prometheus.MustNewConstMetric(pc.pending, prometheus.GaugeValue, pm[p].pending, p)
		}
		if pm[p].running > 0 {
			ch <- prometheus.MustNewConstMetric(pc.running, prometheus.GaugeValue, pm[p].running, p)
		}
		if pm[p].total > 0 {
			ch <- prometheus.MustNewConstMetric(pc.total, prometheus.GaugeValue, pm[p].total, p)
		}
	}
}
