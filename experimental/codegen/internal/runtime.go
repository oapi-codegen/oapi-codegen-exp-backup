package codegen

import (
	"fmt"
	"strings"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/runtimeextract"
	runtime "github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/runtime"
)

// RuntimeOutput holds the generated Go source code for each runtime sub-package.
type RuntimeOutput struct {
	Params  string // params sub-package (style/bind functions, helpers)
	Types   string // types sub-package (Date, Email, UUID, File, Nullable)
	Helpers string // helpers sub-package (MarshalForm)
}

// GenerateRuntime produces standalone Go source files for each of the three
// runtime sub-packages. baseImportPath is the base import path for the runtime
// module (e.g., "github.com/org/project/runtime"). The params sub-package
// imports the types sub-package for Date references.
func GenerateRuntime(baseImportPath string) (*RuntimeOutput, error) {
	if baseImportPath == "" {
		return nil, fmt.Errorf("base import path is required")
	}

	typesCode, err := generateRuntimePackage("types", "types", baseImportPath)
	if err != nil {
		return nil, fmt.Errorf("generating runtime types: %w", err)
	}

	paramsCode, err := generateRuntimePackage("params", "params", baseImportPath)
	if err != nil {
		return nil, fmt.Errorf("generating runtime params: %w", err)
	}

	helpersCode, err := generateRuntimePackage("helpers", "helpers", baseImportPath)
	if err != nil {
		return nil, fmt.Errorf("generating runtime helpers: %w", err)
	}

	return &RuntimeOutput{
		Params:  paramsCode,
		Types:   typesCode,
		Helpers: helpersCode,
	}, nil
}

// generateRuntimePackage produces a standalone Go source file for one runtime
// sub-package by extracting annotated code from the embedded runtime sources.
func generateRuntimePackage(dir, packageName, baseImportPath string) (string, error) {
	code, imports, err := runtimeextract.ExtractPackage(runtime.SourceFS, dir)
	if err != nil {
		return "", fmt.Errorf("extracting %s: %w", dir, err)
	}

	output := NewOutput(packageName)
	output.AddType(code)

	for _, imp := range imports {
		path := imp.Path
		// Rewrite internal runtime import paths to the target base path.
		// e.g., ".../codegen/internal/runtime/types" → "baseImportPath/types"
		if strings.HasPrefix(path, runtimeextract.RuntimeModulePrefix) {
			path = baseImportPath + "/" + strings.TrimPrefix(path, runtimeextract.RuntimeModulePrefix)
		}
		output.AddImport(path, imp.Alias)
	}

	return output.Format()
}
