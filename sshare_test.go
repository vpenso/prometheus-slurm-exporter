/* Copyright 2026 prometheus-slurm-exporter contributors

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadSshareFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Can not open test data: %v", err)
	}
	return data
}

func findSubaccount(subaccounts map[string]*FairShareMetrics, name string) *FairShareMetrics {
	for _, m := range subaccounts {
		if m.name == name {
			return m
		}
	}
	return nil
}

// The fixture is the exact `sshare -n -P -o account,fairshare` output
// reported in issue #93.
func TestParseFairShareMetricsIssueSample(t *testing.T) {
	accounts, subaccounts := ParseFairShareMetrics(loadSshareFixture(t, "test_data/sshare.txt"))

	// Compatibility contract: slurm_account_fairshare keeps exporting
	// exactly root and the top-level accounts with unchanged names/values.
	assert.Len(t, accounts, 3)
	byName := make(map[string]*FairShareMetrics)
	for _, m := range accounts {
		byName[m.name] = m
	}
	assert.Equal(t, 0.5, byName["root"].fairshare)
	assert.Equal(t, 0.999998, byName["top_1"].fairshare)
	assert.Equal(t, 0.481723, byName["top_2"].fairshare)

	// The fix for #93: every nested account is now reported.
	assert.Len(t, subaccounts, 5)
	n121 := findSubaccount(subaccounts, "nested_1_2_1")
	if assert.NotNil(t, n121) {
		assert.Equal(t, "nested_1_2", n121.parent)
		assert.Equal(t, 3, n121.depth)
		assert.Equal(t, 1.0, n121.fairshare)
	}
	n22 := findSubaccount(subaccounts, "nested_2_2")
	if assert.NotNil(t, n22) {
		assert.Equal(t, "nested_2_1", n22.parent)
		assert.Equal(t, 3, n22.depth)
		assert.Equal(t, 0.961831, n22.fairshare)
	}
	n11 := findSubaccount(subaccounts, "nested_1_1")
	if assert.NotNil(t, n11) {
		assert.Equal(t, "top_1", n11.parent)
		assert.Equal(t, 2, n11.depth)
	}
}

// Robustness against the sshare rendering variations seen across clusters:
// blank fairshare fields (fair_tree algorithm, Slurm default since 19.05),
// identically named accounts at different tree positions, CRLF endings,
// tab indentation, header/garbage lines and indent jumps.
func TestParseFairShareMetricsTreeRobustness(t *testing.T) {
	accounts, subaccounts := ParseFairShareMetrics(loadSshareFixture(t, "test_data/sshare_tree.txt"))
	assert.Len(t, accounts, 6)
	assert.Len(t, subaccounts, 3)

	byName := make(map[string]*FairShareMetrics)
	for _, m := range accounts {
		byName[m.name] = m
	}

	// Blank value -> structural position kept, no series exported (not a
	// misleading 0).
	assert.NotContains(t, byName, "interns")
	for _, m := range subaccounts {
		assert.NotEqual(t, "interns", m.name)
	}

	// Top-level set: root plus every parseable depth-1 account, with mixed
	// whitespace renderings treated as equal depth.
	assert.Equal(t, 0.1, byName["next-generation"].fairshare)
	assert.Equal(t, 0.2, byName["tails"].fairshare)
	assert.Equal(t, 0.3, byName["sysadmin"].fairshare)
	assert.Equal(t, 0.5, byName["heads"].fairshare)
	assert.Equal(t, 0.7, byName["other_tab"].fairshare)

	// Same leaf name under a same-named parent must exist as a distinct
	// series (this collision would panic a bare-name implementation).
	nestedSysadmin := findSubaccount(subaccounts, "sysadmin")
	if assert.NotNil(t, nestedSysadmin) {
		assert.Equal(t, 0.4, nestedSysadmin.fairshare)
		assert.Equal(t, "sysadmin", nestedSysadmin.parent)
		assert.Equal(t, 2, nestedSysadmin.depth)
	}

	// CRLF-rendered line parses.
	windows := findSubaccount(subaccounts, "windows")
	if assert.NotNil(t, windows) {
		assert.Equal(t, "heads", windows.parent)
		assert.Equal(t, 2, windows.depth)
	}

	// Indent jump after a tab-indented line attaches to the nearest
	// available ancestor; garbage lines are skipped entirely.
	orphan := findSubaccount(subaccounts, "orphan_jump")
	if assert.NotNil(t, orphan) {
		assert.Equal(t, "other_tab", orphan.parent)
	}
	for _, m := range accounts {
		assert.NotContains(t, m.name, "garbage")
	}
}

func TestParseFairShareMetricsEmpty(t *testing.T) {
	accounts, subaccounts := ParseFairShareMetrics([]byte(""))
	assert.Empty(t, accounts)
	assert.Empty(t, subaccounts)

	accounts, subaccounts = ParseFairShareMetrics([]byte("Account|Fairshare\n\n\n"))
	assert.Empty(t, accounts)
	assert.Empty(t, subaccounts)
}
