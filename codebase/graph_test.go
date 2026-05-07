package codebase

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestSymbolAddress verifies posix relative path + qualified name in address strings.
func TestSymbolAddress(t *testing.T) {
	repo := t.TempDir()
	sub := filepath.Join(repo, "a", "b.go")
	if err := os.MkdirAll(filepath.Dir(sub), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sub, []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name      string
		repoRoot  string
		absFile   string
		qual      string
		want      string
		expectErr bool
	}{
		{
			name:     "nested file",
			repoRoot: repo,
			absFile:  sub,
			qual:     "Foo.Bar",
			want:     "a/b.go::Foo.Bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SymbolAddress(tt.repoRoot, tt.absFile, tt.qual)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("SymbolAddress(...) = %q; want %q", got, tt.want)
			}
		})
	}
}

// TestSaveToJSONFileCodeBase checks JSON round-trip and file creation.
func TestSaveToJSONFileCodeBase(t *testing.T) {
	tests := []struct {
		name     string
		input    *CodeBase
		wantKeys []string
	}{
		{
			name: "minimal",
			input: &CodeBase{
				Symbols: []CodeSymbol{{Name: "X", Kind: "function", Address: "f.go::X"}},
				Calls:   []CallEdge{{CallerAddress: "f.go::A", CalleeAddress: "f.go::B", CallLine: 3}},
			},
			wantKeys: []string{"symbols", "calls"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "out.json")
			if err := SaveToJSONFileCodeBase(tt.input, p); err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			var m map[string]json.RawMessage
			if err := json.Unmarshal(raw, &m); err != nil {
				t.Fatal(err)
			}
			for _, k := range tt.wantKeys {
				if _, ok := m[k]; !ok {
					t.Errorf("missing top-level key %q", k)
				}
			}
		})
	}
}

// TestBuildCodebaseForFiles validates symbol extraction and call resolution for the testdata/pack package.
func TestBuildCodebaseForFiles(t *testing.T) {
	repoRoot := findRepoRootForTest(t)
	packDir := filepath.Join(repoRoot, "codebase", "testdata", "pack")

	tests := []struct {
		name       string
		roots      []string
		wantCallee string // substring that must appear as callee in some call edge
		wantCaller string // substring for caller in that edge
	}{
		{
			name:       "same-file and cross-file calls",
			roots:      []string{packDir},
			wantCallee: "::helper",
			wantCaller: "::Exported",
		},
		{
			name:       "cross-file same package",
			roots:      []string{packDir},
			wantCallee: "::Exported",
			wantCaller: "::FromB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb, err := BuildCodebaseForFiles(tt.roots, BuildOptions{RepoRoot: repoRoot})
			if err != nil {
				t.Fatal(err)
			}
			var found bool
			for _, e := range cb.Calls {
				if strings.Contains(e.CalleeAddress, tt.wantCallee) && strings.Contains(e.CallerAddress, tt.wantCaller) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("no call edge caller %q -> callee %q; calls=%v", tt.wantCaller, tt.wantCallee, cb.Calls)
			}
		})
	}

	t.Run("symbols include Exported and type", func(t *testing.T) {
		cb, err := BuildCodebaseForFiles([]string{packDir}, BuildOptions{RepoRoot: repoRoot})
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for _, s := range cb.Symbols {
			if s.Kind == "function" {
				names = append(names, s.Name)
			}
		}
		for _, need := range []string{"Exported", "FromB", "helper"} {
			if !slices.Contains(names, need) {
				t.Errorf("missing function symbol %q, have %v", need, names)
			}
		}
	})
}

func findRepoRootForTest(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	r, err := findRepoRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
