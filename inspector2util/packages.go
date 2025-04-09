package inspector2util

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/mogo/type/strslices"
)

type Packages []types.VulnerablePackage

func (ps Packages) FilepathsContainsPOMProperities() int {
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

func sliceCondenseSpaceAndJoin(s []string) string {
	if s = stringsutil.SliceCondenseSpace(s, true, true); len(s) == 0 {
		return ""
	} else {
		return strings.Join(s, sepJoinDispolay)
	}
}
