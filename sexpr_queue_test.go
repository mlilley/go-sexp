package sexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := SexprQueue{}

	s1 := NewSexpr("One")
	s2 := NewSexpr("Two")

	q.Enqueue(s1)
	q.Enqueue(s2)

	require.Equal(t, q.Len(), 2)

	one := q.Dequeue()
	two := q.Dequeue()

	require.Equal(t, q.Len(), 0)
	require.Equal(t, "One", one.Name())
	require.Equal(t, "Two", two.Name())
}

func TestDequeueEmpty(t *testing.T) {
	q := SexprQueue{}

	require.Equal(t, q.Len(), 0)
	require.Nil(t, q.Dequeue())
	require.Equal(t, q.Len(), 0)
}
