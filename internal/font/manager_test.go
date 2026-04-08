package font

import (
	"testing"
)

func TestLoadEmbeddedFont(t *testing.T) {
	m := NewManager()

	tests := []struct {
		name     string
		fontName string
		wantErr  bool
	}{
		{
			name:     "GoMono",
			fontName: "GoMono",
			wantErr:  false,
		},
		{
			name:     "GoMonoBold",
			fontName: "GoMonoBold",
			wantErr:  false,
		},
		{
			name:     "GoMonoItalic",
			fontName: "GoMonoItalic",
			wantErr:  false,
		},
		{
			name:     "GoMonoBoldItalic",
			fontName: "GoMonoBoldItalic",
			wantErr:  false,
		},
		{
			name:     "unknown font",
			fontName: "UnknownFont",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			face, err := m.loadEmbeddedFont(tt.fontName, 18, 72)
			if tt.wantErr {
				if err == nil {
					t.Errorf("loadEmbeddedFont(%q) expected error, got nil", tt.fontName)
				}
				return
			}

			if err != nil {
				t.Errorf("loadEmbeddedFont(%q) unexpected error: %v", tt.fontName, err)
				return
			}

			if face == nil {
				t.Errorf("loadEmbeddedFont(%q) returned nil face", tt.fontName)
			}
		})
	}
}

func TestLoadFontFamilyWithEmbedded(t *testing.T) {
	m := NewManager()

	familyConfig := FontFamilyConfig{
		Regular:    []string{"GoMono"},
		Bold:       []string{"GoMonoBold"},
		Italic:     []string{"GoMonoItalic"},
		BoldItalic: []string{"GoMonoBoldItalic"},
	}

	family, err := m.LoadFontFamily(familyConfig, 18, 72)
	if err != nil {
		t.Fatalf("LoadFontFamily unexpected error: %v", err)
	}

	if family.Regular == nil {
		t.Error("Regular font face is nil")
	}
	if family.Bold == nil {
		t.Error("Bold font face is nil")
	}
	if family.Italic == nil {
		t.Error("Italic font face is nil")
	}
	if family.BoldItalic == nil {
		t.Error("BoldItalic font face is nil")
	}
}

func TestLoadFontFamilyWithFallback(t *testing.T) {
	m := NewManager()

	familyConfig := FontFamilyConfig{
		Regular:    []string{"NonExistentFont", "GoMono"},
		Bold:       []string{"AnotherNonExistentFont", "GoMonoBold"},
		Italic:     []string{"YetAnotherNonExistent", "GoMonoItalic"},
		BoldItalic: []string{"StillNonExistent", "GoMonoBoldItalic"},
	}

	family, err := m.LoadFontFamily(familyConfig, 18, 72)
	if err != nil {
		t.Fatalf("LoadFontFamily unexpected error: %v", err)
	}

	if family.Regular == nil {
		t.Error("Regular font face should have fallen back to GoMono")
	}
	if family.Bold == nil {
		t.Error("Bold font face should have fallen back to GoMonoBold")
	}
	if family.Italic == nil {
		t.Error("Italic font face should have fallen back to GoMonoItalic")
	}
	if family.BoldItalic == nil {
		t.Error("BoldItalic font face should have fallen back to GoMonoBoldItalic")
	}
}
