// Package smicli wraps invocations of AMD's GPU management CLI (amd-smi,
// or the older rocm-smi) to collect node-local ROCm device telemetry. It is
// deliberately separate from internal/slurmcli: this data has nothing to
// do with Slurm's own scheduling/GRES accounting, it is raw hardware
// telemetry read directly off the node's GPUs.
package smicli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// RunJSON executes command with args and unmarshals its stdout into out.
func RunJSON(command string, args []string, out interface{}) error {
	cmd := exec.Command(command, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	data, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%s: %w: %s", command, err, bytes.TrimSpace(stderr.Bytes()))
		}
		return fmt.Errorf("%s: %w", command, err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parsing %s output: %w", command, err)
	}
	return nil
}
