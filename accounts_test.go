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

func TestAccountsMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("test_data/squeue_no_accounts.txt")
	data, _ := ioutil.ReadAll(file)
	accounts := ParseAccountsMetrics(data)

	if accounts["(null)"].pending != 449.0 {
		t.Errorf("Miscount of pending account jobs, got: %v, wanted: %f", accounts["(null)"].pending, 449.0)
	}
	if accounts["(null)"].running != 79.0 {
		t.Errorf("Miscount of running account jobs, got: %v, wanted: %f", accounts["(null)"].running, 79.0)
	}
	if accounts["(null)"].running_cpus != 798.0 {
		t.Errorf("Miscount of running_cpus account jobs, got: %v, wanted: %f", accounts["(null)"].running_cpus, 798.0)
	}
	if accounts["(null)"].suspended != 0.0 {
		t.Errorf("Miscount of suspended account jobs, got: %v, wanted: %f", accounts["(null)"].suspended, 0.0)
	}

}
