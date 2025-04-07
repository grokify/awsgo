package inspector2util

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/stringsutil"
)

func ReadFileListFindingsOutput(filename string) (inspector2.ListFindingsOutput, error) {
	out := inspector2.ListFindingsOutput{}
	if b, err := os.ReadFile(filename); err != nil {
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

// ImageRepositoryNames returns a list of unique repo names
func (fs Findings) ImageRepositoryNames() []string {
	var out []string
	for _, f := range fs {
		fx := Finding(f)
		out = append(out, fx.ImageRepositoryNames()...)
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}

// ImageRepoNameVulnID is used as a unique key across images.
func (fs Findings) ImageRepoNameVulnIDs(sep string) []string {
	var out []string
	for _, f := range fs {
		fx := Finding(f)
		out = append(out, fx.ImageRepoNameVulnIDs(sepFilepathVersion)...)
	}
	sort.Strings(out)
	return out
}

func (fs Findings) Stats() FindingsStats {
	return FindingsStats{Findings: fs}
}
