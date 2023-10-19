/* Copyright 2021 Victor Penso
Updated by The Center for AI Safety internal only usage 2023

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
    "bytes"
    "os/exec"
    "log"
    "strings"
    "strconv"
    "github.com/prometheus/client_golang/prometheus"
)

func FairShareData() []byte {
    cmd := exec.Command("sshare", "-n", "-P")

    var out bytes.Buffer
    cmd.Stdout = &out

    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }

    return out.Bytes()
}

// 1. Modify the FairShareMetrics struct to include these two new fields.
type FairShareMetrics struct {
    fairshare   float64
    rawUsage    float64
    effectvUsage float64
}

// 2. Update the ParseFairShareMetrics function to extract and parse the values for RawUsage and EffectvUsage.
func ParseFairShareMetrics() map[string]*FairShareMetrics {
    accounts := make(map[string]*FairShareMetrics)
    lines := strings.Split(strings.TrimSpace(string(FairShareData())), "\n")
    for _, line := range lines {
        parts := strings.Split(line, "|")
        if len(parts) < 7 { // Assuming that the required values are in columns 5 and 6
            continue
        }
        account := strings.TrimSpace(parts[0])
        if _, exists := accounts[account]; !exists {
            accounts[account] = &FairShareMetrics{}
        }

        // Original fairshare parsing code
        fairshare, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
        if err != nil {
            log.Printf("Failed to parse fairshare: %s", err)
            continue
        }
        accounts[account].fairshare = fairshare

        // New code for RawUsage parsing
        rawUsage, err := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
        if err != nil {
            log.Printf("Failed to parse rawUsage: %s", err)
            continue
        }
        accounts[account].rawUsage = rawUsage

        // New code for EffectvUsage parsing
        effectvUsage, err := strconv.ParseFloat(strings.TrimSpace(parts[5]), 64)
        if err != nil {
            log.Printf("Failed to parse effectvUsage: %s", err)
            continue
        }
        accounts[account].effectvUsage = effectvUsage
    }
    return accounts
}

// 3. Update the FairShareCollector struct to describe these new metrics.
type FairShareCollector struct {
    fairshare   *prometheus.Desc
    rawUsage    *prometheus.Desc
    effectvUsage *prometheus.Desc
}

func NewFairShareCollector() *FairShareCollector {
    labels := []string{"account"}
    return &FairShareCollector{
        fairshare:   prometheus.NewDesc("slurm_account_fairshare", "FairShare for account", labels, nil),
        rawUsage:    prometheus.NewDesc("slurm_account_rawUsage", "RawUsage for account", labels, nil),
        effectvUsage: prometheus.NewDesc("slurm_account_effectvUsage", "EffectvUsage for account", labels, nil),
    }
}

// 4. Update the Collect method of the FairShareCollector to send these metrics to Prometheus.
func (fsc *FairShareCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- fsc.fairshare
    ch <- fsc.rawUsage
    ch <- fsc.effectvUsage
}

func (fsc *FairShareCollector) Collect(ch chan<- prometheus.Metric) {
    fsm := ParseFairShareMetrics()
    for f := range fsm {
        ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, fsm[f].fairshare, f)
        ch <- prometheus.MustNewConstMetric(fsc.rawUsage, prometheus.GaugeValue, fsm[f].rawUsage, f)
        ch <- prometheus.MustNewConstMetric(fsc.effectvUsage, prometheus.GaugeValue, fsm[f].effectvUsage, f)
    }
}
