// Package all_styles tests parameter type generation across all locations and styles.
// This monolithic spec covers path (simple/label/matrix), query (form/deepObject),
// header, cookie, and content-based parameters.
package all_styles

//go:generate go run ../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
