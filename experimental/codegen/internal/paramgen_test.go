package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultParamStyle(t *testing.T) {
	tests := []struct {
		location string
		expected string
	}{
		{"path", "simple"},
		{"header", "simple"},
		{"query", "form"},
		{"cookie", "form"},
		{"unknown", "form"},
	}

	for _, tc := range tests {
		t.Run(tc.location, func(t *testing.T) {
			assert.Equal(t, tc.expected, DefaultParamStyle(tc.location))
		})
	}
}

func TestDefaultParamExplode(t *testing.T) {
	tests := []struct {
		location string
		expected bool
	}{
		{"path", false},
		{"header", false},
		{"query", true},
		{"cookie", true},
		{"unknown", false},
	}

	for _, tc := range tests {
		t.Run(tc.location, func(t *testing.T) {
			assert.Equal(t, tc.expected, DefaultParamExplode(tc.location))
		})
	}
}

func TestValidateParamStyle(t *testing.T) {
	validCases := []struct {
		style    string
		location string
	}{
		{"simple", "path"},
		{"label", "path"},
		{"matrix", "path"},
		{"form", "query"},
		{"spaceDelimited", "query"},
		{"pipeDelimited", "query"},
		{"deepObject", "query"},
		{"simple", "header"},
		{"form", "cookie"},
	}

	for _, tc := range validCases {
		t.Run(tc.style+"_in_"+tc.location, func(t *testing.T) {
			err := ValidateParamStyle(tc.style, tc.location)
			assert.NoError(t, err)
		})
	}

	invalidCases := []struct {
		style    string
		location string
	}{
		{"deepObject", "path"},
		{"matrix", "query"},
		{"label", "header"},
		{"simple", "cookie"},
	}

	for _, tc := range invalidCases {
		t.Run(tc.style+"_in_"+tc.location+"_invalid", func(t *testing.T) {
			err := ValidateParamStyle(tc.style, tc.location)
			assert.Error(t, err)
		})
	}

	t.Run("unknown location", func(t *testing.T) {
		err := ValidateParamStyle("simple", "body")
		assert.Error(t, err)
	})
}
