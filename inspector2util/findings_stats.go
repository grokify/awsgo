package inspector2util

import (
	"strconv"

	"github.com/grokify/mogo/type/maputil"
)

type FindingsStats struct {
	Findings Findings
}

func (stats FindingsStats) VendorCreatedAtMonthly() map[string]int {
	out := map[string]int{}
	for _, f := range stats.Findings {
		if dt := Finding(f).VendorCreatedAt(); dt != nil {
			out[dt.Format("2006-01")]++
		} else {
			out[""]++
		}
	}
	return addTotal(out)
}

func (stats FindingsStats) VendorCreatedAtYearly() map[string]int {
	out := map[string]int{}
	for _, f := range stats.Findings {
		if dt := Finding(f).VendorCreatedAt(); dt != nil {
			out[dt.Format("2006")]++
		} else {
			out[""]++
		}
	}
	return addTotal(out)
}

func (stats FindingsStats) VendorSeverities(canonicalSev bool) map[string]int {
	out := map[string]int{}
	for _, f := range stats.Findings {
		out[Finding(f).VendorSeverity(canonicalSev)]++
	}
	return addTotal(out)
}

func (stats FindingsStats) FindingVulnerablePackageCounts() map[string]int {
	out := map[string]int{}
	for _, f := range stats.Findings {
		out[strconv.Itoa(
			len(
				Finding(f).VulnerablePackages(),
			),
		)]++
	}
	return addTotal(out)
}

// ImageRepoNameVulnID is used as a unique key across images.
func (stats FindingsStats) ImagenameVulnidCounts() map[string]int {
	out := map[string]int{}
	for _, f := range stats.Findings {
		imgnameVulnids := Finding(f).ImageRepoNameVulnIDs(sepFilepathVersion)
		for _, id := range imgnameVulnids {
			out[id]++
		}
	}
	return out
}

func (stats FindingsStats) ImagenameVulnidRevCounts() map[string]int {
	counts := stats.ImagenameVulnidCounts()
	return maputil.MapCompInt[string](counts).ReverseCountsString()
}

func addTotal(m map[string]int) map[string]int {
	totalCount := 0
	for _, v := range m {
		totalCount += v
	}
	m["_total"] = totalCount
	return m
}
