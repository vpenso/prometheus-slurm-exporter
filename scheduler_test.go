// +build unit

/* Copyright 2017 Victor Penso, Matteo Dessalvi, Rovanion Luckey

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

func TestSchedulerMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("test_data/sdiag.txt")
	data, _ := ioutil.ReadAll(file)
	schedulerMetrics := ParseSchedulerMetrics(data)

	if schedulerMetrics.threads != 3.0 {
		t.Errorf("Miscount of threads, got %v, wanted: %f", schedulerMetrics.threads, 3.0)
	}
	if schedulerMetrics.queue_size != 0.0 {
		t.Errorf("Miscount of queue_size, got %v, wanted: %f", schedulerMetrics.queue_size, 0.0)
	}
	if schedulerMetrics.dbd_queue_size != 0.0 {
		t.Errorf("Miscount of dbd_queue_size, got %v, wanted: %f", schedulerMetrics.dbd_queue_size, 0.0)
	}
	if schedulerMetrics.last_cycle != 97209.0 {
		t.Errorf("Miscount of last_cycle, got %v, wanted: %f", schedulerMetrics.last_cycle, 97209.0)
	}
	if schedulerMetrics.mean_cycle != 74593.0 {
		t.Errorf("Miscount of mean_cycle, got %v, wanted: %f", schedulerMetrics.mean_cycle, 74593.0)
	}
	if schedulerMetrics.cycle_per_minute != 63.0 {
		t.Errorf("Miscount of cycle_per_minute, got %v, wanted: %f", schedulerMetrics.cycle_per_minute, 63.0)
	}
	if schedulerMetrics.backfill_last_cycle != 1.94289e+06 {
		t.Errorf("Miscount of backfill_last_cycle, got %v, wanted: %f", schedulerMetrics.backfill_last_cycle, 1.94289e+06)
	}
	if schedulerMetrics.backfill_mean_cycle != 1.96082e+06 {
		t.Errorf("Miscount of backfill_mean_cycle, got %v, wanted: %f", schedulerMetrics.backfill_mean_cycle, 1.96082e+06)
	}
	if schedulerMetrics.backfill_depth_mean != 29324.0 {
		t.Errorf("Miscount of backfill_depth_mean, got %v, wanted: %f", schedulerMetrics.backfill_depth_mean, 29324.0)
	}
	if schedulerMetrics.total_backfilled_jobs_since_start != 111544.0 {
		t.Errorf("Miscount of total_backfilled_jobs_since_start, got %v, wanted: %f", schedulerMetrics.total_backfilled_jobs_since_start, 111544.0)
	}
	if schedulerMetrics.total_backfilled_jobs_since_cycle != 793.0 {
		t.Errorf("Miscount of total_backfilled_jobs_since_cycle, got %v, wanted: %f", schedulerMetrics.total_backfilled_jobs_since_cycle, 793.0)
	}
	if schedulerMetrics.total_backfilled_heterogeneous != 10.0 {
		t.Errorf("Miscount of total_backfilled_heterogeneous, got %v, wanted: %f", schedulerMetrics.total_backfilled_heterogeneous, 10.0)
	}
}
