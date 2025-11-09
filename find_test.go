package main

import (
	"strings"
	"testing"
	"time"
)

// TestIndexOf tests the IndexOf function with various element types and scenarios
func TestIndexOf(t *testing.T) {
	tests := []struct {
		name       string
		collection interface{}
		element    interface{}
		expected   int
	}{
		{
			name:       "element found at beginning",
			collection: []int{1, 2, 3, 4, 5},
			element:    1,
			expected:   0,
		},
		{
			name:       "element found in middle",
			collection: []int{1, 2, 3, 4, 5},
			element:    3,
			expected:   2,
		},
		{
			name:       "element found at end",
			collection: []int{1, 2, 3, 4, 5},
			element:    5,
			expected:   4,
		},
		{
			name:       "element not found",
			collection: []int{1, 2, 3, 4, 5},
			element:    10,
			expected:   -1,
		},
		{
			name:       "empty collection",
			collection: []int{},
			element:    1,
			expected:   -1,
		},
		{
			name:       "string collection",
			collection: []string{"apple", "banana", "cherry"},
			element:    "banana",
			expected:   1,
		},
		{
			name:       "duplicate elements returns first",
			collection: []int{1, 2, 3, 2, 5},
			element:    2,
			expected:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result int
			switch col := tt.collection.(type) {
			case []int:
				result = IndexOf(col, tt.element.(int))
			case []string:
				result = IndexOf(col, tt.element.(string))
			}

			if result != tt.expected {
				t.Errorf("IndexOf() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestLastIndexOf tests the LastIndexOf function with various scenarios
func TestLastIndexOf(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		element    int
		expected   int
	}{
		{
			name:       "element found at end",
			collection: []int{1, 2, 3, 4, 5},
			element:    5,
			expected:   4,
		},
		{
			name:       "element found in middle",
			collection: []int{1, 2, 3, 4, 5},
			element:    3,
			expected:   2,
		},
		{
			name:       "element not found",
			collection: []int{1, 2, 3, 4, 5},
			element:    10,
			expected:   -1,
		},
		{
			name:       "empty collection",
			collection: []int{},
			element:    1,
			expected:   -1,
		},
		{
			name:       "duplicate elements returns last",
			collection: []int{1, 2, 3, 2, 5},
			element:    2,
			expected:   3,
		},
		{
			name:       "all same elements",
			collection: []int{5, 5, 5, 5},
			element:    5,
			expected:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastIndexOf(tt.collection, tt.element)
			if result != tt.expected {
				t.Errorf("LastIndexOf() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestFind tests the Find function with various predicate functions
func TestFind(t *testing.T) {
	tests := []struct {
		name          string
		collection    []int
		predicate     func(int) bool
		expectedValue int
		expectedFound bool
	}{
		{
			name:          "find even number",
			collection:    []int{1, 3, 4, 5, 7},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedValue: 4,
			expectedFound: true,
		},
		{
			name:          "find greater than 10",
			collection:    []int{1, 3, 15, 5, 7},
			predicate:     func(n int) bool { return n > 10 },
			expectedValue: 15,
			expectedFound: true,
		},
		{
			name:          "not found",
			collection:    []int{1, 3, 5, 7},
			predicate:     func(n int) bool { return n > 100 },
			expectedValue: 0,
			expectedFound: false,
		},
		{
			name:          "empty collection",
			collection:    []int{},
			predicate:     func(n int) bool { return n > 0 },
			expectedValue: 0,
			expectedFound: false,
		},
		{
			name:          "find first element",
			collection:    []int{10, 20, 30},
			predicate:     func(n int) bool { return n >= 10 },
			expectedValue: 10,
			expectedFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := Find(tt.collection, tt.predicate)
			if value != tt.expectedValue || found != tt.expectedFound {
				t.Errorf("Find() = (%d, %v); expected (%d, %v)",
					value, found, tt.expectedValue, tt.expectedFound)
			}
		})
	}
}

// TestFindIndexOf tests the FindIndexOf function
func TestFindIndexOf(t *testing.T) {
	tests := []struct {
		name          string
		collection    []int
		predicate     func(int) bool
		expectedValue int
		expectedIndex int
		expectedFound bool
	}{
		{
			name:          "find with index",
			collection:    []int{1, 3, 4, 5, 7},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedValue: 4,
			expectedIndex: 2,
			expectedFound: true,
		},
		{
			name:          "not found",
			collection:    []int{1, 3, 5, 7},
			predicate:     func(n int) bool { return n > 100 },
			expectedValue: 0,
			expectedIndex: -1,
			expectedFound: false,
		},
		{
			name:          "empty collection",
			collection:    []int{},
			predicate:     func(n int) bool { return n > 0 },
			expectedValue: 0,
			expectedIndex: -1,
			expectedFound: false,
		},
		{
			name:          "find at beginning",
			collection:    []int{10, 20, 30},
			predicate:     func(n int) bool { return n == 10 },
			expectedValue: 10,
			expectedIndex: 0,
			expectedFound: true,
		},
		{
			name:          "find at end",
			collection:    []int{10, 20, 30},
			predicate:     func(n int) bool { return n == 30 },
			expectedValue: 30,
			expectedIndex: 2,
			expectedFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, index, found := FindIndexOf(tt.collection, tt.predicate)
			if value != tt.expectedValue || index != tt.expectedIndex || found != tt.expectedFound {
				t.Errorf("FindIndexOf() = (%d, %d, %v); expected (%d, %d, %v)",
					value, index, found, tt.expectedValue, tt.expectedIndex, tt.expectedFound)
			}
		})
	}
}

// TestFindLastIndexOf tests the FindLastIndexOf function
func TestFindLastIndexOf(t *testing.T) {
	tests := []struct {
		name          string
		collection    []int
		predicate     func(int) bool
		expectedValue int
		expectedIndex int
		expectedFound bool
	}{
		{
			name:          "find last even number",
			collection:    []int{1, 4, 3, 6, 7},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedValue: 6,
			expectedIndex: 3,
			expectedFound: true,
		},
		{
			name:          "not found",
			collection:    []int{1, 3, 5, 7},
			predicate:     func(n int) bool { return n > 100 },
			expectedValue: 0,
			expectedIndex: -1,
			expectedFound: false,
		},
		{
			name:          "empty collection",
			collection:    []int{},
			predicate:     func(n int) bool { return n > 0 },
			expectedValue: 0,
			expectedIndex: -1,
			expectedFound: false,
		},
		{
			name:          "find last element",
			collection:    []int{10, 20, 30},
			predicate:     func(n int) bool { return n > 15 },
			expectedValue: 30,
			expectedIndex: 2,
			expectedFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, index, found := FindLastIndexOf(tt.collection, tt.predicate)
			if value != tt.expectedValue || index != tt.expectedIndex || found != tt.expectedFound {
				t.Errorf("FindLastIndexOf() = (%d, %d, %v); expected (%d, %d, %v)",
					value, index, found, tt.expectedValue, tt.expectedIndex, tt.expectedFound)
			}
		})
	}
}

// TestFindOrElse tests the FindOrElse function
func TestFindOrElse(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		fallback   int
		predicate  func(int) bool
		expected   int
	}{
		{
			name:       "find existing element",
			collection: []int{1, 3, 4, 5, 7},
			fallback:   999,
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   4,
		},
		{
			name:       "return fallback when not found",
			collection: []int{1, 3, 5, 7},
			fallback:   999,
			predicate:  func(n int) bool { return n > 100 },
			expected:   999,
		},
		{
			name:       "empty collection returns fallback",
			collection: []int{},
			fallback:   -1,
			predicate:  func(n int) bool { return n > 0 },
			expected:   -1,
		},
		{
			name:       "find first matching element",
			collection: []int{10, 20, 30},
			fallback:   0,
			predicate:  func(n int) bool { return n >= 20 },
			expected:   20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindOrElse(tt.collection, tt.fallback, tt.predicate)
			if result != tt.expected {
				t.Errorf("FindOrElse() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestFindKey tests the FindKey function for maps
func TestFindKey(t *testing.T) {
	tests := []struct {
		name          string
		object        map[string]int
		value         int
		expectedKey   string
		expectedFound bool
	}{
		{
			name:          "find existing value",
			object:        map[string]int{"a": 1, "b": 2, "c": 3},
			value:         2,
			expectedKey:   "b",
			expectedFound: true,
		},
		{
			name:          "value not found",
			object:        map[string]int{"a": 1, "b": 2, "c": 3},
			value:         10,
			expectedKey:   "",
			expectedFound: false,
		},
		{
			name:          "empty map",
			object:        map[string]int{},
			value:         1,
			expectedKey:   "",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, found := FindKey(tt.object, tt.value)
			if key != tt.expectedKey || found != tt.expectedFound {
				t.Errorf("FindKey() = (%s, %v); expected (%s, %v)",
					key, found, tt.expectedKey, tt.expectedFound)
			}
		})
	}
}

// TestFindKeyBy tests the FindKeyBy function with custom predicates
func TestFindKeyBy(t *testing.T) {
	tests := []struct {
		name          string
		object        map[string]int
		predicate     func(string, int) bool
		expectedKey   string
		expectedFound bool
	}{
		{
			name:   "find key by value condition",
			object: map[string]int{"a": 1, "b": 2, "c": 3},
			predicate: func(k string, v int) bool {
				return v > 1
			},
			expectedKey:   "b",
			expectedFound: true,
		},
		{
			name:   "find key by key condition",
			object: map[string]int{"apple": 1, "banana": 2, "cherry": 3},
			predicate: func(k string, v int) bool {
				return strings.HasPrefix(k, "b")
			},
			expectedKey:   "banana",
			expectedFound: true,
		},
		{
			name:   "not found",
			object: map[string]int{"a": 1, "b": 2, "c": 3},
			predicate: func(k string, v int) bool {
				return v > 100
			},
			expectedKey:   "",
			expectedFound: false,
		},
		{
			name:   "empty map",
			object: map[string]int{},
			predicate: func(k string, v int) bool {
				return true
			},
			expectedKey:   "",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, found := FindKeyBy(tt.object, tt.predicate)
			if found != tt.expectedFound {
				t.Errorf("FindKeyBy() found = %v; expected %v", found, tt.expectedFound)
			}
			if found && tt.expectedFound {
				// For maps, iteration order is not guaranteed, so we just check the predicate matches
				if !tt.predicate(key, tt.object[key]) {
					t.Errorf("FindKeyBy() returned key %s that doesn't match predicate", key)
				}
			}
		})
	}
}

// TestFindUniques tests the FindUniques function
func TestFindUniques(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   []int
	}{
		{
			name:       "some duplicates",
			collection: []int{1, 2, 3, 2, 4, 5, 5},
			expected:   []int{1, 3, 4},
		},
		{
			name:       "no duplicates",
			collection: []int{1, 2, 3, 4, 5},
			expected:   []int{1, 2, 3, 4, 5},
		},
		{
			name:       "all duplicates",
			collection: []int{1, 1, 2, 2, 3, 3},
			expected:   []int{},
		},
		{
			name:       "empty collection",
			collection: []int{},
			expected:   []int{},
		},
		{
			name:       "single element",
			collection: []int{1},
			expected:   []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindUniques(tt.collection)
			if len(result) != len(tt.expected) {
				t.Errorf("FindUniques() length = %d; expected %d", len(result), len(tt.expected))
				return
			}
			// Check all expected elements are in result (order preserved)
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("FindUniques() at index %d = %d; expected %d", i, result[i], v)
				}
			}
		})
	}
}

// TestFindUniquesBy tests the FindUniquesBy function with custom iteratee
func TestFindUniquesBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name       string
		collection []Person
		iteratee   func(Person) int
		expected   []Person
	}{
		{
			name: "unique by age",
			collection: []Person{
				{"Alice", 25},
				{"Bob", 30},
				{"Charlie", 25},
				{"David", 35},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: []Person{{"Bob", 30}, {"David", 35}},
		},
		{
			name: "all unique ages",
			collection: []Person{
				{"Alice", 25},
				{"Bob", 30},
				{"Charlie", 35},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: []Person{{"Alice", 25}, {"Bob", 30}, {"Charlie", 35}},
		},
		{
			name:       "empty collection",
			collection: []Person{},
			iteratee:   func(p Person) int { return p.Age },
			expected:   []Person{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindUniquesBy(tt.collection, tt.iteratee)
			if len(result) != len(tt.expected) {
				t.Errorf("FindUniquesBy() length = %d; expected %d", len(result), len(tt.expected))
			}
		})
	}
}

// TestFindDuplicates tests the FindDuplicates function
func TestFindDuplicates(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   []int
	}{
		{
			name:       "some duplicates",
			collection: []int{1, 2, 3, 2, 4, 5, 5},
			expected:   []int{2, 5},
		},
		{
			name:       "no duplicates",
			collection: []int{1, 2, 3, 4, 5},
			expected:   []int{},
		},
		{
			name:       "all duplicates",
			collection: []int{1, 1, 2, 2, 3, 3},
			expected:   []int{1, 2, 3},
		},
		{
			name:       "empty collection",
			collection: []int{},
			expected:   []int{},
		},
		{
			name:       "single element",
			collection: []int{1},
			expected:   []int{},
		},
		{
			name:       "multiple occurrences",
			collection: []int{1, 1, 1, 2, 2, 3},
			expected:   []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindDuplicates(tt.collection)
			if len(result) != len(tt.expected) {
				t.Errorf("FindDuplicates() length = %d; expected %d", len(result), len(tt.expected))
				return
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("FindDuplicates() at index %d = %d; expected %d", i, result[i], v)
				}
			}
		})
	}
}

// TestFindDuplicatesBy tests the FindDuplicatesBy function with custom iteratee
func TestFindDuplicatesBy(t *testing.T) {
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
			name: "duplicates by age",
			collection: []Person{
				{"Alice", 25},
				{"Bob", 30},
				{"Charlie", 25},
				{"David", 35},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: 1, // Only one age (25) is duplicated
		},
		{
			name: "no duplicates",
			collection: []Person{
				{"Alice", 25},
				{"Bob", 30},
				{"Charlie", 35},
			},
			iteratee: func(p Person) int { return p.Age },
			expected: 0,
		},
		{
			name:       "empty collection",
			collection: []Person{},
			iteratee:   func(p Person) int { return p.Age },
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindDuplicatesBy(tt.collection, tt.iteratee)
			if len(result) != tt.expected {
				t.Errorf("FindDuplicatesBy() length = %d; expected %d", len(result), tt.expected)
			}
		})
	}
}

// TestMin tests the Min function with various types
func TestMin(t *testing.T) {
	tests := []struct {
		name       string
		collection interface{}
		expected   interface{}
	}{
		{
			name:       "integers",
			collection: []int{5, 2, 8, 1, 9},
			expected:   1,
		},
		{
			name:       "floats",
			collection: []float64{5.5, 2.2, 8.8, 1.1, 9.9},
			expected:   1.1,
		},
		{
			name:       "strings",
			collection: []string{"zebra", "apple", "banana"},
			expected:   "apple",
		},
		{
			name:       "single element",
			collection: []int{42},
			expected:   42,
		},
		{
			name:       "negative numbers",
			collection: []int{-5, -2, -8, -1, -9},
			expected:   -9,
		},
		{
			name:       "empty collection",
			collection: []int{},
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch col := tt.collection.(type) {
			case []int:
				result := Min(col)
				if result != tt.expected.(int) {
					t.Errorf("Min() = %d; expected %d", result, tt.expected.(int))
				}
			case []float64:
				result := Min(col)
				if result != tt.expected.(float64) {
					t.Errorf("Min() = %f; expected %f", result, tt.expected.(float64))
				}
			case []string:
				result := Min(col)
				if result != tt.expected.(string) {
					t.Errorf("Min() = %s; expected %s", result, tt.expected.(string))
				}
			}
		})
	}
}

// TestMinBy tests the MinBy function with custom comparison
func TestMinBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name       string
		collection []Person
		comparison func(Person, Person) bool
		expected   Person
	}{
		{
			name: "minimum by age",
			collection: []Person{
				{"Alice", 30},
				{"Bob", 25},
				{"Charlie", 35},
			},
			comparison: func(a, b Person) bool { return a.Age < b.Age },
			expected:   Person{"Bob", 25},
		},
		{
			name: "minimum by name",
			collection: []Person{
				{"Zebra", 30},
				{"Apple", 25},
				{"Banana", 35},
			},
			comparison: func(a, b Person) bool { return a.Name < b.Name },
			expected:   Person{"Apple", 25},
		},
		{
			name:       "empty collection",
			collection: []Person{},
			comparison: func(a, b Person) bool { return a.Age < b.Age },
			expected:   Person{"", 0},
		},
		{
			name: "single element",
			collection: []Person{
				{"Alice", 30},
			},
			comparison: func(a, b Person) bool { return a.Age < b.Age },
			expected:   Person{"Alice", 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinBy(tt.collection, tt.comparison)
			if result != tt.expected {
				t.Errorf("MinBy() = %+v; expected %+v", result, tt.expected)
			}
		})
	}
}

// TestMax tests the Max function with various types
func TestMax(t *testing.T) {
	tests := []struct {
		name       string
		collection interface{}
		expected   interface{}
	}{
		{
			name:       "integers",
			collection: []int{5, 2, 8, 1, 9},
			expected:   9,
		},
		{
			name:       "floats",
			collection: []float64{5.5, 2.2, 8.8, 1.1, 9.9},
			expected:   9.9,
		},
		{
			name:       "strings",
			collection: []string{"zebra", "apple", "banana"},
			expected:   "zebra",
		},
		{
			name:       "single element",
			collection: []int{42},
			expected:   42,
		},
		{
			name:       "negative numbers",
			collection: []int{-5, -2, -8, -1, -9},
			expected:   -1,
		},
		{
			name:       "empty collection",
			collection: []int{},
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch col := tt.collection.(type) {
			case []int:
				result := Max(col)
				if result != tt.expected.(int) {
					t.Errorf("Max() = %d; expected %d", result, tt.expected.(int))
				}
			case []float64:
				result := Max(col)
				if result != tt.expected.(float64) {
					t.Errorf("Max() = %f; expected %f", result, tt.expected.(float64))
				}
			case []string:
				result := Max(col)
				if result != tt.expected.(string) {
					t.Errorf("Max() = %s; expected %s", result, tt.expected.(string))
				}
			}
		})
	}
}

// TestMaxBy tests the MaxBy function with custom comparison
func TestMaxBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name       string
		collection []Person
		comparison func(Person, Person) bool
		expected   Person
	}{
		{
			name: "maximum by age",
			collection: []Person{
				{"Alice", 30},
				{"Bob", 25},
				{"Charlie", 35},
			},
			comparison: func(a, b Person) bool { return a.Age > b.Age },
			expected:   Person{"Charlie", 35},
		},
		{
			name: "maximum by name",
			collection: []Person{
				{"Zebra", 30},
				{"Apple", 25},
				{"Banana", 35},
			},
			comparison: func(a, b Person) bool { return a.Name > b.Name },
			expected:   Person{"Zebra", 30},
		},
		{
			name:       "empty collection",
			collection: []Person{},
			comparison: func(a, b Person) bool { return a.Age > b.Age },
			expected:   Person{"", 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxBy(tt.collection, tt.comparison)
			if result != tt.expected {
				t.Errorf("MaxBy() = %+v; expected %+v", result, tt.expected)
			}
		})
	}
}

// TestEarliest tests the Earliest function with time.Time values
func TestEarliest(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	lastWeek := now.AddDate(0, 0, -7)

	tests := []struct {
		name     string
		times    []time.Time
		expected time.Time
	}{
		{
			name:     "multiple times",
			times:    []time.Time{now, yesterday, tomorrow, lastWeek},
			expected: lastWeek,
		},
		{
			name:     "single time",
			times:    []time.Time{now},
			expected: now,
		},
		{
			name:     "empty",
			times:    []time.Time{},
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Earliest(tt.times...)
			if !result.Equal(tt.expected) {
				t.Errorf("Earliest() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestEarliestBy tests the EarliestBy function with custom iteratee
func TestEarliestBy(t *testing.T) {
	type Event struct {
		Name string
		Time time.Time
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name       string
		collection []Event
		iteratee   func(Event) time.Time
		expected   Event
	}{
		{
			name: "earliest event",
			collection: []Event{
				{"Today", now},
				{"Yesterday", yesterday},
				{"Tomorrow", tomorrow},
			},
			iteratee: func(e Event) time.Time { return e.Time },
			expected: Event{"Yesterday", yesterday},
		},
		{
			name:       "empty collection",
			collection: []Event{},
			iteratee:   func(e Event) time.Time { return e.Time },
			expected:   Event{"", time.Time{}},
		},
		{
			name: "single event",
			collection: []Event{
				{"Today", now},
			},
			iteratee: func(e Event) time.Time { return e.Time },
			expected: Event{"Today", now},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EarliestBy(tt.collection, tt.iteratee)
			if result.Name != tt.expected.Name {
				t.Errorf("EarliestBy() = %+v; expected %+v", result, tt.expected)
			}
		})
	}
}

// TestLatest tests the Latest function with time.Time values
func TestLatest(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	tests := []struct {
		name     string
		times    []time.Time
		expected time.Time
	}{
		{
			name:     "multiple times",
			times:    []time.Time{now, yesterday, tomorrow, nextWeek},
			expected: nextWeek,
		},
		{
			name:     "single time",
			times:    []time.Time{now},
			expected: now,
		},
		{
			name:     "empty",
			times:    []time.Time{},
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Latest(tt.times...)
			if !result.Equal(tt.expected) {
				t.Errorf("Latest() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestLatestBy tests the LatestBy function with custom iteratee
func TestLatestBy(t *testing.T) {
	type Event struct {
		Name string
		Time time.Time
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name       string
		collection []Event
		iteratee   func(Event) time.Time
		expected   Event
	}{
		{
			name: "latest event",
			collection: []Event{
				{"Today", now},
				{"Yesterday", yesterday},
				{"Tomorrow", tomorrow},
			},
			iteratee: func(e Event) time.Time { return e.Time },
			expected: Event{"Tomorrow", tomorrow},
		},
		{
			name:       "empty collection",
			collection: []Event{},
			iteratee:   func(e Event) time.Time { return e.Time },
			expected:   Event{"", time.Time{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LatestBy(tt.collection, tt.iteratee)
			if result.Name != tt.expected.Name {
				t.Errorf("LatestBy() = %+v; expected %+v", result, tt.expected)
			}
		})
	}
}

// TestFirst tests the First function
func TestFirst(t *testing.T) {
	tests := []struct {
		name          string
		collection    []int
		expectedValue int
		expectedFound bool
	}{
		{
			name:          "non-empty collection",
			collection:    []int{1, 2, 3, 4, 5},
			expectedValue: 1,
			expectedFound: true,
		},
		{
			name:          "single element",
			collection:    []int{42},
			expectedValue: 42,
			expectedFound: true,
		},
		{
			name:          "empty collection",
			collection:    []int{},
			expectedValue: 0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := First(tt.collection)
			if value != tt.expectedValue || found != tt.expectedFound {
				t.Errorf("First() = (%d, %v); expected (%d, %v)",
					value, found, tt.expectedValue, tt.expectedFound)
			}
		})
	}
}

// TestFirstOrEmpty tests the FirstOrEmpty function
func TestFirstOrEmpty(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   int
	}{
		{
			name:       "non-empty collection",
			collection: []int{1, 2, 3, 4, 5},
			expected:   1,
		},
		{
			name:       "empty collection returns zero value",
			collection: []int{},
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FirstOrEmpty(tt.collection)
			if result != tt.expected {
				t.Errorf("FirstOrEmpty() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestFirstOr tests the FirstOr function
func TestFirstOr(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		fallback   int
		expected   int
	}{
		{
			name:       "non-empty collection",
			collection: []int{1, 2, 3, 4, 5},
			fallback:   999,
			expected:   1,
		},
		{
			name:       "empty collection returns fallback",
			collection: []int{},
			fallback:   999,
			expected:   999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FirstOr(tt.collection, tt.fallback)
			if result != tt.expected {
				t.Errorf("FirstOr() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestLast tests the Last function
func TestLast(t *testing.T) {
	tests := []struct {
		name          string
		collection    []int
		expectedValue int
		expectedFound bool
	}{
		{
			name:          "non-empty collection",
			collection:    []int{1, 2, 3, 4, 5},
			expectedValue: 5,
			expectedFound: true,
		},
		{
			name:          "single element",
			collection:    []int{42},
			expectedValue: 42,
			expectedFound: true,
		},
		{
			name:          "empty collection",
			collection:    []int{},
			expectedValue: 0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := Last(tt.collection)
			if value != tt.expectedValue || found != tt.expectedFound {
				t.Errorf("Last() = (%d, %v); expected (%d, %v)",
					value, found, tt.expectedValue, tt.expectedFound)
			}
		})
	}
}

// TestLastOrEmpty tests the LastOrEmpty function
func TestLastOrEmpty(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expected   int
	}{
		{
			name:       "non-empty collection",
			collection: []int{1, 2, 3, 4, 5},
			expected:   5,
		},
		{
			name:       "empty collection returns zero value",
			collection: []int{},
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastOrEmpty(tt.collection)
			if result != tt.expected {
				t.Errorf("LastOrEmpty() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestLastOr tests the LastOr function
func TestLastOr(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		fallback   int
		expected   int
	}{
		{
			name:       "non-empty collection",
			collection: []int{1, 2, 3, 4, 5},
			fallback:   999,
			expected:   5,
		},
		{
			name:       "empty collection returns fallback",
			collection: []int{},
			fallback:   999,
			expected:   999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastOr(tt.collection, tt.fallback)
			if result != tt.expected {
				t.Errorf("LastOr() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestNth tests the Nth function with positive and negative indices
func TestNth(t *testing.T) {
	tests := []struct {
		name        string
		collection  []int
		nth         int
		expected    int
		expectError bool
	}{
		{
			name:        "positive index",
			collection:  []int{1, 2, 3, 4, 5},
			nth:         2,
			expected:    3,
			expectError: false,
		},
		{
			name:        "negative index",
			collection:  []int{1, 2, 3, 4, 5},
			nth:         -1,
			expected:    5,
			expectError: false,
		},
		{
			name:        "negative index from beginning",
			collection:  []int{1, 2, 3, 4, 5},
			nth:         -5,
			expected:    1,
			expectError: false,
		},
		{
			name:        "first element",
			collection:  []int{1, 2, 3, 4, 5},
			nth:         0,
			expected:    1,
			expectError: false,
		},
		{
			name:        "out of bounds positive",
			collection:  []int{1, 2, 3, 4, 5},
			nth:         10,
			expected:    0,
			expectError: true,
		},
		{
			name:        "out of bounds negative",
			collection:  []int{1, 2, 3, 4, 5},
			nth:         -10,
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty collection",
			collection:  []int{},
			nth:         0,
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Nth(tt.collection, tt.nth)
			if tt.expectError {
				if err == nil {
					t.Errorf("Nth() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Nth() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Nth() = %d; expected %d", result, tt.expected)
				}
			}
		})
	}
}

// TestSample tests the Sample function
func TestSample(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		expectZero bool
	}{
		{
			name:       "non-empty collection",
			collection: []int{1, 2, 3, 4, 5},
			expectZero: false,
		},
		{
			name:       "single element",
			collection: []int{42},
			expectZero: false,
		},
		{
			name:       "empty collection",
			collection: []int{},
			expectZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sample(tt.collection)
			if tt.expectZero {
				if result != 0 {
					t.Errorf("Sample() expected zero value for empty collection, got %d", result)
				}
			} else {
				// Check if result is in collection
				found := false
				for _, v := range tt.collection {
					if v == result {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Sample() = %d; not found in collection %v", result, tt.collection)
				}
			}
		})
	}
}

// TestSamples tests the Samples function
func TestSamples(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		count      int
		expected   int
	}{
		{
			name:       "sample less than collection size",
			collection: []int{1, 2, 3, 4, 5},
			count:      3,
			expected:   3,
		},
		{
			name:       "sample equal to collection size",
			collection: []int{1, 2, 3, 4, 5},
			count:      5,
			expected:   5,
		},
		{
			name:       "sample more than collection size",
			collection: []int{1, 2, 3, 4, 5},
			count:      10,
			expected:   5,
		},
		{
			name:       "empty collection",
			collection: []int{},
			count:      3,
			expected:   0,
		},
		{
			name:       "zero count",
			collection: []int{1, 2, 3, 4, 5},
			count:      0,
			expected:   0,
		},
		{
			name:       "single element",
			collection: []int{42},
			count:      1,
			expected:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Samples(tt.collection, tt.count)
			if len(result) != tt.expected {
				t.Errorf("Samples() length = %d; expected %d", len(result), tt.expected)
			}

			// Check all results are from original collection
			for _, v := range result {
				found := false
				for _, orig := range tt.collection {
					if v == orig {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Samples() contains %d which is not in original collection", v)
				}
			}

			// Check uniqueness
			seen := make(map[int]bool)
			for _, v := range result {
				if seen[v] {
					t.Errorf("Samples() contains duplicate value %d", v)
				}
				seen[v] = true
			}
		})
	}
}

// TestSamplesUniqueness tests that Samples returns unique values
func TestSamplesUniqueness(t *testing.T) {
	t.Run("samples are unique", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		count := 5

		// Run multiple times to ensure randomness and uniqueness
		for i := 0; i < 10; i++ {
			result := Samples(collection, count)

			// Check length
			if len(result) != count {
				t.Errorf("Samples() length = %d; expected %d", len(result), count)
			}

			// Check uniqueness
			seen := make(map[int]bool)
			for _, v := range result {
				if seen[v] {
					t.Errorf("Samples() contains duplicate value %d in iteration %d", v, i)
				}
				seen[v] = true
			}
		}
	})
}

// BenchmarkIndexOf benchmarks the IndexOf function
func BenchmarkIndexOf(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IndexOf(collection, 500)
	}
}

// BenchmarkFind benchmarks the Find function
func BenchmarkFind(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}

	predicate := func(n int) bool { return n == 500 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Find(collection, predicate)
	}
}

// BenchmarkMin benchmarks the Min function
func BenchmarkMin(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Min(collection)
	}
}

// BenchmarkMax benchmarks the Max function
func BenchmarkMax(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Max(collection)
	}
}

// BenchmarkFindDuplicates benchmarks the FindDuplicates function
func BenchmarkFindDuplicates(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i % 100 // Create duplicates
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindDuplicates(collection)
	}
}

