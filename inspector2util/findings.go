package inspector2util

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/strslices"
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

// ImageHashes returns a list of unique repo names
func (fs Findings) ImageHashes() []string {
	var out []string
	for _, f := range fs {
		fx := Finding(f)
		out = append(out, fx.ImageHashes()...)
	}
	return strslices.CondenseSpace(out, true, true)
}

// ImageRepositoryNames returns a list of unique repo names
func (fs Findings) ImageRepositoryNames() []string {
	var out []string
	for _, f := range fs {
		fx := Finding(f)
		out = append(out, fx.ImageRepositoryNames()...)
	}
	return strslices.CondenseSpace(out, true, true)
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

func (fs Findings) ImageSet(hashesIncl []string) (*ImageSet, error) {
	if rs, err := fs.ResourceSet([]types.ResourceType{ResourceTypeAwsEcrContainerImage}); err != nil {
		return nil, err
	} else {
		return rs.ImageSet(hashesIncl)
	}
}

func (fs Findings) ResourceSet(inclResourceTypes []types.ResourceType) (*ResourceSet, error) {
	rs := NewResourceSet()
	for _, f := range fs {
		for _, r := range f.Resources {
			if len(inclResourceTypes) == 0 {
				if err := rs.Add(Resource(r)); err != nil {
					return rs, err
				}
			} else {
				for _, rt := range inclResourceTypes {
					if r.Type == rt {
						if err := rs.Add(Resource(r)); err != nil {
							return rs, err
						} else {
							break
						}
					}
				}
			}
		}
	}
	return rs, nil
}

func (fs Findings) Stats() FindingsStats {
	return FindingsStats{Findings: fs}
}
