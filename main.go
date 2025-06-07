/* Copyright 2017-2020 Victor Penso, Matteo Dessalvi, Joeri Hermans

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
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func init() {
	// Metrics have to be registered to be exposed
	prometheus.MustRegister(NewAccountsCollector())   // from accounts.go
	prometheus.MustRegister(NewCPUsCollector())       // from cpus.go
	prometheus.MustRegister(NewNodesCollector())      // from nodes.go
	prometheus.MustRegister(NewNodeCollector())       // from node.go
	prometheus.MustRegister(NewPartitionsCollector()) // from partitions.go
	prometheus.MustRegister(NewQueueCollector())      // from queue.go
	prometheus.MustRegister(NewSchedulerCollector())  // from scheduler.go
	prometheus.MustRegister(NewFairShareCollector())  // from sshare.go
	prometheus.MustRegister(NewUsersCollector())      // from users.go
}

var listenAddress = flag.String(
	"listen-address",
	":8080",
	"The address to listen on for HTTP requests.")

var gpuAcct = flag.Bool(
	"gpus-acct",
	false,
	"Enable GPUs accounting")

var npuAcct = flag.Bool(
	"npus-acct",
	false,
	"Enable NPUs accounting")

// Execute the sinfo command and return its output
func Execute(command string, arguments []string) []byte {
	cmd := exec.Command(command, arguments...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

func main() {
	flag.Parse()

	// Turn on GPUs accounting only if the corresponding command line option is set to true.
	if *gpuAcct {
		log.Infof("GPU accounting enabled, detecting GPU resources...")
		
		// Get GPU metrics for logging
		gpuMetrics := GPUsGetMetrics()
		log.Infof("GPU Detection Results:")
		log.Infof("  Total GPUs detected: %.0f", gpuMetrics.total)
		log.Infof("  Allocated GPUs: %.0f", gpuMetrics.alloc)
		log.Infof("  Idle GPUs: %.0f", gpuMetrics.idle)
		log.Infof("  GPU Utilization: %.2f%%", gpuMetrics.utilization*100)
		
		if gpuMetrics.total == 0 {
			log.Warnf("No GPUs detected in SLURM cluster. Please check:")
			log.Warnf("  1. SLURM GRES configuration for GPUs")
			log.Warnf("  2. 'sinfo' command output format")
			log.Warnf("  3. GPU resource naming in SLURM (should be 'gpu:*')")
		}
		
		prometheus.MustRegister(NewGPUsCollector()) // from gpus.go
	}
	
	if *npuAcct {
		log.Infof("NPU accounting enabled, detecting NPU resources...")
		
		// Get NPU metrics for logging
		npuMetrics := NPUsGetMetrics()
		log.Infof("NPU Detection Results:")
		log.Infof("  Total NPUs detected: %.0f", npuMetrics.total)
		log.Infof("  Allocated NPUs: %.0f", npuMetrics.alloc)
		log.Infof("  Idle NPUs: %.0f", npuMetrics.idle)
		log.Infof("  NPU Utilization: %.2f%%", npuMetrics.utilization*100)
		
		if npuMetrics.total == 0 {
			log.Warnf("No NPUs detected in SLURM cluster. Please check:")
			log.Warnf("  1. SLURM GRES configuration for NPUs")
			log.Warnf("  2. 'sinfo' command output format")
			log.Warnf("  3. NPU resource naming in SLURM (should be 'npu:*')")
		}
		
		prometheus.MustRegister(NewNPUsCollector()) // from npus.go
	}
	
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	log.Infof("Starting Server: %s", *listenAddress)
	log.Infof("GPUs Accounting: %t", *gpuAcct)
	log.Infof("NPUs Accounting: %t", *npuAcct)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
