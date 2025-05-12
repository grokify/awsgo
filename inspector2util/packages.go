package inspector2util

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/mogo/type/strslices"
)

type Packages []types.VulnerablePackage

func (ps Packages) FilepathsContainsPOMProperties() int {
	return strslices.Contains(ps.Filepaths(), filenamePomProperties)
}

func (ps Packages) Filepaths() []string {
	var paths []string
	for _, p := range ps {
		if v := Package(p).FilepathString(); v != "" {
			paths = append(paths, v)
		}
	}
	return stringsutil.SliceCondenseSpace(paths, true, true)
}

func (ps Packages) FilepathsAtVersion() string {
	var paths []string
	for _, p := range ps {
		if v := Package(p).FilepathAtVersion(); v != "" {
			paths = append(paths, v)
		}
	}
	return sliceCondenseSpaceAndJoin(paths)
}

func (ps Packages) FilepathsAtVersionFixed() string {
	var paths []string
	for _, p := range ps {
		if v := Package(p).FilepathAtVersionFixed(); v != "" {
			paths = append(paths, v)
		}
	}
	return sliceCondenseSpaceAndJoin(paths)
}

func (ps Packages) NameAtVersionAtFilepaths() []string {
	var out []string
	for _, pi := range ps {
		px := Package(pi)
		out = append(out, px.NameAtVersionAtFilepath())
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}

func (ps Packages) NamesAndFilepathsAtVersion() string {
	var paths []string
	for _, p := range ps {
		if v := Package(p).NameAndFilepathAtVersion(); v != "" {
			paths = append(paths, v)
		}
	}
	return sliceCondenseSpaceAndJoin(paths)
}

func (ps Packages) NamesAndFilepathsAtVersionFixed() string {
	var paths []string
	for _, p := range ps {
		if v := Package(p).NameAndFilepathAtVersionFixed(); v != "" {
			paths = append(paths, v)
		}
	}
	return sliceCondenseSpaceAndJoin(paths)
}

func (ps Packages) PackagesManagers() []string {
	var out []string
	for _, p := range ps {
		out = append(out, string(p.PackageManager))
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}

func (ps Packages) PackagesTypes() string {
	if len(ps) == 0 {
		return ""
	}
	pkgTypes := map[string]int{}
	for _, p := range ps {
		px := Package(p)
		pkgTypes[px.PackageType()]++
	}
	switch len(pkgTypes) {
	case 0:
		return ""
	case 1:
		for k := range pkgTypes {
			return k
		}
	}
	return "both"
}

func sliceCondenseSpaceAndJoin(s []string) string {
	if s = stringsutil.SliceCondenseSpace(s, true, true); len(s) == 0 {
		return ""
	} else {
		return strings.Join(s, sepJoinDispolay)
	}
}
