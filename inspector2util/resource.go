package inspector2util

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
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
