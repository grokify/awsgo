package rdsutil

import "testing"

func TestParametersResponseReadFile(t *testing.T) {
	got, err := ParametersResponseReadFile("testdata/db-perf-group_mysql8.json")
	if err != nil {
		t.Fatalf("ParametersResponseReadFile() error = %v", err)
	}
	if len(got.Parameters) == 0 {
		t.Fatal("len(got.Parameters) = 0, want > 0")
	}
}

func TestParametersResponseReadFile_NotFound(t *testing.T) {
	_, err := ParametersResponseReadFile("testdata/does-not-exist.json")
	if err == nil {
		t.Fatal("ParametersResponseReadFile() error = nil, want error for missing file")
	}
}
