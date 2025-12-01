package utilities

import (
	"reflect"
	"strings"
	"testing"
)

// TestFilter tests the Filter function with various predicates
func TestFilter(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int, int) bool
		expected   []int
	}{
		{
			name:       "filter even numbers",
			collection: []int{1, 2, 3, 4, 5, 6},
			predicate:  func(item, index int) bool { return item%2 == 0 },
			expected:   []int{2, 4, 6},
		},
		{
			name:       "filter by index",
			collection: []int{10, 20, 30, 40, 50},
			predicate:  func(item, index int) bool { return index%2 == 0 },
			expected:   []int{10, 30, 50},
		},
		{
			name:       "filter greater than value",
			collection: []int{1, 5, 10, 15, 20},
			predicate:  func(item, index int) bool { return item > 10 },
			expected:   []int{15, 20},
		},
		{
			name:       "filter empty collection",
			collection: []int{},
			predicate:  func(item, index int) bool { return true },
			expected:   []int{},
		},
		{
			name:       "filter none matching",
			collection: []int{1, 2, 3},
			predicate:  func(item, index int) bool { return item > 100 },
			expected:   []int{},
		},
		{
			name:       "filter all matching",
			collection: []int{2, 4, 6, 8},
			predicate:  func(item, index int) bool { return item%2 == 0 },
			expected:   []int{2, 4, 6, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.collection, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestMap tests the Map function
func TestMap(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		iteratee   func(int, int) int
		expected   []int
	}{
		{
			name:       "double values",
			collection: []int{1, 2, 3, 4, 5},
			iteratee:   func(item, index int) int { return item * 2 },
			expected:   []int{2, 4, 6, 8, 10},
		},
		{
			name:       "add index to value",
			collection: []int{10, 20, 30},
			iteratee:   func(item, index int) int { return item + index },
			expected:   []int{10, 21, 32},
		},
		{
			name:       "empty collection",
			collection: []int{},
			iteratee:   func(item, index int) int { return item * 2 },
			expected:   []int{},
		},
		{
			name:       "single element",
			collection: []int{5},
			iteratee:   func(item, index int) int { return item * 3 },
			expected:   []int{15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.collection, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Map() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestMapToString tests Map with type transformation
func TestMapToString(t *testing.T) {
	t.Run("convert int to string", func(t *testing.T) {
		collection := []int{1, 2, 3}
		result := Map(collection, func(item, index int) string {
			return string(rune('0' + item))
		})
		expected := []string{"1", "2", "3"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Map() = %v; expected %v", result, expected)
		}
	})
}

// TestFilterMap tests the FilterMap function
func TestFilterMap(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		callback   func(int, int) (int, bool)
		expected   []int
	}{
		{
			name:       "filter and double even numbers",
			collection: []int{1, 2, 3, 4, 5, 6},
			callback: func(item, index int) (int, bool) {
				if item%2 == 0 {
					return item * 2, true
				}
				return 0, false
			},
			expected: []int{4, 8, 12},
		},
		{
			name:       "filter by index and square",
			collection: []int{10, 20, 30, 40, 50},
			callback: func(item, index int) (int, bool) {
				if index%2 == 0 {
					return item * item, true
				}
				return 0, false
			},
			expected: []int{100, 900, 2500},
		},
		{
			name:       "empty collection",
			collection: []int{},
			callback: func(item, index int) (int, bool) {
				return item, true
			},
			expected: []int{},
		},
		{
			name:       "none matching",
			collection: []int{1, 2, 3},
			callback: func(item, index int) (int, bool) {
				return item, false
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterMap(tt.collection, tt.callback)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FilterMap() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestFlatMap tests the FlatMap function
func TestFlatMap(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		iteratee   func(int, int) []int
		expected   []int
	}{
		{
			name:       "duplicate each element",
			collection: []int{1, 2, 3},
			iteratee:   func(item, index int) []int { return []int{item, item} },
			expected:   []int{1, 1, 2, 2, 3, 3},
		},
		{
			name:       "expand to range",
			collection: []int{1, 2},
			iteratee:   func(item, index int) []int { return []int{item, item + 10} },
			expected:   []int{1, 11, 2, 12},
		},
		{
			name:       "empty collection",
			collection: []int{},
			iteratee:   func(item, index int) []int { return []int{item} },
			expected:   []int{},
		},
		{
			name:       "return nil slices",
			collection: []int{1, 2, 3},
			iteratee:   func(item, index int) []int { return nil },
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlatMap(tt.collection, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FlatMap() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestReduce tests the Reduce function
func TestReduce(t *testing.T) {
	tests := []struct {
		name        string
		collection  []int
		accumulator func(int, int, int) int
		initial     int
		expected    int
	}{
		{
			name:        "sum all elements",
			collection:  []int{1, 2, 3, 4, 5},
			accumulator: func(agg, item, index int) int { return agg + item },
			initial:     0,
			expected:    15,
		},
		{
			name:        "multiply all elements",
			collection:  []int{2, 3, 4},
			accumulator: func(agg, item, index int) int { return agg * item },
			initial:     1,
			expected:    24,
		},
		{
			name:        "empty collection",
			collection:  []int{},
			accumulator: func(agg, item, index int) int { return agg + item },
			initial:     10,
			expected:    10,
		},
		{
			name:        "count elements",
			collection:  []int{1, 2, 3},
			accumulator: func(agg, item, index int) int { return agg + 1 },
			initial:     0,
			expected:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reduce(tt.collection, tt.accumulator, tt.initial)
			if result != tt.expected {
				t.Errorf("Reduce() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestReduceRight tests the ReduceRight function
func TestReduceRight(t *testing.T) {
	tests := []struct {
		name        string
		collection  []string
		accumulator func(string, string, int) string
		initial     string
		expected    string
	}{
		{
			name:        "concatenate from right",
			collection:  []string{"a", "b", "c"},
			accumulator: func(agg string, item string, index int) string { return agg + item },
			initial:     "",
			expected:    "cba",
		},
		{
			name:        "concatenate with separator from right",
			collection:  []string{"1", "2", "3"},
			accumulator: func(agg string, item string, index int) string { return agg + "-" + item },
			initial:     "start",
			expected:    "start-3-2-1",
		},
		{
			name:        "empty collection",
			collection:  []string{},
			accumulator: func(agg string, item string, index int) string { return agg + item },
			initial:     "initial",
			expected:    "initial",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReduceRight(tt.collection, tt.accumulator, tt.initial)
			if result != tt.expected {
				t.Errorf("ReduceRight() = %s; expected %s", result, tt.expected)
			}
		})
	}
}

// TestForEach tests the ForEach function
func TestForEach(t *testing.T) {
	t.Run("iterate and accumulate", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		sum := 0
		ForEach(collection, func(item, index int) {
			sum += item
		})
		if sum != 15 {
			t.Errorf("ForEach() sum = %d; expected 15", sum)
		}
	})

	t.Run("empty collection", func(t *testing.T) {
		collection := []int{}
		called := false
		ForEach(collection, func(item, index int) {
			called = true
		})
		if called {
			t.Error("ForEach() should not call iteratee for empty collection")
		}
	})
}

// TestForEachWhile tests the ForEachWhile function
func TestForEachWhile(t *testing.T) {
	tests := []struct {
		name          string
		collection    []int
		iteratee      func(int, int) bool
		expectedCalls int
	}{
		{
			name:       "stop on first false",
			collection: []int{1, 2, 3, 4, 5},
			iteratee: func(item, index int) bool {
				return item < 3
			},
			expectedCalls: 3,
		},
		{
			name:       "process all elements",
			collection: []int{1, 2, 3},
			iteratee: func(item, index int) bool {
				return true
			},
			expectedCalls: 3,
		},
		{
			name:       "stop immediately",
			collection: []int{1, 2, 3},
			iteratee: func(item, index int) bool {
				return false
			},
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			ForEachWhile(tt.collection, func(item, index int) bool {
				calls++
				return tt.iteratee(item, index)
			})
			if calls != tt.expectedCalls {
				t.Errorf("ForEachWhile() calls = %d; expected %d", calls, tt.expectedCalls)
			}
		})
	}
}

// TestTimes tests the Times function
func TestTimes(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		iteratee func(int) int
		expected []int
	}{
		{
			name:     "create sequence",
			count:    5,
			iteratee: func(i int) int { return i },
			expected: []int{0, 1, 2, 3, 4},
		},
		{
			name:     "create doubled sequence",
			count:    3,
			iteratee: func(i int) int { return i * 2 },
			expected: []int{0, 2, 4},
		},
		{
			name:     "zero count",
			count:    0,
			iteratee: func(i int) int { return i },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Times(tt.count, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Times() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestUniq tests the Uniq function
func TestUniq(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   []int
	}{
		{
			name:       "remove duplicates",
			collection: []int{1, 2, 2, 3, 3, 3, 4},
			expected:   []int{1, 2, 3, 4},
		},
		{
			name:       "no duplicates",
			collection: []int{1, 2, 3, 4, 5},
			expected:   []int{1, 2, 3, 4, 5},
		},
		{
			name:       "all duplicates",
			collection: []int{1, 1, 1, 1},
			expected:   []int{1},
		},
		{
			name:       "empty collection",
			collection: []int{},
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uniq(tt.collection)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Uniq() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestUniqBy tests the UniqBy function
func TestUniqBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name       string
		collection []Person
		iteratee   func(Person) int
		expected   int
	}{
		{
			name: "unique by age",
			collection: []Person{
				{"Alice", 25},
				{"Bob", 30},
				{"Charlie", 25},
				{"David", 30},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: 2,
		},
		{
			name: "all unique",
			collection: []Person{
				{"Alice", 25},
				{"Bob", 30},
				{"Charlie", 35},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqBy(tt.collection, tt.iteratee)
			if len(result) != tt.expected {
				t.Errorf("UniqBy() length = %d; expected %d", len(result), tt.expected)
			}
		})
	}
}

// TestGroupBy tests the GroupBy function
func TestGroupBy(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		iteratee   func(int) string
		expected   map[string][]int
	}{
		{
			name:       "group by even/odd",
			collection: []int{1, 2, 3, 4, 5, 6},
			iteratee: func(i int) string {
				if i%2 == 0 {
					return "even"
				}
				return "odd"
			},
			expected: map[string][]int{
				"even": {2, 4, 6},
				"odd":  {1, 3, 5},
			},
		},
		{
			name:       "empty collection",
			collection: []int{},
			iteratee:   func(i int) string { return "group" },
			expected:   map[string][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GroupBy(tt.collection, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GroupBy() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestChunk tests the Chunk function
func TestChunk(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		size       int
		expected   [][]int
	}{
		{
			name:       "even division",
			collection: []int{1, 2, 3, 4, 5, 6},
			size:       2,
			expected:   [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:       "uneven division",
			collection: []int{1, 2, 3, 4, 5},
			size:       2,
			expected:   [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:       "size larger than collection",
			collection: []int{1, 2, 3},
			size:       10,
			expected:   [][]int{{1, 2, 3}},
		},
		{
			name:       "size of 1",
			collection: []int{1, 2, 3},
			size:       1,
			expected:   [][]int{{1}, {2}, {3}},
		},
		{
			name:       "empty collection",
			collection: []int{},
			size:       2,
			expected:   [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Chunk(tt.collection, tt.size)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Chunk() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestChunkPanic tests that Chunk panics with invalid size
func TestChunkPanic(t *testing.T) {
	t.Run("panic on zero size", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Chunk() should panic with size 0")
			}
		}()
		Chunk([]int{1, 2, 3}, 0)
	})

	t.Run("panic on negative size", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Chunk() should panic with negative size")
			}
		}()
		Chunk([]int{1, 2, 3}, -1)
	})
}

// TestFlatten tests the Flatten function
func TestFlatten(t *testing.T) {
	tests := []struct {
		name       string
		collection [][]int
		expected   []int
	}{
		{
			name:       "flatten nested slices",
			collection: [][]int{{1, 2}, {3, 4}, {5, 6}},
			expected:   []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:       "flatten uneven slices",
			collection: [][]int{{1}, {2, 3}, {4, 5, 6}},
			expected:   []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:       "flatten empty nested",
			collection: [][]int{{}, {1, 2}, {}},
			expected:   []int{1, 2},
		},
		{
			name:       "empty collection",
			collection: [][]int{},
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.collection)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Flatten() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestShuffle tests the Shuffle function
func TestShuffle(t *testing.T) {
	t.Run("shuffle maintains length", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := Shuffle(append([]int{}, collection...)) // copy to avoid modifying original
		if len(result) != len(collection) {
			t.Errorf("Shuffle() length = %d; expected %d", len(result), len(collection))
		}
	})

	t.Run("shuffle maintains elements", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := Shuffle(append([]int{}, collection...))

		// Check all original elements are present
		counts := make(map[int]int)
		for _, v := range result {
			counts[v]++
		}
		for _, v := range collection {
			if counts[v] != 1 {
				t.Errorf("Shuffle() missing or duplicate element %d", v)
			}
		}
	})
}

// TestReverse tests the Reverse function
func TestReverse(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   []int
	}{
		{
			name:       "reverse odd length",
			collection: []int{1, 2, 3, 4, 5},
			expected:   []int{5, 4, 3, 2, 1},
		},
		{
			name:       "reverse even length",
			collection: []int{1, 2, 3, 4},
			expected:   []int{4, 3, 2, 1},
		},
		{
			name:       "reverse single element",
			collection: []int{1},
			expected:   []int{1},
		},
		{
			name:       "reverse empty",
			collection: []int{},
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying original
			collection := make([]int, len(tt.collection))
			copy(collection, tt.collection)

			result := Reverse(collection)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Reverse() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestKeyBy tests the KeyBy function
func TestKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	tests := []struct {
		name       string
		collection []User
		iteratee   func(User) int
		expected   map[int]User
	}{
		{
			name: "key by ID",
			collection: []User{
				{1, "Alice"},
				{2, "Bob"},
				{3, "Charlie"},
			},
			iteratee: func(u User) int { return u.ID },
			expected: map[int]User{
				1: {1, "Alice"},
				2: {2, "Bob"},
				3: {3, "Charlie"},
			},
		},
		{
			name:       "empty collection",
			collection: []User{},
			iteratee:   func(u User) int { return u.ID },
			expected:   map[int]User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KeyBy(tt.collection, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("KeyBy() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestDrop tests the Drop function
func TestDrop(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		n          int
		expected   []int
	}{
		{
			name:       "drop first 2",
			collection: []int{1, 2, 3, 4, 5},
			n:          2,
			expected:   []int{3, 4, 5},
		},
		{
			name:       "drop more than length",
			collection: []int{1, 2, 3},
			n:          5,
			expected:   []int{},
		},
		{
			name:       "drop 0",
			collection: []int{1, 2, 3},
			n:          0,
			expected:   []int{1, 2, 3},
		},
		{
			name:       "drop from empty",
			collection: []int{},
			n:          2,
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Drop(tt.collection, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Drop() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestDropRight tests the DropRight function
func TestDropRight(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		n          int
		expected   []int
	}{
		{
			name:       "drop last 2",
			collection: []int{1, 2, 3, 4, 5},
			n:          2,
			expected:   []int{1, 2, 3},
		},
		{
			name:       "drop more than length",
			collection: []int{1, 2, 3},
			n:          5,
			expected:   []int{},
		},
		{
			name:       "drop 0",
			collection: []int{1, 2, 3},
			n:          0,
			expected:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DropRight(tt.collection, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DropRight() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestDropWhile tests the DropWhile function
func TestDropWhile(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int) bool
		expected   []int
	}{
		{
			name:       "drop while less than 3",
			collection: []int{1, 2, 3, 4, 2, 1},
			predicate:  func(i int) bool { return i < 3 },
			expected:   []int{3, 4, 2, 1},
		},
		{
			name:       "drop none",
			collection: []int{3, 4, 5},
			predicate:  func(i int) bool { return i < 3 },
			expected:   []int{3, 4, 5},
		},
		{
			name:       "drop all",
			collection: []int{1, 2, 3},
			predicate:  func(i int) bool { return i < 10 },
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DropWhile(tt.collection, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DropWhile() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestReject tests the Reject function
func TestReject(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int, int) bool
		expected   []int
	}{
		{
			name:       "reject even numbers",
			collection: []int{1, 2, 3, 4, 5, 6},
			predicate:  func(item, index int) bool { return item%2 == 0 },
			expected:   []int{1, 3, 5},
		},
		{
			name:       "reject none",
			collection: []int{1, 2, 3},
			predicate:  func(item, index int) bool { return item > 100 },
			expected:   []int{1, 2, 3},
		},
		{
			name:       "reject all",
			collection: []int{2, 4, 6},
			predicate:  func(item, index int) bool { return item%2 == 0 },
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reject(tt.collection, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Reject() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestCount tests the Count function
func TestCount(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		value      int
		expected   int
	}{
		{
			name:       "count occurrences",
			collection: []int{1, 2, 2, 3, 2, 4},
			value:      2,
			expected:   3,
		},
		{
			name:       "count not found",
			collection: []int{1, 2, 3},
			value:      5,
			expected:   0,
		},
		{
			name:       "count in empty",
			collection: []int{},
			value:      1,
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Count(tt.collection, tt.value)
			if result != tt.expected {
				t.Errorf("Count() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestCountBy tests the CountBy function
func TestCountBy(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int) bool
		expected   int
	}{
		{
			name:       "count even numbers",
			collection: []int{1, 2, 3, 4, 5, 6},
			predicate:  func(i int) bool { return i%2 == 0 },
			expected:   3,
		},
		{
			name:       "count none matching",
			collection: []int{1, 3, 5},
			predicate:  func(i int) bool { return i%2 == 0 },
			expected:   0,
		},
		{
			name:       "count all matching",
			collection: []int{2, 4, 6},
			predicate:  func(i int) bool { return i%2 == 0 },
			expected:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountBy(tt.collection, tt.predicate)
			if result != tt.expected {
				t.Errorf("CountBy() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestCountValues tests the CountValues function
func TestCountValues(t *testing.T) {
	tests := []struct {
		name       string
		collection []string
		expected   map[string]int
	}{
		{
			name:       "count occurrences",
			collection: []string{"a", "b", "a", "c", "b", "a"},
			expected:   map[string]int{"a": 3, "b": 2, "c": 1},
		},
		{
			name:       "all unique",
			collection: []string{"a", "b", "c"},
			expected:   map[string]int{"a": 1, "b": 1, "c": 1},
		},
		{
			name:       "empty collection",
			collection: []string{},
			expected:   map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountValues(tt.collection)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CountValues() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestSubset tests the Subset function
func TestSubset(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		offset     int
		length     uint
		expected   []int
	}{
		{
			name:       "normal subset",
			collection: []int{1, 2, 3, 4, 5},
			offset:     1,
			length:     3,
			expected:   []int{2, 3, 4},
		},
		{
			name:       "negative offset",
			collection: []int{1, 2, 3, 4, 5},
			offset:     -2,
			length:     2,
			expected:   []int{4, 5},
		},
		{
			name:       "length overflow",
			collection: []int{1, 2, 3},
			offset:     1,
			length:     10,
			expected:   []int{2, 3},
		},
		{
			name:       "offset overflow",
			collection: []int{1, 2, 3},
			offset:     10,
			length:     2,
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subset(tt.collection, tt.offset, tt.length)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Subset() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestSlice tests the Slice function
func TestSlice(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		start      int
		end        int
		expected   []int
	}{
		{
			name:       "normal slice",
			collection: []int{1, 2, 3, 4, 5},
			start:      1,
			end:        4,
			expected:   []int{2, 3, 4},
		},
		{
			name:       "start >= end",
			collection: []int{1, 2, 3, 4, 5},
			start:      3,
			end:        2,
			expected:   []int{},
		},
		{
			name:       "end overflow",
			collection: []int{1, 2, 3},
			start:      1,
			end:        10,
			expected:   []int{2, 3},
		},
		{
			name:       "start overflow",
			collection: []int{1, 2, 3},
			start:      10,
			end:        15,
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slice(tt.collection, tt.start, tt.end)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Slice() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestReplace tests the Replace function
func TestReplace(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		old        int
		new        int
		n          int
		expected   []int
	}{
		{
			name:       "replace first 2",
			collection: []int{1, 2, 2, 3, 2},
			old:        2,
			new:        9,
			n:          2,
			expected:   []int{1, 9, 9, 3, 2},
		},
		{
			name:       "replace all with -1",
			collection: []int{1, 2, 2, 3, 2},
			old:        2,
			new:        9,
			n:          -1,
			expected:   []int{1, 9, 9, 3, 9},
		},
		{
			name:       "replace none",
			collection: []int{1, 2, 3},
			old:        5,
			new:        9,
			n:          1,
			expected:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Replace(tt.collection, tt.old, tt.new, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Replace() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestReplaceAll tests the ReplaceAll function
func TestReplaceAll(t *testing.T) {
	t.Run("replace all occurrences", func(t *testing.T) {
		collection := []int{1, 2, 2, 3, 2}
		result := ReplaceAll(collection, 2, 9)
		expected := []int{1, 9, 9, 3, 9}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("ReplaceAll() = %v; expected %v", result, expected)
		}
	})
}

// TestCompact tests the Compact function
func TestCompact(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   []int
	}{
		{
			name:       "remove zeros",
			collection: []int{1, 0, 2, 0, 3, 0, 4},
			expected:   []int{1, 2, 3, 4},
		},
		{
			name:       "no zeros",
			collection: []int{1, 2, 3, 4},
			expected:   []int{1, 2, 3, 4},
		},
		{
			name:       "all zeros",
			collection: []int{0, 0, 0},
			expected:   []int{},
		},
		{
			name:       "empty collection",
			collection: []int{},
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Compact(tt.collection)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Compact() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestCompactStrings tests Compact with strings
func TestCompactStrings(t *testing.T) {
	t.Run("remove empty strings", func(t *testing.T) {
		collection := []string{"a", "", "b", "", "c"}
		result := Compact(collection)
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Compact() = %v; expected %v", result, expected)
		}
	})
}

// TestIsSorted tests the IsSorted function
func TestIsSorted(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   bool
	}{
		{
			name:       "sorted ascending",
			collection: []int{1, 2, 3, 4, 5},
			expected:   true,
		},
		{
			name:       "not sorted",
			collection: []int{1, 3, 2, 4, 5},
			expected:   false,
		},
		{
			name:       "sorted with duplicates",
			collection: []int{1, 2, 2, 3, 4},
			expected:   true,
		},
		{
			name:       "empty collection",
			collection: []int{},
			expected:   true,
		},
		{
			name:       "single element",
			collection: []int{1},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSorted(tt.collection)
			if result != tt.expected {
				t.Errorf("IsSorted() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestIsSortedByKey tests the IsSortedByKey function
func TestIsSortedByKey(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name       string
		collection []Person
		iteratee   func(Person) int
		expected   bool
	}{
		{
			name: "sorted by age",
			collection: []Person{
				{"Alice", 25},
				{"Bob", 30},
				{"Charlie", 35},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: true,
		},
		{
			name: "not sorted by age",
			collection: []Person{
				{"Alice", 30},
				{"Bob", 25},
				{"Charlie", 35},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSortedByKey(tt.collection, tt.iteratee)
			if result != tt.expected {
				t.Errorf("IsSortedByKey() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestSplice tests the Splice function
func TestSplice(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		i          int
		elements   []int
		expected   []int
	}{
		{
			name:       "insert in middle",
			collection: []int{1, 2, 5, 6},
			i:          2,
			elements:   []int{3, 4},
			expected:   []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:       "insert at beginning",
			collection: []int{3, 4, 5},
			i:          0,
			elements:   []int{1, 2},
			expected:   []int{1, 2, 3, 4, 5},
		},
		{
			name:       "insert at end",
			collection: []int{1, 2, 3},
			i:          3,
			elements:   []int{4, 5},
			expected:   []int{1, 2, 3, 4, 5},
		},
		{
			name:       "negative index",
			collection: []int{1, 2, 3, 4},
			i:          -2,
			elements:   []int{9},
			expected:   []int{1, 2, 9, 3, 4},
		},
		{
			name:       "no elements",
			collection: []int{1, 2, 3},
			i:          1,
			elements:   []int{},
			expected:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Splice(tt.collection, tt.i, tt.elements...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Splice() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestFilterMapWithStrings tests FilterMap with string transformation
func TestFilterMapWithStrings(t *testing.T) {
	t.Run("filter and uppercase", func(t *testing.T) {
		collection := []string{"apple", "banana", "apricot", "cherry"}
		result := FilterMap(collection, func(item string, index int) (string, bool) {
			if strings.HasPrefix(item, "a") {
				return strings.ToUpper(item), true
			}
			return "", false
		})
		expected := []string{"APPLE", "APRICOT"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("FilterMap() = %v; expected %v", result, expected)
		}
	})
}

// BenchmarkFilter benchmarks the Filter function
func BenchmarkFilter(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}
	predicate := func(item, index int) bool { return item%2 == 0 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Filter(collection, predicate)
	}
}

// BenchmarkMap benchmarks the Map function
func BenchmarkMap(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}
	iteratee := func(item, index int) int { return item * 2 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(collection, iteratee)
	}
}

// BenchmarkUniq benchmarks the Uniq function
func BenchmarkUniq(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i % 100
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Uniq(collection)
	}
}

// BenchmarkReverse benchmarks the Reverse function
func BenchmarkReverse(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Reverse(append([]int{}, collection...))
	}
}
