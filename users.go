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

func UsersData() []byte {
        cmd := exec.Command("squeue", "-h", "-o '%A|%u|%T|%C'")
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
        running float64
}

func ParseUsersMetrics(input []byte) map[string]*UserJobMetrics {
        users := make(map[string]*UserJobMetrics)
        lines := strings.Split(string(input), "\n")
        for _, line := range lines {
                if strings.Contains(line,"|") {
                        user := strings.Split(line,"|")[1]
                        _,key := users[user]
                        if !key {
                                users[user] = &UserJobMetrics{0}
                        }
                        state := strings.Split(line,"|")[2]
                        state = strings.ToLower(state)
                        running := regexp.MustCompile(`^running`)
                        switch {
                        case running.MatchString(state) == true:
                                users[user].running++
                        }
                }
        }
        return users
}

type UsersCollector struct {
        running *prometheus.Desc
}

func NewUsersCollector() *UsersCollector {
        labels := []string{"user"}
        return &UsersCollector {
                running: prometheus.NewDesc("slurm_user_jobs_running", "Running jobs for user", labels, nil),
        }
}

func (uc *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
        ch <- uc.running
}

func (uc *UsersCollector) Collect(ch chan<- prometheus.Metric) {
        um := ParseUsersMetrics(UsersData())
        for u := range um {
                ch <- prometheus.MustNewConstMetric(uc.running, prometheus.GaugeValue, um[u].running, u)
        }
}

