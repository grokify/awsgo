package inspector2util

import (
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/pointer"
)

type Finding types.Finding

func (f Finding) FilePathsInclPOMProperties() bool {
	fp := f.FilePaths()
	for _, fpx := range fp {
		if strings.Contains(fpx, "pom.properties") {
			return true
		}
	}
	return false
}

func (f Finding) FilePaths() []string {
	var out []string
	if f.PackageVulnerabilityDetails != nil {
		for _, v := range f.PackageVulnerabilityDetails.VulnerablePackages {
			if v.FilePath != nil {
				if fp := strings.TrimSpace(*v.FilePath); fp != "" {
					out = append(out, fp)
				}
			}
		}
	}
	return out
}

func (f Finding) ImageHashes() []string {
	var hashes []string
	for _, res := range f.Resources {
		if res.Details != nil &&
			res.Details.AwsEcrContainerImage != nil &&
			res.Details.AwsEcrContainerImage.ImageHash != nil {
			hashes = append(hashes, pointer.Dereference(res.Details.AwsEcrContainerImage.ImageHash))
		}
	}
	sort.Strings(hashes)
	return hashes
}

func (f Finding) ImageRepositoryNames() []string {
	var out []string
	for _, r := range f.Resources {
		if r.Details != nil && r.Details.AwsEcrContainerImage != nil &&
			r.Details.AwsEcrContainerImage.RepositoryName != nil {
			out = append(out, pointer.Dereference(r.Details.AwsEcrContainerImage.RepositoryName))
		}
	}
	return out
}

func (f Finding) VendorCreatedAt() *time.Time {
	if f.PackageVulnerabilityDetails != nil && f.PackageVulnerabilityDetails.VendorCreatedAt != nil {
		return pointer.Pointer(
			pointer.Dereference(f.PackageVulnerabilityDetails.VendorCreatedAt))
	} else {
		return nil
	}
}

func (f Finding) VendorSeverity() string {
	if f.PackageVulnerabilityDetails != nil && f.PackageVulnerabilityDetails.VendorSeverity != nil {
		return pointer.Dereference(f.PackageVulnerabilityDetails.VendorSeverity)
	} else {
		return ""
	}
}

func (f Finding) VulnerabilityID() string {
	if f.PackageVulnerabilityDetails != nil && f.PackageVulnerabilityDetails.VulnerabilityId != nil {
		return pointer.Dereference(f.PackageVulnerabilityDetails.VulnerabilityId)
	} else {
		return ""
	}
}

func (f Finding) VulnerablePackages() []types.VulnerablePackage {
	if f.PackageVulnerabilityDetails != nil {
		return f.PackageVulnerabilityDetails.VulnerablePackages
	} else {
		return []types.VulnerablePackage{}
	}
}

// ImageRepoNameVulnID is used as a unique key across images.
func (f Finding) ImageRepoNameVulnIDs(sep string) []string {
	names := f.ImageRepositoryNames()
	vulnID := f.VulnerabilityID()
	var ids []string
	for _, n := range names {
		p := []string{}
		p = append(p, n)
		p = append(p, vulnID)
		ids = append(ids, strings.Join(p, sep))
	}
	return ids
}
