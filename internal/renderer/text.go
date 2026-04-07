package renderer

import (
	"image"
	"image/color"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/gjtiquia/txt2ig/internal/config"
	txtfont "github.com/gjtiquia/txt2ig/internal/font"
)

type TextRenderer struct {
	face       font.Face
	config     *config.Config
	fontColor  color.RGBA
	lineHeight int
}

func NewTextRenderer(face font.Face, cfg *config.Config) (*TextRenderer, error) {
	fontColor, err := txtfont.ParseColor(cfg.FontColor)
	if err != nil {
		return nil, err
	}

	lineHeight := txtfont.CalculateLineHeight(face, cfg.LineHeight)

	return &TextRenderer{
		face:       face,
		config:     cfg,
		fontColor:  fontColor,
		lineHeight: lineHeight,
	}, nil
}

func (r *TextRenderer) WrapText(text string, maxWidth int) []string {
	if !r.config.TextWrap {
		return strings.Split(text, "\n")
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " " + word
		} else {
			testLine = word
		}

		lineWidth := txtfont.MeasureTextWidth(r.face, testLine)
		if lineWidth <= maxWidth || currentLine.Len() == 0 {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

func (r *TextRenderer) CalculateTextBoxSize(lines []string) (width, height int) {
	maxWidth := 0
	for _, line := range lines {
		w := txtfont.MeasureTextWidth(r.face, line)
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

func (r *TextRenderer) DrawText(img *image.RGBA, lines []string) error {
	textWidth, textHeight := r.CalculateTextBoxSize(lines)
	startX, startY := r.CalculateTextBoxPosition(textWidth, textHeight)

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(r.fontColor),
		Face: r.face,
	}

	for i, line := range lines {
		y := startY + (i+1)*r.lineHeight - r.lineHeight/4

		lineWidth := txtfont.MeasureTextWidth(r.face, line)

		// Apply text justification within the text box
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

		drawer.Dot = fixed.Point26_6{
			X: fixed.Int26_6(x << 6),
			Y: fixed.Int26_6(y << 6),
		}
		drawer.DrawString(line)
	}

	return nil
}
