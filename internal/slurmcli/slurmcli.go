// Package slurmcli wraps invocations of Slurm's CLI tools (sinfo, squeue,
// sacct, ...) using their --json output mode instead of hand-parsed
// positional text fields. Slurm's text output layout is not a stable API
// and has changed between releases (see the project README's notes on
// GPU/GRES accounting breakage); --json is versioned and documented by
// SchedMD as the interface to build tooling against going forward.
package slurmcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Run executes command with args and returns its stdout. Any stderr output
// is captured and folded into the returned error on failure, since the
// legacy text-parsing collectors in this project silently discarded it.
func Run(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("%s: %w: %s", command, err, bytes.TrimSpace(stderr.Bytes()))
		}
		return nil, fmt.Errorf("%s: %w", command, err)
	}
	return out, nil
}

// RunJSON executes command with args (which must request JSON output, e.g.
// by including "--json") and unmarshals its stdout into out.
func RunJSON(command string, args []string, out interface{}) error {
	data, err := Run(command, args)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parsing %s output: %w", command, err)
	}
	return nil
}
