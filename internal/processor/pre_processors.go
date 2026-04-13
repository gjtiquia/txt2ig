package processor

import (
	"fmt"
	"strings"
	"time"
)

type ExactSearchAndReplace struct {
	SearchString  string `json:"searchString"`
	ReplaceString string `json:"replaceString"`
}

func (p *ExactSearchAndReplace) Process(text string) (string, error) {
	return strings.ReplaceAll(text, p.SearchString, p.ReplaceString), nil
}

func (p *ExactSearchAndReplace) Name() string {
	return "exactSearchAndReplace"
}

type ExactSearchAndReplaceWithDateTimeNow struct {
	SearchString  string `json:"searchString"`
	ReplaceFormat string `json:"replaceFormat"`
}

func (p *ExactSearchAndReplaceWithDateTimeNow) Process(text string) (string, error) {
	now := time.Now()
	formatted := p.formatDateTime(now)
	return strings.ReplaceAll(text, p.SearchString, formatted), nil
}

func (p *ExactSearchAndReplaceWithDateTimeNow) Name() string {
	return "exactSearchAndReplaceWithDateTimeNow"
}

func (p *ExactSearchAndReplaceWithDateTimeNow) formatDateTime(t time.Time) string {
	result := p.ReplaceFormat
	result = strings.ReplaceAll(result, "yyyy", formatInt(t.Year(), 4))
	result = strings.ReplaceAll(result, "mm", formatInt(int(t.Month()), 2))
	result = strings.ReplaceAll(result, "dd", formatInt(t.Day(), 2))
	result = strings.ReplaceAll(result, "hh", formatInt(t.Hour(), 2))
	result = strings.ReplaceAll(result, "MM", formatInt(t.Minute(), 2))
	result = strings.ReplaceAll(result, "ss", formatInt(t.Second(), 2))
	return result
}

func formatInt(n int, width int) string {
	s := fmt.Sprintf("%d", n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

func parseExactSearchAndReplace(config interface{}) (PreProcessor, error) {
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

func parseExactSearchAndReplaceWithDateTimeNow(config interface{}) (PreProcessor, error) {
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

func init() {
	RegisterPreProcessorParser("exactSearchAndReplace", parseExactSearchAndReplace)
	RegisterPreProcessorParser("exactSearchAndReplaceWithDateTimeNow", parseExactSearchAndReplaceWithDateTimeNow)
}
