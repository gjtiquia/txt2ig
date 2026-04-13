package processor

type TextStyle struct {
	Bold      bool
	Italic    bool
	Underline bool
	FontColor string
	Size      *int
}

type StyledSegment struct {
	Text  string
	Style *TextStyle
}

type StyledLine struct {
	Segments []StyledSegment
}

func SingleSegment(text string, style *TextStyle) StyledLine {
	if style == nil {
		return StyledLine{
			Segments: []StyledSegment{{Text: text, Style: nil}},
		}
	}
	return StyledLine{
		Segments: []StyledSegment{{Text: text, Style: style}},
	}
}

func PlainText(text string) StyledLine {
	return StyledLine{
		Segments: []StyledSegment{{Text: text, Style: nil}},
	}
}

func HasStyledSegments(line StyledLine) bool {
	for _, seg := range line.Segments {
		if seg.Style != nil {
			if seg.Style.FontColor != "" || seg.Style.Bold || seg.Style.Italic || seg.Style.Underline || seg.Style.Size != nil {
				return true
			}
		}
	}
	return false
}
