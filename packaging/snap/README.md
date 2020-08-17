# Building the prometheus-slurm-exporter snap
Packaging and delivering the prometheus-slurm-exporter as a snap provides users of prometheus-slurm-exporter
a hardened, streamlined, and idempotent experience when consuming this software. See [snapcraft](https://snapcraft.io/) for more information on snaps. 


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
From the root of this project, build the snap:
```bash
snapcraft --use-lxd
```
Once the snap build has completed, list the current working directory to see the resultant snap artifact.
```bash
$ ls -la *.snap
-rw-r--r-- 1 bdx bdx 5562368 Aug 16 18:19 prometheus-slurm-exporter_0.11-1-g01dd959_amd64.snap
```

### Install locally built snap
```bash
sudo snap install prometheus-slurm-exporter_`git describe --tags`_amd64.snap --classic --dangerous
```
* `--classic` - this snap uses classic confinement to allow it to find the slurm commands in the system.
* `--dangerous` - because we are installing this snap from a local resource and sha can't be verified by the snapstore.

### Verify install
Use `ps` to verify the process is running.
```bash
$ ps aux | grep prometheus | head -1
root     2271391  0.0  0.0 1453596 14012 ?       SLsl 18:32   0:00 /snap/prometheus-slurm-exporter/x1/bin/prometheus-slurm-exporter
```

Use `netstat` to verify that the installed `prometheus-slurm-exporter` snap process is listening on port 8080.
```bash
$ sudo netstat -peanut | grep prometheus
tcp6       0      0 :::8080                 :::*                    LISTEN      0          15042010   2271391/prometheus-slurm-exporter
```

Lastly, curl the metrics endpoint.
```bash
$ curl 127.0.0.1:8080/metrics
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

... 

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
...
```

To uninstall the prometheus-slurm-exporter snap:
```bash
sudo snap remove prometheus-slurm-exporter
```
