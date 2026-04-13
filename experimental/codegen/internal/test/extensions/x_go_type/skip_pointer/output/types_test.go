package output

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

// TestTypeWithOptionalFieldInstantiation verifies that TypeWithOptionalField
// uses uuid.UUID fields (via the googleuuid alias in generated code).
func TestTypeWithOptionalFieldInstantiation(t *testing.T) {
	id1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	id2 := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	tw := TypeWithOptionalField{
		At:         id1,
		AtRequired: id2,
	}

	if tw.At != id1 {
		t.Errorf("At = %v, want %v", tw.At, id1)
	}
	if tw.AtRequired != id2 {
		t.Errorf("AtRequired = %v, want %v", tw.AtRequired, id2)
	}
}

// TestTypeWithOptionalFieldJSONRoundTrip verifies JSON round-trip.
func TestTypeWithOptionalFieldJSONRoundTrip(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	original := TypeWithOptionalField{
		At:         id,
		AtRequired: id,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded TypeWithOptionalField
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.At != id {
		t.Errorf("At = %v, want %v", decoded.At, id)
	}
	if decoded.AtRequired != id {
		t.Errorf("AtRequired = %v, want %v", decoded.AtRequired, id)
	}
}

// TestTypeWithAllOfInstantiation verifies TypeWithAllOf with its ID field.
// After the allOf extension merge fix, ID should be googleuuid.UUID (value type, non-pointer)
// because x-go-type-skip-optional-pointer is merged from the allOf member.
func TestTypeWithAllOfInstantiation(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	tw := TypeWithAllOf{
		ID: id,
	}

	if tw.ID != id {
		t.Errorf("ID = %v, want %v", tw.ID, id)
	}
}

// TestTypeWithAllOfJSONRoundTrip verifies JSON round-trip for TypeWithAllOf.
func TestTypeWithAllOfJSONRoundTrip(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	original := TypeWithAllOf{
		ID: id,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded TypeWithAllOf
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ID != id {
		t.Errorf("ID = %v, want %v", decoded.ID, id)
	}
}

// TestTypeAliases verifies that UUID-related type aliases work correctly.
func TestTypeAliases(t *testing.T) {
	id := uuid.New()

	// ID is an alias for googleuuid.UUID (= uuid.UUID)
	idAlias := ID(id)
	if idAlias != id {
		t.Errorf("ID alias = %v, want %v", idAlias, id)
	}

	// GetRootParameter is an alias for googleuuid.UUID
	param := GetRootParameter(id)
	if param != id {
		t.Errorf("GetRootParameter alias = %v, want %v", param, id)
	}

	// TypeWithOptionalFieldAt is an alias for googleuuid.UUID
	at := TypeWithOptionalFieldAt(id)
	if at != id {
		t.Errorf("TypeWithOptionalFieldAt alias = %v, want %v", at, id)
	}

	// TypeWithOptionalFieldAtRequired is an alias for googleuuid.UUID
	atReq := TypeWithOptionalFieldAtRequired(id)
	if atReq != id {
		t.Errorf("TypeWithOptionalFieldAtRequired alias = %v, want %v", atReq, id)
	}

	// TypeWithAllOfID is now an alias for googleuuid.UUID (after allOf extension merge fix)
	allOfID := TypeWithAllOfID(id)
	if allOfID != id {
		t.Errorf("TypeWithAllOfID alias = %v, want %v", allOfID, id)
	}
}

// TestApplyDefaults verifies ApplyDefaults does not panic on all types.
func TestApplyDefaults(t *testing.T) {
	(&TypeWithOptionalField{}).ApplyDefaults()
	(&TypeWithAllOf{}).ApplyDefaults()
	// TypeWithAllOfID is now a type alias (= googleuuid.UUID), not a struct,
	// so it does not have ApplyDefaults.
}

// TestGetOpenAPISpecJSON verifies the embedded spec can be decoded.
func TestGetOpenAPISpecJSON(t *testing.T) {
	data, err := GetOpenAPISpecJSON()
	if err != nil {
		t.Fatalf("GetOpenAPISpecJSON failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("GetOpenAPISpecJSON returned empty data")
	}
}
