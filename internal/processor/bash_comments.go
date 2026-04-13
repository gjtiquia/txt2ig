package processor

import (
	"fmt"
	"strings"
)

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

func parseBashComments(config interface{}) (PostProcessor, error) {
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

func init() {
	RegisterPostProcessorParser("bash-comments", parseBashComments)
}
