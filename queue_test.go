// +build unit

/* Copyright 2017 Victor Penso, Matteo Dessalvi
   Copyright 2021 Rovanion Luckey

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
	"testing"
)

func TestParseQueueMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("test_data/squeue.txt")
	data, _ := ioutil.ReadAll(file)
	queueMetrics := ParseQueueMetrics(data)

	if queueMetrics.pending != 4.0 {
		t.Errorf("Miscount of pending jobs, got %v, wanted: %f", queueMetrics.pending, 4.0)
	}
	if queueMetrics.pending_dep != 0.0 {
		t.Errorf("Miscount of pending_dep jobs, got %v, wanted: %f", queueMetrics.pending_dep, 0.0)
	}
	if queueMetrics.running != 28.0 {
		t.Errorf("Miscount of running jobs, got %v, wanted: %f", queueMetrics.running, 28.0)
	}
	if queueMetrics.suspended != 1.0 {
		t.Errorf("Miscount of suspended jobs, got %v, wanted: %f", queueMetrics.suspended, 1.0)
	}
	if queueMetrics.cancelled != 1.0 {
		t.Errorf("Miscount of cancelled jobs, got %v, wanted: %f", queueMetrics.cancelled, 1.0)
	}
	if queueMetrics.completing != 2.0 {
		t.Errorf("Miscount of completing jobs, got %v, wanted: %f", queueMetrics.completing, 2.0)
	}
	if queueMetrics.completed != 1.0 {
		t.Errorf("Miscount of completed jobs, got %v, wanted: %f", queueMetrics.completed, 1.0)
	}
	if queueMetrics.configuring != 1.0 {
		t.Errorf("Miscount of configuring jobs, got %v, wanted: %f", queueMetrics.configuring, 1.0)
	}
	if queueMetrics.failed != 1.0 {
		t.Errorf("Miscount of failed jobs, got %v, wanted: %f", queueMetrics.failed, 1.0)
	}
	if queueMetrics.timeout != 1.0 {
		t.Errorf("Miscount of timeout jobs, got %v, wanted: %f", queueMetrics.timeout, 1.0)
	}
	if queueMetrics.preempted != 1.0 {
		t.Errorf("Miscount of preempted jobs, got %v, wanted: %f", queueMetrics.preempted, 1.0)
	}
	if queueMetrics.node_fail != 1.0 {
		t.Errorf("Miscount of node_fail jobs, got %v, wanted: %f", queueMetrics.node_fail, 1.0)
	}
}
