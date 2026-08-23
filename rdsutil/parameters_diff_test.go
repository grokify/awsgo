package rdsutil

import "testing"

func TestDiffParameterValues(t *testing.T) {
	params1 := map[string]string{
		"same":       "1",
		"differs":    "a",
		"only_in_p1": "x",
	}
	params2 := map[string]string{
		"same":       "1",
		"differs":    "b",
		"only_in_p2": "y",
	}

	got := DiffParameterValues(params1, params2)

	if len(got) != 3 {
		t.Fatalf("len(got) = %d, want 3: %+v", len(got), got)
	}

	if _, ok := got["same"]; ok {
		t.Error(`got["same"] present, want no diff for equal values`)
	}

	if d, ok := got["differs"]; !ok {
		t.Error(`got["differs"] missing, want a diff entry`)
	} else if d.Value1 == nil || *d.Value1 != "a" || d.Value2 == nil || *d.Value2 != "b" {
		t.Errorf(`got["differs"] = %+v, want Value1="a" Value2="b"`, d)
	}

	if d, ok := got["only_in_p1"]; !ok {
		t.Error(`got["only_in_p1"] missing`)
	} else if d.Value1 == nil || *d.Value1 != "x" || d.Value2 != nil {
		t.Errorf(`got["only_in_p1"] = %+v, want Value1="x" Value2=nil`, d)
	}

	if d, ok := got["only_in_p2"]; !ok {
		t.Error(`got["only_in_p2"] missing`)
	} else if d.Value2 == nil || *d.Value2 != "y" || d.Value1 != nil {
		t.Errorf(`got["only_in_p2"] = %+v, want Value2="y" Value1=nil`, d)
	}
}

func TestDiffParameterValues_Empty(t *testing.T) {
	got := DiffParameterValues(map[string]string{}, map[string]string{})
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}
