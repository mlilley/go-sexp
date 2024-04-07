package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

type TokenKind int

const (
	TokenKindQuoteOpen TokenKind = iota
	TokenKindQuoteClose
	TokenKindParam
)

type Token struct {
	Kind  TokenKind
	Value string
}

// An SParam can be of type string, int, float64, SExp
type SParam interface{}

// An SExp has zero or more SParams
type SExp struct {
	params []SParam
}

func (sexp *SExp) String() string {
	return "TODO"
}

const MAX_PARAM_BYTES int = 1024

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

	sexp, err := parseSexp(bufio.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}

	log.Print(sexp.String())
}

func parseSexp(reader *bufio.Reader) (*SExp, error) {
	inSexp := false
	sexp := SExp{}

	for {
		token, err := nextToken(reader)
		if err != nil {
			if err == io.EOF {
				// an EOF illegally terminates an sexp
				if inSexp {
					return nil, errors.New("unexpected EOF reading sexp")
				}
			}
			return nil, err
		}

		if inSexp {
			switch t := token.(type) {
			case rune:
				// proper termination of sexp
				if t == ')' {
					return &sexp, nil
				}
				// param is a sexp, recurse into it
				if t == '(' {
					sexpParam, err := parseSexp(reader)
					if err != nil {
						if err == io.EOF {
							return nil, errors.New("unexpected EOF reading sexp")
						}
						return nil, err
					}
					sexp.params = append(sexp.params, sexpParam)
					continue
				}
				// unexpected rune
				return nil, errors.New("unexpected token (rune)")

			default:
				// token is string, int, or float64
				sexp.params = append(sexp.params, token)
				continue
			}
		}

		switch t := token.(type) {
		case rune:
			// only valid token outside of a sexp
			if t == '(' {
				inSexp = true
				continue
			}
		}

		return nil, fmt.Errorf("unexpected token (%q:%T) reading sexp", token, token)
	}
}

// returns rune, string, int, or float
func nextToken(reader *bufio.Reader) (interface{}, error) {
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
				return token, nil
			}

			// append rune to quoted token
			token += string(r)
			continue
		}

		if inUnquoted {
			if r == '(' {
				// illegal sexp start
				return nil, errors.New("illegal sexp start")
			}
			if r == '"' {
				// illegal quoted start
				return nil, errors.New("illegal quoted start")
			}
			if r == ')' {
				// push back the ), end and return the unquoted token
				err = reader.UnreadRune()
				return token, nil
			}
			if unicode.IsSpace(r) {
				// end and return the unquoted token
				return token, nil
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
