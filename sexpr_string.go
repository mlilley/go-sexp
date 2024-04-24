package sexpr

import "unicode"

type SexprString struct {
	parent *Sexpr

	value  string
	quoted bool

	line int
	col  int
}

func NewSexprString(v string, parent *Sexpr, line int, col int) *SexprString {
	return &SexprString{
		value:  v,
		quoted: shouldQuote(v),
		parent: parent,
		line:   line,
		col:    col,
	}
}

func NewSexprStringQuoted(v string, quoted bool, parent *Sexpr, line int, col int) *SexprString {
	return &SexprString{
		value:  v,
		quoted: quoted,
		parent: parent,
		line:   line,
		col:    col,
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

func (ss *SexprString) Value() string {
	return ss.value
}

func (ss *SexprString) Quoted() bool {
	return ss.quoted
}

func (ss *SexprString) Set(v string) {
	ss.value = v
	ss.quoted = shouldQuote(v)
}

func (ss *SexprString) String() string {
	if ss.quoted {
		return `"` + ss.value + `"`
	} else {
		return ss.value
	}
}
