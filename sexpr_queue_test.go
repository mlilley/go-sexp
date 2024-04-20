package sexpr

import "testing"

func TestQueue(t *testing.T) {
	q := SexprQueue{}

	q.Enqueue(NewSexpr("One", nil, nil, 0, 0))
	q.Enqueue(NewSexpr("Two", nil, nil, 0, 0))
	if q.Len() != 2 {
		t.Fail()
	}

	one := q.Dequeue()
	two := q.Dequeue()

	if q.Len() != 0 {
		t.Fail()
	}
	if one.Name() != "One" {
		t.Fail()
	}
	if two.Name() != "Two" {
		t.Fail()
	}
}

func TestDequeueEmpty(t *testing.T) {
	q := SexprQueue{}

	if q.Len() != 0 {
		t.Fail()
	}

	s := q.Dequeue()
	if s != nil {
		t.Fail()
	}

	if q.Len() != 0 {
		t.Fail()
	}
}
