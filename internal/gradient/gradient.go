package gradient

import (
	"fmt"
	"math"
	"strings"
)

// Theme defines a color theme
type Theme struct {
	Name   string
	Colors []Color
}

// Color represents an RGB color
type Color struct {
	R, G, B uint8
}

// Predefined themes
var Themes = map[string]Theme{
	"rainbow": {
		Name: "Rainbow",
		Colors: []Color{
			{255, 0, 0},   // Red
			{255, 127, 0}, // Orange
			{255, 255, 0}, // Yellow
			{0, 255, 0},   // Green
			{0, 0, 255},   // Blue
			{75, 0, 130},  // Indigo
			{148, 0, 211}, // Violet
		},
	},
	"neon": {
		Name: "Neon",
		Colors: []Color{
			{255, 0, 255}, // Magenta
			{0, 255, 255}, // Cyan
			{255, 255, 0}, // Yellow
			{255, 0, 128}, // Hot Pink
			{0, 255, 128}, // Spring Green
		},
	},
	"fire": {
		Name: "Fire",
		Colors: []Color{
			{255, 255, 0}, // Yellow
			{255, 200, 0}, // Gold
			{255, 128, 0}, // Orange
			{255, 0, 0},   // Red
			{128, 0, 0},   // Dark Red
		},
	},
	"ice": {
		Name: "Ice",
		Colors: []Color{
			{255, 255, 255}, // White
			{200, 230, 255}, // Light Blue
			{100, 180, 255}, // Sky Blue
			{0, 100, 255},   // Blue
			{0, 50, 150},    // Dark Blue
		},
	},
	"matrix": {
		Name: "Matrix",
		Colors: []Color{
			{0, 50, 0},      // Dark Green
			{0, 100, 0},     // Green
			{0, 200, 0},     // Bright Green
			{0, 255, 0},     // Lime
			{150, 255, 150}, // Light Green
		},
	},
	"sunset": {
		Name: "Sunset",
		Colors: []Color{
			{255, 100, 100}, // Light Red
			{255, 150, 50},  // Orange
			{255, 200, 100}, // Gold
			{200, 100, 150}, // Purple
			{100, 50, 150},  // Dark Purple
		},
	},
	"ocean": {
		Name: "Ocean",
		Colors: []Color{
			{0, 50, 100},    // Deep Blue
			{0, 100, 150},   // Ocean Blue
			{0, 150, 200},   // Sea Blue
			{100, 200, 230}, // Light Blue
			{200, 255, 255}, // Foam White
		},
	},
	"c64": {
		Name: "Commodore 64",
		Colors: []Color{
			{64, 50, 133},   // Purple
			{102, 90, 255},  // Light Blue
			{134, 199, 65},  // Light Green
			{255, 241, 224}, // Cream
		},
	},
	"dracula": {
		Name: "Dracula",
		Colors: []Color{
			{255, 121, 198}, // Pink
			{189, 147, 249}, // Purple
			{139, 233, 253}, // Cyan
			{80, 250, 123},  // Green
			{255, 184, 108}, // Orange
		},
	},
	"vaporwave": {
		Name: "Vaporwave",
		Colors: []Color{
			{255, 113, 206}, // Pink
			{185, 103, 255}, // Purple
			{1, 205, 254},   // Cyan
			{5, 255, 161},   // Green
			{255, 251, 150}, // Yellow
		},
	},
	"retro": {
		Name: "Retro Green",
		Colors: []Color{
			{0, 64, 0},
			{0, 128, 0},
			{0, 192, 0},
			{0, 255, 0},
			{128, 255, 128},
		},
	},
	"pastel": {
		Name: "Pastel",
		Colors: []Color{
			{255, 179, 186}, // Pink
			{255, 223, 186}, // Peach
			{255, 255, 186}, // Yellow
			{186, 255, 201}, // Green
			{186, 225, 255}, // Blue
		},
	},
}

// Apply applies a gradient to text
func Apply(text string, themeName string, mode string) string {
	theme, ok := Themes[themeName]
	if !ok {
		theme = Themes["rainbow"]
	}

	lines := strings.Split(text, "\n")
	var result strings.Builder

	totalChars := 0
	for _, line := range lines {
		totalChars += len([]rune(line))
	}

	charIndex := 0

	for _, line := range lines {
		runes := []rune(line)
		for _, r := range runes {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				charIndex++
				continue
			}

			var progress float64
			switch mode {
			case "horizontal":
				if len(runes) > 1 {
					progress = float64(charIndex) / float64(totalChars)
				}
			case "vertical":
				progress = float64(charIndex) / float64(totalChars)
			default: // diagonal
				progress = float64(charIndex) / float64(totalChars)
			}

			color := interpolateTheme(theme, progress)
			result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", color.R, color.G, color.B, r))
			charIndex++
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// ApplyPerLine applies gradient per line (resets each line)
func ApplyPerLine(text string, themeName string) string {
	theme, ok := Themes[themeName]
	if !ok {
		theme = Themes["rainbow"]
	}

	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		runes := []rune(line)
		for i, r := range runes {
			if r == ' ' || r == '\t' {
				result.WriteRune(r)
				continue
			}

			progress := 0.0
			if len(runes) > 1 {
				progress = float64(i) / float64(len(runes)-1)
			}

			color := interpolateTheme(theme, progress)
			result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", color.R, color.G, color.B, r))
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// interpolateTheme gets a color from theme at position t (0-1)
func interpolateTheme(theme Theme, t float64) Color {
	if len(theme.Colors) == 0 {
		return Color{255, 255, 255}
	}
	if len(theme.Colors) == 1 {
		return theme.Colors[0]
	}

	t = math.Max(0, math.Min(1, t))

	// Scale t to color segments
	segment := t * float64(len(theme.Colors)-1)
	idx := int(segment)
	if idx >= len(theme.Colors)-1 {
		return theme.Colors[len(theme.Colors)-1]
	}

	// Interpolate between two colors
	localT := segment - float64(idx)
	c1 := theme.Colors[idx]
	c2 := theme.Colors[idx+1]

	return Color{
		R: uint8(float64(c1.R) + localT*(float64(c2.R)-float64(c1.R))),
		G: uint8(float64(c1.G) + localT*(float64(c2.G)-float64(c1.G))),
		B: uint8(float64(c1.B) + localT*(float64(c2.B)-float64(c1.B))),
	}
}

// ListThemes returns available theme names
func ListThemes() []struct{ Name, Desc string } {
	return []struct{ Name, Desc string }{
		{"rainbow", "Classic rainbow gradient"},
		{"neon", "Bright neon colors"},
		{"fire", "Yellow to red fire"},
		{"ice", "White to blue ice"},
		{"matrix", "Green matrix style"},
		{"sunset", "Warm sunset colors"},
		{"ocean", "Deep blue ocean"},
		{"c64", "Commodore 64 palette"},
		{"dracula", "Dracula theme colors"},
		{"vaporwave", "80s vaporwave aesthetic"},
		{"retro", "Retro green terminal"},
		{"pastel", "Soft pastel colors"},
	}
}
