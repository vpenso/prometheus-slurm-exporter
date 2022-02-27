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
	"strconv"
	"regexp"
)

type GPUsMetrics struct {
	alloc       float64
	idle        float64
	other       float64
	total       float64
	utilization float64
}

func GPUsGetMetrics() *GPUsMetrics {
	return ParseGPUsMetrics()
}

/* TODO:
  sinfo has gresUSED since slurm>=19.05.0rc01 https://github.com/SchedMD/slurm/blob/master/NEWS
  revert to old process on slurm<19.05.0rc01
  --format=AllocGRES will return gres/gpu=8
  --format=AllocTRES will return billing=16,cpu=16,gres/gpu=8,mem=256G,node=1
func ParseAllocatedGPUs() float64 {
	var num_gpus = 0.0

	args := []string{"-a", "-X", "--format=Allocgres", "--state=RUNNING", "--noheader", "--parsable2"}
	output := string(Execute("sacct", args))
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			if len(line) > 0 {
				line = strings.Trim(line, "\"")
				descriptor := strings.TrimPrefix(line, "gpu:")
				job_gpus, _ := strconv.ParseFloat(descriptor, 64)
				num_gpus += job_gpus
			}
		}
	}

	return num_gpus
}
*/

func ParseAllocatedGPUs() float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format=Nodes,GresUsed:512 --state=allocated
	// 3                   gpu:2
	// 1                   gpu:(null):3(IDX:0-7)
	// 13                  gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)


	args := []string{"-a", "-h", "--Format=Nodes,GresUsed:512", "--state=allocated"}
	output := string(Execute("sinfo", args))
	re := regexp.MustCompile("gpu:([^:(]*):?([0-9]+)(\\([^)]*\\))?")
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_active_gpus := strings.Fields(line)[1]
				num_node_active_gpus := 0.0
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_active_gpus
			}
		}
	}

	return num_gpus
}

func ParseIdleGPUs() float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format=Nodes,Gres:512,GresUsed:512 --state=idle,allocated
	// 3                   gpu:4                                       gpu:2
	// 1                   gpu:8(S:0-1)                                gpu:(null):3(IDX:0-7)
	// 13                  gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1)          gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)


	args := []string{"-a", "-h", "--Format=Nodes,Gres:512,GresUsed:512", "--state=idle,allocated"}
	output := string(Execute("sinfo", args))
	re := regexp.MustCompile("gpu:([^:(]*):?([0-9]+)(\\([^)]*\\))?")
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_node_active_gpus := 0.0
				node_active_gpus := strings.Fields(line)[2]
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * (num_node_gpus - num_node_active_gpus)
			}
		}
	}

	return num_gpus
}

func ParseTotalGPUs() float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format=Nodes,Gres:512
	// 3                   gpu:4
	// 1                   gpu:8(S:0-1)
	// 13                  gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1)

	args := []string{"-a", "-h", "--Format=Nodes,Gres:512"}
	output := string(Execute("sinfo", args))
	re := regexp.MustCompile("gpu:([^:(]*):?([0-9]+)(\\([^)]*\\))?")
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_gpus
			}
		}
	}

	return num_gpus
}

func ParseGPUsMetrics() *GPUsMetrics {
	var gm GPUsMetrics
	total_gpus := ParseTotalGPUs()
	allocated_gpus := ParseAllocatedGPUs()
	idle_gpus := ParseIdleGPUs()
	other_gpus := total_gpus - allocated_gpus - idle_gpus
	gm.alloc = allocated_gpus
	gm.idle = idle_gpus
	gm.other = other_gpus
	gm.total = total_gpus
	gm.utilization = allocated_gpus / total_gpus
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
		other: prometheus.NewDesc("slurm_gpus_other", "Other GPUs", nil, nil),
		total: prometheus.NewDesc("slurm_gpus_total", "Total GPUs", nil, nil),
		utilization: prometheus.NewDesc("slurm_gpus_utilization", "Total GPU utilization", nil, nil),
	}
}

type GPUsCollector struct {
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	other       *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
}

// Send all metric descriptions
func (cc *GPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
	ch <- cc.utilization
}
func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) {
	cm := GPUsGetMetrics()
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, cm.utilization)
}
