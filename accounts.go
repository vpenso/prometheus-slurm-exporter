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
        "strconv"
        "regexp"
        "github.com/prometheus/client_golang/prometheus"
)

func AccountsData() []byte {
        cmd := exec.Command("squeue","-a","-r","-h","-o %A|%a|%T|%C")
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
        pending_cpus float64
        running float64
        running_cpus float64
        suspended float64
}

func ParseAccountsMetrics(input []byte) map[string]*JobMetrics {
        accounts := make(map[string]*JobMetrics)
        lines := strings.Split(string(input), "\n")
        for _, line := range lines {
                if strings.Contains(line,"|") {
                        account := strings.Split(line,"|")[1]
                        _,key := accounts[account]
                        if !key {
                                accounts[account] = &JobMetrics{0,0,0,0,0}
                        }
                        state := strings.Split(line,"|")[2]
                        state = strings.ToLower(state)
                        cpus,_ := strconv.ParseFloat(strings.Split(line,"|")[3],64)
                        pending := regexp.MustCompile(`^pending`)
                        running := regexp.MustCompile(`^running`)
                        suspended := regexp.MustCompile(`^suspended`)
                        switch {
                        case pending.MatchString(state) == true:
                                accounts[account].pending++
                                accounts[account].pending_cpus += cpus
                        case running.MatchString(state) == true:
                                accounts[account].running++
                                accounts[account].running_cpus += cpus
                        case suspended.MatchString(state) == true:
                                accounts[account].suspended++
                        }
                }
        }
        return accounts
}

type AccountsCollector struct {
        pending *prometheus.Desc
        pending_cpus *prometheus.Desc
        running *prometheus.Desc
        running_cpus *prometheus.Desc
        suspended *prometheus.Desc
}

func NewAccountsCollector() *AccountsCollector {
        labels := []string{"account"}
        return &AccountsCollector{
                pending: prometheus.NewDesc("slurm_account_jobs_pending", "Pending jobs for account", labels, nil),
                pending_cpus: prometheus.NewDesc("slurm_account_cpus_pending", "Pending jobs for account", labels, nil),
                running: prometheus.NewDesc("slurm_account_jobs_running", "Running jobs for account", labels, nil),
                running_cpus: prometheus.NewDesc("slurm_account_cpus_running", "Running cpus for account", labels, nil),
                suspended: prometheus.NewDesc("slurm_account_jobs_suspended", "Suspended jobs for account", labels, nil),
        }
}

func (ac *AccountsCollector) Describe(ch chan<- *prometheus.Desc) {
        ch <- ac.pending
        ch <- ac.pending_cpus
        ch <- ac.running
        ch <- ac.running_cpus
        ch <- ac.suspended
}

func (ac *AccountsCollector) Collect(ch chan<- prometheus.Metric) {
        am := ParseAccountsMetrics(AccountsData())
        for a := range am {
                if am[a].pending > 0 {
                        ch <- prometheus.MustNewConstMetric(ac.pending, prometheus.GaugeValue, am[a].pending, a)
                }
                if am[a].pending_cpus > 0 {
                        ch <- prometheus.MustNewConstMetric(ac.pending_cpus, prometheus.GaugeValue, am[a].pending_cpus, a)
                }
                if am[a].running > 0 {
                        ch <- prometheus.MustNewConstMetric(ac.running, prometheus.GaugeValue, am[a].running, a)
                }
                if am[a].running_cpus > 0 {
                        ch <- prometheus.MustNewConstMetric(ac.running_cpus, prometheus.GaugeValue, am[a].running_cpus, a)
                }
                if am[a].suspended > 0 {
                        ch <- prometheus.MustNewConstMetric(ac.suspended, prometheus.GaugeValue, am[a].suspended, a)
                }
        }
}
