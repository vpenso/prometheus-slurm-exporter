/* Copyright 2021 Chris Read

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
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// NodeMetrics stores metrics for each node
type NodeMetrics struct {
	memAlloc uint64
	memTotal uint64
	cpuAlloc uint64
	cpuIdle  uint64
	cpuOther uint64
	cpuTotal uint64
	nodeStatus string
}

func NodeGetMetrics() map[string]*NodeMetrics {
	return ParseNodeMetrics(NodeData())
}

// ParseNodeMetrics takes the output of sinfo with node data
// It returns a map of metrics per node
func ParseNodeMetrics(input []byte) map[string]*NodeMetrics {
	nodes := make(map[string]*NodeMetrics)
	lines := strings.Split(string(input), "\n")

	// Sort and remove all the duplicates from the 'sinfo' output
	sort.Strings(lines)
	linesUniq := RemoveDuplicates(lines)

	for _, line := range linesUniq {
		node := strings.Fields(line)
		
		if len(node) == 0 { 
			continue 
		}

		nodeName := node[0]
		
		nodes[nodeName] = &NodeMetrics{0, 0, 0, 0, 0, 0, ""}

		if len(node) >= 5 {
			nodeStatus := node[4]
		  nodes[nodeName].nodeStatus = nodeStatus
		} else {
			nodeStatus := "unknown"
			nodes[nodeName].nodeStatus = nodeStatus
		}

		if len(node) >= 3 {
 		  memAlloc, _ := strconv.ParseUint(node[1], 10, 64)
		  memTotal, _ := strconv.ParseUint(node[2], 10, 64)
			nodes[nodeName].memAlloc = memAlloc
			nodes[nodeName].memTotal = memTotal
		} else {
			memAlloc := uint64(0)
			memTotal := uint64(0)
			nodes[nodeName].memAlloc = memAlloc
			nodes[nodeName].memTotal = memTotal
		}

		if len(node) >= 4 {
		  cpuInfo := strings.Split(node[3], "/")
			if len(cpuInfo) >= 1 {
		    cpuAlloc, _ := strconv.ParseUint(cpuInfo[0], 10, 64)
			  nodes[nodeName].cpuAlloc = cpuAlloc
		  } else {
				cpuAlloc := uint64(0)
			  nodes[nodeName].cpuAlloc = cpuAlloc
			}
			if len(cpuInfo) >= 2 {
  		  cpuIdle, _ := strconv.ParseUint(cpuInfo[1], 10, 64)
			  nodes[nodeName].cpuIdle = cpuIdle
			} else {
  		  cpuIdle := uint64(0)
			  nodes[nodeName].cpuIdle = cpuIdle
			}
			if len(cpuInfo) >= 3 {
		    cpuOther, _ := strconv.ParseUint(cpuInfo[2], 10, 64)
	  	  nodes[nodeName].cpuOther = cpuOther
  		} else {
  		  cpuOther := uint64(0)
		    nodes[nodeName].cpuOther = cpuOther
		  }
			if len(cpuInfo) >= 4 {
		    cpuTotal, _ := strconv.ParseUint(cpuInfo[3], 10, 64)
	  	  nodes[nodeName].cpuTotal = cpuTotal
  		} else {
  		  cpuTotal := uint64(0)
	  	  nodes[nodeName].cpuTotal = cpuTotal
		  }
	  } else {
			cpuAlloc := uint64(0)
			cpuIdle := uint64(0)
			cpuOther := uint64(0)
			cpuTotal := uint64(0)
			nodes[nodeName].cpuAlloc = cpuAlloc
			nodes[nodeName].cpuIdle = cpuIdle
		  nodes[nodeName].cpuOther = cpuOther
		  nodes[nodeName].cpuTotal = cpuTotal
	  }
	}

	return nodes
}

// NodeData executes the sinfo command to get data for each node
// It returns the output of the sinfo command
func NodeData() []byte {
	cmd := exec.Command("sinfo", "-h", "-N", "-O", "NodeList,AllocMem,Memory,CPUsState,StateLong")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

type NodeCollector struct {
	cpuAlloc *prometheus.Desc
	cpuIdle  *prometheus.Desc
	cpuOther *prometheus.Desc
	cpuTotal *prometheus.Desc
	memAlloc *prometheus.Desc
	memTotal *prometheus.Desc
}

// NewNodeCollector creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeCollector() *NodeCollector {
	labels := []string{"node","status"}

	return &NodeCollector{
		cpuAlloc: prometheus.NewDesc("slurm_node_cpu_alloc", "Allocated CPUs per node", labels, nil),
		cpuIdle:  prometheus.NewDesc("slurm_node_cpu_idle", "Idle CPUs per node", labels, nil),
		cpuOther: prometheus.NewDesc("slurm_node_cpu_other", "Other CPUs per node", labels, nil),
		cpuTotal: prometheus.NewDesc("slurm_node_cpu_total", "Total CPUs per node", labels, nil),
		memAlloc: prometheus.NewDesc("slurm_node_mem_alloc", "Allocated memory per node", labels, nil),
		memTotal: prometheus.NewDesc("slurm_node_mem_total", "Total memory per node", labels, nil),
	}
}

// Send all metric descriptions
func (nc *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.cpuAlloc
	ch <- nc.cpuIdle
	ch <- nc.cpuOther
	ch <- nc.cpuTotal
	ch <- nc.memAlloc
	ch <- nc.memTotal
}

func (nc *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	nodes := NodeGetMetrics()
	for node := range nodes {
		ch <- prometheus.MustNewConstMetric(nc.cpuAlloc, prometheus.GaugeValue, float64(nodes[node].cpuAlloc), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuIdle,  prometheus.GaugeValue, float64(nodes[node].cpuIdle),  node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuOther, prometheus.GaugeValue, float64(nodes[node].cpuOther), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuTotal, prometheus.GaugeValue, float64(nodes[node].cpuTotal), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memAlloc, prometheus.GaugeValue, float64(nodes[node].memAlloc), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memTotal, prometheus.GaugeValue, float64(nodes[node].memTotal), node, nodes[node].nodeStatus)
	}
}
