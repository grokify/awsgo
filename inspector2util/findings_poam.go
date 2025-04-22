package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/govex"
	"github.com/grokify/govex/poam"
)

func (fs Findings) POAMTable(opts *govex.ValueOptions, overrides func(field poam.POAMField) (*string, error)) (*table.Table, error) {
	items := []poam.POAMItem{}
	for _, fx := range fs {
		items = append(items, Finding(fx))
	}
	return poam.Table(items, opts, overrides)
}
