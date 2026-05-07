package codebase

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

var (
	goParser     *sitter.Parser
	goParserInit sync.Once
	goParserErr  error
)

func parser() (*sitter.Parser, error) {
	goParserInit.Do(func() {
		p := sitter.NewParser()
		lang := sitter.NewLanguage(tree_sitter_go.Language())
		if err := p.SetLanguage(lang); err != nil {
			goParserErr = err
			p.Close()
			return
		}
		goParser = p
	})
	if goParserErr != nil {
		return nil, goParserErr
	}
	return goParser, nil
}

// BuildOptions configures BuildCodebaseForFiles.
type BuildOptions struct {
	// RepoRoot is the module root (directory containing go.mod). Relative paths in JSON are posix paths relative to this.
	RepoRoot string
	// Ignore is absolute or relative file/dir paths; files under them are skipped (same idea as Python ignore=).
	Ignore []string
}

type parsedFile struct {
	abs string
	rel string
	src []byte
	tr  *sitter.Tree
}

func (p *parsedFile) close() {
	if p != nil && p.tr != nil {
		p.tr.Close()
		p.tr = nil
	}
}

func nodeText(src []byte, n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return string(src[n.StartByte():n.EndByte()])
}

func lineStart1(n *sitter.Node) int {
	if n == nil {
		return 0
	}
	return int(n.StartPosition().Row) + 1
}

func lineEnd1(n *sitter.Node) int {
	if n == nil {
		return 0
	}
	return int(n.EndPosition().Row) + 1
}

func lineSnippet(src []byte, n *sitter.Node) string {
	if n == nil {
		return ""
	}
	chunk := nodeText(src, n)
	first := strings.TrimSpace(strings.Split(chunk, "\n")[0])
	if len(first) > 200 {
		return first[:197] + "..."
	}
	return first
}

var constNameGo = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

func shouldIgnorePath(path string, ignore []string) (bool, error) {
	rp, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	for _, raw := range ignore {
		ign, err := filepath.Abs(raw)
		if err != nil {
			return false, err
		}
		if rp == ign {
			return true, nil
		}
		rel, err := filepath.Rel(ign, rp)
		if err != nil {
			continue
		}
		if rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return true, nil
		}
	}
	return false, nil
}

func expandGoRoots(roots []string, ignore []string) ([]string, error) {
	seen := map[string]struct{}{}
	var out []string
	for _, raw := range roots {
		p, err := filepath.Abs(raw)
		if err != nil {
			return nil, err
		}
		fi, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			_ = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if !strings.HasSuffix(strings.ToLower(path), ".go") {
					return nil
				}
				if skip, _ := shouldIgnorePath(path, ignore); skip {
					return nil
				}
				if _, ok := seen[path]; ok {
					return nil
				}
				seen[path] = struct{}{}
				out = append(out, path)
				return nil
			})
		} else {
			if !strings.HasSuffix(strings.ToLower(p), ".go") {
				continue
			}
			if skip, err := shouldIgnorePath(p, ignore); err != nil {
				return nil, err
			} else if skip {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out, nil
}

func readModulePath(repoRoot string) (string, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			mod := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if i := strings.IndexByte(mod, ' '); i >= 0 {
				mod = mod[:i]
			}
			return mod, nil
		}
	}
	return "", fmt.Errorf("no module directive in go.mod under %s", repoRoot)
}

func findRepoRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", start)
		}
		dir = parent
	}
}

// SymbolAddress builds stable address relPath::qualifiedName (posix relPath).
func SymbolAddress(repoRoot, absFile, qualifiedName string) (string, error) {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", err
	}
	f, err := filepath.Abs(absFile)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, f)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)
	return rel + "::" + qualifiedName, nil
}

func toPosixRel(repoRoot, absFile string) (string, error) {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", err
	}
	f, err := filepath.Abs(absFile)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, f)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

var builtinsGo = map[string]struct{}{
	"append": {}, "cap": {}, "close": {}, "complex": {}, "copy": {}, "delete": {},
	"imag": {}, "len": {}, "make": {}, "new": {}, "panic": {}, "print": {}, "println": {},
	"real": {}, "recover": {}, "clear": {}, "min": {}, "max": {},
}

func typeExprShortName(n *sitter.Node, src []byte) string {
	if n == nil {
		return "?"
	}
	switch n.Kind() {
	case "type_identifier":
		return strings.TrimSpace(nodeText(src, n))
	case "field_identifier":
		return strings.TrimSpace(nodeText(src, n))
	case "pointer_type":
		if n.NamedChildCount() > 0 {
			return typeExprShortName(n.NamedChild(0), src)
		}
	case "parenthesized_type":
		if n.NamedChildCount() > 0 {
			return typeExprShortName(n.NamedChild(0), src)
		}
	case "array_type", "slice_type":
		if el := n.ChildByFieldName("element"); el != nil {
			return typeExprShortName(el, src)
		}
	case "qualified_type":
		if name := n.ChildByFieldName("name"); name != nil {
			return typeExprShortName(name, src)
		}
	}
	t := strings.TrimSpace(nodeText(src, n))
	if len(t) > 64 {
		return t[:61] + "..."
	}
	return t
}

func receiverShortName(recv *sitter.Node, src []byte) string {
	if recv == nil || recv.Kind() != "parameter_list" {
		return "?"
	}
	for i := uint(0); i < recv.NamedChildCount(); i++ {
		ch := recv.NamedChild(i)
		if ch == nil {
			continue
		}
		if ch.Kind() != "parameter_declaration" {
			continue
		}
		typ := ch.ChildByFieldName("type")
		return typeExprShortName(typ, src)
	}
	return "?"
}

func godocAbove(anchor *sitter.Node, src []byte) *string {
	p := anchor.Parent()
	if p == nil {
		return nil
	}
	var idx int = -1
	for i := uint(0); i < p.ChildCount(); i++ {
		ch := p.Child(i)
		if ch != nil && ch.Id() == anchor.Id() {
			idx = int(i)
			break
		}
	}
	if idx <= 0 {
		return nil
	}
	var lines []string
	for j := idx - 1; j >= 0; j-- {
		sib := p.Child(uint(j))
		if sib == nil {
			break
		}
		if sib.Kind() != "comment" {
			break
		}
		raw := strings.TrimSpace(nodeText(src, sib))
		if strings.HasPrefix(raw, "//") {
			body := strings.TrimSpace(strings.TrimPrefix(raw, "//"))
			lines = append(lines, body)
		} else {
			lines = append(lines, raw)
		}
	}
	if len(lines) == 0 {
		return nil
	}
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	joined := strings.TrimSpace(strings.Join(lines, "\n"))
	if joined == "" {
		return nil
	}
	if len(joined) > 4000 {
		joined = joined[:3997] + "..."
	}
	return &joined
}

func extractParameters(fn *sitter.Node, src []byte) []ParameterInfo {
	params := fn.ChildByFieldName("parameters")
	if params == nil {
		return nil
	}
	var out []ParameterInfo
	for i := uint(0); i < params.NamedChildCount(); i++ {
		ch := params.NamedChild(i)
		if ch == nil {
			continue
		}
		switch ch.Kind() {
		case "parameter_declaration":
			var nameNode *sitter.Node
			for j := uint(0); j < ch.NamedChildCount(); j++ {
				c := ch.NamedChild(j)
				if c != nil && c.Kind() == "identifier" {
					nameNode = c
					break
				}
			}
			if nameNode == nil {
				continue
			}
			ann := ch.ChildByFieldName("type")
			var annStr *string
			if ann != nil {
				s := strings.TrimSpace(nodeText(src, ann))
				if len(s) > 200 {
					s = s[:197] + "..."
				}
				annStr = &s
			}
			out = append(out, ParameterInfo{
				Name:       strings.TrimSpace(nodeText(src, nameNode)),
				Line:       lineStart1(nameNode),
				Annotation: annStr,
			})
		case "variadic_parameter_declaration":
			nameNode := ch.ChildByFieldName("name")
			if nameNode != nil && nameNode.Kind() == "identifier" {
				out = append(out, ParameterInfo{
					Name: strings.TrimSpace(nodeText(src, nameNode)),
					Line: lineStart1(nameNode),
				})
			}
		}
	}
	return out
}

func resultTypeText(fn *sitter.Node, src []byte) *string {
	rt := fn.ChildByFieldName("result")
	if rt == nil {
		return nil
	}
	s := strings.TrimSpace(nodeText(src, rt))
	if s == "" {
		return nil
	}
	if len(s) > 300 {
		s = s[:297] + "..."
	}
	return &s
}

func importPathString(n *sitter.Node, src []byte) string {
	if n == nil {
		return ""
	}
	raw := strings.TrimSpace(nodeText(src, n))
	raw = strings.Trim(raw, "`")
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		return raw[1 : len(raw)-1]
	}
	return raw
}

type importRow struct {
	path   string // import path string
	local  string // local package name (default base of path)
	dotted string // same as local for from-import parity; optional
}

func collectImports(root *sitter.Node, src []byte) []importRow {
	var rows []importRow
	for i := uint(0); i < root.ChildCount(); i++ {
		ch := root.Child(i)
		if ch == nil || ch.Kind() != "import_declaration" {
			continue
		}
		var walk func(*sitter.Node)
		walk = func(n *sitter.Node) {
			if n == nil {
				return
			}
			switch n.Kind() {
			case "import_spec":
				pathNode := n.ChildByFieldName("path")
				ip := importPathString(pathNode, src)
				if ip == "" {
					return
				}
				local := path.Base(ip)
				if name := n.ChildByFieldName("name"); name != nil {
					switch name.Kind() {
					case "package_identifier":
						local = strings.TrimSpace(nodeText(src, name))
					case "dot":
						local = "."
					case "blank_identifier":
						return
					}
				}
				rows = append(rows, importRow{path: ip, local: local, dotted: ip})
			case "import_spec_list":
				for j := uint(0); j < n.NamedChildCount(); j++ {
					walk(n.NamedChild(j))
				}
			}
		}
		for j := uint(0); j < ch.NamedChildCount(); j++ {
			walk(ch.NamedChild(j))
		}
	}
	return rows
}

func importDirForPath(repoRoot, modulePath, importPath string) (string, bool) {
	if importPath == "" {
		return "", false
	}
	var suffix string
	switch {
	case importPath == modulePath:
		suffix = ""
	case strings.HasPrefix(importPath, modulePath+"/"):
		suffix = strings.TrimPrefix(importPath, modulePath+"/")
	default:
		return "", false
	}
	dir := filepath.Join(repoRoot, filepath.FromSlash(suffix))
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		return "", false
	}
	return dir, true
}

func topLevelNameInFile(tree *sitter.Tree, src []byte, want string) bool {
	root := tree.RootNode()
	if root == nil {
		return false
	}
	for i := uint(0); i < root.ChildCount(); i++ {
		ch := root.Child(i)
		if ch == nil {
			continue
		}
		switch ch.Kind() {
		case "function_declaration":
			if nm := ch.ChildByFieldName("name"); nm != nil && nodeText(src, nm) == want {
				return true
			}
		case "type_declaration":
			for j := uint(0); j < ch.NamedChildCount(); j++ {
				sub := ch.NamedChild(j)
				if sub == nil {
					continue
				}
				if sub.Kind() == "type_spec" {
					if nm := sub.ChildByFieldName("name"); nm != nil && nodeText(src, nm) == want {
						return true
					}
				}
			}
		}
	}
	return false
}

func findDefInPackageDir(dir, simple string, cache map[string]*parsedFile, p *sitter.Parser) (absFile, qual string, ok bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", "", false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		pf, ok2 := cache[path]
		if !ok2 {
			b, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			tr := p.Parse(b, nil)
			if tr == nil {
				continue
			}
			pf = &parsedFile{abs: path, src: b, tr: tr}
			cache[path] = pf
		}
		if topLevelNameInFile(pf.tr, pf.src, simple) {
			return path, simple, true
		}
	}
	return "", "", false
}

func sameFileIndex(syms []CodeSymbol, relFile string) map[string][][2]string {
	m := map[string][][2]string{}
	for _, s := range syms {
		if s.FilePath != relFile {
			continue
		}
		if s.Kind != "function" && s.Kind != "class" {
			continue
		}
		line := fmt.Sprintf("%d", s.LineStart)
		m[s.Name] = append(m[s.Name], [2]string{line, s.Address})
	}
	for k, v := range m {
		sort.Slice(v, func(i, j int) bool {
			return v[i][0] < v[j][0]
		})
		m[k] = v
	}
	return m
}

// packageTopLevelIndex maps absolute package directory -> short name -> (line, address) for top-level
// funcs and types only (qualified name has no '.'), so method names do not collide across receivers.
func packageTopLevelIndex(syms []CodeSymbol, repoRoot string) map[string]map[string][][2]string {
	out := map[string]map[string][][2]string{}
	for _, s := range syms {
		if s.Kind != "function" && s.Kind != "class" {
			continue
		}
		_, qual, ok := strings.Cut(s.Address, "::")
		if !ok || strings.Contains(qual, ".") {
			continue
		}
		dir := filepath.Join(repoRoot, filepath.FromSlash(path.Dir(s.FilePath)))
		if out[dir] == nil {
			out[dir] = map[string][][2]string{}
		}
		line := fmt.Sprintf("%d", s.LineStart)
		out[dir][s.Name] = append(out[dir][s.Name], [2]string{line, s.Address})
	}
	for _, m := range out {
		for k, v := range m {
			sort.Slice(v, func(i, j int) bool {
				return v[i][0] < v[j][0]
			})
			m[k] = v
		}
	}
	return out
}

func sameFileBestAddress(idx map[string][][2]string, name string, usageLine int) string {
	pairs, ok := idx[name]
	if !ok {
		return ""
	}
	var eligible []string
	for _, p := range pairs {
		var ln int
		fmt.Sscanf(p[0], "%d", &ln)
		if ln < usageLine {
			eligible = append(eligible, p[1])
		}
	}
	if len(eligible) == 0 {
		return ""
	}
	return eligible[len(eligible)-1]
}

// packageTopLevelAddress resolves a top-level symbol name within a package directory.
// Unlike same-file resolution, line numbers from different source files are not comparable,
// so we only support the Go rule of one top-level name per package (or pick the last if duplicated).
func packageTopLevelAddress(m map[string][][2]string, name string) string {
	pairs, ok := m[name]
	if !ok || len(pairs) == 0 {
		return ""
	}
	return pairs[len(pairs)-1][1]
}

func calleeFromCall(call *sitter.Node, src []byte) (callee string, pkgAlias string) {
	if call == nil || call.Kind() != "call_expression" {
		return "", ""
	}
	fn := call.ChildByFieldName("function")
	fn = unwrapPrimary(fn)
	if fn == nil {
		return "", ""
	}
	if fn.Kind() == "identifier" {
		return strings.TrimSpace(nodeText(src, fn)), ""
	}
	if fn.Kind() == "selector_expression" {
		fd := fn.ChildByFieldName("field")
		if fd == nil {
			return "", ""
		}
		callee = strings.TrimSpace(nodeText(src, fd))
		if op := fn.ChildByFieldName("operand"); op != nil && op.Kind() == "identifier" {
			pkgAlias = strings.TrimSpace(nodeText(src, op))
		}
		return callee, pkgAlias
	}
	return "", ""
}

func unwrapPrimary(n *sitter.Node) *sitter.Node {
	for n != nil && n.Kind() == "parenthesized_expression" {
		if n.NamedChildCount() > 0 {
			n = n.NamedChild(0)
		} else {
			break
		}
	}
	return n
}

type rawCall struct {
	callee   string
	pkgAlias string
	line     int
}

func collectCallsSkipNested(node *sitter.Node, src []byte) []rawCall {
	var out []rawCall
	var walk func(*sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Kind() == "function_declaration" || n.Kind() == "func_literal" {
			return
		}
		if n.Kind() == "call_expression" {
			callee, pkg := calleeFromCall(n, src)
			if callee != "" {
				out = append(out, rawCall{callee: callee, pkgAlias: pkg, line: lineStart1(n)})
			}
		}
		for i := uint(0); i < n.ChildCount(); i++ {
			walk(n.Child(i))
		}
	}
	walk(node)
	return out
}

func collectCallsInNode(node *sitter.Node, src []byte) []rawCall {
	var out []rawCall
	var walk func(*sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Kind() == "call_expression" {
			callee, pkg := calleeFromCall(n, src)
			if callee != "" {
				out = append(out, rawCall{callee: callee, pkgAlias: pkg, line: lineStart1(n)})
			}
		}
		for i := uint(0); i < n.ChildCount(); i++ {
			walk(n.Child(i))
		}
	}
	walk(node)
	return out
}

func resolveCallee(
	repoRoot, modulePath, importerAbs string,
	callee, pkgAlias string,
	imports []importRow,
	usageLine int,
	fileIdx map[string][][2]string,
	pkgIdx map[string]map[string][][2]string,
	cache map[string]*parsedFile,
	p *sitter.Parser,
) string {
	if _, ok := builtinsGo[callee]; ok {
		return ""
	}
	if usageLine > 0 && fileIdx != nil {
		if addr := sameFileBestAddress(fileIdx, callee, usageLine); addr != "" {
			return addr
		}
	}
	pkgDir := filepath.Dir(importerAbs)
	if pkgIdx != nil {
		if m := pkgIdx[pkgDir]; m != nil {
			if addr := packageTopLevelAddress(m, callee); addr != "" {
				return addr
			}
		}
	}
	if pkgAlias != "" {
		for _, row := range imports {
			if row.local != pkgAlias {
				continue
			}
			dir, ok := importDirForPath(repoRoot, modulePath, row.path)
			if !ok {
				continue
			}
			abs, qual, ok := findDefInPackageDir(dir, callee, cache, p)
			if ok {
				rel, _ := toPosixRel(repoRoot, abs)
				return rel + "::" + qual
			}
		}
	}
	return ""
}

type funcBody struct {
	qual string
	node *sitter.Node
}

func collectPackageLevel(tree *sitter.Tree, src []byte, abs, rel string, repoRoot string, outSyms *[]CodeSymbol, bodies *[]funcBody) error {
	cur := tree.Walk()
	defer cur.Close()
	root := tree.RootNode()
	if root == nil {
		return nil
	}
	for i := uint(0); i < root.ChildCount(); i++ {
		st := root.Child(i)
		if st == nil {
			continue
		}
		switch st.Kind() {
		case "function_declaration":
			name := st.ChildByFieldName("name")
			if name == nil {
				continue
			}
			nm := strings.TrimSpace(nodeText(src, name))
			addr, err := SymbolAddress(repoRoot, abs, nm)
			if err != nil {
				return err
			}
			doc := godocAbove(st, src)
			*outSyms = append(*outSyms, CodeSymbol{
				Name: nm, Kind: "function",
				LineStart: lineStart1(st), LineEnd: lineEnd1(st),
				LineCode: lineSnippet(src, st), FilePath: rel, Address: addr,
				CallsTo: []string{}, CalledBy: []string{},
				Docstring: doc,
			})
			*bodies = append(*bodies, funcBody{qual: nm, node: st})
		case "method_declaration":
			recv := st.ChildByFieldName("receiver")
			rname := receiverShortName(recv, src)
			name := st.ChildByFieldName("name")
			if name == nil {
				continue
			}
			meth := strings.TrimSpace(nodeText(src, name))
			qual := rname + "." + meth
			addr, err := SymbolAddress(repoRoot, abs, qual)
			if err != nil {
				return err
			}
			doc := godocAbove(st, src)
			*outSyms = append(*outSyms, CodeSymbol{
				Name: meth, Kind: "function",
				LineStart: lineStart1(st), LineEnd: lineEnd1(st),
				LineCode: lineSnippet(src, st), FilePath: rel, Address: addr,
				CallsTo: []string{}, CalledBy: []string{},
				Docstring: doc,
			})
			*bodies = append(*bodies, funcBody{qual: qual, node: st})
		case "type_declaration":
			for j := uint(0); j < st.NamedChildCount(); j++ {
				spec := st.NamedChild(j)
				if spec == nil || spec.Kind() != "type_spec" {
					continue
				}
				tdef := spec.ChildByFieldName("type")
				if tdef == nil {
					continue
				}
				k := tdef.Kind()
				if k != "struct_type" && k != "interface_type" {
					continue
				}
				name := spec.ChildByFieldName("name")
				if name == nil {
					continue
				}
				tnm := strings.TrimSpace(nodeText(src, name))
				addr, err := SymbolAddress(repoRoot, abs, tnm)
				if err != nil {
					return err
				}
				doc := godocAbove(st, src)
				*outSyms = append(*outSyms, CodeSymbol{
					Name: tnm, Kind: "class",
					LineStart: lineStart1(spec), LineEnd: lineEnd1(st),
					LineCode: lineSnippet(src, spec), FilePath: rel, Address: addr,
					CallsTo: []string{}, CalledBy: []string{},
					Docstring: doc,
				})
			}
		case "const_declaration":
			for j := uint(0); j < st.NamedChildCount(); j++ {
				sp := st.NamedChild(j)
				if sp == nil || sp.Kind() != "const_spec" {
					continue
				}
				names := sp.ChildrenByFieldName("name", cur)
				for _, nn := range names {
					if nn.Kind() != "identifier" {
						continue
					}
					nm := strings.TrimSpace(nodeText(src, &nn))
					if !constNameGo.MatchString(nm) {
						continue
					}
					addr, err := SymbolAddress(repoRoot, abs, nm)
					if err != nil {
						return err
					}
					val := ""
					if v := sp.ChildByFieldName("value"); v != nil {
						val = strings.TrimSpace(nodeText(src, v))
						if len(val) > 400 {
							val = val[:397] + "..."
						}
					}
					cv := val
					*outSyms = append(*outSyms, CodeSymbol{
						Name: nm, Kind: "constant",
						LineStart: lineStart1(sp), LineEnd: lineEnd1(sp),
						LineCode: lineSnippet(src, sp), FilePath: rel, Address: addr,
						CallsTo: []string{}, CalledBy: []string{},
						ConstantValue: &cv,
					})
				}
			}
		case "var_declaration":
			for j := uint(0); j < st.NamedChildCount(); j++ {
				inner := st.NamedChild(j)
				if inner == nil {
					continue
				}
				var specs []*sitter.Node
				switch inner.Kind() {
				case "var_spec":
					specs = append(specs, inner)
				case "var_spec_list":
					for k := uint(0); k < inner.NamedChildCount(); k++ {
						if ch := inner.NamedChild(k); ch != nil && ch.Kind() == "var_spec" {
							specs = append(specs, ch)
						}
					}
				}
				for _, sp := range specs {
					names := sp.ChildrenByFieldName("name", cur)
					for _, nn := range names {
						if nn.Kind() != "identifier" {
							continue
						}
						nm := strings.TrimSpace(nodeText(src, &nn))
						if !constNameGo.MatchString(nm) {
							continue
						}
						addr, err := SymbolAddress(repoRoot, abs, nm)
						if err != nil {
							return err
						}
						val := ""
						if v := sp.ChildByFieldName("value"); v != nil {
							val = strings.TrimSpace(nodeText(src, v))
							if len(val) > 400 {
								val = val[:397] + "..."
							}
						}
						cv := val
						*outSyms = append(*outSyms, CodeSymbol{
							Name: nm, Kind: "constant",
							LineStart: lineStart1(sp), LineEnd: lineEnd1(sp),
							LineCode: lineSnippet(src, sp), FilePath: rel, Address: addr,
							CallsTo: []string{}, CalledBy: []string{},
							ConstantValue: &cv,
						})
					}
				}
			}
		}
	}
	return nil
}

func edgesFromFunc(repoRoot, modulePath, abs, rel string, fb funcBody, src []byte, imports []importRow, fileIdx map[string][][2]string, pkgIdx map[string]map[string][][2]string, cache map[string]*parsedFile, p *sitter.Parser) ([]CallEdge, error) {
	callerAddr, err := SymbolAddress(repoRoot, abs, fb.qual)
	if err != nil {
		return nil, err
	}
	var edges []CallEdge
	params := fb.node.ChildByFieldName("parameters")
	if params != nil {
		for _, rc := range collectCallsInNode(params, src) {
			if addr := resolveCallee(repoRoot, modulePath, abs, rc.callee, rc.pkgAlias, imports, rc.line, fileIdx, pkgIdx, cache, p); addr != "" {
				edges = append(edges, CallEdge{CallerAddress: callerAddr, CalleeAddress: addr, CallLine: rc.line})
			}
		}
	}
	if rt := fb.node.ChildByFieldName("result"); rt != nil {
		for _, rc := range collectCallsInNode(rt, src) {
			if addr := resolveCallee(repoRoot, modulePath, abs, rc.callee, rc.pkgAlias, imports, rc.line, fileIdx, pkgIdx, cache, p); addr != "" {
				edges = append(edges, CallEdge{CallerAddress: callerAddr, CalleeAddress: addr, CallLine: rc.line})
			}
		}
	}
	if body := fb.node.ChildByFieldName("body"); body != nil {
		for _, rc := range collectCallsSkipNested(body, src) {
			if addr := resolveCallee(repoRoot, modulePath, abs, rc.callee, rc.pkgAlias, imports, rc.line, fileIdx, pkgIdx, cache, p); addr != "" {
				edges = append(edges, CallEdge{CallerAddress: callerAddr, CalleeAddress: addr, CallLine: rc.line})
			}
		}
	}
	return edges, nil
}

// BuildCodebaseForFiles parses .go files (or walks directories) and returns symbols + resolved call edges,
// aligned with Python build_codebase_for_files / CodeBase.
func BuildCodebaseForFiles(paths []string, opt BuildOptions) (*CodeBase, error) {
	pr, err := parser()
	if err != nil {
		return nil, err
	}
	files, err := expandGoRoots(paths, opt.Ignore)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return &CodeBase{Symbols: nil, Calls: nil}, nil
	}
	repoRoot := opt.RepoRoot
	if repoRoot == "" {
		repoRoot, err = findRepoRoot(filepath.Dir(files[0]))
		if err != nil {
			return nil, err
		}
	} else {
		repoRoot, err = filepath.Abs(repoRoot)
		if err != nil {
			return nil, err
		}
	}
	modulePath, err := readModulePath(repoRoot)
	if err != nil {
		return nil, err
	}

	cache := map[string]*parsedFile{}
	defer func() {
		for _, pf := range cache {
			pf.close()
		}
	}()

	for _, abs := range files {
		src, err := os.ReadFile(abs)
		if err != nil {
			return nil, err
		}
		tr := pr.Parse(src, nil)
		if tr == nil {
			return nil, fmt.Errorf("parse failed for %s", abs)
		}
		rel, err := toPosixRel(repoRoot, abs)
		if err != nil {
			tr.Close()
			return nil, err
		}
		cache[abs] = &parsedFile{abs: abs, rel: rel, src: src, tr: tr}
	}

	var allSyms []CodeSymbol
	var allEdges []CallEdge
	perFileBodies := map[string][]funcBody{}

	for abs, pf := range cache {
		var bodies []funcBody
		if err := collectPackageLevel(pf.tr, pf.src, abs, pf.rel, repoRoot, &allSyms, &bodies); err != nil {
			return nil, err
		}
		perFileBodies[abs] = bodies
	}

	pkgIdx := packageTopLevelIndex(allSyms, repoRoot)

	for abs, pf := range cache {
		imports := collectImports(pf.tr.RootNode(), pf.src)
		fileIdx := sameFileIndex(allSyms, pf.rel)
		for _, fb := range perFileBodies[abs] {
			e, err := edgesFromFunc(repoRoot, modulePath, abs, pf.rel, fb, pf.src, imports, fileIdx, pkgIdx, cache, pr)
			if err != nil {
				return nil, err
			}
			allEdges = append(allEdges, e...)
		}
	}

	addrToSym := map[string]*CodeSymbol{}
	for i := range allSyms {
		addrToSym[allSyms[i].Address] = &allSyms[i]
	}
	for i := range allEdges {
		e := &allEdges[i]
		if c, ok := addrToSym[e.CallerAddress]; ok {
			if !containsStr(c.CallsTo, e.CalleeAddress) {
				c.CallsTo = append(c.CallsTo, e.CalleeAddress)
			}
		}
		if t, ok := addrToSym[e.CalleeAddress]; ok {
			if !containsStr(t.CalledBy, e.CallerAddress) {
				t.CalledBy = append(t.CalledBy, e.CallerAddress)
			}
		}
	}
	for i := range allSyms {
		sort.Strings(allSyms[i].CallsTo)
		sort.Strings(allSyms[i].CalledBy)
	}

	return &CodeBase{Symbols: allSyms, Calls: allEdges}, nil
}

func containsStr(sl []string, v string) bool {
	for _, s := range sl {
		if s == v {
			return true
		}
	}
	return false
}
