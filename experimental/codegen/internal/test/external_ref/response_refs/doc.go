// Package response_refs tests external response refs across specs.
// https://github.com/oapi-codegen/oapi-codegen/issues/1182
package response_refs

//go:generate go run ../../../../../cmd/oapi-codegen -config pkg2/config.yaml pkg2.yaml
//go:generate go run ../../../../../cmd/oapi-codegen -config pkg1/config.yaml pkg1.yaml
