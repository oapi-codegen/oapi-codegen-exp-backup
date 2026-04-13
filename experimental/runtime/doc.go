// Package runtime provides shared helper types and functions for code generated
// by oapi-codegen. The source of truth for this code lives in internal/runtime/;
// the files here are generated from that source using GenerateRuntime.
//
// Sub-packages:
//   - types/   — custom Go types for OpenAPI format mappings (Date, Email, UUID, File, Nullable)
//   - params/  — parameter serialization/deserialization functions
//   - helpers/ — utility functions for request body encoding (MarshalForm, JSONMerge)
//
//go:generate go run ../cmd/oapi-codegen --generate-runtime github.com/oapi-codegen/oapi-codegen-exp/experimental/runtime
package runtime
