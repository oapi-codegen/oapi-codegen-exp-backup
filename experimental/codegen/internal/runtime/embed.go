package runtime

import "embed"

// SourceFS contains the Go source files for the runtime sub-packages.
// The codegen binary reads these to extract runtime code for inlining
// into generated files. This package is under internal/ so external
// users importing the public runtime sub-packages don't pay the cost
// of embedding the source files.
//
//go:embed types/*.go params/*.go helpers/*.go
var SourceFS embed.FS
