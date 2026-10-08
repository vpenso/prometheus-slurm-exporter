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

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

func loadSacctFixture(t *testing.T) *slurmcli.SacctResponse {
	t.Helper()
	data, err := ioutil.ReadFile("test_data/sacct.json")
	if err != nil {
		t.Fatalf("Can not open test data: %v", err)
	}
	var resp slurmcli.SacctResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("Can not parse test data: %v", err)
	}
	return &resp
}

// Regression test for issue #38: typed GRES entries (AMD "gpu:mi250:N",
// NVIDIA "gpu:a100:N") and the untyped legacy form ("gpu:N") must all be
// parsed correctly and kept separate by gpu_type.
func TestParseGPUsMetrics(t *testing.T) {
	sinfo := loadSinfoFixture(t)
	sacct := loadSacctFixture(t)

	gm := ParseGPUsMetrics(sinfo, sacct)

	assert.Equal(t, 8.0, gm.total["mi250"])
	assert.Equal(t, 4.0, gm.total["a100"])
	assert.Equal(t, 2.0, gm.total[""])

	assert.Equal(t, 3.0, gm.alloc["mi250"])
	assert.Equal(t, 4.0, gm.alloc["a100"])
	assert.Equal(t, 1.0, gm.alloc[""])

	assert.Equal(t, 5.0, gm.idle["mi250"])
	assert.Equal(t, 0.0, gm.idle["a100"])
	assert.Equal(t, 1.0, gm.idle[""])
}

func TestParseTotalGPUsHandlesTypedAndUntypedGres(t *testing.T) {
	totals := ParseTotalGPUs(loadSinfoFixture(t))
	assert.Equal(t, 8.0, totals["mi250"])
	assert.Equal(t, 4.0, totals["a100"])
	assert.Equal(t, 2.0, totals[""])
}
