package utilities

import (
	"regexp"
	"sync"
)

// Regex represents a regular expression pattern as a string
type Regex string

// Common regex patterns for validation
const (
	VietnamPhone Regex = "^(0[35789]\\d{8}|\\+84[35789]\\d{8})$"
	Email        Regex = "^[a-zA-Z0-9._%+\\-]+@[a-zA-Z0-9.\\-]+\\.[a-zA-Z]{2,}$"
	Alphanumeric Regex = "^[a-zA-Z0-9]+$"
	Alpha        Regex = "^[a-zA-Z]+$"
	Digit        Regex = "^\\d+$"
	// StrongPassword validates passwords with at least 8 characters
	// Note: Go regexp doesn't support lookaheads, so this pattern checks for
	// a mix of character types but cannot enforce all requirements simultaneously
	StrongPassword Regex = "^[a-zA-Z0-9\\W_]{8,}$"
	URL            Regex = "^(https?://)?([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}(/.*)?$"
	IPv4           Regex = "^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	HexColor       Regex = "^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
	CreditCard     Regex = "^[0-9]{13,19}$"
	Date           Regex = "^\\d{4}-\\d{2}-\\d{2}$" // YYYY-MM-DD
)

// regexCache stores compiled regex patterns for performance
var (
	regexCache = make(map[Regex]*regexp.Regexp)
	cacheMutex sync.RWMutex
)

// getCompiledRegex returns a cached compiled regex or compiles and caches a new one
func getCompiledRegex(regex Regex) (*regexp.Regexp, error) {
	// Try to get from cache with read lock
	cacheMutex.RLock()
	if compiled, exists := regexCache[regex]; exists {
		cacheMutex.RUnlock()
		return compiled, nil
	}
	cacheMutex.RUnlock()

	// Compile with write lock
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if compiled, exists := regexCache[regex]; exists {
		return compiled, nil
	}

	// Compile and cache
	compiled, err := regexp.Compile(string(regex))
	if err != nil {
		return nil, err
	}
	regexCache[regex] = compiled
	return compiled, nil
}

// CheckRegex checks if the input string matches the given regex pattern
// Returns true if the input matches, false otherwise
func CheckRegex(regex Regex, input string) bool {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return false
	}
	return compiled.MatchString(input)
}

// CheckRegexWithError checks if the input string matches the given regex pattern
// Returns match result and any compilation error
func CheckRegexWithError(regex Regex, input string) (bool, error) {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return false, err
	}
	return compiled.MatchString(input), nil
}

// FindAllMatches returns all matches of the regex pattern in the input string
// Returns empty slice if no matches are found or if regex is invalid
func FindAllMatches(regex Regex, input string) []string {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return []string{}
	}
	return compiled.FindAllString(input, -1)
}

// FindFirstMatch returns the first match of the regex pattern in the input string
// Returns empty string if no match is found or if regex is invalid
func FindFirstMatch(regex Regex, input string) string {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return ""
	}
	return compiled.FindString(input)
}

// FindSubmatch returns the first match and its submatches (captured groups)
// Returns empty slice if no match is found or if regex is invalid
func FindSubmatch(regex Regex, input string) []string {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return []string{}
	}
	return compiled.FindStringSubmatch(input)
}

// FindAllSubmatches returns all matches and their submatches (captured groups)
// Returns empty slice if no matches are found or if regex is invalid
func FindAllSubmatches(regex Regex, input string) [][]string {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return [][]string{}
	}
	return compiled.FindAllStringSubmatch(input, -1)
}

// ReplaceAllMatches replaces all matches of the regex pattern with the replacement string
// Returns the original string if regex is invalid
func ReplaceAllMatches(regex Regex, input, replacement string) string {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return input
	}
	return compiled.ReplaceAllString(input, replacement)
}

// ReplaceFirstMatch replaces the first match of the regex pattern with the replacement string
// Returns the original string if regex is invalid
func ReplaceFirstMatch(regex Regex, input, replacement string) string {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return input
	}
	match := compiled.FindString(input)
	if match == "" {
		return input
	}
	index := compiled.FindStringIndex(input)
	if index == nil {
		return input
	}
	return input[:index[0]] + replacement + input[index[1]:]
}

// Split splits the input string by the regex pattern
// Returns slice with original string if regex is invalid or no matches found
func Split(regex Regex, input string) []string {
	compiled, err := getCompiledRegex(regex)
	if err != nil {
		return []string{input}
	}
	return compiled.Split(input, -1)
}

// CountMatches returns the number of matches of the regex pattern in the input string
// Returns 0 if no matches are found or if regex is invalid
func CountMatches(regex Regex, input string) int {
	matches := FindAllMatches(regex, input)
	return len(matches)
}

// IsValid checks if the regex pattern is valid (can be compiled)
func IsValid(regex Regex) bool {
	_, err := regexp.Compile(string(regex))
	return err == nil
}

// ClearCache clears the compiled regex cache
// Useful for memory management in long-running applications
func ClearCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	regexCache = make(map[Regex]*regexp.Regexp)
}

// GetCacheSize returns the number of compiled regex patterns in the cache
func GetCacheSize() int {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return len(regexCache)
}
