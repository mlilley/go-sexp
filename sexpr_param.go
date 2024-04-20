package sexpr

import "errors"

type SexprParamKind int

const (
	SexprParamKindString SexprParamKind = iota
	SexprParamKindSexpr
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
	return sp.String()
}
