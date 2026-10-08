/* Copyright 2026 prometheus-slurm-exporter contributors

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
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/vpenso/prometheus-slurm-exporter/internal/smicli"
)

// ROCmCollector exposes node-local AMD/ROCm GPU device telemetry
// (utilization, memory, temperature, power) read directly from amd-smi (or
// rocm-smi). This is deliberately separate from the GRES-based slurm_gpus_*
// metrics in gpus.go: those describe job-scheduling occupancy as tracked by
// Slurm, this describes actual hardware state on the node it runs on, with
// no job-ID correlation - it is meant to run as a node-level exporter
// instance (e.g. one per GPU compute node), not the cluster-wide instance.
type ROCmCollector struct {
	smiCmd      string
	utilization *prometheus.Desc
	memUsed     *prometheus.Desc
	memTotal    *prometheus.Desc
	temperature *prometheus.Desc
	power       *prometheus.Desc
}

// NewROCmCollector creates a collector that invokes smiCmd (e.g. "amd-smi"
// or "rocm-smi") to gather ROCm GPU telemetry.
func NewROCmCollector(smiCmd string) *ROCmCollector {
	labels := []string{"node", "gpu"}
	tempLabels := []string{"node", "gpu", "sensor"}
	return &ROCmCollector{
		smiCmd:      smiCmd,
		utilization: prometheus.NewDesc("slurm_rocm_gpu_utilization_percent", "AMD ROCm GPU utilization percent", labels, nil),
		memUsed:     prometheus.NewDesc("slurm_rocm_gpu_memory_used_bytes", "AMD ROCm GPU memory used, in bytes", labels, nil),
		memTotal:    prometheus.NewDesc("slurm_rocm_gpu_memory_total_bytes", "AMD ROCm GPU memory total, in bytes", labels, nil),
		temperature: prometheus.NewDesc("slurm_rocm_gpu_temperature_celsius", "AMD ROCm GPU temperature, in Celsius", tempLabels, nil),
		power:       prometheus.NewDesc("slurm_rocm_gpu_power_watts", "AMD ROCm GPU socket power draw, in watts", labels, nil),
	}
}

func (rc *ROCmCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.utilization
	ch <- rc.memUsed
	ch <- rc.memTotal
	ch <- rc.temperature
	ch <- rc.power
}

func (rc *ROCmCollector) Collect(ch chan<- prometheus.Metric) {
	var devices []smicli.ROCmDevice
	if err := smicli.RunJSON(rc.smiCmd, []string{"metric", "--json"}, &devices); err != nil {
		// A node without GPUs, or without amd-smi/rocm-smi installed, must
		// not take down the whole exporter process - log and emit nothing.
		log.Errorf("rocm: %v", err)
		return
	}

	node, err := os.Hostname()
	if err != nil {
		node = "unknown"
	}

	for _, d := range devices {
		gpu := strconv.Itoa(d.GPU)
		ch <- prometheus.MustNewConstMetric(rc.utilization, prometheus.GaugeValue, d.Usage.GfxActivity, node, gpu)
		ch <- prometheus.MustNewConstMetric(rc.memUsed, prometheus.GaugeValue, float64(d.Mem.Used), node, gpu)
		ch <- prometheus.MustNewConstMetric(rc.memTotal, prometheus.GaugeValue, float64(d.Mem.Total), node, gpu)
		ch <- prometheus.MustNewConstMetric(rc.temperature, prometheus.GaugeValue, d.Temperature.Edge, node, gpu, "edge")
		ch <- prometheus.MustNewConstMetric(rc.temperature, prometheus.GaugeValue, d.Temperature.Junction, node, gpu, "junction")
		ch <- prometheus.MustNewConstMetric(rc.temperature, prometheus.GaugeValue, d.Temperature.Memory, node, gpu, "memory")
		ch <- prometheus.MustNewConstMetric(rc.power, prometheus.GaugeValue, d.Power.Socket, node, gpu)
	}
}
