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

	var root *Sexpr = nil
	var sexpr *Sexpr
	var token Token

	for {
		lexer.NextToken(&token)

		if token.Kind == TokenWhitespace {
			// ignore

		} else if token.Kind == TokenOpen {
			if sexpr != nil && sexpr.Name() == "" {
				return nil, fmt.Errorf("unexpected open at Line %d, Column %d", token.Line, token.Column)
			}
			if sexpr == nil && root != nil {
				return nil, fmt.Errorf("unexpected open at Line %d, Column %d", token.Line, token.Column)
			}
			p := sexpr
			sexpr = NewSexpr("")
			sexpr.SetLocation(token.Line, token.Column)
			sexpr.SetParent(p)
			if p != nil {
				sp, err := NewSexprParam(sexpr)
				if err != nil {
					return nil, err
				}
				p.AddParam(len(p.Params()), sp)
			}
			if root == nil {
				root = sexpr
			}

		} else if token.Kind == TokenClose {
			if sexpr == nil {
				return nil, fmt.Errorf("unexpected close at Line %d, Column %d", token.Line, token.Column)
			}
			if sexpr.Name() == "" {
				return nil, fmt.Errorf("unexpected close at Line %d, Column %d", token.Line, token.Column)
			}
			sexpr = sexpr.Parent()

		} else if token.Kind == TokenString {
			if sexpr == nil {
				return nil, fmt.Errorf("unexpected string at Line %d, Column %d: '%s'", token.Line, token.Column, token.Content)
			}
			if sexpr.Name() == "" {
				sexpr.SetName(token.Content)
			} else {
				str := NewSexprStringQuoted(token.Content, false)
				str.SetLocation(token.Line, token.Column)
				str.SetParent(sexpr)
				param, err := NewSexprParam(str)
				if err != nil {
					return nil, err
				}
				sexpr.AddParam(len(sexpr.Params()), param)
			}

		} else if token.Kind == TokenQuotedString {
			if sexpr == nil {
				return nil, fmt.Errorf("unexpected quoted string at Line %d, Column %d: '%s'", token.Line, token.Column, token.Content)
			}
			if sexpr.Name() == "" {
				return nil, fmt.Errorf("unexpected quoted string at Line %d, Column %d: '%s'", token.Line, token.Column, token.Content)
			}
			str := NewSexprStringQuoted(token.Content[1:len(token.Content)-1], true)
			str.SetLocation(token.Line, token.Column)
			str.SetParent(sexpr)
			param, err := NewSexprParam(str)
			if err != nil {
				return nil, err
			}
			sexpr.AddParam(len(sexpr.Params()), param)

		} else if token.Kind == TokenEOF {
			if sexpr != nil {
				return nil, fmt.Errorf("unexpected EOF at Line %d, Column %d", token.Line, token.Column)
			}
			return root, nil

		} else if token.Kind == TokenErr {
			return nil, fmt.Errorf("error at Line %d, Column %d: %s", token.Line, token.Column, token.Err.Error())

		}
	}
}
