package processor

import (
	"fmt"
	"regexp"
)

var linkRegex = regexp.MustCompile(`\[([^\]]*)\]\(([^\)]*)\)`)

type MarkdownLinks struct {
	NameColor string `json:"nameColor"`
	LinkColor string `json:"linkColor"`
}

func (p *MarkdownLinks) Name() string {
	return "markdown-links"
}

func (p *MarkdownLinks) ProcessLines(lines []StyledLine) ([]StyledLine, error) {
	result := make([]StyledLine, len(lines))

	for i, line := range lines {
		text := StyledSegmentsToText(line.Segments)

		matches := linkRegex.FindAllStringSubmatchIndex(text, -1)

		if len(matches) == 0 {
			result[i] = line
			continue
		}

		segments := []StyledSegment{}
		lastEnd := 0

		for _, match := range matches {
			fullStart := match[0]
			fullEnd := match[1]
			nameStart := match[2]
			nameEnd := match[3]
			urlStart := match[4]
			urlEnd := match[5]

			name := text[nameStart:nameEnd]
			url := text[urlStart:urlEnd]

			if fullStart > lastEnd {
				segments = append(segments, StyledSegment{
					Text:  text[lastEnd:fullStart],
					Style: getStyleAtPosition(line, lastEnd),
				})
			}

			segments = append(segments, StyledSegment{
				Text:  "[",
				Style: nil,
			})

			if name != "" {
				nameStyle := getStyleAtPosition(line, nameStart)
				if p.NameColor != "" {
					if nameStyle == nil {
						nameStyle = &TextStyle{FontColor: p.NameColor}
					} else {
						merged := *nameStyle
						merged.FontColor = p.NameColor
						nameStyle = &merged
					}
				}
				segments = append(segments, StyledSegment{
					Text:  name,
					Style: nameStyle,
				})
			}

			segments = append(segments, StyledSegment{
				Text:  "](",
				Style: nil,
			})

			if url != "" {
				urlStyle := getStyleAtPosition(line, urlStart)
				if p.LinkColor != "" {
					if urlStyle == nil {
						urlStyle = &TextStyle{FontColor: p.LinkColor}
					} else {
						merged := *urlStyle
						merged.FontColor = p.LinkColor
						urlStyle = &merged
					}
				}
				segments = append(segments, StyledSegment{
					Text:  url,
					Style: urlStyle,
				})
			}

			segments = append(segments, StyledSegment{
				Text:  ")",
				Style: nil,
			})

			lastEnd = fullEnd
		}

		if lastEnd < len(text) {
			segments = append(segments, StyledSegment{
				Text:  text[lastEnd:],
				Style: getStyleAtPosition(line, lastEnd),
			})
		}

		result[i] = StyledLine{Segments: segments}
	}

	return result, nil
}

func getStyleAtPosition(line StyledLine, pos int) *TextStyle {
	currentPos := 0
	for _, seg := range line.Segments {
		segStart := currentPos
		segEnd := currentPos + len(seg.Text)
		if pos >= segStart && pos < segEnd {
			return seg.Style
		}
		currentPos = segEnd
	}
	return nil
}

func parseMarkdownLinks(config interface{}) (PostProcessor, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("markdown-links config must be an object")
	}

	p := &MarkdownLinks{}
	if nameColor, ok := configMap["nameColor"]; ok {
		p.NameColor, ok = nameColor.(string)
		if !ok {
			return nil, fmt.Errorf("nameColor must be a string")
		}
	}
	if linkColor, ok := configMap["linkColor"]; ok {
		p.LinkColor, ok = linkColor.(string)
		if !ok {
			return nil, fmt.Errorf("linkColor must be a string")
		}
	}

	return p, nil
}

func init() {
	RegisterPostProcessorParser("markdown-links", parseMarkdownLinks)
}
