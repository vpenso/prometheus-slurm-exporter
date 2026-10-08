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

	"github.com/stretchr/testify/assert"
)

func TestNodesMetrics(t *testing.T) {
	nm := ParseNodesMetrics(loadSinfoFixture(t))
	assert.Equal(t, 1.0, nm.alloc) // b001: ALLOCATED
	assert.Equal(t, 0.0, nm.comp)
	assert.Equal(t, 1.0, nm.down)  // g003: DOWN
	assert.Equal(t, 1.0, nm.drain) // b003: MIXED+DRAIN
	assert.Equal(t, 0.0, nm.err)
	assert.Equal(t, 0.0, nm.fail)
	assert.Equal(t, 2.0, nm.idle) // b002, g002: IDLE
	assert.Equal(t, 0.0, nm.maint)
	assert.Equal(t, 1.0, nm.mix) // g001: MIXED
	assert.Equal(t, 0.0, nm.resv)
}
