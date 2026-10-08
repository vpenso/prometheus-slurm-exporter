## Changelog

Full commit history per tag: https://github.com/vpenso/prometheus-slurm-exporter/commits/{tag number}

* **0.24**
  - Fix #93: nested Slurm accounts were silently dropped from fairshare metrics (`ParseFairShareMetrics` discarded every indented `sshare` line). The parser now reconstructs the account tree from `sshare` indentation (tolerant of 1-space, 2-space or tab rendering) and exports depth >= 2 accounts on the new `slurm_subaccount_fairshare{account,parent_account,account_depth}` metric. `slurm_account_fairshare` keeps its exact historical series (root + top-level, single `account` label), so existing dashboards and recording rules are unaffected.
  - Behavior change (fair_tree clusters): accounts whose `fairshare` field `sshare` renders blank (fair_tree algorithm, Slurm default since 19.05) are no longer exported as misleading `0` values; they are omitted. Note that on multi-space-indent clusters, top-level accounts previously dropped by the old width-based check now correctly appear in `slurm_account_fairshare`.
  - Behavior change: a failing `sshare` invocation no longer terminates the exporter (was `log.Fatal`); the fairshare metrics go absent and the failure is logged. Alert on `absent(slurm_account_fairshare)` to detect it.
  - Add RockyLinux 9 RPM packaging (spec + rpmbuild tooling) and a manual release workflow that attaches the built RPM to the GitHub Release (ported from PR #110, Go toolchain updated to a current 1.24.x toolchain).
  - Document the systemd service environment and the `-listen-address` override (issue #95).

* **0.23**
  - New metrics, ported from long-open community PRs (authors credited in commit trailers): `slurm_node_gpu_alloc` / `slurm_node_gpu_total` per node and `gpu_type` (#57, #83, #96), `slurm_user_mem_running` in MB (#37, #96), `slurm_user_jobs_running_host{user,host}` (#66), `slurm_partition_jobs_running` (#54), `slurm_nodes_idle_power_save` for idle~ nodes (#52).
  - Run `sinfo` with `-a` so unprivileged exporter accounts see hidden partitions and their nodes (#94).
  - Breaking change (label sets only): `slurm_node_*` CPU/memory series now carry an additional `partition` label (node partitions joined with commas); queries grouping by all labels see the extra dimension.
  - Multi-stage `Dockerfile` (golang builder + slurm-wlm runtime, runs as `nobody`; Slurm/Munge files bind-mounted at run time) (#85).
  - Security: bumped client_golang v1.11.1, prometheus/common v0.26.0, logrus v1.9.3, x/sys v0.8.0, protobuf v1.33.0 - all Dependabot alerts on the default branch cleared. The prometheus/common pin is the last version shipping the `log` package the collectors still use; see DEVELOPMENT.md.

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
