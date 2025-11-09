package dsa

import (
	"reflect"
	"testing"
)

// TestQueue_Len tests the Len method of Queue
// which returns the number of elements in the queue
func TestQueue_Len(t *testing.T) {
	tests := []struct {
		name     string
		queue    Queue[int]
		expected int
	}{
		{
			name:     "empty queue",
			queue:    Queue[int]{},
			expected: 0,
		},
		{
			name:     "queue with one element",
			queue:    Queue[int]{10},
			expected: 1,
		},
		{
			name:     "queue with multiple elements",
			queue:    Queue[int]{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "queue with many elements",
			queue:    Queue[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.queue.Len()

			// Assert
			if result != tt.expected {
				t.Errorf("Queue.Len() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestQueue_IsEmpty tests the IsEmpty method of Queue
// which returns true if the queue has no elements
func TestQueue_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		queue    Queue[string]
		expected bool
	}{
		{
			name:     "empty queue returns true",
			queue:    Queue[string]{},
			expected: true,
		},
		{
			name:     "queue with one element returns false",
			queue:    Queue[string]{"item"},
			expected: false,
		},
		{
			name:     "queue with multiple elements returns false",
			queue:    Queue[string]{"a", "b", "c"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.queue.IsEmpty()

			// Assert
			if result != tt.expected {
				t.Errorf("Queue.IsEmpty() = %v; expected %v (len=%d)", result, tt.expected, tt.queue.Len())
			}
		})
	}
}

// TestQueue_Push tests the Push method of Queue
// which adds elements to the back of the queue
func TestQueue_Push(t *testing.T) {
	tests := []struct {
		name     string
		initial  Queue[int]
		toPush   []int
		expected Queue[int]
	}{
		{
			name:     "push to empty queue",
			initial:  Queue[int]{},
			toPush:   []int{1},
			expected: Queue[int]{1},
		},
		{
			name:     "push to queue with elements",
			initial:  Queue[int]{1, 2},
			toPush:   []int{3},
			expected: Queue[int]{1, 2, 3},
		},
		{
			name:     "push multiple elements",
			initial:  Queue[int]{},
			toPush:   []int{1, 2, 3, 4, 5},
			expected: Queue[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "push zero value",
			initial:  Queue[int]{1},
			toPush:   []int{0},
			expected: Queue[int]{1, 0},
		},
		{
			name:     "push negative numbers",
			initial:  Queue[int]{},
			toPush:   []int{-1, -2, -3},
			expected: Queue[int]{-1, -2, -3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			queue := make(Queue[int], len(tt.initial))
			copy(queue, tt.initial)

			// Act
			for _, item := range tt.toPush {
				queue.Push(item)
			}

			// Assert
			if !reflect.DeepEqual(queue, tt.expected) {
				t.Errorf("Queue.Push() resulted in %v; expected %v", queue, tt.expected)
			}
		})
	}
}

// TestQueue_Pop tests the Pop method of Queue
// which removes and returns the front element from the queue
func TestQueue_Pop(t *testing.T) {
	tests := []struct {
		name          string
		initial       Queue[int]
		expectedItem  int
		expectedOk    bool
		expectedQueue Queue[int]
	}{
		{
			name:          "pop from empty queue",
			initial:       Queue[int]{},
			expectedItem:  0,
			expectedOk:    false,
			expectedQueue: Queue[int]{},
		},
		{
			name:          "pop from queue with one element",
			initial:       Queue[int]{42},
			expectedItem:  42,
			expectedOk:    true,
			expectedQueue: Queue[int]{},
		},
		{
			name:          "pop from queue with multiple elements",
			initial:       Queue[int]{1, 2, 3, 4, 5},
			expectedItem:  1,
			expectedOk:    true,
			expectedQueue: Queue[int]{2, 3, 4, 5},
		},
		{
			name:          "pop zero value from queue",
			initial:       Queue[int]{0, 1, 2},
			expectedItem:  0,
			expectedOk:    true,
			expectedQueue: Queue[int]{1, 2},
		},
		{
			name:          "pop negative number",
			initial:       Queue[int]{-5, 1},
			expectedItem:  -5,
			expectedOk:    true,
			expectedQueue: Queue[int]{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			queue := make(Queue[int], len(tt.initial))
			copy(queue, tt.initial)

			// Act
			item, ok := queue.Pop()

			// Assert
			if item != tt.expectedItem {
				t.Errorf("Queue.Pop() item = %v; expected %v", item, tt.expectedItem)
			}
			if ok != tt.expectedOk {
				t.Errorf("Queue.Pop() ok = %v; expected %v", ok, tt.expectedOk)
			}
			if !reflect.DeepEqual(queue, tt.expectedQueue) {
				t.Errorf("Queue after Pop() = %v; expected %v", queue, tt.expectedQueue)
			}
		})
	}
}

// TestQueue_Peek tests the Peek method of Queue
// which returns the front element without removing it
func TestQueue_Peek(t *testing.T) {
	tests := []struct {
		name          string
		initial       Queue[string]
		expectedItem  string
		expectedOk    bool
		expectedQueue Queue[string]
	}{
		{
			name:          "peek at empty queue",
			initial:       Queue[string]{},
			expectedItem:  "",
			expectedOk:    false,
			expectedQueue: Queue[string]{},
		},
		{
			name:          "peek at queue with one element",
			initial:       Queue[string]{"first"},
			expectedItem:  "first",
			expectedOk:    true,
			expectedQueue: Queue[string]{"first"},
		},
		{
			name:          "peek at queue with multiple elements",
			initial:       Queue[string]{"front", "middle", "back"},
			expectedItem:  "front",
			expectedOk:    true,
			expectedQueue: Queue[string]{"front", "middle", "back"},
		},
		{
			name:          "peek does not modify queue",
			initial:       Queue[string]{"a", "b", "c"},
			expectedItem:  "a",
			expectedOk:    true,
			expectedQueue: Queue[string]{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			queue := make(Queue[string], len(tt.initial))
			copy(queue, tt.initial)

			// Act
			item, ok := queue.Peek()

			// Assert
			if item != tt.expectedItem {
				t.Errorf("Queue.Peek() item = %v; expected %v", item, tt.expectedItem)
			}
			if ok != tt.expectedOk {
				t.Errorf("Queue.Peek() ok = %v; expected %v", ok, tt.expectedOk)
			}
			if !reflect.DeepEqual(queue, tt.expectedQueue) {
				t.Errorf("Queue after Peek() = %v; expected %v (Peek should not modify queue)", queue, tt.expectedQueue)
			}
		})
	}
}

// TestQueue_Clear tests the Clear method of Queue
// which removes all elements from the queue
func TestQueue_Clear(t *testing.T) {
	tests := []struct {
		name     string
		initial  Queue[int]
		expected Queue[int]
	}{
		{
			name:     "clear empty queue",
			initial:  Queue[int]{},
			expected: Queue[int]{},
		},
		{
			name:     "clear queue with one element",
			initial:  Queue[int]{42},
			expected: Queue[int]{},
		},
		{
			name:     "clear queue with multiple elements",
			initial:  Queue[int]{1, 2, 3, 4, 5},
			expected: Queue[int]{},
		},
		{
			name:     "clear large queue",
			initial:  Queue[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: Queue[int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			queue := make(Queue[int], len(tt.initial))
			copy(queue, tt.initial)

			// Act
			queue.Clear()

			// Assert
			if !reflect.DeepEqual(queue, tt.expected) {
				t.Errorf("Queue.Clear() resulted in %v; expected %v", queue, tt.expected)
			}
			if queue.Len() != 0 {
				t.Errorf("Queue.Len() after Clear() = %d; expected 0", queue.Len())
			}
			if !queue.IsEmpty() {
				t.Errorf("Queue.IsEmpty() after Clear() = false; expected true")
			}
		})
	}
}

// TestQueue_Integration tests multiple queue operations together
// to ensure they work correctly in combination
func TestQueue_Integration(t *testing.T) {
	tests := []struct {
		name       string
		operations []string
		values     []int
		expected   Queue[int]
	}{
		{
			name:       "push and pop sequence",
			operations: []string{"push", "push", "pop", "push", "pop", "pop"},
			values:     []int{1, 2, 0, 3, 0, 0},
			expected:   Queue[int]{},
		},
		{
			name:       "push multiple then pop once",
			operations: []string{"push", "push", "push", "pop"},
			values:     []int{10, 20, 30, 0},
			expected:   Queue[int]{20, 30},
		},
		{
			name:       "clear after pushes",
			operations: []string{"push", "push", "push", "clear"},
			values:     []int{1, 2, 3, 0},
			expected:   Queue[int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			queue := Queue[int]{}

			// Act
			for i, op := range tt.operations {
				switch op {
				case "push":
					queue.Push(tt.values[i])
				case "pop":
					queue.Pop()
				case "clear":
					queue.Clear()
				}
			}

			// Assert
			if !reflect.DeepEqual(queue, tt.expected) {
				t.Errorf("Queue after operations = %v; expected %v", queue, tt.expected)
			}
		})
	}
}

// TestQueue_FIFOOrder tests that Queue follows First-In-First-Out order
// which is the fundamental property of a queue data structure
func TestQueue_FIFOOrder(t *testing.T) {
	tests := []struct {
		name          string
		pushSequence  []int
		expectedOrder []int
	}{
		{
			name:          "simple sequence",
			pushSequence:  []int{1, 2, 3},
			expectedOrder: []int{1, 2, 3},
		},
		{
			name:          "single element",
			pushSequence:  []int{42},
			expectedOrder: []int{42},
		},
		{
			name:          "larger sequence",
			pushSequence:  []int{10, 20, 30, 40, 50},
			expectedOrder: []int{10, 20, 30, 40, 50},
		},
		{
			name:          "sequence with duplicates",
			pushSequence:  []int{1, 2, 1, 2, 1},
			expectedOrder: []int{1, 2, 1, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			queue := Queue[int]{}

			// Act: Push all elements
			for _, val := range tt.pushSequence {
				queue.Push(val)
			}

			// Pop all elements and verify FIFO order
			result := make([]int, 0, len(tt.expectedOrder))
			for !queue.IsEmpty() {
				item, ok := queue.Pop()
				if !ok {
					t.Errorf("Pop failed when queue should not be empty")
					break
				}
				result = append(result, item)
			}

			// Assert
			if !reflect.DeepEqual(result, tt.expectedOrder) {
				t.Errorf("FIFO order = %v; expected %v", result, tt.expectedOrder)
			}
		})
	}
}

// TestQueue_MultiplePopOnEmpty tests behavior when popping from empty queue multiple times
// to ensure it consistently returns false without panicking
func TestQueue_MultiplePopOnEmpty(t *testing.T) {
	tests := []struct {
		name     string
		popCount int
		shouldOk bool
	}{
		{
			name:     "pop once from empty queue",
			popCount: 1,
			shouldOk: false,
		},
		{
			name:     "pop multiple times from empty queue",
			popCount: 5,
			shouldOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			queue := Queue[int]{}

			// Act & Assert
			for i := 0; i < tt.popCount; i++ {
				item, ok := queue.Pop()
				if ok != tt.shouldOk {
					t.Errorf("Pop iteration %d: ok = %v; expected %v", i+1, ok, tt.shouldOk)
				}
				if ok && item != 0 {
					t.Errorf("Pop iteration %d: item = %v; expected zero value", i+1, item)
				}
			}
		})
	}
}

// TestQueue_WithDifferentTypes tests Queue with different generic types
// to ensure the generic implementation works correctly
func TestQueue_WithDifferentTypes(t *testing.T) {
	t.Run("with string type", func(t *testing.T) {
		// Arrange
		queue := Queue[string]{}

		// Act
		queue.Push("hello")
		queue.Push("world")
		item, ok := queue.Pop()

		// Assert
		if !ok || item != "hello" {
			t.Errorf("Pop() = (%v, %v); expected (hello, true)", item, ok)
		}
	})

	t.Run("with struct type", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		queue := Queue[Person]{}

		// Act
		queue.Push(Person{Name: "Alice", Age: 30})
		queue.Push(Person{Name: "Bob", Age: 25})
		item, ok := queue.Peek()

		// Assert
		expected := Person{Name: "Alice", Age: 30}
		if !ok || item != expected {
			t.Errorf("Peek() = (%v, %v); expected (%v, true)", item, ok, expected)
		}
	})

	t.Run("with pointer type", func(t *testing.T) {
		// Arrange
		queue := Queue[*int]{}
		val1 := 100
		val2 := 200

		// Act
		queue.Push(&val1)
		queue.Push(&val2)
		item, ok := queue.Pop()

		// Assert
		if !ok || item == nil || *item != 100 {
			t.Errorf("Pop() = (%v, %v); expected pointer to 100", item, ok)
		}
	})

	t.Run("with bool type", func(t *testing.T) {
		// Arrange
		queue := Queue[bool]{}

		// Act
		queue.Push(true)
		queue.Push(false)
		queue.Push(true)
		item1, ok1 := queue.Pop()
		item2, ok2 := queue.Pop()

		// Assert
		if !ok1 || item1 != true {
			t.Errorf("First Pop() = (%v, %v); expected (true, true)", item1, ok1)
		}
		if !ok2 || item2 != false {
			t.Errorf("Second Pop() = (%v, %v); expected (false, true)", item2, ok2)
		}
	})

	t.Run("with float64 type", func(t *testing.T) {
		// Arrange
		queue := Queue[float64]{}

		// Act
		queue.Push(3.14)
		queue.Push(2.71)
		item, ok := queue.Peek()

		// Assert
		if !ok || item != 3.14 {
			t.Errorf("Peek() = (%v, %v); expected (3.14, true)", item, ok)
		}
		if queue.Len() != 2 {
			t.Errorf("Len() = %d; expected 2 (Peek should not remove element)", queue.Len())
		}
	})
}

// TestQueue_SequentialOperations tests that Queue behaves correctly with many operations
// to verify performance and correctness at scale
func TestQueue_SequentialOperations(t *testing.T) {
	tests := []struct {
		name       string
		iterations int
	}{
		{
			name:       "small number of iterations",
			iterations: 10,
		},
		{
			name:       "medium number of iterations",
			iterations: 100,
		},
		{
			name:       "large number of iterations",
			iterations: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			queue := Queue[int]{}

			// Act: Push N items
			for i := 0; i < tt.iterations; i++ {
				queue.Push(i)
			}

			// Assert: Queue has correct length
			if queue.Len() != tt.iterations {
				t.Errorf("After pushing %d items, Len() = %d; expected %d", tt.iterations, queue.Len(), tt.iterations)
			}

			// Act: Pop all items
			for i := 0; i < tt.iterations; i++ {
				item, ok := queue.Pop()
				if !ok {
					t.Errorf("Pop failed at iteration %d", i)
					break
				}
				if item != i {
					t.Errorf("Pop() = %d; expected %d (FIFO order)", item, i)
				}
			}

			// Assert: Queue is empty
			if !queue.IsEmpty() {
				t.Errorf("After popping all items, IsEmpty() = false; expected true")
			}
		})
	}
}

// TestQueue_ClearAndReuse tests that a queue can be reused after clearing
// to ensure the Clear operation properly resets the queue state
func TestQueue_ClearAndReuse(t *testing.T) {
	tests := []struct {
		name           string
		firstSequence  []int
		secondSequence []int
		expectedFinal  Queue[int]
	}{
		{
			name:           "clear and push new values",
			firstSequence:  []int{1, 2, 3},
			secondSequence: []int{10, 20},
			expectedFinal:  Queue[int]{10, 20},
		},
		{
			name:           "clear empty queue and push",
			firstSequence:  []int{},
			secondSequence: []int{5, 6, 7},
			expectedFinal:  Queue[int]{5, 6, 7},
		},
		{
			name:           "clear and push same values",
			firstSequence:  []int{1, 2},
			secondSequence: []int{1, 2},
			expectedFinal:  Queue[int]{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			queue := Queue[int]{}

			// Act: First sequence
			for _, val := range tt.firstSequence {
				queue.Push(val)
			}

			// Clear the queue
			queue.Clear()

			// Second sequence
			for _, val := range tt.secondSequence {
				queue.Push(val)
			}

			// Assert
			if !reflect.DeepEqual(queue, tt.expectedFinal) {
				t.Errorf("Queue after clear and reuse = %v; expected %v", queue, tt.expectedFinal)
			}
		})
	}
}

// TestQueue_MixedOperations tests alternating push and pop operations
// to ensure the queue maintains correct state with interleaved operations
func TestQueue_MixedOperations(t *testing.T) {
	tests := []struct {
		name          string
		operations    []string
		values        []int
		expectedFinal Queue[int]
	}{
		{
			name:          "alternating push and pop",
			operations:    []string{"push", "pop", "push", "pop", "push"},
			values:        []int{1, 0, 2, 0, 3},
			expectedFinal: Queue[int]{3},
		},
		{
			name:          "multiple pushes then pops",
			operations:    []string{"push", "push", "push", "pop", "pop"},
			values:        []int{1, 2, 3, 0, 0},
			expectedFinal: Queue[int]{3},
		},
		{
			name:          "push pop push pattern",
			operations:    []string{"push", "push", "pop", "push", "push", "pop", "pop"},
			values:        []int{10, 20, 0, 30, 40, 0, 0},
			expectedFinal: Queue[int]{40},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			queue := Queue[int]{}

			// Act
			for i, op := range tt.operations {
				switch op {
				case "push":
					queue.Push(tt.values[i])
				case "pop":
					queue.Pop()
				}
			}

			// Assert
			if !reflect.DeepEqual(queue, tt.expectedFinal) {
				t.Errorf("Queue after mixed operations = %v; expected %v", queue, tt.expectedFinal)
			}
		})
	}
}

// TestQueue_PeekConsistency tests that Peek always returns the front element
// and doesn't modify the queue regardless of how many times it's called
func TestQueue_PeekConsistency(t *testing.T) {
	tests := []struct {
		name       string
		initial    Queue[int]
		peekCount  int
		expectedOk bool
	}{
		{
			name:       "peek multiple times on non-empty queue",
			initial:    Queue[int]{1, 2, 3},
			peekCount:  5,
			expectedOk: true,
		},
		{
			name:       "peek multiple times on single element queue",
			initial:    Queue[int]{42},
			peekCount:  3,
			expectedOk: true,
		},
		{
			name:       "peek multiple times on empty queue",
			initial:    Queue[int]{},
			peekCount:  3,
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			queue := make(Queue[int], len(tt.initial))
			copy(queue, tt.initial)
			originalLen := queue.Len()
			var firstItem int

			// Act & Assert
			for i := 0; i < tt.peekCount; i++ {
				item, ok := queue.Peek()
				
				if i == 0 {
					firstItem = item
				}

				// Verify consistency
				if ok != tt.expectedOk {
					t.Errorf("Peek iteration %d: ok = %v; expected %v", i+1, ok, tt.expectedOk)
				}
				if ok && item != firstItem {
					t.Errorf("Peek iteration %d: item = %v; expected %v (should be consistent)", i+1, item, firstItem)
				}
				if queue.Len() != originalLen {
					t.Errorf("Peek iteration %d: Len() = %d; expected %d (Peek should not modify queue)", i+1, queue.Len(), originalLen)
				}
			}

			// Verify queue wasn't modified
			if !reflect.DeepEqual(queue, tt.initial) {
				t.Errorf("Queue after multiple Peeks = %v; expected %v (queue should be unchanged)", queue, tt.initial)
			}
		})
	}
}

