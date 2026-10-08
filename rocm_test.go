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
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vpenso/prometheus-slurm-exporter/internal/smicli"
)

func TestParseROCmMetricFixture(t *testing.T) {
	data, err := ioutil.ReadFile("test_data/rocm_metric.json")
	if err != nil {
		t.Fatalf("Can not open test data: %v", err)
	}
	var devices []smicli.ROCmDevice
	if err := json.Unmarshal(data, &devices); err != nil {
		t.Fatalf("Can not parse test data: %v", err)
	}

	assert.Len(t, devices, 2)
	assert.Equal(t, 0, devices[0].GPU)
	assert.Equal(t, 87.5, devices[0].Usage.GfxActivity)
	assert.Equal(t, uint64(34359738368), devices[0].Mem.Used)
	assert.Equal(t, uint64(68719476736), devices[0].Mem.Total)
	assert.Equal(t, 62.0, devices[0].Temperature.Edge)
	assert.Equal(t, 310.2, devices[0].Power.Socket)
}
