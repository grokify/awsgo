package rdsutil

import (
	"errors"
	"sort"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/grokify/mogo/pointer"
)

type ParametersSet struct {
	Data map[string]rds.Parameter
}

func NewParametersSet() ParametersSet {
	return ParametersSet{Data: map[string]rds.Parameter{}}
}

func (ps *ParametersSet) AddParameters(p ...rds.Parameter) {
	if ps.Data == nil {
		ps.Data = map[string]rds.Parameter{}
	}
	for _, pi := range p {
		name := pointer.Dereference(pi.ParameterName)
		ps.Data[name] = pi
	}
}

func (ps *ParametersSet) Map() map[string]string {
	m := map[string]string{}
	for _, pi := range ps.Data {
		m[pointer.Dereference(pi.ParameterName)] = pointer.Dereference(pi.ParameterValue)
	}
	return m
}

func (ps *ParametersSet) Names() []string {
	var names []string
	for _, pi := range ps.Data {
		names = append(names, pointer.Dereference(pi.ParameterName))
	}
	sort.Strings(names)
	return names
}

var ErrParameterNotFound = errors.New("parameter not found")

func (ps *ParametersSet) GetValue(key string) (string, error) {
	if pi, ok := ps.Data[key]; !ok {
		return "", ErrParameterNotFound
	} else {
		return pointer.Dereference(pi.ParameterValue), nil
	}
}
