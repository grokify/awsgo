package rdsutil

import (
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/maputil"
)

type ParametersSets struct {
	Data map[string]ParametersSet // paramSetName
}

func NewParametersSets() *ParametersSets {
	return &ParametersSets{Data: map[string]ParametersSet{}}
}

func (pss *ParametersSets) init() {
	if pss.Data == nil {
		pss.Data = map[string]ParametersSet{}
	}
}

func (pss *ParametersSets) Names() []string {
	pss.init()
	return maputil.Keys(pss.Data)
}

// GetParamNameValues returns the values for a particular param name where the keys
// are the ParameterSet keys and the value is the value of the param. If the param
// is not present, a `nil` is returned in the map.
func (pss *ParametersSets) GetParamNameValues(paramName string) map[string]*string {
	pss.init()
	comp := map[string]*string{}
	names := pss.Names()
	for _, name := range names {
		ps, ok := pss.Data[name]
		if !ok {
			comp[name] = nil
			continue
		}
		v, err := ps.GetValue(name)
		if err != nil {
			comp[name] = nil
			continue
		}
		comp[name] = pointer.Pointer(v)
	}
	return comp
}
