// Package composition tests complex component schemas including additionalProperties,
// oneOf/anyOf patterns, enums, readOnly/writeOnly, and x-go-name.
package composition

//go:generate go run ../../../../../cmd/oapi-codegen -package output -output output/types.gen.go spec.yaml
