package inspector2util

import (
	"errors"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type ImageSet struct {
	Data map[string]types.AwsEcrContainerImageDetails `json:"data"`
}

func NewImageSet() *ImageSet {
	return &ImageSet{Data: map[string]types.AwsEcrContainerImageDetails{}}
}

func (set *ImageSet) Add(imgs ...types.AwsEcrContainerImageDetails) error {
	if len(imgs) == 0 {
		return nil
	}
	for _, img := range imgs {
		if img.ImageHash == nil {
			return errors.New("image hash cannot be nil")
		}
		if set.Data == nil {
			set.Data = map[string]types.AwsEcrContainerImageDetails{}
		}
		set.Data[*img.ImageHash] = img
	}
	return nil
}

func (set *ImageSet) FilterHashes(hashesIncl []string) (*ImageSet, error) {
	hashesIncl = slicesutil.Dedupe(hashesIncl)
	out := NewImageSet()
	if len(hashesIncl) == 0 {
		return out, nil
	}
	for _, img := range set.Data {
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

func (set *ImageSet) ImageByHash(imgHash string) *types.AwsEcrContainerImageDetails {
	if img, ok := set.Data[imgHash]; ok {
		return &img
	} else {
		return nil
	}
}

func (set *ImageSet) ImageByLatest() *types.AwsEcrContainerImageDetails {
	if len(set.Data) == 0 {
		return nil
	}
	var latest *types.AwsEcrContainerImageDetails
	dt := time.Time{}
	for _, img := range set.Data {
		if img.PushedAt == nil {
			continue
		} else if dt.IsZero() || dt.After(*img.PushedAt) {
			dt = *img.PushedAt
			latest = &img
		}
	}
	return latest
}

func (set *ImageSet) ImagesByTag(imageTags []string) []types.AwsEcrContainerImageDetails {
	var out []types.AwsEcrContainerImageDetails
	imageTags = slicesutil.Dedupe(imageTags)
	sort.Strings(imageTags)
	if len(imageTags) == 0 {
		return out
	}
	for _, img := range set.Data {
		for _, imgTag := range imageTags {
			if slices.Contains(img.ImageTags, imgTag) {
				out = append(out, img)
				continue
			}
		}
	}
	return out
}

func (set *ImageSet) ImagesByTagsOrLatest(imageTagOverrides []string) []types.AwsEcrContainerImageDetails {
	imgs := set.ImagesByTag(imageTagOverrides)
	if len(imgs) > 0 {
		return imgs
	} else if imgLatest := set.ImageByLatest(); imgLatest != nil {
		imgs = append(imgs, *imgLatest)
	}
	return imgs
}

func (set *ImageSet) ImageSetsByRepoName() (*ImageSets, error) {
	sets := NewImageSets()
	for _, img := range set.Data {
		if err := sets.AddImagesByRepositoryName(img); err != nil {
			return nil, err
		}
	}
	return sets, nil
}

func (set *ImageSet) ImageHashes() []string {
	var s []string
	for _, img := range set.Data {
		if img.ImageHash != nil {
			s = append(s, strings.TrimSpace(*img.ImageHash))
		}
	}
	return stringsutil.SliceCondenseSpace(s, true, true)
}

func (set *ImageSet) ImageTags() ([]string, map[string]int) {
	var s []string
	m := map[string]int{}
	for _, img := range set.Data {
		for _, tag := range img.ImageTags {
			if _, ok := m[tag]; !ok {
				s = append(s, tag)
			}
			m[tag]++
		}
	}
	sort.Strings(s)
	return s, m
}

func (set *ImageSet) RepositoryNames() ([]string, map[string]int) {
	var s []string
	m := map[string]int{}
	for _, img := range set.Data {
		if img.RepositoryName != nil {
			if _, ok := m[*img.RepositoryName]; !ok {
				s = append(s, *img.RepositoryName)
			}
			m[*img.RepositoryName]++
		}
	}
	sort.Strings(s)
	return s, m
}

func (set *ImageSet) Table() (*table.Table, error) {
	t := table.NewTable("")
	t.Columns = []string{"repo_name", "image_hash", "push_time", "image_tags"}
	t.FormatMap = map[int]string{2: table.FormatTime}
	for _, img := range set.Data {
		var timeString string
		if img.PushedAt != nil {
			timeString = img.PushedAt.Format(time.RFC3339)
		}
		row := []string{
			pointer.Dereference(img.RepositoryName),
			pointer.Dereference(img.ImageHash),
			timeString,
			strings.Join(img.ImageTags, ", ")}
		t.Rows = append(t.Rows, row)
	}
	return &t, nil
}

func (set *ImageSet) WriteFileJSON(filename, prefix, indent string, perm os.FileMode) error {
	return jsonutil.MarshalFile(filename, set, prefix, indent, perm)
}

func (set *ImageSet) WriteFileXLSX(filename string) error {
	if t, err := set.Table(); err != nil {
		return err
	} else {
		return t.WriteXLSX(filename, "images")
	}
}
