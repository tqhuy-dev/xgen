package utilities

import (
	"reflect"
	"testing"
)

// TestCheckRegex_VietnamPhone tests the CheckRegex function with Vietnam phone number pattern
// which validates Vietnamese phone numbers in both local and international formats
func TestCheckRegex_VietnamPhone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid phone starting with 03",
			input:    "0312345678",
			expected: true,
		},
		{
			name:     "valid phone starting with 05",
			input:    "0512345678",
			expected: true,
		},
		{
			name:     "valid phone starting with 07",
			input:    "0712345678",
			expected: true,
		},
		{
			name:     "valid phone starting with 08",
			input:    "0812345678",
			expected: true,
		},
		{
			name:     "valid phone starting with 09",
			input:    "0912345678",
			expected: true,
		},
		{
			name:     "valid international format",
			input:    "+84912345678",
			expected: true,
		},
		{
			name:     "invalid starting digit 01",
			input:    "0112345678",
			expected: false,
		},
		{
			name:     "too short",
			input:    "091234567",
			expected: false,
		},
		{
			name:     "too long",
			input:    "09123456789",
			expected: false,
		},
		{
			name:     "contains letters",
			input:    "091234567a",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(VietnamPhone, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(VietnamPhone, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_Email tests the CheckRegex function with email pattern
// which validates email addresses in standard format
func TestCheckRegex_Email(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid simple email",
			input:    "test@example.com",
			expected: true,
		},
		{
			name:     "valid email with plus",
			input:    "test+tag@example.com",
			expected: true,
		},
		{
			name:     "valid email with dots",
			input:    "first.last@example.com",
			expected: true,
		},
		{
			name:     "valid email with numbers",
			input:    "user123@example123.com",
			expected: true,
		},
		{
			name:     "valid email with subdomain",
			input:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "missing at symbol",
			input:    "testexample.com",
			expected: false,
		},
		{
			name:     "missing domain",
			input:    "test@",
			expected: false,
		},
		{
			name:     "missing local part",
			input:    "@example.com",
			expected: false,
		},
		{
			name:     "invalid characters",
			input:    "test@#@example.com",
			expected: false,
		},
		{
			name:     "no TLD",
			input:    "test@example",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(Email, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(Email, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_Alphanumeric tests the CheckRegex function with alphanumeric pattern
// which validates strings containing only letters and numbers
func TestCheckRegex_Alphanumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid alphanumeric",
			input:    "abc123",
			expected: true,
		},
		{
			name:     "only letters",
			input:    "abcdef",
			expected: true,
		},
		{
			name:     "only numbers",
			input:    "123456",
			expected: true,
		},
		{
			name:     "mixed case",
			input:    "AbC123XyZ",
			expected: true,
		},
		{
			name:     "contains space",
			input:    "abc 123",
			expected: false,
		},
		{
			name:     "contains special char",
			input:    "abc-123",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(Alphanumeric, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(Alphanumeric, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_Alpha tests the CheckRegex function with alphabetic pattern
// which validates strings containing only letters
func TestCheckRegex_Alpha(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid lowercase",
			input:    "abcdef",
			expected: true,
		},
		{
			name:     "valid uppercase",
			input:    "ABCDEF",
			expected: true,
		},
		{
			name:     "valid mixed case",
			input:    "AbCdEf",
			expected: true,
		},
		{
			name:     "contains numbers",
			input:    "abc123",
			expected: false,
		},
		{
			name:     "contains space",
			input:    "abc def",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(Alpha, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(Alpha, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_Digit tests the CheckRegex function with digit pattern
// which validates strings containing only digits
func TestCheckRegex_Digit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid digits",
			input:    "123456",
			expected: true,
		},
		{
			name:     "single digit",
			input:    "0",
			expected: true,
		},
		{
			name:     "long number",
			input:    "1234567890",
			expected: true,
		},
		{
			name:     "contains letters",
			input:    "123abc",
			expected: false,
		},
		{
			name:     "contains space",
			input:    "123 456",
			expected: false,
		},
		{
			name:     "negative sign",
			input:    "-123",
			expected: false,
		},
		{
			name:     "decimal point",
			input:    "123.45",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(Digit, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(Digit, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_StrongPassword tests the CheckRegex function with strong password pattern
// which validates passwords with at least 8 characters containing alphanumeric and special chars
// Note: Go regexp doesn't support lookaheads, so this cannot enforce all strong password rules
func TestCheckRegex_StrongPassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid strong password",
			input:    "Passw0rd!",
			expected: true,
		},
		{
			name:     "valid with underscore",
			input:    "Test_1aA",
			expected: true,
		},
		{
			name:     "valid long password",
			input:    "MyP@ssw0rd123",
			expected: true,
		},
		{
			name:     "valid lowercase only with length",
			input:    "passw0rd!",
			expected: true,
		},
		{
			name:     "valid uppercase only with length",
			input:    "PASSW0RD!",
			expected: true,
		},
		{
			name:     "valid without digit",
			input:    "Password!",
			expected: true,
		},
		{
			name:     "valid alphanumeric only",
			input:    "Passw0rd",
			expected: true,
		},
		{
			name:     "too short",
			input:    "Pass1!",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "exactly 8 characters",
			input:    "Test1234",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(StrongPassword, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(StrongPassword, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_URL tests the CheckRegex function with URL pattern
// which validates URLs in various formats
func TestCheckRegex_URL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid http URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "valid https URL",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "valid URL without protocol",
			input:    "example.com",
			expected: true,
		},
		{
			name:     "valid URL with path",
			input:    "https://example.com/path/to/page",
			expected: true,
		},
		{
			name:     "valid URL with subdomain",
			input:    "https://www.example.com",
			expected: true,
		},
		{
			name:     "valid URL with query",
			input:    "https://example.com/page?query=test",
			expected: true,
		},
		{
			name:     "missing domain",
			input:    "https://",
			expected: false,
		},
		{
			name:     "invalid protocol",
			input:    "ftp://example.com",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(URL, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(URL, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_IPv4 tests the CheckRegex function with IPv4 pattern
// which validates IPv4 addresses
func TestCheckRegex_IPv4(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid IPv4",
			input:    "192.168.1.1",
			expected: true,
		},
		{
			name:     "valid IPv4 with zeros",
			input:    "0.0.0.0",
			expected: true,
		},
		{
			name:     "valid IPv4 max values",
			input:    "255.255.255.255",
			expected: true,
		},
		{
			name:     "valid IPv4 localhost",
			input:    "127.0.0.1",
			expected: true,
		},
		{
			name:     "invalid octet over 255",
			input:    "256.1.1.1",
			expected: false,
		},
		{
			name:     "invalid missing octet",
			input:    "192.168.1",
			expected: false,
		},
		{
			name:     "invalid too many octets",
			input:    "192.168.1.1.1",
			expected: false,
		},
		{
			name:     "invalid contains letters",
			input:    "192.168.1.a",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(IPv4, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(IPv4, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_HexColor tests the CheckRegex function with hex color pattern
// which validates hexadecimal color codes
func TestCheckRegex_HexColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid 6-digit hex",
			input:    "#FF5733",
			expected: true,
		},
		{
			name:     "valid 3-digit hex",
			input:    "#F53",
			expected: true,
		},
		{
			name:     "valid lowercase",
			input:    "#ff5733",
			expected: true,
		},
		{
			name:     "valid mixed case",
			input:    "#Ff5733",
			expected: true,
		},
		{
			name:     "missing hash",
			input:    "FF5733",
			expected: false,
		},
		{
			name:     "invalid length",
			input:    "#FF57",
			expected: false,
		},
		{
			name:     "invalid characters",
			input:    "#GG5733",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(HexColor, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(HexColor, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_CreditCard tests the CheckRegex function with credit card pattern
// which validates credit card numbers (13-19 digits)
func TestCheckRegex_CreditCard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid 16-digit card",
			input:    "4532015112830366",
			expected: true,
		},
		{
			name:     "valid 13-digit card",
			input:    "4532015112830",
			expected: true,
		},
		{
			name:     "valid 19-digit card",
			input:    "4532015112830366123",
			expected: true,
		},
		{
			name:     "too short",
			input:    "453201511283",
			expected: false,
		},
		{
			name:     "too long",
			input:    "45320151128303661234",
			expected: false,
		},
		{
			name:     "contains spaces",
			input:    "4532 0151 1283 0366",
			expected: false,
		},
		{
			name:     "contains letters",
			input:    "453201511283036a",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(CreditCard, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(CreditCard, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegex_Date tests the CheckRegex function with date pattern
// which validates dates in YYYY-MM-DD format
func TestCheckRegex_Date(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid date",
			input:    "2025-11-09",
			expected: true,
		},
		{
			name:     "valid leap year",
			input:    "2024-02-29",
			expected: true,
		},
		{
			name:     "valid start of year",
			input:    "2025-01-01",
			expected: true,
		},
		{
			name:     "invalid format MM-DD-YYYY",
			input:    "11-09-2025",
			expected: false,
		},
		{
			name:     "invalid separator",
			input:    "2025/11/09",
			expected: false,
		},
		{
			name:     "missing leading zero",
			input:    "2025-1-9",
			expected: false,
		},
		{
			name:     "two-digit year",
			input:    "25-11-09",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CheckRegex(Date, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CheckRegex(Date, %q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckRegexWithError tests the CheckRegexWithError function
// which returns both the match result and any compilation error
func TestCheckRegexWithError(t *testing.T) {
	tests := []struct {
		name          string
		regex         Regex
		input         string
		expectedMatch bool
		expectError   bool
	}{
		{
			name:          "valid regex and matching input",
			regex:         Digit,
			input:         "123",
			expectedMatch: true,
			expectError:   false,
		},
		{
			name:          "valid regex and non-matching input",
			regex:         Digit,
			input:         "abc",
			expectedMatch: false,
			expectError:   false,
		},
		{
			name:          "invalid regex pattern",
			regex:         Regex("[invalid"),
			input:         "test",
			expectedMatch: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			match, err := CheckRegexWithError(tt.regex, tt.input)

			// Assert
			if match != tt.expectedMatch {
				t.Errorf("CheckRegexWithError() match = %v; expected %v", match, tt.expectedMatch)
			}
			if (err != nil) != tt.expectError {
				t.Errorf("CheckRegexWithError() error = %v; expected error = %v", err, tt.expectError)
			}
		})
	}
}

// TestFindAllMatches tests the FindAllMatches function
// which returns all matches of a regex pattern in the input
func TestFindAllMatches(t *testing.T) {
	tests := []struct {
		name     string
		regex    Regex
		input    string
		expected []string
	}{
		{
			name:     "find all digits",
			regex:    Regex("\\d+"),
			input:    "abc123def456ghi789",
			expected: []string{"123", "456", "789"},
		},
		{
			name:     "find all words",
			regex:    Regex("[a-z]+"),
			input:    "hello world test",
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "no matches",
			regex:    Regex("\\d+"),
			input:    "abcdef",
			expected: []string{},
		},
		{
			name:     "empty input",
			regex:    Regex("\\d+"),
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FindAllMatches(tt.regex, tt.input)

			// Assert
			if len(result) != len(tt.expected) {
				t.Errorf("FindAllMatches() length = %d; expected %d", len(result), len(tt.expected))
			} else if len(result) > 0 && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindAllMatches() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestFindFirstMatch tests the FindFirstMatch function
// which returns the first match of a regex pattern
func TestFindFirstMatch(t *testing.T) {
	tests := []struct {
		name     string
		regex    Regex
		input    string
		expected string
	}{
		{
			name:     "find first digit sequence",
			regex:    Regex("\\d+"),
			input:    "abc123def456",
			expected: "123",
		},
		{
			name:     "find first word",
			regex:    Regex("[a-z]+"),
			input:    "hello world",
			expected: "hello",
		},
		{
			name:     "no match",
			regex:    Regex("\\d+"),
			input:    "abcdef",
			expected: "",
		},
		{
			name:     "empty input",
			regex:    Regex("\\d+"),
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FindFirstMatch(tt.regex, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("FindFirstMatch() = %q; expected %q", result, tt.expected)
			}
		})
	}
}

// TestFindSubmatch tests the FindSubmatch function
// which returns the first match and its captured groups
func TestFindSubmatch(t *testing.T) {
	tests := []struct {
		name     string
		regex    Regex
		input    string
		expected []string
	}{
		{
			name:     "capture groups in date",
			regex:    Regex("(\\d{4})-(\\d{2})-(\\d{2})"),
			input:    "Date: 2025-11-09",
			expected: []string{"2025-11-09", "2025", "11", "09"},
		},
		{
			name:     "capture email parts",
			regex:    Regex("([a-z]+)@([a-z]+)\\.([a-z]+)"),
			input:    "test@example.com",
			expected: []string{"test@example.com", "test", "example", "com"},
		},
		{
			name:     "no match",
			regex:    Regex("(\\d+)"),
			input:    "abcdef",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FindSubmatch(tt.regex, tt.input)

			// Assert
			if len(result) != len(tt.expected) {
				t.Errorf("FindSubmatch() length = %d; expected %d", len(result), len(tt.expected))
			} else if len(result) > 0 && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindSubmatch() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestFindAllSubmatches tests the FindAllSubmatches function
// which returns all matches and their captured groups
func TestFindAllSubmatches(t *testing.T) {
	tests := []struct {
		name     string
		regex    Regex
		input    string
		expected [][]string
	}{
		{
			name:  "capture multiple phone patterns",
			regex: Regex("(\\d{3})-(\\d{3})-(\\d{4})"),
			input: "Call 123-456-7890 or 987-654-3210",
			expected: [][]string{
				{"123-456-7890", "123", "456", "7890"},
				{"987-654-3210", "987", "654", "3210"},
			},
		},
		{
			name:     "no matches",
			regex:    Regex("(\\d+)"),
			input:    "abcdef",
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FindAllSubmatches(tt.regex, tt.input)

			// Assert
			if len(result) != len(tt.expected) {
				t.Errorf("FindAllSubmatches() length = %d; expected %d", len(result), len(tt.expected))
			} else if len(result) > 0 && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindAllSubmatches() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestReplaceAllMatches tests the ReplaceAllMatches function
// which replaces all occurrences of a pattern
func TestReplaceAllMatches(t *testing.T) {
	tests := []struct {
		name        string
		regex       Regex
		input       string
		replacement string
		expected    string
	}{
		{
			name:        "replace all digits",
			regex:       Regex("\\d+"),
			input:       "abc123def456",
			replacement: "X",
			expected:    "abcXdefX",
		},
		{
			name:        "replace all spaces",
			regex:       Regex("\\s+"),
			input:       "hello world test",
			replacement: "_",
			expected:    "hello_world_test",
		},
		{
			name:        "no matches to replace",
			regex:       Regex("\\d+"),
			input:       "abcdef",
			replacement: "X",
			expected:    "abcdef",
		},
		{
			name:        "replace with empty string",
			regex:       Regex("\\d+"),
			input:       "abc123def456",
			replacement: "",
			expected:    "abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ReplaceAllMatches(tt.regex, tt.input, tt.replacement)

			// Assert
			if result != tt.expected {
				t.Errorf("ReplaceAllMatches() = %q; expected %q", result, tt.expected)
			}
		})
	}
}

// TestReplaceFirstMatch tests the ReplaceFirstMatch function
// which replaces only the first occurrence of a pattern
func TestReplaceFirstMatch(t *testing.T) {
	tests := []struct {
		name        string
		regex       Regex
		input       string
		replacement string
		expected    string
	}{
		{
			name:        "replace first digit",
			regex:       Regex("\\d+"),
			input:       "abc123def456",
			replacement: "X",
			expected:    "abcXdef456",
		},
		{
			name:        "replace first word",
			regex:       Regex("[a-z]+"),
			input:       "hello world",
			replacement: "hi",
			expected:    "hi world",
		},
		{
			name:        "no match to replace",
			regex:       Regex("\\d+"),
			input:       "abcdef",
			replacement: "X",
			expected:    "abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ReplaceFirstMatch(tt.regex, tt.input, tt.replacement)

			// Assert
			if result != tt.expected {
				t.Errorf("ReplaceFirstMatch() = %q; expected %q", result, tt.expected)
			}
		})
	}
}

// TestSplit tests the Split function
// which splits a string by a regex pattern
func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		regex    Regex
		input    string
		expected []string
	}{
		{
			name:     "split by comma",
			regex:    Regex(","),
			input:    "a,b,c,d",
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "split by whitespace",
			regex:    Regex("\\s+"),
			input:    "hello world test",
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "split by digits",
			regex:    Regex("\\d+"),
			input:    "a123b456c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "no match",
			regex:    Regex(","),
			input:    "abcdef",
			expected: []string{"abcdef"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := Split(tt.regex, tt.input)

			// Assert
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Split() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestCountMatches tests the CountMatches function
// which counts the number of matches of a pattern
func TestCountMatches(t *testing.T) {
	tests := []struct {
		name     string
		regex    Regex
		input    string
		expected int
	}{
		{
			name:     "count digits",
			regex:    Regex("\\d+"),
			input:    "abc123def456ghi789",
			expected: 3,
		},
		{
			name:     "count words",
			regex:    Regex("[a-z]+"),
			input:    "hello world test",
			expected: 3,
		},
		{
			name:     "no matches",
			regex:    Regex("\\d+"),
			input:    "abcdef",
			expected: 0,
		},
		{
			name:     "empty input",
			regex:    Regex("\\d+"),
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CountMatches(tt.regex, tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("CountMatches() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestIsValid tests the IsValid function
// which checks if a regex pattern is valid
func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		regex    Regex
		expected bool
	}{
		{
			name:     "valid simple pattern",
			regex:    Regex("\\d+"),
			expected: true,
		},
		{
			name:     "valid complex pattern",
			regex:    Regex("^[a-zA-Z0-9]+@[a-zA-Z0-9]+\\.[a-zA-Z]{2,}$"),
			expected: true,
		},
		{
			name:     "invalid pattern - unclosed bracket",
			regex:    Regex("[invalid"),
			expected: false,
		},
		{
			name:     "invalid pattern - unclosed group",
			regex:    Regex("(unclosed"),
			expected: false,
		},
		{
			name:     "empty pattern",
			regex:    Regex(""),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := IsValid(tt.regex)

			// Assert
			if result != tt.expected {
				t.Errorf("IsValid(%q) = %v; expected %v", tt.regex, result, tt.expected)
			}
		})
	}
}

// TestCachePerformance tests that the regex cache improves performance
// by verifying that compiled regexes are cached
func TestCachePerformance(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "verify cache is populated after first use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ClearCache()
			initialSize := GetCacheSize()

			// Act: Use a regex pattern multiple times
			pattern := Regex("\\d+")
			CheckRegex(pattern, "123")
			CheckRegex(pattern, "456")
			finalSize := GetCacheSize()

			// Assert
			if initialSize != 0 {
				t.Errorf("Initial cache size = %d; expected 0", initialSize)
			}
			if finalSize != 1 {
				t.Errorf("Final cache size = %d; expected 1 (pattern should be cached)", finalSize)
			}
		})
	}
}

// TestClearCache tests the ClearCache function
// which removes all compiled regex patterns from the cache
func TestClearCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "clear cache removes all entries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Populate cache
			CheckRegex(Digit, "123")
			CheckRegex(Alpha, "abc")
			CheckRegex(Email, "test@example.com")

			// Act
			ClearCache()
			size := GetCacheSize()

			// Assert
			if size != 0 {
				t.Errorf("Cache size after clear = %d; expected 0", size)
			}
		})
	}
}

// TestGetCacheSize tests the GetCacheSize function
// which returns the number of cached compiled regex patterns
func TestGetCacheSize(t *testing.T) {
	tests := []struct {
		name            string
		patternsToCache int
		expected        int
	}{
		{
			name:            "cache three patterns",
			patternsToCache: 3,
			expected:        3,
		},
		{
			name:            "cache one pattern",
			patternsToCache: 1,
			expected:        1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ClearCache()
			patterns := []Regex{Digit, Alpha, Email}

			// Act
			for i := 0; i < tt.patternsToCache && i < len(patterns); i++ {
				CheckRegex(patterns[i], "test")
			}
			size := GetCacheSize()

			// Assert
			if size != tt.expected {
				t.Errorf("GetCacheSize() = %d; expected %d", size, tt.expected)
			}
		})
	}
}

// TestConcurrentAccess tests that the regex cache is thread-safe
// by performing concurrent operations on the cache
func TestConcurrentAccess(t *testing.T) {
	tests := []struct {
		name           string
		goroutines     int
		operationsEach int
	}{
		{
			name:           "concurrent access with 10 goroutines",
			goroutines:     10,
			operationsEach: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ClearCache()
			done := make(chan bool, tt.goroutines)

			// Act: Multiple goroutines accessing cache concurrently
			for i := 0; i < tt.goroutines; i++ {
				go func() {
					for j := 0; j < tt.operationsEach; j++ {
						CheckRegex(Digit, "123")
						CheckRegex(Alpha, "abc")
						CheckRegex(Email, "test@example.com")
					}
					done <- true
				}()
			}

			// Wait for all goroutines to complete
			for i := 0; i < tt.goroutines; i++ {
				<-done
			}

			// Assert: Cache should have 3 entries (not more due to race conditions)
			size := GetCacheSize()
			if size != 3 {
				t.Errorf("Cache size after concurrent access = %d; expected 3", size)
			}
		})
	}
}
