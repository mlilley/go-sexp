package sexpr

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBasic(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b "c c" #$% 1 2.3)`)))

	require.NoError(t, err)
	assertSexpr(t, root, "a", 5)
	assertStringParam(t, root.Params()[0], "b", false)
	assertStringParam(t, root.Params()[1], "c c", true)
	assertStringParam(t, root.Params()[2], "#$%", false)
	assertStringParam(t, root.Params()[3], "1", false)
	assertStringParam(t, root.Params()[4], "2.3", false)
}

func TestParseNested(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b (c d))`)))

	require.Nil(t, err)
	assertSexpr(t, root, "a", 2)
	assertStringParam(t, root.Params()[0], "b", false)
	assertSexprParam(t, root.Params()[1], "c", 1)

	nested, _ := root.Params()[1].Value().(*Sexpr)
	assertStringParam(t, nested.Params()[0], "d", false)
}

func TestParseEmpty(t *testing.T) {
	result, err := Parse(bufio.NewReader(strings.NewReader("")))
	require.Nil(t, err)
	require.Nil(t, result)

	result, err = Parse(bufio.NewReader(strings.NewReader("   \t  ")))
	require.Nil(t, err)
	require.Nil(t, result)
}

func TestParseSexprEmpty(t *testing.T) {
	_, err := Parse(bufio.NewReader(strings.NewReader("()")))
	require.ErrorContains(t, err, "unexpected close")

	_, err = Parse(bufio.NewReader(strings.NewReader("   ( \t )  ")))
	require.ErrorContains(t, err, "unexpected close")
}

func TestParseSexprAsName(t *testing.T) {
	_, err := Parse(bufio.NewReader(strings.NewReader("(a ((b)))")))
	require.ErrorContains(t, err, "unexpected open")
}

func TestParseSexprNoParams(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader("(a)")))
	require.Nil(t, err)
	require.NotNil(t, root)
	require.Equal(t, "a", root.Name())

	root, err = Parse(bufio.NewReader(strings.NewReader("   ( b )  ")))
	require.Nil(t, err)
	require.NotNil(t, root)
	require.Equal(t, "b", root.Name())
}

func TestParseNotNested(t *testing.T) {
	_, err := Parse(bufio.NewReader(strings.NewReader("(a)(b)")))
	require.ErrorContains(t, err, "unexpected open")
}

func TestSerialize(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b "c c" #$% 1 2.3)`)))
	require.Nil(t, err)
	require.NotNil(t, root)
	require.Equal(t, `(a b "c c" #$% 1 2.3)`, root.String())

	root, err = Parse(bufio.NewReader(strings.NewReader(`(a (b c) (d e))`)))
	require.Nil(t, err)
	require.NotNil(t, root)
	require.Equal(t, "(a\n\t(b c)\n\t(d e)\n)", root.String())

	root, err = Parse(bufio.NewReader(strings.NewReader(`(a (b (c d)))`)))
	require.Nil(t, err)
	require.NotNil(t, root)
	require.Equal(t, "(a\n\t(b\n\t\t(c d)\n\t)\n)", root.String())
}

// ---

func assertSexpr(t *testing.T, s *Sexpr, name string, params int) {
	require.NotNil(t, s)
	require.Equal(t, name, s.Name())
	require.Equal(t, params, len(s.Params()))
}

func assertStringParam(t *testing.T, sp *SexprParam, v string, quoted bool) {
	ss, ok := sp.Value().(*SexprString)
	require.True(t, ok)
	require.Equal(t, v, ss.Value())
	require.Equal(t, quoted, ss.Quoted())
}

func assertSexprParam(t *testing.T, sp *SexprParam, name string, params int) {
	s, ok := sp.Value().(*Sexpr)
	require.True(t, ok)
	require.Equal(t, name, s.Name())
	require.Equal(t, params, len(s.Params()))
}
