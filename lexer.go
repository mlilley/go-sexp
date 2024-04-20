package sexpr

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	line        int
	column      int
	prevLine    int
	prevColumn  int
	startLine   int
	startColumn int
	content     string
	input       *bufio.Reader
}

func NewLexer(input *bufio.Reader) *Lexer {
	return &Lexer{
		line:        1,
		column:      1,
		prevLine:    1,
		prevColumn:  1,
		startLine:   1,
		startColumn: 1,
		content:     "",
		input:       input,
	}
}

func (l *Lexer) NextToken(token *Token) {
	l.startLine = l.line
	l.startColumn = l.column
	l.content = ""

	r, err := l.read()
	if err == io.EOF {
		l.emit(token, TokenEOF, nil)
	} else if err != nil {
		l.emit(token, TokenErr, err)
	} else if r == '(' {
		l.emit(token, TokenOpen, nil)
	} else if r == ')' {
		l.emit(token, TokenClose, nil)
	} else if unicode.IsSpace(r) {
		l.acceptWhitespace()
		l.emit(token, TokenWhitespace, nil)
	} else if r == '"' {
		err = l.acceptQuotedString()
		if err != nil {
			l.emit(token, TokenErr, err)
			return
		}
		l.emit(token, TokenQuotedString, nil)
	} else {
		err = l.acceptString()
		if err != nil {
			l.emit(token, TokenErr, err)
			return
		}
		l.emit(token, TokenString, nil)
	}
}

func (l *Lexer) read() (rune, error) {
	r, _, err := l.input.ReadRune()
	if err != nil {
		return 0, err
	}
	l.prevLine = l.line
	l.prevColumn = l.column
	if r == '\n' {
		l.line += 1
		l.column = 1
	} else {
		l.column += 1
	}
	l.content += string(r)
	return r, nil
}

func (l *Lexer) unread() error {
	if l.column == l.prevColumn && l.line == l.prevLine {
		return errors.New("unable to unread")
	}
	err := l.input.UnreadRune()
	if err != nil {
		return err
	}
	l.line = l.prevLine
	l.column = l.prevColumn
	_, lastRuneSize := utf8.DecodeLastRuneInString(l.content)
	l.content = l.content[:len(l.content)-lastRuneSize]
	return nil
}

func (l *Lexer) emit(token *Token, kind TokenKind, err error) {
	token.Kind = kind
	token.Line = l.startLine
	token.Column = l.startColumn
	token.Content = l.content
	token.Err = err
}

func (l *Lexer) acceptWhitespace() error {
	for {
		r, err := l.read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			return l.unread()
		}
	}
}

func (l *Lexer) acceptQuotedString() error {
	for {
		r, err := l.read()
		if err == io.EOF {
			return fmt.Errorf("unterminated quoted string")
		}
		if err != nil {
			return err
		}
		if r == '"' {
			return nil
		}
	}
}

func (l *Lexer) acceptString() error {
	for {
		r, err := l.read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if r == '(' || r == ')' || unicode.IsSpace(r) {
			return l.unread()
		}
	}
}
