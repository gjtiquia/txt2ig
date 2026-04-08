package processor

import (
	"fmt"
)

type PreProcessor interface {
	Process(text string) (string, error)
	Name() string
}

type PostProcessor interface {
	Process(line string) (string, *TextStyle, error)
	Name() string
}

type TextStyle struct {
	Bold      bool
	Italic    bool
	Underline bool
	FontColor string
	Size      *int
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

func ApplyPostProcessors(lines []string, configs []interface{}) ([]string, error) {
	processors, err := ParsePostProcessorConfigs(configs)
	if err != nil {
		return nil, err
	}

	result := lines
	for i, p := range processors {
		newLines := make([]string, 0, len(result))
		for _, line := range result {
			processedLine, _, err := p.Process(line)
			if err != nil {
				return nil, fmt.Errorf("apply post-processor %d (%s): %w", i, p.Name(), err)
			}
			newLines = append(newLines, processedLine)
		}
		result = newLines
	}

	return result, nil
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

	// Try each known processor type
	if mdBold, ok := configMap["markdown-bold-headers"]; ok {
		return parseMarkdownBoldHeaders(mdBold)
	}

	if bashComments, ok := configMap["bash-comments"]; ok {
		return parseBashComments(bashComments)
	}

	// Unknown processor type
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
