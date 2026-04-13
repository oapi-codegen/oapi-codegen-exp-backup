// Package dce implements dead code elimination for generated Go source files.
// It removes top-level declarations from the runtime section that are not
// reachable from non-runtime code.
package dce

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

const (
	runtimeBeginMarker = "// --- oapi-runtime begin ---"
	runtimeEndMarker   = "// --- oapi-runtime end ---"
)

// EliminateDeadCode parses a generated Go source file, identifies runtime
// declarations (between the oapi-runtime markers), and removes any that are
// not reachable from non-runtime code.
//
// The output is re-printed via go/printer, so callers should run goimports
// afterward to normalize formatting.
func EliminateDeadCode(src string) (string, error) {
	beginIdx := strings.Index(src, runtimeBeginMarker)
	endIdx := strings.Index(src, runtimeEndMarker)
	if beginIdx == -1 || endIdx == -1 {
		return src, nil
	}

	// Replace markers with spaces so byte offsets are preserved for AST.
	cleanSrc := strings.Replace(src, runtimeBeginMarker, strings.Repeat(" ", len(runtimeBeginMarker)), 1)
	cleanSrc = strings.Replace(cleanSrc, runtimeEndMarker, strings.Repeat(" ", len(runtimeEndMarker)), 1)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "generated.go", cleanSrc, parser.ParseComments)
	if err != nil {
		return "", err
	}

	runtimeStart := beginIdx

	// Partition declarations.
	var roots []ast.Decl
	type candidate struct {
		names []string
		decl  ast.Decl
	}
	var candidates []candidate

	for _, decl := range f.Decls {
		offset := fset.Position(decl.Pos()).Offset
		if offset >= runtimeStart {
			var names []string
			switch d := decl.(type) {
			case *ast.GenDecl:
				names = genDeclNames(d)
			case *ast.FuncDecl:
				name := d.Name.Name
				if d.Recv != nil && len(d.Recv.List) > 0 {
					name = receiverTypeName(d.Recv.List[0].Type) + "." + name
				}
				names = []string{name}
			}
			candidates = append(candidates, candidate{names: names, decl: decl})
		} else {
			roots = append(roots, decl)
		}
	}

	if len(candidates) == 0 {
		return src, nil
	}

	// Seed reachable set from root declarations.
	reachable := make(map[string]bool)
	for _, d := range roots {
		collectIdents(d, reachable)
	}

	// Transitive closure.
	changed := true
	for changed {
		changed = false
		for _, c := range candidates {
			if isReachable(c.names, reachable) {
				before := len(reachable)
				collectIdents(c.decl, reachable)
				if len(reachable) > before {
					changed = true
				}
			}
		}
	}

	// Keep only reachable declarations.
	var kept []ast.Decl
	for _, d := range roots {
		kept = append(kept, d)
	}
	for _, c := range candidates {
		if isReachable(c.names, reachable) {
			kept = append(kept, c.decl)
		}
	}
	f.Decls = kept

	var buf strings.Builder
	if err := printer.Fprint(&buf, fset, f); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func isReachable(names []string, reachable map[string]bool) bool {
	for _, name := range names {
		if reachable[name] {
			return true
		}
		if before, _, ok := strings.Cut(name, "."); ok && reachable[before] {
			return true
		}
	}
	return false
}

func genDeclNames(d *ast.GenDecl) []string {
	var names []string
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			names = append(names, s.Name.Name)
		case *ast.ValueSpec:
			for _, n := range s.Names {
				names = append(names, n.Name)
			}
		}
	}
	return names
}

func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return receiverTypeName(t.X)
	case *ast.IndexListExpr:
		return receiverTypeName(t.X)
	}
	return ""
}

func collectIdents(node ast.Node, idents map[string]bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			idents[id.Name] = true
		}
		return true
	})
}
