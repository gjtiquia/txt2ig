package renderer

import (
	"testing"

	"github.com/gjtiquia/txt2ig/internal/config"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

func TestWrapTextPreservesNewlines(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		TextWrap:   true,
		LineHeight: 1.4,
		FontColor:  "#FFFFFF", // Add required color
	}

	// Create a font face for testing
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatalf("Failed to create font face: %v", err)
	}

	renderer, err := NewTextRenderer(face, cfg)
	if err != nil {
		t.Fatalf("Failed to create text renderer: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		maxWidth int
		expected []string
	}{
		{
			name:     "single line",
			input:    "Hello World",
			maxWidth: 1000,
			expected: []string{"Hello World"},
		},
		{
			name:     "preserves single newline",
			input:    "Hello\nWorld",
			maxWidth: 1000,
			expected: []string{"Hello", "World"},
		},
		{
			name:     "preserves empty line",
			input:    "Hello\n\nWorld",
			maxWidth: 1000,
			expected: []string{"Hello", "", "World"},
		},
		{
			name:     "preserves multiple empty lines",
			input:    "A\n\n\nB",
			maxWidth: 1000,
			expected: []string{"A", "", "", "B"},
		},
		{
			name:     "empty input",
			input:    "",
			maxWidth: 1000,
			expected: []string{""},
		},
		{
			name:     "only newlines",
			input:    "\n\n",
			maxWidth: 1000,
			expected: []string{"", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderer.WrapText(tt.input, tt.maxWidth)
			if len(result) != len(tt.expected) {
				t.Errorf("WrapText(%q) returned %d lines, want %d lines: got %v, want %v",
					tt.input, len(result), len(tt.expected), result, tt.expected)
				return
			}
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("WrapText(%q) line %d = %q, want %q",
						tt.input, i, line, tt.expected[i])
				}
			}
		})
	}
}

func TestWrapTextWithLongLines(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		TextWrap:   true,
		LineHeight: 1.4,
		FontColor:  "#FFFFFF", // Add required color
	}

	// Create a font face for testing
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatalf("Failed to create font face: %v", err)
	}

	renderer, err := NewTextRenderer(face, cfg)
	if err != nil {
		t.Fatalf("Failed to create text renderer: %v", err)
	}

	tests := []struct {
		name             string
		input            string
		maxWidth         int
		atLeastLines     int
		shouldContainNew bool
	}{
		{
			name:         "long line wraps",
			input:        "This is a very long line that should wrap when the width is small",
			maxWidth:     200,
			atLeastLines: 2, // Should wrap to at least 2 lines
		},
		{
			name:             "wraps with newlines preserved",
			input:            "Paragraph one\n\nParagraph two is longer and needs wrapping",
			maxWidth:         200,
			atLeastLines:     3, // At least: "Paragraph one", "", and wrapped "Paragraph two..."
			shouldContainNew: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderer.WrapText(tt.input, tt.maxWidth)
			if len(result) < tt.atLeastLines {
				t.Errorf("WrapText(%q) returned %d lines, want at least %d lines: got %v",
					tt.input, len(result), tt.atLeastLines, result)
			}
			if tt.shouldContainNew {
				hasEmptyLine := false
				for _, line := range result {
					if line == "" {
						hasEmptyLine = true
						break
					}
				}
				if !hasEmptyLine {
					t.Errorf("WrapText(%q) should preserve empty lines, but got: %v",
						tt.input, result)
				}
			}
		})
	}
}

func TestWrapTextDisabled(t *testing.T) {
	// Create a test config with wrapping disabled
	cfg := &config.Config{
		TextWrap:   false,
		LineHeight: 1.4,
		FontColor:  "#FFFFFF", // Add required color
	}

	// Create a font face for testing
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatalf("Failed to create font face: %v", err)
	}

	renderer, err := NewTextRenderer(face, cfg)
	if err != nil {
		t.Fatalf("Failed to create text renderer: %v", err)
	}

	input := "Hello\n\nWorld\nTest"
	result := renderer.WrapText(input, 1000)

	// Even with textWrap disabled, newlines should be preserved
	expected := []string{"Hello", "", "World", "Test"}
	if len(result) != len(expected) {
		t.Errorf("WrapText with textWrap=false returned %d lines, want %d lines: got %v, want %v",
			len(result), len(expected), result, expected)
		return
	}
	for i, line := range result {
		if line != expected[i] {
			t.Errorf("WrapText with textWrap=false line %d = %q, want %q",
				i, line, expected[i])
		}
	}
}
