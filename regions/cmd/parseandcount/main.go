package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/grokify/awsgo/regions"
	"github.com/grokify/mogo/location"
	"github.com/grokify/mogo/location/country"
	"github.com/grokify/mogo/type/stringsutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	File string `short:"f" long:"file" description:"File of whitespace-separated AWS/Azure region codes to analyze" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	myRegions, err := readRegionCodesFile(opts.File)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("MyRegionCount(%d)\n", len(myRegions))

	myLocs := location.Locations{}

	awsLocs := regions.Regions()
	for _, awsLoc := range awsLocs {
		if slices.Contains(myRegions, awsLoc.RegionCode) {
			myLocs = append(myLocs, awsLoc)
		}
	}

	azureLocs := regions.RegionsAzure()
	for _, azureLoc := range azureLocs {
		if slices.Contains(myRegions, azureLoc.RegionCode) {
			myLocs = append(myLocs, azureLoc)
		}
	}

	fmt.Printf("MyLocationCount(%d)\n", len(myLocs))

	for i, r := range myRegions {
		fmt.Printf("[%d](%s)\n", i, r)
		if !myLocs.ContainsRegionCode(r) {
			fmt.Println("MISSING")
		}
	}

	fmt.Printf("US DC Count(%d)\n", myLocs.SubregionsCountByCountry([]string{"US"}))
	fmt.Printf("CA DC Count(%d)\n", myLocs.SubregionsCountByCountry([]string{"CA"}))

	dach := country.CountriesMapDACH()
	fmt.Printf("DACH DC Count(%d)\n", myLocs.SubregionsCountByCountry(dach.ISO3166P1Alpha2s()))
	nord := country.CountriesMapNordic()
	fmt.Printf("Nord DC Count(%d)\n", myLocs.SubregionsCountByCountry(nord.ISO3166P1Alpha2s()))
	bene := country.CountriesMapBenelux()
	fmt.Printf("Bene DC Count(%d)\n", myLocs.SubregionsCountByCountry(bene.ISO3166P1Alpha2s()))

	fmt.Printf("GB DC Count(%d)\n", myLocs.SubregionsCountByCountry([]string{"GB"}))

	ee := country.CountriesMapEasternEuropeFull()
	fmt.Printf("EEurope DC Count(%d)\n", myLocs.SubregionsCountByCountry(ee.ISO3166P1Alpha2s()))

	me := country.CountriesMapNearAndMiddleEast()
	fmt.Printf("NME DC Count(%d)\n", myLocs.SubregionsCountByCountry(me.ISO3166P1Alpha2s()))

	anz := country.CountriesMapANZ()
	fmt.Printf("ANZ DC Count(%d)\n", myLocs.SubregionsCountByCountry(anz.ISO3166P1Alpha2s()))

	fmt.Printf("CN DC Count(%d)\n", myLocs.SubregionsCountByCountry([]string{"CN"}))
	fmt.Printf("JP DC Count(%d)\n", myLocs.SubregionsCountByCountry([]string{"JP"}))

	capac := country.CountriesMapAPACSales()
	capac = capac.Sub(map[string]string{"CN": "China"})
	capac = capac.Sub(map[string]string{"JP": "Japan"})
	capac = capac.Sub(anz)
	fmt.Printf("APAC REST Count(%d)\n", myLocs.SubregionsCountByCountry(capac.ISO3166P1Alpha2s()))

	ceu := country.CountriesMapEurope()
	ceu = ceu.Sub(dach)
	ceu = ceu.Sub(bene)
	ceu = ceu.Sub(nord)
	ceu = ceu.Sub(map[string]string{"GB": "Great Britain"})
	fmt.Printf("Europe REST Count(%d)\n", myLocs.SubregionsCountByCountry(ceu.ISO3166P1Alpha2s()))
}

// readRegionCodesFile reads whitespace-separated region codes from filename,
// lowercasing and deduplicating them.
func readRegionCodesFile(filename string) ([]string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(string(b))
	for i, fi := range fields {
		fields[i] = strings.ToLower(fi)
	}
	return stringsutil.SliceCondenseSpace(fields, true, true), nil
}
