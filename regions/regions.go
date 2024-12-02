// regions is based on the information here:
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html and
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html
// https://docs.aws.amazon.com/local-zones/latest/ug/available-local-zones.html
package regions

import "github.com/grokify/mogo/location"

const (
	RegionTypeAWS   = "AWS"
	RegionTypeAzure = "Azure"
)

func Regions() location.Locations {
	return []location.Location{
		{
			RegionCode:                   "af-south-1",
			RegionName:                   "Africa (Cape Town)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Cape Town",
			UNLOCODE:                     "ZACPT",
			ISO3166P1A2CountryCode:       "ZA",
			ISO3166P2SubdivisionCodeFull: "ZA-WC",
			ISO3166P2CountryCode:         "ZA",
			ISO3166P2SubdivisionCode:     "WC",
			ISO3166P2SubdivisionName:     "Western Cape",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeProvince,
			Subregions: []string{
				"afs1-az1", "afs1-az2", "afs1-az3",
			},
		},
		{
			RegionCode:                   "ap-east-1",
			RegionName:                   "Asia Pacific (Hong Kong)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Hong Kong",
			UNLOCODE:                     "HKHKG",
			ISO3166P1A2CountryCode:       "CN",
			ISO3166P2SubdivisionCodeFull: "HK-HK",
			ISO3166P2CountryCode:         "HK",
			ISO3166P2SubdivisionCode:     "HK",
			ISO3166P2SubdivisionName:     "Special Administrative Region (SAR) of Hong Kong",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeSAR,
			Subregions: []string{
				"ape1-az1", "ape1-az2", "ape1-az3",
			},
		},
		{
			RegionCode:                   "ap-northeast-1",
			RegionName:                   "Asia Pacific (Tokyo)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Tokyo",
			UNLOCODE:                     "JPTYO",
			ISO3166P1A2CountryCode:       "JP",
			ISO3166P2SubdivisionCodeFull: "JP-13",
			ISO3166P2CountryCode:         "JP",
			ISO3166P2SubdivisionCode:     "13",
			ISO3166P2SubdivisionName:     "Tokyo Metropolis (東京都, Tōkyō-to)",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypePrefecture,
			Subregions: []string{
				"apne1-az1", "apne1-az2", "apne1-az3", "apne1-az4",
			},
		},
		{
			RegionCode:                   "ap-northeast-2",
			RegionName:                   "Asia Pacific (Seoul)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Seoul",
			UNLOCODE:                     "KRSEL",
			ISO3166P1A2CountryCode:       "KR",
			ISO3166P2SubdivisionCodeFull: "KR-11",
			ISO3166P2CountryCode:         "KR",
			ISO3166P2SubdivisionCode:     "11",
			ISO3166P2SubdivisionName:     "Seoul Special City (서울특별시, Seoul-teukbyeolsi)",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypePrefecture,
			Subregions: []string{
				"apne2-az1", "apne2-az2", "apne2-az3", "apne2-az4",
			},
		},
		{
			RegionCode:                   "ap-south-1",
			RegionName:                   "Asia Pacific (Mumbai)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Mumbai",
			UNLOCODE:                     "INBOM",
			ISO3166P1A2CountryCode:       "IN",
			ISO3166P2SubdivisionCodeFull: "IN-MH",
			ISO3166P2CountryCode:         "IN",
			ISO3166P2SubdivisionCode:     "MH",
			ISO3166P2SubdivisionName:     "Maharashtra",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"aps1-az1", "aps1-az2", "aps1-az3",
			},
		},
		{
			RegionCode:                   "ap-southeast-1",
			RegionName:                   "Asia Pacific (Singapore)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Singapore",
			UNLOCODE:                     "SG",
			ISO3166P1A2CountryCode:       "SG",
			ISO3166P2SubdivisionCodeFull: "SG",
			ISO3166P2CountryCode:         "SG",
			ISO3166P2SubdivisionCode:     "",
			ISO3166P2SubdivisionName:     "",
			ISO3166P2SubdivisionCategory: "",
			Subregions: []string{
				"apse1-az1", "apse1-az2", "apse1-az3",
			},
		},
		{
			RegionCode:                   "ap-southeast-2",
			RegionName:                   "Asia Pacific (Sydney)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Sydney",
			UNLOCODE:                     "AUSYD",
			ISO3166P1A2CountryCode:       "AU",
			ISO3166P2SubdivisionCodeFull: "AU-NSW",
			ISO3166P2CountryCode:         "AU",
			ISO3166P2SubdivisionCode:     "NSW",
			ISO3166P2SubdivisionName:     "New South Wales",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"apse2-az1", "apse2-az2", "apse2-az3",
			},
		},
		{
			RegionCode:                   "ap-southeast-3",
			RegionName:                   "Asia Pacific (Jakarta)",
			RegionType:                   RegionTypeAWS,
			CityName:                     "Jakarta",
			UNLOCODE:                     "IDJKT",
			ISO3166P1A2CountryCode:       "ID",
			ISO3166P2SubdivisionCodeFull: "ID-JK",
			ISO3166P2CountryCode:         "ID",
			ISO3166P2SubdivisionCode:     "JK",
			ISO3166P2SubdivisionName:     "Jakarta Special Capital Region",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeSCR,
			Subregions: []string{
				"apse3-az1", "apse3-az2", "apse3-az3",
			},
		},
		{
			RegionCode:                   "ca-central-1",
			RegionName:                   "Canada (Central)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "CA",
			ISO3166P2SubdivisionCodeFull: "CA-ON",
			ISO3166P2CountryCode:         "CA",
			ISO3166P2SubdivisionCode:     "ON",
			ISO3166P2SubdivisionName:     "Ontario",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeProvince,
			UNLOCODE:                     "CATOR",
			CityName:                     "Toronto",
			Subregions: []string{
				"cac1-az1", "cac1-az2", "cac1-az4",
			},
		},
		{
			RegionCode:                   "eu-central-1",
			RegionName:                   "Europe (Frankfurt)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "DE",
			ISO3166P2SubdivisionCodeFull: "DE-HE",
			ISO3166P2CountryCode:         "DE",
			ISO3166P2SubdivisionCode:     "HE", // EN:Hesse, DE:Hessen
			ISO3166P2SubdivisionName:     "Hesse",
			ISO3166P2SubdivisionCategory: "",
			UNLOCODE:                     "DEFRA",
			CityName:                     "Frankfurt",
			Subregions: []string{
				"euc1-az1", "euc1-az2", "euc1-az3",
			},
		},
		{
			RegionCode:                   "eu-west-1",
			RegionName:                   "Europe (Ireland)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "IE",
			ISO3166P2SubdivisionCodeFull: "IE-D",
			ISO3166P2CountryCode:         "IE",
			ISO3166P2SubdivisionCode:     "D",
			ISO3166P2SubdivisionName:     "County Dublin",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeCounty,
			UNLOCODE:                     "IEDUB",
			CityName:                     "Dublin",
			Subregions: []string{
				"euw1-az1", "euw1-az2", "euw1-az3",
			},
		},
		{
			RegionCode:                   "us-east-1",
			RegionName:                   "US East (N. Virginia)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "US",
			ISO3166P2SubdivisionCodeFull: "US-VA",
			ISO3166P2CountryCode:         "US",
			ISO3166P2SubdivisionCode:     "VA",
			ISO3166P2SubdivisionName:     "Virginia",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"use1-az1", "use1-az2", "use1-az3", "use1-az4", "use1-az5", "use1-az6",
			},
		},
		{
			RegionCode:                   "us-east-2",
			RegionName:                   "US East (Ohio)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "US",
			ISO3166P2SubdivisionCodeFull: "US-OH",
			ISO3166P2CountryCode:         "US",
			ISO3166P2SubdivisionCode:     "OH",
			ISO3166P2SubdivisionName:     "Ohio",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"use2-az1", "use2-az2", "use2-az3",
			},
		},
		{
			RegionCode:                   "us-gov-east-1",
			RegionName:                   "AWS GovCloud (US-East)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "US",
			ISO3166P2SubdivisionCodeFull: "US-VA",
			ISO3166P2CountryCode:         "US",
			ISO3166P2SubdivisionCode:     "VA",
			ISO3166P2SubdivisionName:     "Virginia",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"usge1-az1", "usge1-az2", "usge1-az3",
			},
		},
		{
			RegionCode:                   "us-gov-west-1",
			RegionName:                   "AWS GovCloud (US-West)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "US",
			ISO3166P2SubdivisionCodeFull: "US-CA",
			ISO3166P2CountryCode:         "US",
			ISO3166P2SubdivisionCode:     "CA",
			ISO3166P2SubdivisionName:     "California",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"usgw1-az1", "usgw1-az2", "usgw1-az3",
			},
		},
		{
			RegionCode:                   "us-west-1",
			RegionName:                   "US West (N. California)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "US",
			ISO3166P2SubdivisionCodeFull: "US-CA",
			ISO3166P2CountryCode:         "US",
			ISO3166P2SubdivisionCode:     "CA",
			ISO3166P2SubdivisionName:     "California",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"usw1-az1", "usw1-az2", "usw1-az3",
			},
		},
		{
			RegionCode:                   "us-west-2",
			RegionName:                   "US West (Oregon)",
			RegionType:                   RegionTypeAWS,
			ISO3166P1A2CountryCode:       "US",
			ISO3166P2SubdivisionCodeFull: "US-OR",
			ISO3166P2CountryCode:         "US",
			ISO3166P2SubdivisionCode:     "OR",
			ISO3166P2SubdivisionName:     "Oregon",
			ISO3166P2SubdivisionCategory: location.SubdivisionTypeState,
			Subregions: []string{
				"usw2-az1", "usw2-az2", "usw2-az3", "usw2-az4",
			},
		},
	}
}

/*
US East (N. Virginia) – use1-az1 | use1-az2 | use1-az3 | use1-az4 | use1-az5 | use1-az6

US East (Ohio) – use2-az1 | use2-az2 | use2-az3

US West (N. California) – usw1-az1 | usw1-az2 | usw1-az3

US West (Oregon) – usw2-az1 | usw2-az2 | usw2-az3 | usw2-az4

Africa (Cape Town) – afs1-az1 | afs1-az2 | afs1-az3

Asia Pacific (Hong Kong) – ape1-az1 | ape1-az2 | ape1-az3

Asia Pacific (Hyderabad) – aps2-az1 | aps2-az2 | aps2-az3

Asia Pacific (Jakarta) – apse3-az1 | apse3-az2 | apse3-az3

Asia Pacific (Malaysia) – apse5-az1 | apse5-az2 | apse5-az3

Asia Pacific (Melbourne) – apse4-az1 | apse4-az2 | apse4-az3

Asia Pacific (Mumbai) – aps1-az1 | aps1-az2 | aps1-az3

Asia Pacific (Osaka) – apne3-az1 | apne3-az2 | apne3-az3

Asia Pacific (Seoul) – apne2-az1 | apne2-az2 | apne2-az3 | apne2-az4

Asia Pacific (Singapore) – apse1-az1 | apse1-az2 | apse1-az3

Asia Pacific (Sydney) – apse2-az1 | apse2-az2 | apse2-az3

Asia Pacific (Tokyo) – apne1-az1 | apne1-az2 | apne1-az3 | apne1-az4

Canada (Central) – cac1-az1 | cac1-az2 | cac1-az4

Canada West (Calgary) – caw1-az1 | caw1-az2 | caw1-az3

Europe (Frankfurt) – euc1-az1 | euc1-az2 | euc1-az3

Europe (Ireland) – euw1-az1 | euw1-az2 | euw1-az3

Europe (London) – euw2-az1 | euw2-az2 | euw2-az3

Europe (Milan) – eus1-az1 | eus1-az2 | eus1-az3

Europe (Paris) – euw3-az1 | euw3-az2 | euw3-az3

Europe (Spain) – eus2-az1 | eus2-az2 | eus2-az3

Europe (Stockholm) – eun1-az1 | eun1-az2 | eun1-az3

Europe (Zurich) – euc2-az1 | euc2-az2 | euc2-az3

Israel (Tel Aviv) – ilc1-az1 | ilc1-az2 | ilc1-az3

Middle East (Bahrain) – mes1-az1 | mes1-az2 | mes1-az3

Middle East (UAE) – mec1-az1 | mec1-az2 | mec1-az3

South America (São Paulo) – sae1-az1 | sae1-az2 | sae1-az3

AWS GovCloud (US-East) – usge1-az1 | usge1-az2 | usge1-az3

AWS GovCloud (US-West) – usgw1-az1 | usgw1-az2 | usgw1-az3


*/
