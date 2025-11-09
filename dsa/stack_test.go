package dsa

import (
	"reflect"
	"testing"
)

// TestStack_Len tests the Len method of Stack
// which returns the number of elements in the stack
func TestStack_Len(t *testing.T) {
	tests := []struct {
		name     string
		stack    Stack[int]
		expected int
	}{
		{
			name:     "empty stack",
			stack:    Stack[int]{},
			expected: 0,
		},
		{
			name:     "stack with one element",
			stack:    Stack[int]{10},
			expected: 1,
		},
		{
			name:     "stack with multiple elements",
			stack:    Stack[int]{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "stack with many elements",
			stack:    Stack[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.stack.Len()

			// Assert
			if result != tt.expected {
				t.Errorf("Stack.Len() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestStack_IsEmpty tests the IsEmpty method of Stack
// which returns true if the stack has no elements
func TestStack_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		stack    Stack[string]
		expected bool
	}{
		{
			name:     "empty stack returns true",
			stack:    Stack[string]{},
			expected: true,
		},
		{
			name:     "stack with one element returns false",
			stack:    Stack[string]{"item"},
			expected: false,
		},
		{
			name:     "stack with multiple elements returns false",
			stack:    Stack[string]{"a", "b", "c"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.stack.IsEmpty()

			// Assert
			if result != tt.expected {
				t.Errorf("Stack.IsEmpty() = %v; expected %v (len=%d)", result, tt.expected, tt.stack.Len())
			}
		})
	}
}

// TestStack_Push tests the Push method of Stack
// which adds elements to the top of the stack
func TestStack_Push(t *testing.T) {
	tests := []struct {
		name     string
		initial  Stack[int]
		toPush   []int
		expected Stack[int]
	}{
		{
			name:     "push to empty stack",
			initial:  Stack[int]{},
			toPush:   []int{1},
			expected: Stack[int]{1},
		},
		{
			name:     "push to stack with elements",
			initial:  Stack[int]{1, 2},
			toPush:   []int{3},
			expected: Stack[int]{1, 2, 3},
		},
		{
			name:     "push multiple elements",
			initial:  Stack[int]{},
			toPush:   []int{1, 2, 3, 4, 5},
			expected: Stack[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "push zero value",
			initial:  Stack[int]{1},
			toPush:   []int{0},
			expected: Stack[int]{1, 0},
		},
		{
			name:     "push negative numbers",
			initial:  Stack[int]{},
			toPush:   []int{-1, -2, -3},
			expected: Stack[int]{-1, -2, -3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			stack := make(Stack[int], len(tt.initial))
			copy(stack, tt.initial)

			// Act
			for _, item := range tt.toPush {
				stack.Push(item)
			}

			// Assert
			if !reflect.DeepEqual(stack, tt.expected) {
				t.Errorf("Stack.Push() resulted in %v; expected %v", stack, tt.expected)
			}
		})
	}
}

// TestStack_Pop tests the Pop method of Stack
// which removes and returns the top element from the stack
func TestStack_Pop(t *testing.T) {
	tests := []struct {
		name          string
		initial       Stack[int]
		expectedItem  int
		expectedOk    bool
		expectedStack Stack[int]
	}{
		{
			name:          "pop from empty stack",
			initial:       Stack[int]{},
			expectedItem:  0,
			expectedOk:    false,
			expectedStack: Stack[int]{},
		},
		{
			name:          "pop from stack with one element",
			initial:       Stack[int]{42},
			expectedItem:  42,
			expectedOk:    true,
			expectedStack: Stack[int]{},
		},
		{
			name:          "pop from stack with multiple elements",
			initial:       Stack[int]{1, 2, 3, 4, 5},
			expectedItem:  5,
			expectedOk:    true,
			expectedStack: Stack[int]{1, 2, 3, 4},
		},
		{
			name:          "pop zero value from stack",
			initial:       Stack[int]{1, 2, 0},
			expectedItem:  0,
			expectedOk:    true,
			expectedStack: Stack[int]{1, 2},
		},
		{
			name:          "pop negative number",
			initial:       Stack[int]{1, -5},
			expectedItem:  -5,
			expectedOk:    true,
			expectedStack: Stack[int]{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			stack := make(Stack[int], len(tt.initial))
			copy(stack, tt.initial)

			// Act
			item, ok := stack.Pop()

			// Assert
			if item != tt.expectedItem {
				t.Errorf("Stack.Pop() item = %v; expected %v", item, tt.expectedItem)
			}
			if ok != tt.expectedOk {
				t.Errorf("Stack.Pop() ok = %v; expected %v", ok, tt.expectedOk)
			}
			if !reflect.DeepEqual(stack, tt.expectedStack) {
				t.Errorf("Stack after Pop() = %v; expected %v", stack, tt.expectedStack)
			}
		})
	}
}

// TestStack_Peek tests the Peek method of Stack
// which returns the top element without removing it
func TestStack_Peek(t *testing.T) {
	tests := []struct {
		name          string
		initial       Stack[string]
		expectedItem  string
		expectedOk    bool
		expectedStack Stack[string]
	}{
		{
			name:          "peek at empty stack",
			initial:       Stack[string]{},
			expectedItem:  "",
			expectedOk:    false,
			expectedStack: Stack[string]{},
		},
		{
			name:          "peek at stack with one element",
			initial:       Stack[string]{"first"},
			expectedItem:  "first",
			expectedOk:    true,
			expectedStack: Stack[string]{"first"},
		},
		{
			name:          "peek at stack with multiple elements",
			initial:       Stack[string]{"bottom", "middle", "top"},
			expectedItem:  "top",
			expectedOk:    true,
			expectedStack: Stack[string]{"bottom", "middle", "top"},
		},
		{
			name:          "peek does not modify stack",
			initial:       Stack[string]{"a", "b", "c"},
			expectedItem:  "c",
			expectedOk:    true,
			expectedStack: Stack[string]{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			stack := make(Stack[string], len(tt.initial))
			copy(stack, tt.initial)

			// Act
			item, ok := stack.Peek()

			// Assert
			if item != tt.expectedItem {
				t.Errorf("Stack.Peek() item = %v; expected %v", item, tt.expectedItem)
			}
			if ok != tt.expectedOk {
				t.Errorf("Stack.Peek() ok = %v; expected %v", ok, tt.expectedOk)
			}
			if !reflect.DeepEqual(stack, tt.expectedStack) {
				t.Errorf("Stack after Peek() = %v; expected %v (Peek should not modify stack)", stack, tt.expectedStack)
			}
		})
	}
}

// TestStack_Clear tests the Clear method of Stack
// which removes all elements from the stack
func TestStack_Clear(t *testing.T) {
	tests := []struct {
		name     string
		initial  Stack[int]
		expected Stack[int]
	}{
		{
			name:     "clear empty stack",
			initial:  Stack[int]{},
			expected: Stack[int]{},
		},
		{
			name:     "clear stack with one element",
			initial:  Stack[int]{42},
			expected: Stack[int]{},
		},
		{
			name:     "clear stack with multiple elements",
			initial:  Stack[int]{1, 2, 3, 4, 5},
			expected: Stack[int]{},
		},
		{
			name:     "clear large stack",
			initial:  Stack[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: Stack[int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create a copy to avoid modifying test data
			stack := make(Stack[int], len(tt.initial))
			copy(stack, tt.initial)

			// Act
			stack.Clear()

			// Assert
			if !reflect.DeepEqual(stack, tt.expected) {
				t.Errorf("Stack.Clear() resulted in %v; expected %v", stack, tt.expected)
			}
			if stack.Len() != 0 {
				t.Errorf("Stack.Len() after Clear() = %d; expected 0", stack.Len())
			}
			if !stack.IsEmpty() {
				t.Errorf("Stack.IsEmpty() after Clear() = false; expected true")
			}
		})
	}
}

// TestStack_Integration tests multiple stack operations together
// to ensure they work correctly in combination
func TestStack_Integration(t *testing.T) {
	tests := []struct {
		name       string
		operations []string
		values     []int
		expected   Stack[int]
	}{
		{
			name:       "push and pop sequence",
			operations: []string{"push", "push", "pop", "push", "pop", "pop"},
			values:     []int{1, 2, 0, 3, 0, 0},
			expected:   Stack[int]{},
		},
		{
			name:       "push multiple then pop once",
			operations: []string{"push", "push", "push", "pop"},
			values:     []int{10, 20, 30, 0},
			expected:   Stack[int]{10, 20},
		},
		{
			name:       "clear after pushes",
			operations: []string{"push", "push", "push", "clear"},
			values:     []int{1, 2, 3, 0},
			expected:   Stack[int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			stack := Stack[int]{}

			// Act
			for i, op := range tt.operations {
				switch op {
				case "push":
					stack.Push(tt.values[i])
				case "pop":
					stack.Pop()
				case "clear":
					stack.Clear()
				}
			}

			// Assert
			if !reflect.DeepEqual(stack, tt.expected) {
				t.Errorf("Stack after operations = %v; expected %v", stack, tt.expected)
			}
		})
	}
}

// TestStack_LIFOOrder tests that Stack follows Last-In-First-Out order
// which is the fundamental property of a stack data structure
func TestStack_LIFOOrder(t *testing.T) {
	tests := []struct {
		name          string
		pushSequence  []int
		expectedOrder []int
	}{
		{
			name:          "simple sequence",
			pushSequence:  []int{1, 2, 3},
			expectedOrder: []int{3, 2, 1},
		},
		{
			name:          "single element",
			pushSequence:  []int{42},
			expectedOrder: []int{42},
		},
		{
			name:          "larger sequence",
			pushSequence:  []int{10, 20, 30, 40, 50},
			expectedOrder: []int{50, 40, 30, 20, 10},
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
			stack := Stack[int]{}

			// Act: Push all elements
			for _, val := range tt.pushSequence {
				stack.Push(val)
			}

			// Pop all elements and verify LIFO order
			result := make([]int, 0, len(tt.expectedOrder))
			for !stack.IsEmpty() {
				item, ok := stack.Pop()
				if !ok {
					t.Errorf("Pop failed when stack should not be empty")
					break
				}
				result = append(result, item)
			}

			// Assert
			if !reflect.DeepEqual(result, tt.expectedOrder) {
				t.Errorf("LIFO order = %v; expected %v", result, tt.expectedOrder)
			}
		})
	}
}

// TestStack_MultiplePopOnEmpty tests behavior when popping from empty stack multiple times
// to ensure it consistently returns false without panicking
func TestStack_MultiplePopOnEmpty(t *testing.T) {
	tests := []struct {
		name      string
		popCount  int
		shouldOk  bool
	}{
		{
			name:     "pop once from empty stack",
			popCount: 1,
			shouldOk: false,
		},
		{
			name:     "pop multiple times from empty stack",
			popCount: 5,
			shouldOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			stack := Stack[int]{}

			// Act & Assert
			for i := 0; i < tt.popCount; i++ {
				item, ok := stack.Pop()
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

// TestStack_WithDifferentTypes tests Stack with different generic types
// to ensure the generic implementation works correctly
func TestStack_WithDifferentTypes(t *testing.T) {
	t.Run("with string type", func(t *testing.T) {
		// Arrange
		stack := Stack[string]{}

		// Act
		stack.Push("hello")
		stack.Push("world")
		item, ok := stack.Pop()

		// Assert
		if !ok || item != "world" {
			t.Errorf("Pop() = (%v, %v); expected (world, true)", item, ok)
		}
	})

	t.Run("with struct type", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		stack := Stack[Person]{}

		// Act
		stack.Push(Person{Name: "Alice", Age: 30})
		stack.Push(Person{Name: "Bob", Age: 25})
		item, ok := stack.Peek()

		// Assert
		expected := Person{Name: "Bob", Age: 25}
		if !ok || item != expected {
			t.Errorf("Peek() = (%v, %v); expected (%v, true)", item, ok, expected)
		}
	})

	t.Run("with pointer type", func(t *testing.T) {
		// Arrange
		stack := Stack[*int]{}
		val1 := 100
		val2 := 200

		// Act
		stack.Push(&val1)
		stack.Push(&val2)
		item, ok := stack.Pop()

		// Assert
		if !ok || item == nil || *item != 200 {
			t.Errorf("Pop() = (%v, %v); expected pointer to 200", item, ok)
		}
	})

	t.Run("with bool type", func(t *testing.T) {
		// Arrange
		stack := Stack[bool]{}

		// Act
		stack.Push(true)
		stack.Push(false)
		stack.Push(true)
		item1, ok1 := stack.Pop()
		item2, ok2 := stack.Pop()

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
		stack := Stack[float64]{}

		// Act
		stack.Push(3.14)
		stack.Push(2.71)
		item, ok := stack.Peek()

		// Assert
		if !ok || item != 2.71 {
			t.Errorf("Peek() = (%v, %v); expected (2.71, true)", item, ok)
		}
		if stack.Len() != 2 {
			t.Errorf("Len() = %d; expected 2 (Peek should not remove element)", stack.Len())
		}
	})
}

// TestStack_Concurrency tests that Stack behaves correctly with concurrent operations
// Note: Stack is not thread-safe by design, but we test sequential operations
func TestStack_SequentialOperations(t *testing.T) {
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
			stack := Stack[int]{}

			// Act: Push N items
			for i := 0; i < tt.iterations; i++ {
				stack.Push(i)
			}

			// Assert: Stack has correct length
			if stack.Len() != tt.iterations {
				t.Errorf("After pushing %d items, Len() = %d; expected %d", tt.iterations, stack.Len(), tt.iterations)
			}

			// Act: Pop all items
			for i := tt.iterations - 1; i >= 0; i-- {
				item, ok := stack.Pop()
				if !ok {
					t.Errorf("Pop failed at iteration %d", i)
					break
				}
				if item != i {
					t.Errorf("Pop() = %d; expected %d", item, i)
				}
			}

			// Assert: Stack is empty
			if !stack.IsEmpty() {
				t.Errorf("After popping all items, IsEmpty() = false; expected true")
			}
		})
	}
}

// TestStack_ClearAndReuse tests that a stack can be reused after clearing
// to ensure the Clear operation properly resets the stack state
func TestStack_ClearAndReuse(t *testing.T) {
	tests := []struct {
		name            string
		firstSequence   []int
		secondSequence  []int
		expectedFinal   Stack[int]
	}{
		{
			name:           "clear and push new values",
			firstSequence:  []int{1, 2, 3},
			secondSequence: []int{10, 20},
			expectedFinal:  Stack[int]{10, 20},
		},
		{
			name:           "clear empty stack and push",
			firstSequence:  []int{},
			secondSequence: []int{5, 6, 7},
			expectedFinal:  Stack[int]{5, 6, 7},
		},
		{
			name:           "clear and push same values",
			firstSequence:  []int{1, 2},
			secondSequence: []int{1, 2},
			expectedFinal:  Stack[int]{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			stack := Stack[int]{}

			// Act: First sequence
			for _, val := range tt.firstSequence {
				stack.Push(val)
			}

			// Clear the stack
			stack.Clear()

			// Second sequence
			for _, val := range tt.secondSequence {
				stack.Push(val)
			}

			// Assert
			if !reflect.DeepEqual(stack, tt.expectedFinal) {
				t.Errorf("Stack after clear and reuse = %v; expected %v", stack, tt.expectedFinal)
			}
		})
	}
}

