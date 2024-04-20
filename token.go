package sexpr

import "fmt"

type TokenKind int

const (
	TokenWhitespace TokenKind = iota
	TokenOpen
	TokenClose
	TokenString
	TokenQuotedString
	TokenEOF
	TokenErr
)

type Token struct {
	Kind    TokenKind
	Content string
	Line    int
	Column  int
	Err     error
}

func (t *Token) String() string {
	var s string
	switch t.Kind {
	case TokenWhitespace:
		s = "WS"
	case TokenOpen:
		s = "OPEN"
	case TokenClose:
		s = "CLOSE"
	case TokenString:
		s = "STRING"
	case TokenQuotedString:
		s = "QSTRING"
	case TokenEOF:
		s = "EOF"
	case TokenErr:
		s = "ERROR"
	}

	if t.Kind == TokenErr {
		return s + fmt.Sprintf(" %d:%d %s", t.Line, t.Column, t.Err.Error())
	} else {
		return s + fmt.Sprintf(" %d:%d %s", t.Line, t.Column, t.Content)
	}
}
