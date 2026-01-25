package speech

import (
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		style    string
		width    int
		contains []string
	}{
		{
			name:     "round bubble",
			text:     "Hello World",
			style:    "round",
			width:    40,
			contains: []string{"╭", "╮", "╰", "╯", "│", "Hello World"},
		},
		{
			name:     "square bubble",
			text:     "Test",
			style:    "square",
			width:    40,
			contains: []string{"┌", "┐", "└", "┘", "│", "Test"},
		},
		{
			name:     "double bubble",
			text:     "Double",
			style:    "double",
			width:    40,
			contains: []string{"╔", "╗", "╚", "╝", "║", "Double"},
		},
		{
			name:     "thick bubble",
			text:     "Thick",
			style:    "thick",
			width:    40,
			contains: []string{"┏", "┓", "┗", "┛", "┃", "Thick"},
		},
		{
			name:     "ascii bubble",
			text:     "ASCII",
			style:    "ascii",
			width:    40,
			contains: []string{"+", "-", "|", "ASCII"},
		},
		{
			name:     "think bubble",
			text:     "Thinking",
			style:    "think",
			width:    40,
			contains: []string{"╭", "╮", "(", ")", "○", "Thinking"},
		},
		{
			name:     "unknown defaults to round",
			text:     "Default",
			style:    "nonexistent",
			width:    40,
			contains: []string{"╭", "╮", "│", "Default"},
		},
		{
			name:     "zero width defaults to 40",
			text:     "Zero",
			style:    "round",
			width:    0,
			contains: []string{"╭", "Zero"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Wrap(tt.text, tt.style, tt.width)
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Wrap(%q, %q, %d) should contain %q\nGot:\n%s",
						tt.text, tt.style, tt.width, substr, result)
				}
			}
		})
	}
}

func TestWrapWordWrap(t *testing.T) {
	longText := "This is a very long sentence that should be wrapped into multiple lines inside the bubble"
	result := Wrap(longText, "round", 20)

	lines := strings.Split(result, "\n")
	// Should have multiple content lines (more than just top + 1 content + bottom + tail)
	if len(lines) < 5 {
		t.Errorf("Long text should be wrapped into multiple lines, got %d lines", len(lines))
	}
}

func TestWrapThink(t *testing.T) {
	result := WrapThink("Deep thoughts", 40)
	if !strings.Contains(result, "○") {
		t.Error("WrapThink should use ○ tail character")
	}
	if !strings.Contains(result, "(") || !strings.Contains(result, ")") {
		t.Error("WrapThink should use ( and ) for sides")
	}
}

func TestCombine(t *testing.T) {
	bubble := "╭──╮\n│Hi│\n╰──╯\n"
	art := " /\\_/\\\n( o.o )\n"
	result := Combine(bubble, art)

	if !strings.Contains(result, "╭──╮") {
		t.Error("Combine should contain bubble")
	}
	if !strings.Contains(result, "( o.o )") {
		t.Error("Combine should contain art")
	}
	// Bubble should come before art
	bubbleIdx := strings.Index(result, "╭──╮")
	artIdx := strings.Index(result, "( o.o )")
	if bubbleIdx > artIdx {
		t.Error("Bubble should appear before art in output")
	}
}

func TestListStyles(t *testing.T) {
	styles := ListStyles()
	if len(styles) != 6 {
		t.Errorf("ListStyles should return 6 styles, got %d", len(styles))
	}

	expected := map[string]bool{
		"round": true, "square": true, "double": true,
		"thick": true, "ascii": true, "think": true,
	}
	for _, s := range styles {
		if !expected[s] {
			t.Errorf("Unexpected style: %s", s)
		}
	}
}

func TestWordWrap(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected int // expected number of lines
	}{
		{"short text", "Hello", 40, 1},
		{"empty text", "", 40, 1},
		{"exact width", "Hello World", 11, 1},
		{"needs wrap", "Hello World", 5, 2},
		{"zero width", "Hello", 0, 1},
		{"multiple words", "one two three four five", 10, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wordWrap(tt.text, tt.width)
			if len(result) != tt.expected {
				t.Errorf("wordWrap(%q, %d) = %d lines, want %d. Lines: %v",
					tt.text, tt.width, len(result), tt.expected, result)
			}
		})
	}
}

func TestBubbleTail(t *testing.T) {
	result := Wrap("Hi", "round", 40)
	// Should have tail lines at the end
	lines := strings.Split(result, "\n")
	// Find the tail character
	foundTail := false
	for _, line := range lines {
		if strings.Contains(line, "╲") {
			foundTail = true
			break
		}
	}
	if !foundTail {
		t.Error("Round bubble should have ╲ tail")
	}
}

func TestBubbleContentPadding(t *testing.T) {
	result := Wrap("Hi", "round", 40)
	// Content line should have padding: "│ Hi │"
	if !strings.Contains(result, "│ Hi") {
		t.Errorf("Content should be padded with spaces, got:\n%s", result)
	}
}
