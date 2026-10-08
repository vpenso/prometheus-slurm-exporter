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
	"strings"

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

// JobMetrics is the per-key (account or user) job breakdown shared by the
// accounts and users collectors, which were near-identical duplicates of
// each other differing only in which job field they grouped by.
type JobMetrics struct {
	pending      float64
	pending_cpus float64
	running      float64
	running_cpus float64
	suspended    float64
}

// aggregateJobsByKey groups squeue jobs by keyFn(job) and tallies
// pending/running/suspended counts (plus pending/running CPUs) per key.
func aggregateJobsByKey(squeue *slurmcli.SqueueResponse, keyFn func(slurmcli.SqueueJob) string) map[string]*JobMetrics {
	out := make(map[string]*JobMetrics)
	for _, j := range squeue.Jobs {
		key := keyFn(j)
		if key == "" {
			continue
		}
		if _, ok := out[key]; !ok {
			out[key] = &JobMetrics{}
		}
		switch strings.ToLower(j.State()) {
		case "pending":
			out[key].pending++
			out[key].pending_cpus += float64(j.CPUs)
		case "running":
			out[key].running++
			out[key].running_cpus += float64(j.CPUs)
		case "suspended":
			out[key].suspended++
		}
	}
	return out
}
