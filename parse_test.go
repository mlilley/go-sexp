package sexpr

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b "c c" #$% 1 2.3)`)))

	assertNoError(t, err)
	assertSexpr(t, root, "a", 5)
	assertStringParam(t, root.Params()[0], "b", false)
	assertStringParam(t, root.Params()[1], "c c", true)
	assertStringParam(t, root.Params()[2], "#$%", false)
	assertStringParam(t, root.Params()[3], "1", false)
	assertStringParam(t, root.Params()[4], "2.3", false)
}

func TestParseNested(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b (c d))`)))

	assertNoError(t, err)
	assertSexpr(t, root, "a", 2)
	assertStringParam(t, root.Params()[0], "b", false)
	assertSexprParam(t, root.Params()[1], "c", 1)

	nested, _ := root.Params()[1].Value().(*Sexpr)
	assertStringParam(t, nested.Params()[0], "d", false)
}

func TestParseEmpty(t *testing.T) {
	result, err := Parse(bufio.NewReader(strings.NewReader("")))
	assertNoError(t, err)
	assertTrue(t, result == nil)

	result, err = Parse(bufio.NewReader(strings.NewReader("   \t  ")))
	assertNoError(t, err)
	assertTrue(t, result == nil)
}

func TestParseSexprEmpty(t *testing.T) {
	_, err := Parse(bufio.NewReader(strings.NewReader("()")))
	assertError(t, err, "expression without name")

	_, err = Parse(bufio.NewReader(strings.NewReader("   ( \t )  ")))
	assertError(t, err, "expression without name")
}

func TestParseSexprAsName(t *testing.T) {
	_, err := Parse(bufio.NewReader(strings.NewReader("(a ((b)))")))
	assertError(t, err, "unexpected open")
}

func TestParseSexprNoParams(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader("(a)")))
	assertNoError(t, err)
	assertTrue(t, root != nil)
	assertEqualsStr(t, root.Name(), "a")

	root, err = Parse(bufio.NewReader(strings.NewReader("   ( b )  ")))
	assertNoError(t, err)
	assertTrue(t, root != nil)
	assertEqualsStr(t, root.Name(), "b")
}

func TestParseNotNested(t *testing.T) {
	_, err := Parse(bufio.NewReader(strings.NewReader("(a)(b)")))
	assertError(t, err, "unexpected open")
}

func TestSerialize(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b "c c" #$% 1 2.3)`)))
	assertNoError(t, err)
	assertTrue(t, root != nil)
	assertEqualsStr(t, root.String(), `(a b "c c" #$% 1 2.3)`)

	root, err = Parse(bufio.NewReader(strings.NewReader(`(a (b c) (d e))`)))
	assertNoError(t, err)
	assertTrue(t, root != nil)
	assertEqualsStr(t, root.String(), "(a\n\t(b c)\n\t(d e)\n)")

	root, err = Parse(bufio.NewReader(strings.NewReader(`(a (b (c d)))`)))
	assertNoError(t, err)
	assertTrue(t, root != nil)
	assertEqualsStr(t, root.String(), "(a\n\t(b\n\t\t(c d)\n\t)\n)")
}

// ---

func assertSexpr(t *testing.T, s *Sexpr, name string, params int) {
	assertTrue(t, s != nil)
	assertEqualsStr(t, s.Name(), name)
	assertEqualsInt(t, len(s.Params()), params)
}

func assertStringParam(t *testing.T, sp *SexprParam, v string, quoted bool) {
	ss, ok := sp.Value().(*SexprString)
	assertTrue(t, ok)
	assertEqualsStr(t, ss.Value(), v)
	assertEqualsBool(t, ss.Quoted(), quoted)
}

func assertSexprParam(t *testing.T, sp *SexprParam, name string, params int) {
	s, ok := sp.Value().(*Sexpr)
	assertTrue(t, ok)
	assertEqualsStr(t, s.Name(), name)
	assertEqualsInt(t, len(s.Params()), params)
}

func assertTrue(t *testing.T, actual bool) {
	if !actual {
		t.Error("expected true, but got false")
	}
}

func assertFalse(t *testing.T, actual bool) {
	if actual {
		t.Error("expected false, but got true")
	}
}

func assertEqualsStr(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("expected '%s', but got '%s'", expected, actual)
	}
}

func assertEqualsInt(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Errorf("expected '%d', but got '%d'", expected, actual)
	}
}

func assertEqualsBool(t *testing.T, actual bool, expected bool) {
	if actual != expected {
		t.Errorf("expected %t, but got %t", expected, actual)
	}
}

func assertError(t *testing.T, actual error, containing string) {
	if actual == nil {
		t.Errorf("expected error, but got nil")
	}
	if !strings.Contains(actual.Error(), containing) {
		t.Errorf("expected error containing '%s', but got '%s'", containing, actual.Error())
	}
}

func assertNoError(t *testing.T, actual error) {
	if actual != nil {
		t.Errorf("expected no error, but got error containing '%s'", actual.Error())
	}
}
