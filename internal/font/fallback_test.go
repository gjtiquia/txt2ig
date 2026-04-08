package font

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected color.RGBA
		wantErr  bool
	}{
		{
			name:     "3-digit hex",
			input:    "#FFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			wantErr:  false,
		},
		{
			name:     "4-digit hex with alpha",
			input:    "#FFF0",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 0},
			wantErr:  false,
		},
		{
			name:     "6-digit hex",
			input:    "#FFFFFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			wantErr:  false,
		},
		{
			name:     "8-digit hex with alpha",
			input:    "#FFFFFFFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			wantErr:  false,
		},
		{
			name:     "lowercase hex",
			input:    "#ffffff",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			wantErr:  false,
		},
		{
			name:     "mixed case hex",
			input:    "#FfFfFf",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			wantErr:  false,
		},
		{
			name:     "black",
			input:    "#000000",
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			wantErr:  false,
		},
		{
			name:     "red",
			input:    "#FF0000",
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			wantErr:  false,
		},
		{
			name:     "green",
			input:    "#00FF00",
			expected: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			wantErr:  false,
		},
		{
			name:     "blue",
			input:    "#0000FF",
			expected: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			wantErr:  false,
		},
		{
			name:    "no hash prefix",
			input:   "FFF",
			wantErr: true,
		},
		{
			name:     "invalid hex character",
			input:    "#GGG",
			wantErr:  false, // Returns black with alpha 255
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "too short",
			input:   "#FF",
			wantErr: true, // Invalid length
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseColor(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseColor(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseColor(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("ParseColor(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		})
	}
}
