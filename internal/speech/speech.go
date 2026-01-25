package speech

import (
	"strings"
)

// BubbleStyle defines a speech bubble style
type BubbleStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	Left        string
	Right       string
	Tail        string
}

var BubbleStyles = map[string]BubbleStyle{
	"round": {
		TopLeft: "╭", TopRight: "╮",
		BottomLeft: "╰", BottomRight: "╯",
		Horizontal: "─", Vertical: "│",
		Left: "│", Right: "│",
		Tail: "╲",
	},
	"square": {
		TopLeft: "┌", TopRight: "┐",
		BottomLeft: "└", BottomRight: "┘",
		Horizontal: "─", Vertical: "│",
		Left: "│", Right: "│",
		Tail: "\\",
	},
	"double": {
		TopLeft: "╔", TopRight: "╗",
		BottomLeft: "╚", BottomRight: "╝",
		Horizontal: "═", Vertical: "║",
		Left: "║", Right: "║",
		Tail: "╲",
	},
	"thick": {
		TopLeft: "┏", TopRight: "┓",
		BottomLeft: "┗", BottomRight: "┛",
		Horizontal: "━", Vertical: "┃",
		Left: "┃", Right: "┃",
		Tail: "╲",
	},
	"ascii": {
		TopLeft: "+", TopRight: "+",
		BottomLeft: "+", BottomRight: "+",
		Horizontal: "-", Vertical: "|",
		Left: "|", Right: "|",
		Tail: "\\",
	},
	"think": {
		TopLeft: "╭", TopRight: "╮",
		BottomLeft: "╰", BottomRight: "╯",
		Horizontal: "─", Vertical: "│",
		Left: "(", Right: ")",
		Tail: "○",
	},
}

// Wrap wraps text in a speech bubble
func Wrap(text string, style string, maxWidth int) string {
	bs, ok := BubbleStyles[style]
	if !ok {
		bs = BubbleStyles["round"]
	}

	if maxWidth <= 0 {
		maxWidth = 40
	}

	// Word wrap the text
	lines := wordWrap(text, maxWidth)

	// Find max line length
	maxLen := 0
	for _, line := range lines {
		if len([]rune(line)) > maxLen {
			maxLen = len([]rune(line))
		}
	}

	var sb strings.Builder

	// Top border
	sb.WriteString(bs.TopLeft)
	sb.WriteString(strings.Repeat(bs.Horizontal, maxLen+2))
	sb.WriteString(bs.TopRight)
	sb.WriteString("\n")

	// Content lines
	for _, line := range lines {
		runes := []rune(line)
		padding := maxLen - len(runes)
		sb.WriteString(bs.Left)
		sb.WriteString(" ")
		sb.WriteString(line)
		sb.WriteString(strings.Repeat(" ", padding))
		sb.WriteString(" ")
		sb.WriteString(bs.Right)
		sb.WriteString("\n")
	}

	// Bottom border
	sb.WriteString(bs.BottomLeft)
	sb.WriteString(strings.Repeat(bs.Horizontal, maxLen+2))
	sb.WriteString(bs.BottomRight)
	sb.WriteString("\n")

	// Tail
	sb.WriteString("        " + bs.Tail + "\n")
	sb.WriteString("         " + bs.Tail + "\n")

	return sb.String()
}

// WrapThink creates a thought bubble (with ○ tail)
func WrapThink(text string, maxWidth int) string {
	return Wrap(text, "think", maxWidth)
}

// wordWrap breaks text into lines of maxWidth
func wordWrap(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)

	if len(words) == 0 {
		return []string{""}
	}

	currentLine := words[0]

	for _, word := range words[1:] {
		if len([]rune(currentLine))+1+len([]rune(word)) <= maxWidth {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}

// Combine puts speech bubble above ASCII art
func Combine(bubble string, art string) string {
	return bubble + art
}

// ListStyles returns available bubble styles
func ListStyles() []string {
	return []string{"round", "square", "double", "thick", "ascii", "think"}
}
