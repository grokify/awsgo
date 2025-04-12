package inspector2util

import (
	"errors"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/type/slicesutil"
)

type ImageSet struct {
	Set map[string]types.AwsEcrContainerImageDetails
}

func NewImageSet() *ImageSet {
	return &ImageSet{Set: map[string]types.AwsEcrContainerImageDetails{}}
}

func (is *ImageSet) Add(img types.AwsEcrContainerImageDetails) error {
	if img.ImageHash == nil {
		return errors.New("image hash cannot be nil")
	}
	if is.Set == nil {
		is.Set = map[string]types.AwsEcrContainerImageDetails{}
	}
	is.Set[*img.ImageHash] = img
	return nil
}

func (is *ImageSet) FilterHashes(hashesIncl []string) (*ImageSet, error) {
	hashesIncl = slicesutil.Dedupe(hashesIncl)
	out := NewImageSet()
	if len(hashesIncl) == 0 {
		return out, nil
	}
	for _, img := range is.Set {
		if img.ImageHash == nil {
			continue
		} else if !slices.Contains(hashesIncl, *img.ImageHash) {
			continue
		} else if err := out.Add(img); err != nil {
			return nil, err
		}
	}
	return out, nil
}
