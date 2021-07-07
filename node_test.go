/* Copyright 2021 Chris Read

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
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
For this example data line:

a048,79384,193000,3/13/0/16,mix

We want output that looks like:

slurm_node_cpus_allocated{name="a048",status="mix"} 3
slurm_node_cpus_idle{name="a048",status="mix"} 3
slurm_node_cpus_other{name="a048",status="mix"} 0
slurm_node_cpus_total{name="a048",status="mix"} 16
slurm_node_mem_allocated{name="a048",status="mix"} 179384
slurm_node_mem_total{name="a048",status="mix"} 193000
slurm_node_gpu_allocated{gputype="rtx5000",name="a048",status="mix"} 2
slurm_node_gpu_total{gputype="rtx5000",name="a048",status="mix"} 4
*/

func TestNodeMetrics(t *testing.T) {
	// Read the input data from a file
	data, err := ioutil.ReadFile("test_data/sinfo_mem.txt")
	if err != nil {
		t.Fatalf("Can not open test data: %v", err)
	}
	metrics := ParseNodeMetrics(data)
	t.Logf("%+v", metrics)

	assert.Contains(t, metrics, "b001")
	assert.Equal(t, uint64(327680), metrics["b001"].memAlloc)
	assert.Equal(t, uint64(386000), metrics["b001"].memTotal)
	assert.Equal(t, uint64(32), metrics["b001"].cpuAlloc)
	assert.Equal(t, uint64(0), metrics["b001"].cpuIdle)
	assert.Equal(t, uint64(0), metrics["b001"].cpuOther)
	assert.Equal(t, uint64(32), metrics["b001"].cpuTotal)
	assert.Equal(t, uint64(4), metrics["b001"].gpuAlloc)
	assert.Equal(t, uint64(4), metrics["b001"].gpuTotal)
}
