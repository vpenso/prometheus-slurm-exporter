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
        "os/exec"
        "log"
        "strings"
        "regexp"
        "github.com/prometheus/client_golang/prometheus"
)

func AccountsData() []byte {
        cmd := exec.Command("squeue", "-h", "-o '%A|%a|%T'")
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

type JobMetrics struct {
        pending float64
        running float64
}

func ParseAccountsMetrics(input []byte) map[string]*JobMetrics {
        accounts := make(map[string]*JobMetrics)
        lines := strings.Split(string(input), "\n")
        for _, line := range lines {
                if strings.Contains(line,"|") {
                        account := strings.Split(line,"|")[1]
                        _,key := accounts[account]
                        if !key {
                                accounts[account] = &JobMetrics{0,0}
                        }
                        state := strings.Split(line,"|")[2]
                        state = strings.ToLower(state)
                        running := regexp.MustCompile(`^running`)
                        pending := regexp.MustCompile(`^pending`)
                        switch {
                        case running.MatchString(state) == true:
                                accounts[account].running++
                        case pending.MatchString(state) == true:
                                accounts[account].pending++
                        }
                }
        }
        return accounts
}

type AccountsCollector struct {
        running *prometheus.Desc
        pending *prometheus.Desc
}

func NewAccountsCollector() *AccountsCollector {
        labels := []string{"account"}
        return &AccountsCollector{
                running: prometheus.NewDesc("slurm_accounts_jobs_running", "Running jobs for account", labels, nil),
                pending: prometheus.NewDesc("slurm_accounts_jobs_pending", "Running jobs for account", labels, nil),
        }
}

func (ac *AccountsCollector) Describe(ch chan<- *prometheus.Desc) {
        ch <- ac.running
        ch <- ac.pending
}

func (ac *AccountsCollector) Collect(ch chan<- prometheus.Metric) {
        am := ParseAccountsMetrics(AccountsData())
        for a := range am {
                log.Print(a)
                ch <- prometheus.MustNewConstMetric(ac.running, prometheus.GaugeValue, am[a].running, a)
                ch <- prometheus.MustNewConstMetric(ac.pending, prometheus.GaugeValue, am[a].pending, a)
        }
}
