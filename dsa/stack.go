package dsa

type Stack[T any] []T

// Len returns the number of elements in the stack
func (s *Stack[T]) Len() int {
	return len(*s)
}

// IsEmpty returns true if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

// Push adds an item to the top of the stack
func (s *Stack[T]) Push(item T) {
	*s = append(*s, item)
}

// Pop removes and returns the top item from the stack
// Returns zero value of T and false if stack is empty
func (s *Stack[T]) Pop() (item T, ok bool) {
	if s.IsEmpty() {
		return item, false
	}
	item = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item, true
}

// Peek returns the top item without removing it
// Returns zero value of T and false if stack is empty
func (s *Stack[T]) Peek() (item T, ok bool) {
	if s.IsEmpty() {
		return item, false
	}
	return (*s)[len(*s)-1], true
}

// Clear removes all items from the stack
func (s *Stack[T]) Clear() {
	*s = (*s)[:0]
}
