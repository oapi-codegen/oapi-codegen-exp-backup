package helpers

//oapi-runtime:function helpers/JSONMerge

import (
	"encoding/json"
	"fmt"
)

// JSONMerge merges two JSON-encoded objects. Fields from patch override
// fields in base. Both arguments must be valid JSON objects (or nil/null).
func JSONMerge(base, patch json.RawMessage) (json.RawMessage, error) {
	if len(base) == 0 || string(base) == "null" {
		return patch, nil
	}
	if len(patch) == 0 || string(patch) == "null" {
		return base, nil
	}

	var baseMap map[string]json.RawMessage
	if err := json.Unmarshal(base, &baseMap); err != nil {
		return nil, fmt.Errorf("JSONMerge: unmarshaling base: %w", err)
	}

	var patchMap map[string]json.RawMessage
	if err := json.Unmarshal(patch, &patchMap); err != nil {
		return nil, fmt.Errorf("JSONMerge: unmarshaling patch: %w", err)
	}

	for k, v := range patchMap {
		baseMap[k] = v
	}

	return json.Marshal(baseMap)
}
