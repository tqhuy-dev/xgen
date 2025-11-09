package dsa

import "container/heap"

type HeapNode[T any] struct {
	Data  T
	Point int
}

type HeapList[T any] []HeapNode[T]

func (h HeapList[T]) Len() int {
	return len(h)
}

func (h HeapList[T]) Less(i, j int) bool {
	return h[i].Point < h[j].Point
}

func (h HeapList[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *HeapList[T]) Push(x any) {
	*h = append(*h, x.(HeapNode[T]))
}

func (h *HeapList[T]) PushNode(point int, data T) {}

func (h *HeapList[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TopMaxPoint[T any](data []HeapNode[T], nPoints int) []HeapNode[T] {
	h := make(HeapList[T], 0)
	heap.Init(&h)

	for _, i := range data {
		heap.Push(&h, i)
	}
	result := make([]HeapNode[T], h.Len())
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = heap.Pop(&h).(HeapNode[T])
	}
	return result[:nPoints]
}
