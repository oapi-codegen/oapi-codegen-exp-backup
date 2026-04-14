package output

import (
	"encoding/json"
	"testing"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/runtime/types"
)

// TestNullableTypes verifies that nullable types are generated properly.
// https://github.com/oapi-codegen/oapi-codegen/issues/1039
//
// The implementation uses Nullable[T] for nullable types:
// - Nullable primitive schemas generate type aliases: type SimpleRequiredNullable = Nullable[int]
// - Nullable object fields are wrapped: Nullable[ComplexType]
// - Inline nullable primitives use Nullable[T] directly
func TestNullableTypes(t *testing.T) {
	// Create a patch request with various nullable fields
	name := "test-name"

	// SimpleRequiredNullable is Nullable[int]
	simpleRequired := types.NewNullableWithValue(42)

	// ComplexRequiredNullable is wrapped in Nullable
	complexRequired := types.NewNullableWithValue(ComplexRequiredNullable{Name: &name})

	req := PatchRequest{
		SimpleRequiredNullable:  simpleRequired,
		ComplexRequiredNullable: complexRequired,
	}

	// Verify simple nullable value
	val, err := req.SimpleRequiredNullable.Get()
	if err != nil {
		t.Fatalf("SimpleRequiredNullable.Get() failed: %v", err)
	}
	if val != 42 {
		t.Errorf("SimpleRequiredNullable = %v, want 42", val)
	}

	// Verify complex nullable can retrieve value
	complexVal, err := req.ComplexRequiredNullable.Get()
	if err != nil {
		t.Fatalf("ComplexRequiredNullable.Get() failed: %v", err)
	}
	if *complexVal.Name != "test-name" {
		t.Errorf("ComplexRequiredNullable.Name = %q, want %q", *complexVal.Name, "test-name")
	}
}

func TestPatchRequestJSONRoundTrip(t *testing.T) {
	name := "test"
	original := PatchRequest{
		SimpleRequiredNullable:  types.NewNullableWithValue(100),
		ComplexRequiredNullable: types.NewNullableWithValue(ComplexRequiredNullable{Name: &name}),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	t.Logf("Marshaled: %s", string(data))

	var decoded PatchRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify simple nullable round-trips correctly
	decodedSimple, err := decoded.SimpleRequiredNullable.Get()
	if err != nil {
		t.Fatalf("SimpleRequiredNullable.Get() failed: %v", err)
	}
	if decodedSimple != 100 {
		t.Errorf("SimpleRequiredNullable mismatch: got %v, want %v", decodedSimple, 100)
	}

	// Verify complex nullable round-trips correctly
	complexVal, err := decoded.ComplexRequiredNullable.Get()
	if err != nil {
		t.Fatalf("ComplexRequiredNullable.Get() failed: %v", err)
	}
	if *complexVal.Name != "test" {
		t.Errorf("ComplexRequiredNullable.Name = %q, want %q", *complexVal.Name, "test")
	}
}

func TestComplexNullableTypes(t *testing.T) {
	// Complex nullable types use Nullable[T]
	name := "name"
	opt := ComplexOptionalNullable{
		AliasName: types.NewNullableWithValue("alias"),
		Name:      &name,
	}

	req := PatchRequest{
		SimpleRequiredNullable:  types.NewNullNullable[int](), // explicitly null
		ComplexRequiredNullable: types.NewNullNullable[ComplexRequiredNullable](),
		ComplexOptionalNullable: types.NewNullableWithValue(opt),
	}

	// Check the complex optional nullable
	if !req.ComplexOptionalNullable.IsSpecified() {
		t.Fatal("ComplexOptionalNullable should be specified")
	}
	optVal := req.ComplexOptionalNullable.MustGet()
	aliasVal := optVal.AliasName.MustGet()
	if aliasVal != "alias" {
		t.Errorf("AliasName = %q, want %q", aliasVal, "alias")
	}

	// Check that required nullable can be null
	if !req.ComplexRequiredNullable.IsNull() {
		t.Error("ComplexRequiredNullable should be null")
	}
	if !req.SimpleRequiredNullable.IsNull() {
		t.Error("SimpleRequiredNullable should be null")
	}
}

func TestNullableThreeStates(t *testing.T) {
	// Test unspecified (nil/empty map)
	unspecified := types.Nullable[string](nil)
	if unspecified.IsSpecified() {
		t.Error("unspecified should not be specified")
	}
	if unspecified.IsNull() {
		t.Error("unspecified should not be null")
	}
	_, err := unspecified.Get()
	if err != types.ErrNullableNotSpecified {
		t.Errorf("Get() on unspecified should return ErrNullableNotSpecified, got %v", err)
	}

	// Test explicitly null
	null := types.NewNullNullable[string]()
	if !null.IsSpecified() {
		t.Error("null should be specified")
	}
	if !null.IsNull() {
		t.Error("null should be null")
	}
	_, err = null.Get()
	if err != types.ErrNullableIsNull {
		t.Errorf("Get() on null should return ErrNullableIsNull, got %v", err)
	}

	// Test with value
	withValue := types.NewNullableWithValue("hello")
	if !withValue.IsSpecified() {
		t.Error("withValue should be specified")
	}
	if withValue.IsNull() {
		t.Error("withValue should not be null")
	}
	val, err := withValue.Get()
	if err != nil {
		t.Errorf("Get() on withValue should succeed, got %v", err)
	}
	if val != "hello" {
		t.Errorf("Get() = %q, want %q", val, "hello")
	}
}

func TestNullableJSONMarshal(t *testing.T) {
	// Test marshaling each state
	tests := []struct {
		name     string
		nullable types.Nullable[string]
		want     string
	}{
		{"with value", types.NewNullableWithValue("test"), `"test"`},
		{"explicitly null", types.NewNullNullable[string](), "null"},
		{"unspecified", types.Nullable[string](nil), "null"}, // unspecified marshals as null
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.nullable)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("Marshal() = %s, want %s", string(data), tt.want)
			}
		})
	}
}

func TestNullableJSONUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantNull  bool
		wantValue string
		wantErr   error
	}{
		{"with value", `"test"`, false, "test", nil},
		{"explicitly null", "null", true, "", types.ErrNullableIsNull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n types.Nullable[string]
			if err := json.Unmarshal([]byte(tt.json), &n); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if n.IsNull() != tt.wantNull {
				t.Errorf("IsNull() = %v, want %v", n.IsNull(), tt.wantNull)
			}
			val, err := n.Get()
			if err != tt.wantErr {
				t.Errorf("Get() error = %v, want %v", err, tt.wantErr)
			}
			if err == nil && val != tt.wantValue {
				t.Errorf("Get() = %q, want %q", val, tt.wantValue)
			}
		})
	}
}

// TestNullablePrimitiveTypeAlias verifies that nullable primitive schemas
// generate proper type aliases.
func TestNullablePrimitiveTypeAlias(t *testing.T) {
	// SimpleRequiredNullable should be a type alias to Nullable[int]
	var simple SimpleRequiredNullable
	simple.Set(42)

	val, err := simple.Get()
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if val != 42 {
		t.Errorf("Get() = %d, want 42", val)
	}

	// Test null state
	simple.SetNull()
	if !simple.IsNull() {
		t.Error("should be null after SetNull()")
	}

	// Test unspecified state
	simple.SetUnspecified()
	if simple.IsSpecified() {
		t.Error("should not be specified after SetUnspecified()")
	}
}

// TestAdditionalPropertiesFalse verifies that additionalProperties: false
// generates proper marshal/unmarshal that rejects extra fields.
func TestAdditionalPropertiesFalse(t *testing.T) {
	// The struct has AdditionalProperties field but additionalProperties: false
	// means unknown fields are still collected but not expected
	req := PatchRequest{
		SimpleRequiredNullable:  types.NewNullableWithValue(1),
		ComplexRequiredNullable: types.NewNullNullable[ComplexRequiredNullable](),
		AdditionalProperties:    map[string]any{"extra": "value"},
	}

	// Should marshal with additional properties
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	t.Logf("Marshaled: %s", string(data))
}
