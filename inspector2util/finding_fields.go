package inspector2util

import (
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/pointer"
)

const (
	FindingDescription                = "find_desc"
	FindingImageHash                  = "image_hash"
	FindingImageRepositoryName        = "image_repo_name"
	PackageVulnerabilityCreated       = "vuln_created"
	PackageVulnerabilityID            = "vuln_id"
	PackageVulnerabilitySeverity      = "vuln_severity"
	PackageVulnerabilitySourceURL     = "vuln_url"
	PackageVulnerabilityReferenceURLs = "vuln_ref_urls"
	PackageInfoFilepath               = "pkg_filepath"
	PackageInfoName                   = "pkg_name"
	PackageInfoVersion                = "pkg_version"
	PackageInfoVersionFixed           = "pkg_version_fixed"
)

func TableColumnsImagePackages() []string {
	return []string{
		FindingImageRepositoryName,
		FindingImageHash,
		PackageVulnerabilitySeverity,
		PackageVulnerabilityCreated,
		PackageVulnerabilityID,
		PackageInfoFilepath,
		PackageInfoName,
		PackageInfoVersion,
		PackageInfoVersionFixed,
		FindingDescription,
		PackageVulnerabilitySourceURL,
		PackageVulnerabilityReferenceURLs,
	}
}

func (f Finding) PackageVulnerabilityField(field string) string {
	if f.PackageVulnerabilityDetails == nil {
		return ""
	}
	switch field {
	case PackageVulnerabilityCreated:
		return pointer.Dereference(f.PackageVulnerabilityDetails.VendorCreatedAt).Format(time.RFC3339)
	case PackageVulnerabilityID:
		return pointer.Dereference(f.PackageVulnerabilityDetails.VulnerabilityId)
	case PackageVulnerabilityReferenceURLs:
		return strings.Join(f.PackageVulnerabilityDetails.ReferenceUrls, ", ")
	case PackageVulnerabilitySeverity:
		return pointer.Dereference(f.PackageVulnerabilityDetails.VendorSeverity)
	case PackageVulnerabilitySourceURL:
		return pointer.Dereference(f.PackageVulnerabilityDetails.SourceUrl)
	default:
		return ""
	}
}

func (f Finding) PackageSlices(fields []string) [][]string {
	var rows [][]string
	for _, res := range f.Resources {
		r2 := Resource(res)
		imgHash := r2.ImageHash()
		imgName := r2.ImageRepositoryName()
		for _, pkg := range f.VulnerablePackages() {
			var row []string
			fmtMap := map[int]string{}
			for i, field := range fields {
				switch field {
				case FindingDescription:
					row = append(row, pointer.Dereference(f.Description))
				case FindingImageHash:
					row = append(row, imgHash)
				case FindingImageRepositoryName:
					row = append(row, imgName)
				case PackageVulnerabilityCreated:
					row = append(row, f.PackageVulnerabilityField(PackageVulnerabilityCreated))
				case PackageVulnerabilityID:
					row = append(row, f.PackageVulnerabilityField(PackageVulnerabilityID))
				case PackageVulnerabilitySeverity:
					row = append(row, f.PackageVulnerabilityField(PackageVulnerabilitySeverity))
					fmtMap[i] = table.FormatDate
				case PackageVulnerabilitySourceURL:
					row = append(row, f.PackageVulnerabilityField(PackageVulnerabilitySourceURL))
				case PackageVulnerabilityReferenceURLs:
					row = append(row, f.PackageVulnerabilityField(PackageVulnerabilityReferenceURLs))
				case PackageInfoFilepath:
					row = append(row, pointer.Dereference(pkg.FilePath))
				case PackageInfoName:
					row = append(row, pointer.Dereference(pkg.Name))
				case PackageInfoVersion:
					row = append(row, pointer.Dereference(pkg.Version))
				case PackageInfoVersionFixed:
					row = append(row, pointer.Dereference(pkg.FixedInVersion))
				}
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
	}
	return rows
}
