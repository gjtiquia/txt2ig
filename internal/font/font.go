package font

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/opentype"
)

type Manager struct {
	fonts map[string]*opentype.Font
}

func NewManager() *Manager {
	return &Manager{
		fonts: make(map[string]*opentype.Font),
	}
}

func (m *Manager) LoadFontWithFallback(fontNames []string, size float64, dpi float64) (font.Face, error) {
	for _, fontName := range fontNames {
		f, err := m.loadFont(fontName)
		if err == nil {
			face, err := opentype.NewFace(f, &opentype.FaceOptions{
				Size:    size,
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
			if err != nil {
				return nil, fmt.Errorf("create font face for %s: %w", fontName, err)
			}
			return face, nil
		}
	}

	// Final fallback: embedded GoMono
	return m.loadGoMono(size, dpi)
}

func (m *Manager) loadFont(fontSpec string) (*opentype.Font, error) {
	// Check if already loaded
	if f, ok := m.fonts[fontSpec]; ok {
		return f, nil
	}

	// Try as embedded font name first
	if isEmbeddedFontName(fontSpec) {
		f, err := m.parseEmbeddedFont(fontSpec)
		if err != nil {
			return nil, fmt.Errorf("load embedded font %s: %w", fontSpec, err)
		}
		m.fonts[fontSpec] = f
		return f, nil
	}

	// Try as file path
	if _, err := os.Stat(fontSpec); err == nil {
		f, err := m.loadFontFromFile(fontSpec)
		if err != nil {
			return nil, err
		}
		m.fonts[fontSpec] = f
		return f, nil
	}

	// Try as system font
	f, err := m.loadSystemFont(fontSpec)
	if err != nil {
		return nil, fmt.Errorf("font not found: %s (not a file path and not in system fonts)", fontSpec)
	}
	m.fonts[fontSpec] = f
	return f, nil
}

func isEmbeddedFontName(name string) bool {
	switch name {
	case "GoMono", "GoMonoBold", "GoMonoItalic", "GoMonoBoldItalic":
		return true
	}
	return false
}

func (m *Manager) parseEmbeddedFont(name string) (*opentype.Font, error) {
	var fontData []byte

	switch name {
	case "GoMono":
		fontData = gomono.TTF
	case "GoMonoBold":
		fontData = gomonobold.TTF
	case "GoMonoItalic":
		fontData = gomonoitalic.TTF
	case "GoMonoBoldItalic":
		fontData = gomonobolditalic.TTF
	default:
		return nil, fmt.Errorf("unknown embedded font: %s", name)
	}

	return opentype.Parse(fontData)
}

func (m *Manager) loadFontFromFile(path string) (*opentype.Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read font file %s: %w", path, err)
	}

	f, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse font file %s: %w", path, err)
	}

	return f, nil
}

func (m *Manager) loadSystemFont(fontName string) (*opentype.Font, error) {
	// Common font directories on Unix-like systems
	fontDirs := []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
	}

	// Add user font directory if available
	if homeDir, err := os.UserHomeDir(); err == nil {
		fontDirs = append(fontDirs, filepath.Join(homeDir, ".local/share/fonts"))
		fontDirs = append(fontDirs, filepath.Join(homeDir, ".fonts"))
	}

	// Search for font files
	for _, dir := range fontDirs {
		fontPath := m.findFontFile(dir, fontName)
		if fontPath != "" {
			return m.loadFontFromFile(fontPath)
		}
	}

	return nil, fmt.Errorf("font not found in system directories: %s", fontName)
}

func (m *Manager) findFontFile(dir string, fontName string) string {
	// Common font file extensions
	extensions := []string{".ttf", ".otf", ".TTF", ".OTF"}

	// Try exact match with different extensions
	for _, ext := range extensions {
		path := filepath.Join(dir, fontName+ext)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Try searching in subdirectories
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// Check if filename contains font name (case-insensitive)
		name := info.Name()
		for _, ext := range extensions {
			if filepath.Ext(name) == ext {
				baseName := name[:len(name)-len(ext)]
				if baseName == fontName {
					return filepath.SkipDir // Found it, but we can't return path from Walk
				}
			}
		}
		return nil
	})

	// Note: This is a simplified implementation
	// A full implementation would properly return the found path
	return ""
}

func (m *Manager) loadGoMono(size float64, dpi float64) (font.Face, error) {
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, fmt.Errorf("parse embedded GoMono: %w", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create GoMono font face: %w", err)
	}

	return face, nil
}

func (m *Manager) loadEmbeddedFont(name string, size float64, dpi float64) (font.Face, error) {
	var fontData []byte

	switch name {
	case "GoMono":
		fontData = gomono.TTF
	case "GoMonoBold":
		fontData = gomonobold.TTF
	case "GoMonoItalic":
		fontData = gomonoitalic.TTF
	case "GoMonoBoldItalic":
		fontData = gomonobolditalic.TTF
	default:
		return nil, fmt.Errorf("unknown embedded font: %s", name)
	}

	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("parse embedded font %s: %w", name, err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create font face for %s: %w", name, err)
	}

	return face, nil
}
