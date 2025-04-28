package inspector2util

import (
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/govex"
	"github.com/grokify/govex/severity"
)

// Table is used as a unique key across images.
func (fs Findings) TablePackages(cols []string, opts *govex.ValueOptions) (*table.Table, error) {
	t := table.NewTable("")
	if len(cols) == 0 {
		cols = TableColumnsImageVulnerabilityPackages()
	}
	t.Columns = cols
	for _, f := range fs {
		f2 := Finding(f)
		if rows, err := f2.PackageSlices(cols, opts); err != nil {
			return nil, err
		} else {
			t.Rows = append(t.Rows, rows...)
		}
	}
	return &t, nil
}

func (fs Findings) TableImageVulnerabilities(cols []string, fmtMap map[int]string, opts *ReportOptions) (*table.Table, error) {
	if opts == nil {
		opts = &ReportOptions{}
	}
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
		if row, err := f2.VulnerabilityFields(cols, opts.VulnerabilityValueOptions); err != nil {
			return nil, err
		} else {
			t.Rows = append(t.Rows, row)
		}
	}
	if len(opts.ColumnInsertOptions) == 0 {
		return &t, nil
	}
	for _, copts := range opts.ColumnInsertOptions {
		if err := t.ColumnInsert(copts); err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func (fs Findings) TablePivotImagenameSeverityCounts(opts *ReportOptions) (*table.Table, error) {
	if opts == nil {
		opts = &ReportOptions{}
	}
	hsets := fs.HistogramSets()
	hset := hsets.HistogramSetOneTwo()
	hset.BinsOrder = severity.SeveritiesAll()

	tbl, err := hset.TablePivot("Image Severities", "Image Name",
		&histogram.SetTablePivotOpts{
			ColTotalRight:  true,
			RowTotalBottom: true})
	if err != nil {
		return nil, err
	}
	if len(opts.ColumnInsertOptions) == 0 {
		return tbl, nil
	}
	for _, copts := range opts.ColumnInsertOptions {
		if err := tbl.ColumnInsert(copts); err != nil {
			return nil, err
		}
	}
	return tbl, nil
}

func (fs Findings) TableImagenameSeverityYear(opts *ReportOptions) (*table.Table, error) {
	h := histogram.NewHistogram("")
	if opts == nil {
		opts = &ReportOptions{}
	}
	for _, f := range fs {
		f2 := Finding(f)
		m := map[string]string{
			ImageRepositoryName:      f2.MustVulnerabilityField(ImageRepositoryName, "", opts.VulnerabilityValueOptions),
			FindingSeverity:          f2.MustVulnerabilityField(FindingSeverity, "", opts.VulnerabilityValueOptions),
			VulnerabilityCreatedYear: f2.MustVulnerabilityField(VulnerabilityCreatedYear, "", opts.VulnerabilityValueOptions),
		}
		h.AddMap(m, 1)
	}
	tbl, err := h.TableMap([]string{
		ImageRepositoryName,
		FindingSeverity,
		VulnerabilityCreatedYear},
		"vuln_count", nil)
	if err != nil {
		return nil, err
	}
	tbl.FormatMap = map[int]string{
		2: table.FormatInt,
		3: table.FormatInt,
	}
	if len(opts.ColumnInsertOptions) == 0 {
		return tbl, nil
	}
	for _, copts := range opts.ColumnInsertOptions {
		if err := tbl.ColumnInsert(copts); err != nil {
			return nil, err
		}
	}
	return tbl, nil
}

// HistogramSets returns histogram sets using the fields
func (fs Findings) HistogramSets() *histogram.HistogramSets {
	hsets := histogram.NewHistogramSets("")
	for _, f := range fs {
		fx := Finding(f)
		imgNames := fx.ImageRepositoryNames()
		vulnSev := fx.FindingSeverity(true)
		vulnID := fx.VulnerabilityID()
		for _, imgName := range imgNames {
			hsets.Add(imgName, vulnSev, vulnID, 1, true)
		}
	}
	return hsets
}
