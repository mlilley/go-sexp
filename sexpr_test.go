package sexpr

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChildAnyDepth(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	assertNoError(t, err)

	var sexpr *Sexpr

	// doesn't find root
	sexpr = root.FindChildByName("aaa", -1)
	assert.Nil(t, sexpr)

	// finds 2nd ccc (first in breadth-first)
	sexpr = root.FindChildByName("ccc", -1)
	assert.Equal(t, sexpr.Name(), "ccc")
	assert.Equal(t, len(sexpr.Params()), 2)

	// finds more than 1 level deep
	sexpr = root.FindChildByName("eee", -1)
	assert.NotNil(t, sexpr)
	assert.Equal(t, sexpr.Name(), "eee")

	// returns nil where no child exists by the name
	sexpr = root.FindChildByName("doesnotexist", -1)
	assert.Nil(t, sexpr)
}

func TestGetChildMaxDepth(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	assertNoError(t, err)

	var sexpr *Sexpr

	// doesn't find root
	sexpr = root.FindChildByName("aaa", 1)
	assert.Nil(t, sexpr)

	// finds 1st level
	sexpr = root.FindChildByName("bbb", 1)
	assert.Equal(t, sexpr.Name(), "bbb")

	// does not find at depth greater than max
	sexpr = root.FindChildByName("eee", 1)
	assert.Nil(t, sexpr)

	// returns nil wher eno child exists by the name
	sexpr = root.FindChildByName("doesnotexist", 1)
	assert.Nil(t, sexpr)
}

func TestGetChildrenByName(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	assertNoError(t, err)

	var sexprs []*Sexpr

	// doesn't find root
	sexprs = root.FindChildrenByName("aaa", -1)
	assert.Equal(t, len(sexprs), 0)

	// finds at multiple levels
	sexprs = root.FindChildrenByName("ccc", -1)
	assert.Equal(t, len(sexprs), 2)
	assert.Equal(t, sexprs[0].Name(), "ccc")
	assert.Equal(t, len(sexprs[0].Params()), 2)
	assert.Equal(t, sexprs[1].Name(), "ccc")
	assert.Equal(t, len(sexprs[1].Params()), 0)
}
