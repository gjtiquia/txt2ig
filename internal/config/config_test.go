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

func TestDefaultConfig_FromEmbeddedFile(t *testing.T) {
	cfg := DefaultConfig()

	// Verify values match the embedded default.jsonc from README
	if cfg.FontSize != 32 {
		t.Errorf("DefaultConfig FontSize = %d, want 32", cfg.FontSize)
	}

	if cfg.FontColor != "#FFFFFF" {
		t.Errorf("DefaultConfig FontColor = %s, want #FFFFFF", cfg.FontColor)
	}

	if cfg.BgColor != "#000000" {
		t.Errorf("DefaultConfig BgColor = %s, want #000000", cfg.BgColor)
	}

	if cfg.TextJustify != "left" {
		t.Errorf("DefaultConfig TextJustify = %s, want left", cfg.TextJustify)
	}

	if cfg.TextBoxJustify != "center" {
		t.Errorf("DefaultConfig TextBoxJustify = %s, want center", cfg.TextBoxJustify)
	}

	if cfg.TextBoxAlign != "center" {
		t.Errorf("DefaultConfig TextBoxAlign = %s, want center", cfg.TextBoxAlign)
	}

	if !reflect.DeepEqual(cfg.TextBoxOffset, []int{0, 0}) {
		t.Errorf("DefaultConfig TextBoxOffset = %v, want [0 0]", cfg.TextBoxOffset)
	}

	if cfg.TextBoxMaxWidth != 972 {
		t.Errorf("DefaultConfig TextBoxMaxWidth = %d, want 972", cfg.TextBoxMaxWidth)
	}

	if !reflect.DeepEqual(cfg.ScreenSize, []int{1080, 1920}) {
		t.Errorf("DefaultConfig ScreenSize = %v, want [1080 1920]", cfg.ScreenSize)
	}

	if !cfg.TextWrap {
		t.Error("DefaultConfig TextWrap should be true")
	}

	if cfg.LineHeight != 1.4 {
		t.Errorf("DefaultConfig LineHeight = %f, want 1.4", cfg.LineHeight)
	}

	// Verify postProcessors from embedded config
	if len(cfg.PostProcessors) != 2 {
		t.Errorf("DefaultConfig PostProcessors should have 2 items, got %d", len(cfg.PostProcessors))
	}
}

func TestDefaultConfigContent_ReturnsValidJSONC(t *testing.T) {
	content := DefaultConfigContent()
	if len(content) == 0 {
		t.Error("DefaultConfigContent should not be empty")
	}

	cfg := &Config{}
	if err := ParseJSONC(content, cfg); err != nil {
		t.Errorf("DefaultConfigContent should be valid JSONC: %v", err)
	}

	if cfg.FontSize != 32 {
		t.Errorf("Parsed config FontSize = %d, want 32", cfg.FontSize)
	}

	if len(cfg.PostProcessors) != 2 {
		t.Errorf("Parsed config PostProcessors should have 2 items, got %d", len(cfg.PostProcessors))
	}
}
