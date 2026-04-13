// Package inline_identifiers tests that inline schemas generate valid Go identifiers.
// https://github.com/oapi-codegen/oapi-codegen/issues/1496
package inline_identifiers

//go:generate go run ../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
