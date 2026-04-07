package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/gjtiquia/txt2ig/internal/font"
)

type ImageRenderer struct {
	width   int
	height  int
	bgColor color.RGBA
}

func NewImageRenderer(width, height int, bgColor string) (*ImageRenderer, error) {
	bg, err := font.ParseColor(bgColor)
	if err != nil {
		return nil, fmt.Errorf("parse background color: %w", err)
	}

	return &ImageRenderer{
		width:   width,
		height:  height,
		bgColor: bg,
	}, nil
}

func (r *ImageRenderer) CreateCanvas() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			img.Set(x, y, r.bgColor)
		}
	}

	return img
}

func SaveImage(img *image.RGBA, outputPath string) error {
	ext := strings.ToLower(filepath.Ext(outputPath))

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file %s: %w", outputPath, err)
	}
	defer f.Close()

	switch ext {
	case ".png":
		return png.Encode(f, img)
	case ".jpg", ".jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	default:
		return fmt.Errorf("unsupported output format: %s (use .jpg or .png)", ext)
	}
}
