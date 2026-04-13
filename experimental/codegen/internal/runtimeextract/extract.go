// Package runtimeextract reads Go source files from the embedded runtime FS
// and extracts code bodies for inlining into generated output or assembling
// into standalone runtime packages.
package runtimeextract

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// Import represents a Go import with optional alias.
type Import struct {
	Path  string
	Alias string
}

// RuntimeModulePrefix is the import path prefix for the internal runtime sub-packages.
// Used by the extractor to strip internal imports when inlining, and by the runtime
// generator to rewrite them to the target base path.
const RuntimeModulePrefix = "github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/runtime/"

// qualifierReplacements maps qualified references in params code (which imports
// the types sub-package) to their unqualified forms for inlining.
var qualifierReplacements = []struct{ old, new string }{
	{"types.DateFormat", "DateFormat"},
	{"types.Date", "Date"},
	{"types.Email", "Email"},
	{"types.UUID", "UUID"},
	{"types.File", "File"},
	{"types.Nullable", "Nullable"},
	{"types.NewNullableWithValue", "NewNullableWithValue"},
	{"types.NewNullNullable", "NewNullNullable"},
	{"types.ErrNullableIsNull", "ErrNullableIsNull"},
	{"types.ErrNullableNotSpecified", "ErrNullableNotSpecified"},
	{"types.ErrValidationEmail", "ErrValidationEmail"},
}

// ExtractPackage reads all .go files from a sub-directory of the given FS
// that contain an //oapi-runtime:function annotation and returns the
// concatenated code bodies and merged imports. No qualifier substitution
// is performed — the output is suitable for a standalone runtime package.
func ExtractPackage(fsys fs.FS, dir string) (code string, imports []Import, err error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return "", nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	importSet := make(map[Import]bool)
	var codeParts []string

	// Process files in sorted order for deterministic output.
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		filePath := path.Join(dir, entry.Name())
		data, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return "", nil, fmt.Errorf("reading %s: %w", filePath, err)
		}

		content := string(data)

		// Skip files without the runtime function annotation.
		if !strings.Contains(content, "//oapi-runtime:function ") {
			continue
		}

		fileImports, body, err := parseGoFile(filePath, data)
		if err != nil {
			return "", nil, fmt.Errorf("parsing %s: %w", filePath, err)
		}

		for _, imp := range fileImports {
			importSet[imp] = true
		}
		if body != "" {
			codeParts = append(codeParts, body)
		}
	}

	// Collect and sort imports.
	imports = make([]Import, 0, len(importSet))
	for imp := range importSet {
		imports = append(imports, imp)
	}
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	return strings.Join(codeParts, "\n"), imports, nil
}

// ExtractAllInline reads ALL runtime .go files (types/*, params/*, helpers/*),
// strips package qualifiers for inlining (types.Date → Date, etc.), merges
// imports, and returns a single code body ready to be inserted into a
// generated file. Internal runtime import paths are removed from the import
// list.
//
// The returned code is wrapped in marker comments so the DCE pass can
// identify which declarations are runtime candidates.
func ExtractAllInline(fsys fs.FS) (code string, imports []Import, err error) {
	// Order matters: types first (they define types used by params), then
	// params, then helpers. This ensures declarations appear before use.
	dirs := []string{"types", "params", "helpers"}

	importSet := make(map[Import]bool)
	var codeParts []string

	for _, dir := range dirs {
		entries, err := fs.ReadDir(fsys, dir)
		if err != nil {
			return "", nil, fmt.Errorf("reading directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				continue
			}
			if strings.HasSuffix(entry.Name(), "_test.go") {
				continue
			}

			filePath := path.Join(dir, entry.Name())
			data, err := fs.ReadFile(fsys, filePath)
			if err != nil {
				return "", nil, fmt.Errorf("reading %s: %w", filePath, err)
			}

			content := string(data)

			// Skip files without the runtime function annotation.
			if !strings.Contains(content, "//oapi-runtime:function ") {
				continue
			}

			fileImports, body, err := parseGoFile(filePath, data)
			if err != nil {
				return "", nil, fmt.Errorf("parsing %s: %w", filePath, err)
			}

			// Strip type qualifiers for inlining.
			body = dequalifyTypes(body)

			for _, imp := range fileImports {
				// Skip internal runtime imports — when inlined,
				// all code is in the same package.
				if strings.HasPrefix(imp.Path, RuntimeModulePrefix) {
					continue
				}
				importSet[imp] = true
			}
			if body != "" {
				codeParts = append(codeParts, body)
			}
		}
	}

	// Collect and sort imports.
	imports = make([]Import, 0, len(importSet))
	for imp := range importSet {
		imports = append(imports, imp)
	}
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	// Wrap in markers for DCE.
	joined := strings.Join(codeParts, "\n")
	wrapped := "// --- oapi-runtime begin ---\n" + joined + "\n// --- oapi-runtime end ---\n"

	return wrapped, imports, nil
}

// parseGoFile parses a Go source file and returns its imports and code body
// (everything after the import block, minus the //oapi-runtime: annotations).
func parseGoFile(filePath string, data []byte) ([]Import, string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, data, parser.ImportsOnly)
	if err != nil {
		return nil, "", fmt.Errorf("parsing imports: %w", err)
	}

	// Collect imports.
	var imports []Import
	for _, imp := range f.Imports {
		impPath := strings.Trim(imp.Path.Value, `"`)
		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		imports = append(imports, Import{Path: impPath, Alias: alias})
	}

	// Find where the code body starts (after all imports and package/import declarations).
	// We scan for the last import's closing position, then take everything after.
	bodyStart := findBodyStart(data)
	if bodyStart >= len(data) {
		return imports, "", nil
	}

	body := string(data[bodyStart:])

	// Strip //oapi-runtime: annotation lines from the body.
	body = stripAnnotations(body)

	return imports, strings.TrimSpace(body), nil
}

// findBodyStart returns the byte offset where the code body begins,
// after the package clause and any import declarations.
func findBodyStart(data []byte) int {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	offset := 0
	inImportBlock := false
	foundImport := false
	lastImportEnd := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineEnd := offset + len(line) + 1 // +1 for newline

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "import (") {
			inImportBlock = true
			foundImport = true
		} else if inImportBlock && trimmed == ")" {
			inImportBlock = false
			lastImportEnd = lineEnd
		} else if strings.HasPrefix(trimmed, "import ") && !inImportBlock {
			// Single-line import.
			foundImport = true
			lastImportEnd = lineEnd
		}

		offset = lineEnd
	}

	if foundImport {
		return lastImportEnd
	}

	// No imports — find end of package clause.
	scanner = bufio.NewScanner(strings.NewReader(string(data)))
	offset = 0
	for scanner.Scan() {
		line := scanner.Text()
		lineEnd := offset + len(line) + 1
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			return lineEnd
		}
		offset = lineEnd
	}

	return 0
}

// stripAnnotations removes //oapi-runtime: lines from the code body.
func stripAnnotations(code string) string {
	var lines []string
	for _, line := range strings.Split(code, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//oapi-runtime:") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// dequalifyTypes replaces qualified type references (types.Date, etc.)
// with unqualified forms for inlining into a single package.
func dequalifyTypes(code string) string {
	for _, r := range qualifierReplacements {
		code = strings.ReplaceAll(code, r.old, r.new)
	}
	return code
}
