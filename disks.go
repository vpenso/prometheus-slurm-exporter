package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// JobsMetrics stores metrics for each node
type DiskMetrics struct {
	fsize      float64
	size_avail float64
	size_used  float64
	size       float64
	hostname   string
}

type Jobio struct {
	read     float64
	write    float64
	hostname string
}

func DiskGetMetrics() (map[string]*DiskMetrics, map[string]*Jobio) {
	return ParseDiskMetrics(DiskData())
}

func readProcIO(pid string) string {
	// Формируем путь к файлу /proc/{{pid}}/io
	procIOPath := fmt.Sprintf("/proc/%s/io", pid)

	// Читаем содержимое файла
	data, err := ioutil.ReadFile(procIOPath)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

// ParseNodeMetrics takes the output of sinfo with node data
// It returns a map of metrics per node
func ParseDiskMetrics(input []byte) (map[string]*DiskMetrics, map[string]*Jobio) {
	disk_info := make(map[string]*DiskMetrics)
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		split := strings.Split(line, " ")
		disk_name := strings.Split(split[0], "=")[1]
		disk_name = disk_name[1 : len(disk_name)-1]
		disk_info[disk_name] = &DiskMetrics{0, 0, 0, 0, ""}
		if strings.Split(split[1], "=")[1] != "\"\"" {
			disk_info[disk_name].size_avail, _ = strconv.ParseFloat(strings.Split(split[1], "=")[1][1:len(strings.Split(split[1], "=")[1])-1], 64)
		} else {
			disk_info[disk_name].size_avail = -1
		}
		if strings.Split(split[2], "=")[1] != "\"\"" {
			disk_info[disk_name].fsize, _ = strconv.ParseFloat(strings.Split(split[2], "=")[1][1:len(strings.Split(split[2], "=")[1])-1], 64)
		} else {
			disk_info[disk_name].fsize = -10
		}
		disk_info[disk_name].size, _ = strconv.ParseFloat(strings.Split(split[3], "=")[1][1:len(strings.Split(split[3], "=")[1])-1], 64)
		disk_info[disk_name].size_used = disk_info[disk_name].fsize - disk_info[disk_name].size_avail
		hostname := string(GetHostName())
		if strings.Contains(hostname, ".") {
			hostname = strings.Split(hostname, ".")[0]
		}
		disk_info[disk_name].hostname = hostname

	}

	pids_lines, err := ShowPids()
	jobs_io := make(map[string]*Jobio)
	if err == nil {
		lines = strings.Split(string(pids_lines), "\n")
		lines = lines[1 : len(lines)-1]
		fmt.Println(lines)
		for _, line := range lines {
			split := strings.Fields(line)
			pid_io_lines := strings.Split(readProcIO(split[0]), "\n")
			fmt.Println(pid_io_lines)
			fmt.Println("JOB id is ", split[1])
			if _, exists := jobs_io[split[1]]; !exists {
				// Если ключа нет, создаем новый элемент
				jobs_io[split[1]] = &Jobio{0, 0, ""}
				jobs_io[split[1]].read, _ = strconv.ParseFloat(strings.Split(pid_io_lines[0], " ")[1], 64)
				jobs_io[split[1]].write, _ = strconv.ParseFloat(strings.Split(pid_io_lines[1], " ")[1], 64)
				hostname := string(GetHostName())
				if strings.Contains(hostname, ".") {
					hostname = strings.Split(hostname, ".")[0]
				}
				jobs_io[split[1]].hostname = hostname
			} else {
				// Если ключ уже существует, суммируем значения
				read_sum, _ := strconv.ParseFloat(strings.Split(pid_io_lines[0], " ")[1], 64)
				jobs_io[split[1]].read = jobs_io[split[1]].read + read_sum
				write_sum, _ := strconv.ParseFloat(strings.Split(pid_io_lines[1], " ")[1], 64)
				jobs_io[split[1]].write = jobs_io[split[1]].write + write_sum
			}

		}
	}

	return disk_info, jobs_io
}

// NodeData executes the sinfo command to get data for each node
// It returns the output of the sinfo command
func DiskData() []byte {
	cmd := exec.Command("lsblk", "-Pb", "-o", "NAME,FSAVAIL,FSSIZE,SIZE")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}
func GetHostName() []byte {
	cmd := exec.Command("hostname")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func ShowPids() ([]byte, error) {
	cmd := exec.Command("scontrol", "listpids")
	out, err := cmd.Output()
	return out, err
}

type DiskCollector struct {
	disk_fsize      *prometheus.Desc
	disk_size_avail *prometheus.Desc
	disk_size_used  *prometheus.Desc
	disk_size       *prometheus.Desc
	jobs_read_disk  *prometheus.Desc
	jobs_write_disk *prometheus.Desc
}

// NewNodeCollector creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewDiskCollector() *DiskCollector {
	labels := []string{"DISK", "HOSTNAME"}
	labels2 := []string{"JOBID", "HOSTNAME"}
	return &DiskCollector{
		disk_fsize:      prometheus.NewDesc("slurm_disk_filesystemsize", "DISK INFOf", labels, nil),
		disk_size_avail: prometheus.NewDesc("slurm_disk_size_avail", "DISK INFOd", labels, nil),
		disk_size_used:  prometheus.NewDesc("slurm_disk_size_used", "DISK INFOs", labels, nil),
		disk_size:       prometheus.NewDesc("slurm_disk_size", "DISK INFOf", labels, nil),
		jobs_read_disk:  prometheus.NewDesc("slurm_disk_jobs_read", "DISK INFOf", labels2, nil),
		jobs_write_disk: prometheus.NewDesc("slurm_disk_jobs_write", "DISK INFOf", labels2, nil),
	}
}

// Send all metric descriptions
func (nc *DiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.disk_fsize
	ch <- nc.disk_size_avail
	ch <- nc.disk_size_used
	ch <- nc.disk_size
	ch <- nc.jobs_read_disk
	ch <- nc.jobs_write_disk
}

func (nc *DiskCollector) Collect(ch chan<- prometheus.Metric) {
	disks, jobs_io := DiskGetMetrics()
	for disk := range disks {
		ch <- prometheus.MustNewConstMetric(nc.disk_fsize, prometheus.GaugeValue, disks[disk].fsize, disk, disks[disk].hostname)
		ch <- prometheus.MustNewConstMetric(nc.disk_size, prometheus.GaugeValue, disks[disk].size, disk, disks[disk].hostname)
		ch <- prometheus.MustNewConstMetric(nc.disk_size_avail, prometheus.GaugeValue, disks[disk].size_avail, disk, disks[disk].hostname)
		ch <- prometheus.MustNewConstMetric(nc.disk_size_used, prometheus.GaugeValue, disks[disk].size_used, disk, disks[disk].hostname)
	}
	for job := range jobs_io {
		ch <- prometheus.MustNewConstMetric(nc.jobs_read_disk, prometheus.GaugeValue, jobs_io[job].read, job, jobs_io[job].hostname)
		ch <- prometheus.MustNewConstMetric(nc.jobs_write_disk, prometheus.GaugeValue, jobs_io[job].write, job, jobs_io[job].hostname)
	}
}
