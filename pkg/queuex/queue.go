package queuex

import "sync"

type Node[T any] struct {
	Value T
	Next  *Node[T]
}
type Queue[T any] struct {
	Head  *Node[T]
	Tail  *Node[T]
	Size  int
	mutex *sync.RWMutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		Head:  nil,
		Tail:  nil,
		Size:  0,
		mutex: new(sync.RWMutex),
	}
}

// Enqueue adds an element to the end of the queue (tail).
func (q *Queue[T]) Enqueue(v T) {

	node := &Node[T]{
		Value: v,
		Next:  nil,
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.Size++

	if q.Tail == nil {
		q.Head = node
		q.Tail = node
		return
	}

	q.Tail.Next = node
	q.Tail = node
}

// Dequeue removes and returns the element at the front of the queue (head). FIFO
func (q *Queue[T]) Dequeue() (T, bool) {

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.Head != nil {

		v := q.Head.Value
		q.Head = q.Head.Next

		if q.Head == nil {
			q.Tail = nil
		}

		q.Size--

		return v, true
	}

	var zeroValue T
	return zeroValue, false
}

// IsEmpty checks if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.Size == 0
}

// Print to get all values each node.
func (q *Queue[T]) Print() []T {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	var values []T
	current := q.Head
	for current != nil {
		values = append(values, current.Value)
		current = current.Next
	}
	return values
}
