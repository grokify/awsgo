package inspector2util

import (
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/stringsutil"
)

type Resources []types.Resource

func (res Resources) ImageHashes() []string {
	var out []string
	for _, ri := range res {
		rx := Resource(ri)
		out = append(out, rx.ImageHash())
	}
	return stringsutil.SliceCondenseSpace(out, true, true)
}
