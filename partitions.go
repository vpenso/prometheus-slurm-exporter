/* Copyright 2020 Victor Penso

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
	"strconv"
	"strings"
)

func PartitionsData() []byte {
	cmd := exec.Command("sinfo", "-h", "-o%R,%C")
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

func PartitionsPendingJobsData() []byte {
	cmd := exec.Command("squeue", "-a", "-r", "-h", "-o%P", "--states=PENDING")
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

type PartitionMetrics struct {
	allocated float64
	idle      float64
	other     float64
	pending   float64
	total     float64
}

func ParsePartitionsMetrics() map[string]*PartitionMetrics {
	partitions := make(map[string]*PartitionMetrics)
	lines := strings.Split(string(PartitionsData()), "\n")
	for _, line := range lines {
		if strings.Contains(line, ",") {
			// name of a partition
			partition := strings.Split(line, ",")[0]
			_, key := partitions[partition]
			if !key {
				partitions[partition] = &PartitionMetrics{0, 0, 0, 0, 0}
			}
			states := strings.Split(line, ",")[1]
			allocated, _ := strconv.ParseFloat(strings.Split(states, "/")[0], 64)
			idle, _ := strconv.ParseFloat(strings.Split(states, "/")[1], 64)
			other, _ := strconv.ParseFloat(strings.Split(states, "/")[2], 64)
			total, _ := strconv.ParseFloat(strings.Split(states, "/")[3], 64)
			partitions[partition].allocated = allocated
			partitions[partition].idle = idle
			partitions[partition].other = other
			partitions[partition].total = total
		}
	}
	// get list of pending jobs by partition name
	list := strings.Split(string(PartitionsPendingJobsData()), "\n")
	for _, partition := range list {
		// accumulate the number of pending jobs
		_, key := partitions[partition]
		if key {
			partitions[partition].pending += 1
		}
	}

	return partitions
}

type PartitionsCollector struct {
	allocated *prometheus.Desc
	idle      *prometheus.Desc
	other     *prometheus.Desc
	pending   *prometheus.Desc
	total     *prometheus.Desc
}

func NewPartitionsCollector() *PartitionsCollector {
	labels := []string{"partition"}
	return &PartitionsCollector{
		allocated: prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUs for partition", labels, nil),
		idle:      prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUs for partition", labels, nil),
		other:     prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUs for partition", labels, nil),
		pending:   prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", labels, nil),
		total:     prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUs for partition", labels, nil),
	}
}

func (pc *PartitionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pc.allocated
	ch <- pc.idle
	ch <- pc.other
	ch <- pc.pending
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
		if pm[p].total > 0 {
			ch <- prometheus.MustNewConstMetric(pc.total, prometheus.GaugeValue, pm[p].total, p)
		}
	}
}
