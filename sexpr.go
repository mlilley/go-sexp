package main

import (
	"strings"
)

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
	wasSexpr := false
	indent := strings.Repeat("\t", level)
	acc.WriteString(indent)
	acc.WriteString("(")
	acc.WriteString(sp.Name)
	if len(sp.Params) == 0 {
		acc.WriteString(")\n")
	} else {
		for _, param := range sp.Params {
			if sparam, ok := param.(*Sexpr); ok {
				acc.WriteString("\n")
				sparam.string_(acc, level+1)
				wasSexpr = true
			} else {
				acc.WriteString(" ")
				acc.WriteString(param.String())
				wasSexpr = false
			}
		}
		if wasSexpr {
			acc.WriteString("\n")
			acc.WriteString(indent)
			acc.WriteString(")")
		} else {
			acc.WriteString(")")
		}
	}
}
