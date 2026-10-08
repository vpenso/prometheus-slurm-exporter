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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

type NodesMetrics struct {
	alloc float64
	comp  float64
	down  float64
	drain float64
	err   float64
	fail  float64
	idle  float64
	maint float64
	mix   float64
	resv  float64
}

func NodesGetMetrics() *NodesMetrics {
	sinfo, err := SinfoData()
	if err != nil {
		log.Errorf("nodes: %v", err)
		return &NodesMetrics{}
	}
	return ParseNodesMetrics(sinfo)
}

// classifyNodeState maps a node's --json state flag set (e.g.
// ["MIXED","DRAIN"]) to this collector's single legacy bucket, in the same
// priority order the old sinfo "%T"-prefix regexes used. Under the old text
// format sinfo already emitted one composite state token per node; the
// JSON state array can carry multiple flags, so ties are broken by this
// fixed priority list rather than double-counting a node into two buckets.
func classifyNodeState(flags []string) string {
	joined := strings.ToLower(strings.Join(flags, ","))
	switch {
	case strings.Contains(joined, "alloc"):
		return "alloc"
	case strings.Contains(joined, "comp"):
		return "comp"
	case strings.Contains(joined, "down"):
		return "down"
	case strings.Contains(joined, "drain"):
		return "drain"
	case strings.Contains(joined, "fail"):
		return "fail"
	case strings.Contains(joined, "err"):
		return "err"
	case strings.Contains(joined, "idle"):
		return "idle"
	case strings.Contains(joined, "maint"):
		return "maint"
	case strings.Contains(joined, "mix"):
		return "mix"
	case strings.Contains(joined, "resv"), strings.Contains(joined, "reserved"):
		return "resv"
	default:
		return ""
	}
}

func ParseNodesMetrics(sinfo *slurmcli.SinfoResponse) *NodesMetrics {
	var nm NodesMetrics
	for _, n := range sinfo.Nodes {
		switch classifyNodeState(n.State) {
		case "alloc":
			nm.alloc++
		case "comp":
			nm.comp++
		case "down":
			nm.down++
		case "drain":
			nm.drain++
		case "err":
			nm.err++
		case "fail":
			nm.fail++
		case "idle":
			nm.idle++
		case "maint":
			nm.maint++
		case "mix":
			nm.mix++
		case "resv":
			nm.resv++
		}
	}
	return &nm
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

func NewNodesCollector() *NodesCollector {
	return &NodesCollector{
		alloc: prometheus.NewDesc("slurm_nodes_alloc", "Allocated nodes", nil, nil),
		comp:  prometheus.NewDesc("slurm_nodes_comp", "Completing nodes", nil, nil),
		down:  prometheus.NewDesc("slurm_nodes_down", "Down nodes", nil, nil),
		drain: prometheus.NewDesc("slurm_nodes_drain", "Drain nodes", nil, nil),
		err:   prometheus.NewDesc("slurm_nodes_err", "Error nodes", nil, nil),
		fail:  prometheus.NewDesc("slurm_nodes_fail", "Fail nodes", nil, nil),
		idle:  prometheus.NewDesc("slurm_nodes_idle", "Idle nodes", nil, nil),
		maint: prometheus.NewDesc("slurm_nodes_maint", "Maint nodes", nil, nil),
		mix:   prometheus.NewDesc("slurm_nodes_mix", "Mix nodes", nil, nil),
		resv:  prometheus.NewDesc("slurm_nodes_resv", "Reserved nodes", nil, nil),
	}
}

type NodesCollector struct {
	alloc *prometheus.Desc
	comp  *prometheus.Desc
	down  *prometheus.Desc
	drain *prometheus.Desc
	err   *prometheus.Desc
	fail  *prometheus.Desc
	idle  *prometheus.Desc
	maint *prometheus.Desc
	mix   *prometheus.Desc
	resv  *prometheus.Desc
}

// Send all metric descriptions
func (nc *NodesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.alloc
	ch <- nc.comp
	ch <- nc.down
	ch <- nc.drain
	ch <- nc.err
	ch <- nc.fail
	ch <- nc.idle
	ch <- nc.maint
	ch <- nc.mix
	ch <- nc.resv
}
func (nc *NodesCollector) Collect(ch chan<- prometheus.Metric) {
	nm := NodesGetMetrics()
	ch <- prometheus.MustNewConstMetric(nc.alloc, prometheus.GaugeValue, nm.alloc)
	ch <- prometheus.MustNewConstMetric(nc.comp, prometheus.GaugeValue, nm.comp)
	ch <- prometheus.MustNewConstMetric(nc.down, prometheus.GaugeValue, nm.down)
	ch <- prometheus.MustNewConstMetric(nc.drain, prometheus.GaugeValue, nm.drain)
	ch <- prometheus.MustNewConstMetric(nc.err, prometheus.GaugeValue, nm.err)
	ch <- prometheus.MustNewConstMetric(nc.fail, prometheus.GaugeValue, nm.fail)
	ch <- prometheus.MustNewConstMetric(nc.idle, prometheus.GaugeValue, nm.idle)
	ch <- prometheus.MustNewConstMetric(nc.maint, prometheus.GaugeValue, nm.maint)
	ch <- prometheus.MustNewConstMetric(nc.mix, prometheus.GaugeValue, nm.mix)
	ch <- prometheus.MustNewConstMetric(nc.resv, prometheus.GaugeValue, nm.resv)
}
