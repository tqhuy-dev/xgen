// Package codebase mirrors base/schema.py: JSON shape for symbols and call edges.
package codebase

import (
	"encoding/json"
	"os"
)

// CodeSymbol is one symbol: function, type (as "class"), or module-level constant-like name.
type CodeSymbol struct {
	Name            string   `json:"name"`
	Kind            string   `json:"kind"`
	LineStart       int      `json:"line_start"`
	LineEnd         int      `json:"line_end"`
	LineCode        string   `json:"line_code"`
	FilePath        string   `json:"file_path"`
	Address         string   `json:"address"`
	CallsTo         []string `json:"calls_to"`
	CalledBy        []string `json:"called_by"`
	ConstantValue   *string  `json:"constant_value"`
	Docstring       *string  `json:"docstring"`
	LeadingComment  *string  `json:"leading_comment"`
}

// ParameterInfo describes one formal parameter.
type ParameterInfo struct {
	Name         string  `json:"name"`
	Line         int     `json:"line"`
	Annotation   *string `json:"annotation"`
	DefaultValue *string `json:"default_value"`
}

// NameUsageSite is one identifier occurrence (for future RepoUsageReport parity).
type NameUsageSite struct {
	FilePath        string  `json:"file_path"`
	Line            int     `json:"line"`
	Column          int     `json:"column"`
	Name            string  `json:"name"`
	Kind            string  `json:"kind"`
	ResolvedAddress *string `json:"resolved_address"`
}

// ListedFunction is a flat function listing entry.
type ListedFunction struct {
	Name              string          `json:"name"`
	QualifiedName     string          `json:"qualified_name"`
	FilePath          string          `json:"file_path"`
	Line              int             `json:"line"`
	Parameters        []ParameterInfo `json:"parameters"`
	ReturnAnnotation  *string         `json:"return_annotation"`
	Docstring         *string         `json:"docstring"`
	LeadingComment    *string         `json:"leading_comment"`
}

// ListedClass lists a type (struct/interface) as a "class" for schema compatibility.
type ListedClass struct {
	Name           string  `json:"name"`
	QualifiedName  string  `json:"qualified_name"`
	FilePath       string  `json:"file_path"`
	Line           int     `json:"line"`
	Docstring      *string `json:"docstring"`
	LeadingComment *string `json:"leading_comment"`
}

// ListedConstant is a const (or ALL_CAPS var) at package scope.
type ListedConstant struct {
	Name          string `json:"name"`
	QualifiedName string `json:"qualified_name"`
	FilePath      string `json:"file_path"`
	Line          int    `json:"line"`
	Value         string `json:"value"`
}

// RepoSymbolListing matches Python RepoSymbolListing.
type RepoSymbolListing struct {
	Functions []ListedFunction `json:"functions"`
	Classes   []ListedClass    `json:"classes"`
	Constants []ListedConstant `json:"constants"`
}

// RepoUsageReport matches Python RepoUsageReport.
type RepoUsageReport struct {
	Functions []ListedFunction `json:"functions"`
	Classes   []ListedClass    `json:"classes"`
	Constants []ListedConstant `json:"constants"`
	Usages    []NameUsageSite  `json:"usages"`
}

// CallEdge is a resolved call from caller to callee at a line (1-based).
type CallEdge struct {
	CallerAddress string `json:"caller_address"`
	CalleeAddress string `json:"callee_address"`
	CallLine      int    `json:"call_line"`
}

// CodeBase is the full graph payload (symbols + calls), same as Python CodeBase.
type CodeBase struct {
	Symbols []CodeSymbol `json:"symbols"`
	Calls   []CallEdge   `json:"calls"`
}

// SaveToJSONFileCodeBase writes codeBase to filePath with UTF-8 and indentation (like Python).
func SaveToJSONFileCodeBase(cb *CodeBase, filePath string) error {
	b, err := json.MarshalIndent(cb, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, b, 0o644)
}
