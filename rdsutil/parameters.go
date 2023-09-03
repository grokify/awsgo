package rdsutil

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/strconv/strconvutil"
)

type Parameters []rds.Parameter

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
			pointer.Dereference(pi.DataType),
			strconvutil.FormatBoolMore(pointer.Dereference(pi.IsModifiable), "Modifiable", "Non Modifiable"),
			pointer.Dereference(pi.DataType),
		})
	}
	return tbl
}
