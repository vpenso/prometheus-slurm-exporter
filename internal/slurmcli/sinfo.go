package slurmcli

// SinfoNode is one node record from `sinfo --json`. Field names/types here
// follow Slurm's data_parser JSON schema (the same schema slurmrestd's
// OpenAPI spec documents) as of the 23.02-26.05 release line; verify
// against `sinfo --json | jq '.nodes[0]'` on the target cluster if a field
// is missing or empty after an upgrade, since SchedMD does not treat this
// text-of-a-manpage as authoritative the way it does the OpenAPI spec.
type SinfoNode struct {
	Name        string   `json:"name"`
	State       []string `json:"state"` // e.g. ["IDLE"], ["MIXED","DRAIN"] - a flag set, not a single token
	CPUs        int64    `json:"cpus"`
	AllocCPUs   int64    `json:"alloc_cpus"`
	IdleCPUs    int64    `json:"idle_cpus"`
	RealMemory  int64    `json:"real_memory"`
	AllocMemory int64    `json:"alloc_memory"`
	// Gres/GresUsed remain colon-delimited descriptor strings (e.g.
	// "gpu:mi250:8") even under --json - Slurm does not decompose GRES
	// into a structured array here. Use ParseGresString on these.
	Gres       string   `json:"gres"`
	GresUsed   string   `json:"gres_used"`
	Partitions []string `json:"partitions"`
}

// SinfoResponse is the top-level shape of `sinfo --json`.
type SinfoResponse struct {
	Nodes []SinfoNode `json:"nodes"`
}

// OtherCPUs returns CPUs neither allocated nor idle (reserved, draining,
// etc.), matching the legacy sinfo "%C" alloc/idle/other/total breakdown.
func (n SinfoNode) OtherCPUs() int64 {
	other := n.CPUs - n.AllocCPUs - n.IdleCPUs
	if other < 0 {
		return 0
	}
	return other
}
