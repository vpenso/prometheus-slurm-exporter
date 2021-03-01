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

func TestCPUsMetrics(t *testing.T) {
	file, _ := os.Open("test_data/sinfo_cpus.txt")
	data, _ := ioutil.ReadAll(file)
	cpus := ParseCPUsMetrics(data)
	if cpus.alloc != 5725.0 {
		t.Errorf("Miscount of alloc CPUs, got: %v, expected: %f", cpus.alloc, 5725.0)
	}
	if cpus.idle != 877.0 {
		t.Errorf("Miscount of idle CPUs, got: %v, expected: %f", cpus.idle, 877.0)
	}
	if cpus.other != 34.0 {
		t.Errorf("Miscount of other CPUs, got: %v, expected: %f", cpus.other, 34.0)
	}
	if cpus.total != 6636.0 {
		t.Errorf("Miscount of total CPUs, got: %v, expected: %f", cpus.total, 6636.0)
	}
}
