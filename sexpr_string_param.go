package main

type SexprStringParam struct {
	Value  string
	Line   int
	Column int
}

func (sp *SexprStringParam) __sexprParamMemberDummy__() {}

func (sp *SexprStringParam) String() string {
	return sp.Value
}
