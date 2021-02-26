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

func TestParsePartitionsMetrics(t *testing.T) {
	// Read the input data from a file
	sinfoFile, _ := os.Open("test_data/sinfo_partitions.txt")
	squeueFile, _ := os.Open("test_data/squeue_partitions.txt")
	sinfoData, _ := ioutil.ReadAll(sinfoFile)
	squeueData, _ := ioutil.ReadAll(squeueFile)
	partitionMetrics := ParsePartitionsMetrics(sinfoData, squeueData)
	if partitionMetrics["amdcpu"].allocated != 60 {
		t.Errorf("Miscount of allocated CPUS, got: %v, wanted: %d", partitionMetrics["amdcpu"].allocated, 60)
	}
	if partitionMetrics["amdcpu"].idle != 36 {
		t.Errorf("Miscount of idle CPUS, got: %v, wanted: %d", partitionMetrics["amdcpu"].idle, 36)
	}
	if partitionMetrics["amdcpu"].other != 672 {
		t.Errorf("Miscount of other CPUS, got: %v, wanted: %d", partitionMetrics["amdcpu"].other, 672)
	}
	if partitionMetrics["amdcpu"].total != 768 {
		t.Errorf("Miscount of total CPUS, got: %v, wanted: %d", partitionMetrics["amdcpu"].total, 768)
	}
}
