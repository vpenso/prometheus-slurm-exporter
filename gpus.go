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

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

// GPUsMetrics holds per-GPU-type (e.g. "mi250", "a100", or "" for untyped
// legacy GRES) allocation counts. Vendor is not special-cased here: AMD and
// NVIDIA GPUs are both just typed Slurm GRES entries, so the same code path
// handles both.
type GPUsMetrics struct {
	alloc       map[string]float64
	idle        map[string]float64
	total       map[string]float64
	utilization map[string]float64
}

func GPUsGetMetrics() *GPUsMetrics {
	sinfo, err := SinfoData()
	if err != nil {
		log.Errorf("gpus: %v", err)
		return emptyGPUsMetrics()
	}
	sacct, err := SacctData()
	if err != nil {
		log.Errorf("gpus: %v", err)
		return emptyGPUsMetrics()
	}
	return ParseGPUsMetrics(sinfo, sacct)
}

func emptyGPUsMetrics() *GPUsMetrics {
	return &GPUsMetrics{
		alloc:       map[string]float64{},
		idle:        map[string]float64{},
		total:       map[string]float64{},
		utilization: map[string]float64{},
	}
}

// ParseTotalGPUs returns configured GPU counts per type, summed across all
// nodes' `gres` field.
func ParseTotalGPUs(sinfo *slurmcli.SinfoResponse) map[string]float64 {
	totals := map[string]float64{}
	for _, n := range sinfo.Nodes {
		for _, g := range slurmcli.ParseGresString(n.Gres) {
			if g.Kind != "gpu" {
				continue
			}
			totals[g.Type] += float64(g.Count)
		}
	}
	return totals
}

// ParseAllocatedGPUs returns allocated GPU counts per type, summed across
// every currently running job's allocated GRES tres entries.
func ParseAllocatedGPUs(sacct *slurmcli.SacctResponse) map[string]float64 {
	allocated := map[string]float64{}
	for _, j := range sacct.Jobs {
		for _, g := range j.AllocatedGres() {
			allocated[g.Type] += float64(g.Count)
		}
	}
	return allocated
}

func ParseGPUsMetrics(sinfo *slurmcli.SinfoResponse, sacct *slurmcli.SacctResponse) *GPUsMetrics {
	gm := emptyGPUsMetrics()
	gm.total = ParseTotalGPUs(sinfo)
	gm.alloc = ParseAllocatedGPUs(sacct)

	for gpuType, total := range gm.total {
		idle := total - gm.alloc[gpuType]
		if idle < 0 {
			log.Warnf("gpus: allocated GPU count for type %q exceeds total, clamping idle to 0", gpuType)
			idle = 0
		}
		gm.idle[gpuType] = idle
		if total > 0 {
			gm.utilization[gpuType] = gm.alloc[gpuType] / total
		}
	}
	return gm
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

func NewGPUsCollector() *GPUsCollector {
	labels := []string{"gpu_type"}
	return &GPUsCollector{
		alloc:       prometheus.NewDesc("slurm_gpus_alloc", "Allocated GPUs", labels, nil),
		idle:        prometheus.NewDesc("slurm_gpus_idle", "Idle GPUs", labels, nil),
		total:       prometheus.NewDesc("slurm_gpus_total", "Total GPUs", labels, nil),
		utilization: prometheus.NewDesc("slurm_gpus_utilization", "GPU allocation utilization (alloc/total)", labels, nil),
	}
}

type GPUsCollector struct {
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
}

// Send all metric descriptions
func (cc *GPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.total
	ch <- cc.utilization
}
func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) {
	cm := GPUsGetMetrics()
	for gpuType, total := range cm.total {
		ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc[gpuType], gpuType)
		ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle[gpuType], gpuType)
		ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, total, gpuType)
		if total > 0 {
			ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, cm.utilization[gpuType], gpuType)
		}
	}
}
