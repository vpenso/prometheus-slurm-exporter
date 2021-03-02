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

func TestNodesMetrics(t *testing.T) {
	// Read the input data from a file
	file, _ := os.Open("test_data/sinfo_nodes.txt")
	data, _ := ioutil.ReadAll(file)
	nodes := ParseNodesMetrics(data)

	if nodes.alloc != 1.0 {
		t.Errorf("Miscount of alloc nodes, got: %v, wanted: %f", nodes.alloc, 1.0)
	}
	if nodes.comp != 0.0 {
		t.Errorf("Miscount of comp nodes, got: %v, wanted: %f", nodes.comp, 0.0)
	}
	if nodes.down != 0.0 {
		t.Errorf("Miscount of down nodes, got: %v, wanted: %f", nodes.down, 0.0)
	}
	if nodes.drain != 7.0 {
		t.Errorf("Miscount of drain nodes, got: %v, wanted: %f", nodes.drain, 7.0)
	}
	if nodes.err != 0.0 {
		t.Errorf("Miscount of err nodes, got: %v, wanted: %f", nodes.err, 0.0)
	}
	if nodes.fail != 0.0 {
		t.Errorf("Miscount of fail nodes, got: %v, wanted: %f", nodes.fail, 0.0)
	}
	if nodes.idle != 9.0 {
		t.Errorf("Miscount of idle nodes, got: %v, wanted: %f", nodes.idle, 9.0)
	}
	if nodes.maint != 0.0 {
		t.Errorf("Miscount of maint nodes, got: %v, wanted: %f", nodes.maint, 0.0)
	}
	if nodes.mix != 20.0 {
		t.Errorf("Miscount of mix nodes, got: %v, wanted: %f", nodes.mix, 20.0)
	}
	if nodes.resv != 0.0 {
		t.Errorf("Miscount of resv nodes, got: %v, wanted: %f", nodes.resv, 0.0)
	}

}
