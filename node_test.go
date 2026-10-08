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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeMetrics(t *testing.T) {
	metrics := ParseNodeMetrics(loadSinfoFixture(t))
	t.Logf("%+v", metrics)

	assert.Contains(t, metrics, "b001")
	assert.Equal(t, uint64(327680), metrics["b001"].memAlloc)
	assert.Equal(t, uint64(386000), metrics["b001"].memTotal)
	assert.Equal(t, uint64(32), metrics["b001"].cpuAlloc)
	assert.Equal(t, uint64(0), metrics["b001"].cpuIdle)
	assert.Equal(t, uint64(0), metrics["b001"].cpuOther)
	assert.Equal(t, uint64(32), metrics["b001"].cpuTotal)
	assert.Equal(t, "allocated", metrics["b001"].nodeStatus)
	assert.Equal(t, "batch", metrics["b001"].partitions)

	// g001: gres "gpu:mi250:8" with gres_used "gpu:mi250:3(IDX:0-2)"
	if assert.Len(t, metrics["g001"].gpus, 1) {
		assert.Equal(t, "mi250", metrics["g001"].gpus[0].gpuType)
		assert.Equal(t, uint64(8), metrics["g001"].gpus[0].total)
		assert.Equal(t, uint64(3), metrics["g001"].gpus[0].alloc)
	}

	// g003: untyped legacy GRES "gpu:2" yields an empty gpu_type
	if assert.Len(t, metrics["g003"].gpus, 1) {
		assert.Equal(t, "", metrics["g003"].gpus[0].gpuType)
		assert.Equal(t, uint64(2), metrics["g003"].gpus[0].total)
		assert.Equal(t, uint64(0), metrics["g003"].gpus[0].alloc)
	}

	assert.Empty(t, metrics["b001"].gpus)
}
