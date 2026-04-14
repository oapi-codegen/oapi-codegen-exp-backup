// Package custom_json tests custom JSON content types (application/test+json).
// https://github.com/oapi-codegen/oapi-codegen/issues/1208
// https://github.com/oapi-codegen/oapi-codegen/issues/1209
package custom_json

//go:generate go run ../../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
