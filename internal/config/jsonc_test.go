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

func TestStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no comments",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "single-line comment only",
			input:    `{"key": "value"} // comment`,
			expected: `{"key": "value"} `,
		},
		{
			name:     "multi-line comment only",
			input:    `{"key": /* comment */ "value"}`,
			expected: `{"key":  "value"}`,
		},
		{
			name:     "comment in string preserved",
			input:    `{"key": "value // comment"}`,
			expected: `{"key": "value // comment"}`,
		},
		{
			name:     "comment at end",
			input:    `{"key": "value"} // trailing`,
			expected: `{"key": "value"} `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripComments([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("stripComments() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}
