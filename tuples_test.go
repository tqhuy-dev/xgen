package main

import (
	"testing"
)

// TestT2 tests the T2 tuple constructor
func TestT2(t *testing.T) {
	tests := []struct {
		name      string
		val0      interface{}
		val1      interface{}
		expected0 interface{}
		expected1 interface{}
	}{
		{
			name:      "int and string",
			val0:      42,
			val1:      "hello",
			expected0: 42,
			expected1: "hello",
		},
		{
			name:      "two strings",
			val0:      "first",
			val1:      "second",
			expected0: "first",
			expected1: "second",
		},
		{
			name:      "mixed types",
			val0:      true,
			val1:      3.14,
			expected0: true,
			expected1: 3.14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := T2(tt.val0, tt.val1)
			if result.T0 != tt.expected0 {
				t.Errorf("T2().T0 = %v; expected %v", result.T0, tt.expected0)
			}
			if result.T1 != tt.expected1 {
				t.Errorf("T2().T1 = %v; expected %v", result.T1, tt.expected1)
			}
		})
	}
}

// TestT3 tests the T3 tuple constructor
func TestT3(t *testing.T) {
	tests := []struct {
		name      string
		val0      int
		val1      string
		val2      bool
		expected0 int
		expected1 string
		expected2 bool
	}{
		{
			name:      "mixed types",
			val0:      1,
			val1:      "test",
			val2:      true,
			expected0: 1,
			expected1: "test",
			expected2: true,
		},
		{
			name:      "zero values",
			val0:      0,
			val1:      "",
			val2:      false,
			expected0: 0,
			expected1: "",
			expected2: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := T3(tt.val0, tt.val1, tt.val2)
			if result.T0 != tt.expected0 || result.T1 != tt.expected1 || result.T2 != tt.expected2 {
				t.Errorf("T3() = (%v, %v, %v); expected (%v, %v, %v)",
					result.T0, result.T1, result.T2, tt.expected0, tt.expected1, tt.expected2)
			}
		})
	}
}

// TestUnpack2 tests the Unpack2 function
func TestUnpack2(t *testing.T) {
	tests := []struct {
		name      string
		tuple     Tuple2[int, string]
		expected0 int
		expected1 string
	}{
		{
			name:      "basic unpack",
			tuple:     T2(42, "hello"),
			expected0: 42,
			expected1: "hello",
		},
		{
			name:      "zero values",
			tuple:     T2(0, ""),
			expected0: 0,
			expected1: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val0, val1 := Unpack2(tt.tuple)
			if val0 != tt.expected0 {
				t.Errorf("Unpack2() val0 = %v; expected %v", val0, tt.expected0)
			}
			if val1 != tt.expected1 {
				t.Errorf("Unpack2() val1 = %v; expected %v", val1, tt.expected1)
			}
		})
	}
}

// TestUnpack3 tests the Unpack3 function
func TestUnpack3(t *testing.T) {
	tests := []struct {
		name      string
		tuple     Tuple3[int, string, bool]
		expected0 int
		expected1 string
		expected2 bool
	}{
		{
			name:      "basic unpack",
			tuple:     T3(1, "test", true),
			expected0: 1,
			expected1: "test",
			expected2: true,
		},
		{
			name:      "zero values",
			tuple:     T3(0, "", false),
			expected0: 0,
			expected1: "",
			expected2: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val0, val1, val2 := Unpack3(tt.tuple)
			if val0 != tt.expected0 || val1 != tt.expected1 || val2 != tt.expected2 {
				t.Errorf("Unpack3() = (%v, %v, %v); expected (%v, %v, %v)",
					val0, val1, val2, tt.expected0, tt.expected1, tt.expected2)
			}
		})
	}
}

// TestZip2 tests the Zip2 function
func TestZip2(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []string
		expected []Tuple2[int, string]
	}{
		{
			name:   "equal length slices",
			slice1: []int{1, 2, 3},
			slice2: []string{"a", "b", "c"},
			expected: []Tuple2[int, string]{
				{T0: 1, T1: "a"},
				{T0: 2, T1: "b"},
				{T0: 3, T1: "c"},
			},
		},
		{
			name:   "first slice longer",
			slice1: []int{1, 2, 3, 4},
			slice2: []string{"a", "b"},
			expected: []Tuple2[int, string]{
				{T0: 1, T1: "a"},
				{T0: 2, T1: "b"},
				{T0: 3, T1: ""},
				{T0: 4, T1: ""},
			},
		},
		{
			name:   "second slice longer",
			slice1: []int{1, 2},
			slice2: []string{"a", "b", "c", "d"},
			expected: []Tuple2[int, string]{
				{T0: 1, T1: "a"},
				{T0: 2, T1: "b"},
				{T0: 0, T1: "c"},
				{T0: 0, T1: "d"},
			},
		},
		{
			name:     "empty slices",
			slice1:   []int{},
			slice2:   []string{},
			expected: []Tuple2[int, string]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Zip2(tt.slice1, tt.slice2)
			if len(result) != len(tt.expected) {
				t.Errorf("Zip2() length = %d; expected %d", len(result), len(tt.expected))
				return
			}
			for i, tuple := range result {
				if tuple.T0 != tt.expected[i].T0 || tuple.T1 != tt.expected[i].T1 {
					t.Errorf("Zip2()[%d] = (%v, %v); expected (%v, %v)",
						i, tuple.T0, tuple.T1, tt.expected[i].T0, tt.expected[i].T1)
				}
			}
		})
	}
}

// TestZip3 tests the Zip3 function
func TestZip3(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []string
		slice3   []bool
		expected []Tuple3[int, string, bool]
	}{
		{
			name:   "equal length slices",
			slice1: []int{1, 2},
			slice2: []string{"a", "b"},
			slice3: []bool{true, false},
			expected: []Tuple3[int, string, bool]{
				{T0: 1, T1: "a", T2: true},
				{T0: 2, T1: "b", T2: false},
			},
		},
		{
			name:   "unequal length slices",
			slice1: []int{1, 2, 3},
			slice2: []string{"a"},
			slice3: []bool{true, false},
			expected: []Tuple3[int, string, bool]{
				{T0: 1, T1: "a", T2: true},
				{T0: 2, T1: "", T2: false},
				{T0: 3, T1: "", T2: false},
			},
		},
		{
			name:     "empty slices",
			slice1:   []int{},
			slice2:   []string{},
			slice3:   []bool{},
			expected: []Tuple3[int, string, bool]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Zip3(tt.slice1, tt.slice2, tt.slice3)
			if len(result) != len(tt.expected) {
				t.Errorf("Zip3() length = %d; expected %d", len(result), len(tt.expected))
				return
			}
			for i, tuple := range result {
				if tuple.T0 != tt.expected[i].T0 || tuple.T1 != tt.expected[i].T1 || tuple.T2 != tt.expected[i].T2 {
					t.Errorf("Zip3()[%d] = (%v, %v, %v); expected (%v, %v, %v)",
						i, tuple.T0, tuple.T1, tuple.T2, tt.expected[i].T0, tt.expected[i].T1, tt.expected[i].T2)
				}
			}
		})
	}
}

// TestZip4 tests the Zip4 function with different length slices
func TestZip4(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []int
		slice3   []int
		slice4   []int
		expected []Tuple4[int, int, int, int]
	}{
		{
			name:   "equal length slices",
			slice1: []int{1, 2},
			slice2: []int{3, 4},
			slice3: []int{5, 6},
			slice4: []int{7, 8},
			expected: []Tuple4[int, int, int, int]{
				{T0: 1, T1: 3, T2: 5, T3: 7},
				{T0: 2, T1: 4, T2: 6, T3: 8},
			},
		},
		{
			name:   "unequal length slices",
			slice1: []int{1, 2, 3},
			slice2: []int{4},
			slice3: []int{5, 6},
			slice4: []int{7, 8, 9, 10},
			expected: []Tuple4[int, int, int, int]{
				{T0: 1, T1: 4, T2: 5, T3: 7},
				{T0: 2, T1: 0, T2: 6, T3: 8},
				{T0: 3, T1: 0, T2: 0, T3: 9},
				{T0: 0, T1: 0, T2: 0, T3: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Zip4(tt.slice1, tt.slice2, tt.slice3, tt.slice4)
			if len(result) != len(tt.expected) {
				t.Errorf("Zip4() length = %d; expected %d", len(result), len(tt.expected))
				return
			}
			for i, tuple := range result {
				if tuple.T0 != tt.expected[i].T0 || tuple.T1 != tt.expected[i].T1 ||
					tuple.T2 != tt.expected[i].T2 || tuple.T3 != tt.expected[i].T3 {
					t.Errorf("Zip4()[%d] = (%v, %v, %v, %v); expected (%v, %v, %v, %v)",
						i, tuple.T0, tuple.T1, tuple.T2, tuple.T3,
						tt.expected[i].T0, tt.expected[i].T1, tt.expected[i].T2, tt.expected[i].T3)
				}
			}
		})
	}
}

// TestZipT1y2 tests the ZipT1y2 function with iteratee
func TestZipT1y2(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []int
		iteratee func(a int, b int) int
		expected []int
	}{
		{
			name:     "sum two slices",
			slice1:   []int{1, 2, 3},
			slice2:   []int{4, 5, 6},
			iteratee: func(a int, b int) int { return a + b },
			expected: []int{5, 7, 9},
		},
		{
			name:     "multiply two slices",
			slice1:   []int{2, 3, 4},
			slice2:   []int{5, 6, 7},
			iteratee: func(a int, b int) int { return a * b },
			expected: []int{10, 18, 28},
		},
		{
			name:     "unequal length slices",
			slice1:   []int{1, 2, 3, 4},
			slice2:   []int{5, 6},
			iteratee: func(a int, b int) int { return a + b },
			expected: []int{6, 8, 3, 4}, // Zero values for missing elements
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZipT1y2(tt.slice1, tt.slice2, tt.iteratee)
			if len(result) != len(tt.expected) {
				t.Errorf("ZipT1y2() length = %d; expected %d", len(result), len(tt.expected))
				return
			}
			for i, val := range result {
				if val != tt.expected[i] {
					t.Errorf("ZipT1y2()[%d] = %v; expected %v", i, val, tt.expected[i])
				}
			}
		})
	}
}

// TestZipT1y3 tests the ZipT1y3 function with iteratee
func TestZipT1y3(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []int
		slice3   []int
		iteratee func(a int, b int, c int) int
		expected []int
	}{
		{
			name:     "sum three slices",
			slice1:   []int{1, 2},
			slice2:   []int{3, 4},
			slice3:   []int{5, 6},
			iteratee: func(a int, b int, c int) int { return a + b + c },
			expected: []int{9, 12},
		},
		{
			name:     "unequal length slices",
			slice1:   []int{1, 2, 3},
			slice2:   []int{4},
			slice3:   []int{5, 6},
			iteratee: func(a int, b int, c int) int { return a + b + c },
			expected: []int{10, 8, 3}, // Zero values for missing elements
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZipT1y3(tt.slice1, tt.slice2, tt.slice3, tt.iteratee)
			if len(result) != len(tt.expected) {
				t.Errorf("ZipT1y3() length = %d; expected %d", len(result), len(tt.expected))
				return
			}
			for i, val := range result {
				if val != tt.expected[i] {
					t.Errorf("ZipT1y3()[%d] = %v; expected %v", i, val, tt.expected[i])
				}
			}
		})
	}
}

// TestUnzip2 tests the Unzip2 function
func TestUnzip2(t *testing.T) {
	tests := []struct {
		name      string
		tuples    []Tuple2[int, string]
		expected1 []int
		expected2 []string
	}{
		{
			name: "basic unzip",
			tuples: []Tuple2[int, string]{
				{T0: 1, T1: "a"},
				{T0: 2, T1: "b"},
				{T0: 3, T1: "c"},
			},
			expected1: []int{1, 2, 3},
			expected2: []string{"a", "b", "c"},
		},
		{
			name:      "empty tuples",
			tuples:    []Tuple2[int, string]{},
			expected1: []int{},
			expected2: []string{},
		},
		{
			name: "single tuple",
			tuples: []Tuple2[int, string]{
				{T0: 42, T1: "hello"},
			},
			expected1: []int{42},
			expected2: []string{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, result2 := Unzip2(tt.tuples)
			if len(result1) != len(tt.expected1) || len(result2) != len(tt.expected2) {
				t.Errorf("Unzip2() lengths = (%d, %d); expected (%d, %d)",
					len(result1), len(result2), len(tt.expected1), len(tt.expected2))
				return
			}
			for i := range result1 {
				if result1[i] != tt.expected1[i] {
					t.Errorf("Unzip2() result1[%d] = %v; expected %v", i, result1[i], tt.expected1[i])
				}
			}
			for i := range result2 {
				if result2[i] != tt.expected2[i] {
					t.Errorf("Unzip2() result2[%d] = %v; expected %v", i, result2[i], tt.expected2[i])
				}
			}
		})
	}
}

// TestUnzip3 tests the Unzip3 function
func TestUnzip3(t *testing.T) {
	tests := []struct {
		name      string
		tuples    []Tuple3[int, string, bool]
		expected1 []int
		expected2 []string
		expected3 []bool
	}{
		{
			name: "basic unzip",
			tuples: []Tuple3[int, string, bool]{
				{T0: 1, T1: "a", T2: true},
				{T0: 2, T1: "b", T2: false},
			},
			expected1: []int{1, 2},
			expected2: []string{"a", "b"},
			expected3: []bool{true, false},
		},
		{
			name:      "empty tuples",
			tuples:    []Tuple3[int, string, bool]{},
			expected1: []int{},
			expected2: []string{},
			expected3: []bool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, result2, result3 := Unzip3(tt.tuples)
			if len(result1) != len(tt.expected1) {
				t.Errorf("Unzip3() result1 length = %d; expected %d", len(result1), len(tt.expected1))
			}
			for i := range result1 {
				if result1[i] != tt.expected1[i] || result2[i] != tt.expected2[i] || result3[i] != tt.expected3[i] {
					t.Errorf("Unzip3()[%d] = (%v, %v, %v); expected (%v, %v, %v)",
						i, result1[i], result2[i], result3[i], tt.expected1[i], tt.expected2[i], tt.expected3[i])
				}
			}
		})
	}
}

// TestUnzipT1y2 tests the UnzipT1y2 function with iteratee
func TestUnzipT1y2(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		iteratee  func(int) (int, string)
		expected1 []int
		expected2 []string
	}{
		{
			name:  "split into value and string",
			items: []int{1, 2, 3},
			iteratee: func(n int) (int, string) {
				return n * 2, string(rune('a' + n - 1))
			},
			expected1: []int{2, 4, 6},
			expected2: []string{"a", "b", "c"},
		},
		{
			name:      "empty items",
			items:     []int{},
			iteratee:  func(n int) (int, string) { return n, "" },
			expected1: []int{},
			expected2: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, result2 := UnzipT1y2(tt.items, tt.iteratee)
			if len(result1) != len(tt.expected1) || len(result2) != len(tt.expected2) {
				t.Errorf("UnzipT1y2() lengths = (%d, %d); expected (%d, %d)",
					len(result1), len(result2), len(tt.expected1), len(tt.expected2))
				return
			}
			for i := range result1 {
				if result1[i] != tt.expected1[i] {
					t.Errorf("UnzipT1y2() result1[%d] = %v; expected %v", i, result1[i], tt.expected1[i])
				}
			}
			for i := range result2 {
				if result2[i] != tt.expected2[i] {
					t.Errorf("UnzipT1y2() result2[%d] = %v; expected %v", i, result2[i], tt.expected2[i])
				}
			}
		})
	}
}

// TestUnzipT1y3 tests the UnzipT1y3 function with iteratee
func TestUnzipT1y3(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		iteratee  func(int) (int, int, bool)
		expected1 []int
		expected2 []int
		expected3 []bool
	}{
		{
			name:  "split into multiple values",
			items: []int{1, 2, 3},
			iteratee: func(n int) (int, int, bool) {
				return n, n * 2, n%2 == 0
			},
			expected1: []int{1, 2, 3},
			expected2: []int{2, 4, 6},
			expected3: []bool{false, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, result2, result3 := UnzipT1y3(tt.items, tt.iteratee)
			if len(result1) != len(tt.expected1) {
				t.Errorf("UnzipT1y3() result1 length = %d; expected %d", len(result1), len(tt.expected1))
				return
			}
			for i := range result1 {
				if result1[i] != tt.expected1[i] || result2[i] != tt.expected2[i] || result3[i] != tt.expected3[i] {
					t.Errorf("UnzipT1y3()[%d] = (%v, %v, %v); expected (%v, %v, %v)",
						i, result1[i], result2[i], result3[i], tt.expected1[i], tt.expected2[i], tt.expected3[i])
				}
			}
		})
	}
}

// TestT4ToT9 tests tuple constructors T4 through T9
func TestT4ToT9(t *testing.T) {
	t.Run("T4 constructor", func(t *testing.T) {
		result := T4(1, 2, 3, 4)
		if result.T0 != 1 || result.T1 != 2 || result.T2 != 3 || result.T3 != 4 {
			t.Errorf("T4() failed")
		}
	})

	t.Run("T5 constructor", func(t *testing.T) {
		result := T5(1, 2, 3, 4, 5)
		if result.T0 != 1 || result.T4 != 5 {
			t.Errorf("T5() failed")
		}
	})

	t.Run("T6 constructor", func(t *testing.T) {
		result := T6(1, 2, 3, 4, 5, 6)
		if result.T0 != 1 || result.T5 != 6 {
			t.Errorf("T6() failed")
		}
	})

	t.Run("T7 constructor", func(t *testing.T) {
		result := T7(1, 2, 3, 4, 5, 6, 7)
		if result.T0 != 1 || result.T6 != 7 {
			t.Errorf("T7() failed")
		}
	})

	t.Run("T8 constructor", func(t *testing.T) {
		result := T8(1, 2, 3, 4, 5, 6, 7, 8)
		if result.T0 != 1 || result.T7 != 8 {
			t.Errorf("T8() failed")
		}
	})

	t.Run("T9 constructor", func(t *testing.T) {
		result := T9(1, 2, 3, 4, 5, 6, 7, 8, 9)
		if result.T0 != 1 || result.T8 != 9 {
			t.Errorf("T9() failed")
		}
	})
}

// TestUnpack4ToUnpack9 tests unpacking functions for tuples 4-9
func TestUnpack4ToUnpack9(t *testing.T) {
	t.Run("Unpack4", func(t *testing.T) {
		tuple := T4(1, 2, 3, 4)
		v0, v1, v2, v3 := Unpack4(tuple)
		if v0 != 1 || v1 != 2 || v2 != 3 || v3 != 4 {
			t.Errorf("Unpack4() failed")
		}
	})

	t.Run("Unpack5", func(t *testing.T) {
		tuple := T5(1, 2, 3, 4, 5)
		v0, _, _, _, v4 := Unpack5(tuple)
		if v0 != 1 || v4 != 5 {
			t.Errorf("Unpack5() failed")
		}
	})

	t.Run("Unpack6", func(t *testing.T) {
		tuple := T6(1, 2, 3, 4, 5, 6)
		v0, _, _, _, _, v5 := Unpack6(tuple)
		if v0 != 1 || v5 != 6 {
			t.Errorf("Unpack6() failed")
		}
	})

	t.Run("Unpack7", func(t *testing.T) {
		tuple := T7(1, 2, 3, 4, 5, 6, 7)
		v0, _, _, _, _, _, v6 := Unpack7(tuple)
		if v0 != 1 || v6 != 7 {
			t.Errorf("Unpack7() failed")
		}
	})

	t.Run("Unpack8", func(t *testing.T) {
		tuple := T8(1, 2, 3, 4, 5, 6, 7, 8)
		v0, _, _, _, _, _, _, v7 := Unpack8(tuple)
		if v0 != 1 || v7 != 8 {
			t.Errorf("Unpack8() failed")
		}
	})

	t.Run("Unpack9", func(t *testing.T) {
		tuple := T9(1, 2, 3, 4, 5, 6, 7, 8, 9)
		v0, _, _, _, _, _, _, _, v8 := Unpack9(tuple)
		if v0 != 1 || v8 != 9 {
			t.Errorf("Unpack9() failed")
		}
	})
}

// TestZip5ToZip9 tests zip functions for 5-9 slices
func TestZip5ToZip9(t *testing.T) {
	t.Run("Zip5", func(t *testing.T) {
		result := Zip5([]int{1}, []int{2}, []int{3}, []int{4}, []int{5})
		if len(result) != 1 || result[0].T0 != 1 || result[0].T4 != 5 {
			t.Errorf("Zip5() failed")
		}
	})

	t.Run("Zip6", func(t *testing.T) {
		result := Zip6([]int{1}, []int{2}, []int{3}, []int{4}, []int{5}, []int{6})
		if len(result) != 1 || result[0].T0 != 1 || result[0].T5 != 6 {
			t.Errorf("Zip6() failed")
		}
	})

	t.Run("Zip7", func(t *testing.T) {
		result := Zip7([]int{1}, []int{2}, []int{3}, []int{4}, []int{5}, []int{6}, []int{7})
		if len(result) != 1 || result[0].T0 != 1 || result[0].T6 != 7 {
			t.Errorf("Zip7() failed")
		}
	})

	t.Run("Zip8", func(t *testing.T) {
		result := Zip8([]int{1}, []int{2}, []int{3}, []int{4}, []int{5}, []int{6}, []int{7}, []int{8})
		if len(result) != 1 || result[0].T0 != 1 || result[0].T7 != 8 {
			t.Errorf("Zip8() failed")
		}
	})

	t.Run("Zip9", func(t *testing.T) {
		result := Zip9([]int{1}, []int{2}, []int{3}, []int{4}, []int{5}, []int{6}, []int{7}, []int{8}, []int{9})
		if len(result) != 1 || result[0].T0 != 1 || result[0].T8 != 9 {
			t.Errorf("Zip9() failed")
		}
	})
}

// TestZipT1y4ToZipT1y9 tests ZipT1y functions with iteratee for 4-9 slices
func TestZipT1y4ToZipT1y9(t *testing.T) {
	t.Run("ZipT1y4", func(t *testing.T) {
		result := ZipT1y4(
			[]int{1}, []int{2}, []int{3}, []int{4},
			func(a, b, c, d int) int { return a + b + c + d },
		)
		if len(result) != 1 || result[0] != 10 {
			t.Errorf("ZipT1y4() = %v; expected [10]", result)
		}
	})

	t.Run("ZipT1y5", func(t *testing.T) {
		result := ZipT1y5(
			[]int{1}, []int{2}, []int{3}, []int{4}, []int{5},
			func(a, b, c, d, e int) int { return a + b + c + d + e },
		)
		if len(result) != 1 || result[0] != 15 {
			t.Errorf("ZipT1y5() = %v; expected [15]", result)
		}
	})

	t.Run("ZipT1y6", func(t *testing.T) {
		result := ZipT1y6(
			[]int{1}, []int{2}, []int{3}, []int{4}, []int{5}, []int{6},
			func(a, b, c, d, e, f int) int { return a + b + c + d + e + f },
		)
		if len(result) != 1 || result[0] != 21 {
			t.Errorf("ZipT1y6() = %v; expected [21]", result)
		}
	})

	t.Run("ZipT1y7 with fixed bug", func(t *testing.T) {
		// Test that the bug fix for size calculation works correctly
		result := ZipT1y7(
			[]int{1, 2}, []int{3}, []int{4}, []int{5}, []int{6}, []int{7}, []int{8, 9, 10},
			func(a, b, c, d, e, f, g int) int { return a + b + c + d + e + f + g },
		)
		// Should use max length (3) due to the bug fix
		if len(result) != 3 {
			t.Errorf("ZipT1y7() length = %d; expected 3", len(result))
		}
	})

	t.Run("ZipT1y8 with fixed bug", func(t *testing.T) {
		// Test that the bug fix for size calculation works correctly
		result := ZipT1y8(
			[]int{1}, []int{2}, []int{3}, []int{4}, []int{5}, []int{6}, []int{7}, []int{8, 9},
			func(a, b, c, d, e, f, g, h int) int { return a + b + c + d + e + f + g + h },
		)
		// Should use max length (2) due to the bug fix
		if len(result) != 2 {
			t.Errorf("ZipT1y8() length = %d; expected 2", len(result))
		}
	})

	t.Run("ZipT1y9", func(t *testing.T) {
		result := ZipT1y9(
			[]int{1}, []int{2}, []int{3}, []int{4}, []int{5}, []int{6}, []int{7}, []int{8}, []int{9},
			func(a, b, c, d, e, f, g, h, i int) int { return a + b + c + d + e + f + g + h + i },
		)
		if len(result) != 1 || result[0] != 45 {
			t.Errorf("ZipT1y9() = %v; expected [45]", result)
		}
	})
}

// TestUnzip4ToUnzip9 tests unzip functions for tuples 4-9
func TestUnzip4ToUnzip9(t *testing.T) {
	t.Run("Unzip4", func(t *testing.T) {
		tuples := []Tuple4[int, int, int, int]{
			{T0: 1, T1: 2, T2: 3, T3: 4},
			{T0: 5, T1: 6, T2: 7, T3: 8},
		}
		r1, _, _, r4 := Unzip4(tuples)
		if len(r1) != 2 || r1[0] != 1 || r4[1] != 8 {
			t.Errorf("Unzip4() failed")
		}
	})

	t.Run("Unzip5", func(t *testing.T) {
		tuples := []Tuple5[int, int, int, int, int]{
			{T0: 1, T1: 2, T2: 3, T3: 4, T4: 5},
		}
		r1, _, _, _, r5 := Unzip5(tuples)
		if len(r1) != 1 || r1[0] != 1 || r5[0] != 5 {
			t.Errorf("Unzip5() failed")
		}
	})

	t.Run("Unzip6", func(t *testing.T) {
		tuples := []Tuple6[int, int, int, int, int, int]{
			{T0: 1, T1: 2, T2: 3, T3: 4, T4: 5, T5: 6},
		}
		r1, _, _, _, _, r6 := Unzip6(tuples)
		if len(r1) != 1 || r1[0] != 1 || r6[0] != 6 {
			t.Errorf("Unzip6() failed")
		}
	})

	t.Run("Unzip7", func(t *testing.T) {
		tuples := []Tuple7[int, int, int, int, int, int, int]{
			{T0: 1, T1: 2, T2: 3, T3: 4, T4: 5, T5: 6, T6: 7},
		}
		r1, _, _, _, _, _, r7 := Unzip7(tuples)
		if len(r1) != 1 || r1[0] != 1 || r7[0] != 7 {
			t.Errorf("Unzip7() failed")
		}
	})

	t.Run("Unzip8", func(t *testing.T) {
		tuples := []Tuple8[int, int, int, int, int, int, int, int]{
			{T0: 1, T1: 2, T2: 3, T3: 4, T4: 5, T5: 6, T6: 7, T7: 8},
		}
		r1, _, _, _, _, _, _, r8 := Unzip8(tuples)
		if len(r1) != 1 || r1[0] != 1 || r8[0] != 8 {
			t.Errorf("Unzip8() failed")
		}
	})

	t.Run("Unzip9", func(t *testing.T) {
		tuples := []Tuple9[int, int, int, int, int, int, int, int, int]{
			{T0: 1, T1: 2, T2: 3, T3: 4, T4: 5, T5: 6, T6: 7, T7: 8, T8: 9},
		}
		r1, _, _, _, _, _, _, _, r9 := Unzip9(tuples)
		if len(r1) != 1 || r1[0] != 1 || r9[0] != 9 {
			t.Errorf("Unzip9() failed")
		}
	})
}

// TestUnzipT1y4ToUnzipT1y9 tests UnzipT1y functions with iteratee for 4-9 elements
func TestUnzipT1y4ToUnzipT1y9(t *testing.T) {
	t.Run("UnzipT1y4", func(t *testing.T) {
		items := []int{1, 2}
		r1, _, _, r4 := UnzipT1y4(items, func(n int) (int, int, int, int) {
			return n, n * 2, n * 3, n * 4
		})
		if len(r1) != 2 || r1[0] != 1 || r4[1] != 8 {
			t.Errorf("UnzipT1y4() failed")
		}
	})

	t.Run("UnzipT1y5", func(t *testing.T) {
		items := []int{1}
		r1, _, _, _, r5 := UnzipT1y5(items, func(n int) (int, int, int, int, int) {
			return n, n * 2, n * 3, n * 4, n * 5
		})
		if r1[0] != 1 || r5[0] != 5 {
			t.Errorf("UnzipT1y5() failed")
		}
	})

	t.Run("UnzipT1y6", func(t *testing.T) {
		items := []int{1}
		r1, _, _, _, _, r6 := UnzipT1y6(items, func(n int) (int, int, int, int, int, int) {
			return n, n * 2, n * 3, n * 4, n * 5, n * 6
		})
		if r1[0] != 1 || r6[0] != 6 {
			t.Errorf("UnzipT1y6() failed")
		}
	})

	t.Run("UnzipT1y7", func(t *testing.T) {
		items := []int{1}
		r1, _, _, _, _, _, r7 := UnzipT1y7(items, func(n int) (int, int, int, int, int, int, int) {
			return n, n * 2, n * 3, n * 4, n * 5, n * 6, n * 7
		})
		if r1[0] != 1 || r7[0] != 7 {
			t.Errorf("UnzipT1y7() failed")
		}
	})

	t.Run("UnzipT1y8", func(t *testing.T) {
		items := []int{1}
		r1, _, _, _, _, _, _, r8 := UnzipT1y8(items, func(n int) (int, int, int, int, int, int, int, int) {
			return n, n * 2, n * 3, n * 4, n * 5, n * 6, n * 7, n * 8
		})
		if r1[0] != 1 || r8[0] != 8 {
			t.Errorf("UnzipT1y8() failed")
		}
	})

	t.Run("UnzipT1y9", func(t *testing.T) {
		items := []int{1}
		r1, _, _, _, _, _, _, _, r9 := UnzipT1y9(items, func(n int) (int, int, int, int, int, int, int, int, int) {
			return n, n * 2, n * 3, n * 4, n * 5, n * 6, n * 7, n * 8, n * 9
		})
		if r1[0] != 1 || r9[0] != 9 {
			t.Errorf("UnzipT1y9() failed")
		}
	})
}

