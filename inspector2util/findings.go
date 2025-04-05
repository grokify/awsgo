package inspector2util

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

func ReadFileListFindingsOutput(filename string) (inspector2.ListFindingsOutput, error) {
	out := inspector2.ListFindingsOutput{}
	b, err := os.ReadFile(filename)
	if err != nil {
		return out, err
	} else if err := json.Unmarshal(b, &out); err != nil {
		return out, err
	} else {
		return out, nil
	}
}

type Findings []types.Finding

func (fs Findings) FilterImageHashes(hashesIncl []string) Findings {
	var out Findings
	hashesInclMap := map[string]int{}
	for _, h := range hashesIncl {
		hashesInclMap[h]++
	}
	for _, f := range fs {
		fx := Finding(f)
		imgHashes := fx.ImageHashes()
		for _, h := range imgHashes {
			if _, ok := hashesInclMap[h]; ok {
				out = append(out, f)
				break
			}
		}
	}
	return out
}

func (fs Findings) FilterPOMPropertiesExcl() Findings {
	var out Findings
	for _, f := range fs {
		fx := Finding(f)
		if !fx.FilePathsInclPOMProperties() {
			out = append(out, f)
		}
	}
	return out
}

func (fs Findings) FindingOneRawMatch(s string) *Finding {
	for _, f := range fs {
		b, err := json.Marshal(f)
		if err != nil {
			continue
		}
		if strings.Index(string(b), s) > 0 {
			f2 := Finding(f)
			return &f2
		}
	}
	return nil
}

func (fs Findings) VendorCreatedAtMonthly() map[string]int {
	out := map[string]int{}
	for _, f := range fs {
		if dt := Finding(f).VendorCreatedAt(); dt != nil {
			out[dt.Format("2006-01")]++
		} else {
			out[""]++
		}
	}
	return addTotal(out)
}

func (fs Findings) VendorCreatedAtYearly() map[string]int {
	out := map[string]int{}
	for _, f := range fs {
		if dt := Finding(f).VendorCreatedAt(); dt != nil {
			out[dt.Format("2006")]++
		} else {
			out[""]++
		}
	}
	return addTotal(out)
}

func (fs Findings) VendorSeverities() map[string]int {
	out := map[string]int{}
	for _, f := range fs {
		out[CanonicalSeverity(Finding(f).VendorSeverity())]++
	}
	return addTotal(out)
}

func addTotal(m map[string]int) map[string]int {
	totalCount := 0
	for _, v := range m {
		totalCount += v
	}
	m["_total"] = totalCount
	return m
}

// ImageRepoNameVulnID is used as a unique key across images.
func (fs Findings) ImageRepoNameVulnIDs(sep string) []string {
	var out []string
	for _, f := range fs {
		fx := Finding(f)
		out = append(out, fx.ImageRepoNameVulnIDs(":")...)
	}
	sort.Strings(out)
	return out
}
