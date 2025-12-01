package utilities

import (
	"strings"
	"testing"
)

// TestRandomString tests the RandomString function with various charsets and sizes
func TestRandomString(t *testing.T) {
	tests := []struct {
		name        string
		size        int
		charset     []rune
		expectPanic bool
	}{
		{
			name:        "valid lowercase letters",
			size:        10,
			charset:     LowerCaseLettersCharset,
			expectPanic: false,
		},
		{
			name:        "valid uppercase letters",
			size:        5,
			charset:     UpperCaseLettersCharset,
			expectPanic: false,
		},
		{
			name:        "valid alphanumeric",
			size:        15,
			charset:     AlphanumericCharset,
			expectPanic: false,
		},
		{
			name:        "valid numbers",
			size:        8,
			charset:     NumbersCharset,
			expectPanic: false,
		},
		{
			name:        "size zero should panic",
			size:        0,
			charset:     LowerCaseLettersCharset,
			expectPanic: true,
		},
		{
			name:        "negative size should panic",
			size:        -5,
			charset:     LowerCaseLettersCharset,
			expectPanic: true,
		},
		{
			name:        "empty charset should panic",
			size:        10,
			charset:     []rune{},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("RandomString() should have panicked")
					}
				}()
				RandomString(tt.size, tt.charset)
			} else {
				result := RandomString(tt.size, tt.charset)
				if len(result) != tt.size {
					t.Errorf("RandomString() length = %d; expected %d", len(result), tt.size)
				}
				// Verify all characters are from the charset
				for _, r := range result {
					found := false
					for _, c := range tt.charset {
						if r == c {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("RandomString() contains character '%c' not in charset", r)
					}
				}
			}
		})
	}
}

// TestRandomStringWithPrefix tests the RandomStringWithPrefix function
func TestRandomStringWithPrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		size    int
		charset []rune
	}{
		{
			name:    "with prefix",
			prefix:  "test_",
			size:    10,
			charset: AlphanumericCharset,
		},
		{
			name:    "empty prefix",
			prefix:  "",
			size:    5,
			charset: LowerCaseLettersCharset,
		},
		{
			name:    "long prefix",
			prefix:  "prefix_12345_",
			size:    8,
			charset: NumbersCharset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomStringWithPrefix(tt.prefix, tt.size, tt.charset)
			if !strings.HasPrefix(result, tt.prefix) {
				t.Errorf("RandomStringWithPrefix() = %v; should have prefix %v", result, tt.prefix)
			}
			if len(result) != len(tt.prefix)+tt.size {
				t.Errorf("RandomStringWithPrefix() length = %d; expected %d", len(result), len(tt.prefix)+tt.size)
			}
		})
	}
}

// TestSubstring tests the Substring function with various offset and length values
func TestSubstring(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		offset   int
		length   uint
		expected string
	}{
		{
			name:     "basic substring",
			str:      "hello world",
			offset:   0,
			length:   5,
			expected: "hello",
		},
		{
			name:     "substring from middle",
			str:      "hello world",
			offset:   6,
			length:   5,
			expected: "world",
		},
		{
			name:     "negative offset",
			str:      "hello world",
			offset:   -5,
			length:   5,
			expected: "world",
		},
		{
			name:     "negative offset beyond start",
			str:      "hello",
			offset:   -10,
			length:   3,
			expected: "hel",
		},
		{
			name:     "length exceeds remaining",
			str:      "hello",
			offset:   2,
			length:   100,
			expected: "llo",
		},
		{
			name:     "offset beyond string",
			str:      "hello",
			offset:   10,
			length:   5,
			expected: "",
		},
		{
			name:     "empty string",
			str:      "",
			offset:   0,
			length:   5,
			expected: "",
		},
		{
			name:     "unicode characters",
			str:      "こんにちは",
			offset:   0,
			length:   3,
			expected: "こんに",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Substring(tt.str, tt.offset, tt.length)
			if result != tt.expected {
				t.Errorf("Substring(%q, %d, %d) = %q; expected %q", tt.str, tt.offset, tt.length, result, tt.expected)
			}
		})
	}
}

// TestChunkString tests the ChunkString function
func TestChunkString(t *testing.T) {
	tests := []struct {
		name        string
		str         string
		size        int
		expected    []string
		expectPanic bool
	}{
		{
			name:        "even chunks",
			str:         "abcdef",
			size:        2,
			expected:    []string{"ab", "cd", "ef"},
			expectPanic: false,
		},
		{
			name:        "uneven chunks",
			str:         "abcdefgh",
			size:        3,
			expected:    []string{"abc", "def", "gh"},
			expectPanic: false,
		},
		{
			name:        "size equals string length",
			str:         "hello",
			size:        5,
			expected:    []string{"hello"},
			expectPanic: false,
		},
		{
			name:        "size greater than string length",
			str:         "hi",
			size:        10,
			expected:    []string{"hi"},
			expectPanic: false,
		},
		{
			name:        "empty string",
			str:         "",
			size:        3,
			expected:    []string{""},
			expectPanic: false,
		},
		{
			name:        "size zero should panic",
			str:         "hello",
			size:        0,
			expected:    nil,
			expectPanic: true,
		},
		{
			name:        "negative size should panic",
			str:         "hello",
			size:        -1,
			expected:    nil,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("ChunkString() should have panicked")
					}
				}()
				ChunkString(tt.str, tt.size)
			} else {
				result := ChunkString(tt.str, tt.size)
				if len(result) != len(tt.expected) {
					t.Errorf("ChunkString() length = %d; expected %d", len(result), len(tt.expected))
					return
				}
				for i, chunk := range result {
					if chunk != tt.expected[i] {
						t.Errorf("ChunkString() chunk[%d] = %q; expected %q", i, chunk, tt.expected[i])
					}
				}
			}
		})
	}
}

// TestRuneLength tests the RuneLength function
func TestRuneLength(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected int
	}{
		{
			name:     "ascii string",
			str:      "hello",
			expected: 5,
		},
		{
			name:     "unicode string",
			str:      "こんにちは",
			expected: 5,
		},
		{
			name:     "mixed ascii and unicode",
			str:      "hello世界",
			expected: 7,
		},
		{
			name:     "empty string",
			str:      "",
			expected: 0,
		},
		{
			name:     "emojis",
			str:      "😀😃😄",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RuneLength(tt.str)
			if result != tt.expected {
				t.Errorf("RuneLength(%q) = %d; expected %d", tt.str, result, tt.expected)
			}
		})
	}
}

// TestPascalCase tests the PascalCase function
func TestPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "snake case",
			input:    "hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "kebab case",
			input:    "hello-world",
			expected: "HelloWorld",
		},
		{
			name:     "camel case",
			input:    "helloWorld",
			expected: "HelloWorld",
		},
		{
			name:     "multiple words",
			input:    "the_quick_brown_fox",
			expected: "TheQuickBrownFox",
		},
		{
			name:     "with numbers",
			input:    "hello2world",
			expected: "Hello2World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("PascalCase(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCamelCase tests the CamelCase function
func TestCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "snake case",
			input:    "hello_world",
			expected: "helloWorld",
		},
		{
			name:     "kebab case",
			input:    "hello-world",
			expected: "helloWorld",
		},
		{
			name:     "pascal case",
			input:    "HelloWorld",
			expected: "helloWorld",
		},
		{
			name:     "multiple words",
			input:    "the_quick_brown_fox",
			expected: "theQuickBrownFox",
		},
		{
			name:     "with numbers",
			input:    "hello2world",
			expected: "hello2World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("CamelCase(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestKebabCase tests the KebabCase function
func TestKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "snake case",
			input:    "hello_world",
			expected: "hello-world",
		},
		{
			name:     "camel case",
			input:    "helloWorld",
			expected: "hello-world",
		},
		{
			name:     "pascal case",
			input:    "HelloWorld",
			expected: "hello-world",
		},
		{
			name:     "multiple words",
			input:    "the_quick_brown_fox",
			expected: "the-quick-brown-fox",
		},
		{
			name:     "with numbers",
			input:    "hello2World",
			expected: "hello-2-world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("KebabCase(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSnakeCase tests the SnakeCase function
func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camel case",
			input:    "helloWorld",
			expected: "hello_world",
		},
		{
			name:     "pascal case",
			input:    "HelloWorld",
			expected: "hello_world",
		},
		{
			name:     "kebab case",
			input:    "hello-world",
			expected: "hello_world",
		},
		{
			name:     "multiple words",
			input:    "theQuickBrownFox",
			expected: "the_quick_brown_fox",
		},
		{
			name:     "with numbers",
			input:    "hello2World",
			expected: "hello_2_world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("SnakeCase(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestWords tests the Words function
func TestWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "camel case",
			input:    "helloWorld",
			expected: []string{"hello", "World"},
		},
		{
			name:     "pascal case",
			input:    "HelloWorld",
			expected: []string{"Hello", "World"},
		},
		{
			name:     "snake case",
			input:    "hello_world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "kebab case",
			input:    "hello-world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "with numbers",
			input:    "hello2World",
			expected: []string{"hello", "2", "World"},
		},
		{
			name:     "multiple separators",
			input:    "hello__world--test",
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Words(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Words(%q) length = %d; expected %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i, word := range result {
				if word != tt.expected[i] {
					t.Errorf("Words(%q)[%d] = %q; expected %q", tt.input, i, word, tt.expected[i])
				}
			}
		})
	}
}

// TestCapitalize tests the Capitalize function
func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase word",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "uppercase word",
			input:    "HELLO",
			expected: "Hello",
		},
		{
			name:     "mixed case",
			input:    "hElLo",
			expected: "Hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Capitalize(tt.input)
			if result != tt.expected {
				t.Errorf("Capitalize(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestEllipsis tests the Ellipsis function
func TestEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected string
	}{
		{
			name:     "string shorter than length",
			input:    "hello",
			length:   10,
			expected: "hello",
		},
		{
			name:     "string equal to length",
			input:    "hello",
			length:   5,
			expected: "hello",
		},
		{
			name:     "string longer than length",
			input:    "hello world",
			length:   8,
			expected: "hello...",
		},
		{
			name:     "very short length",
			input:    "hello",
			length:   2,
			expected: "...",
		},
		{
			name:     "length of 3",
			input:    "hello world",
			length:   3,
			expected: "...",
		},
		{
			name:     "string with spaces",
			input:    "  hello world  ",
			length:   8,
			expected: "hello...",
		},
		{
			name:     "empty string",
			input:    "",
			length:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Ellipsis(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("Ellipsis(%q, %d) = %q; expected %q", tt.input, tt.length, result, tt.expected)
			}
		})
	}
}

// TestHashString tests the HashString function
func TestHashString(t *testing.T) {
	tests := []struct {
		name        string
		payload     interface{}
		expectError bool
	}{
		{
			name:        "string payload",
			payload:     "hello world",
			expectError: false,
		},
		{
			name:        "number payload",
			payload:     42,
			expectError: false,
		},
		{
			name: "struct payload",
			payload: struct {
				Name string
				Age  int
			}{"John", 30},
			expectError: false,
		},
		{
			name:        "map payload",
			payload:     map[string]int{"a": 1, "b": 2},
			expectError: false,
		},
		{
			name:        "nil payload",
			payload:     nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HashString(tt.payload)
			if tt.expectError && err == nil {
				t.Errorf("HashString() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("HashString() unexpected error: %v", err)
			}
			if !tt.expectError && result == "" {
				t.Errorf("HashString() returned empty string")
			}
			// Verify hash is consistent
			if !tt.expectError {
				result2, _ := HashString(tt.payload)
				if result != result2 {
					t.Errorf("HashString() not consistent: %s != %s", result, result2)
				}
			}
		})
	}
}

// TestReverseString tests the ReverseString function
func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: "olleh",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "palindrome",
			input:    "racecar",
			expected: "racecar",
		},
		{
			name:     "unicode characters",
			input:    "こんにちは",
			expected: "はちにんこ",
		},
		{
			name:     "with spaces",
			input:    "hello world",
			expected: "dlrow olleh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReverseString(tt.input)
			if result != tt.expected {
				t.Errorf("ReverseString(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestBuildQueryUri tests the BuildQueryUri struct methods
func TestBuildQueryUri(t *testing.T) {
	tests := []struct {
		name     string
		baseUri  string
		params   map[string]interface{}
		expected string
	}{
		{
			name:     "no parameters",
			baseUri:  "/api/users",
			params:   nil,
			expected: "/api/users",
		},
		{
			name:    "single parameter",
			baseUri: "/api/users",
			params: map[string]interface{}{
				"id": 1,
			},
			expected: "/api/users?id=1",
		},
		{
			name:    "multiple parameters",
			baseUri: "/api/users",
			params: map[string]interface{}{
				"id":   1,
				"name": "john",
			},
			expected: "", // Will check if contains both params
		},
		{
			name:    "various types",
			baseUri: "/api/test",
			params: map[string]interface{}{
				"str":  "value",
				"int":  42,
				"bool": true,
			},
			expected: "", // Will check if contains all params
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &BuildQueryUri{QueryUri: tt.baseUri}
			for k, v := range tt.params {
				builder.AddParam(k, v)
			}
			result := builder.Build()

			if tt.expected != "" {
				if result != tt.expected {
					t.Errorf("BuildQueryUri.Build() = %q; expected %q", result, tt.expected)
				}
			} else {
				// For multiple params, just check they are present
				if !strings.HasPrefix(result, tt.baseUri+"?") && len(tt.params) > 0 {
					t.Errorf("BuildQueryUri.Build() should start with %q?", tt.baseUri)
				}
				for k := range tt.params {
					if !strings.Contains(result, k+"=") {
						t.Errorf("BuildQueryUri.Build() should contain param %q", k)
					}
				}
			}
		})
	}
}

// TestBuildQueryUriAddParam tests adding parameters incrementally
func TestBuildQueryUriAddParam(t *testing.T) {
	t.Run("add params to nil map", func(t *testing.T) {
		builder := &BuildQueryUri{QueryUri: "/api/test"}
		builder.AddParam("key", "value")
		result := builder.Build()
		if !strings.Contains(result, "key=value") {
			t.Errorf("BuildQueryUri should contain key=value, got %s", result)
		}
	})

	t.Run("add multiple params", func(t *testing.T) {
		builder := &BuildQueryUri{QueryUri: "/api/test"}
		builder.AddParam("key1", "value1")
		builder.AddParam("key2", "value2")
		result := builder.Build()
		if !strings.Contains(result, "key1=value1") {
			t.Errorf("BuildQueryUri should contain key1=value1")
		}
		if !strings.Contains(result, "key2=value2") {
			t.Errorf("BuildQueryUri should contain key2=value2")
		}
	})
}
