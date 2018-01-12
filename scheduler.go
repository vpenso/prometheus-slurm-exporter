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
  "log"
  "strings"
  "strconv"
  "io/ioutil"
  "os/exec"
  "github.com/prometheus/client_golang/prometheus"
)

/*
 * Execute the Slurm sdiag command to read the current statistics
 * from the Slurm scheduler. It will be repreatedly called by the
 * collector.
 */

// Basic metrics for the scheduler
type SchedulerMetrics struct {
  threads             float64
  queue_size          float64
  last_cycle          float64
  mean_cycle          float64
  cycle_per_minute    float64
  backfill_last_cycle float64
  backfill_mean_cycle float64
  backfill_depth_mean float64
}

// Execute the sdiag command and return its output
func SchedulerData() []byte {
  cmd := exec.Command("/usr/bin/sdiag")
  stdout, err := cmd.StdoutPipe()
  if err != nil { log.Fatal(err) }
  if err := cmd.Start(); err != nil { log.Fatal(err) }
  out, _ := ioutil.ReadAll(stdout)
  if err := cmd.Wait(); err != nil { log.Fatal(err) }
  return out
}

// Helper function to split a single line from the sdiag output
func SplitColonValueToFloat(input string) float64 {
  // Verify if the string extracted from the sdiag output is empty
  if input == "" { return 0 }
  str := strings.Split(input,":")
  rvalue := strings.TrimSpace(str[1])
  flt, _ := strconv.ParseFloat(rvalue,64)
  return flt
}

// Extract the relevant metrics from the sdiag output
func ParseSchedulerMetrics(input []byte) *SchedulerMetrics {
  lines := strings.Split(string(input),"\n")
  var sm SchedulerMetrics
  sm.threads             = SplitColonValueToFloat(lines[4])
  sm.queue_size          = SplitColonValueToFloat(lines[5])
  sm.last_cycle          = SplitColonValueToFloat(lines[14])
  sm.mean_cycle          = SplitColonValueToFloat(lines[17])
  sm.cycle_per_minute    = SplitColonValueToFloat(lines[19])
  sm.backfill_last_cycle = SplitColonValueToFloat(lines[27])
  sm.backfill_mean_cycle = SplitColonValueToFloat(lines[29])
  sm.backfill_depth_mean = SplitColonValueToFloat(lines[32])
  return &sm
}

// Returns the scheduler metrics
func SchedulerGetMetrics() *SchedulerMetrics {
  return ParseSchedulerMetrics(SchedulerData())
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */


// Collector strcture
type SchedulerCollector struct {
  threads             *prometheus.Desc
  queue_size          *prometheus.Desc
  last_cycle          *prometheus.Desc
  mean_cycle          *prometheus.Desc
  cycle_per_minute    *prometheus.Desc
  backfill_last_cycle *prometheus.Desc
  backfill_mean_cycle *prometheus.Desc
  backfill_depth_mean *prometheus.Desc
}

// Send all metric descriptions
func (c *SchedulerCollector) Describe(ch chan<- *prometheus.Desc) {
  ch <- c.threads
  ch <- c.queue_size
  ch <- c.last_cycle
  ch <- c.mean_cycle
  ch <- c.cycle_per_minute
  ch <- c.backfill_last_cycle
  ch <- c.backfill_mean_cycle
  ch <- c.backfill_depth_mean
}

// Send the values of all metrics
func (sc *SchedulerCollector) Collect(ch chan<- prometheus.Metric) {
  sm := SchedulerGetMetrics()
  ch <- prometheus.MustNewConstMetric(sc.threads, prometheus.GaugeValue, sm.threads)
  ch <- prometheus.MustNewConstMetric(sc.queue_size, prometheus.GaugeValue, sm.queue_size)
  ch <- prometheus.MustNewConstMetric(sc.last_cycle, prometheus.GaugeValue, sm.last_cycle)
  ch <- prometheus.MustNewConstMetric(sc.mean_cycle, prometheus.GaugeValue, sm.mean_cycle)
  ch <- prometheus.MustNewConstMetric(sc.cycle_per_minute, prometheus.GaugeValue, sm.cycle_per_minute)
  ch <- prometheus.MustNewConstMetric(sc.backfill_last_cycle, prometheus.GaugeValue, sm.backfill_last_cycle)
  ch <- prometheus.MustNewConstMetric(sc.backfill_mean_cycle, prometheus.GaugeValue, sm.backfill_mean_cycle)
  ch <- prometheus.MustNewConstMetric(sc.backfill_depth_mean, prometheus.GaugeValue, sm.backfill_depth_mean)
}

// Returns the Slurm scheduler collector, used to register with the prometheus client
func NewSchedulerCollector() *SchedulerCollector {
  return &SchedulerCollector{
    threads: prometheus.NewDesc(
      "slurm_scheduler_threads",
      "Information provided by the Slurm sdiag command, number of scheduler threads ",
      nil,
      nil),
    queue_size: prometheus.NewDesc(
      "slurm_scheduler_queue_size",
      "Information provided by the Slurm sdiag command, length of the scheduler queue",
      nil,
      nil),
    last_cycle: prometheus.NewDesc(
      "slurm_scheduler_last_cycle",
      "Information provided by the Slurm sdiag command, scheduler last cycle time in (microseconds)",
      nil,
      nil),
    mean_cycle: prometheus.NewDesc(
      "slurm_scheduler_mean_cycle",
      "Information provided by the Slurm sdiag command, scheduler mean cycle time in (microseconds)",
      nil,
      nil),
    cycle_per_minute: prometheus.NewDesc(
      "slurm_scheduler_cycle_per_minute",
      "Information provided by the Slurm sdiag command, number scheduler cycles per minute",
      nil,
      nil),
    backfill_last_cycle: prometheus.NewDesc(
      "slurm_scheduler_backfill_last_cycle",
      "Information provided by the Slurm sdiag command, scheduler backfill last cycle time in (microseconds)",
      nil,
      nil),
    backfill_mean_cycle: prometheus.NewDesc(
      "slurm_scheduler_backfill_mean_cycle",
      "Information provided by the Slurm sdiag command, scheduler backfill mean cycle time in (microseconds)",
      nil,
      nil),
    backfill_depth_mean: prometheus.NewDesc(
      "slurm_scheduler_backfill_depth_mean",
      "Information provided by the Slurm sdiag command, scheduler backfill mean depth",
      nil,
      nil),
  }
}
