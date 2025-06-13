package main

import (
	"bufio"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	layoutISO = "2006-01-02T15:04:05" // Slurm’s %V default format
)

type PartitionWaitCollector struct {
	desc    *prometheus.Desc
	timeout time.Duration
	clock   func() time.Time
}

func NewPartitionWaitCollector() prometheus.Collector {
	return &PartitionWaitCollector{
		desc: prometheus.NewDesc(
			"slurm_partition_pending_wait_seconds_avg",
			"Average wait (seconds) of currently PENDING jobs, grouped by partition.",
			[]string{"partition"}, nil,
		),
		timeout: 5 * time.Second,
		clock:   time.Now,
	}
}

func (c *PartitionWaitCollector) Describe(ch chan<- *prometheus.Desc) { ch <- c.desc }

func (c *PartitionWaitCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "squeue",
		"-h",              // no header
		"--states=PD",     // pending only
		"--format=%P|%V", // partition | submit-time
	)

	out, err := cmd.Output()
	if err != nil {
		log.Errorf("partition_wait_collector: squeue failed: %v", err)
		return
	}

	type agg struct{ sum float64; n uint64 }
	partMap := make(map[string]agg)

	now := c.clock()

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			log.Warnf("partition_wait_collector: bad line %q", line)
			continue
		}
		part, ts := parts[0], parts[1]

		submit, err := time.ParseInLocation(layoutISO, ts, time.Local)
		if err != nil {
			log.Warnf("partition_wait_collector: cannot parse time %q: %v", ts, err)
			continue
		}
		wait := now.Sub(submit).Seconds()
		if wait < 0 { // clock skew safety
			wait = 0
		}

		a := partMap[part]
		a.sum += wait
		a.n++
		partMap[part] = a
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("partition_wait_collector: scanner: %v", err)
	}

	for part, a := range partMap {
		if a.n == 0 {
			continue
		}
		avg := a.sum / float64(a.n)
		metric, err := prometheus.NewConstMetric(c.desc, prometheus.GaugeValue, avg, part)
		if err == nil {
			ch <- metric
		}
	}
}
