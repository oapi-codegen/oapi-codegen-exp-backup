// Package multi_spec tests multi-spec cross-package imports.
// https://github.com/oapi-codegen/oapi-codegen/issues/1093
package multi_spec

//go:generate go run ../../../../../cmd/oapi-codegen -config parent.config.yaml parent.api.yaml
//go:generate go run ../../../../../cmd/oapi-codegen -config child.config.yaml child.api.yaml
