package codegen

import (
	"text/template"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/templates"
)

// RuntimePrefixes holds the package-qualifier prefixes for the three runtime sub-packages.
// When embedded (no runtime), all fields are empty strings.
type RuntimePrefixes struct {
	Params  string // "params." or ""
	Types   string // "types." or ""
	Helpers string // "helpers." or ""
}

// FuncMap returns a template.FuncMap that exposes runtime prefix accessors to templates.
func (rp RuntimePrefixes) FuncMap() template.FuncMap {
	return template.FuncMap{
		"runtimeParamsPrefix":  func() string { return rp.Params },
		"runtimeTypesPrefix":   func() string { return rp.Types },
		"runtimeHelpersPrefix": func() string { return rp.Helpers },
	}
}

// CodegenContext is a centralized tracker for imports needed during code generation.
// Code at any depth can call its registration methods; the final output assembly
// queries it to collect all required imports.
type CodegenContext struct {
	imports map[string]string // path -> alias

	// Runtime sub-package prefixes — non-empty when a runtime package is configured.
	runtimeParamsPrefix  string // "params." or ""
	runtimeTypesPrefix   string // "types." or ""
	runtimeHelpersPrefix string // "helpers." or ""
}

// NewCodegenContext creates a new CodegenContext.
func NewCodegenContext() *CodegenContext {
	return &CodegenContext{
		imports: make(map[string]string),
	}
}

// --- Runtime prefixes ---

// SetRuntimePrefixes sets the package prefixes for the three runtime sub-packages.
// When non-empty, generated code references runtime helpers via these prefixes
// (e.g., "params.", "types.", "helpers.") instead of embedding them.
func (c *CodegenContext) SetRuntimePrefixes(params, types, helpers string) {
	c.runtimeParamsPrefix = params
	c.runtimeTypesPrefix = types
	c.runtimeHelpersPrefix = helpers
}

// RuntimeParamsPrefix returns the params sub-package prefix (e.g., "params.").
func (c *CodegenContext) RuntimeParamsPrefix() string {
	return c.runtimeParamsPrefix
}

// RuntimeTypesPrefix returns the types sub-package prefix (e.g., "types.").
func (c *CodegenContext) RuntimeTypesPrefix() string {
	return c.runtimeTypesPrefix
}

// RuntimeHelpersPrefix returns the helpers sub-package prefix (e.g., "helpers.").
func (c *CodegenContext) RuntimeHelpersPrefix() string {
	return c.runtimeHelpersPrefix
}

// HasRuntimePackage returns true when an external runtime package is configured.
func (c *CodegenContext) HasRuntimePackage() bool {
	return c.runtimeTypesPrefix != ""
}

// --- Import registration ---

// AddImport records an import path needed by the generated code.
func (c *CodegenContext) AddImport(path string) {
	if path != "" {
		c.imports[path] = ""
	}
}

// AddImportAlias records an import path with an alias.
func (c *CodegenContext) AddImportAlias(path, alias string) {
	if path != "" {
		c.imports[path] = alias
	}
}

// AddImports adds multiple imports from a map[path]alias.
func (c *CodegenContext) AddImports(imports map[string]string) {
	for path, alias := range imports {
		c.AddImportAlias(path, alias)
	}
}

// --- Query methods ---

// Imports returns the collected imports as a map[path]alias.
func (c *CodegenContext) Imports() map[string]string {
	return c.imports
}

// --- Convenience methods (mirror TypeGenerator's API for easy migration) ---

// AddJSONImport adds encoding/json import.
func (c *CodegenContext) AddJSONImport() {
	c.AddImport("encoding/json")
}

// AddJSONImports adds encoding/json and fmt imports.
func (c *CodegenContext) AddJSONImports() {
	c.AddImport("encoding/json")
	c.AddImport("fmt")
}

// AddTemplateImports adds all imports declared by the given template import slices.
func (c *CodegenContext) AddTemplateImports(imports []templates.Import) {
	for _, imp := range imports {
		c.AddImportAlias(imp.Path, imp.Alias)
	}
}
