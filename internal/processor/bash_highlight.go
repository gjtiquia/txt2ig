package processor

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type BashCodeHighlighter struct {
	StyleName    string
	DefaultColor string
}

func (p *BashCodeHighlighter) Process(line string) (string, *TextStyle, error) {
	return line, nil, nil
}

func (p *BashCodeHighlighter) Name() string {
	return "bash-code-highlighting"
}

func (p *BashCodeHighlighter) ProcessLines(lines []StyledLine) ([]StyledLine, error) {
	result := make([]StyledLine, len(lines))
	copy(result, lines)

	i := 0
	for i < len(lines) {
		lineText := StyledSegmentsToText(lines[i].Segments)
		if strings.TrimSpace(lineText) == "```bash" {
			blockStart := i
			blockEnd := -1
			for j := i + 1; j < len(lines); j++ {
				jText := StyledSegmentsToText(lines[j].Segments)
				if strings.TrimSpace(jText) == "```" {
					blockEnd = j
					break
				}
			}

			if blockEnd == -1 {
				i++
				continue
			}

			codeLines := make([]string, 0)
			for k := blockStart + 1; k < blockEnd; k++ {
				codeLines = append(codeLines, StyledSegmentsToText(lines[k].Segments))
			}

			highlighted := p.highlightBashCode(codeLines)

			for k := 0; k < len(highlighted); k++ {
				result[blockStart+1+k] = highlighted[k]
			}

			i = blockEnd + 1
		} else {
			i++
		}
	}

	return result, nil
}

func (p *BashCodeHighlighter) highlightBashCode(lines []string) []StyledLine {
	lexer := lexers.Get("bash")
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(p.StyleName)
	if style == nil {
		style = styles.Fallback
	}

	result := make([]StyledLine, len(lines))

	for i, line := range lines {
		iterator, err := lexer.Tokenise(nil, line)
		if err != nil {
			result[i] = PlainText(line)
			continue
		}

		tokens := iterator.Tokens()

		segments := make([]StyledSegment, 0, len(tokens))
		for _, token := range tokens {
			entry := style.Get(token.Type)

			var segStyle *TextStyle
			if token.Type != chroma.TokenType(chroma.Whitespace) || p.DefaultColor != "" {
				segStyle = p.chromaEntryToTextStyle(entry)
			}

			segments = append(segments, StyledSegment{
				Text:  token.Value,
				Style: segStyle,
			})
		}

		result[i] = StyledLine{Segments: segments}
	}

	return result
}

func (p *BashCodeHighlighter) chromaEntryToTextStyle(entry chroma.StyleEntry) *TextStyle {
	ts := &TextStyle{}

	// Chroma's Colour.String() already includes the # prefix
	// e.g., returns "#f8f8f2" not "f8f8f2"
	if entry.Colour.IsSet() {
		ts.FontColor = entry.Colour.String()
	} else {
		ts.FontColor = p.DefaultColor
	}

	if entry.Bold == chroma.Yes {
		ts.Bold = true
	}
	if entry.Italic == chroma.Yes {
		ts.Italic = true
	}
	if entry.Underline == chroma.Yes {
		ts.Underline = true
	}

	return ts
}
