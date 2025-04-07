package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/table"
)

func (fs Findings) TableSetVulnerabilities(opts *table.ColumnInsertOpts) (*table.TableSet, error) {
	ts := table.NewTableSet("")
	if tbl, err := fs.TablePivotImagenameSeverityCounts(opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Image Severity Count"
		ts.Add(tbl)
	}
	if tbl, err := fs.TableImagenameSeverityYear(opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Image Severity Age Counts"
		ts.Add(tbl)
	}
	if tbl, err := fs.TableImageVulnerabilities([]string{}, map[int]string{}, opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Vulnerabilities"
		ts.Add(tbl)
	}
	return ts, nil
}
