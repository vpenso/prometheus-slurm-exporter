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
	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

// SinfoData runs `sinfo --json` and returns the parsed response. Every
// sinfo-based collector (cpus, nodes, node, partitions, gpus) shares this
// one call shape instead of each shelling out with its own -o field list.
func SinfoData() (*slurmcli.SinfoResponse, error) {
	var resp slurmcli.SinfoResponse
	if err := slurmcli.RunJSON("sinfo", []string{"--json"}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SqueueData runs `squeue --json --states=all` and returns the parsed
// response. Shared by the queue, accounts, users and partitions (pending
// jobs) collectors.
func SqueueData() (*slurmcli.SqueueResponse, error) {
	var resp slurmcli.SqueueResponse
	if err := slurmcli.RunJSON("squeue", []string{"--json", "--states=all"}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SacctData runs `sacct --json` scoped to currently running jobs across all
// users/accounts, used by gpus.go to compute allocated GRES.
func SacctData() (*slurmcli.SacctResponse, error) {
	var resp slurmcli.SacctResponse
	args := []string{"--json", "-a", "-X", "--state=RUNNING"}
	if err := slurmcli.RunJSON("sacct", args, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
