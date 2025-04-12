package inspector2util

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/type/maputil"
)

type ResourceSet struct {
	Set map[string]Resource
}

func NewResourceSet() *ResourceSet {
	return &ResourceSet{Set: map[string]Resource{}}
}

func (rs *ResourceSet) Add(r Resource) error {
	if r.Id == nil {
		return errors.New("no resource id")
	}
	rs.Set[*r.Id] = Resource(r)
	return nil
}

func (rs *ResourceSet) ImageRepositoryNames() []string {
	m := map[string]int{}
	for _, r := range rs.Set {
		if name := r.ImageRepositoryName(); name != "" {
			m[name]++
		}
	}
	return maputil.Keys(m)
}

func (rs *ResourceSet) ImageTags() []string {
	m := map[string]int{}
	for _, r := range rs.Set {
		if r.Details == nil || r.Details.AwsEcrContainerImage == nil {
			continue
		}
		for _, tag := range r.Details.AwsEcrContainerImage.ImageTags {
			m[tag]++
		}
	}
	return maputil.Keys(m)
}

func (rs *ResourceSet) FilterImageHash(hashesIncl []string) *ResourceSet {
	out := NewResourceSet()
	for _, r := range rs.Set {
		for _, h := range hashesIncl {
			if h == r.ImageHash() {

			}
		}
	}
	return out
}

func (rs *ResourceSet) FilterImageTags(tagsAny []string) *ResourceSet {
	out := NewResourceSet()
	for k, r := range rs.Set {
		if r.HasTagsAny(tagsAny) {
			out.Set[k] = r
		}
	}
	return out
}

func (rs *ResourceSet) FilterResourceTypes(inclResourceTypes []types.ResourceType) (*ResourceSet, error) {
	out := NewResourceSet()
	for _, r := range rs.Set {
		if len(inclResourceTypes) == 0 {
			if err := out.Add(r); err != nil {
				return rs, err
			}
		} else {
			for _, rt := range inclResourceTypes {
				if r.Type == rt {
					if err := out.Add(r); err != nil {
						return rs, err
					} else {
						break
					}
				}
			}
		}
	}
	return out, nil
}

func ResourceTableColumns() ([]string, map[int]string) {
	return []string{
		ImageRepositoryName,
		ImageTags,
		ImageHash,
	}, map[int]string{}
}

func (rs *ResourceSet) Table(cols []string, fmtMap map[int]string) (*table.Table, error) {
	if len(cols) == 0 {
		cols, fmtMap = ResourceTableColumns()
	}
	t := table.NewTable("")
	t.Columns = cols
	t.FormatMap = fmtMap
	for _, r := range rs.Set {
		if row, err := r.Values(cols); err != nil {
			return nil, err
		} else {
			t.Rows = append(t.Rows, row)
		}
	}
	return &t, nil
}
