package dsa

type LinkListNode[T any] struct {
	Data T
	Next *LinkListNode[T]
}

type LinkList[T any] struct {
	Head *LinkListNode[T]
	Tail *LinkListNode[T]
	Len  int
}

// NewLinkList Init new link list
func NewLinkList[T any]() *LinkList[T] {
	return new(LinkList[T])
}

// AddNode Add node to link list
func (ll *LinkList[T]) AddNode(node *LinkListNode[T]) {
	if ll.Head == nil {
		ll.Head = node
		ll.Tail = node
	} else {
		ll.Tail.Next = node
		ll.Tail = node
	}
	ll.Len += 1
}

func (ll *LinkList[T]) Iteration(handle func(node *LinkListNode[T]) (*LinkListNode[T], bool, error)) (*LinkListNode[T], bool, error) {
	currentNode := ll.Head
	for currentNode != nil {
		nextNode, stop, err := handle(currentNode)
		if err != nil {
			return nil, true, err
		}
		if stop {
			return nextNode, true, nil
		}
		currentNode = currentNode.Next
	}
	return currentNode, true, nil
}

// Reverse link list
// 1 - 3 - 5 - 7
// 7 - 5- 3 - 1
func (ll *LinkList[T]) Reverse() {
	cursorNode := ll.Head
	var prevNode *LinkListNode[T]
	for cursorNode != nil {
		nextNode := cursorNode.Next
		cursorNode.Next = prevNode
		prevNode = cursorNode
		cursorNode = nextNode
	}
	ll.Tail = ll.Head
	ll.Head = prevNode
}
