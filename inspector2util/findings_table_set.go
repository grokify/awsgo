package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/table"
)

func (fs Findings) TableSetVulnerabilities(opts *ReportOptions) (*table.TableSet, error) {
	ts := table.NewTableSet("")
	if tbl, err := fs.TablePivotImagenameSeverityCounts(opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Image Severity Count"
		if err := ts.Add(tbl); err != nil {
			return nil, err
		}
	}
	if tbl, err := fs.TableImagenameSeverityYear(opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Image Severity Age Counts"
		if err := ts.Add(tbl); err != nil {
			return nil, err
		}
	}
	if tbl, err := fs.TableImageVulnerabilities([]string{}, map[int]string{}, opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Vulnerabilities"
		if err := ts.Add(tbl); err != nil {
			return nil, err
		}
	}
	return ts, nil
}
