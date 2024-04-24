package sexpr

import (
	"errors"
	"strconv"
)

type SexprParam struct {
	kind  SexprParamKind
	value interface{}
}

func NewSexprParam(v interface{}) (*SexprParam, error) {
	var sp SexprParam
	switch v.(type) {
	case *SexprString:
		sp = SexprParam{value: v, kind: SexprParamKindString}
	case *Sexpr:
		sp = SexprParam{value: v, kind: SexprParamKindSexpr}
	default:
		return nil, errors.New("value must be string or sexpr")
	}
	return &sp, nil
}

func (sp *SexprParam) Value() any {
	return sp.value
}

func (sp *SexprParam) Kind() SexprParamKind {
	return sp.kind
}

func (sp *SexprParam) String() string {
	switch sp.kind {
	case SexprParamKindString:
		return sp.value.(*SexprString).String()
	case SexprParamKindSexpr:
		return sp.value.(*Sexpr).String()
	default:
		return ""
	}
}

func (sp *SexprParam) AsSexpr() (*Sexpr, error) {
	ss, ok := sp.value.(*Sexpr)
	if !ok {
		return nil, errors.New("value is not a sexpr")
	}
	return ss, nil
}

func (sp *SexprParam) AsString() (string, error) {
	ss, ok := sp.value.(*SexprString)
	if !ok {
		return "", errors.New("value is not a string")
	}
	return ss.Value(), nil
}

func (sp *SexprParam) AsInt() (int64, error) {
	ss, ok := sp.value.(*SexprString)
	if !ok {
		return 0, errors.New("value is not an int")
	}
	i, err := strconv.ParseInt(ss.Value(), 10, 64)
	if err != nil {
		return 0, errors.New("value is not an int")
	}
	return i, nil
}

func (sp *SexprParam) AsFloat() (float64, error) {
	ss, ok := sp.value.(*SexprString)
	if !ok {
		return 0, errors.New("value is not a float")
	}
	f, err := strconv.ParseFloat(ss.Value(), 64)
	if err != nil {
		return 0, errors.New("value is not a float")
	}
	return f, nil
}
