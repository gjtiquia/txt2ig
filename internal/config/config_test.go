package config

import (
	"reflect"
	"testing"
)

func TestDefaultConfig_HasFontFamily(t *testing.T) {
	cfg := DefaultConfig()

	if len(cfg.FontFamily.Regular) == 0 {
		t.Error("DefaultConfig FontFamily.Regular should not be empty")
	}
	if len(cfg.FontFamily.Bold) == 0 {
		t.Error("DefaultConfig FontFamily.Bold should not be empty")
	}
	if len(cfg.FontFamily.Italic) == 0 {
		t.Error("DefaultConfig FontFamily.Italic should not be empty")
	}
	if len(cfg.FontFamily.BoldItalic) == 0 {
		t.Error("DefaultConfig FontFamily.BoldItalic should not be empty")
	}

	expectedRegular := []string{"GoMono"}
	if !reflect.DeepEqual(cfg.FontFamily.Regular, expectedRegular) {
		t.Errorf("DefaultConfig FontFamily.Regular = %v, want %v", cfg.FontFamily.Regular, expectedRegular)
	}

	expectedBold := []string{"GoMonoBold"}
	if !reflect.DeepEqual(cfg.FontFamily.Bold, expectedBold) {
		t.Errorf("DefaultConfig FontFamily.Bold = %v, want %v", cfg.FontFamily.Bold, expectedBold)
	}

	expectedItalic := []string{"GoMonoItalic"}
	if !reflect.DeepEqual(cfg.FontFamily.Italic, expectedItalic) {
		t.Errorf("DefaultConfig FontFamily.Italic = %v, want %v", cfg.FontFamily.Italic, expectedItalic)
	}

	expectedBoldItalic := []string{"GoMonoBoldItalic"}
	if !reflect.DeepEqual(cfg.FontFamily.BoldItalic, expectedBoldItalic) {
		t.Errorf("DefaultConfig FontFamily.BoldItalic = %v, want %v", cfg.FontFamily.BoldItalic, expectedBoldItalic)
	}
}

func TestFontFamilyJSONTags(t *testing.T) {
	cfg := FontFamily{
		Regular:    []string{"CustomFont"},
		Bold:       []string{"CustomFontBold"},
		Italic:     []string{"CustomFontItalic"},
		BoldItalic: []string{"CustomFontBoldItalic"},
	}

	if len(cfg.Regular) == 0 || cfg.Regular[0] != "CustomFont" {
		t.Errorf("FontFamily.Regular not properly set")
	}
	if len(cfg.Bold) == 0 || cfg.Bold[0] != "CustomFontBold" {
		t.Errorf("FontFamily.Bold not properly set")
	}
	if len(cfg.Italic) == 0 || cfg.Italic[0] != "CustomFontItalic" {
		t.Errorf("FontFamily.Italic not properly set")
	}
	if len(cfg.BoldItalic) == 0 || cfg.BoldItalic[0] != "CustomFontBoldItalic" {
		t.Errorf("FontFamily.BoldItalic not properly set")
	}
}
