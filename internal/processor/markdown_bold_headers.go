package processor

import (
	"fmt"
	"strings"
)

type MarkdownBoldHeaders struct {
	Bold      bool   `json:"bold"`
	FontColor string `json:"fontColor"`
}

func (p *MarkdownBoldHeaders) ProcessLines(lines []StyledLine) ([]StyledLine, error) {
	result := make([]StyledLine, len(lines))
	for i, line := range lines {
		text := StyledSegmentsToText(line.Segments)
		if strings.HasPrefix(text, "#") {
			style := &TextStyle{
				Bold:      p.Bold,
				FontColor: p.FontColor,
			}
			for j := range line.Segments {
				line.Segments[j].Style = mergeStyles(line.Segments[j].Style, style)
			}
		}
		result[i] = line
	}
	return result, nil
}

func (p *MarkdownBoldHeaders) Name() string {
	return "markdown-bold-headers"
}

func parseMarkdownBoldHeaders(config interface{}) (PostProcessor, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("markdown-bold-headers config must be an object")
	}

	p := &MarkdownBoldHeaders{}
	if bold, ok := configMap["bold"]; ok {
		p.Bold, ok = bold.(bool)
		if !ok {
			return nil, fmt.Errorf("bold must be a boolean")
		}
	}
	if fontColor, ok := configMap["fontColor"]; ok {
		p.FontColor, ok = fontColor.(string)
		if !ok {
			return nil, fmt.Errorf("fontColor must be a string")
		}
	}

	return p, nil
}

func init() {
	RegisterPostProcessorParser("markdown-bold-headers", parseMarkdownBoldHeaders)
}
