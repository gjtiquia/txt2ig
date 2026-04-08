package font

import (
	"fmt"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type FontFamily struct {
	Regular    font.Face
	Bold       font.Face
	Italic     font.Face
	BoldItalic font.Face
}

type FontFamilyConfig struct {
	Regular    []string
	Bold       []string
	Italic     []string
	BoldItalic []string
}

func (m *Manager) LoadFontFamily(familyConfig FontFamilyConfig, size float64, dpi float64) (*FontFamily, error) {
	regular, err := m.loadFontWithFallback(familyConfig.Regular, size, dpi, "GoMono")
	if err != nil {
		return nil, fmt.Errorf("load regular font: %w", err)
	}

	bold, err := m.loadFontWithFallback(familyConfig.Bold, size, dpi, "GoMonoBold")
	if err != nil {
		return nil, fmt.Errorf("load bold font: %w", err)
	}

	italic, err := m.loadFontWithFallback(familyConfig.Italic, size, dpi, "GoMonoItalic")
	if err != nil {
		return nil, fmt.Errorf("load italic font: %w", err)
	}

	boldItalic, err := m.loadFontWithFallback(familyConfig.BoldItalic, size, dpi, "GoMonoBoldItalic")
	if err != nil {
		return nil, fmt.Errorf("load bold-italic font: %w", err)
	}

	return &FontFamily{
		Regular:    regular,
		Bold:       bold,
		Italic:     italic,
		BoldItalic: boldItalic,
	}, nil
}

func (m *Manager) loadFontWithFallback(fontNames []string, size float64, dpi float64, embeddedFallback string) (font.Face, error) {
	for _, fontName := range fontNames {
		f, err := m.loadFont(fontName)
		if err == nil {
			face, err := opentype.NewFace(f, &opentype.FaceOptions{
				Size:    size,
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
			if err != nil {
				log.Printf("warning: create font face for %s: %v", fontName, err)
				continue
			}
			return face, nil
		}
		log.Printf("warning: failed to load font %s: %v", fontName, err)
	}

	log.Printf("warning: using embedded %s as fallback", embeddedFallback)
	return m.loadEmbeddedFont(embeddedFallback, size, dpi)
}
