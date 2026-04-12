package processor

import (
	"testing"
)

func TestMarkdownBoldHeaders(t *testing.T) {
	tests := []struct {
		name          string
		bold          bool
		fontColor     string
		input         string
		expectedLine  string
		expectedBold  bool
		expectedColor string
	}{
		{
			name:          "single hash header",
			bold:          true,
			fontColor:     "#FF0000",
			input:         "# Title",
			expectedLine:  "# Title", // # is kept
			expectedBold:  true,
			expectedColor: "#FF0000",
		},
		{
			name:          "double hash header",
			bold:          true,
			fontColor:     "",
			input:         "## Heading",
			expectedLine:  "## Heading", // # is kept
			expectedBold:  true,
			expectedColor: "",
		},
		{
			name:          "triple hash header",
			bold:          false,
			fontColor:     "#00FF00",
			input:         "### Subheading",
			expectedLine:  "### Subheading", // # is kept
			expectedBold:  false,
			expectedColor: "#00FF00",
		},
		{
			name:          "not a header",
			bold:          true,
			fontColor:     "#FF0000",
			input:         "Normal text",
			expectedLine:  "Normal text",
			expectedBold:  false,
			expectedColor: "",
		},
		{
			name:          "header with spaces after hash",
			bold:          true,
			fontColor:     "",
			input:         "#   Heading with spaces",
			expectedLine:  "#   Heading with spaces", // # is kept
			expectedBold:  true,
			expectedColor: "",
		},
		{
			name:          "just hash",
			bold:          true,
			fontColor:     "#000000",
			input:         "#",
			expectedLine:  "#",
			expectedBold:  true,
			expectedColor: "#000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MarkdownBoldHeaders{
				Bold:      tt.bold,
				FontColor: tt.fontColor,
			}
			line, style, err := p.Process(tt.input)
			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}
			if line != tt.expectedLine {
				t.Errorf("Process() line = %q, want %q", line, tt.expectedLine)
			}
			if style != nil {
				if style.Bold != tt.expectedBold {
					t.Errorf("Style.Bold = %v, want %v", style.Bold, tt.expectedBold)
				}
				if style.FontColor != tt.expectedColor {
					t.Errorf("Style.FontColor = %q, want %q", style.FontColor, tt.expectedColor)
				}
			} else if tt.expectedBold || tt.expectedColor != "" {
				t.Errorf("Style is nil, expected non-nil")
			}
		})
	}
}

func TestBashComments(t *testing.T) {
	tests := []struct {
		name           string
		italic         bool
		fontColor      string
		input          string
		expectedLine   string
		expectedItalic bool
		expectedColor  string
	}{
		{
			name:           "bash comment",
			italic:         true,
			fontColor:      "#888888",
			input:          "# This is a comment",
			expectedLine:   "# This is a comment", // Bash comments preserve the #
			expectedItalic: true,
			expectedColor:  "#888888",
		},
		{
			name:           "not a comment",
			italic:         true,
			fontColor:      "#888888",
			input:          "echo hello",
			expectedLine:   "echo hello",
			expectedItalic: false,
			expectedColor:  "",
		},
		{
			name:           "comment at end",
			italic:         false,
			fontColor:      "#CCCCCC",
			input:          "# comment",
			expectedLine:   "# comment",
			expectedItalic: false,
			expectedColor:  "#CCCCCC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &BashComments{
				Italic:    tt.italic,
				FontColor: tt.fontColor,
			}
			line, style, err := p.Process(tt.input)
			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}
			if line != tt.expectedLine {
				t.Errorf("Process() line = %q, want %q", line, tt.expectedLine)
			}
			if style != nil {
				if style.Italic != tt.expectedItalic {
					t.Errorf("Style.Italic = %v, want %v", style.Italic, tt.expectedItalic)
				}
				if style.FontColor != tt.expectedColor {
					t.Errorf("Style.FontColor = %q, want %q", style.FontColor, tt.expectedColor)
				}
			} else if tt.expectedItalic || tt.expectedColor != "" {
				t.Errorf("Style is nil, expected non-nil")
			}
		})
	}
}

func TestMarkdownLinks_BasicLink(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		PlainText("Check out [my blog](https://example.com) for more"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(result))
	}

	segs := result[0].Segments
	expectedText := "Check out [my blog](https://example.com) for more"
	actualText := ""
	for _, s := range segs {
		actualText += s.Text
	}
	if actualText != expectedText {
		t.Errorf("Text = %q, want %q", actualText, expectedText)
	}

	nameSeg := segs[2]
	if nameSeg.Text != "my blog" {
		t.Errorf("Name segment text = %q, want %q", nameSeg.Text, "my blog")
	}
	if nameSeg.Style == nil || nameSeg.Style.FontColor != "#3B82F6" {
		t.Errorf("Name segment color = %v, want #3B82F6", nameSeg.Style)
	}

	urlSeg := segs[4]
	if urlSeg.Text != "https://example.com" {
		t.Errorf("URL segment text = %q, want %q", urlSeg.Text, "https://example.com")
	}
	if urlSeg.Style == nil || urlSeg.Style.FontColor != "#60A5FA" {
		t.Errorf("URL segment color = %v, want #60A5FA", urlSeg.Style)
	}
}

func TestMarkdownLinks_MultipleLinks(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		PlainText("See [link1](url1) and [link2](url2)"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	segs := result[0].Segments

	if len(segs) < 10 {
		t.Fatalf("Expected at least 10 segments, got %d", len(segs))
	}

	if segs[2].Text != "link1" {
		t.Errorf("First link name = %q, want %q", segs[2].Text, "link1")
	}
	if segs[2].Style == nil || segs[2].Style.FontColor != "#3B82F6" {
		t.Errorf("First link name color = %v, want #3B82F6", segs[2].Style)
	}

	if segs[4].Text != "url1" {
		t.Errorf("First link URL = %q, want %q", segs[4].Text, "url1")
	}
	if segs[4].Style == nil || segs[4].Style.FontColor != "#60A5FA" {
		t.Errorf("First link URL color = %v, want #60A5FA", segs[4].Style)
	}
}

func TestMarkdownLinks_NoLinks(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		PlainText("Just plain text"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	if len(result[0].Segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(result[0].Segments))
	}

	if result[0].Segments[0].Style != nil {
		t.Errorf("Expected no style, got %v", result[0].Segments[0].Style)
	}
}

func TestMarkdownLinks_EmptyColors(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "",
		LinkColor: "",
	}

	lines := []StyledLine{
		PlainText("Check out [my blog](https://example.com)"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	for _, seg := range result[0].Segments {
		if seg.Style != nil && seg.Style.FontColor != "" {
			t.Errorf("Expected no color styling, but segment %q has color %s", seg.Text, seg.Style.FontColor)
		}
	}
}

func TestMarkdownLinks_OnlyNameColor(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "",
	}

	lines := []StyledLine{
		PlainText("[name](url)"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	nameSeg := result[0].Segments[1]
	if nameSeg.Text != "name" {
		t.Errorf("Name text = %q, want %q", nameSeg.Text, "name")
	}
	if nameSeg.Style == nil || nameSeg.Style.FontColor != "#3B82F6" {
		t.Errorf("Name color = %v, want #3B82F6", nameSeg.Style)
	}

	urlSeg := result[0].Segments[3]
	if urlSeg.Text != "url" {
		t.Errorf("URL text = %q, want %q", urlSeg.Text, "url")
	}
	if urlSeg.Style != nil && urlSeg.Style.FontColor != "" {
		t.Errorf("URL should not have color, got %s", urlSeg.Style.FontColor)
	}
}

func TestMarkdownLinks_OnlyLinkColor(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		PlainText("[name](url)"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	nameSeg := result[0].Segments[1]
	if nameSeg.Text != "name" {
		t.Errorf("Name text = %q, want %q", nameSeg.Text, "name")
	}
	if nameSeg.Style != nil && nameSeg.Style.FontColor != "" {
		t.Errorf("Name should not have color, got %s", nameSeg.Style.FontColor)
	}

	urlSeg := result[0].Segments[3]
	if urlSeg.Text != "url" {
		t.Errorf("URL text = %q, want %q", urlSeg.Text, "url")
	}
	if urlSeg.Style == nil || urlSeg.Style.FontColor != "#60A5FA" {
		t.Errorf("URL color = %v, want #60A5FA", urlSeg.Style)
	}
}

func TestMarkdownLinks_MergeWithExistingStyles(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		{
			Segments: []StyledSegment{
				{Text: "# Check out ", Style: &TextStyle{Bold: true}},
				{Text: "[my blog](https://example.com)", Style: &TextStyle{Bold: true}},
			},
		},
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	nameSeg := result[0].Segments[2]
	if nameSeg.Text != "my blog" {
		t.Errorf("Name text = %q, want %q", nameSeg.Text, "my blog")
	}
	if nameSeg.Style == nil {
		t.Error("Name segment style is nil")
	} else {
		if !nameSeg.Style.Bold {
			t.Error("Name should be bold")
		}
		if nameSeg.Style.FontColor != "#3B82F6" {
			t.Errorf("Name color = %s, want #3B82F6", nameSeg.Style.FontColor)
		}
	}
}

func TestMarkdownLinks_EmptyBrackets(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		PlainText("Empty: []() here"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	actualText := ""
	for _, s := range result[0].Segments {
		actualText += s.Text
	}
	expectedText := "Empty: []() here"
	if actualText != expectedText {
		t.Errorf("Text = %q, want %q", actualText, expectedText)
	}
}

func TestMarkdownLinks_LinkAtStart(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		PlainText("[start](url1) middle"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	if result[0].Segments[1].Text != "start" {
		t.Errorf("Name segment = %q, want %q", result[0].Segments[1].Text, "start")
	}
}

func TestMarkdownLinks_LinkAtEnd(t *testing.T) {
	p := &MarkdownLinks{
		NameColor: "#3B82F6",
		LinkColor: "#60A5FA",
	}

	lines := []StyledLine{
		PlainText("middle [end](url2)"),
	}

	result, err := p.ProcessLines(lines)
	if err != nil {
		t.Fatalf("ProcessLines() error: %v", err)
	}

	actualText := ""
	for _, s := range result[0].Segments {
		actualText += s.Text
	}
	expectedText := "middle [end](url2)"
	if actualText != expectedText {
		t.Errorf("Text = %q, want %q", actualText, expectedText)
	}

	found := false
	for _, seg := range result[0].Segments {
		if seg.Text == "url2" {
			found = true
			if seg.Style == nil || seg.Style.FontColor != "#60A5FA" {
				t.Errorf("URL color = %v, want #60A5FA", seg.Style)
			}
		}
	}
	if !found {
		t.Error("URL segment 'url2' not found")
	}
}
