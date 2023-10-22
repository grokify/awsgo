package rdsutil

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/strconv/strconvutil"
)

// ParametersResponse represents the response from the AWS CLI utility for reading
// database parameter groups.
type ParametersResponse struct {
	Parameters Parameters
}

func ParametersResponseReadFile(name string) (ParametersResponse, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return ParametersResponse{}, err
	}
	return ParametersResponseReadBytes(b)
}

func ParametersResponseReadBytes(b []byte) (ParametersResponse, error) {
	var params ParametersResponse
	return params, json.Unmarshal(b, &params)
}

type Parameters []rds.Parameter

func (p Parameters) Map() map[string]string {
	m := map[string]string{}
	for _, pi := range p {
		m[pointer.Dereference(pi.ParameterName)] = pointer.Dereference(pi.ParameterValue)
	}
	return m
}

func (p Parameters) Names() []string {
	var names []string
	for _, pi := range p {
		names = append(names, pointer.Dereference(pi.ParameterName))
	}
	sort.Strings(names)
	return names
}

func (p Parameters) ToPointers() []*rds.Parameter {
	ptrs := []*rds.Parameter{}
	for _, px := range p {
		ptrs = append(ptrs, &px)
	}
	return ptrs
}

func (p Parameters) Table() table.Table {
	tbl := table.NewTable("")
	tbl.Columns = []string{
		"Name",
		"Value",
		"Description",
		"Apply type",
		"Data type",
		"Value type",
		"Source",
	}
	for _, pi := range p {
		tbl.Rows = append(tbl.Rows, []string{
			pointer.Dereference(pi.ParameterName),
			pointer.Dereference(pi.ParameterValue),
			pointer.Dereference(pi.Description),
			pointer.Dereference(pi.ApplyType),
			pointer.Dereference(pi.DataType),
			strconvutil.FormatBoolMore(pointer.Dereference(pi.IsModifiable), "Modifiable", "Non Modifiable"),
			pointer.Dereference(pi.Source),
		})
	}
	return tbl
}

func MapToParameters(m map[string]string) Parameters {
	params := Parameters{}
	for k, v := range m {
		params = append(params, rds.Parameter{
			ParameterName:  pointer.Pointer(k),
			ParameterValue: pointer.Pointer(v)})
	}
	return params
}
