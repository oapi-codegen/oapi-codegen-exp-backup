// Package overlays tests spec overlays and external refs.
// https://github.com/oapi-codegen/oapi-codegen/issues/1825
package overlays

//go:generate go run ../../../../../cmd/oapi-codegen -config packageA/config.yaml packageA/spec.yaml
//go:generate go run ../../../../../cmd/oapi-codegen -config spec/config.yaml spec/spec.yaml
