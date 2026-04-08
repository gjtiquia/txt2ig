package config

import (
	"reflect"
	"testing"
)

func TestParseJSONC(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "simple JSON",
			input:    `{"key": "value"}`,
			expected: map[string]interface{}{"key": "value"},
			wantErr:  false,
		},
		{
			name:     "JSON with number",
			input:    `{"fontSize": 18}`,
			expected: map[string]interface{}{"fontSize": float64(18)},
			wantErr:  false,
		},
		{
			name:     "JSON with boolean",
			input:    `{"textWrap": true}`,
			expected: map[string]interface{}{"textWrap": true},
			wantErr:  false,
		},
		{
			name:     "JSON with array",
			input:    `{"screenSize": [1080, 1920]}`,
			expected: map[string]interface{}{"screenSize": []interface{}{float64(1080), float64(1920)}},
			wantErr:  false,
		},
		{
			name: "single-line comment",
			input: `{"key": "value" // this is a comment
}`,
			expected: map[string]interface{}{"key": "value"},
			wantErr:  false,
		},
		{
			name:     "multi-line comment",
			input:    `{"key": /* comment */ "value"}`,
			expected: map[string]interface{}{"key": "value"},
			wantErr:  false,
		},
		{
			name:     "comment in string should be preserved",
			input:    `{"key": "value // not a comment"}`,
			expected: map[string]interface{}{"key": "value // not a comment"},
			wantErr:  false,
		},
		{
			name:     "multi-line comment in string should be preserved",
			input:    `{"key": "value /* not a comment */"}`,
			expected: map[string]interface{}{"key": "value /* not a comment */"},
			wantErr:  false,
		},
		{
			name: "multiple single-line comments",
			input: `{
	// Comment 1
	"key": "value"
	// Comment 2
}`,
			expected: map[string]interface{}{"key": "value"},
			wantErr:  false,
		},
		{
			name:     "escaped quote in string",
			input:    `{"key": "value \"quoted\""}`,
			expected: map[string]interface{}{"key": `value "quoted"`},
			wantErr:  false,
		},
		{
			name:    "invalid JSON should error",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:     "empty object",
			input:    `{}`,
			expected: map[string]interface{}{},
			wantErr:  false,
		},
		// Trailing comma tests
		{
			name:     "trailing comma in object",
			input:    `{"key": "value",}`,
			expected: map[string]interface{}{"key": "value"},
			wantErr:  false,
		},
		{
			name:     "trailing comma in array",
			input:    `{"arr": [1, 2, 3,]}`,
			expected: map[string]interface{}{"arr": []interface{}{float64(1), float64(2), float64(3)}},
			wantErr:  false,
		},
		{
			name:     "multiple trailing commas",
			input:    `{"a": 1, "b": 2,}`,
			expected: map[string]interface{}{"a": float64(1), "b": float64(2)},
			wantErr:  false,
		},
		{
			name: "trailing comma with newlines",
			input: `{
	"key": "value",
	"num": 42,
}`,
			expected: map[string]interface{}{"key": "value", "num": float64(42)},
			wantErr:  false,
		},
		{
			name: "trailing comma with comments",
			input: `{
	"key": "value", // comment
	"num": 42, /* another comment */
}`,
			expected: map[string]interface{}{"key": "value", "num": float64(42)},
			wantErr:  false,
		},
		{
			name: "complex config with trailing commas",
			input: `{
	"fontSize": 18,
	"fontColor": "#FFFFFF",
	"bgColor": "#000000",
}`,
			expected: map[string]interface{}{
				"fontSize":  float64(18),
				"fontColor": "#FFFFFF",
				"bgColor":   "#000000",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}
			err := ParseJSONC([]byte(tt.input), &result)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseJSONC() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseJSONC() unexpected error: %v", err)
				return
			}

			// Compare results
			if len(result) != len(tt.expected) {
				t.Errorf("ParseJSONC() result length = %d, want %d", len(result), len(tt.expected))
				return
			}

			for key, expectedVal := range tt.expected {
				resultVal, ok := result[key]
				if !ok {
					t.Errorf("ParseJSONC() missing key %q", key)
					continue
				}
				if !reflect.DeepEqual(resultVal, expectedVal) {
					t.Errorf("ParseJSONC()[%q] = %v, want %v", key, resultVal, expectedVal)
				}
			}
		})
	}
}
