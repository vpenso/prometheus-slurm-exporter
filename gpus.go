/* Copyright 2022 Iztok Lebar Bajec

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
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type GPUsMetrics struct {
	alloc         float64
	idle          float64
	other         float64
	total         float64
	utilization   float64
	UserGPUsDCGM  map[string]float64
	UserGPUsSLURM map[string]float64
}

func GPUsGetMetrics() *GPUsMetrics {
	return ParseGPUsMetrics()
}

/* TODO:
  sinfo has gresUSED since slurm>=19.05.0rc01 https://github.com/SchedMD/slurm/blob/master/NEWS
  revert to old process on slurm<19.05.0rc01
  --format=AllocGRES will return gres/gpu=8
  --format=AllocTRES will return billing=16,cpu=16,gres/gpu=8,mem=256G,node=1
func ParseAllocatedGPUs() float64 {
	var num_gpus = 0.0

	args := []string{"-a", "-X", "--format=Allocgres", "--state=RUNNING", "--noheader", "--parsable2"}
	output := string(Execute("sacct", args))
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			if len(line) > 0 {
				line = strings.Trim(line, "\"")
				descriptor := strings.TrimPrefix(line, "gpu:")
				job_gpus, _ := strconv.ParseFloat(descriptor, 64)
				num_gpus += job_gpus
			}
		}
	}

	return num_gpus
}
*/

func ParseAllocatedGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,GresUsed:" --state=allocated
	// 3 gpu:2                                       # slurm>=20.11.8
	// 1 gpu:(null):3(IDX:0-7)                       # slurm 21.08.5
	// 13 gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)      # slurm 21.08.5	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(data) > 0 {
		for _, line := range strings.Split(string(data), "\n") {
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_active_gpus := strings.Fields(line)[1]
				num_node_active_gpus := 0.0
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_active_gpus
			}
		}
	}
	return num_gpus
}

func ParseIdleGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,Gres: ,GresUsed:" --state=idle,allocated
	// 3 gpu:4 gpu:2                                       																# slurm 20.11.8
	// 1 gpu:8(S:0-1) gpu:(null):3(IDX:0-7)                       												# slurm 21.08.5
	// 13 gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1) gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(data) > 0 {
		for _, line := range strings.Split(string(data), "\n") {
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_node_active_gpus := 0.0
				node_active_gpus := strings.Fields(line)[2]
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * (num_node_gpus - num_node_active_gpus)
			}
		}
	}
	return num_gpus
}

func ParseTotalGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,Gres:"
	// 3 gpu:4                                       	# slurm 20.11.8
	// 1 gpu:8(S:0-1)                                	# slurm 21.08.5
	// 13 gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1)        	# slurm 21.08.5
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(data) > 0 {
		for _, line := range strings.Split(string(data), "\n") {
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_gpus
			}
		}
	}
	return num_gpus
}

func ParseGPUsMetrics() *GPUsMetrics {
	var gm GPUsMetrics
	total_gpus := ParseTotalGPUs(TotalGPUsData())
	allocated_gpus := ParseAllocatedGPUs(AllocatedGPUsData())
	idle_gpus := ParseIdleGPUs(IdleGPUsData())
	other_gpus := total_gpus - allocated_gpus - idle_gpus
	gm.alloc = allocated_gpus
	gm.idle = idle_gpus
	gm.other = other_gpus
	gm.total = total_gpus
	gm.utilization = allocated_gpus / total_gpus
	gm.UserGPUsDCGM = ParseUserGPUsDCGM()
	gm.UserGPUsSLURM = ParseUserGPUsSLURM()
	return &gm
}

func AllocatedGPUsData() []byte {
	args := []string{"-a", "-h", "--Format=Nodes: ,GresUsed:", "--state=allocated"}
	return Execute("sinfo", args)
}

func IdleGPUsData() []byte {
	args := []string{"-a", "-h", "--Format=Nodes: ,Gres: ,GresUsed:", "--state=idle,allocated"}
	return Execute("sinfo", args)
}

func TotalGPUsData() []byte {
	args := []string{"-a", "-h", "--Format=Nodes: ,Gres:"}
	return Execute("sinfo", args)
}

func Execute(command string, arguments []string) []byte {
	cmd := exec.Command(command, arguments...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

 func ParseUserGPUsDCGM() map[string]float64 {
    userGPUs := make(map[string]float64)
    
    // Execute a command to get GPU usage information
    cmd := exec.Command("dcgmi", "dmon", "-c", "1", "-e", "203")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("Failed to execute dcgmi command: %v", err)
    }
    
    // Parse the command output
    lines := strings.Split(string(out), "\n")
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) >= 3 && fields[0] != "#" {
            user := fields[1]
            gpuUsage, err := strconv.ParseFloat(fields[2], 64)
            if err != nil {
                log.Errorf("Failed to parse GPU usage for user %s: %v", user, err)
                continue
            }
            userGPUs[user] = gpuUsage
        }
    }
    
    return userGPUs
}


func ParseUserGPUsSLURM() map[string]float64 {
    userGPUs := make(map[string]float64)
    
    // Execute a command to get GPU usage information
    cmd := exec.Command("sacct", "-P", "--format=User,GresUsed", "-a")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("Failed to execute sacct command: %v", err)
    }
    
    // Parse the command output
    lines := strings.Split(string(out), "\n")
    for _, line := range lines {
        fields := strings.Split(line, "|")
        if len(fields) >= 2 {
            user := fields[0]
            gpuUsage, err := strconv.ParseFloat(fields[1], 64)
            if err != nil {
                log.Errorf("Failed to parse GPU usage for user %s: %v", user, err)
                continue
            }
            userGPUs[user] = gpuUsage
        }
    }
    
    return userGPUs
}

type GPUsCollector struct {
    alloc       *prometheus.Desc
    idle        *prometheus.Desc
    other       *prometheus.Desc
    total       *prometheus.Desc
    utilization *prometheus.Desc
    userGPUsDCGM *prometheus.Desc
    userGPUsSLURM *prometheus.Desc
}

func NewGPUsCollector() *GPUsCollector {
    return &GPUsCollector{
        alloc:       prometheus.NewDesc("slurm_gpus_alloc", "Allocated GPUs", nil, nil),
        idle:        prometheus.NewDesc("slurm_gpus_idle", "Idle GPUs", nil, nil),
        other:       prometheus.NewDesc("slurm_gpus_other", "Other GPUs", nil, nil),
        total:       prometheus.NewDesc("slurm_gpus_total", "Total GPUs", nil, nil),
        utilization: prometheus.NewDesc("slurm_gpus_utilization", "Total GPU utilization", nil, nil),
        userGPUsDCGM: prometheus.NewDesc("slurm_user_gpus_dcgm", "Number of GPUs used per user over time, obtained using DCGM and Linux tools", []string{"user"}, nil),
        userGPUsSLURM: prometheus.NewDesc("slurm_user_gpus_slurm", "Number of GPUs used per user over time, obtained using SLURM commands", []string{"user"}, nil),
    }
}

func (cc *GPUsCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- cc.alloc
    ch <- cc.idle
    ch <- cc.other
    ch <- cc.total
    ch <- cc.utilization
    ch <- cc.userGPUsDCGM
    ch <- cc.userGPUsSLURM
}

func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) {
    cm := GPUsGetMetrics()
    ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
    ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
    ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
    ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
    ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, cm.utilization)
    for user, gpus := range cm.UserGPUsDCGM {
        ch <- prometheus.MustNewConstMetric(cc.userGPUsDCGM, prometheus.GaugeValue, gpus, user)
    }
    for user, gpus := range cm.UserGPUsSLURM {
        ch <- prometheus.MustNewConstMetric(cc.userGPUsSLURM, prometheus.GaugeValue, gpus, user)
    }
}