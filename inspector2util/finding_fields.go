package inspector2util

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/govex"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/timeutil"
)

const (
	FindingDescription                      = "find_desc"
	FindingSeverity                         = "find_severity"
	FindingTitle                            = "find_title"
	ImageHash                               = "image_hash"
	ImageRepositoryName                     = "image_repo_name"
	ImageNameVulnerabilityID                = "image_name_vuln_id"
	ImageTags                               = "image_tags"
	VulnerabilityCreated                    = "vuln_created"
	VulnerabilityCreatedYear                = "vuln_created_year"
	VulnerabilityCreatedAgeMonthsInt        = "vuln_created_age_months_int"
	VulnerabilitySLADueDate                 = "vuln_sla_due_date"
	VulnerabilityID                         = "vuln_id"
	VulnerabilitySeverity                   = "vuln_severity"
	VulnerabilitySourceURL                  = "vuln_url"
	VulnerabilityReferenceURLs              = "vuln_ref_urls"
	PackageInfoFilepath                     = "pkg_filepath"
	PackageInfoName                         = "pkg_name"
	PackageInfoVersion                      = "pkg_version"
	PackageInfoVersionFixed                 = "pkg_version_fixed"
	PackagesFilepathsAtVersion              = "pkgs_filepaths_version"
	PackagesFilepathsAtVersionFixed         = "pkgs_filepaths_version_fixed"
	PackagesFilepathsPOMProperites          = "pkgs_filepaths_pom_properties"
	PackagesManagers                        = "pkgs_mgrs"
	PackagesNamesAtVersion                  = "pkgs_names_version"
	PackagesNamesAtVersionsFixed            = "pkgs_names_version_fixed"
	PackagesNamesAndFilepathsAtVersion      = "pkgs_names_filepaths_version"
	PackagesNamesAndFilepathsAtVersionFixed = "pkgs_names_filepaths_version_fixed"
	PackagesTypes                           = "pkgs_types"
)

// TableColumnsImageVulnerabilityPackages returns rows where
// each row is an image+vulnerability+package.
func TableColumnsImageVulnerabilityPackages() []string {
	return []string{
		ImageRepositoryName,
		ImageHash,
		FindingSeverity,
		VulnerabilityCreatedYear,
		VulnerabilityCreated,
		VulnerabilityID,
		PackageInfoFilepath,
		PackageInfoName,
		PackageInfoVersion,
		PackageInfoVersionFixed,
		FindingDescription,
		VulnerabilitySourceURL,
		VulnerabilityReferenceURLs,
	}
}

// TableColumnsImageVulnerabilities returns rows where
// each row is an image+vulnerability.
func TableColumnsImageVulnerabilities() ([]string, map[int]string) {
	return []string{
			ImageRepositoryName,
			FindingSeverity,
			VulnerabilityCreatedYear, // table.FormatInt
			VulnerabilityCreated,     // table.FormatDate
			VulnerabilitySLADueDate,  // table.FormatDate
			ImageNameVulnerabilityID,
			ImageHash,
			VulnerabilityID,
			VulnerabilitySourceURL, // table.FormatURL
			PackagesManagers,
			PackagesTypes,
			PackagesFilepathsPOMProperites,
			PackagesNamesAndFilepathsAtVersion,
			PackagesNamesAndFilepathsAtVersionFixed,
		}, map[int]string{
			2: table.FormatInt,
			3: table.FormatDate,
			4: table.FormatDate,
			8: table.FormatURL,
		}
}

func (f Finding) MustVulnerabilityField(field, def string, opts *govex.ValueOptions) string {
	if v, err := f.VulnerabilityField(field, opts); err != nil {
		return def
	} else {
		return v
	}
}

func (f Finding) VulnerabilityField(field string, opts *govex.ValueOptions) (string, error) {
	if f.PackageVulnerabilityDetails == nil {
		return "", nil
	}
	switch field {
	case FindingDescription:
		return pointer.Dereference(f.Description), nil
	case FindingSeverity:
		return f.FindingSeverity(true), nil
	case FindingTitle:
		return pointer.Dereference(f.Title), nil
	case ImageHash:
		hashes := f.ImageHashes()
		return strings.Join(hashes, ", "), nil
	case ImageRepositoryName:
		names := f.ImageRepositoryNames()
		return strings.Join(names, ", "), nil
	case ImageNameVulnerabilityID:
		names := f.ImageRepositoryNames()
		return strings.Join(names, "+") + sepFilepathVersion + pointer.Dereference(f.PackageVulnerabilityDetails.VulnerabilityId), nil
	case VulnerabilityCreated:
		return pointer.Dereference(f.PackageVulnerabilityDetails.VendorCreatedAt).Format(time.RFC3339), nil
	case VulnerabilityCreatedYear:
		return pointer.Dereference(f.PackageVulnerabilityDetails.VendorCreatedAt).Format("2006"), nil
	case VulnerabilityCreatedAgeMonthsInt:
		if f.PackageVulnerabilityDetails.VendorCreatedAt == nil {
			return "", nil
		} else {
			return fmt.Sprintf("%.0f",
				float64(time.Since(*f.PackageVulnerabilityDetails.VendorCreatedAt))/
					float64(timeutil.Day*30)), nil
		}
	case VulnerabilityID:
		return pointer.Dereference(f.PackageVulnerabilityDetails.VulnerabilityId), nil
	case VulnerabilityReferenceURLs:
		return strings.Join(f.PackageVulnerabilityDetails.ReferenceUrls, ", "), nil
	case VulnerabilitySeverity:
		return pointer.Dereference(f.PackageVulnerabilityDetails.VendorSeverity), nil
	case VulnerabilitySLADueDate:
		if opts == nil || opts.SLAOptions == nil {
			return "", nil
		} else if dueDate, err := opts.SLAOptions.DueDate(
			f.FindingOrVendorSeverity(true),
			f.PackageVulnerabilityDetails.VendorCreatedAt,
			nil); err != nil {
			return "", err
		} else if dueDate == nil {
			return "", nil
		} else {
			return dueDate.Format(time.RFC3339), nil
		}
	case VulnerabilitySourceURL:
		return pointer.Dereference(f.PackageVulnerabilityDetails.SourceUrl), nil
	case PackagesFilepathsAtVersion:
		if f.PackageVulnerabilityDetails != nil {
			pkgs := Packages(f.PackageVulnerabilityDetails.VulnerablePackages)
			return pkgs.FilepathsAtVersion(), nil
		} else {
			return "", nil
		}
	case PackagesFilepathsAtVersionFixed:
		if f.PackageVulnerabilityDetails != nil {
			pkgs := Packages(f.PackageVulnerabilityDetails.VulnerablePackages)
			return pkgs.FilepathsAtVersionFixed(), nil
		} else {
			return "", nil
		}
	case PackagesFilepathsPOMProperites:
		if f.PackageVulnerabilityDetails != nil {
			pkgs := Packages(f.PackageVulnerabilityDetails.VulnerablePackages)
			if havePOMProperties := pkgs.FilepathsContainsPOMProperties(); havePOMProperties > 0 {
				return "all", nil
			} else if havePOMProperties < 0 {
				return "none", nil
			} else {
				return "some", nil
			}
		} else {
			return "false", nil
		}
	case PackagesManagers:
		if f.PackageVulnerabilityDetails != nil {
			pkgs := Packages(f.PackageVulnerabilityDetails.VulnerablePackages)
			return strings.Join(pkgs.PackagesManagers(), ","), nil
		} else {
			return "", nil
		}
	case PackagesNamesAndFilepathsAtVersion:
		if f.PackageVulnerabilityDetails != nil {
			pkgs := Packages(f.PackageVulnerabilityDetails.VulnerablePackages)
			return pkgs.NamesAndFilepathsAtVersion(), nil
		} else {
			return "", nil
		}
	case PackagesNamesAndFilepathsAtVersionFixed:
		if f.PackageVulnerabilityDetails != nil {
			pkgs := Packages(f.PackageVulnerabilityDetails.VulnerablePackages)
			return pkgs.NamesAndFilepathsAtVersionFixed(), nil
		} else {
			return "", nil
		}
	case PackagesTypes:
		if f.PackageVulnerabilityDetails != nil {
			pkgs := Packages(f.PackageVulnerabilityDetails.VulnerablePackages)
			return pkgs.PackagesTypes(), nil
		} else {
			return "", nil
		}
	default:
		return "", fmt.Errorf("field unknown or not supported (%s)", field)
	}
}

// VulnerabilitySlices returns one slice per vulnerable package.
func (f Finding) VulnerabilityFields(fields []string, opts *govex.ValueOptions) ([]string, error) {
	var row []string
	for _, field := range fields {
		if v, err := f.VulnerabilityField(field, opts); err != nil {
			return []string{}, err
		} else {
			row = append(row, v)
		}
	}
	if len(fields) != len(row) {
		panic("output row len mismatch")
	}
	return row, nil
}

// PackageSlices returns one slice per vulnerable package.
func (f Finding) PackageSlices(fields []string, opts *govex.ValueOptions) ([][]string, error) {
	var rows [][]string
	fmtMap := map[int]string{}
	for _, res := range f.Resources {
		r2 := Resource(res)
		imgHash := r2.ImageHash()
		imgName := r2.ImageRepositoryName()
		for _, pkg := range f.VulnerablePackages() {
			var row []string
			for i, field := range fields {
				switch field {
				case FindingDescription:
					row = append(row, pointer.Dereference(f.Description))
				case ImageHash:
					row = append(row, imgHash)
				case ImageRepositoryName:
					row = append(row, imgName)
				case VulnerabilityCreated:
					if v, err := f.VulnerabilityField(VulnerabilityCreated, opts); err != nil {
						return rows, err
					} else {
						row = append(row, v)
						fmtMap[i] = table.FormatTime
					}
				case VulnerabilityCreatedYear:
					if v, err := f.VulnerabilityField(VulnerabilityCreatedYear, opts); err != nil {
						return rows, err
					} else {
						row = append(row, v)
						fmtMap[i] = table.FormatInt
					}
				case VulnerabilityID:
					if v, err := f.VulnerabilityField(VulnerabilityID, opts); err != nil {
						return rows, err
					} else {
						row = append(row, v)
					}
				case VulnerabilitySeverity:
					if v, err := f.VulnerabilityField(VulnerabilitySeverity, opts); err != nil {
						return rows, err
					} else {
						row = append(row, v)
					}
				case VulnerabilitySourceURL:
					if v, err := f.VulnerabilityField(VulnerabilitySourceURL, opts); err != nil {
						return rows, err
					} else {
						row = append(row, v)
					}
				case VulnerabilityReferenceURLs:
					if v, err := f.VulnerabilityField(VulnerabilityReferenceURLs, opts); err != nil {
						return rows, err
					} else {
						row = append(row, v)
					}
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
	return rows, nil
}
