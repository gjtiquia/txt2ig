package renderer

import (
	"testing"
)

func TestGetDefaultOutputPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "markdown extension with dot",
			input:    "test.md",
			expected: "test.jpg",
		},
		{
			name:     "txt extension with dot",
			input:    "post.txt",
			expected: "post.jpg",
		},
		{
			name:     "no extension",
			input:    "file",
			expected: "file.jpg",
		},
		{
			name:     "path with directory",
			input:    "dir/post.md",
			expected: "dir/post.jpg",
		},
		{
			name:     "nested path",
			input:    "a/b/c/file.md",
			expected: "a/b/c/file.jpg",
		},
		{
			name:     "file with multiple dots",
			input:    "my.post.md",
			expected: "my.post.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDefaultOutputPath(tt.input)
			if result != tt.expected {
				t.Errorf("GetDefaultOutputPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetermineOutputPath(t *testing.T) {
	tests := []struct {
		name       string
		inputPath  string
		outputFlag string
		expected   string
	}{
		{
			name:       "no output flag uses default",
			inputPath:  "test.md",
			outputFlag: "",
			expected:   "test.jpg",
		},
		{
			name:       "output flag overrides default",
			inputPath:  "test.md",
			outputFlag: "output.png",
			expected:   "output.png",
		},
		{
			name:       "custom output name",
			inputPath:  "test.md",
			outputFlag: "custom.jpg",
			expected:   "custom.jpg",
		},
		{
			name:       "path with directory no output flag",
			inputPath:  "dir/test.md",
			outputFlag: "",
			expected:   "dir/test.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineOutputPath(tt.inputPath, tt.outputFlag)
			if result != tt.expected {
				t.Errorf("DetermineOutputPath(%q, %q) = %q, want %q",
					tt.inputPath, tt.outputFlag, result, tt.expected)
			}
		})
	}
}
