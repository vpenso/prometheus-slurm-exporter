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

func TestParseAllocatedGPUs(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("test_data/sacct_gpus.txt")
	data, _ := ioutil.ReadAll(file)
	allocatedGPUs := ParseAllocatedGPUs(data)
	if allocatedGPUs != 81 {
		t.Errorf("Miscount of allocated GPUs, got: %f, wanted: %d", allocatedGPUs, 81)
	}
}

func TestParseTotalGPUs(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("test_data/sinfo_gpus.txt")
	data, _ := ioutil.ReadAll(file)
	totalGPUs := ParseTotalGPUs(data)
	if totalGPUs != 164 {
		t.Errorf("Miscount of total GPUs, got: %f, wanted: %d", totalGPUs, 164)
	}
}
