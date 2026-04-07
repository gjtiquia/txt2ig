package processor

import (
	"strings"
)

type MarkdownBoldHeaders struct {
	Bold      bool   `json:"bold"`
	FontColor string `json:"fontColor"`
}

func (p *MarkdownBoldHeaders) Process(line string) (string, *TextStyle, error) {
	if strings.HasPrefix(line, "#") {
		// Remove leading # characters
		cleanedLine := strings.TrimLeft(line, "# ")
		style := &TextStyle{
			Bold:      p.Bold,
			FontColor: p.FontColor,
		}
		return cleanedLine, style, nil
	}
	return line, nil, nil
}

func (p *MarkdownBoldHeaders) Name() string {
	return "markdown-bold-headers"
}

type BashComments struct {
	Italic    bool   `json:"italic"`
	FontColor string `json:"fontColor"`
}

func (p *BashComments) Process(line string) (string, *TextStyle, error) {
	if strings.HasPrefix(line, "#") {
		style := &TextStyle{
			Italic:    p.Italic,
			FontColor: p.FontColor,
		}
		return line, style, nil
	}
	return line, nil, nil
}

func (p *BashComments) Name() string {
	return "bash-comments"
}
