package main

import (
	"bufio"
	"fmt"
)

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
			s := &Sexpr{
				Line:   token.Line,
				Column: token.Column,
				Parent: sexpr,
			}
			if depth == 0 {
				if root != nil {
					return nil, fmt.Errorf("unexpected open at Line %d, Column %d", token.Line, token.Column)
				}
				root = s
			}
			if depth > 0 && !hasName {
				return nil, fmt.Errorf("unexpected open at Line %d, Column %d", token.Line, token.Column)
			}
			hasName = false
			sexpr = s
			depth++

		} else if token.Kind == TokenClose {
			if depth == 0 {
				return nil, fmt.Errorf("unexpected close at Line %d, Column %d", token.Line, token.Column)
			}
			if !hasName {
				return nil, fmt.Errorf("sexpression without name at Line %d, Column %d", token.Line, token.Column)
			}
			if depth > 1 {
				sexpr.Parent.AddParam(sexpr)
			}
			sexpr = sexpr.Parent
			depth--

		} else if token.Kind == TokenString {
			if depth == 0 {
				return nil, fmt.Errorf("unexpected string at Line %d, Column %d: '%s'", token.Line, token.Column, token.Content)
			}
			if hasName {
				sexpr.AddParam(NewSexprStringParamQuoted(token.Content, false, token.Line, token.Column))
			} else {
				sexpr.Name = token.Content
				hasName = true
			}

		} else if token.Kind == TokenQuotedString {
			if depth == 0 {
				return nil, fmt.Errorf("unexpected string at Line %d, Column %d: '%s'", token.Line, token.Column, token.Content)
			}
			if hasName {
				sexpr.AddParam(NewSexprStringParamQuoted(token.Content[1:len(token.Content)-1], true, token.Line, token.Column))
			} else {
				sexpr.Name = token.Content
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
