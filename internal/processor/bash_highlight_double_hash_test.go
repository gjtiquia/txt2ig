package processor

import (
	"strings"
	"testing"
)

func TestBashCodeHighlighter_ChromaColorsNotDoubleHashed(t *testing.T) {
	p := &BashCodeHighlighter{
		StyleName:    "monokai",
		DefaultColor: "#FFFFFF",
	}

	input := []StyledLine{
		PlainText("```bash"),
		PlainText(`echo "hello"`),
		PlainText("```"),
	}

	result, err := p.ProcessLines(input)
	if err != nil {
		t.Fatalf("ProcessLines() error = %v", err)
	}

	echoLine := result[1]
	if len(echoLine.Segments) == 0 {
		t.Fatalf("Expected segments in echo line, got 0")
	}

	// Find the "echo" token segment
	for _, seg := range echoLine.Segments {
		if strings.TrimSpace(seg.Text) == "" {
			continue // Skip whitespace
		}

		if seg.Style == nil {
			t.Errorf("Segment %q has nil style", seg.Text)
			continue
		}

		if seg.Style.FontColor == "" {
			t.Errorf("Segment %q has empty FontColor", seg.Text)
			continue
		}

		// Check that FontColor does NOT start with double hash
		if strings.HasPrefix(seg.Style.FontColor, "##") {
			t.Errorf("Segment %q has double-hashed FontColor: %s (should be single hash like #f8f8f2)",
				seg.Text, seg.Style.FontColor)
		}

		// Check that FontColor starts with single hash
		if !strings.HasPrefix(seg.Style.FontColor, "#") {
			t.Errorf("Segment %q FontColor should start with #, got: %s", seg.Text, seg.Style.FontColor)
		}

		// Check that FontColor is not default (means Chroma color was applied)
		if seg.Text == "echo" && seg.Style.FontColor == "#FFFFFF" {
			t.Errorf("Segment %q has default color #FFFFFF (Chroma color should be applied, not default)", seg.Text)
		}
	}
}
