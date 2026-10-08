## Description

<!-- What does this PR change, and why? Link issues as fixes #NNN. -->

## Checklist

- [ ] Functional change only: no unrelated reformatting of lines this PR does
      not need to change (the `diff-noise-gate` CI job fails cosmetic-dominated
      diffs during current limited-maintenance mode)
- [ ] `gofmt` run with the Go version pinned in `go.mod` — not a newer toolchain
- [ ] `go build ./... && go vet ./... && go test ./...` pass (the test suite is
      fixture-based and needs no live Slurm install)
- [ ] Any new/renamed metrics or label changes are listed below and called out
      as breaking changes (this repo's Grafana dashboards and users' alert
      rules depend on stable metric names and label sets)

<!-- This project is currently maintained on a limited basis: responses and
merges may be slow, and PRs that cannot be reviewed without ongoing capacity
(large refactors, new subsystems) may be closed with rationale. -->
