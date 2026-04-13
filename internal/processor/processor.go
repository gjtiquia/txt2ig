package processor

import (
	"fmt"
	"strings"
)

type PreProcessorParserFunc func(config interface{}) (PreProcessor, error)
type PostProcessorParserFunc func(config interface{}) (PostProcessor, error)

var (
	preProcessorParsers  = make(map[string]PreProcessorParserFunc)
	postProcessorParsers = make(map[string]PostProcessorParserFunc)
)

func RegisterPreProcessorParser(name string, parser PreProcessorParserFunc) {
	preProcessorParsers[name] = parser
}

func RegisterPostProcessorParser(name string, parser PostProcessorParserFunc) {
	postProcessorParsers[name] = parser
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
		result, err = p.ProcessLines(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func mergeStyles(base, override *TextStyle) *TextStyle {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := &TextStyle{
		Bold:      base.Bold || override.Bold,
		Italic:    base.Italic || override.Italic,
		Underline: base.Underline || override.Underline,
	}

	if override.FontColor != "" {
		merged.FontColor = override.FontColor
	} else {
		merged.FontColor = base.FontColor
	}

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

	for name, parser := range preProcessorParsers {
		if cfg, ok := configMap[name]; ok {
			return parser(cfg)
		}
	}

	for k := range configMap {
		return nil, fmt.Errorf("unknown pre-processor: %s", k)
	}

	return nil, fmt.Errorf("empty processor config")
}

func ParsePostProcessorConfig(config interface{}) (PostProcessor, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("processor config must be an object")
	}

	for name, parser := range postProcessorParsers {
		if cfg, ok := configMap[name]; ok {
			return parser(cfg)
		}
	}

	for k := range configMap {
		return nil, fmt.Errorf("unknown post-processor: %s", k)
	}

	return nil, fmt.Errorf("empty processor config")
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

func StyledSegmentsToText(segments []StyledSegment) string {
	var sb strings.Builder
	for _, seg := range segments {
		sb.WriteString(seg.Text)
	}
	return sb.String()
}
