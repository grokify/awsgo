package inspector2util

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	findingspec "github.com/plexusone/findingspec"
	"github.com/plexusone/findingspec/security"
)

// FindingSpecOptions configures conversion of AWS Inspector2 findings to
// findingspec.
type FindingSpecOptions struct {
	// Repo sets Location.Repo for multi-repo/workspace sweeps, optional.
	Repo string
}

// Findings converts a list of AWS Inspector2 findings into a
// findingspec.FindingSet of normalized security vulnerability findings. Each
// vulnerable package within a finding produces its own findingspec.Finding.
func (fs Findings) Findings(opts FindingSpecOptions) *findingspec.FindingSet {
	set := findingspec.NewFindingSet()
	for _, f := range fs {
		set.Add(Finding(f).ToFinding(opts)...)
	}
	return set
}

// ToFinding converts a single AWS Inspector2 finding into a slice of
// findingspec.Finding values, one per vulnerable package. Findings that are not
// package vulnerabilities return nil.
func (f Finding) ToFinding(opts FindingSpecOptions) []findingspec.Finding {
	if f.Type != types.FindingTypePackageVulnerability {
		return nil
	}
	pvd := f.PackageVulnerabilityDetails
	if pvd == nil {
		return nil
	}

	vulnID := strings.TrimSpace(deref(pvd.VulnerabilityId))
	cvss := bestCVSS(pvd.Cvss)
	var score float64
	if cvss != nil {
		score = cvss.Score
	}
	sev := security.SeverityOrCVSS(string(f.Severity), score)
	refs := referencesFrom(pvd)

	title := strings.TrimSpace(deref(f.Title))
	if title == "" {
		title = vulnID
	}
	desc := strings.TrimSpace(deref(f.Description))

	resType, artifact := artifactFrom(f.Resources)
	resID := firstResourceID(f.Resources)

	var out []findingspec.Finding
	for _, vp := range pvd.VulnerablePackages {
		pkgName := strings.TrimSpace(deref(vp.Name))
		pkgVer := strings.TrimSpace(deref(vp.Version))

		pkg := &security.Package{
			Name:      pkgName,
			Version:   pkgVer,
			Ecosystem: string(vp.PackageManager),
			Path:      strings.TrimSpace(deref(vp.FilePath)),
			Arch:      strings.TrimSpace(deref(vp.Arch)),
		}

		// Copy the artifact so each package can carry its own source layer.
		var art *security.Artifact
		if artifact != nil {
			a := *artifact
			if layer := strings.TrimSpace(deref(vp.SourceLayerHash)); layer != "" {
				a.Layer = layer
			}
			art = &a
		}

		vuln := security.Vulnerability{
			ID:          strings.Join([]string{vulnID, pkgName + "@" + pkgVer, resID}, ":"),
			Title:       title,
			Description: desc,
			Type:        resType,
			CVEs:        cveIDs(vulnID),
			CVSS:        cvss,
			Severity:    sev,
			Package:     pkg,
			Fix:         fixFrom(vp, f.FixAvailable),
			Artifact:    art,
			Component:   pkgName + "@" + pkgVer,
			References:  refs,
			DetectedAt:  f.FirstObservedAt,
		}
		if opts.Repo != "" {
			vuln.Location = &findingspec.Location{Repo: opts.Repo}
		}

		fnd := vuln.ToFinding()
		fnd.RuleID = vulnID // preserve the Inspector vulnerability ID (CVE, ALAS, …)
		fnd.Source = findingspec.Source{Tool: "aws-inspector"}
		out = append(out, fnd)
	}
	return out
}

// cveIDs returns the vuln ID as a CVE list only when it is a CVE identifier.
func cveIDs(id string) []string {
	if strings.HasPrefix(strings.ToUpper(id), "CVE-") {
		return []string{id}
	}
	return nil
}

// bestCVSS selects the highest-scoring CVSS entry and maps it to security.CVSS.
func bestCVSS(list []types.CvssScore) *security.CVSS {
	var best *types.CvssScore
	for i := range list {
		if list[i].BaseScore == nil {
			continue
		}
		if best == nil || deref(list[i].BaseScore) > deref(best.BaseScore) {
			best = &list[i]
		}
	}
	if best == nil {
		return nil
	}
	return &security.CVSS{
		Version: strings.TrimSpace(deref(best.Version)),
		Vector:  strings.TrimSpace(deref(best.ScoringVector)),
		Score:   deref(best.BaseScore),
	}
}

// fixFrom maps a vulnerable package's fix data to a security.Fix. A fixed
// version yields FixStateFixed; otherwise FixAvailable == "NO" yields
// FixStateWontFix. When neither is known it returns nil.
func fixFrom(vp types.VulnerablePackage, avail types.FixAvailable) *security.Fix {
	if fixedVer := strings.TrimSpace(deref(vp.FixedInVersion)); fixedVer != "" {
		return &security.Fix{
			State:    security.FixStateFixed,
			Versions: []string{fixedVer},
		}
	}
	if avail == types.FixAvailableNo {
		return &security.Fix{State: security.FixStateWontFix}
	}
	return nil
}

// referencesFrom collects reference URLs and the source URL from the package
// vulnerability details.
func referencesFrom(pvd *types.PackageVulnerabilityDetails) []findingspec.Reference {
	var refs []findingspec.Reference
	if src := strings.TrimSpace(deref(pvd.SourceUrl)); src != "" {
		refs = append(refs, findingspec.Reference{Title: strings.TrimSpace(deref(pvd.Source)), URL: src})
	}
	for _, u := range pvd.ReferenceUrls {
		if u = strings.TrimSpace(u); u != "" {
			refs = append(refs, findingspec.Reference{URL: u})
		}
	}
	return refs
}

// artifactFrom derives the security finding Type and Artifact from the first
// resource. ECR container images map to TypeContainer; EC2 instances map to
// TypeRuntime.
func artifactFrom(resources []types.Resource) (string, *security.Artifact) {
	for i := range resources {
		r := resources[i]
		switch r.Type {
		case types.ResourceTypeAwsEcrContainerImage:
			art := &security.Artifact{}
			if d := r.Details; d != nil && d.AwsEcrContainerImage != nil {
				ci := d.AwsEcrContainerImage
				img := strings.TrimSpace(deref(ci.RepositoryName))
				if len(ci.ImageTags) > 0 {
					if tag := strings.TrimSpace(ci.ImageTags[0]); tag != "" {
						img = img + ":" + tag
					}
				}
				art.Image = img
				art.ImageDigest = strings.TrimSpace(deref(ci.ImageHash))
				art.OS = strings.TrimSpace(deref(ci.Platform))
			}
			return security.TypeContainer, art
		case types.ResourceTypeAwsEc2Instance:
			art := &security.Artifact{Image: strings.TrimSpace(deref(r.Id))}
			return security.TypeRuntime, art
		}
	}
	return "", nil
}

// firstResourceID returns the ID of the first resource, if any.
func firstResourceID(resources []types.Resource) string {
	for i := range resources {
		if id := strings.TrimSpace(deref(resources[i].Id)); id != "" {
			return id
		}
	}
	return ""
}

// deref safely dereferences a pointer, returning the zero value when nil.
func deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
