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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

func TestParseAccountsMetrics(t *testing.T) {
	am := ParseAccountsMetrics(loadSqueueFixture(t))

	// physics: 101 pending (4 cpus), 102 pending (8 cpus), 105 cancelled (ignored)
	assert.Equal(t, 2.0, am["physics"].pending)
	assert.Equal(t, 12.0, am["physics"].pending_cpus)
	assert.Equal(t, 0.0, am["physics"].running)

	// chemistry: 103 running (16 cpus), 104 suspended
	assert.Equal(t, 1.0, am["chemistry"].running)
	assert.Equal(t, 16.0, am["chemistry"].running_cpus)
	assert.Equal(t, 1.0, am["chemistry"].suspended)
}

func TestParseUsersMetrics(t *testing.T) {
	um := ParseUsersMetrics(loadSqueueFixture(t))

	// alice: 101 pending (4 cpus), 103 running (16 cpus)
	assert.Equal(t, 1.0, um["alice"].pending)
	assert.Equal(t, 4.0, um["alice"].pending_cpus)
	assert.Equal(t, 1.0, um["alice"].running)
	assert.Equal(t, 16.0, um["alice"].running_cpus)

	// bob: 102 pending (8 cpus), 105 cancelled (ignored)
	assert.Equal(t, 1.0, um["bob"].pending)
	assert.Equal(t, 8.0, um["bob"].pending_cpus)
}

func TestAggregateJobsByKeySkipsEmptyKey(t *testing.T) {
	squeue := &slurmcli.SqueueResponse{Jobs: []slurmcli.SqueueJob{
		{JobState: []string{"RUNNING"}, Account: ""},
	}}
	out := aggregateJobsByKey(squeue, func(j slurmcli.SqueueJob) string { return j.Account })
	assert.Empty(t, out)
}
