package styles

import (
	"strings"
)

// ANSI color codes - Basic
const (
	Reset       = "\x1b[0m"
	Bold        = "\x1b[1m"
	Dim         = "\x1b[2m"
	Italic      = "\x1b[3m"
	Underline   = "\x1b[4m"
	Blink       = "\x1b[5m"
	Reverse     = "\x1b[7m"
	Red         = "\x1b[91m"
	Yellow      = "\x1b[93m"
	Green       = "\x1b[92m"
	Cyan        = "\x1b[96m"
	Blue        = "\x1b[94m"
	Magenta     = "\x1b[95m"
	White       = "\x1b[97m"
	DarkYellow  = "\x1b[33m"
	DarkRed     = "\x1b[31m"
	DarkBlue    = "\x1b[34m"
	DarkGreen   = "\x1b[32m"
	DarkCyan    = "\x1b[36m"
	DarkMagenta = "\x1b[35m"
	Black       = "\x1b[30m"
	Gray        = "\x1b[90m"
)

// Rainbow applies rainbow colors to text
func Rainbow(text string) string {
	colors := []string{Red, Yellow, Green, Cyan, Blue, Magenta}
	var result strings.Builder
	colorIdx := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				result.WriteRune(ch)
			} else {
				result.WriteString(colors[colorIdx%len(colors)])
				result.WriteRune(ch)
				colorIdx++
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Gradient applies a top-to-bottom gradient
func Gradient(text string) string {
	colors := []string{Magenta, Blue, Cyan, Green, Yellow, Red}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		colorIdx := (i * len(colors)) / max(len(lines), 1)
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		result.WriteString(colors[colorIdx])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Fire applies a fire effect (yellow to red gradient)
func Fire(text string) string {
	colors := []string{Yellow, DarkYellow, Red, DarkRed}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		colorIdx := (i * len(colors)) / max(len(lines), 1)
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		result.WriteString(colors[colorIdx])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Ice applies an ice effect (white to blue gradient)
func Ice(text string) string {
	colors := []string{White, Cyan, Blue, DarkBlue}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		colorIdx := (i * len(colors)) / max(len(lines), 1)
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		result.WriteString(colors[colorIdx])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Matrix applies a matrix green effect
func Matrix(text string) string {
	colors := []string{Green, DarkGreen, Green, DarkGreen}
	var result strings.Builder
	colorIdx := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				result.WriteRune(ch)
			} else {
				result.WriteString(colors[colorIdx%len(colors)])
				result.WriteRune(ch)
				colorIdx++
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Neon applies a bright neon pink/cyan effect
func Neon(text string) string {
	colors := []string{Magenta, Cyan, Magenta, Cyan}
	var result strings.Builder
	colorIdx := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				result.WriteRune(ch)
			} else {
				result.WriteString(Bold)
				result.WriteString(colors[colorIdx%len(colors)])
				result.WriteRune(ch)
				colorIdx++
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Ocean applies an ocean blue gradient
func Ocean(text string) string {
	colors := []string{Cyan, Blue, DarkBlue, DarkCyan, Blue, Cyan}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		colorIdx := (i * len(colors)) / max(len(lines), 1)
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		result.WriteString(colors[colorIdx])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Sunset applies a sunset gradient (purple -> orange -> red)
func Sunset(text string) string {
	colors := []string{Magenta, DarkMagenta, Red, DarkRed, DarkYellow, Yellow}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		colorIdx := (i * len(colors)) / max(len(lines), 1)
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		result.WriteString(colors[colorIdx])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Cyberpunk applies a cyberpunk neon effect (magenta/cyan alternating lines)
func Cyberpunk(text string) string {
	colors := []string{Magenta, Cyan}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		result.WriteString(Bold)
		result.WriteString(colors[i%len(colors)])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Lava applies a lava effect (red/orange/black)
func Lava(text string) string {
	colors := []string{Red, DarkRed, DarkYellow, Red, DarkRed}
	var result strings.Builder
	colorIdx := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				result.WriteRune(ch)
			} else {
				result.WriteString(colors[colorIdx%len(colors)])
				result.WriteRune(ch)
				colorIdx++
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Toxic applies a toxic green/yellow effect
func Toxic(text string) string {
	colors := []string{Green, Yellow, DarkGreen, DarkYellow}
	var result strings.Builder
	colorIdx := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				result.WriteRune(ch)
			} else {
				result.WriteString(colors[colorIdx%len(colors)])
				result.WriteRune(ch)
				colorIdx++
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Galaxy applies a galaxy purple/blue effect
func Galaxy(text string) string {
	colors := []string{Magenta, DarkMagenta, Blue, DarkBlue, Magenta}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		colorIdx := (i * len(colors)) / max(len(lines), 1)
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		result.WriteString(colors[colorIdx])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Gold applies a gold/bronze gradient
func Gold(text string) string {
	colors := []string{Yellow, DarkYellow, Yellow, DarkYellow}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		result.WriteString(Bold)
		result.WriteString(colors[i%len(colors)])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Hacker applies a hacker-style dim green
func Hacker(text string) string {
	var result strings.Builder

	for _, line := range strings.Split(text, "\n") {
		result.WriteString(DarkGreen)
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Vaporwave applies a vaporwave aesthetic (pink/cyan)
func Vaporwave(text string) string {
	colors := []string{Magenta, Cyan, Magenta, Cyan}
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		result.WriteString(colors[i%len(colors)])
		result.WriteString(line)
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Christmas applies red/green alternating
func Christmas(text string) string {
	colors := []string{Red, Green}
	var result strings.Builder
	colorIdx := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				result.WriteRune(ch)
			} else {
				result.WriteString(Bold)
				result.WriteString(colors[colorIdx%len(colors)])
				result.WriteRune(ch)
				colorIdx++
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// USA applies red/white/blue
func USA(text string) string {
	colors := []string{Red, White, Blue}
	var result strings.Builder
	colorIdx := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				result.WriteRune(ch)
			} else {
				result.WriteString(Bold)
				result.WriteString(colors[colorIdx%len(colors)])
				result.WriteRune(ch)
				colorIdx++
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(Reset)
	return result.String()
}

// Mono applies a single bright white
func Mono(text string) string {
	var result strings.Builder
	result.WriteString(Bold)
	result.WriteString(White)
	result.WriteString(text)
	result.WriteString(Reset)
	return result.String()
}

// Apply applies a named style to text
func Apply(text, style string) string {
	switch strings.ToLower(style) {
	case "rainbow":
		return Rainbow(text)
	case "gradient":
		return Gradient(text)
	case "fire":
		return Fire(text)
	case "ice":
		return Ice(text)
	case "matrix":
		return Matrix(text)
	case "neon":
		return Neon(text)
	case "ocean":
		return Ocean(text)
	case "sunset":
		return Sunset(text)
	case "cyberpunk", "cyber":
		return Cyberpunk(text)
	case "lava":
		return Lava(text)
	case "toxic":
		return Toxic(text)
	case "galaxy":
		return Galaxy(text)
	case "gold":
		return Gold(text)
	case "hacker":
		return Hacker(text)
	case "vaporwave", "vapor":
		return Vaporwave(text)
	case "christmas", "xmas":
		return Christmas(text)
	case "usa", "america":
		return USA(text)
	case "mono", "white":
		return Mono(text)
	default:
		return text
	}
}

// StyleInfo contains style information
type StyleInfo struct {
	Name string
	Desc string
}

// ListStyles returns available styles with descriptions
func ListStyles() []StyleInfo {
	return []StyleInfo{
		{"none", "No color (default)"},
		{"rainbow", "Rainbow colors"},
		{"gradient", "Top-to-bottom gradient"},
		{"fire", "Fire effect (yellow to red)"},
		{"ice", "Ice effect (white to blue)"},
		{"matrix", "Matrix green"},
		{"neon", "Bright neon pink/cyan"},
		{"ocean", "Ocean blue waves"},
		{"sunset", "Sunset purple to orange"},
		{"cyberpunk", "Cyberpunk neon lines"},
		{"lava", "Lava red/orange flow"},
		{"toxic", "Toxic green/yellow"},
		{"galaxy", "Galaxy purple/blue"},
		{"gold", "Gold/bronze metallic"},
		{"hacker", "Hacker dim green"},
		{"vaporwave", "Vaporwave aesthetic"},
		{"christmas", "Red/green holiday"},
		{"usa", "Red/white/blue patriotic"},
		{"mono", "Bright white monochrome"},
	}
}

// Border styles
type BorderStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

var (
	BorderSingle = BorderStyle{"┌", "┐", "└", "┘", "─", "│"}
	BorderDouble = BorderStyle{"╔", "╗", "╚", "╝", "═", "║"}
	BorderRound  = BorderStyle{"╭", "╮", "╰", "╯", "─", "│"}
	BorderBold   = BorderStyle{"┏", "┓", "┗", "┛", "━", "┃"}
	BorderAscii  = BorderStyle{"+", "+", "+", "+", "-", "|"}
	BorderStars  = BorderStyle{"*", "*", "*", "*", "*", "*"}
	BorderHash   = BorderStyle{"#", "#", "#", "#", "#", "#"}
)

// GetBorderStyle returns a border style by name
func GetBorderStyle(name string) BorderStyle {
	switch strings.ToLower(name) {
	case "single":
		return BorderSingle
	case "double":
		return BorderDouble
	case "round", "rounded":
		return BorderRound
	case "bold", "thick":
		return BorderBold
	case "ascii":
		return BorderAscii
	case "stars", "star":
		return BorderStars
	case "hash":
		return BorderHash
	default:
		return BorderStyle{}
	}
}

// ApplyBorder adds a border around text
func ApplyBorder(text, borderName string) string {
	if borderName == "" || borderName == "none" {
		return text
	}

	border := GetBorderStyle(borderName)
	if border.Horizontal == "" {
		return text
	}

	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")

	// Find max width
	maxWidth := 0
	for _, line := range lines {
		// Count runes, not bytes for proper width
		w := runeWidth(line)
		if w > maxWidth {
			maxWidth = w
		}
	}

	var result strings.Builder

	// Top border
	result.WriteString(border.TopLeft)
	for i := 0; i < maxWidth+2; i++ {
		result.WriteString(border.Horizontal)
	}
	result.WriteString(border.TopRight)
	result.WriteString("\n")

	// Content lines
	for _, line := range lines {
		result.WriteString(border.Vertical)
		result.WriteString(" ")
		result.WriteString(line)
		// Pad to max width
		padding := maxWidth - runeWidth(line)
		for i := 0; i < padding; i++ {
			result.WriteString(" ")
		}
		result.WriteString(" ")
		result.WriteString(border.Vertical)
		result.WriteString("\n")
	}

	// Bottom border
	result.WriteString(border.BottomLeft)
	for i := 0; i < maxWidth+2; i++ {
		result.WriteString(border.Horizontal)
	}
	result.WriteString(border.BottomRight)
	result.WriteString("\n")

	return result.String()
}

// runeWidth calculates display width (simplified - doesn't handle all unicode)
func runeWidth(s string) int {
	count := 0
	for range s {
		count++
	}
	return count
}

// ListBorders returns available border styles
func ListBorders() []StyleInfo {
	return []StyleInfo{
		{"none", "No border"},
		{"single", "Single line ┌─┐"},
		{"double", "Double line ╔═╗"},
		{"round", "Rounded corners ╭─╮"},
		{"bold", "Bold/thick ┏━┓"},
		{"ascii", "ASCII +--+"},
		{"stars", "Stars ****"},
		{"hash", "Hash ####"},
	}
}

// Align constants
const (
	AlignLeft   = "left"
	AlignCenter = "center"
	AlignRight  = "right"
)

// ApplyAlignment aligns text within a given width
func ApplyAlignment(text, align string, width int) string {
	if width <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		lineWidth := runeWidth(line)
		if lineWidth >= width {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		padding := width - lineWidth
		switch strings.ToLower(align) {
		case AlignCenter:
			leftPad := padding / 2
			rightPad := padding - leftPad
			result.WriteString(strings.Repeat(" ", leftPad))
			result.WriteString(line)
			result.WriteString(strings.Repeat(" ", rightPad))
		case AlignRight:
			result.WriteString(strings.Repeat(" ", padding))
			result.WriteString(line)
		default: // left
			result.WriteString(line)
			result.WriteString(strings.Repeat(" ", padding))
		}
		result.WriteString("\n")
	}

	return strings.TrimRight(result.String(), "\n") + "\n"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
