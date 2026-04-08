package processor

import (
	"strings"
	"testing"
)

func TestBashCodeHighlighter_ProcessLines(t *testing.T) {
	tests := []struct {
		name              string
		input             []StyledLine
		expectedTexts     []string
		expectHighlighted []bool
	}{
		{
			name: "code block detection",
			input: []StyledLine{
				PlainText("Some text"),
				PlainText("```bash"),
				PlainText("echo hello"),
				PlainText("```"),
				PlainText("More text"),
			},
			expectedTexts:     []string{"Some text", "```bash", "echo hello", "```", "More text"},
			expectHighlighted: []bool{false, false, true, false, false},
		},
		{
			name: "lines outside code blocks unchanged",
			input: []StyledLine{
				PlainText("#normal comment"),
				PlainText("regular text"),
			},
			expectedTexts:     []string{"#normal comment", "regular text"},
			expectHighlighted: []bool{false, false},
		},
		{
			name: "multiple code blocks",
			input: []StyledLine{
				PlainText("```bash"),
				PlainText("ls -la"),
				PlainText("```"),
				PlainText("Some text"),
				PlainText("```bash"),
				PlainText("pwd"),
				PlainText("```"),
			},
			expectedTexts:     []string{"```bash", "ls -la", "```", "Some text", "```bash", "pwd", "```"},
			expectHighlighted: []bool{false, true, false, false, false, true, false},
		},
		{
			name: "unclosed code block",
			input: []StyledLine{
				PlainText("```bash"),
				PlainText("echo hello"),
				PlainText("more code"),
			},
			expectedTexts:     []string{"```bash", "echo hello", "more code"},
			expectHighlighted: []bool{false, false, false},
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

			if len(result) != len(tt.expectedTexts) {
				t.Fatalf("ProcessLines() returned %d lines, expected %d", len(result), len(tt.expectedTexts))
			}

			for i := range result {
				text := StyledSegmentsToText(result[i].Segments)
				if text != tt.expectedTexts[i] {
					t.Errorf("Line %d: text = %q, expected %q", i, text, tt.expectedTexts[i])
				}

				hasHighlight := HasStyledSegments(result[i])
				if hasHighlight != tt.expectHighlighted[i] {
					t.Errorf("Line %d: hasHighlight = %v, expected %v", i, hasHighlight, tt.expectHighlighted[i])
				}
			}
		})
	}
}

func TestBashCodeHighlighter_SegmentsPerToken(t *testing.T) {
	p := &BashCodeHighlighter{
		StyleName:    "monokai",
		DefaultColor: "#FFFFFF",
	}

	input := []StyledLine{
		PlainText("```bash"),
		PlainText("echo hello"),
		PlainText("```"),
	}

	result, err := p.ProcessLines(input)
	if err != nil {
		t.Fatalf("ProcessLines() error = %v", err)
	}

	echoLine := result[1]
	if len(echoLine.Segments) == 0 {
		t.Fatalf("Expected segments in echo line, got 0")
	}

	fullText := StyledSegmentsToText(echoLine.Segments)
	if fullText != "echo hello" {
		t.Errorf("Full text = %q, expected %q", fullText, "echo hello")
	}

	nonWhitespaceCount := 0
	for _, seg := range echoLine.Segments {
		if strings.TrimSpace(seg.Text) != "" && seg.Style != nil && seg.Style.FontColor != "" {
			nonWhitespaceCount++
		}
	}

	if nonWhitespaceCount < 1 {
		t.Errorf("Expected at least 1 non-whitespace segment with color, got %d", nonWhitespaceCount)
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
