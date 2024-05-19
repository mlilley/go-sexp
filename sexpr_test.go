package sexpr

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetByNameFindPredicate(name string) FindPredicate {
	return func(sexpr *Sexpr, depth int) bool {
		if strings.EqualFold(sexpr.Name(), name) {
			return true
		} else {
			return false
		}
	}
}

func TestFindChild(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	require.NoError(t, err)

	var sexpr *Sexpr

	// finds 2nd ccc (first in breadth-first)
	sexpr = root.FindChild(GetByNameFindPredicate("ccc"), -1)
	require.Equal(t, "ccc", sexpr.Name())
	require.Equal(t, 2, len(sexpr.Params()))

	// finds more than 1 level deep
	sexpr = root.FindChild(GetByNameFindPredicate("eee"), -1)
	require.Equal(t, "eee", sexpr.Name())

	// doesn't find root
	sexpr = root.FindChild(GetByNameFindPredicate("aaa"), -1)
	require.Nil(t, sexpr)
}

func TestFindChildren(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	require.NoError(t, err)

	var sexprs []*Sexpr

	// finds all ccc's
	sexprs = root.FindChildren(GetByNameFindPredicate("ccc"), -1)
	require.Equal(t, 2, len(sexprs))
	require.Equal(t, "ccc", sexprs[0].Name())
	require.Equal(t, 2, len(sexprs[0].Params()))
	require.Equal(t, "ccc", sexprs[1].Name())
	require.Equal(t, 0, len(sexprs[1].Params()))

	// finds more than 1 level deep
	sexprs = root.FindChildren(GetByNameFindPredicate("eee"), -1)
	require.Equal(t, 1, len(sexprs))
	require.Equal(t, "eee", sexprs[0].Name())

	// doesn't find root
	sexprs = root.FindChildren(GetByNameFindPredicate("aaa"), -1)
	require.Equal(t, 0, len(sexprs))
}

func TestFindChildByName(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	require.NoError(t, err)

	sexpr := root.FindChildByName("eee", -1)
	require.Equal(t, sexpr.Name(), "eee")
}

func TestFindChildrenByName(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	require.NoError(t, err)

	sexprs := root.FindChildrenByName("ccc", -1)
	require.Equal(t, 2, len(sexprs))
	require.Equal(t, "ccc", sexprs[0].Name())
	require.Equal(t, 2, len(sexprs[0].Params()))
	require.Equal(t, "ccc", sexprs[1].Name())
	require.Equal(t, 0, len(sexprs[1].Params()))
}

func TestFindDirectChildByName(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	require.NoError(t, err)

	sexpr := root.FindDirectChildByName("eee")
	require.Nil(t, sexpr)
}

func TestFindDirectChildrenByName(t *testing.T) {
	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa (bbb (ccc)) (ccc a b) (ddd (eee)))`)))
	require.NoError(t, err)

	sexprs := root.FindDirectChildrenByName("ccc")
	require.Equal(t, 1, len(sexprs))
	require.Equal(t, "ccc", sexprs[0].Name())
	require.Equal(t, 2, len(sexprs[0].Params()))
}

// func TestSetParam(t *testing.T) {
// 	root, err := Parse(bufio.NewReader(strings.NewReader(`(aaa bbb ccc ddd)`)))
// 	assertNoError(t, err)

// 	newParam := NewSexprParam(NewSexprString("eee", root, -1, -1))

// 	root.SetParam(0, newParam)
// }
