package slurmcli

import "strings"

// TresEntry is one Trackable RESource entry, as returned in a sacct job's
// tres.allocated/tres.requested arrays. GPU GRES allocations appear here as
// entries with Type "gres" and Name "gpu:<gputype>" (e.g. "gpu:mi250").
type TresEntry struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// SacctJob is one job record from `sacct --json`.
type SacctJob struct {
	JobID int64    `json:"job_id"`
	State []string `json:"state"`
	Tres  struct {
		Allocated []TresEntry `json:"allocated"`
	} `json:"tres"`
}

// SacctResponse is the top-level shape of `sacct --json`.
type SacctResponse struct {
	Jobs []SacctJob `json:"jobs"`
}

// AllocatedGres returns the job's allocated GPU GRES, parsed from its
// tres.allocated entries (name "gpu:<type>", or the untyped "gpu").
func (j SacctJob) AllocatedGres() []GresEntry {
	var entries []GresEntry
	for _, t := range j.Tres.Allocated {
		if t.Type != "gres" {
			continue
		}
		kind, typ := t.Name, ""
		if i := strings.IndexByte(t.Name, ':'); i >= 0 {
			kind, typ = t.Name[:i], t.Name[i+1:]
		}
		if kind != "gpu" {
			continue
		}
		entries = append(entries, GresEntry{Kind: kind, Type: typ, Count: t.Count})
	}
	return entries
}
