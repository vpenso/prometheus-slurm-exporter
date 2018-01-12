/* Copyright 2017 Victor Penso, Matteo Dessalvi

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
  "os"
  "io/ioutil"
)

func TestParseQueueMetrics(t *testing.T) {
  // Read the input data from a file
  file, err := os.Open("test_data/squeue.txt")
  if err != nil { t.Fatalf("Can not open test data: %v", err) }
  data, err := ioutil.ReadAll(file)
  t.Logf("%+v", ParseQueueMetrics(data))
}

func TestQueueGetMetrics(t *testing.T) {
  t.Logf("%+v", QueueGetMetrics())
}
