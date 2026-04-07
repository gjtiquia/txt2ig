package font

import (
	"image/color"
	"strconv"

	"golang.org/x/image/font"
)

func ParseColor(hex string) (color.RGBA, error) {
	if len(hex) == 0 || hex[0] != '#' {
		return color.RGBA{}, ParseColorError{hex, "must start with #"}
	}

	hex = hex[1:]

	var r, g, b uint8
	var a uint8 = 255

	switch len(hex) {
	case 3: // #RGB
		r = parseHexDigit(hex[0])
		r = (r << 4) | r
		g = parseHexDigit(hex[1])
		g = (g << 4) | g
		b = parseHexDigit(hex[2])
		b = (b << 4) | b
	case 4: // #RGBA
		r = parseHexDigit(hex[0])
		r = (r << 4) | r
		g = parseHexDigit(hex[1])
		g = (g << 4) | g
		b = parseHexDigit(hex[2])
		b = (b << 4) | b
		a = parseHexDigit(hex[3])
		a = (a << 4) | a
	case 6: // #RRGGBB
		r = parseHexByte(hex[0:2])
		g = parseHexByte(hex[2:4])
		b = parseHexByte(hex[4:6])
	case 8: // #RRGGBBAA
		r = parseHexByte(hex[0:2])
		g = parseHexByte(hex[2:4])
		b = parseHexByte(hex[4:6])
		a = parseHexByte(hex[6:8])
	default:
		return color.RGBA{}, ParseColorError{hex, "invalid length"}
	}

	return color.RGBA{R: r, G: g, B: b, A: a}, nil
}

func parseHexDigit(c byte) uint8 {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}

func parseHexByte(s string) uint8 {
	v, _ := strconv.ParseUint(s, 16, 8)
	return uint8(v)
}

type ParseColorError struct {
	Input string
	Msg   string
}

func (e ParseColorError) Error() string {
	return "parse color '" + e.Input + "': " + e.Msg
}

func MeasureTextWidth(face font.Face, text string) int {
	width := 0
	for _, r := range text {
		advance, _ := face.GlyphAdvance(r)
		width += advance.Ceil()
	}
	return width
}

func CalculateLineHeight(face font.Face, multiplier float64) int {
	metrics := face.Metrics()
	lineHeight := float64(metrics.Height.Ceil()) * multiplier
	return int(lineHeight)
}
