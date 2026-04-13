// Package recursive tests that recursive types are handled properly.
package recursive

//go:generate go run ../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
