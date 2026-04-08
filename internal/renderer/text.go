package renderer

import (
	"image"
	"image/color"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/gjtiquia/txt2ig/internal/config"
	txtfont "github.com/gjtiquia/txt2ig/internal/font"
	"github.com/gjtiquia/txt2ig/internal/processor"
)

type TextRenderer struct {
	fontFamily *txtfont.FontFamily
	config     *config.Config
	fontColor  color.RGBA
	lineHeight int
}

func NewTextRendererWithFamily(fontFamily *txtfont.FontFamily, cfg *config.Config) (*TextRenderer, error) {
	fontColor, err := txtfont.ParseColor(cfg.FontColor)
	if err != nil {
		return nil, err
	}

	lineHeight := txtfont.CalculateLineHeight(fontFamily.Regular, cfg.LineHeight)

	return &TextRenderer{
		fontFamily: fontFamily,
		config:     cfg,
		fontColor:  fontColor,
		lineHeight: lineHeight,
	}, nil
}

func NewTextRenderer(face font.Face, cfg *config.Config) (*TextRenderer, error) {
	fontColor, err := txtfont.ParseColor(cfg.FontColor)
	if err != nil {
		return nil, err
	}

	lineHeight := txtfont.CalculateLineHeight(face, cfg.LineHeight)

	return &TextRenderer{
		fontFamily: &txtfont.FontFamily{Regular: face, Bold: face, Italic: face, BoldItalic: face},
		config:     cfg,
		fontColor:  fontColor,
		lineHeight: lineHeight,
	}, nil
}

func (r *TextRenderer) WrapText(text string, maxWidth int) []string {
	paragraphs := strings.Split(text, "\n")

	if !r.config.TextWrap {
		return paragraphs
	}

	var result []string

	for _, paragraph := range paragraphs {
		if paragraph == "" {
			result = append(result, "")
			continue
		}

		words := strings.Fields(paragraph)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		var currentLine strings.Builder
		for _, word := range words {
			testLine := currentLine.String()
			if testLine != "" {
				testLine += " " + word
			} else {
				testLine = word
			}

			lineWidth := txtfont.MeasureTextWidth(r.fontFamily.Regular, testLine)
			if lineWidth <= maxWidth || currentLine.Len() == 0 {
				if currentLine.Len() > 0 {
					currentLine.WriteString(" ")
				}
				currentLine.WriteString(word)
			} else {
				if currentLine.Len() > 0 {
					result = append(result, currentLine.String())
					currentLine.Reset()
				}
				currentLine.WriteString(word)
			}
		}

		if currentLine.Len() > 0 {
			result = append(result, currentLine.String())
		}
	}

	return result
}

func (r *TextRenderer) CalculateTextBoxSize(lines []string) (width, height int) {
	maxWidth := 0
	for _, line := range lines {
		w := txtfont.MeasureTextWidth(r.fontFamily.Regular, line)
		if w > maxWidth {
			maxWidth = w
		}
	}

	height = len(lines) * r.lineHeight

	return maxWidth, height
}

func (r *TextRenderer) CalculateTextBoxPosition(textWidth, textHeight int) (x, y int) {
	screenWidth := r.config.ScreenSize[0]
	screenHeight := r.config.ScreenSize[1]

	// Calculate horizontal position
	switch r.config.TextBoxJustify {
	case "left":
		x = r.config.TextBoxOffset[0]
	case "right":
		x = screenWidth - textWidth - r.config.TextBoxOffset[0]
	case "center":
		x = (screenWidth-textWidth)/2 + r.config.TextBoxOffset[0]
	default:
		x = (screenWidth-textWidth)/2 + r.config.TextBoxOffset[0]
	}

	// Calculate vertical position
	switch r.config.TextBoxAlign {
	case "top":
		y = r.config.TextBoxOffset[1]
	case "bottom":
		y = screenHeight - textHeight - r.config.TextBoxOffset[1]
	case "center":
		y = (screenHeight-textHeight)/2 + r.config.TextBoxOffset[1]
	default:
		y = (screenHeight-textHeight)/2 + r.config.TextBoxOffset[1]
	}

	return x, y
}

func (r *TextRenderer) selectFace(style *processor.TextStyle) font.Face {
	if style == nil {
		return r.fontFamily.Regular
	}

	if style.Bold && style.Italic {
		return r.fontFamily.BoldItalic
	}
	if style.Bold {
		return r.fontFamily.Bold
	}
	if style.Italic {
		return r.fontFamily.Italic
	}

	return r.fontFamily.Regular
}

func (r *TextRenderer) DrawText(img *image.RGBA, styledLines []processor.StyledLine) error {
	lines := make([]string, len(styledLines))
	for i, sl := range styledLines {
		lines[i] = sl.Text
	}

	textWidth, textHeight := r.CalculateTextBoxSize(lines)
	startX, startY := r.CalculateTextBoxPosition(textWidth, textHeight)

	for i, sl := range styledLines {
		y := startY + (i+1)*r.lineHeight - r.lineHeight/4

		face := r.selectFace(sl.Style)
		lineWidth := txtfont.MeasureTextWidth(face, sl.Text)

		var x int
		switch r.config.TextJustify {
		case "left":
			x = startX
		case "right":
			x = startX + (textWidth - lineWidth)
		case "center":
			x = startX + (textWidth-lineWidth)/2
		default:
			x = startX
		}

		lineColor := r.fontColor
		if sl.Style != nil && sl.Style.FontColor != "" {
			customColor, err := txtfont.ParseColor(sl.Style.FontColor)
			if err == nil {
				lineColor = customColor
			}
		}

		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(lineColor),
			Face: face,
			Dot: fixed.Point26_6{
				X: fixed.Int26_6(x << 6),
				Y: fixed.Int26_6(y << 6),
			},
		}
		drawer.DrawString(sl.Text)
	}

	return nil
}
