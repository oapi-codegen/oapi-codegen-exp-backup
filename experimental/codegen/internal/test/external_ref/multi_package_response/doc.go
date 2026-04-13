// Package multi_package_response tests multi-package response schemas.
// https://github.com/oapi-codegen/oapi-codegen/issues/1212
package multi_package_response

//go:generate go run ../../../../../cmd/oapi-codegen -config pkg2/config.yaml pkg2.yaml
//go:generate go run ../../../../../cmd/oapi-codegen -config pkg1/config.yaml pkg1.yaml
