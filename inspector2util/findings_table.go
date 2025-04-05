package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/table"
)

// Table is used as a unique key across images.
func (fs Findings) TablePackages(cols []string) *table.Table {
	t := table.NewTable((""))
	if len(cols) == 0 {
		cols = TableColumnsImagePackages()
	}
	t.Columns = cols
	for _, f := range fs {
		f2 := Finding(f)
		rows := f2.PackageSlices(cols)
		t.Rows = append(t.Rows, rows...)
	}
	return &t
}
