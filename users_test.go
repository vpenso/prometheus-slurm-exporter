// +build unit

/* Copyright 2021 Rovanion Luckey

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

func TestParseUsersMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("test_data/squeue_users.txt")
	data, _ := ioutil.ReadAll(file)
	users := ParseUsersMetrics(data)

	if users["slurm"].pending != 434.0 {
		t.Errorf("Miscount of pending user jobs, got: %v, wanted: %f", users["slurm"].pending, 434.0)
	}
	if users["slurm"].running != 81.0 {
		t.Errorf("Miscount of running user jobs, got: %v, wanted: %f", users["slurm"].running, 81.0)
	}
	if users["slurm"].running_cpus != 806.0 {
		t.Errorf("Miscount of running_cpus user jobs, got: %v, wanted: %f", users["slurm"].running_cpus, 806.0)
	}
	if users["slurm"].suspended != 0.0 {
		t.Errorf("Miscount of suspended user jobs, got: %v, wanted: %f", users["slurm"].suspended, 0.0)
	}
}
