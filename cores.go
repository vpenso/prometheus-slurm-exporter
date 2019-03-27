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
        "fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type CoresMetrics struct {
	alloc float64
	idle  float64
	other float64
	total float64
}

func CoresGetMetrics() *CoresMetrics {
	return ParseCoresMetrics(CoresData())
}

func ParseCoresMetrics(input []byte) *CoresMetrics {
	var cm CoresMetrics
	if strings.Contains(string(input), "/") {
	    splitted := strings.Split(strings.TrimSpace(string(input)), "/")
	    cm.alloc, _ = strconv.ParseFloat(splitted[0], 64)
	    cm.idle, _  = strconv.ParseFloat(splitted[1], 64)
	    cm.other, _ = strconv.ParseFloat(splitted[2], 64)
	    cm.total, _ = strconv.ParseFloat(splitted[3], 64)
	} 
	return &cm
}

// Execute the sinfo command and return its output
func CoresData() []byte {
	cmd := exec.Command("sinfo", "-h", "-o %C")
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

func NewCoresCollector() *CoresCollector {
	return &CoresCollector{
		alloc: prometheus.NewDesc("slurm_cores_alloc", "Allocated cores", nil, nil),
		idle:  prometheus.NewDesc("slurm_cores_idle", "Idle cores", nil, nil),
		other: prometheus.NewDesc("slurm_cores_other", "Mix cores", nil, nil),
		total: prometheus.NewDesc("slurm_cores_total", "Total cores", nil, nil),
	}
}

type CoresCollector struct {
	alloc *prometheus.Desc
	idle  *prometheus.Desc
	other *prometheus.Desc
	total *prometheus.Desc
}

// Send all metric descriptions
func (cc *CoresCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
}
func (cc *CoresCollector) Collect(ch chan<- prometheus.Metric) {
	cm := CoresGetMetrics()
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
}
