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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

func ParseUsersMetrics(squeue *slurmcli.SqueueResponse) map[string]*JobMetrics {
	return aggregateJobsByKey(squeue, func(j slurmcli.SqueueJob) string { return j.UserName })
}

// ParseUserHostMetrics counts running jobs per user and host. Host tokens
// are taken verbatim from the job's node list expression, so range
// expressions such as "b00[1-3]" become a single label value.
func ParseUserHostMetrics(squeue *slurmcli.SqueueResponse) map[string]map[string]int {
	hosts := make(map[string]map[string]int)
	for _, j := range squeue.Jobs {
		if strings.ToLower(j.State()) != "running" || j.UserName == "" {
			continue
		}
		if _, ok := hosts[j.UserName]; !ok {
			hosts[j.UserName] = make(map[string]int)
		}
		for _, host := range strings.Split(j.Nodes, ",") {
			if host = strings.TrimSpace(host); host != "" {
				hosts[j.UserName][host]++
			}
		}
	}
	return hosts
}

type UsersCollector struct {
	pending      *prometheus.Desc
	pending_cpus *prometheus.Desc
	running      *prometheus.Desc
	running_cpus *prometheus.Desc
	running_mem  *prometheus.Desc
	hosts        *prometheus.Desc
	suspended    *prometheus.Desc
}

func NewUsersCollector() *UsersCollector {
	labels := []string{"user"}
	return &UsersCollector{
		pending:      prometheus.NewDesc("slurm_user_jobs_pending", "Pending jobs for user", labels, nil),
		pending_cpus: prometheus.NewDesc("slurm_user_cpus_pending", "Pending cpus for user", labels, nil),
		running:      prometheus.NewDesc("slurm_user_jobs_running", "Running jobs for user", labels, nil),
		running_cpus: prometheus.NewDesc("slurm_user_cpus_running", "Running cpus for user", labels, nil),
		running_mem:  prometheus.NewDesc("slurm_user_mem_running", "Running memory for user (MB)", labels, nil),
		hosts:        prometheus.NewDesc("slurm_user_jobs_running_host", "Running jobs for user on host", []string{"user", "host"}, nil),
		suspended:    prometheus.NewDesc("slurm_user_jobs_suspended", "Suspended jobs for user", labels, nil),
	}
}

func (uc *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- uc.pending
	ch <- uc.pending_cpus
	ch <- uc.running
	ch <- uc.running_cpus
	ch <- uc.running_mem
	ch <- uc.hosts
	ch <- uc.suspended
}

func (uc *UsersCollector) Collect(ch chan<- prometheus.Metric) {
	squeue, err := SqueueData()
	if err != nil {
		log.Errorf("users: %v", err)
		return
	}
	um := ParseUsersMetrics(squeue)
	for u := range um {
		if um[u].pending > 0 {
			ch <- prometheus.MustNewConstMetric(uc.pending, prometheus.GaugeValue, um[u].pending, u)
		}
		if um[u].pending_cpus > 0 {
			ch <- prometheus.MustNewConstMetric(uc.pending_cpus, prometheus.GaugeValue, um[u].pending_cpus, u)
		}
		if um[u].running > 0 {
			ch <- prometheus.MustNewConstMetric(uc.running, prometheus.GaugeValue, um[u].running, u)
		}
		if um[u].running_cpus > 0 {
			ch <- prometheus.MustNewConstMetric(uc.running_cpus, prometheus.GaugeValue, um[u].running_cpus, u)
		}
		if um[u].running_mem > 0 {
			ch <- prometheus.MustNewConstMetric(uc.running_mem, prometheus.GaugeValue, um[u].running_mem, u)
		}
		if um[u].suspended > 0 {
			ch <- prometheus.MustNewConstMetric(uc.suspended, prometheus.GaugeValue, um[u].suspended, u)
		}
	}
	for u, byHost := range ParseUserHostMetrics(squeue) {
		for host, count := range byHost {
			ch <- prometheus.MustNewConstMetric(uc.hosts, prometheus.GaugeValue, float64(count), u, host)
		}
	}
}
