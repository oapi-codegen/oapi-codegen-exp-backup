package codegen

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/templates"
)

// structTemplateProperty is a pre-computed property for struct templates.
// This mirrors the original oapi-codegen's Property type, providing methods
// like GoFieldName and RequiresNilCheck as pre-computed fields so templates
// can use them directly without needing access to the NameConverter.
type structTemplateProperty struct {
	GoFieldName         string // Go field name (e.g., "AuthType")
	JSONFieldName       string // JSON property name (e.g., "auth_type")
	Type                string // Full Go type (e.g., "*string")
	BaseType            string // Go type without pointer (e.g., "string")
	Pointer             bool   // Whether this is a pointer type
	Required            bool   // Whether the field is required
	RequiresNilCheck    bool   // Whether marshal needs a nil guard
	Default             string // Go literal for default value (empty if none)
	NeedsTypeConversion bool   // Whether default needs explicit type conversion
	IsStruct            bool   // Whether this is a struct type (for recursive ApplyDefaults)
	IsExternal          bool   // Whether this references an external type
}

// structTemplateData is the data passed to struct-related templates.
type structTemplateData struct {
	TypeName     string
	AddPropsType string // Go type for additional properties (e.g., "any", "int")
	Properties   []structTemplateProperty
	NeedsReflect bool // Whether the ApplyDefaults method needs the reflect package
}

// buildStructTemplateData converts StructFields into the enriched template data.
func buildStructTemplateData(typeName string, fields []StructField, addPropsType string) structTemplateData {
	data := structTemplateData{
		TypeName:     typeName,
		AddPropsType: addPropsType,
	}

	for _, f := range fields {
		baseType := strings.TrimPrefix(f.Type, "*")
		prop := structTemplateProperty{
			GoFieldName:         f.Name,
			JSONFieldName:       f.JSONName,
			Type:                f.Type,
			BaseType:            baseType,
			Pointer:             f.Pointer,
			Required:            f.Required,
			RequiresNilCheck:    f.Pointer,
			Default:             f.Default,
			NeedsTypeConversion: needsTypeConversion(baseType),
			IsStruct:            f.IsStruct,
			IsExternal:          f.IsExternal,
		}
		if f.IsExternal && f.Pointer {
			data.NeedsReflect = true
		}
		data.Properties = append(data.Properties, prop)
	}

	return data
}

// loadStructTemplates loads and parses the struct-related templates.
func loadStructTemplates() (*template.Template, error) {
	entries := []string{
		"files/struct/additional-properties.go.tmpl",
		"files/struct/apply-defaults.go.tmpl",
	}

	tmpl := template.New("struct")

	for _, entry := range entries {
		content, err := templates.TemplateFS.ReadFile(entry)
		if err != nil {
			return nil, fmt.Errorf("reading template %s: %w", entry, err)
		}
		if _, err := tmpl.Parse(string(content)); err != nil {
			return nil, fmt.Errorf("parsing template %s: %w", entry, err)
		}
	}

	return tmpl, nil
}

// GenerateAdditionalPropertiesCode generates Get/Set + MarshalJSON/UnmarshalJSON
// for structs with additionalProperties.
func GenerateAdditionalPropertiesCode(typeName string, fields []StructField, addPropsType string) (string, error) {
	data := buildStructTemplateData(typeName, fields, addPropsType)

	tmpl, err := loadStructTemplates()
	if err != nil {
		return "", err
	}

	var buf strings.Builder

	if err := tmpl.ExecuteTemplate(&buf, "additional_properties_accessors", data); err != nil {
		return "", fmt.Errorf("executing additional_properties_accessors: %w", err)
	}
	if err := tmpl.ExecuteTemplate(&buf, "additional_properties_unmarshal", data); err != nil {
		return "", fmt.Errorf("executing additional_properties_unmarshal: %w", err)
	}
	if err := tmpl.ExecuteTemplate(&buf, "additional_properties_marshal", data); err != nil {
		return "", fmt.Errorf("executing additional_properties_marshal: %w", err)
	}

	return buf.String(), nil
}

// GenerateApplyDefaultsCode generates the ApplyDefaults method for a struct.
// Returns the generated code and whether the reflect package is needed.
func GenerateApplyDefaultsCode(typeName string, fields []StructField) (string, bool, error) {
	data := buildStructTemplateData(typeName, fields, "")

	tmpl, err := loadStructTemplates()
	if err != nil {
		return "", false, err
	}

	var buf strings.Builder
	if err := tmpl.ExecuteTemplate(&buf, "apply_defaults", data); err != nil {
		return "", false, fmt.Errorf("executing apply_defaults: %w", err)
	}

	return buf.String(), data.NeedsReflect, nil
}
