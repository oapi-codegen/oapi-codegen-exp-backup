// Package skip_pointer tests x-go-type with skip-optional-pointer and x-go-type-import.
package skip_pointer

//go:generate go run ../../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
