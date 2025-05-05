package inspector2util

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/govex/severity"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type Finding types.Finding

func (f Finding) FilePathsInclPOMProperties() bool {
	fp := f.FilePaths()
	for _, fpx := range fp {
		if strings.Contains(fpx, filenamePomProperties) {
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

func (f Finding) FindingSeverity(canonical bool) string {
	if !canonical {
		return string(f.Severity)
	} else {
		if can, _, err := severity.ParseSeverity(string(f.Severity)); err != nil {
			return string(f.Severity)
		} else {
			return strings.TrimSpace(can)
		}
	}
}

func (f Finding) FindingOrVendorSeverity(canonical bool) string {
	if sev := f.FindingSeverity(canonical); sev != "" {
		return sev
	} else {
		return f.VendorSeverity(canonical)
	}
}

func (f Finding) ImageHashes() []string {
	var out []string
	for _, res := range f.Resources {
		if res.Details != nil &&
			res.Details.AwsEcrContainerImage != nil &&
			res.Details.AwsEcrContainerImage.ImageHash != nil {
			out = append(out, pointer.Dereference(res.Details.AwsEcrContainerImage.ImageHash))
		}
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}

func (f Finding) ImageRepositoryNames() []string {
	var out []string
	for _, r := range f.Resources {
		if r.Details != nil && r.Details.AwsEcrContainerImage != nil &&
			r.Details.AwsEcrContainerImage.RepositoryName != nil {
			out = append(out, pointer.Dereference(r.Details.AwsEcrContainerImage.RepositoryName))
		}
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}

func (f Finding) ImageTags() []string {
	var out []string
	for _, res := range f.Resources {
		if res.Details != nil && res.Details.AwsEcrContainerImage != nil {
			out = append(out, res.Details.AwsEcrContainerImage.ImageTags...)
		}
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}

func (f Finding) VendorCreatedAt() *time.Time {
	if f.PackageVulnerabilityDetails != nil && f.PackageVulnerabilityDetails.VendorCreatedAt != nil {
		return pointer.Pointer(
			pointer.Dereference(f.PackageVulnerabilityDetails.VendorCreatedAt))
	} else {
		return nil
	}
}

func (f Finding) VendorCreatedAtAgeMonths() *float32 {
	if dt := f.VendorCreatedAt(); dt == nil {
		return nil
	} else {
		return pointer.Pointer(float32(time.Since(*dt)) / float32(timeutil.Day))
	}
}

func (f Finding) VendorSeverity(canonical bool) string {
	if f.PackageVulnerabilityDetails != nil && f.PackageVulnerabilityDetails.VendorSeverity != nil {
		if rawSeverity := strings.TrimSpace(
			pointer.Dereference(
				f.PackageVulnerabilityDetails.VendorSeverity)); rawSeverity == "" {
			return ""
		} else if !canonical {
			return rawSeverity
		} else if canonicalSev, _, err := severity.ParseSeverity(rawSeverity); err != nil {
			return rawSeverity
		} else {
			return canonicalSev
		}
	} else {
		return ""
	}
}

func (f Finding) CVEID() *string {
	if vulnID := f.VulnerabilityID(); vulnID == "" {
		return nil
	} else if vulnIDUpper := strings.ToUpper(vulnID); !strings.HasPrefix(vulnIDUpper, "CVE-") {
		return nil
	} else {
		return &vulnIDUpper
	}
}

func (f Finding) VulnerabilityID() string {
	if f.PackageVulnerabilityDetails != nil && f.PackageVulnerabilityDetails.VulnerabilityId != nil {
		return strings.TrimSpace(pointer.Dereference(f.PackageVulnerabilityDetails.VulnerabilityId))
	} else {
		return ""
	}
}

func (f Finding) VulnerablePackages() Packages {
	if f.PackageVulnerabilityDetails != nil {
		return f.PackageVulnerabilityDetails.VulnerablePackages
	} else {
		return []types.VulnerablePackage{}
	}
}

func (f Finding) VulnerablePackagesX() []Package {
	var out []Package
	if f.PackageVulnerabilityDetails != nil {
		for _, pi := range f.PackageVulnerabilityDetails.VulnerablePackages {
			out = append(out, Package(pi))
		}
	}
	return out
}

func (f Finding) VulnerablePackagesIDs() []string {
	pkgs := f.VulnerablePackages()
	return pkgs.NameAtVersionAtFilepaths()
}

// ImageRepoNameVulnID is used as a unique key across images.
func (f Finding) ImageRepoNameVulnIDs(sep string) []string {
	if sep == "" {
		sep = sepFilepathVersion
	}
	names := f.ImageRepositoryNames()
	vulnID := f.VulnerabilityID()
	var ids []string
	for _, n := range names {
		p := []string{}
		p = append(p, n)
		p = append(p, vulnID)
		ids = append(ids, strings.Join(p, sep))
	}
	return slicesutil.Dedupe(ids)
}
