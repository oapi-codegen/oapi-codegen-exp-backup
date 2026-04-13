// Package multiple tests multiple content types in responses.
// https://github.com/oapi-codegen/oapi-codegen/issues/1127
package multiple

//go:generate go run ../../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
