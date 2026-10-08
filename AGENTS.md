# AGENTS.md

Guidance for AI coding agents (and humans) working on this repository.
Read this before proposing changes. Keep changes small, targeted, and
verified; broad rewrites are not wanted.

## Hard constraints (violating these wastes everyone's time)

- **Do not bump `github.com/prometheus/common` past v0.26.0** (nor
  `client_golang` past v1.11.1). Collectors import `prometheus/common/log`,
  which newer versions drop. Rationale and migration path:
  DEVELOPMENT.md "Known limitation: logging dependency pin".
- **Do not "migrate" `sshare.go` or `scheduler.go` to `--json`.** `sshare`
  and `sdiag` have no JSON mode; their text parsing is deliberate and
  documented in the file headers. Do not "fix" the tolerant
  skip-malformed-lines behavior in `ParseFairShareMetrics`.
- **No drive-by reformatting or comment churn.** CI rejects
  cosmetic-dominated diffs (`.github/workflows/diff-noise-gate.yml`).
  Keep diffs minimal and focused; `gofmt -l .` must stay clean.
- **Metric compatibility is API.** Do not rename metrics, change label sets,
  or change which series are emitted without calling it out explicitly in
  the PR description and CHANGELOG.md.

## Verification (there is no live Slurm cluster)

- `go build ./... && go vet ./... && go test ./...` must pass (this is CI).
- Parser changes need fixtures in `test_data/` + table-style tests
  (see `sshare_test.go` for the pattern).
- End-to-end smoke test without a cluster: put stub executables named
  `sinfo`/`squeue`/`sacct`/`sdiag`/`sshare` on `PATH` that `cat` fixtures,
  run the binary with `-listen-address 127.0.0.1:PORT`, curl `/metrics`.

## Release flow (maintainers only)

Update CHANGELOG.md, commit, lightweight tag matching the version
(`git tag 0.24`), `gh release create`, then manually run the
"Release Prometheus Slurm Exporter RPM" workflow with the tag.

## References

- DEVELOPMENT.md: build instructions, dependency pin details
- README.md: metrics documented per collector
- CHANGELOG.md: user-visible changes per release
