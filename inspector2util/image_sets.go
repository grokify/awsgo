package inspector2util

import (
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/mogo/pointer"
)

type ImageSets struct {
	Data map[string]ImageSet
}

func NewImageSets() *ImageSets {
	return &ImageSets{
		Data: map[string]ImageSet{},
	}
}

func (sets *ImageSets) AddImagesByRepositoryName(imgs ...types.AwsEcrContainerImageDetails) error {
	if len(imgs) == 0 {
		return nil
	}
	for _, img := range imgs {
		repoName := pointer.Dereference(img.RepositoryName)
		if _, ok := sets.Data[repoName]; !ok {
			sets.Data[repoName] = pointer.Dereference(NewImageSet())
		}
		if set, ok := sets.Data[repoName]; !ok {
			panic("set not found")
		} else if err := set.Add(img); err != nil {
			return err
		} else {
			sets.Data[repoName] = set
		}
	}
	return nil
}

func (sets *ImageSets) FilterLatestByRepositoryName(imageTagOverrides []string) (*ImageSets, error) {
	latestSets := NewImageSets()
	for _, set := range sets.Data {
		if imgs := set.ImagesByTagsOrLatest(imageTagOverrides); len(imgs) > 0 {
			if err := latestSets.AddImagesByRepositoryName(imgs...); err != nil {
				return nil, err
			}
		}
	}
	return latestSets, nil
}

func (sets *ImageSets) ImageSetByRepositoryNameTagsOrLatest(repoTagOverrides []string) (*ImageSet, error) {
	latestSet := NewImageSet()
	for _, set := range sets.Data {
		if imgs := set.ImagesByTagsOrLatest(repoTagOverrides); len(imgs) > 0 {
			if err := latestSet.Add(imgs...); err != nil {
				return nil, err
			}
		}
	}
	return latestSet, nil
}
