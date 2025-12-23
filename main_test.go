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

func TestUnflatten(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
	}{
		{
			name: "simple map",
			lines: []string{
				"key1 = value1",
				"key2 = value2",
			},
		},
		{
			name: "nested map",
			lines: []string{
				"outer.inner = value",
			},
		},
		{
			name: "array",
			lines: []string{
				"items[0] = a",
				"items[1] = b",
				"items[2] = c",
			},
		},
		{
			name: "complex nested structure",
			lines: []string{
				"level1.level2.level3 = value",
			},
		},
		{
			name: "mixed types",
			lines: []string{
				"string = text",
				"number = 42",
				"bool = true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unflatten the lines
			result, err := unflatten(tt.lines)
			if err != nil {
				t.Errorf("unflatten error: %v", err)
				return
			}

			// Flatten again to compare
			var flattened []string
			flatten(result, "", &flattened)

			// Create maps for comparison (order doesn't matter)
			expectedMap := make(map[string]bool)
			for _, line := range tt.lines {
				expectedMap[line] = true
			}

			resultMap := make(map[string]bool)
			for _, line := range flattened {
				resultMap[line] = true
			}

			// Check all expected lines are present
			for line := range expectedMap {
				if !resultMap[line] {
					t.Errorf("expected line not found after round-trip: %s", line)
				}
			}
		})
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		path     string
		expected []interface{}
	}{
		{
			path:     "simple",
			expected: []interface{}{"simple"},
		},
		{
			path:     "nested.path",
			expected: []interface{}{"nested", "path"},
		},
		{
			path:     "array[0]",
			expected: []interface{}{"array", 0},
		},
		{
			path:     "complex.array[0].field",
			expected: []interface{}{"complex", "array", 0, "field"},
		},
		{
			path:     "multi[0][1]",
			expected: []interface{}{"multi", 0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := parsePath(tt.path)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("token %d: expected %v (type %T), got %v (type %T)",
						i, expected, expected, result[i], result[i])
				}
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"42", int64(42)},
		{"3.14", 3.14},
		{"true", true},
		{"false", false},
		{"text", "text"},
		{"hello world", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseValue(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v (type %T), got %v (type %T)",
					tt.expected, tt.expected, result, result)
			}
		})
	}
}
