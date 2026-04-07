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
