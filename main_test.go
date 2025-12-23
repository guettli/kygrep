package main

import (
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		prefix   string
		expected []string
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			prefix:   "",
			expected: []string{"key1 = value1", "key2 = value2"},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
			prefix:   "",
			expected: []string{"outer.inner = value"},
		},
		{
			name: "array",
			input: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			prefix:   "",
			expected: []string{"items[0] = a", "items[1] = b", "items[2] = c"},
		},
		{
			name: "complex nested structure",
			input: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": "value",
					},
				},
			},
			prefix:   "",
			expected: []string{"level1.level2.level3 = value"},
		},
		{
			name: "mixed types",
			input: map[string]interface{}{
				"string": "text",
				"number": 42,
				"bool":   true,
			},
			prefix:   "",
			expected: []string{"string = text", "number = 42", "bool = true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []string
			flatten(tt.input, tt.prefix, &result)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d lines, got %d", len(tt.expected), len(result))
				return
			}

			// Create a map for easier comparison (order doesn't matter for maps)
			resultMap := make(map[string]bool)
			for _, line := range result {
				resultMap[line] = true
			}

			for _, expected := range tt.expected {
				if !resultMap[expected] {
					t.Errorf("expected line not found: %s\nGot: %v", expected, result)
				}
			}
		})
	}
}

func TestFlattenWithPrefix(t *testing.T) {
	input := map[string]interface{}{
		"key": "value",
	}
	var result []string
	flatten(input, "root", &result)

	expected := "root.key = value"
	if len(result) != 1 || result[0] != expected {
		t.Errorf("expected %s, got %v", expected, result)
	}
}

func TestFlattenArrayWithPrefix(t *testing.T) {
	input := []interface{}{"a", "b"}
	var result []string
	flatten(input, "arr", &result)

	expectedLines := []string{"arr[0] = a", "arr[1] = b"}
	if len(result) != len(expectedLines) {
		t.Errorf("expected %d lines, got %d", len(expectedLines), len(result))
		return
	}

	for i, expected := range expectedLines {
		if result[i] != expected {
			t.Errorf("line %d: expected %s, got %s", i, expected, result[i])
		}
	}
}
