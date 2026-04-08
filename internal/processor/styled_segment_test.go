package processor

import (
	"testing"
)

func TestSingleSegment(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		style     *TextStyle
		wantLen   int
		wantText  string
		wantStyle *TextStyle
	}{
		{
			name:      "with nil style",
			text:      "hello",
			style:     nil,
			wantLen:   1,
			wantText:  "hello",
			wantStyle: nil,
		},
		{
			name:      "with style",
			text:      "world",
			style:     &TextStyle{Bold: true, FontColor: "#FF0000"},
			wantLen:   1,
			wantText:  "world",
			wantStyle: &TextStyle{Bold: true, FontColor: "#FF0000"},
		},
		{
			name:      "empty text",
			text:      "",
			style:     nil,
			wantLen:   1,
			wantText:  "",
			wantStyle: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SingleSegment(tt.text, tt.style)

			if len(result.Segments) != tt.wantLen {
				t.Errorf("SingleSegment() returned %d segments, expected %d", len(result.Segments), tt.wantLen)
				return
			}

			if result.Segments[0].Text != tt.wantText {
				t.Errorf("SingleSegment() text = %q, expected %q", result.Segments[0].Text, tt.wantText)
			}

			if tt.wantStyle == nil {
				if result.Segments[0].Style != nil {
					t.Errorf("SingleSegment() style = %v, expected nil", result.Segments[0].Style)
				}
			} else {
				if result.Segments[0].Style == nil {
					t.Errorf("SingleSegment() style = nil, expected non-nil")
					return
				}
				if result.Segments[0].Style.Bold != tt.wantStyle.Bold {
					t.Errorf("SingleSegment() Bold = %v, expected %v", result.Segments[0].Style.Bold, tt.wantStyle.Bold)
				}
				if result.Segments[0].Style.FontColor != tt.wantStyle.FontColor {
					t.Errorf("SingleSegment() FontColor = %q, expected %q", result.Segments[0].Style.FontColor, tt.wantStyle.FontColor)
				}
			}
		})
	}
}

func TestPlainText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantLen  int
		wantText string
	}{
		{
			name:     "simple text",
			text:     "hello world",
			wantLen:  1,
			wantText: "hello world",
		},
		{
			name:     "empty text",
			text:     "",
			wantLen:  1,
			wantText: "",
		},
		{
			name:     "multiline text",
			text:     "line1\nline2",
			wantLen:  1,
			wantText: "line1\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PlainText(tt.text)

			if len(result.Segments) != tt.wantLen {
				t.Errorf("PlainText() returned %d segments, expected %d", len(result.Segments), tt.wantLen)
				return
			}

			if result.Segments[0].Text != tt.wantText {
				t.Errorf("PlainText() text = %q, expected %q", result.Segments[0].Text, tt.wantText)
			}

			if result.Segments[0].Style != nil {
				t.Errorf("PlainText() style = %v, expected nil", result.Segments[0].Style)
			}
		})
	}
}

func TestStyledSegmentsToText(t *testing.T) {
	tests := []struct {
		name     string
		segments []StyledSegment
		want     string
	}{
		{
			name:     "empty segments",
			segments: []StyledSegment{},
			want:     "",
		},
		{
			name: "single segment",
			segments: []StyledSegment{
				{Text: "hello", Style: nil},
			},
			want: "hello",
		},
		{
			name: "multiple segments",
			segments: []StyledSegment{
				{Text: "echo", Style: &TextStyle{FontColor: "#66D9EF"}},
				{Text: " ", Style: nil},
				{Text: "hello", Style: &TextStyle{FontColor: "#E6DB74"}},
			},
			want: "echo hello",
		},
		{
			name: "segments with styles",
			segments: []StyledSegment{
				{Text: "func", Style: &TextStyle{Bold: true, FontColor: "#66D9EF"}},
				{Text: " ", Style: nil},
				{Text: "main", Style: &TextStyle{FontColor: "#A6E22E"}},
				{Text: "()", Style: &TextStyle{FontColor: "#F8F8F2"}},
			},
			want: "func main()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StyledSegmentsToText(tt.segments)
			if result != tt.want {
				t.Errorf("StyledSegmentsToText() = %q, expected %q", result, tt.want)
			}
		})
	}
}

func TestHasStyledSegments(t *testing.T) {
	tests := []struct {
		name string
		line StyledLine
		want bool
	}{
		{
			name: "plain text no style",
			line: PlainText("hello world"),
			want: false,
		},
		{
			name: "empty segments",
			line: StyledLine{Segments: []StyledSegment{}},
			want: false,
		},
		{
			name: "single segment with font color",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "hello", Style: &TextStyle{FontColor: "#FF0000"}},
				},
			},
			want: true,
		},
		{
			name: "single segment with bold",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "hello", Style: &TextStyle{Bold: true}},
				},
			},
			want: true,
		},
		{
			name: "single segment with italic",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "hello", Style: &TextStyle{Italic: true}},
				},
			},
			want: true,
		},
		{
			name: "single segment with underline",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "hello", Style: &TextStyle{Underline: true}},
				},
			},
			want: true,
		},
		{
			name: "single segment with size",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "hello", Style: &TextStyle{Size: func() *int { i := 18; return &i }()}},
				},
			},
			want: true,
		},
		{
			name: "multiple segments one styled",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "echo", Style: &TextStyle{FontColor: "#66D9EF"}},
					{Text: " ", Style: nil},
					{Text: "hello", Style: nil},
				},
			},
			want: true,
		},
		{
			name: "multiple segments none styled",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "echo", Style: nil},
					{Text: " ", Style: nil},
					{Text: "hello", Style: nil},
				},
			},
			want: false,
		},
		{
			name: "style with empty fields",
			line: StyledLine{
				Segments: []StyledSegment{
					{Text: "hello", Style: &TextStyle{}},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasStyledSegments(tt.line)
			if result != tt.want {
				t.Errorf("HasStyledSegments() = %v, expected %v", result, tt.want)
			}
		})
	}
}
