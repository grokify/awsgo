package inspector2util

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/gocharts/v2/data/table"
)

func (fs Findings) ImageVulnerabilitesReporter() ImageVulnerabilitesReporter {
	return ImageVulnerabilitesReporter{Findings: &fs}
}

func (vex ImageVulnerabilitesReporter) TableSet(opts *ReportOptions) (*table.TableSet, error) {
	if vex.Findings == nil {
		return nil, errors.New("findings cannot be nil")
	}
	fs := vex.Findings
	ts := table.NewTableSet("")
	// Sheet for Image Severity Counts
	if tbl, err := fs.TablePivotImagenameSeverityCounts(opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Image Severity Count"
		if err := ts.Add(tbl); err != nil {
			return nil, err
		}
	}
	// Sheet for Image Severity Year Counts
	if tbl, err := fs.TableImagenameSeverityYear(opts); err != nil {
		return nil, err
	} else {
		tbl.Name = "Image Severity Age Counts"
		if err := ts.Add(tbl); err != nil {
			return nil, err
		}
	}
	// Sheet for Images
	if rs, err := fs.ResourceSet([]types.ResourceType{ResourceTypeAwsEcrContainerImage}); err != nil {
		return nil, err
	} else {
		if tbl, err := rs.Table([]string{}, map[int]string{}); err != nil {
			return nil, err
		} else {
			tbl.Name = "Images"
			if err := ts.Add(tbl); err != nil {
				return nil, err
			}
		}
	}
	// Sheet for Vulnerabilities
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

func (vex ImageVulnerabilitesReporter) WriteFileXLSX(filename string, opts *ReportOptions) error {
	if ts, err := vex.TableSet(opts); err != nil {
		return err
	} else {
		return ts.WriteXLSX(filename)
	}
}
