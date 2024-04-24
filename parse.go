package sexpr

import (
	"bufio"
	"fmt"
)

type tmpSexpr struct {
	parent *tmpSexpr
	name   string
	params *[]*SexprParam
	line   int
	col    int
}

func Parse(input *bufio.Reader) (*Sexpr, error) {
	lexer := NewLexer(input)
	depth := 0

	var root *Sexpr = nil
	var sexpr *Sexpr
	var token Token
	var hasName bool

	for {
		lexer.NextToken(&token)

		if token.Kind == TokenWhitespace {
			// ignore

		} else if token.Kind == TokenOpen {
			tmpsexpr := NewSexpr("", []*SexprParam{}, sexpr, token.Line, token.Column)
			if depth == 0 {
				if root != nil {
					return nil, fmt.Errorf("unexpected open at Line %d, Column %d", token.Line, token.Column)
				}
				root = tmpsexpr
			}
			if depth > 0 && !hasName {
				return nil, fmt.Errorf("unexpected open at Line %d, Column %d", token.Line, token.Column)
			}
			hasName = false
			sexpr = tmpsexpr
			depth++

		} else if token.Kind == TokenClose {
			if depth == 0 {
				return nil, fmt.Errorf("unexpected close at Line %d, Column %d", token.Line, token.Column)
			}
			if !hasName {
				return nil, fmt.Errorf("sexpression without name at Line %d, Column %d", token.Line, token.Column)
			}
			if depth > 1 {
				sp, err := NewSexprParam(sexpr)
				if err != nil {
					return nil, err
				}
				sexpr.Parent().params = append(sexpr.Parent().params, sp)
			}
			sexpr = sexpr.Parent()
			depth--

		} else if token.Kind == TokenString {
			if depth == 0 {
				return nil, fmt.Errorf("unexpected string at Line %d, Column %d: '%s'", token.Line, token.Column, token.Content)
			}
			if hasName {
				ss := NewSexprStringQuoted(token.Content, false, sexpr, token.Line, token.Column)
				sp, err := NewSexprParam(ss)
				if err != nil {
					return nil, err
				}
				sexpr.params = append(sexpr.Params(), sp)
			} else {
				sexpr.name = token.Content
				hasName = true
			}

		} else if token.Kind == TokenQuotedString {
			if depth == 0 {
				return nil, fmt.Errorf("unexpected string at Line %d, Column %d: '%s'", token.Line, token.Column, token.Content)
			}
			if hasName {
				ss := NewSexprStringQuoted(token.Content[1:len(token.Content)-1], true, sexpr, token.Line, token.Column)
				sp, err := NewSexprParam(ss)
				if err != nil {
					return nil, err
				}
				sexpr.params = append(sexpr.Params(), sp)
			} else {
				sexpr.name = token.Content[1 : len(token.Content)-1]
				hasName = true
			}

		} else if token.Kind == TokenEOF {
			if depth != 0 {
				return nil, fmt.Errorf("unexpected EOF at Line %d, Column %d", token.Line, token.Column)
			}
			return root, nil

		} else if token.Kind == TokenErr {
			return nil, fmt.Errorf("error at Line %d, Column %d: %s", token.Line, token.Column, token.Err.Error())

		}
	}
}
