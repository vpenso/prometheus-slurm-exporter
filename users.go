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

func UsersData() []byte {
        cmd := exec.Command("squeue","-a","-r","-h","-o %A|%u|%T|%C|%N")
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

type UserJobMetrics struct {
        pending float64
        running float64
        running_cpus float64
        suspended float64
	hosts map[string]int
}

func ParseUsersMetrics(input []byte) map[string]*UserJobMetrics {
        users := make(map[string]*UserJobMetrics)
        lines := strings.Split(string(input), "\n")
        for _, line := range lines {
                if strings.Contains(line,"|") {
                        user := strings.Split(line,"|")[1]
                        _,key := users[user]
                        if !key {
                                users[user] = &UserJobMetrics{0,0,0,0,map[string]int{}}
                        }
                        state := strings.Split(line,"|")[2]
                        state = strings.ToLower(state)
                        cpus,_ := strconv.ParseFloat(strings.Split(line,"|")[3],64)
                        pending := regexp.MustCompile(`^pending`)
                        running := regexp.MustCompile(`^running`)
                        suspended := regexp.MustCompile(`^suspended`)
                        hosts := strings.Split(strings.Split(line,"|")[4], ",")
                        switch {
                        case pending.MatchString(state) == true:
                                users[user].pending++
                        case running.MatchString(state) == true:
                                users[user].running++
                                users[user].running_cpus += cpus
				for _, host := range hosts {
					users[user].hosts[host] += 1
				}
                        case suspended.MatchString(state) == true:
                                users[user].suspended++
                        }
                }
        }
        return users
}

type UsersCollector struct {
        pending *prometheus.Desc
        running *prometheus.Desc
        running_cpus *prometheus.Desc
        suspended *prometheus.Desc
        hosts *prometheus.Desc
}

func NewUsersCollector() *UsersCollector {
        labels := []string{"user"}
        runningLabels := []string{"user","host"}
        return &UsersCollector {
                pending: prometheus.NewDesc("slurm_user_jobs_pending", "Pending jobs for user", labels, nil),
                running: prometheus.NewDesc("slurm_user_jobs_running", "Running jobs for user", labels, nil),
                running_cpus: prometheus.NewDesc("slurm_user_cpus_running", "Running cpus for user", labels, nil),
                suspended: prometheus.NewDesc("slurm_user_jobs_suspended", "Suspended jobs for user", labels, nil),
                hosts: prometheus.NewDesc("slurm_user_jobs_running_host", "Running jobs for user on host", runningLabels, nil),
        }
}

func (uc *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
        ch <- uc.pending
        ch <- uc.running
        ch <- uc.running_cpus
        ch <- uc.suspended
        ch <- uc.hosts
}

func (uc *UsersCollector) Collect(ch chan<- prometheus.Metric) {
        um := ParseUsersMetrics(UsersData())
        for u := range um {
                if um[u].pending > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.pending, prometheus.GaugeValue, um[u].pending, u)
                }
                if um[u].running > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.running, prometheus.GaugeValue, um[u].running, u)
                        for h, cnt := range um[u].hosts {
				ch <- prometheus.MustNewConstMetric(uc.hosts, prometheus.GaugeValue, float64(cnt), u, h)
			}
                }
                if um[u].running_cpus > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.running_cpus, prometheus.GaugeValue, um[u].running_cpus, u)
                }
                if um[u].suspended > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.suspended, prometheus.GaugeValue, um[u].suspended, u)
                }
        }
}

