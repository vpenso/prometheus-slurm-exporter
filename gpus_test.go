package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseTotalGPUs(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"sinfo_gpus_19.txt", 15}, // slurm version 19
	}
	for _, test := range tests {
		// Read the input data from a file
		file, err := os.Open("test_data/" + test.input)
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		data, err := ioutil.ReadAll(file)
		got := ParseTotalGPUs(data)
		if got != test.want {
			t.Fatalf("got %v; want %v for file %s", got, test.want, test.input)
		}
		t.Logf("%v %+v %f", test, data, got)
	}
}

func TestParseAllocatedGPUs(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"sacct_gpus_19.txt", 12}, // slurm version 19
	}
	for _, test := range tests {
		// Read the input data from a file
		file, err := os.Open("test_data/" + test.input)
		if err != nil {
			t.Fatalf("Can not open test data: %v", err)
		}
		data, err := ioutil.ReadAll(file)
		got := ParseAllocatedGPUs(data)
		if got != test.want {
			t.Fatalf("got %v; want %v for file %s", got, test.want, test.input)
		}
		t.Logf("%v %+v %f", test, data, got)
	}
}
