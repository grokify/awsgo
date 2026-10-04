package inspector2util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	findingspec "github.com/plexusone/findingspec"
	"github.com/plexusone/findingspec/security"
)

func sampleFindings() Findings {
	return Findings{
		{
			Type:         types.FindingTypePackageVulnerability,
			Severity:     types.SeverityHigh,
			Title:        aws.String("CVE-2024-1234 - openssl"),
			Description:  aws.String("A flaw in openssl"),
			FixAvailable: types.FixAvailableYes,
			PackageVulnerabilityDetails: &types.PackageVulnerabilityDetails{
				VulnerabilityId: aws.String("CVE-2024-1234"),
				Source:          aws.String("NVD"),
				SourceUrl:       aws.String("https://nvd.nist.gov/vuln/detail/CVE-2024-1234"),
				ReferenceUrls:   []string{"https://example.com/advisory"},
				Cvss: []types.CvssScore{
					{
						BaseScore:     aws.Float64(7.5),
						ScoringVector: aws.String("CVSS:3.1/AV:N/AC:L"),
						Version:       aws.String("3.1"),
						Source:        aws.String("NVD"),
					},
				},
				VulnerablePackages: []types.VulnerablePackage{
					{
						Name:            aws.String("openssl"),
						Version:         aws.String("1.1.1k"),
						Arch:            aws.String("x86_64"),
						FilePath:        aws.String("/usr/lib/openssl"),
						FixedInVersion:  aws.String("1.1.1w"),
						PackageManager:  types.PackageManagerOs,
						SourceLayerHash: aws.String("sha256:layer1"),
					},
					{
						Name:           aws.String("libssl"),
						Version:        aws.String("1.1.1k"),
						PackageManager: types.PackageManagerOs,
					},
				},
			},
			Resources: []types.Resource{
				{
					Id:   aws.String("arn:aws:ecr:us-east-1:123:repository/app"),
					Type: types.ResourceTypeAwsEcrContainerImage,
					Details: &types.ResourceDetails{
						AwsEcrContainerImage: &types.AwsEcrContainerImageDetails{
							ImageHash:      aws.String("sha256:abc"),
							RepositoryName: aws.String("app"),
							Architecture:   aws.String("amd64"),
							Platform:       aws.String("linux"),
							ImageTags:      []string{"v1.0.0"},
						},
					},
				},
			},
		},
	}
}

func TestFindings(t *testing.T) {
	set := sampleFindings().Findings(FindingSpecOptions{Repo: "acme/app"})
	if set.Len() != 2 {
		t.Fatalf("Len = %d; want 2 (one per vulnerable package)", set.Len())
	}

	for _, f := range set.Findings {
		if err := f.Validate(); err != nil {
			t.Fatalf("invalid finding: %v", err)
		}
		if f.Domain != findingspec.DomainSecurity || f.Type != security.TypeContainer {
			t.Errorf("domain/type = %s/%s; want security/container", f.Domain, f.Type)
		}
		if f.Source.Tool != "aws-inspector" {
			t.Errorf("source = %q; want aws-inspector", f.Source.Tool)
		}
		if f.RuleID != "CVE-2024-1234" {
			t.Errorf("ruleID = %q; want CVE-2024-1234", f.RuleID)
		}
		if f.Severity != findingspec.SeverityHigh {
			t.Errorf("severity = %q; want high", f.Severity)
		}
	}

	first := set.Findings[0]
	d, err := findingspec.DetailAs[security.VulnerabilityDetail](first)
	if err != nil {
		t.Fatalf("DetailAs: %v", err)
	}
	if d.Package == nil || d.Package.Name != "openssl" || d.Package.Version != "1.1.1k" {
		t.Errorf("package detail = %+v", d.Package)
	}
	if d.Package == nil || d.Package.Path != "/usr/lib/openssl" || d.Package.Arch != "x86_64" {
		t.Errorf("package path/arch detail = %+v", d.Package)
	}
	if d.CVSS == nil || d.CVSS.Score != 7.5 || d.CVSS.Version != "3.1" {
		t.Errorf("cvss detail = %+v", d.CVSS)
	}
	if d.Fix == nil || d.Fix.State != security.FixStateFixed || len(d.Fix.Versions) != 1 {
		t.Errorf("fix detail = %+v", d.Fix)
	}
	if d.Artifact == nil || d.Artifact.Image != "app:v1.0.0" || d.Artifact.ImageDigest != "sha256:abc" {
		t.Errorf("artifact detail = %+v", d.Artifact)
	}
	if d.Artifact == nil || d.Artifact.Layer != "sha256:layer1" {
		t.Errorf("artifact layer detail = %+v", d.Artifact)
	}

	// The second package has no fix version, so Fix should be nil (FixAvailable=YES).
	second := set.Findings[1]
	d2, err := findingspec.DetailAs[security.VulnerabilityDetail](second)
	if err != nil {
		t.Fatalf("DetailAs (second): %v", err)
	}
	if d2.Package == nil || d2.Package.Name != "libssl" {
		t.Errorf("second package detail = %+v", d2.Package)
	}
	if d2.Fix != nil {
		t.Errorf("second fix detail = %+v; want nil", d2.Fix)
	}
}

func TestToFindingSkipsNonPackageVuln(t *testing.T) {
	f := Finding{Type: types.FindingTypeNetworkReachability}
	if got := f.ToFinding(FindingSpecOptions{}); got != nil {
		t.Errorf("ToFinding for non-package-vuln = %+v; want nil", got)
	}
}
