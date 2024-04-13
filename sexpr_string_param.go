package main

import "unicode"

type SexprStringParam struct {
	Value  string
	Quoted bool
	Line   int
	Column int
}

func (sp *SexprStringParam) __sexprParamMemberDummy__() {}

func NewSexprStringParam(v string, line int, column int) *SexprStringParam {
	quoted := false
	for _, r := range v {
		if r == '(' || r == ')' || unicode.IsSpace(r) {
			quoted = true
			break
		}
	}
	return NewSexprStringParamQuoted(v, quoted, line, column)
}

func NewSexprStringParamQuoted(v string, quoted bool, line int, column int) *SexprStringParam {
	return &SexprStringParam{
		Quoted: quoted,
		Value:  v,
		Line:   line,
		Column: column,
	}
}

func (sp *SexprStringParam) String() string {
	if sp.Quoted {
		return `"` + sp.Value + `"`
	}
	return sp.Value
}

func (sp *SexprStringParam) Set(v string) {
	for _, r := range v {
		if r == '(' || r == ')' || unicode.IsSpace(r) {
			sp.Quoted = true
			break
		}
	}
	sp.Value = v
}

func (sp *SexprStringParam) SetQuoted(v string, quoted bool) {
	sp.Quoted = true
	sp.Value = v
}
