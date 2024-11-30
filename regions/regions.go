// regions is based on the information here:
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html and
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html
package regions

type RegionGeography struct {
	AWSRegionCode      string
	AWSRegionName      string
	ISOCountryCode     string
	ISOSubdivisionCode string
	ISOSubdivisionName string
	UNLOCODE           string
	CityName           string
}

func RegionGeographies() []RegionGeography {
	return []RegionGeography{
		{
			AWSRegionCode:      "eu-central-1",
			AWSRegionName:      "Europe (Frankfurt)",
			ISOCountryCode:     "DE",
			ISOSubdivisionCode: "HE", // EN:Hesse, DE:Hessen
			ISOSubdivisionName: "Hesse",
			UNLOCODE:           "DEFRA",
			CityName:           "Frankfurt",
		},
		{
			AWSRegionCode:      "eu-west-1",
			AWSRegionName:      "Europe (Ireland)",
			ISOCountryCode:     "IE",
			ISOSubdivisionCode: "D",
			ISOSubdivisionName: "County Dublin",
			UNLOCODE:           "IEDUB",
			CityName:           "Dublin",
		},
		{
			AWSRegionCode:      "us-east-1",
			AWSRegionName:      "US East (N. Virginia)",
			ISOCountryCode:     "US",
			ISOSubdivisionCode: "VA",
			ISOSubdivisionName: "Virginia",
		},
		{
			AWSRegionCode:      "us-east-2",
			AWSRegionName:      "US East (Ohio)",
			ISOCountryCode:     "US",
			ISOSubdivisionCode: "OH",
			ISOSubdivisionName: "Ohio",
		},
		{
			AWSRegionCode:      "us-gov-east-1",
			AWSRegionName:      "AWS GovCloud (US-East)",
			ISOCountryCode:     "US",
			ISOSubdivisionCode: "VA",
			ISOSubdivisionName: "Virginia",
		},
		{
			AWSRegionCode:      "us-gov-west-1",
			AWSRegionName:      "AWS GovCloud (US-West)",
			ISOCountryCode:     "US",
			ISOSubdivisionCode: "CA",
			ISOSubdivisionName: "California",
		},
		{
			AWSRegionCode:      "us-west-1",
			AWSRegionName:      "Northern California",
			ISOCountryCode:     "US",
			ISOSubdivisionCode: "CA",
			ISOSubdivisionName: "California",
		},
		{
			AWSRegionCode:      "us-west-2",
			AWSRegionName:      "US West (Oregon)",
			ISOCountryCode:     "US",
			ISOSubdivisionCode: "OR",
			ISOSubdivisionName: "Oregon",
		},
	}
}
