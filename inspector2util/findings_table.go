package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/govex/severity"
)

// Table is used as a unique key across images.
func (fs Findings) TablePackages(cols []string) (*table.Table, error) {
	t := table.NewTable("")
	if len(cols) == 0 {
		cols = TableColumnsImageVulnerabilityPackages()
	}
	t.Columns = cols
	for _, f := range fs {
		f2 := Finding(f)
		if rows, err := f2.PackageSlices(cols); err != nil {
			return nil, err
		} else {
			t.Rows = append(t.Rows, rows...)
		}
	}
	return &t, nil
}

func (fs Findings) TableImageVulnerabilities(cols []string, fmtMap map[int]string, opts *table.ColumnInsertOpts) (*table.Table, error) {
	t := table.NewTable("")
	if len(cols) == 0 {
		cols, fmtMap = TableColumnsImageVulnerabilities()
		t.Columns = cols
		t.FormatMap = fmtMap
	} else {
		t.Columns = cols
		t.FormatMap = fmtMap
	}
	for _, f := range fs {
		f2 := Finding(f)
		if row, err := f2.VulnerabilityFields(cols); err != nil {
			return nil, err
		} else {
			t.Rows = append(t.Rows, row)
		}
	}
	if opts == nil {
		return &t, nil
	} else if err := t.ColumnInsert(*opts); err != nil {
		return nil, err
	} else {
		return &t, nil
	}
}

func (fs Findings) TablePivotImagenameSeverityCounts(opts *table.ColumnInsertOpts) (*table.Table, error) {
	hsets := fs.HistogramSets()
	hset := hsets.HistogramSetOneTwo()
	hset.BinsOrder = severity.SeveritiesAll()
	if tbl, err := hset.TablePivot("Image Severities", "Image Name", false, true, false, true); err != nil {
		return nil, err
	} else if opts == nil {
		return tbl, nil
	} else if err := tbl.ColumnInsert(*opts); err != nil {
		return nil, err
	} else {
		return tbl, nil
	}
}

func (fs Findings) TableImagenameSeverityYear(opts *table.ColumnInsertOpts) (*table.Table, error) {
	h := histogram.NewHistogram("")
	for _, f := range fs {
		f2 := Finding(f)
		m := map[string]string{
			ImageRepositoryName:              f2.MustVulnerabilityField(ImageRepositoryName, ""),
			FindingSeverity:                  f2.MustVulnerabilityField(FindingSeverity, ""),
			VulnerabilityCreatedAgeMonthsInt: f2.MustVulnerabilityField(VulnerabilityCreatedYear, ""),
		}
		h.AddMap(m, 1)
	}
	tbl, err := h.TableMap([]string{ImageRepositoryName, FindingSeverity, VulnerabilityCreatedAgeMonthsInt},
		"vuln_count", nil)
	if err != nil {
		return nil, err
	}
	tbl.FormatMap = map[int]string{2: table.FormatInt}
	if opts == nil {
		return tbl, nil
	} else if err := tbl.ColumnInsert(*opts); err != nil {
		return nil, err
	} else {
		return tbl, nil
	}
}

// HistogramSets returns histogram sets using the fields
func (fs Findings) HistogramSets() *histogram.HistogramSets {
	hsets := histogram.NewHistogramSets("")
	for _, f := range fs {
		fx := Finding(f)
		imgNames := fx.ImageRepositoryNames()
		vulnSev := fx.VendorSeverity(true)
		vulnID := fx.VulnerabilityID()
		for _, imgName := range imgNames {
			hsets.Add(imgName, vulnSev, vulnID, 1, true)
		}
	}
	return hsets
}
