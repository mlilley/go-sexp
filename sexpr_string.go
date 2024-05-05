package sexpr

import "unicode"

type SexprString struct {
	value  string
	quoted bool

	parent *Sexpr
	line   int
	col    int
}

func NewSexprString(v string) *SexprString {
	return &SexprString{
		value:  v,
		quoted: shouldQuote(v),
	}
}

func NewSexprStringQuoted(v string, quoted bool) *SexprString {
	return &SexprString{
		value:  v,
		quoted: quoted,
	}
}

func (ss *SexprString) Value() string {
	return ss.value
}

func (ss *SexprString) SetValue(v string) {
	ss.value = v
	ss.quoted = shouldQuote(v)
}

func (ss *SexprString) SetValueQuoted(v string, quoted bool) {
	ss.value = v
	ss.quoted = quoted
}

func (ss *SexprString) Quoted() bool {
	return ss.quoted
}

func (ss *SexprString) Parent() *Sexpr {
	return ss.parent
}

func (ss *SexprString) SetParent(parent *Sexpr) {
	ss.parent = parent
}

func (ss *SexprString) Location() (int, int) {
	return ss.line, ss.col
}

func (ss *SexprString) SetLocation(line int, col int) {
	ss.line = line
	ss.col = col
}

func (ss *SexprString) String() string {
	if ss.quoted {
		return `"` + ss.value + `"`
	} else {
		return ss.value
	}
}

func shouldQuote(v string) bool {
	for _, r := range v {
		if r == '(' || r == ')' || unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
