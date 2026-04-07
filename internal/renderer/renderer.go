package renderer

import (
	"fmt"
	"os"

	"github.com/gjtiquia/txt2ig/internal/config"
	txtfont "github.com/gjtiquia/txt2ig/internal/font"

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

	// 2. Load font
	face, err := r.fontManager.LoadFontWithFallback(r.config.Font, float64(r.config.FontSize), 72)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	// 3. Create image renderer
	imgRenderer, err := NewImageRenderer(r.config.ScreenSize[0], r.config.ScreenSize[1], r.config.BgColor)
	if err != nil {
		return fmt.Errorf("create image renderer: %w", err)
	}

	// 4. Create canvas
	img := imgRenderer.CreateCanvas()

	// 5. Create text renderer
	textRenderer, err := NewTextRenderer(face, r.config)
	if err != nil {
		return fmt.Errorf("create text renderer: %w", err)
	}

	// 6. Wrap text
	maxWidth := r.config.TextBoxMaxWidth
	if maxWidth == 0 {
		maxWidth = int(float64(r.config.ScreenSize[0]) * 0.9)
	}
	lines := textRenderer.WrapText(string(text), maxWidth)

	// 7. Draw text
	if err := textRenderer.DrawText(img, lines); err != nil {
		return fmt.Errorf("draw text: %w", err)
	}

	// 8. Save image
	if err := SaveImage(img, outputPath); err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	return nil
}

func (r *Renderer) RenderWithProcessors(inputPath, outputPath string, preProcessors, postProcessors []interface{}) error {
	// TODO: Implement processor pipeline
	// For now, just call the basic render
	return r.Render(inputPath, outputPath)
}

func GetDefaultOutputPath(inputPath string) string {
	// If input is my-post.md, output will be my-post.jpg
	if len(inputPath) > 3 {
		ext := inputPath[len(inputPath)-3:]
		if ext == ".md" || ext == "txt" {
			return inputPath[:len(inputPath)-3] + "jpg"
		}
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
