package sexpr

import (
	"bufio"
	"strings"
	"testing"
)

func TestGetChildAnyDepth(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b (c d (e f g)))`)))
	assertNoError(t, err)

	var sexpr *Sexpr

	sexpr = root.GetChildByName("a", -1)
	if sexpr != nil {
		t.Fatal()
	}

	sexpr = root.GetChildByName("c", -1)
	if sexpr == nil {
		t.Fatal()
	}
	if sexpr.Name() != "c" {
		t.Fatal()
	}

	sexpr = root.GetChildByName("e", -1)
	if sexpr == nil {
		t.Fatal()
	}
	if sexpr.Name() != "e" {
		t.Fatal()
	}

	sexpr = root.GetChildByName("doesnotexist", -1)
	if sexpr != nil {
		t.Fatal()
	}
}

func TestGetChildMaxDepth(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(a b (c d (e f g)))`)))
	assertNoError(t, err)

	var sexpr *Sexpr

	sexpr = root.GetChildByName("a", 1)
	if sexpr != nil {
		t.Fatal()
	}

	sexpr = root.GetChildByName("c", 1)
	if sexpr == nil {
		t.Fatal()
	}
	if sexpr.Name() != "c" {
		t.Fatal()
	}

	sexpr = root.GetChildByName("e", 1)
	if sexpr != nil {
		t.Fatal()
	}

	sexpr = root.GetChildByName("doesnotexist", 1)
	if sexpr != nil {
		t.Fatal()
	}
}
