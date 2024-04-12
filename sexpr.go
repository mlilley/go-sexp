package main

import "strings"

type Sexpr struct {
	Name   string
	Params []SexprParam
	Line   int
	Column int
	Parent *Sexpr
}

func (sp *Sexpr) __sexprParamMemberDummy__() {}

func (sp *Sexpr) AddParam(p SexprParam) {
	sp.Params = append(sp.Params, p)
}

func (sp *Sexpr) String() string {
	var b strings.Builder
	sp.string_(&b, 0)
	return b.String()
}

func (sp *Sexpr) string_(acc *strings.Builder, level int) {
	indent := strings.Repeat("\t", level)
	acc.WriteString(indent)
	acc.WriteString("(\n")
	acc.WriteString(indent)
	acc.WriteString("\t")
	acc.WriteString(sp.Name)
	acc.WriteString("\n")
	for _, param := range sp.Params {
		if sparam, ok := param.(*Sexpr); ok {
			sparam.string_(acc, level+1)
		} else {
			acc.WriteString(indent)
			acc.WriteString("\t")
			acc.WriteString(param.String())
			acc.WriteString("\n")
		}
	}
	acc.WriteString(indent)
	acc.WriteString(")\n")
}
