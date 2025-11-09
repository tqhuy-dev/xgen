package dsa

// DLLNode is a node in a doubly linked list with generic type T
type DLLNode[T any] struct {
	next, prev *DLLNode[T]
	list       *DLList[T]
	Value      T
}

// Next returns the next list element or nil
func (e *DLLNode[T]) Next() *DLLNode[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil
func (e *DLLNode[T]) Prev() *DLLNode[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// DLList is a doubly linked list with generic type T
type DLList[T any] struct {
	root DLLNode[T]
	len  int
}

// Init initializes or clears list l
func (l *DLList[T]) Init() *DLList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewDLList returns an initialized doubly linked list
func NewDLList[T any]() *DLList[T] { return new(DLList[T]).Init() }

// Len returns the number of elements in the list
func (l *DLList[T]) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty
func (l *DLList[T]) Front() *DLLNode[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty
func (l *DLList[T]) Back() *DLLNode[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value
func (l *DLList[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e
func (l *DLList[T]) insert(e, at *DLLNode[T]) *DLLNode[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&DLLNode[T]{Value: v}, at)
func (l *DLList[T]) insertValue(v T, at *DLLNode[T]) *DLLNode[T] {
	return l.insert(&DLLNode[T]{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *DLList[T]) remove(e *DLLNode[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	l.len--
}

// move moves e to next to at
func (l *DLList[T]) move(e, at *DLLNode[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes e from l if e is an element of list l and returns the element value
func (l *DLList[T]) Remove(e *DLLNode[T]) T {
	if e.list == l {
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e
func (l *DLList[T]) PushFront(v T) *DLLNode[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e
func (l *DLList[T]) PushBack(v T) *DLLNode[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e
func (l *DLList[T]) InsertBefore(v T, mark *DLLNode[T]) *DLLNode[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e
func (l *DLList[T]) InsertAfter(v T, mark *DLLNode[T]) *DLLNode[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l
func (l *DLList[T]) MoveToFront(e *DLLNode[T]) {
	if e.list != l || l.root.next == e {
		return
	}
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l
func (l *DLList[T]) MoveToBack(e *DLLNode[T]) {
	if e.list != l || l.root.prev == e {
		return
	}
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark
func (l *DLList[T]) MoveBefore(e, mark *DLLNode[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark
func (l *DLList[T]) MoveAfter(e, mark *DLLNode[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of another list at the back of list l
func (l *DLList[T]) PushBackList(other *DLList[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l
func (l *DLList[T]) PushFrontList(other *DLList[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
