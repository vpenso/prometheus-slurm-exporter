/* Copyright 2017 Victor Penso, Matteo Dessalvi

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

func loadSinfoFixture(t *testing.T) *slurmcli.SinfoResponse {
	t.Helper()
	data, err := ioutil.ReadFile("test_data/sinfo.json")
	if err != nil {
		t.Fatalf("Can not open test data: %v", err)
	}
	var resp slurmcli.SinfoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("Can not parse test data: %v", err)
	}
	return &resp
}

func TestCPUsMetrics(t *testing.T) {
	cm := ParseCPUsMetrics(loadSinfoFixture(t))
	assert.Equal(t, 68.0, cm.alloc)
	assert.Equal(t, 148.0, cm.idle)
	assert.Equal(t, 40.0, cm.other)
	assert.Equal(t, 256.0, cm.total)
}
