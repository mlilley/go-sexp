package sexpr

type SexprQueue struct {
	head *sexprQueueNode
	tail *sexprQueueNode
	len  int
}

type sexprQueueNode struct {
	value *Sexpr
	next  *sexprQueueNode
}

func NewSexprQueue() *SexprQueue {
	return &SexprQueue{}
}

func (sq *SexprQueue) Enqueue(s *Sexpr) {
	node := sexprQueueNode{value: s, next: nil}
	if sq.tail != nil {
		sq.tail.next = &node
	} else {
		sq.head = &node
	}
	sq.tail = &node
	sq.len += 1
}

func (sq *SexprQueue) Dequeue() *Sexpr {
	if sq.head == nil {
		return nil
	}
	head := sq.head
	sq.head = sq.head.next
	if sq.head == nil {
		sq.tail = nil
	}
	sq.len -= 1
	return head.value
}

func (sq *SexprQueue) Len() int {
	return sq.len
}
