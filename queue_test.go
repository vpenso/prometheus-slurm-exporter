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

func loadSqueueFixture(t *testing.T) *slurmcli.SqueueResponse {
	t.Helper()
	data, err := ioutil.ReadFile("test_data/squeue.json")
	if err != nil {
		t.Fatalf("Can not open test data: %v", err)
	}
	var resp slurmcli.SqueueResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("Can not parse test data: %v", err)
	}
	return &resp
}

func TestParseQueueMetrics(t *testing.T) {
	qm := ParseQueueMetrics(loadSqueueFixture(t))
	assert.Equal(t, 2.0, qm.pending)
	assert.Equal(t, 1.0, qm.pending_dep)
	assert.Equal(t, 1.0, qm.running)
	assert.Equal(t, 1.0, qm.suspended)
	assert.Equal(t, 1.0, qm.cancelled)
	assert.Equal(t, 0.0, qm.completing)
	assert.Equal(t, 0.0, qm.completed)
	assert.Equal(t, 0.0, qm.configuring)
	assert.Equal(t, 0.0, qm.failed)
	assert.Equal(t, 0.0, qm.timeout)
	assert.Equal(t, 0.0, qm.preempted)
	assert.Equal(t, 0.0, qm.node_fail)
}
