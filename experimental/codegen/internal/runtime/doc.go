// Package runtime contains the source-of-truth Go source files for the
// oapi-codegen runtime helpers. These files are embedded by the codegen
// binary and either inlined into generated output (with dead code
// elimination) or used to generate the public runtime/ package.
//
// This package is internal — external consumers import the generated
// runtime/ package instead.
package runtime
