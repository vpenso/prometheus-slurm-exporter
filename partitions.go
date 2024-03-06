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
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
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
	allocated     float64
	idle          float64
	other         float64
	pending       float64
	total         float64
	totalGPUs     float64
	allocatedGPUs float64
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
				partitions[partition] = &PartitionMetrics{0, 0, 0, 0, 0, 0, 0}
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

	totalGPUs := TotalPartitionGPUsData()
	allocatedGPUs := AllocatedPartitionGPUsData()

	// Iterate over partitions to update GPU data
	for partition, metrics := range partitions {
		if gpuTotal, ok := totalGPUs[partition]; ok {
			metrics.totalGPUs = gpuTotal
		}
		if gpuAllocated, ok := allocatedGPUs[partition]; ok {
			metrics.allocatedGPUs = gpuAllocated
		}
	}

	return partitions
}

func TotalPartitionGPUsData() map[string]float64 {
	cmd := exec.Command("bash", "-c", `sinfo -h -o "%R,%G,%n" | grep "gpu:" | awk -F':' '{print $1,$3}' | awk -F'[(,]' '{print $1,$2}' | awk '{print $1, $3}' | awk '{count[$1]+=$2} END {for (partition in count) print partition, count[partition]}'`)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	totalGPUs := make(map[string]float64)
	lines := strings.Split(string(stdout), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			partition := parts[0]
			count, err := strconv.ParseFloat(parts[1], 64)
			if err == nil {
				totalGPUs[partition] = count
			}
		}
	}
	return totalGPUs
}

func AllocatedPartitionGPUsData() map[string]float64 {
	cmd := exec.Command("bash", "-c", `squeue -h -t R -o "%P %b" | grep "gres:gpu:" | awk -F'[: ]' '{print $1, $4}' | awk '{count[$1]+=$2} END {for (partition in count) print partition, count[partition]}' | sort`)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	allocatedGPUs := make(map[string]float64)
	lines := strings.Split(string(stdout), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			partition := parts[0]
			count, err := strconv.ParseFloat(parts[1], 64)
			if err == nil {
				allocatedGPUs[partition] = count
			}
		}
	}
	return allocatedGPUs
}

type PartitionsCollector struct {
	allocated     *prometheus.Desc
	idle          *prometheus.Desc
	other         *prometheus.Desc
	pending       *prometheus.Desc
	total         *prometheus.Desc
	totalGPUs     *prometheus.Desc
	allocatedGPUs *prometheus.Desc
}

func NewPartitionsCollector() *PartitionsCollector {
	labels := []string{"partition"}
	return &PartitionsCollector{
		allocated:     prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUs for partition", labels, nil),
		idle:          prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUs for partition", labels, nil),
		other:         prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUs for partition", labels, nil),
		pending:       prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", labels, nil),
		total:         prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUs for partition", labels, nil),
		totalGPUs:     prometheus.NewDesc("slurm_partition_gpus_total", "Total GPUs available for partition", labels, nil), // New
		allocatedGPUs: prometheus.NewDesc("slurm_partition_gpus_allocated", "GPUs allocated for partition", labels, nil),   // New
	}
}

func (pc *PartitionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pc.allocated
	ch <- pc.idle
	ch <- pc.other
	ch <- pc.pending
	ch <- pc.total
	ch <- pc.totalGPUs
	ch <- pc.allocatedGPUs
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
		if pm[p].totalGPUs > 0 {
			ch <- prometheus.MustNewConstMetric(pc.totalGPUs, prometheus.GaugeValue, pm[p].totalGPUs, p)
		}
		if pm[p].allocatedGPUs > 0 {
			ch <- prometheus.MustNewConstMetric(pc.allocatedGPUs, prometheus.GaugeValue, pm[p].allocatedGPUs, p)
		}
	}
}
