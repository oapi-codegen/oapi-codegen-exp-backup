// Package precedence tests operation-level parameters overriding path-level parameters.
// https://github.com/oapi-codegen/oapi-codegen/issues/1180
package precedence

//go:generate go run ../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
