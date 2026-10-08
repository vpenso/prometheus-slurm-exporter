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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

// NodeGpuMetrics holds per-GPU-type counts for one node.
type NodeGpuMetrics struct {
	gpuType string
	total   uint64
	alloc   uint64
}

// NodeMetrics stores metrics for each node
type NodeMetrics struct {
	memAlloc   uint64
	memTotal   uint64
	cpuAlloc   uint64
	cpuIdle    uint64
	cpuOther   uint64
	cpuTotal   uint64
	nodeStatus string
	partitions string
	gpus       []NodeGpuMetrics
}

func NodeGetMetrics() map[string]*NodeMetrics {
	sinfo, err := SinfoData()
	if err != nil {
		log.Errorf("node: %v", err)
		return map[string]*NodeMetrics{}
	}
	return ParseNodeMetrics(sinfo)
}

// ParseNodeMetrics takes the parsed `sinfo --json` response and returns a
// map of metrics per node. Unlike the old text-based sinfo output, --json
// already returns exactly one entry per node, so no dedup pass is needed.
func ParseNodeMetrics(sinfo *slurmcli.SinfoResponse) map[string]*NodeMetrics {
	nodes := make(map[string]*NodeMetrics)
	for _, n := range sinfo.Nodes {
		nodes[n.Name] = &NodeMetrics{
			memAlloc:   uint64(n.AllocMemory),
			memTotal:   uint64(n.RealMemory),
			cpuAlloc:   uint64(n.AllocCPUs),
			cpuIdle:    uint64(n.IdleCPUs),
			cpuOther:   uint64(n.OtherCPUs()),
			cpuTotal:   uint64(n.CPUs),
			nodeStatus: strings.ToLower(strings.Join(n.State, "+")),
			partitions: strings.Join(n.Partitions, ","),
			gpus:       parseNodeGpus(n.Gres, n.GresUsed),
		}
	}
	return nodes
}

// parseNodeGpus pairs the total and used GRES descriptor strings of a node
// into per-type GPU counts, e.g. gres "gpu:mi250:8" with gres_used
// "gpu:mi250:3(IDX:0-2)" yields one {mi250, 8, 3} entry.
func parseNodeGpus(gres, gresUsed string) []NodeGpuMetrics {
	var gpus []NodeGpuMetrics
	for _, total := range slurmcli.ParseGresString(gres) {
		if total.Kind != "gpu" {
			continue
		}
		var alloc int64
		for _, used := range slurmcli.ParseGresString(gresUsed) {
			if used.Kind == total.Kind && used.Type == total.Type {
				alloc = used.Count
				break
			}
		}
		gpus = append(gpus, NodeGpuMetrics{gpuType: total.Type, total: uint64(total.Count), alloc: uint64(alloc)})
	}
	return gpus
}

type NodeCollector struct {
	cpuAlloc *prometheus.Desc
	cpuIdle  *prometheus.Desc
	cpuOther *prometheus.Desc
	cpuTotal *prometheus.Desc
	memAlloc *prometheus.Desc
	memTotal *prometheus.Desc
	gpuAlloc *prometheus.Desc
	gpuTotal *prometheus.Desc
}

// NewNodeCollector creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeCollector() *NodeCollector {
	labels := []string{"node", "status", "partition"}
	gpuLabels := []string{"node", "status", "gpu_type"}

	return &NodeCollector{
		cpuAlloc: prometheus.NewDesc("slurm_node_cpu_alloc", "Allocated CPUs per node", labels, nil),
		cpuIdle:  prometheus.NewDesc("slurm_node_cpu_idle", "Idle CPUs per node", labels, nil),
		cpuOther: prometheus.NewDesc("slurm_node_cpu_other", "Other CPUs per node", labels, nil),
		cpuTotal: prometheus.NewDesc("slurm_node_cpu_total", "Total CPUs per node", labels, nil),
		memAlloc: prometheus.NewDesc("slurm_node_mem_alloc", "Allocated memory per node", labels, nil),
		memTotal: prometheus.NewDesc("slurm_node_mem_total", "Total memory per node", labels, nil),
		gpuAlloc: prometheus.NewDesc("slurm_node_gpu_alloc", "Allocated GPUs per node", gpuLabels, nil),
		gpuTotal: prometheus.NewDesc("slurm_node_gpu_total", "Total GPUs per node", gpuLabels, nil),
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
	ch <- nc.gpuAlloc
	ch <- nc.gpuTotal
}

func (nc *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	nodes := NodeGetMetrics()
	for node := range nodes {
		ch <- prometheus.MustNewConstMetric(nc.cpuAlloc, prometheus.GaugeValue, float64(nodes[node].cpuAlloc), node, nodes[node].nodeStatus, nodes[node].partitions)
		ch <- prometheus.MustNewConstMetric(nc.cpuIdle, prometheus.GaugeValue, float64(nodes[node].cpuIdle), node, nodes[node].nodeStatus, nodes[node].partitions)
		ch <- prometheus.MustNewConstMetric(nc.cpuOther, prometheus.GaugeValue, float64(nodes[node].cpuOther), node, nodes[node].nodeStatus, nodes[node].partitions)
		ch <- prometheus.MustNewConstMetric(nc.cpuTotal, prometheus.GaugeValue, float64(nodes[node].cpuTotal), node, nodes[node].nodeStatus, nodes[node].partitions)
		ch <- prometheus.MustNewConstMetric(nc.memAlloc, prometheus.GaugeValue, float64(nodes[node].memAlloc), node, nodes[node].nodeStatus, nodes[node].partitions)
		ch <- prometheus.MustNewConstMetric(nc.memTotal, prometheus.GaugeValue, float64(nodes[node].memTotal), node, nodes[node].nodeStatus, nodes[node].partitions)
		for _, gpu := range nodes[node].gpus {
			ch <- prometheus.MustNewConstMetric(nc.gpuAlloc, prometheus.GaugeValue, float64(gpu.alloc), node, nodes[node].nodeStatus, gpu.gpuType)
			ch <- prometheus.MustNewConstMetric(nc.gpuTotal, prometheus.GaugeValue, float64(gpu.total), node, nodes[node].nodeStatus, gpu.gpuType)
		}
	}
}
