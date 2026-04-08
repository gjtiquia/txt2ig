package renderer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/gjtiquia/txt2ig/internal/config"
	txtfont "github.com/gjtiquia/txt2ig/internal/font"
	"github.com/gjtiquia/txt2ig/internal/processor"

	"golang.org/x/image/font"
)

type Renderer struct {
	config      *config.Config
	fontManager *txtfont.Manager
}

func NewRenderer(cfg *config.Config) *Renderer {
	return &Renderer{
		config:      cfg,
		fontManager: txtfont.NewManager(),
	}
}

func (r *Renderer) Render(inputPath, outputPath string) error {
	// 1. Read input file
	text, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file %s: %w", inputPath, err)
	}

	// 2. Apply pre-processors
	processedText, err := processor.ApplyPreProcessors(string(text), r.config.PreProcessors)
	if err != nil {
		return fmt.Errorf("apply pre-processors: %w", err)
	}

	// 3. Load font
	fontFamily, err := r.fontManager.LoadFontFamily(txtfont.FontFamilyConfig{
		Regular:    r.config.FontFamily.Regular,
		Bold:       r.config.FontFamily.Bold,
		Italic:     r.config.FontFamily.Italic,
		BoldItalic: r.config.FontFamily.BoldItalic,
	}, float64(r.config.FontSize), 72)
	if err != nil {
		return fmt.Errorf("load font family: %w", err)
	}

	// 4. Create image renderer
	imgRenderer, err := NewImageRenderer(r.config.ScreenSize[0], r.config.ScreenSize[1], r.config.BgColor)
	if err != nil {
		return fmt.Errorf("create image renderer: %w", err)
	}

	// 5. Create canvas
	img := imgRenderer.CreateCanvas()

	// 6. Create text renderer
	textRenderer, err := NewTextRendererWithFamily(fontFamily, r.config)
	if err != nil {
		return fmt.Errorf("create text renderer: %w", err)
	}

	// 7. Wrap text
	maxWidth := r.config.TextBoxMaxWidth
	if maxWidth == 0 {
		maxWidth = int(float64(r.config.ScreenSize[0]) * 0.9)
	}
	lines := textRenderer.WrapText(processedText, maxWidth)

	// 8. Apply post-processors
	processedLines, err := processor.ApplyPostProcessors(lines, r.config.PostProcessors)
	if err != nil {
		return fmt.Errorf("apply post-processors: %w", err)
	}

	// 9. Draw text
	if err := textRenderer.DrawText(img, processedLines); err != nil {
		return fmt.Errorf("draw text: %w", err)
	}

	// 10. Save image
	if err := SaveImage(img, outputPath); err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	return nil
}

func (r *Renderer) RenderWithProcessors(inputPath, outputPath string, preProcessors, postProcessors []interface{}) error {
	// Override config processors if provided
	if len(preProcessors) > 0 {
		oldPre := r.config.PreProcessors
		r.config.PreProcessors = preProcessors
		defer func() { r.config.PreProcessors = oldPre }()
	}
	if len(postProcessors) > 0 {
		oldPost := r.config.PostProcessors
		r.config.PostProcessors = postProcessors
		defer func() { r.config.PostProcessors = oldPost }()
	}

	return r.Render(inputPath, outputPath)
}

func GetDefaultOutputPath(inputPath string) string {
	// If input is my-post.md, output will be my-post.jpg
	// Handle .md and .txt extensions
	if len(inputPath) > 3 && inputPath[len(inputPath)-3:] == ".md" {
		return inputPath[:len(inputPath)-3] + ".jpg"
	}
	if len(inputPath) > 4 && inputPath[len(inputPath)-4:] == ".txt" {
		return inputPath[:len(inputPath)-4] + ".jpg"
	}
	return inputPath + ".jpg"
}

func DetermineOutputPath(inputPath, outputFlag string) string {
	if outputFlag != "" {
		return outputFlag
	}
	return GetDefaultOutputPath(inputPath)
}

// Close releases resources
func (r *Renderer) Close() error {
	// Face resources are managed by font.Manager
	// No explicit cleanup needed for now
	return nil
}

// Ensure Renderer implements cleanup
var _ interface{ Close() error } = (*Renderer)(nil)

// Helper function (not exported)
func mustNewFace(f *txtfont.Manager, names []string, size float64) font.Face {
	face, err := f.LoadFontWithFallback(names, size, 72)
	if err != nil {
		panic(err)
	}
	return face
}

func (r *Renderer) RenderString(text string) (*image.RGBA, error) {
	processedText, err := processor.ApplyPreProcessors(text, r.config.PreProcessors)
	if err != nil {
		return nil, fmt.Errorf("apply pre-processors: %w", err)
	}

	fontFamily, err := r.fontManager.LoadFontFamily(txtfont.FontFamilyConfig{
		Regular:    r.config.FontFamily.Regular,
		Bold:       r.config.FontFamily.Bold,
		Italic:     r.config.FontFamily.Italic,
		BoldItalic: r.config.FontFamily.BoldItalic,
	}, float64(r.config.FontSize), 72)
	if err != nil {
		return nil, fmt.Errorf("load font family: %w", err)
	}

	imgRenderer, err := NewImageRenderer(r.config.ScreenSize[0], r.config.ScreenSize[1], r.config.BgColor)
	if err != nil {
		return nil, fmt.Errorf("create image renderer: %w", err)
	}

	img := imgRenderer.CreateCanvas()

	textRenderer, err := NewTextRendererWithFamily(fontFamily, r.config)
	if err != nil {
		return nil, fmt.Errorf("create text renderer: %w", err)
	}

	maxWidth := r.config.TextBoxMaxWidth
	if maxWidth == 0 {
		maxWidth = int(float64(r.config.ScreenSize[0]) * 0.9)
	}
	lines := textRenderer.WrapText(processedText, maxWidth)

	processedLines, err := processor.ApplyPostProcessors(lines, r.config.PostProcessors)
	if err != nil {
		return nil, fmt.Errorf("apply post-processors: %w", err)
	}

	if err := textRenderer.DrawText(img, processedLines); err != nil {
		return nil, fmt.Errorf("draw text: %w", err)
	}

	return img, nil
}

func EncodeImage(img *image.RGBA, format string) (string, error) {
	var buf bytes.Buffer
	switch format {
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return "", fmt.Errorf("encode png: %w", err)
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	case "jpg", "jpeg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
		if err != nil {
			return "", fmt.Errorf("encode jpeg: %w", err)
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
