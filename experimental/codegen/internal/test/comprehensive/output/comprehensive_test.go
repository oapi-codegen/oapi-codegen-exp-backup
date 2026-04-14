package output

import (
	"testing"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/runtime/types"
)

// TestNullableRefOneOf verifies that oneOf: [$ref, {type: "null"}] generates
// Nullable[T] instead of a full union struct. This is the standard OpenAPI 3.1
// pattern for making a $ref'd type nullable.
// https://github.com/oapi-codegen/oapi-codegen-exp/issues/19
func TestNullableRefOneOf(t *testing.T) {
	obj := NullableRefOneOf{
		NullableObject: types.NewNullableWithValue(SimpleObject{ID: 1}),
	}
	val := obj.NullableObject.MustGet()
	if val.ID != 1 {
		t.Errorf("NullableObject.ID = %d, want 1", val.ID)
	}

	// Test null state
	obj.NullableObject.SetNull()
	if !obj.NullableObject.IsNull() {
		t.Error("NullableObject should be null after SetNull()")
	}

	// Optional nullable ref should also be Nullable[T] (not *Nullable[T]),
	// because Nullable already handles the "unspecified" state.
	obj2 := NullableRefOneOf{
		NullableObject:         types.NewNullableWithValue(SimpleObject{ID: 2}),
		NullableObjectOptional: types.NewNullableWithValue(SimpleObject{ID: 3}),
	}
	_ = obj2
}

// TestNullableRefAnyOf verifies the same pattern using anyOf instead of oneOf.
func TestNullableRefAnyOf(t *testing.T) {
	obj := NullableRefAnyOf{
		NullableObject: types.NewNullableWithValue(SimpleObject{ID: 1}),
	}
	val := obj.NullableObject.MustGet()
	if val.ID != 1 {
		t.Errorf("NullableObject.ID = %d, want 1", val.ID)
	}
}
