/* Copyright 2021 Victor Penso

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

// This collector intentionally still parses sshare's legacy pipe-delimited
// text output (-P) rather than `sshare --json`, unlike the sinfo/squeue
// -based collectors elsewhere in this project. sshare has no --json mode at
// all, and its text format here (tree indentation plus pipe-delimited
// fields) is the only interface available; the parser below derives the
// account tree from the indentation *relative* between lines rather than
// assuming a specific indent width, so clusters rendering 1-space, 2-space
// or tab indentation all yield the same logical depths. Unparseable or
// blank lines are skipped, never fatal. Revisit if sshare's text format is
// ever seen to drift beyond what this tolerant parsing absorbs.
package main

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/vpenso/prometheus-slurm-exporter/internal/slurmcli"
)

// FairShareData runs `sshare -n -P -o account,fairshare` and returns its
// raw output. Unlike earlier revisions of this file, a failing sshare call
// is reported as an error to the caller instead of terminating the whole
// exporter; the collector then simply exports nothing for this scrape.
func FairShareData() ([]byte, error) {
	return slurmcli.Run("sshare", []string{"-n", "-P", "-o", "account,fairshare"})
}

// FairShareMetrics holds one parsed sshare row. parent and depth describe
// the position in the account tree reconstructed from sshare's tree
// indentation: depth is the number of ancestors (root = 0, top-level
// accounts = 1), independent of the concrete indent width used by the
// local sshare rendering.
type FairShareMetrics struct {
	name      string
	parent    string
	depth     int
	fairshare float64
}

// ParseFairShareMetrics walks sshare's tree-rendered output and splits the
// accounts into two sets: root and top-level accounts (depth <= 1), which
// preserve the historical slurm_account_fairshare contract, and nested
// subaccounts (depth >= 2), reported with parent/depth context.
//
// Robustness rules (sshare output varies across clusters and Slurm
// versions): lines without a '|' separator (headers, garbage) are skipped;
// lines whose fairshare field is blank or non-numeric are kept in the tree
// structure but export no series - notably the blank account-level
// fairshare that sshare renders when the fair_tree algorithm (Slurm
// default since 19.05) is active, which must not be reported as a
// misleading 0.
func ParseFairShareMetrics(input []byte) (map[string]*FairShareMetrics, map[string]*FairShareMetrics) {
	accounts := make(map[string]*FairShareMetrics)
	subaccounts := make(map[string]*FairShareMetrics)

	type stackEntry struct {
		path  string // full ancestor path, used as map key so identically named accounts under different parents cannot collide
		name  string
		depth int // raw indent width (chars), only compared relative to neighbors
	}
	var stack []stackEntry

	for _, line := range strings.Split(string(input), "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.SplitN(line, "|", 2)
		if len(fields) != 2 {
			continue
		}
		rawName := fields[0]
		name := strings.TrimSpace(rawName)
		if name == "" {
			continue
		}

		// Maintain the ancestor stack: drop entries at or deeper than the
		// current line's indent; the remaining top is the parent.
		indent := leadingWhitespaceLen(rawName)
		for len(stack) > 0 && stack[len(stack)-1].depth >= indent {
			stack = stack[:len(stack)-1]
		}
		parent := ""
		path := name
		if len(stack) > 0 {
			parent = stack[len(stack)-1].name
			path = stack[len(stack)-1].path + "\x00" + name
		}
		level := len(stack)
		stack = append(stack, stackEntry{path: path, name: name, depth: indent})

		fairshare, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			// Blank/invalid fairshare (e.g. fair_tree clusters): structural
			// position was recorded above; export nothing for this row.
			continue
		}

		m := &FairShareMetrics{name: name, parent: parent, depth: level, fairshare: fairshare}
		if level <= 1 {
			accounts[path] = m
		} else {
			subaccounts[path] = m
		}
	}
	return accounts, subaccounts
}

func leadingWhitespaceLen(s string) int {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return i
}

type FairShareCollector struct {
	fairshare  *prometheus.Desc
	subaccount *prometheus.Desc
}

func NewFairShareCollector() *FairShareCollector {
	return &FairShareCollector{
		fairshare:  prometheus.NewDesc("slurm_account_fairshare", "FairShare for root and top-level Slurm accounts", []string{"account"}, nil),
		subaccount: prometheus.NewDesc("slurm_subaccount_fairshare", "FairShare for nested Slurm subaccounts", []string{"account", "parent_account", "account_depth"}, nil),
	}
}

func (fsc *FairShareCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- fsc.fairshare
	ch <- fsc.subaccount
}

func (fsc *FairShareCollector) Collect(ch chan<- prometheus.Metric) {
	data, err := FairShareData()
	if err != nil {
		log.Errorf("fairshare: %v", err)
		return
	}
	accounts, subaccounts := ParseFairShareMetrics(data)
	for _, m := range accounts {
		ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, m.fairshare, m.name)
	}
	for _, m := range subaccounts {
		ch <- prometheus.MustNewConstMetric(fsc.subaccount, prometheus.GaugeValue, m.fairshare, m.name, m.parent, strconv.Itoa(m.depth))
	}
}
