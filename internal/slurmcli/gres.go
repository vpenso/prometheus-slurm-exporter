package slurmcli

import (
	"strconv"
	"strings"
)

// GresEntry is one parsed element of a Slurm GRES descriptor string, e.g.
// the "mi250" entry of "gpu:mi250:8,gpu:mi300x:2".
type GresEntry struct {
	Kind  string // e.g. "gpu"
	Type  string // e.g. "mi250", "a100" - empty for untyped legacy GRES ("gpu:2")
	Count int64
}

// ParseGresString parses a Slurm GRES descriptor string such as
// "gpu:mi250:8,gpu:mi300x:2(IDX:0-1)", "gpu:(null):3(IDX:0-7)" (Slurm's
// literal placeholder when no GPU type is configured), or the untyped
// legacy form "gpu:2" into individual entries. It is vendor-agnostic:
// AMD/ROCm and NVIDIA GPUs are both represented by Slurm as
// "gpu:<type>:<count>" GRES entries, only the type token differs (e.g.
// "mi250" vs "a100").
//
// This replaces the historical bug where the exporter assumed GRES always
// had exactly the untyped two-field form and blindly parsed everything
// after a "gpu:" prefix as a single float, which silently produced zero
// counts for any typed GRES entry.
func ParseGresString(s string) []GresEntry {
	s = strings.TrimSpace(s)
	if s == "" || s == "(null)" {
		return nil
	}

	var entries []GresEntry
	for _, raw := range strings.Split(s, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		// Drop a trailing socket/index annotation, e.g. "(IDX:0-2)" or
		// "(S:0-1)". Only strip a "(...)" group anchored at the end of the
		// string - the type field itself can legitimately be the literal
		// "(null)" (Slurm's placeholder for "no GPU type configured"), so a
		// naive "cut at the first '('" would truncate that case.
		if strings.HasSuffix(raw, ")") {
			if idx := strings.LastIndex(raw, "("); idx >= 0 {
				raw = raw[:idx]
			}
		}

		parts := strings.Split(raw, ":")
		var e GresEntry
		switch len(parts) {
		case 3: // kind:type:count
			e.Kind, e.Type = parts[0], parts[1]
			if e.Type == "(null)" {
				e.Type = ""
			}
			e.Count, _ = strconv.ParseInt(parts[2], 10, 64)
		case 2: // kind:count (untyped, legacy)
			e.Kind = parts[0]
			e.Count, _ = strconv.ParseInt(parts[1], 10, 64)
		default:
			continue
		}
		entries = append(entries, e)
	}
	return entries
}
