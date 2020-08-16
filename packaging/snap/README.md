# Building the prometheus-slurm-exporter snap

### Prereqs
* [snapcraft](https://snapcraft.io)
    ```bash
    sudo snap install snapcraft --classic
    ```
* [lxd](https://linuxcontainers.org/)
    ```bash
    sudo snap install lxd
    ```

### Build
From the root of this project:
```bash
snapcraft --use-lxd
```

### Install locally built snap
```bash
sudo snap install prometheus-slurm-exporter_`git describe --tags`_amd64.snap
```

### Verify install
Curl the metrics endpoint to verify things are working.
```bash
$ curl 127.0.0.1:8080/metrics
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 11
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.14.7"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.639e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 2.639e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 3698
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 1668
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 3.436808e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 2.639e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.2619648e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.899392e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 17262
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.258688e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.651904e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 18930
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 13888
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 89624
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 98304
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.473924e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.771662e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 589824
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 589824
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.243572e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 10
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.22
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1024
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 16
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.179648e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.59760273488e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.33564416e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP slurm_cpus_alloc Allocated CPUs
# TYPE slurm_cpus_alloc gauge
slurm_cpus_alloc 0
# HELP slurm_cpus_idle Idle CPUs
# TYPE slurm_cpus_idle gauge
slurm_cpus_idle 8
# HELP slurm_cpus_other Mix CPUs
# TYPE slurm_cpus_other gauge
slurm_cpus_other 0
# HELP slurm_cpus_total Total CPUs
# TYPE slurm_cpus_total gauge
slurm_cpus_total 8
# HELP slurm_nodes_alloc Allocated nodes
# TYPE slurm_nodes_alloc gauge
slurm_nodes_alloc 0
# HELP slurm_nodes_comp Completing nodes
# TYPE slurm_nodes_comp gauge
slurm_nodes_comp 0
# HELP slurm_nodes_down Down nodes
# TYPE slurm_nodes_down gauge
slurm_nodes_down 0
# HELP slurm_nodes_drain Drain nodes
# TYPE slurm_nodes_drain gauge
slurm_nodes_drain 0
# HELP slurm_nodes_err Error nodes
# TYPE slurm_nodes_err gauge
slurm_nodes_err 0
# HELP slurm_nodes_fail Fail nodes
# TYPE slurm_nodes_fail gauge
slurm_nodes_fail 0
# HELP slurm_nodes_idle Idle nodes
# TYPE slurm_nodes_idle gauge
slurm_nodes_idle 1
# HELP slurm_nodes_maint Maint nodes
# TYPE slurm_nodes_maint gauge
slurm_nodes_maint 0
# HELP slurm_nodes_mix Mix nodes
# TYPE slurm_nodes_mix gauge
slurm_nodes_mix 0
# HELP slurm_nodes_resv Reserved nodes
# TYPE slurm_nodes_resv gauge
slurm_nodes_resv 0
# HELP slurm_queue_cancelled Cancelled jobs in the cluster
# TYPE slurm_queue_cancelled gauge
slurm_queue_cancelled 0
# HELP slurm_queue_completed Completed jobs in the cluster
# TYPE slurm_queue_completed gauge
slurm_queue_completed 0
# HELP slurm_queue_completing Completing jobs in the cluster
# TYPE slurm_queue_completing gauge
slurm_queue_completing 0
# HELP slurm_queue_configuring Configuring jobs in the cluster
# TYPE slurm_queue_configuring gauge
slurm_queue_configuring 0
# HELP slurm_queue_failed Number of failed jobs
# TYPE slurm_queue_failed gauge
slurm_queue_failed 0
# HELP slurm_queue_node_fail Number of jobs stopped due to node fail
# TYPE slurm_queue_node_fail gauge
slurm_queue_node_fail 0
# HELP slurm_queue_pending Pending jobs in queue
# TYPE slurm_queue_pending gauge
slurm_queue_pending 0
# HELP slurm_queue_pending_dependency Pending jobs because of dependency in queue
# TYPE slurm_queue_pending_dependency gauge
slurm_queue_pending_dependency 0
# HELP slurm_queue_preempted Number of preempted jobs
# TYPE slurm_queue_preempted gauge
slurm_queue_preempted 0
# HELP slurm_queue_running Running jobs in the cluster
# TYPE slurm_queue_running gauge
slurm_queue_running 0
# HELP slurm_queue_suspended Suspended jobs in the cluster
# TYPE slurm_queue_suspended gauge
slurm_queue_suspended 0
# HELP slurm_queue_timeout Jobs stopped by timeout
# TYPE slurm_queue_timeout gauge
slurm_queue_timeout 0
# HELP slurm_scheduler_backfill_depth_mean Information provided by the Slurm sdiag command, scheduler backfill mean depth
# TYPE slurm_scheduler_backfill_depth_mean gauge
slurm_scheduler_backfill_depth_mean 0
# HELP slurm_scheduler_backfill_last_cycle Information provided by the Slurm sdiag command, scheduler backfill last cycle time in (microseconds)
# TYPE slurm_scheduler_backfill_last_cycle gauge
slurm_scheduler_backfill_last_cycle 0
# HELP slurm_scheduler_backfill_mean_cycle Information provided by the Slurm sdiag command, scheduler backfill mean cycle time in (microseconds)
# TYPE slurm_scheduler_backfill_mean_cycle gauge
slurm_scheduler_backfill_mean_cycle 481
# HELP slurm_scheduler_backfilled_heterogeneous_total Information provided by the Slurm sdiag command, number of heterogeneous job components started thanks to backfilling since last Slurm start
# TYPE slurm_scheduler_backfilled_heterogeneous_total gauge
slurm_scheduler_backfilled_heterogeneous_total 0
# HELP slurm_scheduler_backfilled_jobs_since_cycle_total Information provided by the Slurm sdiag command, number of jobs started thanks to backfilling since last time stats where reset
# TYPE slurm_scheduler_backfilled_jobs_since_cycle_total gauge
slurm_scheduler_backfilled_jobs_since_cycle_total 0
# HELP slurm_scheduler_backfilled_jobs_since_start_total Information provided by the Slurm sdiag command, number of jobs started thanks to backfilling since last slurm start
# TYPE slurm_scheduler_backfilled_jobs_since_start_total gauge
slurm_scheduler_backfilled_jobs_since_start_total 0
# HELP slurm_scheduler_cycle_per_minute Information provided by the Slurm sdiag command, number scheduler cycles per minute
# TYPE slurm_scheduler_cycle_per_minute gauge
slurm_scheduler_cycle_per_minute 1
# HELP slurm_scheduler_dbd_queue_size Information provided by the Slurm sdiag command, length of the DBD agent queue
# TYPE slurm_scheduler_dbd_queue_size gauge
slurm_scheduler_dbd_queue_size 0
# HELP slurm_scheduler_last_cycle Information provided by the Slurm sdiag command, scheduler last cycle time in (microseconds)
# TYPE slurm_scheduler_last_cycle gauge
slurm_scheduler_last_cycle 40
# HELP slurm_scheduler_mean_cycle Information provided by the Slurm sdiag command, scheduler mean cycle time in (microseconds)
# TYPE slurm_scheduler_mean_cycle gauge
slurm_scheduler_mean_cycle 481
# HELP slurm_scheduler_queue_size Information provided by the Slurm sdiag command, length of the scheduler queue
# TYPE slurm_scheduler_queue_size gauge
slurm_scheduler_queue_size 0
# HELP slurm_scheduler_threads Information provided by the Slurm sdiag command, number of scheduler threads
# TYPE slurm_scheduler_threads gauge
slurm_scheduler_threads 4
```
