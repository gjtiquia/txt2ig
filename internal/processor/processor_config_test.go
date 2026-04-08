package processor

import (
	"encoding/json"
	"testing"
)

func TestParsePreProcessorConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected PreProcessor
		wantErr  bool
	}{
		{
			name: "exact search and replace",
			input: map[string]interface{}{
				"exactSearchAndReplace": map[string]interface{}{
					"searchString":  "apple",
					"replaceString": "banana",
				},
			},
			expected: &ExactSearchAndReplace{
				SearchString:  "apple",
				ReplaceString: "banana",
			},
		},
		{
			name: "date time replacement",
			input: map[string]interface{}{
				"exactSearchAndReplaceWithDateTimeNow": map[string]interface{}{
					"searchString":  "@date",
					"replaceFormat": "yyyy-mm-dd",
				},
			},
			expected: &ExactSearchAndReplaceWithDateTimeNow{
				SearchString:  "@date",
				ReplaceFormat: "yyyy-mm-dd",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePreProcessorConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParsePreProcessorConfig() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ParsePreProcessorConfig() unexpected error: %v", err)
				return
			}
			if result.Name() != tt.expected.Name() {
				t.Errorf("ParsePreProcessorConfig() name = %s, want %s", result.Name(), tt.expected.Name())
			}
		})
	}
}

func TestParsePostProcessorConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected PostProcessor
		wantErr  bool
	}{
		{
			name: "markdown bold headers",
			input: map[string]interface{}{
				"markdown-bold-headers": map[string]interface{}{
					"bold":      true,
					"fontColor": "#EC9006",
				},
			},
			expected: &MarkdownBoldHeaders{
				Bold:      true,
				FontColor: "#EC9006",
			},
		},
		{
			name: "markdown bold headers with default color",
			input: map[string]interface{}{
				"markdown-bold-headers": map[string]interface{}{
					"bold": true,
				},
			},
			expected: &MarkdownBoldHeaders{
				Bold:      true,
				FontColor: "",
			},
		},
		{
			name: "bash comments",
			input: map[string]interface{}{
				"bash-comments": map[string]interface{}{
					"italic":    true,
					"fontColor": "#CCCCCC",
				},
			},
			expected: &BashComments{
				Italic:    true,
				FontColor: "#CCCCCC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePostProcessorConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParsePostProcessorConfig() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ParsePostProcessorConfig() unexpected error: %v", err)
				return
			}
			if result.Name() != tt.expected.Name() {
				t.Errorf("ParsePostProcessorConfig() name = %s, want %s", result.Name(), tt.expected.Name())
			}
		})
	}
}

func TestParsePreProcessorConfigs(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"exactSearchAndReplace": map[string]interface{}{
				"searchString":  "foo",
				"replaceString": "bar",
			},
		},
		map[string]interface{}{
			"exactSearchAndReplaceWithDateTimeNow": map[string]interface{}{
				"searchString":  "@date",
				"replaceFormat": "yyyy-mm-dd",
			},
		},
	}

	results, err := ParsePreProcessorConfigs(input)
	if err != nil {
		t.Errorf("ParsePreProcessorConfigs() unexpected error: %v", err)
		return
	}

	if len(results) != 2 {
		t.Errorf("ParsePreProcessorConfigs() got %d processors, want 2", len(results))
		return
	}

	if results[0].Name() != "exactSearchAndReplace" {
		t.Errorf("ParsePreProcessorConfigs()[0].Name() = %s, want exactSearchAndReplace", results[0].Name())
	}

	if results[1].Name() != "exactSearchAndReplaceWithDateTimeNow" {
		t.Errorf("ParsePreProcessorConfigs()[1].Name() = %s, want exactSearchAndReplaceWithDateTimeNow", results[1].Name())
	}
}

func TestParsePostProcessorConfigs(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"markdown-bold-headers": map[string]interface{}{
				"bold": true,
			},
		},
		map[string]interface{}{
			"bash-comments": map[string]interface{}{
				"italic":    true,
				"fontColor": "#888888",
			},
		},
	}

	results, err := ParsePostProcessorConfigs(input)
	if err != nil {
		t.Errorf("ParsePostProcessorConfigs() unexpected error: %v", err)
		return
	}

	if len(results) != 2 {
		t.Errorf("ParsePostProcessorConfigs() got %d processors, want 2", len(results))
		return
	}

	if results[0].Name() != "markdown-bold-headers" {
		t.Errorf("ParsePostProcessorConfigs()[0].Name() = %s, want markdown-bold-headers", results[0].Name())
	}

	if results[1].Name() != "bash-comments" {
		t.Errorf("ParsePostProcessorConfigs()[1].Name() = %s, want bash-comments", results[1].Name())
	}
}

// Helper to marshal then unmarshal for realistic testing
func marshalUnmarshal(t *testing.T, v interface{}) map[string]interface{} {
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	return result
}
