package main

import (
    "fmt"
    "strconv"
    "time"
    "bufio"
    "os/exec"
    "strings"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"
)

type JobIdMetrics struct {
    JobID   string
    JobName string
    User    string
    Elapsed float64
}

func GetJobMetrics() []*JobIdMetrics {
    jobIDs := GetJobIDs()
    var metrics []*JobIdMetrics

    for _, jobID := range jobIDs {
        jobMetrics := ParseJobMetrics(jobID)
        if jobMetrics != nil {
            metrics = append(metrics, jobMetrics)
        }
    }

    if len(metrics) == 0 {
        log.Warn("No job metrics found or retrieved.")
    } else {
        log.Infof("Retrieved %d job metrics.", len(metrics))
    }

    return metrics
}

func GetJobIDs() []string {
    // Calculate the time one hour ago
    oneHourAgoTime := time.Now().Add(-30 * time.Hour)
    currentTime := time.Now()
    log.Infof("Fetching job information from sacct for jobs completed within the last 30 hours...")

    cmd := exec.Command("sacct", "--state=COMPLETED",
        "-S"+oneHourAgoTime.Format("2006-01-02T15:04:05"),
        "-E"+currentTime.Format("2006-01-02T15:04:05"),
        "-X", "-n", "-a",
        "--format=JobID,JobName,User,Elapsed")

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Error("Error creating stdout pipe:", err)
        return nil
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        log.Error("Error creating stderr pipe:", err)
        return nil
    }

    if err := cmd.Start(); err != nil {
        log.Error("Error starting the command:", err)
        return nil
    }

    var completedJobInfo []string

    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        line := scanner.Text()
        completedJobInfo = append(completedJobInfo, line)
    }

    if err := scanner.Err(); err != nil {
        log.Error("Error scanning stdout:", err)
    }

    errScanner := bufio.NewScanner(stderr)
    for errScanner.Scan() {
        errLine := errScanner.Text()
        log.Error("Error from sacct command:", errLine)
    }

    if err := cmd.Wait(); err != nil {
        log.Error("Error waiting for the command to complete:", err)
    }

    log.Infof("Found %d completed jobs within the last 30 hours.", len(completedJobInfo))
    return completedJobInfo
}

func ParseJobMetrics(line string) *JobIdMetrics {
    fields := strings.Fields(line)
    if len(fields) >= 4 {

        elapsed, err := parseElapsedTime(fields[3])
        if err != nil {
            log.Error("Error parsing Elapsed:", err)
            return nil
        }

        return &JobIdMetrics{
            JobID:   fields[0],
            JobName: fields[1],
            User:    fields[2],
            Elapsed: elapsed,
        }
    } else {
        log.Warnf("Ignoring line: %s", line)
        log.Warnf("Fields: %v", fields)
    }
    return nil
}

func parseElapsedTime(elapsedStr string) (float64, error) {
    parts := strings.Split(elapsedStr, ":")
    if len(parts) != 3 {
        return 0, fmt.Errorf("invalid elapsed time format: %s", elapsedStr)
    }

    hours, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, err
    }

    minutes, err := strconv.Atoi(parts[1])
    if err != nil {
        return 0, err
    }

    seconds, err := strconv.Atoi(parts[2])
    if err != nil {
        return 0, err
    }

    elapsedSeconds := float64(hours*3600 + minutes*60 + seconds)
    return elapsedSeconds, nil
}

func NewJobsCollector() *JobsCollector {
    log.Infof("Initializing JobsCollector")
    return &JobsCollector{
        jobInfo: prometheus.NewDesc("slurm_job_info", "Slurm Job Information", []string{"JobID", "JobName", "User"}, nil),
    }
}

type JobsCollector struct {
    jobInfo *prometheus.Desc
}

// Send all metric descriptions
func (jc *JobsCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- jc.jobInfo
}

func (jc *JobsCollector) Collect(ch chan<- prometheus.Metric) {
    log.Infof("Collecting job metrics")

    jobMetrics := GetJobMetrics()
    for _, metric := range jobMetrics {
        if metric == nil {
            log.Warnf("Skipping nil metric")
            continue
        }
        ch <- prometheus.MustNewConstMetric(jc.jobInfo, prometheus.GaugeValue, metric.Elapsed, metric.JobID, metric.JobName, metric.User)

        log.Infof("Exported metrics for JobID: %s, JobName: %s, User: %s, Elapsed: %f", metric.JobID, metric.JobName, metric.User, metric.Elapsed)
    }
}
