// Package imports tests external dependencies with import resolution.
// https://github.com/oapi-codegen/oapi-codegen/issues/1087
package imports

//go:generate go run ../../../../../cmd/oapi-codegen -config deps/config.yaml deps/my-deps.json
//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
