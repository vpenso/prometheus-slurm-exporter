package smicli

// ROCmDevice is one GPU's record from `amd-smi metric --json` (or the
// equivalent rocm-smi --json invocation an operator may configure via
// -rocm-smi-cmd). Field names/nesting here are representative of AMD's
// published schema as of ROCm 6.x; AMD has changed this schema across
// releases, so verify with `amd-smi metric --json | jq '.[0]'` (or the
// rocm-smi equivalent) on the target node before relying on it, and adjust
// this struct - it is the single place that mapping lives.
type ROCmDevice struct {
	GPU   int    `json:"gpu"`
	UUID  string `json:"unique_id"`
	Usage struct {
		GfxActivity float64 `json:"gfx_activity"`
	} `json:"usage"`
	Mem struct {
		Used  uint64 `json:"used"`
		Total uint64 `json:"total"`
	} `json:"mem_usage"`
	Temperature struct {
		Edge     float64 `json:"edge"`
		Junction float64 `json:"junction"`
		Memory   float64 `json:"memory"`
	} `json:"temperature"`
	Power struct {
		Socket float64 `json:"socket"`
	} `json:"power"`
}
