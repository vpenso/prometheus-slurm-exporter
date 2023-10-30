/* Copyright 2022 Iztok Lebar Bajec

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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGPUsMetrics(t *testing.T) {
	test_data_paths, _ := filepath.Glob("test_data/slurm-*")
	for _, test_data_path := range test_data_paths {
		slurm_version := strings.TrimPrefix(test_data_path, "test_data/slurm-")
		t.Logf("slurm-%s", slurm_version)

		// Read the allocated GPUs data and log the result
		allocatedData, err := ioutil.ReadFile(filepath.Join(test_data_path, "sinfo_gpus_allocated.txt"))
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		allocatedGPUs := ParseAllocatedGPUs(allocatedData)
		t.Logf("Allocated GPUs: %v", allocatedGPUs)

		// Read the idle GPUs data and log the result
		idleData, err := ioutil.ReadFile(filepath.Join(test_data_path, "sinfo_gpus_idle.txt"))
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		idleGPUs := ParseIdleGPUs(idleData)
		t.Logf("Idle GPUs: %v", idleGPUs)

		// Read the total GPUs data and log the result
		totalData, err := ioutil.ReadFile(filepath.Join(test_data_path, "sinfo_gpus_total.txt"))
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		totalGPUs := ParseTotalGPUs(totalData)
		t.Logf("Total GPUs: %v", totalGPUs)
	}
}

func TestGPUsGetMetrics(t *testing.T) {
	metrics := GPUsGetMetrics()
	t.Logf("Allocated GPUs: %v", metrics.alloc)
	t.Logf("Idle GPUs: %v", metrics.idle)
	t.Logf("Other GPUs: %v", metrics.other)
	t.Logf("Total GPUs: %v", metrics.total)
	t.Logf("Utilization: %v", metrics.utilization)
	t.Logf("User GPUs DCGM: %v", metrics.UserGPUsDCGM)
	t.Logf("User GPUs SLURM: %v", metrics.UserGPUsSLURM)
}
