package inspector2util

import (
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/govex"
	"github.com/grokify/govex/reports/poam"
	"github.com/grokify/govex/severity"
	"github.com/grokify/mogo/type/slicesutil"
)

func (f Finding) POAMItemOpen() bool {
	sev := f.FindingSeverity(true)
	sevs := severity.SeveritiesFinding()
	return slices.Contains(sevs, sev)
}

func (f Finding) POAMItemClosed() bool {
	return false
}

func (f Finding) POAMItemValue(field poam.POAMField, opts *govex.ValueOptions, overrides func(field poam.POAMField) (*string, error)) (string, error) {
	var dateFormat string
	if opts != nil && opts.DateFormat != "" {
		dateFormat = opts.DateFormat
	} else {
		dateFormat = time.DateOnly
	}
	// Step 1: check overrides
	if overrides != nil {
		if v, err := overrides(field); err != nil {
			return "", err
		} else if v != nil {
			return *v, nil
		}
	}
	// Step 2: check mappings
	m := MappingPOAM2Inspector()
	if ifield, ok := m[field]; ok {
		return f.VulnerabilityField(ifield, opts)
	}
	switch field {
	case poam.FieldWeaknessDetectorSource:
		return "Inspector", nil
	case poam.FieldBindingOperationalDirective2201Tracking:
		if opts == nil || opts.CISAKEVC == nil {
			return "", nil
		} else if cveID := f.CVEID(); cveID == nil {
			return yesnoNo, nil
		} else if kev := opts.CISAKEVC.CVE(*cveID); kev == nil {
			return yesnoNo, nil
		} else {
			return yesnoYes, nil
		}
	case poam.FieldBindingOperationalDirective2201DueDate:
		if opts == nil || opts.CISAKEVC == nil {
			return "", nil
		} else if cveID := f.CVEID(); cveID == nil {
			return "", nil
		} else if kev := opts.CISAKEVC.CVE(*cveID); kev == nil {
			return "", nil
		} else {
			return kev.DueDate, nil
		}
	case poam.FieldControls:
		if f.Type == types.FindingTypePackageVulnerability {
			return "RA-5", nil
		} else {
			return "", nil
		}
	case poam.FieldCVE:
		if vulnID := f.VulnerabilityID(); strings.HasPrefix(vulnID, "CVE-") {
			return vulnID, nil
		} else {
			return "", nil
		}
	case poam.FieldOriginalDetectionDate:
		if opts != nil && opts.SLAOptions != nil && opts.SLAOptions.SLAStartDateFixed != nil {
			return opts.SLAOptions.SLAStartDateFixed.Format(dateFormat), nil
		} else if f.FirstObservedAt != nil {
			return f.FirstObservedAt.Format(dateFormat), nil
		} else {
			return "", nil
		}
	case poam.FieldOriginalRiskRating:
		return f.FindingSeverity(true), nil
	case poam.FieldOverallRemediationPlan:
		remInfo := f.POAMItemUpgradeRemedationInfo(opts)
		return remInfo.String()
	case poam.FieldScheduledCompletionDate:
		var slaStart *time.Time
		if opts != nil && opts.SLAOptions != nil && opts.SLAOptions.SLAStartDateFixed != nil {
			slaStart = opts.SLAOptions.SLAStartDateFixed
		} else if f.FirstObservedAt != nil {
			slaStart = f.FirstObservedAt
		}
		if opts.SLAOptions.SLAPolicy != nil {
			sev := f.FindingOrVendorSeverity(true)
			if dueDate, err := opts.SLAOptions.SLAPolicy.DueDate(sev, *slaStart); err != nil {
				return "", err
			} else {
				return dueDate.Format(dateFormat), nil
			}
		} else {
			return "", nil
		}
	case poam.FieldServiceName:
		names := slicesutil.Dedupe(f.ImageRepositoryNames())
		return strings.Join(names, ", "), nil
	case poam.FieldVendorDependency:
		return "Yes", nil
	default:
		return "", nil
	}
}

func MappingPOAM2Inspector() map[poam.POAMField]string {
	return map[poam.POAMField]string{
		poam.FieldWeaknessName:        FindingTitle,
		poam.FieldWeaknessDescription: FindingDescription,
		poam.FieldAssetIdentifier:     ImageHash,
	}
}

func (f Finding) POAMItemValues(fields []poam.POAMField, opts *govex.ValueOptions, overrides func(field poam.POAMField) (*string, error)) ([]string, error) {
	var out []string
	for _, field := range fields {
		if v, err := f.POAMItemValue(field, opts, overrides); err != nil {
			return out, err
		} else {
			out = append(out, v)
		}
	}
	return out, nil
}

func (f Finding) POAMItemUpgradeRemedationInfo(opts *govex.ValueOptions) poam.POAMItemUpgradeRemedationInfo {
	info := poam.POAMItemUpgradeRemedationInfo{
		VulnerabilityID: f.VulnerabilityID(),
		Packages:        poam.POAMItemUpgradeRemedationPackages{},
		SLADays:         0,
	}
	if opts != nil && opts.SLAOptions != nil && opts.SLAOptions.SLAPolicy != nil {
		sev := f.FindingOrVendorSeverity(true)
		info.SLADays = opts.SLAOptions.SLAPolicy.SeveritySLADays(sev)
	}
	if f.PackageVulnerabilityDetails != nil {
		for _, vpkg := range f.PackageVulnerabilityDetails.VulnerablePackages {
			vpkg2 := Package(vpkg)
			info.Packages = append(info.Packages, vpkg2.POAMItem())
		}
	}
	return info
}
