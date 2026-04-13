package processor

import (
	"testing"
)

func TestApplyPostProcessors_RespectsConfigOrder_MixedProcessors(t *testing.T) {
	configs := []interface{}{
		map[string]interface{}{
			"markdown-links": map[string]interface{}{
				"nameColor": "#3B82F6",
				"linkColor": "#60A5FA",
			},
		},
		map[string]interface{}{
			"markdown-bold-headers": map[string]interface{}{
				"bold":      true,
				"fontColor": "#FF0000",
			},
		},
	}

	lines := []string{"# Check out [my link](https://example.com)"}

	result, err := ApplyPostProcessors(lines, configs)
	if err != nil {
		t.Fatalf("ApplyPostProcessors() error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(result))
	}

	segs := result[0].Segments

	var nameSeg *StyledSegment
	for i := range segs {
		if segs[i].Text == "my link" {
			nameSeg = &segs[i]
			break
		}
	}

	if nameSeg == nil {
		t.Fatal("Could not find 'my link' segment")
	}

	if nameSeg.Style == nil {
		t.Fatal("Name segment style is nil")
	}

	if nameSeg.Style.FontColor != "#FF0000" {
		t.Errorf("Name segment FontColor = %q, want #FF0000 (bold-headers should override since it runs AFTER markdown-links per config order)", nameSeg.Style.FontColor)
	}

	if !nameSeg.Style.Bold {
		t.Error("Name segment should be Bold=true")
	}
}

func TestApplyPostProcessors_RespectsConfigOrder_SimpleAfterStateful(t *testing.T) {
	configs := []interface{}{
		map[string]interface{}{
			"markdown-links": map[string]interface{}{
				"nameColor": "#3B82F6",
				"linkColor": "#60A5FA",
			},
		},
		map[string]interface{}{
			"markdown-bold-headers": map[string]interface{}{
				"bold": true,
			},
		},
	}

	lines := []string{"Click [here](url) for more"}

	result, err := ApplyPostProcessors(lines, configs)
	if err != nil {
		t.Fatalf("ApplyPostProcessors() error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(result))
	}

	segs := result[0].Segments

	var nameSeg *StyledSegment
	for i := range segs {
		if segs[i].Text == "here" {
			nameSeg = &segs[i]
			break
		}
	}

	if nameSeg == nil {
		t.Fatal("Could not find 'here' segment")
	}

	if nameSeg.Style == nil {
		t.Fatal("Name segment style is nil")
	}

	if nameSeg.Style.FontColor != "#3B82F6" {
		t.Errorf("Name segment FontColor = %q, want #3B82F6 (markdown-links sets this)", nameSeg.Style.FontColor)
	}
}

func TestApplyPostProcessors_RespectsConfigOrder_StatefulAfterSimple(t *testing.T) {
	configs := []interface{}{
		map[string]interface{}{
			"markdown-bold-headers": map[string]interface{}{
				"bold":      true,
				"fontColor": "#FF0000",
			},
		},
		map[string]interface{}{
			"markdown-links": map[string]interface{}{
				"nameColor": "#3B82F6",
				"linkColor": "#60A5FA",
			},
		},
	}

	lines := []string{"# Check out [my link](https://example.com)"}

	result, err := ApplyPostProcessors(lines, configs)
	if err != nil {
		t.Fatalf("ApplyPostProcessors() error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(result))
	}

	segs := result[0].Segments

	var nameSeg *StyledSegment
	for i := range segs {
		if segs[i].Text == "my link" {
			nameSeg = &segs[i]
			break
		}
	}

	if nameSeg == nil {
		t.Fatal("Could not find 'my link' segment")
	}

	if nameSeg.Style == nil {
		t.Fatal("Name segment style is nil")
	}

	if nameSeg.Style.FontColor != "#3B82F6" {
		t.Errorf("Name segment FontColor = %q, want #3B82F6 (markdown-links runs AFTER bold-headers and should override color)", nameSeg.Style.FontColor)
	}

	if !nameSeg.Style.Bold {
		t.Error("Name segment should be Bold=true (bold-headers set it first, markdown-links should preserve)")
	}
}
