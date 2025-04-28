package inspector2util

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/stringsutil"
)

type FindingMerger struct {
	Data map[string]Findings
}

func NewFindingMerger() FindingMerger {
	return FindingMerger{Data: map[string]Findings{}}
}

func (set *FindingMerger) Add(f ...types.Finding) error {
	for _, fi := range f {
		fx := Finding(fi)
		rnvids := fx.ImageRepoNameVulnIDs("")
		rnvids = stringsutil.SliceCondenseSpace(rnvids, true, true)
		if len(rnvids) == 0 {
			return errors.New("no repo vuln ")
		}
		for _, rnvid := range rnvids {
			set.Data[rnvid] = append(set.Data[rnvid], fi)
		}
	}
	return nil
}

func (set *FindingMerger) Merge() (Findings, error) {
	out := Findings{}
	for _, fs := range set.Data {
		if merged, err := fs.MergeFilteredByImageRepoNameAndVulnID(); err != nil {
			return out, err
		} else if merged != nil {
			out = append(out, *merged)
		}
	}
	return out, nil
}
