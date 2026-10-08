package slurmcli

// SqueueJob is one job record from `squeue --json`. See the field-name
// verification note in sinfo.go - the same caveat applies here.
type SqueueJob struct {
	JobID       int64    `json:"job_id"`
	UserName    string   `json:"user_name"`
	Account     string   `json:"account"`
	Partition   string   `json:"partition"`
	JobState    []string `json:"job_state"` // e.g. ["PENDING"], ["RUNNING"]
	StateReason string   `json:"state_reason"`
	CPUs        int64    `json:"cpus"`
}

// SqueueResponse is the top-level shape of `squeue --json`.
type SqueueResponse struct {
	Jobs []SqueueJob `json:"jobs"`
}

// State returns the job's base state token. job_state is a flag array; the
// base run state is conventionally its first entry.
func (j SqueueJob) State() string {
	if len(j.JobState) == 0 {
		return ""
	}
	return j.JobState[0]
}
