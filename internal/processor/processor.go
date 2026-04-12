package processor

import (
	"fmt"
	"strings"
)

type PreProcessor interface {
	Process(text string) (string, error)
	Name() string
}

type PostProcessor interface {
	Process(line string) (string, *TextStyle, error)
	Name() string
}

type StatefulPostProcessor interface {
	PostProcessor
	ProcessLines(lines []StyledLine) ([]StyledLine, error)
}

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

type ProcessorRegistry struct {
	preProcessors  map[string]PreProcessor
	postProcessors map[string]PostProcessor
}

func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{
		preProcessors:  make(map[string]PreProcessor),
		postProcessors: make(map[string]PostProcessor),
	}
}

func (r *ProcessorRegistry) RegisterPreProcessor(name string, processor PreProcessor) {
	r.preProcessors[name] = processor
}

func (r *ProcessorRegistry) RegisterPostProcessor(name string, processor PostProcessor) {
	r.postProcessors[name] = processor
}

func (r *ProcessorRegistry) GetPreProcessor(name string) (PreProcessor, bool) {
	p, ok := r.preProcessors[name]
	return p, ok
}

func (r *ProcessorRegistry) GetPostProcessor(name string) (PostProcessor, bool) {
	p, ok := r.postProcessors[name]
	return p, ok
}

func ApplyPreProcessors(text string, configs []interface{}) (string, error) {
	processors, err := ParsePreProcessorConfigs(configs)
	if err != nil {
		return "", err
	}

	result := text
	for i, p := range processors {
		result, err = p.Process(result)
		if err != nil {
			return "", fmt.Errorf("apply pre-processor %d (%s): %w", i, p.Name(), err)
		}
	}

	return result, nil
}

func ApplyPostProcessors(lines []string, configs []interface{}) ([]StyledLine, error) {
	processors, err := ParsePostProcessorConfigs(configs)
	if err != nil {
		return nil, err
	}

	result := make([]StyledLine, len(lines))
	for i, line := range lines {
		result[i] = PlainText(line)
	}

	for _, p := range processors {
		if sp, ok := p.(StatefulPostProcessor); ok {
			result, err = sp.ProcessLines(result)
			if err != nil {
				return nil, err
			}
		}
	}

	for i := range result {
		if HasStyledSegments(result[i]) {
			continue
		}

		var currentStyle *TextStyle
		lineText := StyledSegmentsToText(result[i].Segments)

		for _, p := range processors {
			if _, ok := p.(StatefulPostProcessor); ok {
				continue
			}

			var style *TextStyle
			lineText, style, err = p.Process(lineText)
			if err != nil {
				return nil, fmt.Errorf("apply post-processor: %w", err)
			}
			if style != nil {
				currentStyle = mergeStyles(currentStyle, style)
			}
		}

		if currentStyle != nil {
			for j := range result[i].Segments {
				result[i].Segments[j].Style = mergeStyles(result[i].Segments[j].Style, currentStyle)
			}
		}
	}

	return result, nil
}

func hasNonDefaultStyle(style *TextStyle) bool {
	if style == nil {
		return false
	}
	return style.FontColor != "" || style.Bold || style.Italic || style.Underline || style.Size != nil
}

func mergeStyles(base, override *TextStyle) *TextStyle {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	// Create a merged style, with override taking precedence
	merged := &TextStyle{
		Bold:      base.Bold || override.Bold,
		Italic:    base.Italic || override.Italic,
		Underline: base.Underline || override.Underline,
	}

	// FontColor: override takes precedence if set
	if override.FontColor != "" {
		merged.FontColor = override.FontColor
	} else {
		merged.FontColor = base.FontColor
	}

	// Size: override takes precedence if set
	if override.Size != nil {
		merged.Size = override.Size
	} else {
		merged.Size = base.Size
	}

	return merged
}

func ParsePreProcessorConfig(config interface{}) (PreProcessor, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("processor config must be an object")
	}

	// Try each known processor type
	if exactReplace, ok := configMap["exactSearchAndReplace"]; ok {
		return parseExactSearchAndReplace(exactReplace)
	}

	if dateTimeReplace, ok := configMap["exactSearchAndReplaceWithDateTimeNow"]; ok {
		return parseExactSearchAndReplaceWithDateTimeNow(dateTimeReplace)
	}

	// Unknown processor type
	// Get the first key to report in error message
	for k := range configMap {
		return nil, fmt.Errorf("unknown pre-processor: %s", k)
	}

	return nil, fmt.Errorf("empty processor config")
}

func parseExactSearchAndReplace(config interface{}) (*ExactSearchAndReplace, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("exactSearchAndReplace config must be an object")
	}

	p := &ExactSearchAndReplace{}
	if search, ok := configMap["searchString"]; ok {
		p.SearchString, ok = search.(string)
		if !ok {
			return nil, fmt.Errorf("searchString must be a string")
		}
	}
	if replace, ok := configMap["replaceString"]; ok {
		p.ReplaceString, ok = replace.(string)
		if !ok {
			return nil, fmt.Errorf("replaceString must be a string")
		}
	}

	return p, nil
}

func parseExactSearchAndReplaceWithDateTimeNow(config interface{}) (*ExactSearchAndReplaceWithDateTimeNow, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("exactSearchAndReplaceWithDateTimeNow config must be an object")
	}

	p := &ExactSearchAndReplaceWithDateTimeNow{}
	if search, ok := configMap["searchString"]; ok {
		p.SearchString, ok = search.(string)
		if !ok {
			return nil, fmt.Errorf("searchString must be a string")
		}
	}
	if format, ok := configMap["replaceFormat"]; ok {
		p.ReplaceFormat, ok = format.(string)
		if !ok {
			return nil, fmt.Errorf("replaceFormat must be a string")
		}
	}

	return p, nil
}

func ParsePostProcessorConfig(config interface{}) (PostProcessor, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("processor config must be an object")
	}

	if mdBold, ok := configMap["markdown-bold-headers"]; ok {
		return parseMarkdownBoldHeaders(mdBold)
	}

	if mdLinks, ok := configMap["markdown-links"]; ok {
		return parseMarkdownLinks(mdLinks)
	}

	if bashComments, ok := configMap["bash-comments"]; ok {
		return parseBashComments(bashComments)
	}

	if bashHighlight, ok := configMap["bash-code-highlighting"]; ok {
		return parseBashCodeHighlighting(bashHighlight)
	}

	for k := range configMap {
		return nil, fmt.Errorf("unknown post-processor: %s", k)
	}

	return nil, fmt.Errorf("empty processor config")
}

func parseMarkdownBoldHeaders(config interface{}) (*MarkdownBoldHeaders, error) {
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

func parseMarkdownLinks(config interface{}) (*MarkdownLinks, error) {
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

func parseBashComments(config interface{}) (*BashComments, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("bash-comments config must be an object")
	}

	p := &BashComments{}
	if italic, ok := configMap["italic"]; ok {
		p.Italic, ok = italic.(bool)
		if !ok {
			return nil, fmt.Errorf("italic must be a boolean")
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

func parseBashCodeHighlighting(config interface{}) (*BashCodeHighlighter, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("bash-code-highlighting config must be an object")
	}

	p := &BashCodeHighlighter{
		StyleName:    "monokai",
		DefaultColor: "#FFFFFF",
	}

	if style, ok := configMap["style"]; ok {
		p.StyleName, ok = style.(string)
		if !ok {
			return nil, fmt.Errorf("style must be a string")
		}
	}

	if defaultColor, ok := configMap["defaultColor"]; ok {
		p.DefaultColor, ok = defaultColor.(string)
		if !ok {
			return nil, fmt.Errorf("defaultColor must be a string")
		}
	}

	return p, nil
}

func ParsePreProcessorConfigs(configs []interface{}) ([]PreProcessor, error) {
	processors := make([]PreProcessor, 0, len(configs))
	for i, config := range configs {
		p, err := ParsePreProcessorConfig(config)
		if err != nil {
			return nil, fmt.Errorf("parse pre-processor %d: %w", i, err)
		}
		processors = append(processors, p)
	}
	return processors, nil
}

func ParsePostProcessorConfigs(configs []interface{}) ([]PostProcessor, error) {
	processors := make([]PostProcessor, 0, len(configs))
	for i, config := range configs {
		p, err := ParsePostProcessorConfig(config)
		if err != nil {
			return nil, fmt.Errorf("parse post-processor %d: %w", i, err)
		}
		processors = append(processors, p)
	}
	return processors, nil
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

func StyledSegmentsToText(segments []StyledSegment) string {
	var sb strings.Builder
	for _, seg := range segments {
		sb.WriteString(seg.Text)
	}
	return sb.String()
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
