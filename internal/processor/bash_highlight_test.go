package processor

import (
	"testing"
)

func TestBashCodeHighlighter_ProcessLines(t *testing.T) {
	tests := []struct {
		name     string
		input    []StyledLine
		expected []StyledLine
	}{
		{
			name: "code block detection",
			input: []StyledLine{
				{Text: "Some text"},
				{Text: "```bash"},
				{Text: "echo hello"},
				{Text: "```"},
				{Text: "More text"},
			},
			expected: []StyledLine{
				{Text: "Some text"},
				{Text: "```bash"},
				{Text: "echo hello", Style: &TextStyle{FontColor: "#66D9EF"}},
				{Text: "```"},
				{Text: "More text"},
			},
		},
		{
			name: "lines outside code blocks unchanged",
			input: []StyledLine{
				{Text: "#normal comment"},
				{Text: "regular text"},
			},
			expected: []StyledLine{
				{Text: "#normal comment"},
				{Text: "regular text"},
			},
		},
		{
			name: "multiple code blocks",
			input: []StyledLine{
				{Text: "```bash"},
				{Text: "ls -la"},
				{Text: "```"},
				{Text: "Some text"},
				{Text: "```bash"},
				{Text: "pwd"},
				{Text: "```"},
			},
			expected: []StyledLine{
				{Text: "```bash"},
				{Text: "ls -la", Style: &TextStyle{FontColor: "#66D9EF"}},
				{Text: "```"},
				{Text: "Some text"},
				{Text: "```bash"},
				{Text: "pwd", Style: &TextStyle{FontColor: "#66D9EF"}},
				{Text: "```"},
			},
		},
		{
			name: "unclosed code block",
			input: []StyledLine{
				{Text: "```bash"},
				{Text: "echo hello"},
				{Text: "more code"},
			},
			expected: []StyledLine{
				{Text: "```bash"},
				{Text: "echo hello"},
				{Text: "more code"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &BashCodeHighlighter{
				StyleName:    "monokai",
				DefaultColor: "#FFFFFF",
			}

			result, err := p.ProcessLines(tt.input)
			if err != nil {
				t.Fatalf("ProcessLines() error = %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("ProcessLines() returned %d lines, expected %d", len(result), len(tt.expected))
			}

			for i := range result {
				if result[i].Text != tt.expected[i].Text {
					t.Errorf("Line %d: text = %q, expected %q", i, result[i].Text, tt.expected[i].Text)
				}

				if tt.expected[i].Style != nil {
					if result[i].Style == nil {
						t.Errorf("Line %d: expected style, got nil", i)
					} else if result[i].Style.FontColor == "" {
						t.Errorf("Line %d: expected FontColor, got empty", i)
					}
				} else {
					if result[i].Style != nil && result[i].Style.FontColor != "" {
						t.Errorf("Line %d: expected no style, got FontColor = %s", i, result[i].Style.FontColor)
					}
				}
			}
		})
	}
}

func TestBashCodeHighlighter_Name(t *testing.T) {
	p := &BashCodeHighlighter{}
	if p.Name() != "bash-code-highlighting" {
		t.Errorf("Name() = %q, expected %q", p.Name(), "bash-code-highlighting")
	}
}

func TestBashCodeHighlighter_Process(t *testing.T) {
	p := &BashCodeHighlighter{}
	line := "test"
	result, style, err := p.Process(line)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if result != line {
		t.Errorf("Process() = %q, expected %q", result, line)
	}
	if style != nil {
		t.Errorf("Process() style = %v, expected nil", style)
	}
}

func TestParseBashCodeHighlighting(t *testing.T) {
	tests := []struct {
		name    string
		config  interface{}
		want    *BashCodeHighlighter
		wantErr bool
	}{
		{
			name: "default values",
			config: map[string]interface{}{
				"bash-code-highlighting": map[string]interface{}{},
			},
			want: &BashCodeHighlighter{
				StyleName:    "monokai",
				DefaultColor: "#FFFFFF",
			},
		},
		{
			name: "custom values",
			config: map[string]interface{}{
				"bash-code-highlighting": map[string]interface{}{
					"style":        "dracula",
					"defaultColor": "#F8F8F2",
				},
			},
			want: &BashCodeHighlighter{
				StyleName:    "dracula",
				DefaultColor: "#F8F8F2",
			},
		},
		{
			name: "invalid config type",
			config: map[string]interface{}{
				"bash-code-highlighting": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid style type",
			config: map[string]interface{}{
				"bash-code-highlighting": map[string]interface{}{
					"style": 123,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePostProcessorConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePostProcessorConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			p, ok := result.(*BashCodeHighlighter)
			if !ok {
				t.Errorf("ParsePostProcessorConfig() returned %T, expected *BashCodeHighlighter", result)
				return
			}

			if p.StyleName != tt.want.StyleName {
				t.Errorf("StyleName = %q, expected %q", p.StyleName, tt.want.StyleName)
			}
			if p.DefaultColor != tt.want.DefaultColor {
				t.Errorf("DefaultColor = %q, expected %q", p.DefaultColor, tt.want.DefaultColor)
			}
		})
	}
}
