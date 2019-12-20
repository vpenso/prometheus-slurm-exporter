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
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
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

func NodesGetMetrics() map[string]*NodesMetrics {
	return ParseNodesMetrics(NodesData())
}

func RemoveDuplicates(s []string) []string {
	m := map[string]bool{}
	t := []string{}

	// Walk through the slice 's' and for each value we haven't seen so far, append it to 't'.
	for _, v := range s {
		if _, seen := m[v]; !seen {
			t = append(t, v)
			m[v] = true
		}
	}

	return t
}

/*
 * Get the names of the partitions so that we can register the relevant collectors.
 */
func GetPartitionNames() []string {
	return ParsePartitionNames(GetSinfo())
}

// Execute the sinfo command and return its output.
func GetSinfo() []byte {
	cmd := exec.Command("sinfo", "-s")
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

func ParsePartitionNames(input []byte) []string {
	partitions := make([]string, 0)
	lines := strings.Split(string(input), "\n")
	for _, line := range lines[1:] {
		if len(line) > 2 {
			partitions = append(partitions, strings.Trim(strings.Fields(line)[0], "*"))
		}
	}
	return partitions
}

func ParseNodesMetrics(input []byte) map[string]*NodesMetrics {
	qnm := make(map[string]*NodesMetrics)
	// initialize the metrics
	for _, part := range GetPartitionNames() {
		qnm[part] = &NodesMetrics{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}

	lines := strings.Split(string(input), "\n")
	// Sort and remove all the duplicates from the 'sinfo' output
	sort.Strings(lines)
	lines_uniq := RemoveDuplicates(lines)

	for _, line := range lines_uniq {
		if strings.Contains(line, ",") {
			state := strings.Split(line, ",")[1]
			part := strings.Split(line, ",")[2]
			alloc := regexp.MustCompile(`^alloc`)
			comp := regexp.MustCompile(`^comp`)
			down := regexp.MustCompile(`^down`)
			drain := regexp.MustCompile(`^drain`)
			fail := regexp.MustCompile(`^fail`)
			err := regexp.MustCompile(`^err`)
			idle := regexp.MustCompile(`^idle`)
			maint := regexp.MustCompile(`^maint`)
			mix := regexp.MustCompile(`^mix`)
			resv := regexp.MustCompile(`^res`)
			switch {
			case alloc.MatchString(state) == true:
				qnm[part].alloc++
			case comp.MatchString(state) == true:
				qnm[part].comp++
			case down.MatchString(state) == true:
				qnm[part].down++
			case drain.MatchString(state) == true:
				qnm[part].drain++
			case fail.MatchString(state) == true:
				qnm[part].fail++
			case err.MatchString(state) == true:
				qnm[part].err++
			case idle.MatchString(state) == true:
				qnm[part].idle++
			case maint.MatchString(state) == true:
				qnm[part].maint++
			case mix.MatchString(state) == true:
				qnm[part].mix++
			case resv.MatchString(state) == true:
				qnm[part].resv++
			}
		}
	}
	return qnm
}

// Execute the sinfo command and return its output
func NodesData() []byte {
	cmd := exec.Command("sinfo", "-h", "-o %n,%T,%R")
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

func NewNodesCollector() *NodesCollector {
	labels := []string{"partition"}
	return &NodesCollector{
		alloc: prometheus.NewDesc("slurm_nodes_alloc", "Allocated nodes", labels, nil),
		comp:  prometheus.NewDesc("slurm_nodes_comp", "Completing nodes", labels, nil),
		down:  prometheus.NewDesc("slurm_nodes_down", "Down nodes", labels, nil),
		drain: prometheus.NewDesc("slurm_nodes_drain", "Drain nodes", labels, nil),
		err:   prometheus.NewDesc("slurm_nodes_err", "Error nodes", labels, nil),
		fail:  prometheus.NewDesc("slurm_nodes_fail", "Fail nodes", labels, nil),
		idle:  prometheus.NewDesc("slurm_nodes_idle", "Idle nodes", labels, nil),
		maint: prometheus.NewDesc("slurm_nodes_maint", "Maint nodes", labels, nil),
		mix:   prometheus.NewDesc("slurm_nodes_mix", "Mix nodes", labels, nil),
		resv:  prometheus.NewDesc("slurm_nodes_resv", "Reserved nodes", labels, nil),
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
	for k := range nm {
		ch <- prometheus.MustNewConstMetric(nc.alloc, prometheus.GaugeValue, nm[k].alloc, k)
		ch <- prometheus.MustNewConstMetric(nc.comp, prometheus.GaugeValue, nm[k].comp, k)
		ch <- prometheus.MustNewConstMetric(nc.down, prometheus.GaugeValue, nm[k].down, k)
		ch <- prometheus.MustNewConstMetric(nc.drain, prometheus.GaugeValue, nm[k].drain, k)
		ch <- prometheus.MustNewConstMetric(nc.err, prometheus.GaugeValue, nm[k].err, k)
		ch <- prometheus.MustNewConstMetric(nc.fail, prometheus.GaugeValue, nm[k].fail, k)
		ch <- prometheus.MustNewConstMetric(nc.idle, prometheus.GaugeValue, nm[k].idle, k)
		ch <- prometheus.MustNewConstMetric(nc.maint, prometheus.GaugeValue, nm[k].maint, k)
		ch <- prometheus.MustNewConstMetric(nc.mix, prometheus.GaugeValue, nm[k].mix, k)
		ch <- prometheus.MustNewConstMetric(nc.resv, prometheus.GaugeValue, nm[k].resv, k)
	}
}
