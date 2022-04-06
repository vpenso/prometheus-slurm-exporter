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
        cmd := exec.Command("squeue","-a","-r","-h","-o %A|%u|%T|%C|%m")
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
        running_mems float64
        suspended float64
}

func prase_unit(input string) float64 {
        reg := `^(\d+)([KMGT])$`
        r := regexp.MustCompile(reg)
        matchs := r.FindStringSubmatch(  string(input) )
        num,_ := strconv.Atoi( matchs[1] )
        unit := matchs[2]
        ret := 0
        if ("K" == unit ) {
                ret = num * 1024
        } else if ("M" == unit ) {
                ret = num * 1024 * 1024
        } else if ("G" == unit ) {
                ret = num * 1024 * 1024 * 1024
        } else if ("T" == unit ) {
                ret = num * 1024 * 1024 * 1024 * 1024
        } else {
                ret = num
        }
        return float64(ret)
}

func ParseUsersMetrics(input []byte) map[string]*UserJobMetrics {
        users := make(map[string]*UserJobMetrics)
        lines := strings.Split(string(input), "\n")
        for _, line := range lines {
                if strings.Contains(line,"|") {
                        user := strings.Split(line,"|")[1]
                        _,key := users[user]
                        if !key {
                                users[user] = &UserJobMetrics{0,0,0,0,0}
                        }
                        state := strings.Split(line,"|")[2]
                        state = strings.ToLower(state)
                        cpus,_ := strconv.ParseFloat(strings.Split(line,"|")[3],64)
                        mems := prase_unit(strings.Split(line,"|")[4])
                        pending := regexp.MustCompile(`^pending`)
                        running := regexp.MustCompile(`^running`)
                        suspended := regexp.MustCompile(`^suspended`)
                        switch {
                        case pending.MatchString(state) == true:
                                users[user].pending++
                        case running.MatchString(state) == true:
                                users[user].running++
                                users[user].running_cpus += cpus
                                users[user].running_mems += mems
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
        running_mems *prometheus.Desc
        suspended *prometheus.Desc
}

func NewUsersCollector() *UsersCollector {
        labels := []string{"user"}
        return &UsersCollector {
                pending: prometheus.NewDesc("slurm_user_jobs_pending", "Pending jobs for user", labels, nil), 
                running: prometheus.NewDesc("slurm_user_jobs_running", "Running jobs for user", labels, nil),
                running_cpus: prometheus.NewDesc("slurm_user_cpus_running", "Running cpus for user", labels, nil),
                running_mems: prometheus.NewDesc("slurm_user_mems_running", "Running mems for user", labels, nil),
                suspended: prometheus.NewDesc("slurm_user_jobs_suspended", "Suspended jobs for user", labels, nil),
        }
}

func (uc *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
        ch <- uc.pending
        ch <- uc.running
        ch <- uc.running_cpus
        ch <- uc.running_mems
        ch <- uc.suspended
}

func (uc *UsersCollector) Collect(ch chan<- prometheus.Metric) {
        um := ParseUsersMetrics(UsersData())
        for u := range um {
                if um[u].pending > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.pending, prometheus.GaugeValue, um[u].pending, u)
                }
                if um[u].running > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.running, prometheus.GaugeValue, um[u].running, u)
                }
                if um[u].running_cpus > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.running_cpus, prometheus.GaugeValue, um[u].running_cpus, u)
                }
                if um[u].running_mems > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.running_mems, prometheus.GaugeValue, um[u].running_mems, u)
                }
                if um[u].suspended > 0 {
                        ch <- prometheus.MustNewConstMetric(uc.suspended, prometheus.GaugeValue, um[u].suspended, u)
                }
        }
}

