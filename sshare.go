/* Copyright 2021 Victor Penso

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

func FairShareData() []byte {
	cmd := exec.Command("sshare", "-n", "-P", "-o", "account,fairshare")
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

type FairShareMetrics struct {
	fairshare float64
}

func ParseFairShareMetrics() map[string]*FairShareMetrics {
	accounts := make(map[string]*FairShareMetrics)
	lines := strings.Split(string(FairShareData()), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "  ") {
			if strings.Contains(line, "|") {
				account := strings.Trim(strings.Split(line, "|")[0], " ")
				_, key := accounts[account]
				if !key {
					accounts[account] = &FairShareMetrics{0}
				}
				fairshare, _ := strconv.ParseFloat(strings.Split(line, "|")[1], 64)
				accounts[account].fairshare = fairshare
			}
		}
	}
	return accounts
}

type FairShareCollector struct {
	fairshare *prometheus.Desc
}

func NewFairShareCollector() *FairShareCollector {
	labels := []string{"account"}
	return &FairShareCollector{
		fairshare: prometheus.NewDesc("slurm_account_fairshare", "FairShare for account", labels, nil),
	}
}

func (fsc *FairShareCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- fsc.fairshare
}

func (fsc *FairShareCollector) Collect(ch chan<- prometheus.Metric) {
	fsm := ParseFairShareMetrics()
	for f := range fsm {
		ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, fsm[f].fairshare, f)
	}
}
