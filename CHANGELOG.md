## Changelog

Full commit history per tag: https://github.com/vpenso/prometheus-slurm-exporter/commits/{tag number}

* **0.22**
  - Add `slurm_account_cpus_pending` / `slurm_user_cpus_pending` metrics (pending CPUs per account/user, alongside the existing pending/running/suspended job counts). Ported from the `development` branch's pre-0.21 work, reapplied on top of the JSON-based `accounts.go`/`users.go`.
  - Fix `internal/slurmcli.ParseGresString` to correctly handle Slurm's `gpu:(null):N` GRES form (literal placeholder for "no GPU type configured") - the previous trailing-annotation stripping cut the string at the first `(`, which truncated this case since `(null)` itself contains a `(`.

* **0.21**
  - Migrate sinfo/squeue/sacct-based collectors from hand-parsed positional text output to `--json`, fixing the output-format-drift bug tracked as issue #38. Minimum supported Slurm version is now 23.02.
  - Rewrite GPU/GRES accounting (`gpus.go`) to correctly parse typed GRES entries (`gpu:<type>:<count>`), fixing silent zero-counting on any GRES beyond the untyped legacy form. `slurm_gpus_{alloc,idle,total,utilization}` now carry a `gpu_type` label.
  - Add an optional AMD ROCm GPU telemetry collector (`-rocm-acct`, `-rocm-smi-cmd`) exposing real device utilization/memory/temperature/power via `amd-smi`/`rocm-smi`, as `slurm_rocm_gpu_*` metrics separate from the GRES-based scheduling metrics.
  - `sdiag`/`sshare`-based collectors (scheduler, fair-share) remain on legacy text parsing (documented in-code), and go.mod's Go version floor is bumped from 1.12 to 1.21.

* **0.19**
  - Merge PR#50

* **0.18**
  - Add CPU/Memory info per node (see PR#47)

* **0.17**
  - Add fair share collector

* **0.16**
  - Export more data per account/partition, fix squeue for pending jobs
  - Merge PR#34

* **0.15**
  - CPU allocation status per partition

* **0.14**
  - add stats about jobs per account/per user

* **0.13**
  - Merge pull request #32 from pdtpartners/faster-node-metrics

* **0.12**
  - Merge pull request #30 from omnivector-solutions/add_snap_packaging

* **0.11**
  - Merge PR#29
  - Add more backfill stats (see PR#27)

* **0.10**
  - Scheduler: keep track of the DBD agent queue size

* **0.9**
  - README: update to fix build problem raised with issue #26

* **0.8**
  - Merge pull request #21 from cleargray/command-paths

* **0.7**
  - Update scheduler.go (fix issue #18)

* **0.6**
  - Merge pull request #13 from rug-cit-hpc/master

* **0.5**
  - [BUG]: count all job states (issue #9)

* **0.4**
  - Merge pull request #8 from MatMaul/pending-dep

* **0.3**
  - Fix issue #4

* **0.2**
  - Fix issue #3

* **0.1** 
  - Basic prototype
  - Merge PR#2
  - Add Grafana dashboard
