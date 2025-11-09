package dsa

import (
	"testing"
)

// TestNew tests the NewDLList function for creating a doubly linked list
func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		expectedLen   int
		expectedFront interface{}
		expectedBack  interface{}
	}{
		{
			name:          "new int list",
			expectedLen:   0,
			expectedFront: nil,
			expectedBack:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			if list == nil {
				t.Errorf("NewDLList() returned nil")
			}
			if list.Len() != tt.expectedLen {
				t.Errorf("NewDLList() list length = %d; expected %d", list.Len(), tt.expectedLen)
			}
			if list.Front() != nil {
				t.Errorf("NewDLList() list Front() should be nil")
			}
			if list.Back() != nil {
				t.Errorf("NewDLList() list Back() should be nil")
			}
		})
	}
}

// TestPushFront tests the PushFront function
func TestPushFront(t *testing.T) {
	tests := []struct {
		name               string
		values             []int
		expectedLen        int
		expectedFrontValue int
	}{
		{
			name:               "push single element",
			values:             []int{1},
			expectedLen:        1,
			expectedFrontValue: 1,
		},
		{
			name:               "push multiple elements",
			values:             []int{1, 2, 3},
			expectedLen:        3,
			expectedFrontValue: 3,
		},
		{
			name:               "push zero value",
			values:             []int{0},
			expectedLen:        1,
			expectedFrontValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			for _, v := range tt.values {
				list.PushFront(v)
			}
			if list.Len() != tt.expectedLen {
				t.Errorf("PushFront() list length = %d; expected %d", list.Len(), tt.expectedLen)
			}
			if list.Front() == nil {
				t.Errorf("PushFront() Front() returned nil")
				return
			}
			if list.Front().Value != tt.expectedFrontValue {
				t.Errorf("PushFront() Front().Value = %v; expected %v", list.Front().Value, tt.expectedFrontValue)
			}
		})
	}
}

// TestPushBack tests the PushBack function
func TestPushBack(t *testing.T) {
	tests := []struct {
		name              string
		values            []int
		expectedLen       int
		expectedBackValue int
	}{
		{
			name:              "push single element",
			values:            []int{1},
			expectedLen:       1,
			expectedBackValue: 1,
		},
		{
			name:              "push multiple elements",
			values:            []int{1, 2, 3},
			expectedLen:       3,
			expectedBackValue: 3,
		},
		{
			name:              "push negative values",
			values:            []int{-1, -2, -3},
			expectedLen:       3,
			expectedBackValue: -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			for _, v := range tt.values {
				list.PushBack(v)
			}
			if list.Len() != tt.expectedLen {
				t.Errorf("PushBack() list length = %d; expected %d", list.Len(), tt.expectedLen)
			}
			if list.Back() == nil {
				t.Errorf("PushBack() Back() returned nil")
				return
			}
			if list.Back().Value != tt.expectedBackValue {
				t.Errorf("PushBack() Back().Value = %v; expected %v", list.Back().Value, tt.expectedBackValue)
			}
		})
	}
}

// TestFront tests the Front function
func TestFront(t *testing.T) {
	tests := []struct {
		name          string
		values        []int
		expectedValue int
		expectNil     bool
	}{
		{
			name:          "empty list",
			values:        []int{},
			expectedValue: 0,
			expectNil:     true,
		},
		{
			name:          "single element",
			values:        []int{42},
			expectedValue: 42,
			expectNil:     false,
		},
		{
			name:          "multiple elements",
			values:        []int{1, 2, 3},
			expectedValue: 1,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			for _, v := range tt.values {
				list.PushBack(v)
			}
			front := list.Front()
			if tt.expectNil {
				if front != nil {
					t.Errorf("Front() = %v; expected nil", front)
				}
			} else {
				if front == nil {
					t.Errorf("Front() = nil; expected non-nil")
					return
				}
				if front.Value != tt.expectedValue {
					t.Errorf("Front().Value = %v; expected %v", front.Value, tt.expectedValue)
				}
			}
		})
	}
}

// TestBack tests the Back function
func TestBack(t *testing.T) {
	tests := []struct {
		name          string
		values        []int
		expectedValue int
		expectNil     bool
	}{
		{
			name:          "empty list",
			values:        []int{},
			expectedValue: 0,
			expectNil:     true,
		},
		{
			name:          "single element",
			values:        []int{42},
			expectedValue: 42,
			expectNil:     false,
		},
		{
			name:          "multiple elements",
			values:        []int{1, 2, 3},
			expectedValue: 3,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			for _, v := range tt.values {
				list.PushBack(v)
			}
			back := list.Back()
			if tt.expectNil {
				if back != nil {
					t.Errorf("Back() = %v; expected nil", back)
				}
			} else {
				if back == nil {
					t.Errorf("Back() = nil; expected non-nil")
					return
				}
				if back.Value != tt.expectedValue {
					t.Errorf("Back().Value = %v; expected %v", back.Value, tt.expectedValue)
				}
			}
		})
	}
}

// TestRemove tests the Remove function
func TestRemove(t *testing.T) {
	tests := []struct {
		name            string
		values          []int
		removeIndex     int
		expectedLen     int
		expectedRemoved int
	}{
		{
			name:            "remove from single element list",
			values:          []int{1},
			removeIndex:     0,
			expectedLen:     0,
			expectedRemoved: 1,
		},
		{
			name:            "remove first element",
			values:          []int{1, 2, 3},
			removeIndex:     0,
			expectedLen:     2,
			expectedRemoved: 1,
		},
		{
			name:            "remove middle element",
			values:          []int{1, 2, 3},
			removeIndex:     1,
			expectedLen:     2,
			expectedRemoved: 2,
		},
		{
			name:            "remove last element",
			values:          []int{1, 2, 3},
			removeIndex:     2,
			expectedLen:     2,
			expectedRemoved: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.values {
				nodes = append(nodes, list.PushBack(v))
			}

			removed := list.Remove(nodes[tt.removeIndex])
			if removed != tt.expectedRemoved {
				t.Errorf("Remove() = %v; expected %v", removed, tt.expectedRemoved)
			}
			if list.Len() != tt.expectedLen {
				t.Errorf("Remove() list length = %d; expected %d", list.Len(), tt.expectedLen)
			}
		})
	}
}

// TestLen tests the Len function
func TestLen(t *testing.T) {
	tests := []struct {
		name        string
		values      []int
		expectedLen int
	}{
		{
			name:        "empty list",
			values:      []int{},
			expectedLen: 0,
		},
		{
			name:        "single element",
			values:      []int{1},
			expectedLen: 1,
		},
		{
			name:        "multiple elements",
			values:      []int{1, 2, 3, 4, 5},
			expectedLen: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			for _, v := range tt.values {
				list.PushBack(v)
			}
			if list.Len() != tt.expectedLen {
				t.Errorf("Len() = %d; expected %d", list.Len(), tt.expectedLen)
			}
		})
	}
}

// TestInsertBefore tests the InsertBefore function
func TestInsertBefore(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		insertValue    int
		markIndex      int
		expectedValues []int
	}{
		{
			name:           "insert before first element",
			initialValues:  []int{2, 3},
			insertValue:    1,
			markIndex:      0,
			expectedValues: []int{1, 2, 3},
		},
		{
			name:           "insert before middle element",
			initialValues:  []int{1, 3},
			insertValue:    2,
			markIndex:      1,
			expectedValues: []int{1, 2, 3},
		},
		{
			name:           "insert before last element",
			initialValues:  []int{1, 2, 4},
			insertValue:    3,
			markIndex:      2,
			expectedValues: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.initialValues {
				nodes = append(nodes, list.PushBack(v))
			}

			list.InsertBefore(tt.insertValue, nodes[tt.markIndex])

			if list.Len() != len(tt.expectedValues) {
				t.Errorf("InsertBefore() list length = %d; expected %d", list.Len(), len(tt.expectedValues))
			}

			i := 0
			for e := list.Front(); e != nil; e = e.Next() {
				if e.Value != tt.expectedValues[i] {
					t.Errorf("InsertBefore() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestInsertAfter tests the InsertAfter function
func TestInsertAfter(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		insertValue    int
		markIndex      int
		expectedValues []int
	}{
		{
			name:           "insert after first element",
			initialValues:  []int{1, 3},
			insertValue:    2,
			markIndex:      0,
			expectedValues: []int{1, 2, 3},
		},
		{
			name:           "insert after middle element",
			initialValues:  []int{1, 2, 4},
			insertValue:    3,
			markIndex:      1,
			expectedValues: []int{1, 2, 3, 4},
		},
		{
			name:           "insert after last element",
			initialValues:  []int{1, 2},
			insertValue:    3,
			markIndex:      1,
			expectedValues: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.initialValues {
				nodes = append(nodes, list.PushBack(v))
			}

			list.InsertAfter(tt.insertValue, nodes[tt.markIndex])

			if list.Len() != len(tt.expectedValues) {
				t.Errorf("InsertAfter() list length = %d; expected %d", list.Len(), len(tt.expectedValues))
			}

			i := 0
			for e := list.Front(); e != nil; e = e.Next() {
				if e.Value != tt.expectedValues[i] {
					t.Errorf("InsertAfter() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestMoveToFront tests the MoveToFront function
func TestMoveToFront(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		moveIndex      int
		expectedValues []int
	}{
		{
			name:           "move last to front",
			initialValues:  []int{1, 2, 3},
			moveIndex:      2,
			expectedValues: []int{3, 1, 2},
		},
		{
			name:           "move middle to front",
			initialValues:  []int{1, 2, 3},
			moveIndex:      1,
			expectedValues: []int{2, 1, 3},
		},
		{
			name:           "move already front",
			initialValues:  []int{1, 2, 3},
			moveIndex:      0,
			expectedValues: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.initialValues {
				nodes = append(nodes, list.PushBack(v))
			}

			list.MoveToFront(nodes[tt.moveIndex])

			i := 0
			for e := list.Front(); e != nil; e = e.Next() {
				if e.Value != tt.expectedValues[i] {
					t.Errorf("MoveToFront() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestMoveToBack tests the MoveToBack function
func TestMoveToBack(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		moveIndex      int
		expectedValues []int
	}{
		{
			name:           "move first to back",
			initialValues:  []int{1, 2, 3},
			moveIndex:      0,
			expectedValues: []int{2, 3, 1},
		},
		{
			name:           "move middle to back",
			initialValues:  []int{1, 2, 3},
			moveIndex:      1,
			expectedValues: []int{1, 3, 2},
		},
		{
			name:           "move already back",
			initialValues:  []int{1, 2, 3},
			moveIndex:      2,
			expectedValues: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.initialValues {
				nodes = append(nodes, list.PushBack(v))
			}

			list.MoveToBack(nodes[tt.moveIndex])

			i := 0
			for e := list.Front(); e != nil; e = e.Next() {
				if e.Value != tt.expectedValues[i] {
					t.Errorf("MoveToBack() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestMoveBefore tests the MoveBefore function
func TestMoveBefore(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		moveIndex      int
		markIndex      int
		expectedValues []int
	}{
		{
			name:           "move last before first",
			initialValues:  []int{1, 2, 3},
			moveIndex:      2,
			markIndex:      0,
			expectedValues: []int{3, 1, 2},
		},
		{
			name:           "move first before last",
			initialValues:  []int{1, 2, 3},
			moveIndex:      0,
			markIndex:      2,
			expectedValues: []int{2, 1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.initialValues {
				nodes = append(nodes, list.PushBack(v))
			}

			list.MoveBefore(nodes[tt.moveIndex], nodes[tt.markIndex])

			i := 0
			for e := list.Front(); e != nil; e = e.Next() {
				if e.Value != tt.expectedValues[i] {
					t.Errorf("MoveBefore() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestMoveAfter tests the MoveAfter function
func TestMoveAfter(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		moveIndex      int
		markIndex      int
		expectedValues []int
	}{
		{
			name:           "move first after last",
			initialValues:  []int{1, 2, 3},
			moveIndex:      0,
			markIndex:      2,
			expectedValues: []int{2, 3, 1},
		},
		{
			name:           "move last after first",
			initialValues:  []int{1, 2, 3},
			moveIndex:      2,
			markIndex:      0,
			expectedValues: []int{1, 3, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.initialValues {
				nodes = append(nodes, list.PushBack(v))
			}

			list.MoveAfter(nodes[tt.moveIndex], nodes[tt.markIndex])

			i := 0
			for e := list.Front(); e != nil; e = e.Next() {
				if e.Value != tt.expectedValues[i] {
					t.Errorf("MoveAfter() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestPushBackList tests the PushBackList function
func TestPushBackList(t *testing.T) {
	tests := []struct {
		name           string
		list1Values    []int
		list2Values    []int
		expectedValues []int
		expectedLen    int
	}{
		{
			name:           "push empty list",
			list1Values:    []int{1, 2},
			list2Values:    []int{},
			expectedValues: []int{1, 2},
			expectedLen:    2,
		},
		{
			name:           "push to empty list",
			list1Values:    []int{},
			list2Values:    []int{1, 2},
			expectedValues: []int{1, 2},
			expectedLen:    2,
		},
		{
			name:           "push non-empty lists",
			list1Values:    []int{1, 2},
			list2Values:    []int{3, 4},
			expectedValues: []int{1, 2, 3, 4},
			expectedLen:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list1 := NewDLList[int]()
			for _, v := range tt.list1Values {
				list1.PushBack(v)
			}

			list2 := NewDLList[int]()
			for _, v := range tt.list2Values {
				list2.PushBack(v)
			}

			list1.PushBackList(list2)

			if list1.Len() != tt.expectedLen {
				t.Errorf("PushBackList() list length = %d; expected %d", list1.Len(), tt.expectedLen)
			}

			i := 0
			for e := list1.Front(); e != nil; e = e.Next() {
				if i >= len(tt.expectedValues) {
					t.Errorf("PushBackList() has more elements than expected")
					break
				}
				if e.Value != tt.expectedValues[i] {
					t.Errorf("PushBackList() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestPushFrontList tests the PushFrontList function
func TestPushFrontList(t *testing.T) {
	tests := []struct {
		name           string
		list1Values    []int
		list2Values    []int
		expectedValues []int
		expectedLen    int
	}{
		{
			name:           "push empty list",
			list1Values:    []int{3, 4},
			list2Values:    []int{},
			expectedValues: []int{3, 4},
			expectedLen:    2,
		},
		{
			name:           "push to empty list",
			list1Values:    []int{},
			list2Values:    []int{1, 2},
			expectedValues: []int{1, 2},
			expectedLen:    2,
		},
		{
			name:           "push non-empty lists",
			list1Values:    []int{3, 4},
			list2Values:    []int{1, 2},
			expectedValues: []int{1, 2, 3, 4},
			expectedLen:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list1 := NewDLList[int]()
			for _, v := range tt.list1Values {
				list1.PushBack(v)
			}

			list2 := NewDLList[int]()
			for _, v := range tt.list2Values {
				list2.PushBack(v)
			}

			list1.PushFrontList(list2)

			if list1.Len() != tt.expectedLen {
				t.Errorf("PushFrontList() list length = %d; expected %d", list1.Len(), tt.expectedLen)
			}

			i := 0
			for e := list1.Front(); e != nil; e = e.Next() {
				if i >= len(tt.expectedValues) {
					t.Errorf("PushFrontList() has more elements than expected")
					break
				}
				if e.Value != tt.expectedValues[i] {
					t.Errorf("PushFrontList() element[%d] = %v; expected %v", i, e.Value, tt.expectedValues[i])
				}
				i++
			}
		})
	}
}

// TestNodeNext tests the Next method on DLLNode
func TestNodeNext(t *testing.T) {
	tests := []struct {
		name         string
		values       []int
		startIndex   int
		expectedNext *int
	}{
		{
			name:         "next of first element",
			values:       []int{1, 2, 3},
			startIndex:   0,
			expectedNext: intPtr(2),
		},
		{
			name:         "next of middle element",
			values:       []int{1, 2, 3},
			startIndex:   1,
			expectedNext: intPtr(3),
		},
		{
			name:         "next of last element",
			values:       []int{1, 2, 3},
			startIndex:   2,
			expectedNext: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.values {
				nodes = append(nodes, list.PushBack(v))
			}

			next := nodes[tt.startIndex].Next()
			if tt.expectedNext == nil {
				if next != nil {
					t.Errorf("Next() = %v; expected nil", next)
				}
			} else {
				if next == nil {
					t.Errorf("Next() = nil; expected non-nil")
					return
				}
				if next.Value != *tt.expectedNext {
					t.Errorf("Next().Value = %v; expected %v", next.Value, *tt.expectedNext)
				}
			}
		})
	}
}

// TestNodePrev tests the Prev method on DLLNode
func TestNodePrev(t *testing.T) {
	tests := []struct {
		name         string
		values       []int
		startIndex   int
		expectedPrev *int
	}{
		{
			name:         "prev of last element",
			values:       []int{1, 2, 3},
			startIndex:   2,
			expectedPrev: intPtr(2),
		},
		{
			name:         "prev of middle element",
			values:       []int{1, 2, 3},
			startIndex:   1,
			expectedPrev: intPtr(1),
		},
		{
			name:         "prev of first element",
			values:       []int{1, 2, 3},
			startIndex:   0,
			expectedPrev: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDLList[int]()
			var nodes []*DLLNode[int]
			for _, v := range tt.values {
				nodes = append(nodes, list.PushBack(v))
			}

			prev := nodes[tt.startIndex].Prev()
			if tt.expectedPrev == nil {
				if prev != nil {
					t.Errorf("Prev() = %v; expected nil", prev)
				}
			} else {
				if prev == nil {
					t.Errorf("Prev() = nil; expected non-nil")
					return
				}
				if prev.Value != *tt.expectedPrev {
					t.Errorf("Prev().Value = %v; expected %v", prev.Value, *tt.expectedPrev)
				}
			}
		})
	}
}

// TestGenericTypes tests that the list works with different types
func TestGenericTypes(t *testing.T) {
	t.Run("string list", func(t *testing.T) {
		list := NewDLList[string]()
		list.PushBack("hello")
		list.PushBack("world")

		if list.Len() != 2 {
			t.Errorf("string list length = %d; expected 2", list.Len())
		}
		if list.Front().Value != "hello" {
			t.Errorf("Front().Value = %s; expected 'hello'", list.Front().Value)
		}
	})

	t.Run("struct list", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		list := NewDLList[Person]()
		list.PushBack(Person{Name: "Alice", Age: 30})
		list.PushBack(Person{Name: "Bob", Age: 25})

		if list.Len() != 2 {
			t.Errorf("struct list length = %d; expected 2", list.Len())
		}
		if list.Front().Value.Name != "Alice" {
			t.Errorf("Front().Value.Name = %s; expected 'Alice'", list.Front().Value.Name)
		}
	})

	t.Run("pointer list", func(t *testing.T) {
		list := NewDLList[*int]()
		val1 := 10
		val2 := 20
		list.PushBack(&val1)
		list.PushBack(&val2)

		if list.Len() != 2 {
			t.Errorf("pointer list length = %d; expected 2", list.Len())
		}
		if *list.Front().Value != 10 {
			t.Errorf("*Front().Value = %d; expected 10", *list.Front().Value)
		}
	})
}

// Helper function to create pointer to int
func intPtr(i int) *int {
	return &i
}
