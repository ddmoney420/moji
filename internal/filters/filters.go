package filters

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Filter represents a text filter function
type Filter func(string) string

// Registry of available filters
var registry = map[string]Filter{
	"metal":     Metal,
	"gay":       Rainbow,
	"rainbow":   Rainbow,
	"crop":      Crop,
	"flip":      Flip,
	"flop":      Flop,
	"rotate":    Rotate180,
	"border":    Border,
	"shadow":    Shadow,
	"3d":        Shadow3D,
	"glitch":    Glitch,
	"matrix":    Matrix,
	"fire":      Fire,
	"ice":       Ice,
	"neon":      Neon,
	"retro":     RetroGreen,
	"bold":      Bold,
	"italic":    Italic,
	"underline": Underline,
	"strike":    Strikethrough,
	"blink":     Blink,
	"dim":       Dim,
	"invert":    Invert,
}

// Get returns a filter by name
func Get(name string) (Filter, bool) {
	f, ok := registry[strings.ToLower(name)]
	return f, ok
}

// List returns available filter names
func List() []string {
	var names []string
	for k := range registry {
		names = append(names, k)
	}
	return names
}

// Chain applies multiple filters in sequence
func Chain(text string, filterNames []string) string {
	result := text
	for _, name := range filterNames {
		if f, ok := Get(name); ok {
			result = f(result)
		}
	}
	return result
}

// ParseChain parses comma-separated filter names
func ParseChain(spec string) []string {
	if spec == "" {
		return nil
	}
	parts := strings.Split(spec, ",")
	var filters []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			filters = append(filters, p)
		}
	}
	return filters
}

// Metal applies a metallic blue/gray effect
func Metal(text string) string {
	colors := []struct{ r, g, b uint8 }{
		{100, 100, 120},
		{140, 140, 160},
		{180, 180, 200},
		{220, 220, 240},
		{180, 180, 200},
		{140, 140, 160},
	}

	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		c := colors[i%len(colors)]
		for _, r := range line {
			if r != ' ' && r != '\t' {
				result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", c.r, c.g, c.b, r))
			} else {
				result.WriteRune(r)
			}
		}
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// Rainbow applies horizontal rainbow gradient (lolcat style)
func Rainbow(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	freq := 0.1
	for lineIdx, line := range lines {
		for i, r := range line {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				continue
			}
			phase := float64(i+lineIdx) * freq
			red := uint8(math.Sin(phase)*127 + 128)
			green := uint8(math.Sin(phase+2*math.Pi/3)*127 + 128)
			blue := uint8(math.Sin(phase+4*math.Pi/3)*127 + 128)
			result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", red, green, blue, r))
		}
		if lineIdx < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// Crop removes empty lines and leading/trailing whitespace
func Crop(text string) string {
	lines := strings.Split(text, "\n")
	var cropped []string

	// Remove empty lines at start/end
	start, end := 0, len(lines)-1
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	for end >= start && strings.TrimSpace(lines[end]) == "" {
		end--
	}

	if start > end {
		return ""
	}

	// Find minimum indent
	minIndent := -1
	for i := start; i <= end; i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}

	// Crop
	for i := start; i <= end; i++ {
		line := lines[i]
		if len(line) > minIndent {
			cropped = append(cropped, line[minIndent:])
		} else {
			cropped = append(cropped, "")
		}
	}

	return strings.Join(cropped, "\n")
}

// Flip flips text vertically (upside down)
func Flip(text string) string {
	lines := strings.Split(text, "\n")
	var flipped []string

	// Reverse line order
	for i := len(lines) - 1; i >= 0; i-- {
		flipped = append(flipped, lines[i])
	}

	return strings.Join(flipped, "\n")
}

// Flop mirrors text horizontally
func Flop(text string) string {
	lines := strings.Split(text, "\n")
	var flopped []string

	// Find max width
	maxWidth := 0
	for _, line := range lines {
		if len([]rune(line)) > maxWidth {
			maxWidth = len([]rune(line))
		}
	}

	for _, line := range lines {
		runes := []rune(line)
		// Pad to max width
		for len(runes) < maxWidth {
			runes = append(runes, ' ')
		}
		// Reverse
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		flopped = append(flopped, strings.TrimRight(string(runes), " "))
	}

	return strings.Join(flopped, "\n")
}

// Rotate180 rotates text 180 degrees
func Rotate180(text string) string {
	return Flip(Flop(text))
}

// Border adds a simple border around text
func Border(text string) string {
	lines := strings.Split(text, "\n")

	// Find max width
	maxWidth := 0
	for _, line := range lines {
		if len([]rune(line)) > maxWidth {
			maxWidth = len([]rune(line))
		}
	}

	var result strings.Builder
	border := strings.Repeat("─", maxWidth+2)
	result.WriteString("┌" + border + "┐\n")

	for _, line := range lines {
		runes := []rune(line)
		padding := maxWidth - len(runes)
		result.WriteString("│ ")
		result.WriteString(line)
		result.WriteString(strings.Repeat(" ", padding))
		result.WriteString(" │\n")
	}

	result.WriteString("└" + border + "┘")

	return result.String()
}

// Shadow adds a drop shadow effect
func Shadow(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	shadowColor := "\033[90m" // Dark gray
	reset := "\033[0m"

	for _, line := range lines {
		result.WriteString(line)
		result.WriteString("\n")
	}

	// Add shadow (offset by 1)
	for _, line := range lines {
		result.WriteString(" ")
		result.WriteString(shadowColor)
		for _, r := range line {
			if r != ' ' && r != '\t' {
				result.WriteRune('░')
			} else {
				result.WriteRune(' ')
			}
		}
		result.WriteString(reset)
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// Shadow3D creates a 3D shadow effect
func Shadow3D(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	colors := []string{
		"\033[38;2;60;60;60m",
		"\033[38;2;80;80;80m",
		"\033[38;2;100;100;100m",
	}
	reset := "\033[0m"

	// Draw shadows first (back to front)
	for layer := len(colors) - 1; layer >= 0; layer-- {
		for _, line := range lines {
			result.WriteString(strings.Repeat(" ", layer+1))
			result.WriteString(colors[layer])
			result.WriteString(line)
			result.WriteString(reset)
			result.WriteString("\n")
		}
	}

	// Draw main text on top
	for _, line := range lines {
		result.WriteString(line)
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// Glitch applies a digital glitch effect
func Glitch(text string) string {
	glitchChars := []rune("░▒▓█▄▀■□▪▫")
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		runes := []rune(line)
		for i, r := range runes {
			// Random glitch chance
			if rand.Float64() < 0.1 && r != ' ' && r != '\t' {
				// Glitch this character
				if rand.Float64() < 0.5 {
					result.WriteString("\033[31m") // Red
				} else {
					result.WriteString("\033[36m") // Cyan
				}
				result.WriteRune(glitchChars[rand.Intn(len(glitchChars))])
				result.WriteString("\033[0m")
			} else if rand.Float64() < 0.05 && i > 0 {
				// Offset glitch
				result.WriteRune(runes[i-1])
			} else {
				result.WriteRune(r)
			}
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// Matrix applies Matrix-style green effect
func Matrix(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	greens := []string{
		"\033[38;2;0;80;0m",
		"\033[38;2;0;120;0m",
		"\033[38;2;0;180;0m",
		"\033[38;2;0;220;0m",
		"\033[38;2;0;255;0m",
	}
	reset := "\033[0m"

	for _, line := range lines {
		for i, r := range line {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				continue
			}
			green := greens[rand.Intn(len(greens))]
			result.WriteString(green)
			result.WriteRune(r)
			result.WriteString(reset)
			_ = i
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// Fire applies a fire/flame color effect
func Fire(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for lineIdx, line := range lines {
		progress := float64(lineIdx) / float64(len(lines))
		for _, r := range line {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				continue
			}
			// Yellow at top, red at bottom
			red := uint8(255)
			green := uint8(255 * (1 - progress))
			blue := uint8(0)
			result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", red, green, blue, r))
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// Ice applies a cold/ice color effect
func Ice(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for lineIdx, line := range lines {
		progress := float64(lineIdx) / float64(len(lines))
		for _, r := range line {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				continue
			}
			// White at top, blue at bottom
			red := uint8(200 * (1 - progress))
			green := uint8(220*(1-progress) + 100*progress)
			blue := uint8(255)
			result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", red, green, blue, r))
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// Neon applies a neon glow effect
func Neon(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	neonColors := []struct{ r, g, b uint8 }{
		{255, 0, 255}, // Magenta
		{0, 255, 255}, // Cyan
		{255, 255, 0}, // Yellow
		{255, 0, 128}, // Hot pink
		{0, 255, 128}, // Spring green
	}

	colorIdx := 0
	for _, line := range lines {
		c := neonColors[colorIdx%len(neonColors)]
		for _, r := range line {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				continue
			}
			result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm\033[1m%c\033[0m", c.r, c.g, c.b, r))
		}
		result.WriteString("\n")
		colorIdx++
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// RetroGreen applies retro green terminal effect
func RetroGreen(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		for _, r := range line {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				continue
			}
			// Phosphor green with slight variation
			green := uint8(180 + rand.Intn(75))
			result.WriteString(fmt.Sprintf("\033[38;2;0;%d;0m%c\033[0m", green, r))
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// ANSI style filters

// Bold applies bold styling
func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

// Italic applies italic styling
func Italic(text string) string {
	return "\033[3m" + text + "\033[0m"
}

// Underline applies underline styling
func Underline(text string) string {
	return "\033[4m" + text + "\033[0m"
}

// Strikethrough applies strikethrough styling
func Strikethrough(text string) string {
	return "\033[9m" + text + "\033[0m"
}

// Blink applies blink effect
func Blink(text string) string {
	return "\033[5m" + text + "\033[0m"
}

// Dim applies dim/faint styling
func Dim(text string) string {
	return "\033[2m" + text + "\033[0m"
}

// Invert applies inverse colors
func Invert(text string) string {
	return "\033[7m" + text + "\033[0m"
}

// ListFilters returns filter names with descriptions
func ListFilters() []struct{ Name, Desc string } {
	return []struct{ Name, Desc string }{
		{"metal", "Metallic blue/gray effect"},
		{"rainbow", "Horizontal rainbow gradient"},
		{"gay", "Alias for rainbow"},
		{"crop", "Remove empty space around text"},
		{"flip", "Flip text vertically"},
		{"flop", "Mirror text horizontally"},
		{"rotate", "Rotate 180 degrees"},
		{"border", "Add border around text"},
		{"shadow", "Add drop shadow"},
		{"3d", "3D shadow effect"},
		{"glitch", "Digital glitch effect"},
		{"matrix", "Matrix green style"},
		{"fire", "Fire/flame colors"},
		{"ice", "Cold ice colors"},
		{"neon", "Neon glow effect"},
		{"retro", "Retro green terminal"},
		{"bold", "Bold text"},
		{"italic", "Italic text"},
		{"underline", "Underlined text"},
		{"strike", "Strikethrough text"},
		{"blink", "Blinking text"},
		{"dim", "Dim/faint text"},
		{"invert", "Inverse colors"},
	}
}
