package terminal

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// ColorLevel represents terminal color capability
type ColorLevel int

const (
	NoColor   ColorLevel = 0
	Basic     ColorLevel = 1 // 16 colors
	Color256  ColorLevel = 2 // 256 colors
	TrueColor ColorLevel = 3 // 24-bit color
)

// Capabilities holds detected terminal capabilities
type Capabilities struct {
	ColorLevel    ColorLevel
	Width         int
	Height        int
	Unicode       bool
	Sixel         bool
	Kitty         bool
	ITerm2        bool
	WezTerm       bool
	Terminology   bool
	IsInteractive bool
	IsTmux        bool
	IsScreen      bool
	IsSSH         bool
	Term          string
}

// Detect detects terminal capabilities
func Detect() Capabilities {
	caps := Capabilities{
		Term:          os.Getenv("TERM"),
		IsInteractive: term.IsTerminal(int(os.Stdout.Fd())),
	}

	// Detect color level
	caps.ColorLevel = detectColorLevel()

	// Detect terminal size
	caps.Width, caps.Height = getSize()

	// Detect Unicode support
	caps.Unicode = detectUnicode()

	// Detect terminal multiplexers
	caps.IsTmux = os.Getenv("TMUX") != ""
	caps.IsScreen = strings.Contains(os.Getenv("TERM"), "screen")

	// Detect SSH
	caps.IsSSH = os.Getenv("SSH_CLIENT") != "" || os.Getenv("SSH_TTY") != ""

	// Detect graphics protocols
	caps.Sixel = detectSixel()
	caps.Kitty = detectKitty()
	caps.ITerm2 = detectITerm2()
	caps.WezTerm = detectWezTerm()
	caps.Terminology = detectTerminology()

	return caps
}

// detectColorLevel determines the color capability of the terminal
func detectColorLevel() ColorLevel {
	// Check COLORTERM first (highest priority)
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return TrueColor
	}

	// Check for explicit NO_COLOR
	if os.Getenv("NO_COLOR") != "" {
		return NoColor
	}

	// Check FORCE_COLOR
	forceColor := os.Getenv("FORCE_COLOR")
	if forceColor != "" {
		level, err := strconv.Atoi(forceColor)
		if err == nil {
			switch level {
			case 0:
				return NoColor
			case 1:
				return Basic
			case 2:
				return Color256
			case 3:
				return TrueColor
			}
		}
		// Non-empty FORCE_COLOR without valid number means basic colors
		return Basic
	}

	// Check TERM
	termVar := os.Getenv("TERM")
	if termVar == "" || termVar == "dumb" {
		return NoColor
	}

	// Known truecolor terminals
	trueColorTerms := []string{
		"xterm-truecolor",
		"xterm-direct",
		"iterm",
		"iterm2",
		"vte",
		"gnome",
		"konsole",
		"alacritty",
		"kitty",
		"wezterm",
		"rio",
	}
	for _, t := range trueColorTerms {
		if strings.Contains(strings.ToLower(termVar), t) {
			return TrueColor
		}
	}

	// Check for common 256-color terminals
	if strings.Contains(termVar, "256color") ||
		strings.Contains(termVar, "256") {
		return Color256
	}

	// Check terminal program env vars
	termProgram := os.Getenv("TERM_PROGRAM")
	trueColorPrograms := []string{
		"iTerm.app",
		"Apple_Terminal",
		"Hyper",
		"vscode",
		"Terminus",
	}
	for _, p := range trueColorPrograms {
		if termProgram == p {
			return TrueColor
		}
	}

	// Check WT_SESSION (Windows Terminal)
	if os.Getenv("WT_SESSION") != "" {
		return TrueColor
	}

	// Default to 256 colors for xterm-like terminals
	if strings.HasPrefix(termVar, "xterm") ||
		strings.HasPrefix(termVar, "screen") ||
		strings.HasPrefix(termVar, "tmux") ||
		strings.HasPrefix(termVar, "rxvt") {
		return Color256
	}

	// Fallback to basic colors
	return Basic
}

// getSize returns terminal width and height
func getSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Defaults
		return 80, 24
	}
	return width, height
}

// detectUnicode checks if terminal supports Unicode
func detectUnicode() bool {
	// Check LANG/LC_ALL for UTF-8
	lang := os.Getenv("LANG")
	lcAll := os.Getenv("LC_ALL")
	lcCtype := os.Getenv("LC_CTYPE")

	for _, l := range []string{lcAll, lcCtype, lang} {
		if strings.Contains(strings.ToUpper(l), "UTF-8") ||
			strings.Contains(strings.ToUpper(l), "UTF8") {
			return true
		}
	}

	// Modern terminals generally support Unicode
	termVar := os.Getenv("TERM")
	unicodeTerms := []string{
		"xterm", "rxvt", "screen", "tmux", "vt100",
		"linux", "alacritty", "kitty", "iterm", "konsole",
	}
	for _, t := range unicodeTerms {
		if strings.Contains(termVar, t) {
			return true
		}
	}

	return false
}

// detectSixel checks for Sixel graphics support
func detectSixel() bool {
	// Known Sixel-capable terminals
	termVar := os.Getenv("TERM")
	sixelTerms := []string{
		"xterm", "mlterm", "mintty", "foot", "contour",
	}
	for _, t := range sixelTerms {
		if strings.Contains(termVar, t) {
			return true
		}
	}
	return false
}

// detectKitty checks for Kitty graphics protocol support
func detectKitty() bool {
	termVar := os.Getenv("TERM")
	return strings.Contains(termVar, "kitty") ||
		os.Getenv("KITTY_WINDOW_ID") != ""
}

// detectITerm2 checks for iTerm2 inline images support
func detectITerm2() bool {
	termProgram := os.Getenv("TERM_PROGRAM")
	return termProgram == "iTerm.app" ||
		os.Getenv("ITERM_SESSION_ID") != ""
}

// detectWezTerm checks for WezTerm graphics protocol support
func detectWezTerm() bool {
	termProgram := os.Getenv("TERM_PROGRAM")
	return termProgram == "WezTerm" ||
		strings.Contains(os.Getenv("TERM"), "wezterm")
}

// detectTerminology checks for Terminology image protocol support
func detectTerminology() bool {
	// Terminology sets the TERMINOLOGY environment variable
	return os.Getenv("TERMINOLOGY") != ""
}

// ColorString returns a human-readable color level description
func (c ColorLevel) String() string {
	switch c {
	case NoColor:
		return "no color"
	case Basic:
		return "16 colors"
	case Color256:
		return "256 colors"
	case TrueColor:
		return "truecolor (24-bit)"
	default:
		return "unknown"
	}
}

// SupportsColor returns true if the terminal supports any color
func (c Capabilities) SupportsColor() bool {
	return c.ColorLevel > NoColor
}

// SupportsTrueColor returns true if the terminal supports 24-bit color
func (c Capabilities) SupportsTrueColor() bool {
	return c.ColorLevel >= TrueColor
}

// Supports256Color returns true if the terminal supports at least 256 colors
func (c Capabilities) Supports256Color() bool {
	return c.ColorLevel >= Color256
}

// FitWidth returns a width that fits the terminal (with margin)
func (c Capabilities) FitWidth(desired int, margin int) int {
	if c.Width <= 0 {
		return desired
	}
	maxWidth := c.Width - margin
	if desired > maxWidth {
		return maxWidth
	}
	return desired
}

// ANSI color code helpers

// Color16 returns ANSI escape for 16-color mode
func Color16(fg, bg int) string {
	if fg < 0 && bg < 0 {
		return ""
	}
	var codes []string
	if fg >= 0 && fg < 8 {
		codes = append(codes, strconv.Itoa(30+fg))
	} else if fg >= 8 && fg < 16 {
		codes = append(codes, strconv.Itoa(90+fg-8))
	}
	if bg >= 0 && bg < 8 {
		codes = append(codes, strconv.Itoa(40+bg))
	} else if bg >= 8 && bg < 16 {
		codes = append(codes, strconv.Itoa(100+bg-8))
	}
	if len(codes) == 0 {
		return ""
	}
	return "\033[" + strings.Join(codes, ";") + "m"
}

// Color256Code returns ANSI escape for 256-color mode
func Color256Code(fg, bg int) string {
	var parts []string
	if fg >= 0 && fg < 256 {
		parts = append(parts, "38;5;"+strconv.Itoa(fg))
	}
	if bg >= 0 && bg < 256 {
		parts = append(parts, "48;5;"+strconv.Itoa(bg))
	}
	if len(parts) == 0 {
		return ""
	}
	return "\033[" + strings.Join(parts, ";") + "m"
}

// TrueColorCode returns ANSI escape for 24-bit color
func TrueColorCode(r, g, b uint8, background bool) string {
	if background {
		return "\033[48;2;" + strconv.Itoa(int(r)) + ";" + strconv.Itoa(int(g)) + ";" + strconv.Itoa(int(b)) + "m"
	}
	return "\033[38;2;" + strconv.Itoa(int(r)) + ";" + strconv.Itoa(int(g)) + ";" + strconv.Itoa(int(b)) + "m"
}

// Reset returns ANSI reset code
func Reset() string {
	return "\033[0m"
}

// RGBTo256 converts RGB to closest 256-color palette index
func RGBTo256(r, g, b uint8) int {
	// Check grayscale ramp first (232-255)
	if r == g && g == b {
		if r < 8 {
			return 16 // black
		}
		if r > 248 {
			return 231 // white
		}
		return 232 + int((float64(r)-8)/10)
	}

	// 6x6x6 color cube (16-231)
	ri := int(float64(r)/255.0*5.0 + 0.5)
	gi := int(float64(g)/255.0*5.0 + 0.5)
	bi := int(float64(b)/255.0*5.0 + 0.5)

	return 16 + 36*ri + 6*gi + bi
}

// RGBTo16 converts RGB to closest 16-color palette index
func RGBTo16(r, g, b uint8) int {
	// Simple threshold-based conversion
	bright := (int(r) + int(g) + int(b)) > 382

	red := r > 127
	green := g > 127
	blue := b > 127

	color := 0
	if red {
		color |= 1
	}
	if green {
		color |= 2
	}
	if blue {
		color |= 4
	}

	if bright && color > 0 {
		color |= 8
	}

	return color
}
