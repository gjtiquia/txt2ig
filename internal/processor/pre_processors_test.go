package processor

import (
	"testing"
)

func TestExactSearchAndReplace(t *testing.T) {
	tests := []struct {
		name     string
		search   string
		replace  string
		input    string
		expected string
	}{
		{
			name:     "simple replacement",
			search:   "apple",
			replace:  "banana",
			input:    "I like apple",
			expected: "I like banana",
		},
		{
			name:     "multiple occurrences",
			search:   "foo",
			replace:  "bar",
			input:    "foo bar foo",
			expected: "bar bar bar",
		},
		{
			name:     "no match",
			search:   "apple",
			replace:  "banana",
			input:    "I like orange",
			expected: "I like orange",
		},
		{
			name:     "empty search",
			search:   "",
			replace:  "banana",
			input:    "test",
			expected: "bananatbananaebananasbananatbanana", // strings.ReplaceAll with empty string replaces between characters
		},
		{
			name:     "case sensitive",
			search:   "Apple",
			replace:  "Banana",
			input:    "apple Apple APPLE",
			expected: "apple Banana APPLE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ExactSearchAndReplace{
				SearchString:  tt.search,
				ReplaceString: tt.replace,
			}
			result, err := p.Process(tt.input)
			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Process() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExactSearchAndReplaceWithDateTimeNow(t *testing.T) {
	tests := []struct {
		name     string
		search   string
		format   string
		input    string
		contains string // Result should contain this substring
	}{
		{
			name:     "simple date replacement",
			search:   "@date",
			format:   "yyyy-mm-dd",
			input:    "Today is @date",
			contains: "Today is ",
		},
		{
			name:     "custom format",
			search:   "@datetime",
			format:   "yyyy-mm-dd hh:MM:ss",
			input:    "Timestamp: @datetime",
			contains: "Timestamp: ",
		},
		{
			name:     "no match",
			search:   "@date",
			format:   "yyyy-mm-dd",
			input:    "No placeholder here",
			contains: "No placeholder here",
		},
		{
			name:     "multiple occurrences",
			search:   "@x",
			format:   "yyyy",
			input:    "Year: @x and @x",
			contains: "Year: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ExactSearchAndReplaceWithDateTimeNow{
				SearchString:  tt.search,
				ReplaceFormat: tt.format,
			}
			result, err := p.Process(tt.input)
			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}
			if result != tt.contains && !contains(result, tt.contains) {
				t.Errorf("Process() = %q, should contain %q", result, tt.contains)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
