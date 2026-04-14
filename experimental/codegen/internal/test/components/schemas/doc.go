// Package schemas tests comprehensive schema generation including generic objects,
// nullable properties, custom formats, extra-tags, deprecated fields, and x-go-type-name.
package schemas

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
