package inspector2util

import (
	"slices"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/slicesutil"
)

type Resource types.Resource

func (r Resource) ImageHash() string {
	if r.Details != nil &&
		r.Details.AwsEcrContainerImage != nil {
		return strings.TrimSpace(*r.Details.AwsEcrContainerImage.ImageHash)
	} else {
		return ""
	}
}

func (r Resource) ImageRepositoryName() string {
	if r.Details != nil &&
		r.Details.AwsEcrContainerImage != nil {
		return strings.TrimSpace(*r.Details.AwsEcrContainerImage.RepositoryName)
	} else {
		return ""
	}
}

func (r Resource) ImageTagFirst() string {
	if r.Details != nil &&
		r.Details.AwsEcrContainerImage != nil &&
		len(r.Details.AwsEcrContainerImage.ImageTags) > 0 {
		return r.Details.AwsEcrContainerImage.ImageTags[0]
	} else {
		return ""
	}
}

func (r Resource) ImageTags() []string {
	if r.Details != nil &&
		r.Details.AwsEcrContainerImage != nil &&
		len(r.Details.AwsEcrContainerImage.ImageTags) > 0 {
		tags := slices.Clone(r.Details.AwsEcrContainerImage.ImageTags)
		sort.Strings(tags)
		return slicesutil.Dedupe(tags)
	} else {
		return []string{}
	}
}

func (r Resource) HasTagsAny(tagsAny []string) bool {
	if len(tagsAny) == 0 {
		return false
	}
	rImgTags := r.ImageTags()
	if len(rImgTags) == 0 {
		return false
	}
	return slicesutil.MatchAny(rImgTags, tagsAny)
}

func (r Resource) Values(fields []string) ([]string, error) {
	var out []string
	for _, f := range fields {
		switch f {
		case ImageHash:
			out = append(out, r.ImageHash())
		case ImageRepositoryName:
			out = append(out, r.ImageRepositoryName())
		case ImageTags:
			out = append(out, strings.Join(r.ImageTags(), ", "))
		}
	}
	return out, nil
}
