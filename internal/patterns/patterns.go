package patterns

import (
	"strings"
)

// Border styles
var Borders = map[string]struct {
	TL, T, TR string
	L, R      string
	BL, B, BR string
}{
	"single": {"┌", "─", "┐", "│", "│", "└", "─", "┘"},
	"double": {"╔", "═", "╗", "║", "║", "╚", "═", "╝"},
	"round":  {"╭", "─", "╮", "│", "│", "╰", "─", "╯"},
	"thick":  {"┏", "━", "┓", "┃", "┃", "┗", "━", "┛"},
	"ascii":  {"+", "-", "+", "|", "|", "+", "-", "+"},
	"stars":  {"*", "*", "*", "*", "*", "*", "*", "*"},
	"dots":   {"·", "·", "·", "·", "·", "·", "·", "·"},
	"blocks": {"█", "▀", "█", "█", "█", "█", "▄", "█"},
	"shadow": {"┌", "─", "┒", "│", "┃", "┕", "━", "┛"},
	"fancy":  {"╔", "╦", "╗", "╠", "╣", "╚", "╩", "╝"},
}

// Dividers - horizontal line styles
var Dividers = map[string]string{
	"single":   "─",
	"double":   "═",
	"thick":    "━",
	"dotted":   "┄",
	"dashed":   "┅",
	"wavy":     "～",
	"stars":    "✦",
	"hearts":   "♥",
	"arrows":   "→",
	"blocks":   "█",
	"gradient": "░▒▓",
	"sparkle":  "✧✦★",
	"flower":   "❀✿❁",
	"music":    "♪♫♬",
	"snow":     "❄❅❆",
	"zigzag":   "⌇",
}

// CreateBorder wraps text in a border
func CreateBorder(text string, style string, padding int) string {
	b, ok := Borders[style]
	if !ok {
		b = Borders["single"]
	}

	lines := strings.Split(text, "\n")

	// Find max width
	maxWidth := 0
	for _, line := range lines {
		if len([]rune(line)) > maxWidth {
			maxWidth = len([]rune(line))
		}
	}

	width := maxWidth + padding*2

	var sb strings.Builder

	// Top border
	sb.WriteString(b.TL)
	sb.WriteString(strings.Repeat(b.T, width))
	sb.WriteString(b.TR)
	sb.WriteString("\n")

	// Content
	pad := strings.Repeat(" ", padding)
	for _, line := range lines {
		runes := []rune(line)
		rightPad := maxWidth - len(runes)
		sb.WriteString(b.L)
		sb.WriteString(pad)
		sb.WriteString(line)
		sb.WriteString(strings.Repeat(" ", rightPad))
		sb.WriteString(pad)
		sb.WriteString(b.R)
		sb.WriteString("\n")
	}

	// Bottom border
	sb.WriteString(b.BL)
	sb.WriteString(strings.Repeat(b.B, width))
	sb.WriteString(b.BR)
	sb.WriteString("\n")

	return sb.String()
}

// CreateDivider creates a horizontal divider
func CreateDivider(style string, width int) string {
	pattern, ok := Dividers[style]
	if !ok {
		pattern = Dividers["single"]
	}

	runes := []rune(pattern)
	var sb strings.Builder

	for i := 0; i < width; i++ {
		sb.WriteRune(runes[i%len(runes)])
	}

	return sb.String()
}

// CreatePattern creates a repeating pattern block
func CreatePattern(pattern string, width, height int) string {
	runes := []rune(pattern)
	if len(runes) == 0 {
		return ""
	}

	var sb strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sb.WriteRune(runes[(x+y)%len(runes)])
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// CreateSymmetric creates horizontally symmetric pattern
func CreateSymmetric(half string) string {
	lines := strings.Split(half, "\n")
	var sb strings.Builder

	for _, line := range lines {
		runes := []rune(line)
		sb.WriteString(line)
		// Mirror
		for i := len(runes) - 1; i >= 0; i-- {
			sb.WriteRune(mirrorChar(runes[i]))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// mirrorChar returns the mirrored version of a character
func mirrorChar(r rune) rune {
	mirrors := map[rune]rune{
		'/': '\\', '\\': '/',
		'(': ')', ')': '(',
		'[': ']', ']': '[',
		'{': '}', '}': '{',
		'<': '>', '>': '<',
		'd': 'b', 'b': 'd',
		'p': 'q', 'q': 'p',
		'╱': '╲', '╲': '╱',
		'╭': '╮', '╮': '╭',
		'╰': '╯', '╯': '╰',
		'┌': '┐', '┐': '┌',
		'└': '┘', '┘': '└',
	}
	if m, ok := mirrors[r]; ok {
		return m
	}
	return r
}

// Preset patterns
var PresetPatterns = map[string]string{
	"checker":   "█ ",
	"dots":      "· ",
	"grid":      "┼─",
	"wave":      "∿∾",
	"zigzag":    "/\\",
	"diamond":   "◇◆",
	"hearts":    "♥♡",
	"stars":     "★☆",
	"blocks":    "░▒▓█",
	"braille":   "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏",
	"shades":    "░▒▓",
	"arrows":    "→↘↓↙←↖↑↗",
	"box":       "┌┬┐├┼┤└┴┘",
	"circles":   "○◎●",
	"triangles": "▲△▼▽",
}

// GetPreset returns a preset pattern
func GetPreset(name string, width, height int) string {
	pattern, ok := PresetPatterns[name]
	if !ok {
		pattern = PresetPatterns["checker"]
	}
	return CreatePattern(pattern, width, height)
}

// ListBorders returns available border styles
func ListBorders() []string {
	return []string{"single", "double", "round", "thick", "ascii", "stars", "dots", "blocks", "shadow", "fancy"}
}

// ListDividers returns available divider styles
func ListDividers() []string {
	return []string{"single", "double", "thick", "dotted", "dashed", "wavy", "stars", "hearts", "arrows", "blocks", "gradient", "sparkle", "flower", "music", "snow", "zigzag"}
}

// ListPatterns returns available pattern presets
func ListPatterns() []string {
	return []string{"checker", "dots", "grid", "wave", "zigzag", "diamond", "hearts", "stars", "blocks", "braille", "shades", "arrows", "box", "circles", "triangles"}
}
