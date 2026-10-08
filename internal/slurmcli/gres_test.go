package slurmcli

import "testing"

func TestParseGresString(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []GresEntry
	}{
		{"empty", "", nil},
		{"null", "(null)", nil},
		{"untyped legacy", "gpu:2", []GresEntry{{Kind: "gpu", Count: 2}}},
		{"typed amd", "gpu:mi250:8", []GresEntry{{Kind: "gpu", Type: "mi250", Count: 8}}},
		{"typed nvidia with idx", "gpu:a100:4(IDX:0-3)", []GresEntry{{Kind: "gpu", Type: "a100", Count: 4}}},
		{
			"multiple typed entries",
			"gpu:mi250:8,gpu:mi300x:2",
			[]GresEntry{{Kind: "gpu", Type: "mi250", Count: 8}, {Kind: "gpu", Type: "mi300x", Count: 2}},
		},
		{"socket affinity suffix", "gpu:4(S:0-1)", []GresEntry{{Kind: "gpu", Count: 4}}},
		{"null type placeholder with idx suffix", "gpu:(null):3(IDX:0-7)", []GresEntry{{Kind: "gpu", Count: 3}}},
		{
			"mixed typed entries with idx suffix",
			"gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)",
			[]GresEntry{{Kind: "gpu", Type: "A30", Count: 4}, {Kind: "gpu", Type: "Q6K", Count: 4}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ParseGresString(c.in)
			if len(got) != len(c.want) {
				t.Fatalf("ParseGresString(%q) = %+v, want %+v", c.in, got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("entry %d: got %+v, want %+v", i, got[i], c.want[i])
				}
			}
		})
	}
}
