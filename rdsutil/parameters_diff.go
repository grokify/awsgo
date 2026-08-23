package rdsutil

import "github.com/grokify/mogo/pointer"

// ParameterValueDiff holds the values for a parameter that differs, or is
// present in only one of the two parameter maps, compared by DiffParameterValues.
type ParameterValueDiff struct {
	Value1 *string
	Value2 *string
}

// ParameterValueDiffs is a map of parameter name to ParameterValueDiff.
type ParameterValueDiffs map[string]ParameterValueDiff

// DiffParameterValues compares two parameter name/value maps, such as those
// returned by ParametersSet.Map(), and returns a map of parameters that
// differ in value or are present in only one of the two maps.
func DiffParameterValues(params1, params2 map[string]string) ParameterValueDiffs {
	diffs := ParameterValueDiffs{}
	for name1, val1 := range params1 {
		if val2, ok := params2[name1]; ok {
			if val1 != val2 {
				diffs[name1] = ParameterValueDiff{Value1: pointer.Pointer(val1), Value2: pointer.Pointer(val2)}
			}
		} else {
			diffs[name1] = ParameterValueDiff{Value1: pointer.Pointer(val1)}
		}
	}
	for name2, val2 := range params2 {
		if _, ok := params1[name2]; !ok {
			diffs[name2] = ParameterValueDiff{Value2: pointer.Pointer(val2)}
		}
	}
	return diffs
}
