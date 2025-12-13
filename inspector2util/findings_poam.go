package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/govex"
	"github.com/grokify/govex/reports/poam"
)

func (fs Findings) POAMItems() []poam.POAMItem {
	var items []poam.POAMItem
	for _, fx := range fs {
		items = append(items, Finding(fx))
	}
	return items
}

func (fs Findings) POAMTable(opts *govex.ValueOptions, overrides func(field poam.POAMField) (*string, error)) (*table.Table, error) {
	return poam.Table(fs.POAMItems(), opts, overrides)
}

func (vex ImageVulnerabilitesReporter) POAMWriteFile(filename string, opts *ReportOptions) error {
	if ts, err := vex.TableSet(opts); err != nil {
		return err
	} else {
		return ts.WriteXLSX(filename)
	}
}
