// Package encoding tests path parameter escaping and special characters.
// https://github.com/oapi-codegen/oapi-codegen/issues/312
package encoding

//go:generate go run ../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
