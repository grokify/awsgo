package regions

import (
	"strings"
	"testing"

	"github.com/grokify/mogo/location"
)

var validSubdivisionCategories = map[string]bool{
	location.SubdivisionTypeCanton:     true,
	location.SubdivisionTypeACT:        true,
	location.SubdivisionTypeCapital:    true,
	location.SubdivisonTypeCountry:     true,
	location.SubdivisionTypeCounty:     true,
	location.SubdivisionTypePrefecture: true,
	location.SubdivisionTypeProvince:   true,
	location.SubdivisionTypeSAR:        true,
	location.SubdivisionTypeSCR:        true,
	location.SubdivisionTypeState:      true,
	"":                                 true, // e.g. southeastasia has no subdivision
}

func TestRegionsAzure(t *testing.T) {
	azureRegions := RegionsAzure()
	if len(azureRegions) == 0 {
		t.Fatal("RegionsAzure() returned no regions")
	}

	seen := map[string]bool{}
	for _, r := range azureRegions {
		if r.RegionType != RegionTypeAzure {
			t.Errorf("region %q: RegionType = %q, want %q", r.RegionCode, r.RegionType, RegionTypeAzure)
		}
		if r.RegionCode == "" {
			t.Error("region has empty RegionCode")
			continue
		}
		if seen[r.RegionCode] {
			t.Errorf("duplicate RegionCode %q", r.RegionCode)
		}
		seen[r.RegionCode] = true

		if r.RegionName == "" {
			t.Errorf("region %q: RegionName is empty", r.RegionCode)
		}
		if r.CityName == "" {
			t.Errorf("region %q: CityName is empty", r.RegionCode)
		}
		if r.ISO3166P1A2CountryCode == "" {
			t.Errorf("region %q: ISO3166P1A2CountryCode is empty", r.RegionCode)
		}
		if r.SubregionsCount <= 0 {
			t.Errorf("region %q: SubregionsCount = %d, want > 0", r.RegionCode, r.SubregionsCount)
		}

		// UN/LOCODEs are always 5 characters: 2-letter country code + 3-letter location code.
		if len(r.UNLOCODE) != 5 {
			t.Errorf("region %q: UNLOCODE = %q, want 5 characters", r.RegionCode, r.UNLOCODE)
		} else if !strings.HasPrefix(r.UNLOCODE, r.ISO3166P1A2CountryCode) {
			t.Errorf("region %q: UNLOCODE %q does not start with country code %q",
				r.RegionCode, r.UNLOCODE, r.ISO3166P1A2CountryCode)
		}

		if !validSubdivisionCategories[r.ISO3166P2SubdivisionCategory] {
			t.Errorf("region %q: unknown ISO3166P2SubdivisionCategory %q", r.RegionCode, r.ISO3166P2SubdivisionCategory)
		}

		if len(r.ReferenceURLs) == 0 {
			t.Errorf("region %q: no ReferenceURLs", r.RegionCode)
		}
		for _, u := range r.ReferenceURLs {
			if !strings.HasPrefix(u, "https://") {
				t.Errorf("region %q: ReferenceURL %q does not look like a URL", r.RegionCode, u)
			}

			// Regression guard: a prior copy-paste bug attached Australia's
			// Border Force UN/LOCODE reference to non-AU regions.
			if strings.Contains(u, "abf.gov.au") && r.ISO3166P1A2CountryCode != "AU" {
				t.Errorf("region %q (country %q): non-AU region references Australian Border Force URL %q",
					r.RegionCode, r.ISO3166P1A2CountryCode, u)
			}

			// Regression guard: UNECE locode reference URLs must match the
			// region's own country, not a copy-pasted neighbor's.
			const unecePrefix = "https://service.unece.org/trade/locode/"
			if strings.HasPrefix(u, unecePrefix) {
				wantSuffix := strings.ToLower(r.ISO3166P1A2CountryCode) + ".htm"
				if !strings.HasSuffix(u, wantSuffix) {
					t.Errorf("region %q (country %q): UNECE locode URL %q does not match country",
						r.RegionCode, r.ISO3166P1A2CountryCode, u)
				}
			}
		}
	}
}
