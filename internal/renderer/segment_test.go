package renderer

import (
	"image"
	"testing"

	"github.com/gjtiquia/txt2ig/internal/config"
	txtfont "github.com/gjtiquia/txt2ig/internal/font"
	"github.com/gjtiquia/txt2ig/internal/processor"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

func TestCalculateLineSegmentsWidth(t *testing.T) {
	cfg := &config.Config{
		TextWrap:   true,
		LineHeight: 1.4,
		FontColor:  "#FFFFFF",
	}

	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatalf("Failed to create font face: %v", err)
	}

	renderer, err := NewTextRenderer(face, cfg)
	if err != nil {
		t.Fatalf("Failed to create text renderer: %v", err)
	}

	tests := []struct {
		name     string
		segments []processor.StyledSegment
	}{
		{
			name:     "empty segments",
			segments: []processor.StyledSegment{},
		},
		{
			name: "single segment plain",
			segments: []processor.StyledSegment{
				{Text: "hello", Style: nil},
			},
		},
		{
			name: "single segment bold",
			segments: []processor.StyledSegment{
				{Text: "hello", Style: &processor.TextStyle{Bold: true}},
			},
		},
		{
			name: "multiple segments",
			segments: []processor.StyledSegment{
				{Text: "echo", Style: &processor.TextStyle{FontColor: "#66D9EF"}},
				{Text: " ", Style: nil},
				{Text: "world", Style: &processor.TextStyle{FontColor: "#E6DB74"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := renderer.calculateLineSegmentsWidth(tt.segments)

			expectedWidth := 0
			for _, seg := range tt.segments {
				face := renderer.selectFace(seg.Style)
				expectedWidth += txtfont.MeasureTextWidth(face, seg.Text)
			}

			if width != expectedWidth {
				t.Errorf("calculateLineSegmentsWidth() = %d, expected %d", width, expectedWidth)
			}
		})
	}
}

func TestDrawTextWithSegments(t *testing.T) {
	cfg := &config.Config{
		ScreenSize:     []int{800, 600},
		TextBoxOffset:  []int{50, 50},
		TextBoxJustify: "center",
		TextBoxAlign:   "center",
		TextJustify:    "left",
		LineHeight:     1.4,
		FontColor:      "#FFFFFF",
	}

	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatalf("Failed to create font face: %v", err)
	}

	renderer, err := NewTextRenderer(face, cfg)
	if err != nil {
		t.Fatalf("Failed to create text renderer: %v", err)
	}

	tests := []struct {
		name  string
		lines []processor.StyledLine
	}{
		{
			name: "single plain line",
			lines: []processor.StyledLine{
				processor.PlainText("hello world"),
			},
		},
		{
			name: "single line with colored segment",
			lines: []processor.StyledLine{
				{
					Segments: []processor.StyledSegment{
						{Text: "echo", Style: &processor.TextStyle{FontColor: "#66D9EF"}},
						{Text: " ", Style: nil},
						{Text: "hello", Style: &processor.TextStyle{FontColor: "#E6DB74"}},
					},
				},
			},
		},
		{
			name: "multiple lines with mixed styles",
			lines: []processor.StyledLine{
				{
					Segments: []processor.StyledSegment{
						{Text: "#!/bin/bash", Style: &processor.TextStyle{FontColor: "#75715E"}},
					},
				},
				processor.PlainText(""),
				{
					Segments: []processor.StyledSegment{
						{Text: "echo", Style: &processor.TextStyle{FontColor: "#66D9EF"}},
						{Text: " ", Style: nil},
						{Text: "\"hello\"", Style: &processor.TextStyle{FontColor: "#E6DB74"}},
					},
				},
			},
		},
		{
			name: "bold and italic segments",
			lines: []processor.StyledLine{
				{
					Segments: []processor.StyledSegment{
						{Text: "Title", Style: &processor.TextStyle{Bold: true, FontColor: "#FF0000"}},
					},
				},
				{
					Segments: []processor.StyledSegment{
						{Text: "Italic text", Style: &processor.TextStyle{Italic: true}},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := image.NewRGBA(image.Rect(0, 0, 800, 600))

			err := renderer.DrawText(img, tt.lines)
			if err != nil {
				t.Errorf("DrawText() error = %v", err)
			}
		})
	}
}

func TestSegmentsWithDifferentColors(t *testing.T) {
	cfg := &config.Config{
		ScreenSize:     []int{800, 600},
		TextBoxOffset:  []int{50, 50},
		TextBoxJustify: "center",
		TextBoxAlign:   "center",
		TextJustify:    "left",
		LineHeight:     1.4,
		FontColor:      "#FFFFFF",
	}

	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatalf("Failed to create font face: %v", err)
	}

	renderer, err := NewTextRenderer(face, cfg)
	if err != nil {
		t.Fatalf("Failed to create text renderer: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	lines := []processor.StyledLine{
		{
			Segments: []processor.StyledSegment{
				{Text: "func", Style: &processor.TextStyle{FontColor: "#66D9EF"}},
				{Text: " ", Style: nil},
				{Text: "main", Style: &processor.TextStyle{FontColor: "#A6E22E"}},
				{Text: "()", Style: &processor.TextStyle{FontColor: "#F8F8F2"}},
				{Text: " ", Style: nil},
				{Text: "{", Style: &processor.TextStyle{FontColor: "#F8F8F2"}},
			},
		},
	}

	err = renderer.DrawText(img, lines)
	if err != nil {
		t.Errorf("DrawText() error = %v", err)
	}
}

func TestSegmentsWithDifferentFontFaces(t *testing.T) {
	cfg := &config.Config{
		ScreenSize:     []int{800, 600},
		TextBoxOffset:  []int{50, 50},
		TextBoxJustify: "center",
		TextBoxAlign:   "center",
		TextJustify:    "left",
		LineHeight:     1.4,
		FontColor:      "#FFFFFF",
	}

	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		t.Fatalf("Failed to create font face: %v", err)
	}

	renderer, err := NewTextRenderer(face, cfg)
	if err != nil {
		t.Fatalf("Failed to create text renderer: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	lines := []processor.StyledLine{
		{
			Segments: []processor.StyledSegment{
				{Text: "Bold", Style: &processor.TextStyle{Bold: true}},
				{Text: " ", Style: nil},
				{Text: "Italic", Style: &processor.TextStyle{Italic: true}},
				{Text: " ", Style: nil},
				{Text: "BoldItalic", Style: &processor.TextStyle{Bold: true, Italic: true}},
			},
		},
	}

	err = renderer.DrawText(img, lines)
	if err != nil {
		t.Errorf("DrawText() error = %v", err)
	}
}
