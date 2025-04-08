package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/govex"
)

type ReportOptions struct {
	ColumnInsertOptions       []table.ColumnInsertOpts
	VulnerabilityValueOptions *govex.ValueOpts
}
