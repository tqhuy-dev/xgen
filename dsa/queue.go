package dsa

type Queue[T any] []T

// Len returns the number of elements in the queue
func (q *Queue[T]) Len() int {
	return len(*q)
}

// IsEmpty returns true if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return len(*q) == 0
}

// Push adds an item to the back of the queue (enqueue)
func (q *Queue[T]) Push(item T) {
	*q = append(*q, item)
}

// Pop removes and returns the front item from the queue (dequeue)
// Returns zero value of T and false if queue is empty
func (q *Queue[T]) Pop() (item T, ok bool) {
	if q.IsEmpty() {
		return item, false
	}
	item = (*q)[0]
	*q = (*q)[1:]
	return item, true
}

// Peek returns the front item without removing it
// Returns zero value of T and false if queue is empty
func (q *Queue[T]) Peek() (item T, ok bool) {
	if q.IsEmpty() {
		return item, false
	}
	return (*q)[0], true
}

// Clear removes all items from the queue
func (q *Queue[T]) Clear() {
	*q = (*q)[:0]
}
