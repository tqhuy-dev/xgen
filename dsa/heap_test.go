package dsa

import (
	"container/heap"
	"reflect"
	"testing"
)

// TestHeapList_Len tests the Len method of HeapList
// which returns the number of elements in the heap
func TestHeapList_Len(t *testing.T) {
	tests := []struct {
		name     string
		heap     HeapList[int]
		expected int
	}{
		{
			name:     "empty heap",
			heap:     HeapList[int]{},
			expected: 0,
		},
		{
			name: "heap with one element",
			heap: HeapList[int]{
				{Data: 1, Point: 5},
			},
			expected: 1,
		},
		{
			name: "heap with multiple elements",
			heap: HeapList[int]{
				{Data: 1, Point: 5},
				{Data: 2, Point: 10},
				{Data: 3, Point: 3},
			},
			expected: 3,
		},
		{
			name: "heap with many elements",
			heap: HeapList[int]{
				{Data: 1, Point: 1},
				{Data: 2, Point: 2},
				{Data: 3, Point: 3},
				{Data: 4, Point: 4},
				{Data: 5, Point: 5},
				{Data: 6, Point: 6},
				{Data: 7, Point: 7},
			},
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.heap.Len()

			// Assert
			if result != tt.expected {
				t.Errorf("HeapList.Len() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestHeapList_Less tests the Less method of HeapList
// which compares two elements based on their Point values (min-heap)
func TestHeapList_Less(t *testing.T) {
	tests := []struct {
		name     string
		heap     HeapList[string]
		i        int
		j        int
		expected bool
	}{
		{
			name: "first element less than second",
			heap: HeapList[string]{
				{Data: "a", Point: 5},
				{Data: "b", Point: 10},
			},
			i:        0,
			j:        1,
			expected: true,
		},
		{
			name: "first element greater than second",
			heap: HeapList[string]{
				{Data: "a", Point: 15},
				{Data: "b", Point: 10},
			},
			i:        0,
			j:        1,
			expected: false,
		},
		{
			name: "equal points",
			heap: HeapList[string]{
				{Data: "a", Point: 10},
				{Data: "b", Point: 10},
			},
			i:        0,
			j:        1,
			expected: false,
		},
		{
			name: "negative points comparison",
			heap: HeapList[string]{
				{Data: "a", Point: -5},
				{Data: "b", Point: -10},
			},
			i:        0,
			j:        1,
			expected: false,
		},
		{
			name: "zero and positive comparison",
			heap: HeapList[string]{
				{Data: "a", Point: 0},
				{Data: "b", Point: 5},
			},
			i:        0,
			j:        1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.heap.Less(tt.i, tt.j)

			// Assert
			if result != tt.expected {
				t.Errorf("HeapList.Less(%d, %d) = %v; expected %v (Point[%d]=%d, Point[%d]=%d)",
					tt.i, tt.j, result, tt.expected, tt.i, tt.heap[tt.i].Point, tt.j, tt.heap[tt.j].Point)
			}
		})
	}
}

// TestHeapList_Swap tests the Swap method of HeapList
// which exchanges the positions of two elements in the heap
func TestHeapList_Swap(t *testing.T) {
	tests := []struct {
		name     string
		heap     HeapList[int]
		i        int
		j        int
		expected HeapList[int]
	}{
		{
			name: "swap first and second elements",
			heap: HeapList[int]{
				{Data: 1, Point: 5},
				{Data: 2, Point: 10},
			},
			i: 0,
			j: 1,
			expected: HeapList[int]{
				{Data: 2, Point: 10},
				{Data: 1, Point: 5},
			},
		},
		{
			name: "swap in middle of heap",
			heap: HeapList[int]{
				{Data: 1, Point: 5},
				{Data: 2, Point: 10},
				{Data: 3, Point: 15},
			},
			i: 0,
			j: 2,
			expected: HeapList[int]{
				{Data: 3, Point: 15},
				{Data: 2, Point: 10},
				{Data: 1, Point: 5},
			},
		},
		{
			name: "swap same element",
			heap: HeapList[int]{
				{Data: 1, Point: 5},
				{Data: 2, Point: 10},
			},
			i: 0,
			j: 0,
			expected: HeapList[int]{
				{Data: 1, Point: 5},
				{Data: 2, Point: 10},
			},
		},
		{
			name: "swap adjacent middle elements",
			heap: HeapList[int]{
				{Data: 1, Point: 5},
				{Data: 2, Point: 10},
				{Data: 3, Point: 15},
				{Data: 4, Point: 20},
			},
			i: 1,
			j: 2,
			expected: HeapList[int]{
				{Data: 1, Point: 5},
				{Data: 3, Point: 15},
				{Data: 2, Point: 10},
				{Data: 4, Point: 20},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			heapCopy := make(HeapList[int], len(tt.heap))
			copy(heapCopy, tt.heap)

			// Act
			heapCopy.Swap(tt.i, tt.j)

			// Assert
			if !reflect.DeepEqual(heapCopy, tt.expected) {
				t.Errorf("HeapList.Swap(%d, %d) resulted in %v; expected %v",
					tt.i, tt.j, heapCopy, tt.expected)
			}
		})
	}
}

// TestHeapList_Push tests the Push method of HeapList
// which adds a new element to the end of the heap
func TestHeapList_Push(t *testing.T) {
	tests := []struct {
		name     string
		initial  HeapList[string]
		toPush   HeapNode[string]
		expected HeapList[string]
	}{
		{
			name:    "push to empty heap",
			initial: HeapList[string]{},
			toPush:  HeapNode[string]{Data: "first", Point: 10},
			expected: HeapList[string]{
				{Data: "first", Point: 10},
			},
		},
		{
			name: "push to heap with one element",
			initial: HeapList[string]{
				{Data: "first", Point: 5},
			},
			toPush: HeapNode[string]{Data: "second", Point: 10},
			expected: HeapList[string]{
				{Data: "first", Point: 5},
				{Data: "second", Point: 10},
			},
		},
		{
			name: "push to heap with multiple elements",
			initial: HeapList[string]{
				{Data: "first", Point: 5},
				{Data: "second", Point: 10},
				{Data: "third", Point: 15},
			},
			toPush: HeapNode[string]{Data: "fourth", Point: 7},
			expected: HeapList[string]{
				{Data: "first", Point: 5},
				{Data: "second", Point: 10},
				{Data: "third", Point: 15},
				{Data: "fourth", Point: 7},
			},
		},
		{
			name: "push element with negative point",
			initial: HeapList[string]{
				{Data: "first", Point: 0},
			},
			toPush: HeapNode[string]{Data: "negative", Point: -5},
			expected: HeapList[string]{
				{Data: "first", Point: 0},
				{Data: "negative", Point: -5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			heapCopy := make(HeapList[string], len(tt.initial))
			copy(heapCopy, tt.initial)

			// Act
			heapCopy.Push(tt.toPush)

			// Assert
			if !reflect.DeepEqual(heapCopy, tt.expected) {
				t.Errorf("HeapList.Push() resulted in %v; expected %v", heapCopy, tt.expected)
			}
		})
	}
}

// TestHeapList_Pop tests the Pop method of HeapList
// which removes and returns the last element from the heap
func TestHeapList_Pop(t *testing.T) {
	tests := []struct {
		name         string
		initial      HeapList[int]
		expectedPop  HeapNode[int]
		expectedHeap HeapList[int]
	}{
		{
			name: "pop from heap with one element",
			initial: HeapList[int]{
				{Data: 100, Point: 10},
			},
			expectedPop:  HeapNode[int]{Data: 100, Point: 10},
			expectedHeap: HeapList[int]{},
		},
		{
			name: "pop from heap with two elements",
			initial: HeapList[int]{
				{Data: 100, Point: 5},
				{Data: 200, Point: 10},
			},
			expectedPop: HeapNode[int]{Data: 200, Point: 10},
			expectedHeap: HeapList[int]{
				{Data: 100, Point: 5},
			},
		},
		{
			name: "pop from heap with multiple elements",
			initial: HeapList[int]{
				{Data: 100, Point: 5},
				{Data: 200, Point: 10},
				{Data: 300, Point: 15},
				{Data: 400, Point: 20},
			},
			expectedPop: HeapNode[int]{Data: 400, Point: 20},
			expectedHeap: HeapList[int]{
				{Data: 100, Point: 5},
				{Data: 200, Point: 10},
				{Data: 300, Point: 15},
			},
		},
		{
			name: "pop from heap with negative point",
			initial: HeapList[int]{
				{Data: 100, Point: 5},
				{Data: 200, Point: -10},
			},
			expectedPop: HeapNode[int]{Data: 200, Point: -10},
			expectedHeap: HeapList[int]{
				{Data: 100, Point: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			heapCopy := make(HeapList[int], len(tt.initial))
			copy(heapCopy, tt.initial)

			// Act
			popped := heapCopy.Pop().(HeapNode[int])

			// Assert
			if !reflect.DeepEqual(popped, tt.expectedPop) {
				t.Errorf("HeapList.Pop() returned %v; expected %v", popped, tt.expectedPop)
			}
			if !reflect.DeepEqual(heapCopy, tt.expectedHeap) {
				t.Errorf("HeapList after Pop() = %v; expected %v", heapCopy, tt.expectedHeap)
			}
		})
	}
}

// TestTopMaxPoint tests the TopMaxPoint function
// which returns the top N nodes with highest points in descending order
func TestTopMaxPoint(t *testing.T) {
	tests := []struct {
		name     string
		data     []HeapNode[string]
		nPoints  int
		expected []HeapNode[string]
	}{
		{
			name: "get top 3 from 5 elements",
			data: []HeapNode[string]{
				{Data: "a", Point: 5},
				{Data: "b", Point: 10},
				{Data: "c", Point: 3},
				{Data: "d", Point: 15},
				{Data: "e", Point: 7},
			},
			nPoints: 3,
			expected: []HeapNode[string]{
				{Data: "d", Point: 15},
				{Data: "b", Point: 10},
				{Data: "e", Point: 7},
			},
		},
		{
			name: "get top 1 from multiple elements",
			data: []HeapNode[string]{
				{Data: "a", Point: 100},
				{Data: "b", Point: 50},
				{Data: "c", Point: 25},
			},
			nPoints: 1,
			expected: []HeapNode[string]{
				{Data: "a", Point: 100},
			},
		},
		{
			name: "get all elements when nPoints equals data length",
			data: []HeapNode[string]{
				{Data: "a", Point: 10},
				{Data: "b", Point: 20},
				{Data: "c", Point: 5},
			},
			nPoints: 3,
			expected: []HeapNode[string]{
				{Data: "b", Point: 20},
				{Data: "a", Point: 10},
				{Data: "c", Point: 5},
			},
		},
		{
			name: "single element get top 1",
			data: []HeapNode[string]{
				{Data: "only", Point: 42},
			},
			nPoints: 1,
			expected: []HeapNode[string]{
				{Data: "only", Point: 42},
			},
		},
		{
			name: "elements with different points",
			data: []HeapNode[string]{
				{Data: "a", Point: 10},
				{Data: "b", Point: 15},
				{Data: "c", Point: 5},
			},
			nPoints: 2,
			expected: []HeapNode[string]{
				{Data: "b", Point: 15},
				{Data: "a", Point: 10},
			},
		},
		{
			name: "elements with negative points",
			data: []HeapNode[string]{
				{Data: "a", Point: -5},
				{Data: "b", Point: -10},
				{Data: "c", Point: 0},
				{Data: "d", Point: 5},
			},
			nPoints: 2,
			expected: []HeapNode[string]{
				{Data: "d", Point: 5},
				{Data: "c", Point: 0},
			},
		},
		{
			name: "get top 5 from larger dataset",
			data: []HeapNode[string]{
				{Data: "a", Point: 100},
				{Data: "b", Point: 20},
				{Data: "c", Point: 50},
				{Data: "d", Point: 30},
				{Data: "e", Point: 10},
				{Data: "f", Point: 70},
				{Data: "g", Point: 40},
				{Data: "h", Point: 60},
			},
			nPoints: 5,
			expected: []HeapNode[string]{
				{Data: "a", Point: 100},
				{Data: "f", Point: 70},
				{Data: "h", Point: 60},
				{Data: "c", Point: 50},
				{Data: "g", Point: 40},
			},
		},
		{
			name: "nPoints is zero",
			data: []HeapNode[string]{
				{Data: "a", Point: 5},
				{Data: "b", Point: 10},
			},
			nPoints: 0,
			expected: []HeapNode[string]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := TopMaxPoint(tt.data, tt.nPoints)

			// Assert
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TopMaxPoint() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestTopMaxPoint_EmptyInput tests TopMaxPoint with empty input data
// to ensure it handles edge cases gracefully
// Note: TopMaxPoint will panic if nPoints > len(data) due to slice bounds check
func TestTopMaxPoint_EmptyInput(t *testing.T) {
	tests := []struct {
		name     string
		data     []HeapNode[int]
		nPoints  int
		expected []HeapNode[int]
	}{
		{
			name:     "empty data with nPoints 0",
			data:     []HeapNode[int]{},
			nPoints:  0,
			expected: []HeapNode[int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := TopMaxPoint(tt.data, tt.nPoints)

			// Assert
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TopMaxPoint() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestHeapList_Integration tests the integration of all heap operations
// using the standard container/heap package to ensure proper heap behavior
func TestHeapList_Integration(t *testing.T) {
	tests := []struct {
		name           string
		operations     []string
		data           []HeapNode[int]
		expectedOrder  []HeapNode[int]
		expectMinFirst bool
	}{
		{
			name:       "push and pop maintains min-heap property",
			operations: []string{"push_all", "pop_all"},
			data: []HeapNode[int]{
				{Data: 1, Point: 50},
				{Data: 2, Point: 20},
				{Data: 3, Point: 80},
				{Data: 4, Point: 10},
				{Data: 5, Point: 30},
			},
			expectedOrder: []HeapNode[int]{
				{Data: 4, Point: 10},
				{Data: 2, Point: 20},
				{Data: 5, Point: 30},
				{Data: 1, Point: 50},
				{Data: 3, Point: 80},
			},
			expectMinFirst: true,
		},
		{
			name:       "sequential push and pop",
			operations: []string{"push_all", "pop_all"},
			data: []HeapNode[int]{
				{Data: 1, Point: 5},
				{Data: 2, Point: 3},
				{Data: 3, Point: 7},
			},
			expectedOrder: []HeapNode[int]{
				{Data: 2, Point: 3},
				{Data: 1, Point: 5},
				{Data: 3, Point: 7},
			},
			expectMinFirst: true,
		},
		{
			name:       "single element heap",
			operations: []string{"push_all", "pop_all"},
			data: []HeapNode[int]{
				{Data: 42, Point: 42},
			},
			expectedOrder: []HeapNode[int]{
				{Data: 42, Point: 42},
			},
			expectMinFirst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			h := make(HeapList[int], 0)
			heap.Init(&h)

			// Act: Push all elements
			for _, node := range tt.data {
				heap.Push(&h, node)
			}

			// Pop all elements and verify order
			result := make([]HeapNode[int], 0, len(tt.data))
			for h.Len() > 0 {
				result = append(result, heap.Pop(&h).(HeapNode[int]))
			}

			// Assert
			if !reflect.DeepEqual(result, tt.expectedOrder) {
				t.Errorf("Heap pop order = %v; expected %v", result, tt.expectedOrder)
			}

			// Verify min-heap property (smallest element comes first)
			if tt.expectMinFirst && len(result) > 0 {
				for i := 1; i < len(result); i++ {
					if result[i].Point < result[i-1].Point {
						t.Errorf("Heap order violated: element at index %d (Point=%d) is less than element at index %d (Point=%d)",
							i, result[i].Point, i-1, result[i-1].Point)
					}
				}
			}
		})
	}
}

// TestTopMaxPoint_WithDifferentTypes tests TopMaxPoint with different generic types
// to ensure the generic implementation works correctly
func TestTopMaxPoint_WithDifferentTypes(t *testing.T) {
	t.Run("with struct data type", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		data := []HeapNode[Person]{
			{Data: Person{Name: "Alice", Age: 30}, Point: 100},
			{Data: Person{Name: "Bob", Age: 25}, Point: 50},
			{Data: Person{Name: "Charlie", Age: 35}, Point: 75},
		}

		// Act
		result := TopMaxPoint(data, 2)

		// Assert
		expected := []HeapNode[Person]{
			{Data: Person{Name: "Alice", Age: 30}, Point: 100},
			{Data: Person{Name: "Charlie", Age: 35}, Point: 75},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("TopMaxPoint() with struct = %v; expected %v", result, expected)
		}
	})

	t.Run("with float64 data type", func(t *testing.T) {
		// Arrange
		data := []HeapNode[float64]{
			{Data: 3.14, Point: 10},
			{Data: 2.71, Point: 5},
			{Data: 1.41, Point: 15},
		}

		// Act
		result := TopMaxPoint(data, 2)

		// Assert
		expected := []HeapNode[float64]{
			{Data: 1.41, Point: 15},
			{Data: 3.14, Point: 10},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("TopMaxPoint() with float64 = %v; expected %v", result, expected)
		}
	})

	t.Run("with bool data type", func(t *testing.T) {
		// Arrange
		data := []HeapNode[bool]{
			{Data: true, Point: 20},
			{Data: false, Point: 10},
			{Data: true, Point: 30},
		}

		// Act
		result := TopMaxPoint(data, 2)

		// Assert
		expected := []HeapNode[bool]{
			{Data: true, Point: 30},
			{Data: true, Point: 20},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("TopMaxPoint() with bool = %v; expected %v", result, expected)
		}
	})
}

