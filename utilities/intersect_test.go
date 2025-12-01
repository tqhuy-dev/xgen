package utilities

import (
	"reflect"
	"strings"
	"testing"
)

// TestContains tests the Contains function with various element types
func TestContains(t *testing.T) {
	tests := []struct {
		name       string
		collection interface{}
		element    interface{}
		expected   bool
	}{
		{
			name:       "element found in integers",
			collection: []int{1, 2, 3, 4, 5},
			element:    3,
			expected:   true,
		},
		{
			name:       "element not found in integers",
			collection: []int{1, 2, 3, 4, 5},
			element:    10,
			expected:   false,
		},
		{
			name:       "element found in strings",
			collection: []string{"apple", "banana", "cherry"},
			element:    "banana",
			expected:   true,
		},
		{
			name:       "element not found in strings",
			collection: []string{"apple", "banana", "cherry"},
			element:    "orange",
			expected:   false,
		},
		{
			name:       "empty collection",
			collection: []int{},
			element:    1,
			expected:   false,
		},
		{
			name:       "element at beginning",
			collection: []int{1, 2, 3, 4, 5},
			element:    1,
			expected:   true,
		},
		{
			name:       "element at end",
			collection: []int{1, 2, 3, 4, 5},
			element:    5,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch col := tt.collection.(type) {
			case []int:
				result = Contains(col, tt.element.(int))
			case []string:
				result = Contains(col, tt.element.(string))
			}

			if result != tt.expected {
				t.Errorf("Contains() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestContainsBy tests the ContainsBy function with various predicates
func TestContainsBy(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int) bool
		expected   bool
	}{
		{
			name:       "find even number",
			collection: []int{1, 3, 4, 5, 7},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
		{
			name:       "no even number",
			collection: []int{1, 3, 5, 7, 9},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   false,
		},
		{
			name:       "find greater than 10",
			collection: []int{1, 5, 15, 7},
			predicate:  func(n int) bool { return n > 10 },
			expected:   true,
		},
		{
			name:       "empty collection",
			collection: []int{},
			predicate:  func(n int) bool { return n > 0 },
			expected:   false,
		},
		{
			name:       "all match predicate",
			collection: []int{2, 4, 6, 8},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsBy(tt.collection, tt.predicate)
			if result != tt.expected {
				t.Errorf("ContainsBy() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestEvery tests the Every function
func TestEvery(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		subset     []int
		expected   bool
	}{
		{
			name:       "all elements present",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{2, 3, 4},
			expected:   true,
		},
		{
			name:       "some elements missing",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{2, 3, 10},
			expected:   false,
		},
		{
			name:       "empty subset",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{},
			expected:   true,
		},
		{
			name:       "empty collection non-empty subset",
			collection: []int{},
			subset:     []int{1, 2},
			expected:   false,
		},
		{
			name:       "both empty",
			collection: []int{},
			subset:     []int{},
			expected:   true,
		},
		{
			name:       "single element present",
			collection: []int{1, 2, 3},
			subset:     []int{2},
			expected:   true,
		},
		{
			name:       "subset equals collection",
			collection: []int{1, 2, 3},
			subset:     []int{1, 2, 3},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Every(tt.collection, tt.subset)
			if result != tt.expected {
				t.Errorf("Every() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestEveryBy tests the EveryBy function
func TestEveryBy(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int) bool
		expected   bool
	}{
		{
			name:       "all even numbers",
			collection: []int{2, 4, 6, 8},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
		{
			name:       "mixed even and odd",
			collection: []int{2, 4, 5, 8},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   false,
		},
		{
			name:       "empty collection",
			collection: []int{},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
		{
			name:       "all positive",
			collection: []int{1, 2, 3, 4},
			predicate:  func(n int) bool { return n > 0 },
			expected:   true,
		},
		{
			name:       "one negative",
			collection: []int{1, 2, -3, 4},
			predicate:  func(n int) bool { return n > 0 },
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EveryBy(tt.collection, tt.predicate)
			if result != tt.expected {
				t.Errorf("EveryBy() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestSome tests the Some function
func TestSome(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		subset     []int
		expected   bool
	}{
		{
			name:       "at least one element present",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{3, 10, 20},
			expected:   true,
		},
		{
			name:       "no elements present",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{10, 20, 30},
			expected:   false,
		},
		{
			name:       "empty subset",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{},
			expected:   false,
		},
		{
			name:       "empty collection",
			collection: []int{},
			subset:     []int{1, 2},
			expected:   false,
		},
		{
			name:       "both empty",
			collection: []int{},
			subset:     []int{},
			expected:   false,
		},
		{
			name:       "all elements present",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{2, 3, 4},
			expected:   true,
		},
		{
			name:       "first element matches",
			collection: []int{1, 2, 3},
			subset:     []int{1, 10, 20},
			expected:   true,
		},
		{
			name:       "last element matches",
			collection: []int{1, 2, 3},
			subset:     []int{10, 20, 3},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Some(tt.collection, tt.subset)
			if result != tt.expected {
				t.Errorf("Some() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestSomeBy tests the SomeBy function
func TestSomeBy(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int) bool
		expected   bool
	}{
		{
			name:       "at least one even",
			collection: []int{1, 3, 4, 5, 7},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
		{
			name:       "no even numbers",
			collection: []int{1, 3, 5, 7},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   false,
		},
		{
			name:       "empty collection",
			collection: []int{},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   false,
		},
		{
			name:       "all match",
			collection: []int{2, 4, 6},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
		{
			name:       "first element matches",
			collection: []int{2, 3, 5, 7},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SomeBy(tt.collection, tt.predicate)
			if result != tt.expected {
				t.Errorf("SomeBy() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestNone tests the None function
func TestNone(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		subset     []int
		expected   bool
	}{
		{
			name:       "no elements present",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{10, 20, 30},
			expected:   true,
		},
		{
			name:       "at least one element present",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{3, 10, 20},
			expected:   false,
		},
		{
			name:       "empty subset",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{},
			expected:   true,
		},
		{
			name:       "empty collection non-empty subset",
			collection: []int{},
			subset:     []int{1, 2},
			expected:   true,
		},
		{
			name:       "both empty",
			collection: []int{},
			subset:     []int{},
			expected:   true,
		},
		{
			name:       "all elements present",
			collection: []int{1, 2, 3, 4, 5},
			subset:     []int{2, 3, 4},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := None(tt.collection, tt.subset)
			if result != tt.expected {
				t.Errorf("None() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestNoneBy tests the NoneBy function
func TestNoneBy(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		predicate  func(int) bool
		expected   bool
	}{
		{
			name:       "no even numbers",
			collection: []int{1, 3, 5, 7},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
		{
			name:       "at least one even",
			collection: []int{1, 3, 4, 5, 7},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   false,
		},
		{
			name:       "empty collection",
			collection: []int{},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   true,
		},
		{
			name:       "all match predicate",
			collection: []int{2, 4, 6, 8},
			predicate:  func(n int) bool { return n%2 == 0 },
			expected:   false,
		},
		{
			name:       "all negative",
			collection: []int{-1, -2, -3},
			predicate:  func(n int) bool { return n > 0 },
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NoneBy(tt.collection, tt.predicate)
			if result != tt.expected {
				t.Errorf("NoneBy() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestIntersect tests the Intersect function
func TestIntersect(t *testing.T) {
	tests := []struct {
		name     string
		list1    []int
		list2    []int
		expected []int
	}{
		{
			name:     "common elements",
			list1:    []int{1, 2, 3, 4, 5},
			list2:    []int{3, 4, 5, 6, 7},
			expected: []int{3, 4, 5},
		},
		{
			name:     "no common elements",
			list1:    []int{1, 2, 3},
			list2:    []int{4, 5, 6},
			expected: []int{},
		},
		{
			name:     "all elements common",
			list1:    []int{1, 2, 3},
			list2:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty first list",
			list1:    []int{},
			list2:    []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "empty second list",
			list1:    []int{1, 2, 3},
			list2:    []int{},
			expected: []int{},
		},
		{
			name:     "both empty",
			list1:    []int{},
			list2:    []int{},
			expected: []int{},
		},
		{
			name:     "single common element",
			list1:    []int{1, 2, 3},
			list2:    []int{3, 4, 5},
			expected: []int{3},
		},
		{
			name:     "duplicates in list2",
			list1:    []int{1, 2, 3},
			list2:    []int{2, 2, 3, 3},
			expected: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersect(tt.list1, tt.list2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Intersect() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestIntersectWithStrings tests Intersect with string slices
func TestIntersectWithStrings(t *testing.T) {
	tests := []struct {
		name     string
		list1    []string
		list2    []string
		expected []string
	}{
		{
			name:     "common strings",
			list1:    []string{"apple", "banana", "cherry"},
			list2:    []string{"banana", "cherry", "date"},
			expected: []string{"banana", "cherry"},
		},
		{
			name:     "no common strings",
			list1:    []string{"apple", "banana"},
			list2:    []string{"cherry", "date"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersect(tt.list1, tt.list2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Intersect() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestDifference tests the Difference function
func TestDifference(t *testing.T) {
	tests := []struct {
		name          string
		list1         []int
		list2         []int
		expectedLeft  []int
		expectedRight []int
	}{
		{
			name:          "some differences",
			list1:         []int{1, 2, 3, 4, 5},
			list2:         []int{3, 4, 5, 6, 7},
			expectedLeft:  []int{1, 2},
			expectedRight: []int{6, 7},
		},
		{
			name:          "no common elements",
			list1:         []int{1, 2, 3},
			list2:         []int{4, 5, 6},
			expectedLeft:  []int{1, 2, 3},
			expectedRight: []int{4, 5, 6},
		},
		{
			name:          "all elements common",
			list1:         []int{1, 2, 3},
			list2:         []int{1, 2, 3},
			expectedLeft:  []int{},
			expectedRight: []int{},
		},
		{
			name:          "empty first list",
			list1:         []int{},
			list2:         []int{1, 2, 3},
			expectedLeft:  []int{},
			expectedRight: []int{1, 2, 3},
		},
		{
			name:          "empty second list",
			list1:         []int{1, 2, 3},
			list2:         []int{},
			expectedLeft:  []int{1, 2, 3},
			expectedRight: []int{},
		},
		{
			name:          "both empty",
			list1:         []int{},
			list2:         []int{},
			expectedLeft:  []int{},
			expectedRight: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, right := Difference(tt.list1, tt.list2)
			if !reflect.DeepEqual(left, tt.expectedLeft) {
				t.Errorf("Difference() left = %v; expected %v", left, tt.expectedLeft)
			}
			if !reflect.DeepEqual(right, tt.expectedRight) {
				t.Errorf("Difference() right = %v; expected %v", right, tt.expectedRight)
			}
		})
	}
}

// TestUnion tests the Union function
func TestUnion(t *testing.T) {
	tests := []struct {
		name     string
		lists    [][]int
		expected []int
	}{
		{
			name:     "two lists with overlap",
			lists:    [][]int{{1, 2, 3}, {3, 4, 5}},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "two lists no overlap",
			lists:    [][]int{{1, 2, 3}, {4, 5, 6}},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "single list",
			lists:    [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty lists",
			lists:    [][]int{{}, {}},
			expected: []int{},
		},
		{
			name:     "three lists",
			lists:    [][]int{{1, 2}, {2, 3}, {3, 4}},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "lists with duplicates",
			lists:    [][]int{{1, 1, 2}, {2, 3, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "no lists",
			lists:    [][]int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Union(tt.lists...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Union() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestUnionOrderPreservation tests that Union preserves element order
func TestUnionOrderPreservation(t *testing.T) {
	t.Run("order is preserved", func(t *testing.T) {
		list1 := []int{5, 3, 1}
		list2 := []int{2, 4}
		expected := []int{5, 3, 1, 2, 4}

		result := Union(list1, list2)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Union() = %v; expected %v (order should be preserved)", result, expected)
		}
	})
}

// TestWithout tests the Without function
func TestWithout(t *testing.T) {
	tests := []struct {
		name       string
		collection []int
		exclude    []int
		expected   []int
	}{
		{
			name:       "exclude some elements",
			collection: []int{1, 2, 3, 4, 5},
			exclude:    []int{2, 4},
			expected:   []int{1, 3, 5},
		},
		{
			name:       "exclude all elements",
			collection: []int{1, 2, 3},
			exclude:    []int{1, 2, 3},
			expected:   []int{},
		},
		{
			name:       "exclude no elements",
			collection: []int{1, 2, 3},
			exclude:    []int{4, 5, 6},
			expected:   []int{1, 2, 3},
		},
		{
			name:       "empty collection",
			collection: []int{},
			exclude:    []int{1, 2},
			expected:   []int{},
		},
		{
			name:       "empty exclude list",
			collection: []int{1, 2, 3},
			exclude:    []int{},
			expected:   []int{1, 2, 3},
		},
		{
			name:       "both empty",
			collection: []int{},
			exclude:    []int{},
			expected:   []int{},
		},
		{
			name:       "exclude duplicates",
			collection: []int{1, 2, 3, 2, 4},
			exclude:    []int{2},
			expected:   []int{1, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Without(tt.collection, tt.exclude...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Without() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestWithoutWithStrings tests Without with string slices
func TestWithoutWithStrings(t *testing.T) {
	tests := []struct {
		name       string
		collection []string
		exclude    []string
		expected   []string
	}{
		{
			name:       "exclude some strings",
			collection: []string{"apple", "banana", "cherry", "date"},
			exclude:    []string{"banana", "date"},
			expected:   []string{"apple", "cherry"},
		},
		{
			name:       "exclude nothing",
			collection: []string{"apple", "banana"},
			exclude:    []string{"cherry"},
			expected:   []string{"apple", "banana"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Without(tt.collection, tt.exclude...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Without() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestContainsWithCustomTypes tests Contains with custom types
func TestContainsWithCustomTypes(t *testing.T) {
	type CustomInt int

	tests := []struct {
		name       string
		collection []CustomInt
		element    CustomInt
		expected   bool
	}{
		{
			name:       "custom int found",
			collection: []CustomInt{1, 2, 3, 4, 5},
			element:    CustomInt(3),
			expected:   true,
		},
		{
			name:       "custom int not found",
			collection: []CustomInt{1, 2, 3, 4, 5},
			element:    CustomInt(10),
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.collection, tt.element)
			if result != tt.expected {
				t.Errorf("Contains() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestContainsByWithStrings tests ContainsBy with string operations
func TestContainsByWithStrings(t *testing.T) {
	tests := []struct {
		name       string
		collection []string
		predicate  func(string) bool
		expected   bool
	}{
		{
			name:       "contains string starting with prefix",
			collection: []string{"apple", "banana", "cherry"},
			predicate:  func(s string) bool { return strings.HasPrefix(s, "ban") },
			expected:   true,
		},
		{
			name:       "no string with prefix",
			collection: []string{"apple", "banana", "cherry"},
			predicate:  func(s string) bool { return strings.HasPrefix(s, "xyz") },
			expected:   false,
		},
		{
			name:       "contains long string",
			collection: []string{"cat", "dog", "elephant"},
			predicate:  func(s string) bool { return len(s) > 5 },
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsBy(tt.collection, tt.predicate)
			if result != tt.expected {
				t.Errorf("ContainsBy() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestEveryByEdgeCases tests EveryBy with edge cases
func TestEveryByEdgeCases(t *testing.T) {
	t.Run("single element matching", func(t *testing.T) {
		collection := []int{5}
		predicate := func(n int) bool { return n == 5 }
		result := EveryBy(collection, predicate)
		if !result {
			t.Error("EveryBy() should return true for single matching element")
		}
	})

	t.Run("single element not matching", func(t *testing.T) {
		collection := []int{5}
		predicate := func(n int) bool { return n == 10 }
		result := EveryBy(collection, predicate)
		if result {
			t.Error("EveryBy() should return false for single non-matching element")
		}
	})
}

// TestIntersectNoDuplicatesInResult tests that Intersect doesn't produce duplicates
func TestIntersectNoDuplicatesInResult(t *testing.T) {
	t.Run("no duplicates in result", func(t *testing.T) {
		list1 := []int{1, 1, 2, 2, 3}
		list2 := []int{2, 2, 3, 3, 4}
		result := Intersect(list1, list2)

		// Check for duplicates
		seen := make(map[int]bool)
		for _, v := range result {
			if seen[v] {
				t.Errorf("Intersect() produced duplicate value %d in result", v)
			}
			seen[v] = true
		}
	})
}

// TestDifferenceSymmetry tests that Difference correctly handles asymmetric differences
func TestDifferenceSymmetry(t *testing.T) {
	t.Run("difference is correctly asymmetric", func(t *testing.T) {
		list1 := []int{1, 2, 3}
		list2 := []int{3, 4, 5}

		left, right := Difference(list1, list2)

		// left should contain elements in list1 but not in list2
		expectedLeft := []int{1, 2}
		if !reflect.DeepEqual(left, expectedLeft) {
			t.Errorf("Difference() left = %v; expected %v", left, expectedLeft)
		}

		// right should contain elements in list2 but not in list1
		expectedRight := []int{4, 5}
		if !reflect.DeepEqual(right, expectedRight) {
			t.Errorf("Difference() right = %v; expected %v", right, expectedRight)
		}
	})
}

// BenchmarkContains benchmarks the Contains function
func BenchmarkContains(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Contains(collection, 500)
	}
}

// BenchmarkIntersect benchmarks the Intersect function
func BenchmarkIntersect(b *testing.B) {
	list1 := make([]int, 1000)
	list2 := make([]int, 1000)
	for i := range list1 {
		list1[i] = i
		list2[i] = i + 500
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Intersect(list1, list2)
	}
}

// BenchmarkUnion benchmarks the Union function
func BenchmarkUnion(b *testing.B) {
	list1 := make([]int, 500)
	list2 := make([]int, 500)
	for i := range list1 {
		list1[i] = i
		list2[i] = i + 250
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Union(list1, list2)
	}
}

// BenchmarkWithout benchmarks the Without function
func BenchmarkWithout(b *testing.B) {
	collection := make([]int, 1000)
	for i := range collection {
		collection[i] = i
	}
	exclude := []int{100, 200, 300, 400, 500}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Without(collection, exclude...)
	}
}
