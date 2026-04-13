package output

import (
	"encoding/json"
	"testing"
)

// TestValidIdentifiers verifies that all generated type names are valid Go identifiers.
// Issue 1496: Inline schemas in responses were generating identifiers starting with numbers.
func TestValidIdentifiers(t *testing.T) {
	// Verify union types exist and can be used via As/From
	var union GetSomething200ResponseJSON2
	err := union.FromGetSomething200ResponseJSONAnyOf0(GetSomething200ResponseJSONAnyOf0{
		Order: ptr("order-123"),
	})
	if err != nil {
		t.Fatalf("FromGetSomething200ResponseJSONAnyOf0 failed: %v", err)
	}

	// Should be able to marshal the union
	data, err := json.Marshal(union)
	if err != nil {
		t.Fatalf("Failed to marshal union: %v", err)
	}
	t.Logf("Marshaled union: %s", string(data))

	// Verify the overall response type compiles and works
	response := GetSomethingJSONResponse{
		Results: []GetSomething200ResponseJSON2{union},
	}

	data, err = json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	t.Logf("Marshaled response: %s", string(data))
}

func ptr[T any](v T) *T {
	return &v
}
