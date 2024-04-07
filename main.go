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

	sexp, err := parse(bufio.NewReader(f), false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sexp.String())
}

type SExp struct {
	params []SParam
}

type SParam interface{} // (string, int, float64, or SExp)

func (sexp *SExp) String() string {
	var sb strings.Builder
	stringSexp(sexp, &sb, 0)
	return sb.String()
}

func stringSexp(sexp *SExp, acc *strings.Builder, level int) {
	indent := strings.Repeat("  ", level)
	acc.WriteString(indent)
	acc.WriteString("(\n")
	for _, param := range sexp.params {
		switch t := param.(type) {
		case string:
			acc.WriteString(indent)
			acc.WriteString(t)
			acc.WriteString("\n")
		case int:
			acc.WriteString(indent)
			acc.WriteString(strconv.Itoa(t))
			acc.WriteString("\n")
		case float64:
			acc.WriteString(indent)
			acc.WriteString(strconv.FormatFloat(t, 'f', -1, 64))
			acc.WriteString("\n")
		case *SExp:
			stringSexp(t, acc, level+1)
		}
	}
	acc.WriteString(indent)
	acc.WriteString(")\n")
}

func parse(reader *bufio.Reader, parseNumerics bool) (*SExp, error) {
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
			sexp, err := parseSexp(reader, parseNumerics)
			if err != nil {
				return nil, err
			}
			return sexp, nil
		}
	}

	return nil, errors.New("unexpected token")
}

func parseSexp(reader *bufio.Reader, parseNumerics bool) (*SExp, error) {
	sexp := SExp{}

	for {
		token, err := nextToken(reader, parseNumerics)
		if err != nil {
			if err == io.EOF {
				// an EOF illegally terminates an sexp
				return nil, errors.New("unexpected EOF reading sexp")
			}
			return nil, err
		}

		switch t := token.(type) {
		case rune:
			// proper termination of sexp
			if t == ')' {
				return &sexp, nil
			}
			// param is a sexp, recurse into it
			if t == '(' {
				sexpParam, err := parseSexp(reader, parseNumerics)
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
}

// returns rune, string, int, or float
func nextToken(reader *bufio.Reader, parseNumerics bool) (interface{}, error) {
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
				if err != nil {
					return nil, err
				}

				if parseNumerics {
					if v, err := strconv.Atoi(token); err == nil {
						return v, nil
					}
					if v, err := strconv.ParseFloat(token, 64); err == nil {
						return v, nil
					}
				}
				return token, nil
			}
			if unicode.IsSpace(r) {
				// end and return the unquoted token
				if parseNumerics {
					if v, err := strconv.Atoi(token); err == nil {
						return v, nil
					}
					if v, err := strconv.ParseFloat(token, 64); err == nil {
						return v, nil
					}
				}
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
