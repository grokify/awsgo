package inspector2util

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/encoding/jsonutil/jsonraw"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/stringsutil"
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

// ImageRepoNameVulnID is used as a unique key across images.
func (fs Findings) ImageRepoNameVulnIDsMapSeverity() map[string][]string {
	out := maputil.MapStringSlice{}
	for _, f := range fs {
		fx := Finding(f)
		rvids := fx.ImageRepoNameVulnIDs(sepFilepathVersion)
		fsev := fx.FindingSeverity(true)
		if _, ok := out[fsev]; !ok {
			out[fsev] = []string{}
		}
		out[fsev] = append(out[fsev], rvids...)
	}
	out.Sort(true)
	return out
}

func (fs Findings) ImageSet(hashesIncl []string) (*ImageSet, error) {
	if rs, err := fs.ResourceSet([]types.ResourceType{ResourceTypeAwsEcrContainerImage}); err != nil {
		return nil, err
	} else {
		return rs.ImageSet(hashesIncl)
	}
}

func (fs Findings) ImageSetRepoNameByTagsOrLatest(imageTagOverrides []string) (*ImageSet, error) {
	if is, err := fs.ImageSet([]string{}); err != nil {
		return nil, err
	} else if iss, err := is.ImageSetsByRepoName(); err != nil {
		return nil, err
	} else {
		return iss.ImageSetByRepositoryNameTagsOrLatest(imageTagOverrides)
	}
}

// MergeByImageRepoNameAndVulnID merges select fields including images and packages.
func (fs Findings) MergeByImageRepoNameAndVulnID() (Findings, error) {
	merger := NewFindingMerger()
	if err := merger.Add(fs...); err != nil {
		return Findings{}, err
	} else {
		return merger.Merge()
	}
}

// MergeFilteredByImageRepoNameAndVulnID merges select fields including images and packages.
// Image repo name and vulnid must be the same.
func (fs Findings) MergeFilteredByImageRepoNameAndVulnID() (*types.Finding, error) {
	if len(fs) == 0 {
		return nil, errors.New("no findings to merge")
	} else if len(fs) == 1 {
		f, err := jsonraw.Clone(fs[0])
		return &f, err
	}
	vulnIDs := fs.VulnerabilityIDs()
	if len(vulnIDs) == 0 {
		return nil, errors.New("no vuln ids, requires one")
	} else if len(vulnIDs) > 1 {
		return nil, errors.New("multiple vuln ids, requires one")
	}
	names := fs.ImageRepositoryNames()
	if len(names) == 0 {
		return nil, errors.New("no image repo names, requires one")
	} else if len(names) > 1 {
		return nil, errors.New("multiple image repo names, requires one")
	}
	var merged Finding
	for i, fi := range fs {
		if i == 0 {
			if try, err := jsonraw.Clone(fi); err != nil {
				return nil, err
			} else {
				if try.PackageVulnerabilityDetails == nil {
					try.PackageVulnerabilityDetails = &types.PackageVulnerabilityDetails{}
				}
				merged = Finding(try)
			}
			continue
		}
		fx := Finding(fi)
		imageHashesMerged := merged.ImageHashes()
		for _, ri := range fi.Resources {
			rx := Resource(ri)
			if !slices.Contains(imageHashesMerged, rx.ImageHash()) {
				merged.Resources = append(merged.Resources, ri)
			}
		}
		mergedPkgIDs := merged.VulnerablePackagesIDs()
		ps := fx.VulnerablePackages()
		for _, pi := range ps {
			px := Package(pi)
			if !slices.Contains(mergedPkgIDs, px.NameAtVersionAtFilepath()) {
				merged.PackageVulnerabilityDetails.VulnerablePackages = append(
					merged.PackageVulnerabilityDetails.VulnerablePackages, pi)
			}
		}
	}
	mergedT := types.Finding(merged)
	return &mergedT, nil
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

func (fs Findings) VulnerabilityIDs() []string {
	var out []string
	for _, fi := range fs {
		fx := Finding(fi)
		out = append(out, fx.VulnerabilityID())
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}
