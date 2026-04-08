package processor

import (
	"testing"
)

func TestMarkdownBoldHeaders(t *testing.T) {
	tests := []struct {
		name          string
		bold          bool
		fontColor     string
		input         string
		expectedLine  string
		expectedBold  bool
		expectedColor string
	}{
		{
			name:          "single hash header",
			bold:          true,
			fontColor:     "#FF0000",
			input:         "# Title",
			expectedLine:  "# Title", // # is kept
			expectedBold:  true,
			expectedColor: "#FF0000",
		},
		{
			name:          "double hash header",
			bold:          true,
			fontColor:     "",
			input:         "## Heading",
			expectedLine:  "## Heading", // # is kept
			expectedBold:  true,
			expectedColor: "",
		},
		{
			name:          "triple hash header",
			bold:          false,
			fontColor:     "#00FF00",
			input:         "### Subheading",
			expectedLine:  "### Subheading", // # is kept
			expectedBold:  false,
			expectedColor: "#00FF00",
		},
		{
			name:          "not a header",
			bold:          true,
			fontColor:     "#FF0000",
			input:         "Normal text",
			expectedLine:  "Normal text",
			expectedBold:  false,
			expectedColor: "",
		},
		{
			name:          "header with spaces after hash",
			bold:          true,
			fontColor:     "",
			input:         "#   Heading with spaces",
			expectedLine:  "#   Heading with spaces", // # is kept
			expectedBold:  true,
			expectedColor: "",
		},
		{
			name:          "just hash",
			bold:          true,
			fontColor:     "#000000",
			input:         "#",
			expectedLine:  "#",
			expectedBold:  true,
			expectedColor: "#000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MarkdownBoldHeaders{
				Bold:      tt.bold,
				FontColor: tt.fontColor,
			}
			line, style, err := p.Process(tt.input)
			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}
			if line != tt.expectedLine {
				t.Errorf("Process() line = %q, want %q", line, tt.expectedLine)
			}
			if style != nil {
				if style.Bold != tt.expectedBold {
					t.Errorf("Style.Bold = %v, want %v", style.Bold, tt.expectedBold)
				}
				if style.FontColor != tt.expectedColor {
					t.Errorf("Style.FontColor = %q, want %q", style.FontColor, tt.expectedColor)
				}
			} else if tt.expectedBold || tt.expectedColor != "" {
				t.Errorf("Style is nil, expected non-nil")
			}
		})
	}
}

func TestBashComments(t *testing.T) {
	tests := []struct {
		name           string
		italic         bool
		fontColor      string
		input          string
		expectedLine   string
		expectedItalic bool
		expectedColor  string
	}{
		{
			name:           "bash comment",
			italic:         true,
			fontColor:      "#888888",
			input:          "# This is a comment",
			expectedLine:   "# This is a comment", // Bash comments preserve the #
			expectedItalic: true,
			expectedColor:  "#888888",
		},
		{
			name:           "not a comment",
			italic:         true,
			fontColor:      "#888888",
			input:          "echo hello",
			expectedLine:   "echo hello",
			expectedItalic: false,
			expectedColor:  "",
		},
		{
			name:           "comment at end",
			italic:         false,
			fontColor:      "#CCCCCC",
			input:          "# comment",
			expectedLine:   "# comment",
			expectedItalic: false,
			expectedColor:  "#CCCCCC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &BashComments{
				Italic:    tt.italic,
				FontColor: tt.fontColor,
			}
			line, style, err := p.Process(tt.input)
			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}
			if line != tt.expectedLine {
				t.Errorf("Process() line = %q, want %q", line, tt.expectedLine)
			}
			if style != nil {
				if style.Italic != tt.expectedItalic {
					t.Errorf("Style.Italic = %v, want %v", style.Italic, tt.expectedItalic)
				}
				if style.FontColor != tt.expectedColor {
					t.Errorf("Style.FontColor = %q, want %q", style.FontColor, tt.expectedColor)
				}
			} else if tt.expectedItalic || tt.expectedColor != "" {
				t.Errorf("Style is nil, expected non-nil")
			}
		})
	}
}
