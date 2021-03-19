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
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

func init() {
	// Metrics have to be registered to be exposed
	prometheus.MustRegister(NewAccountsCollector())   // from accounts.go
	prometheus.MustRegister(NewCPUsCollector())       // from cpus.go
	prometheus.MustRegister(NewGPUsCollector())       // from gpus.go
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

func main() {
	flag.Parse()
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	log.Infof("Starting Server: %s", *listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
