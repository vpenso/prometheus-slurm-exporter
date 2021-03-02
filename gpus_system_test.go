// +build system

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
	"testing"
)

func TestGetAllocatedGPUs(t *testing.T) {
	t.Logf("The number of Allocated GPUs: %+v", GetAllocatedGPUs())
}

func TestGetTotalGPUs(t *testing.T) {
	t.Logf("The total number of GPUs: %+v", GetTotalGPUs())
}

func TestGetGPUsMetrics(t *testing.T) {
	t.Logf("All collected GPU metrics: %+v", GetGPUsMetrics())
}
