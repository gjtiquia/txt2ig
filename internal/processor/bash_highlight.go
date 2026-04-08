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
		if strings.TrimSpace(lines[i].Text) == "```bash" {
			blockStart := i
			blockEnd := -1
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j].Text) == "```" {
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
				codeLines = append(codeLines, lines[k].Text)
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
			result[i] = StyledLine{Text: line, Style: nil}
			continue
		}

		tokens := iterator.Tokens()
		if len(tokens) > 0 {
			var dominantToken chroma.Token
			maxLen := 0
			for _, token := range tokens {
				if len(token.Value) > maxLen && token.Type != chroma.TokenType(chroma.Whitespace) {
					maxLen = len(token.Value)
					dominantToken = token
				}
			}

			if maxLen > 0 {
				entry := style.Get(dominantToken.Type)
				result[i] = StyledLine{
					Text:  line,
					Style: p.chromaEntryToTextStyle(entry),
				}
			} else {
				result[i] = StyledLine{Text: line, Style: nil}
			}
		} else {
			result[i] = StyledLine{Text: line, Style: nil}
		}
	}

	return result
}

func (p *BashCodeHighlighter) chromaEntryToTextStyle(entry chroma.StyleEntry) *TextStyle {
	ts := &TextStyle{}

	if entry.Colour.IsSet() {
		ts.FontColor = "#" + entry.Colour.String()
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
