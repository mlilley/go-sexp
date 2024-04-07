package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Sexpr struct {
	params []SexprParam
}

type SexprParam interface {
	String() string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("expect filename as argument")
	}
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sexpr, err := ParseSexpr(bufio.NewReader(f), false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sexpr.String())
}

// ---

func (s *Sexpr) GetParam(index int) SexprParam {
	return s.params[index]
}

func (s *Sexpr) AddParam(param SexprParam) {
	s.params = append(s.params, param)
}

func (s *Sexpr) UpdateParam(index int, param SexprParam) {
	s.params[index] = param
}

func (s *Sexpr) Size() int {
	return len(s.params)
}

func (s Sexpr) String() string {
	var sb strings.Builder
	s.string_(&sb, 0)
	return sb.String()
}

func (s *Sexpr) string_(acc *strings.Builder, level int) {
	indent := strings.Repeat("  ", level)
	acc.WriteString(indent)
	acc.WriteString("(\n")
	for _, param := range s.params {
		if sexpr, ok := param.(*Sexpr); ok {
			sexpr.string_(acc, level+1)
		} else {
			acc.WriteString(indent)
			acc.WriteString("  ")
			acc.WriteString(param.String())
			acc.WriteString("\n")
		}
	}
	acc.WriteString(indent)
	acc.WriteString(")\n")
}

// ---

type StringParam struct {
	value string
}

func (p *StringParam) String() string {
	return p.value
}
func (p *StringParam) SetValue(value string) {
	p.value = value // todo: check valid unquoted string param value
}
func NewStringParam(value string) *StringParam {
	return &StringParam{value} // todo: check valid unquoted string param value
}

// ---

type QuotedStringParam struct {
	value string
}

func (p *QuotedStringParam) String() string {
	return `"` + p.value + `"`
}
func (p *QuotedStringParam) SetValue(unquotedValue string) {
	p.value = unquotedValue
}
func NewQuotedStringParam(unquotedValue string) *QuotedStringParam {
	return &QuotedStringParam{unquotedValue}
}

// ---

type IntParam struct {
	value int
}

func (p *IntParam) String() string {
	return strconv.Itoa(p.value)
}
func (p *IntParam) GetValue() int {
	return p.value
}
func (p *IntParam) SetValue(value int) {
	p.value = value
}
func NewIntParam(value int) *IntParam {
	return &IntParam{value}
}
func NewIntParamFromInput(value string) (*IntParam, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &IntParam{v}, nil
}

// ---

type FloatParam struct {
	value      float64
	changed    bool
	inputValue string // to prevent reserializing float values if unchanged
}

func (p *FloatParam) String() string {
	if p.changed {
		return strconv.FormatFloat(p.value, 'f', -1, 64)
	}
	return p.inputValue
}
func (p *FloatParam) GetValue() float64 {
	return p.value
}
func (p *FloatParam) Setvalue(value float64) {
	p.changed = true
	p.value = value
}
func (p *FloatParam) NewFloatParam(value float64) *FloatParam {
	return &FloatParam{value, true, ""}
}
func NewFloatParamFromInput(value string) (*FloatParam, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	return &FloatParam{v, false, value}, nil
}

// ---

func ParseSexpr(reader *bufio.Reader, parseNumerics bool) (*Sexpr, error) {
	token, err := nextToken(reader, parseNumerics)
	if err != nil {
		// an EOF terminates an empty input file
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	switch t := token.(type) {
	case rune:
		if t == '(' {
			sexpr, err := parseSexpr(reader, parseNumerics)
			if err != nil {
				return nil, err
			}
			return sexpr, nil
		}
	}

	return nil, errors.New("unexpected token")
}

func parseSexpr(reader *bufio.Reader, parseNumerics bool) (*Sexpr, error) {
	sexpr := &Sexpr{}

	for {
		token, err := nextToken(reader, parseNumerics)
		if err != nil {
			if err == io.EOF {
				// an EOF illegally terminates an sexpr
				return nil, errors.New("unexpected EOF reading sexpr")
			}
			return nil, err
		}

		switch t := token.(type) {
		case rune:
			// proper termination of sexp
			if t == ')' {
				return sexpr, nil
			}
			// param is a sexpr, recurse into it
			if t == '(' {
				sexpr2, err := parseSexpr(reader, parseNumerics)
				if err != nil {
					if err == io.EOF {
						return nil, errors.New("unexpected EOF reading sexp")
					}
					return nil, err
				}
				sexpr.AddParam(sexpr2)
				continue
			}
			// unexpected rune
			return nil, errors.New("unexpected token (rune)")

		case SexprParam:
			sexpr.AddParam(t)
		}
	}
}

// returns rune or SexprParam
func nextToken(reader *bufio.Reader, parseNumerics bool) (any, error) {
	inQuoted := false
	inUnquoted := false
	var token string

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// an EOF illegally terminates a quoted token
				if inQuoted {
					return nil, errors.New("unexpected EOF reading quoted token")
				}
				// an EOF illegally terminates an unquoted token
				if inUnquoted {
					return nil, errors.New("unexpected EOF reading unquoted token")
				}
			}
			return nil, err
		}

		if inQuoted {
			// only a quote terminates a quoted token
			if r == '"' {
				return NewQuotedStringParam(token), nil
			}

			// append rune to quoted token
			token += string(r)
			continue
		}

		if inUnquoted {
			if r == '(' {
				// illegal sexpr start
				return nil, errors.New("illegal sexpr start")
			}
			if r == '"' {
				// illegal quoted start
				return nil, errors.New("illegal quoted start")
			}
			if r == ')' {
				// push back the ), end token, and return it
				err = reader.UnreadRune()
				if err != nil {
					return nil, err
				}

				if parseNumerics {
					if v, err := NewIntParamFromInput(token); err == nil {
						return v, nil
					}
					if v, err := NewFloatParamFromInput(token); err == nil {
						return v, nil
					}
				}

				return NewStringParam(token), nil
			}
			if unicode.IsSpace(r) {
				// end token, and return it
				if parseNumerics {
					if v, err := NewIntParamFromInput(token); err == nil {
						return v, nil
					}
					if v, err := NewFloatParamFromInput(token); err == nil {
						return v, nil
					}
				}

				return NewStringParam(token), nil
			}

			// add the non-(, non-), non-", non-ws rune to the unquoted token
			token += string(r)
			continue
		}

		// return parens as tokens
		if r == '(' || r == ')' {
			return r, nil
		}

		// start reading a quoted token
		if r == '"' {
			token = ""
			inQuoted = true
			continue
		}

		// skip over leading whitespace
		if unicode.IsSpace(r) {
			continue
		}

		// start an unquoted token
		token = string(r)
		inUnquoted = true
	}
}
