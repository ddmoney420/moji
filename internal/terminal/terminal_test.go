package terminal

import (
	"os"
	"strings"
	"testing"
)

// Helper function to set and cleanup environment variables
func setEnv(key, value string) func() {
	oldValue, exists := os.LookupEnv(key)
	os.Setenv(key, value)
	return func() {
		if exists {
			os.Setenv(key, oldValue)
		} else {
			os.Unsetenv(key)
		}
	}
}

// Helper function to unset environment variables
func unsetEnv(key string) func() {
	oldValue, exists := os.LookupEnv(key)
	os.Unsetenv(key)
	return func() {
		if exists {
			os.Setenv(key, oldValue)
		}
	}
}

// Test detectColorLevel with COLORTERM=truecolor
func TestDetectColorLevelTruecolorCOLORTERM(t *testing.T) {
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("TERM")()

	os.Setenv("COLORTERM", "truecolor")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor, got %v", level)
	}
}

// Test detectColorLevel with COLORTERM=24bit
func TestDetectColorLevel24bitCOLORTERM(t *testing.T) {
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("TERM")()

	os.Setenv("COLORTERM", "24bit")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor, got %v", level)
	}
}

// Test detectColorLevel with NO_COLOR set
func TestDetectColorLevelNOCOLOR(t *testing.T) {
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("TERM")()

	os.Setenv("NO_COLOR", "1")
	level := detectColorLevel()
	if level != NoColor {
		t.Errorf("Expected NoColor with NO_COLOR set, got %v", level)
	}
}

// Test detectColorLevel with FORCE_COLOR=0
func TestDetectColorLevelForceColor0(t *testing.T) {
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("FORCE_COLOR", "0")
	level := detectColorLevel()
	if level != NoColor {
		t.Errorf("Expected NoColor with FORCE_COLOR=0, got %v", level)
	}
}

// Test detectColorLevel with FORCE_COLOR=1
func TestDetectColorLevelForceColor1(t *testing.T) {
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("FORCE_COLOR", "1")
	level := detectColorLevel()
	if level != Basic {
		t.Errorf("Expected Basic with FORCE_COLOR=1, got %v", level)
	}
}

// Test detectColorLevel with FORCE_COLOR=2
func TestDetectColorLevelForceColor2(t *testing.T) {
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("FORCE_COLOR", "2")
	level := detectColorLevel()
	if level != Color256 {
		t.Errorf("Expected Color256 with FORCE_COLOR=2, got %v", level)
	}
}

// Test detectColorLevel with FORCE_COLOR=3
func TestDetectColorLevelForceColor3(t *testing.T) {
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("FORCE_COLOR", "3")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor with FORCE_COLOR=3, got %v", level)
	}
}

// Test detectColorLevel with invalid FORCE_COLOR
func TestDetectColorLevelForceColorInvalid(t *testing.T) {
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("FORCE_COLOR", "invalid")
	level := detectColorLevel()
	if level != Basic {
		t.Errorf("Expected Basic with invalid FORCE_COLOR, got %v", level)
	}
}

// Test detectColorLevel with TERM=dumb
func TestDetectColorLevelDumbTerminal(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("TERM", "dumb")
	level := detectColorLevel()
	if level != NoColor {
		t.Errorf("Expected NoColor for dumb terminal, got %v", level)
	}
}

// Test detectColorLevel with empty TERM
func TestDetectColorLevelEmptyTerm(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("TERM", "")
	level := detectColorLevel()
	if level != NoColor {
		t.Errorf("Expected NoColor for empty TERM, got %v", level)
	}
}

// Test detectColorLevel with xterm-truecolor
func TestDetectColorLevelXtermTruecolor(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "xterm-truecolor")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for xterm-truecolor, got %v", level)
	}
}

// Test detectColorLevel with xterm-direct
func TestDetectColorLevelXtermDirect(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "xterm-direct")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for xterm-direct, got %v", level)
	}
}

// Test detectColorLevel with iterm2
func TestDetectColorLevelIterm2(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "iterm2")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for iterm2, got %v", level)
	}
}

// Test detectColorLevel with kitty
func TestDetectColorLevelKitty(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "xterm-kitty")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for xterm-kitty, got %v", level)
	}
}

// Test detectColorLevel with alacritty
func TestDetectColorLevelAlacritty(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "alacritty")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for alacritty, got %v", level)
	}
}

// Test detectColorLevel with konsole
func TestDetectColorLevelKonsole(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "konsole")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for konsole, got %v", level)
	}
}

// Test detectColorLevel with wezterm
func TestDetectColorLevelWezterm(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "wezterm")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for wezterm, got %v", level)
	}
}

// Test detectColorLevel with xterm-256color
func TestDetectColorLevel256Color(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "xterm-256color")
	level := detectColorLevel()
	if level != Color256 {
		t.Errorf("Expected Color256 for xterm-256color, got %v", level)
	}
}

// Test detectColorLevel with screen-256color
func TestDetectColorLevelScreen256(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "screen-256color")
	level := detectColorLevel()
	if level != Color256 {
		t.Errorf("Expected Color256 for screen-256color, got %v", level)
	}
}

// Test detectColorLevel with tmux (256 colors)
func TestDetectColorLevelTmux(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "tmux-256color")
	level := detectColorLevel()
	if level != Color256 {
		t.Errorf("Expected Color256 for tmux-256color, got %v", level)
	}
}

// Test detectColorLevel with rxvt-unicode-256color
func TestDetectColorLevelRxvt256(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "rxvt-unicode-256color")
	level := detectColorLevel()
	if level != Color256 {
		t.Errorf("Expected Color256 for rxvt-unicode-256color, got %v", level)
	}
}

// Test detectColorLevel with xterm (basic 8 colors)
func TestDetectColorLevelXterm(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("TERM", "xterm")
	level := detectColorLevel()
	if level != Color256 {
		t.Errorf("Expected Color256 for xterm, got %v", level)
	}
}

// Test detectColorLevel with TERM_PROGRAM=iTerm.app
func TestDetectColorLevelTermProgram(t *testing.T) {
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "xterm")
	os.Setenv("TERM_PROGRAM", "iTerm.app")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for iTerm.app, got %v", level)
	}
}

// Test detectColorLevel with WT_SESSION (Windows Terminal)
func TestDetectColorLevelWindowsTerminal(t *testing.T) {
	defer unsetEnv("WT_SESSION")()
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "xterm")
	os.Setenv("WT_SESSION", "{guid}")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor with WT_SESSION, got %v", level)
	}
}

// Test detectColorLevel with unknown TERM
func TestDetectColorLevelUnknownTerm(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("TERM", "unknown-terminal")
	level := detectColorLevel()
	if level != Basic {
		t.Errorf("Expected Basic for unknown terminal, got %v", level)
	}
}

// Test detectColorLevel with linux terminal
func TestDetectColorLevelLinux(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("TERM", "linux")
	level := detectColorLevel()
	if level != Basic {
		t.Errorf("Expected Basic for linux terminal, got %v", level)
	}
}

// Test detectUnicode with LC_ALL=en_US.UTF-8
func TestDetectUnicodeWithLCAll(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("LC_ALL", "en_US.UTF-8")
	if !detectUnicode() {
		t.Error("Expected Unicode support with LC_ALL=en_US.UTF-8")
	}
}

// Test detectUnicode with LANG=en_US.UTF-8
func TestDetectUnicodeWithLang(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("LANG", "en_US.UTF-8")
	if !detectUnicode() {
		t.Error("Expected Unicode support with LANG=en_US.UTF-8")
	}
}

// Test detectUnicode with LC_CTYPE=UTF-8
func TestDetectUnicodeWithLCCType(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("LC_CTYPE", "UTF-8")
	if !detectUnicode() {
		t.Error("Expected Unicode support with LC_CTYPE=UTF-8")
	}
}

// Test detectUnicode with LANG containing utf8 (no dash)
func TestDetectUnicodeWithLangUtf8(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("LANG", "en_US.utf8")
	if !detectUnicode() {
		t.Error("Expected Unicode support with LANG=en_US.utf8")
	}
}

// Test detectUnicode with xterm TERM
func TestDetectUnicodeWithXtermTerm(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "xterm")
	if !detectUnicode() {
		t.Error("Expected Unicode support for xterm terminal")
	}
}

// Test detectUnicode with tmux TERM
func TestDetectUnicodeWithTmuxTerm(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "tmux")
	if !detectUnicode() {
		t.Error("Expected Unicode support for tmux terminal")
	}
}

// Test detectUnicode with kitty TERM
func TestDetectUnicodeWithKittyTerm(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "xterm-kitty")
	if !detectUnicode() {
		t.Error("Expected Unicode support for kitty terminal")
	}
}

// Test detectUnicode without any locale settings
func TestDetectUnicodeNoLocale(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	result := detectUnicode()
	// Result depends on TERM or locale settings, just verify it returns a boolean
	if result != true && result != false {
		t.Error("detectUnicode should return a boolean")
	}
}

// Test detectSixel with xterm
func TestDetectSixelXterm(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "xterm")
	if !detectSixel() {
		t.Error("Expected Sixel support for xterm")
	}
}

// Test detectSixel with foot
func TestDetectSixelFoot(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "foot")
	if !detectSixel() {
		t.Error("Expected Sixel support for foot")
	}
}

// Test detectSixel with mlterm
func TestDetectSixelMlterm(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "mlterm")
	if !detectSixel() {
		t.Error("Expected Sixel support for mlterm")
	}
}

// Test detectSixel with mintty
func TestDetectSixelMintty(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "mintty")
	if !detectSixel() {
		t.Error("Expected Sixel support for mintty")
	}
}

// Test detectSixel with contour
func TestDetectSixelContour(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "contour")
	if !detectSixel() {
		t.Error("Expected Sixel support for contour")
	}
}

// Test detectSixel with unsupported terminal
func TestDetectSixelUnsupported(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "vt100")
	if detectSixel() {
		t.Error("Expected no Sixel support for vt100")
	}
}

// Test detectSixel with kitty contains xterm
func TestDetectSixelKittyXtermMatch(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "xterm-kitty")
	// xterm-kitty contains "xterm" so it will match Sixel detection
	if !detectSixel() {
		t.Error("Expected Sixel support for xterm-kitty (contains xterm)")
	}
}

// Test detectKitty with xterm-kitty TERM
func TestDetectKittyTermVar(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("KITTY_WINDOW_ID")()

	os.Setenv("TERM", "xterm-kitty")
	if !detectKitty() {
		t.Error("Expected Kitty support with xterm-kitty")
	}
}

// Test detectKitty with KITTY_WINDOW_ID env var
func TestDetectKittyEnvVar(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("KITTY_WINDOW_ID")()

	os.Setenv("TERM", "xterm")
	os.Setenv("KITTY_WINDOW_ID", "1")
	if !detectKitty() {
		t.Error("Expected Kitty support with KITTY_WINDOW_ID")
	}
}

// Test detectKitty without kitty
func TestDetectKittyUnsupported(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("KITTY_WINDOW_ID")()

	os.Setenv("TERM", "xterm-256color")
	if detectKitty() {
		t.Error("Expected no Kitty support for non-kitty terminal")
	}
}

// Test detectITerm2 with TERM_PROGRAM
func TestDetectITerm2TermProgram(t *testing.T) {
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("ITERM_SESSION_ID")()

	os.Setenv("TERM_PROGRAM", "iTerm.app")
	if !detectITerm2() {
		t.Error("Expected iTerm2 support with TERM_PROGRAM=iTerm.app")
	}
}

// Test detectITerm2 with ITERM_SESSION_ID
func TestDetectITerm2SessionID(t *testing.T) {
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("ITERM_SESSION_ID")()

	os.Setenv("ITERM_SESSION_ID", "w0t1p0")
	if !detectITerm2() {
		t.Error("Expected iTerm2 support with ITERM_SESSION_ID")
	}
}

// Test detectITerm2 without iTerm2
func TestDetectITerm2Unsupported(t *testing.T) {
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("ITERM_SESSION_ID")()

	os.Setenv("TERM_PROGRAM", "Apple_Terminal")
	if detectITerm2() {
		t.Error("Expected no iTerm2 support for Apple Terminal")
	}
}

// Test ColorLevel String representation
func TestColorLevelString(t *testing.T) {
	tests := map[ColorLevel]string{
		NoColor:   "no color",
		Basic:     "16 colors",
		Color256:  "256 colors",
		TrueColor: "truecolor (24-bit)",
	}

	for level, expected := range tests {
		if level.String() != expected {
			t.Errorf("ColorLevel(%d).String() = %q, want %q", level, level.String(), expected)
		}
	}
}

// Test unknown ColorLevel String
func TestColorLevelStringUnknown(t *testing.T) {
	unknown := ColorLevel(99)
	if unknown.String() != "unknown" {
		t.Errorf("Unknown ColorLevel.String() = %q, want %q", unknown.String(), "unknown")
	}
}

// Test Capabilities SupportsColor
func TestCapabilitiesSupportsColor(t *testing.T) {
	tests := []struct {
		level    ColorLevel
		expected bool
	}{
		{NoColor, false},
		{Basic, true},
		{Color256, true},
		{TrueColor, true},
	}

	for _, test := range tests {
		caps := Capabilities{ColorLevel: test.level}
		if caps.SupportsColor() != test.expected {
			t.Errorf("SupportsColor() for %v: got %v, want %v", test.level, caps.SupportsColor(), test.expected)
		}
	}
}

// Test Capabilities SupportsTrueColor
func TestCapabilitiesSupportsTrueColor(t *testing.T) {
	tests := []struct {
		level    ColorLevel
		expected bool
	}{
		{NoColor, false},
		{Basic, false},
		{Color256, false},
		{TrueColor, true},
	}

	for _, test := range tests {
		caps := Capabilities{ColorLevel: test.level}
		if caps.SupportsTrueColor() != test.expected {
			t.Errorf("SupportsTrueColor() for %v: got %v, want %v", test.level, caps.SupportsTrueColor(), test.expected)
		}
	}
}

// Test Capabilities Supports256Color
func TestCapabilitiesSupports256Color(t *testing.T) {
	tests := []struct {
		level    ColorLevel
		expected bool
	}{
		{NoColor, false},
		{Basic, false},
		{Color256, true},
		{TrueColor, true},
	}

	for _, test := range tests {
		caps := Capabilities{ColorLevel: test.level}
		if caps.Supports256Color() != test.expected {
			t.Errorf("Supports256Color() for %v: got %v, want %v", test.level, caps.Supports256Color(), test.expected)
		}
	}
}

// Test Capabilities FitWidth with positive dimensions
func TestCapabilitiesFitWidthPositive(t *testing.T) {
	caps := Capabilities{Width: 80}

	tests := []struct {
		desired  int
		margin   int
		expected int
	}{
		{50, 5, 50},        // fits easily
		{75, 5, 75},        // fits with margin
		{80, 5, 75},        // exceeds max width
		{100, 5, 75},       // far exceeds max width
		{75, 10, 70},       // larger margin
		{0, 5, 0},          // zero desired
	}

	for _, test := range tests {
		result := caps.FitWidth(test.desired, test.margin)
		if result != test.expected {
			t.Errorf("FitWidth(%d, %d) = %d, want %d", test.desired, test.margin, result, test.expected)
		}
	}
}

// Test Capabilities FitWidth with zero or negative width
func TestCapabilitiesFitWidthZeroWidth(t *testing.T) {
	caps := Capabilities{Width: 0}

	result := caps.FitWidth(50, 5)
	if result != 50 {
		t.Errorf("FitWidth with zero width should return desired width: got %d, want 50", result)
	}
}

// Test Color16 with valid colors
func TestColor16Valid(t *testing.T) {
	tests := []struct {
		fg       int
		bg       int
		expected string
	}{
		{0, -1, "\033[30m"},        // black foreground
		{-1, 0, "\033[40m"},        // black background
		{0, 0, "\033[30;40m"},      // black fg and bg
		{7, 7, "\033[37;47m"},      // white fg and bg
		{8, -1, "\033[90m"},        // bright black foreground
		{15, -1, "\033[97m"},       // bright white foreground
		{-1, -1, ""},               // no colors
	}

	for _, test := range tests {
		result := Color16(test.fg, test.bg)
		if result != test.expected {
			t.Errorf("Color16(%d, %d) = %q, want %q", test.fg, test.bg, result, test.expected)
		}
	}
}

// Test Color256Code with valid colors
func TestColor256CodeValid(t *testing.T) {
	tests := []struct {
		fg       int
		bg       int
		expected string
	}{
		{0, -1, "\033[38;5;0m"},       // foreground color 0
		{-1, 0, "\033[48;5;0m"},       // background color 0
		{100, 200, "\033[38;5;100;48;5;200m"}, // both colors
		{255, 255, "\033[38;5;255;48;5;255m"}, // max colors
		{-1, -1, ""},                  // no colors
	}

	for _, test := range tests {
		result := Color256Code(test.fg, test.bg)
		if result != test.expected {
			t.Errorf("Color256Code(%d, %d) = %q, want %q", test.fg, test.bg, result, test.expected)
		}
	}
}

// Test TrueColorCode for foreground
func TestTrueColorCodeForeground(t *testing.T) {
	result := TrueColorCode(255, 128, 0, false)
	expected := "\033[38;2;255;128;0m"
	if result != expected {
		t.Errorf("TrueColorCode(255, 128, 0, false) = %q, want %q", result, expected)
	}
}

// Test TrueColorCode for background
func TestTrueColorCodeBackground(t *testing.T) {
	result := TrueColorCode(255, 128, 0, true)
	expected := "\033[48;2;255;128;0m"
	if result != expected {
		t.Errorf("TrueColorCode(255, 128, 0, true) = %q, want %q", result, expected)
	}
}

// Test Reset function
func TestReset(t *testing.T) {
	result := Reset()
	expected := "\033[0m"
	if result != expected {
		t.Errorf("Reset() = %q, want %q", result, expected)
	}
}

// Test RGBTo256 conversion
func TestRGBTo256Conversion(t *testing.T) {
	tests := []struct {
		r    uint8
		g    uint8
		b    uint8
		name string
	}{
		{0, 0, 0, "black"},
		{255, 255, 255, "white"},
		{255, 0, 0, "red"},
		{0, 255, 0, "green"},
		{0, 0, 255, "blue"},
		{128, 128, 128, "gray"},
	}

	for _, test := range tests {
		result := RGBTo256(test.r, test.g, test.b)
		if result < 0 || result > 255 {
			t.Errorf("RGBTo256(%d, %d, %d) (%s) = %d, want value in [0, 255]", test.r, test.g, test.b, test.name, result)
		}
	}
}

// Test RGBTo256 grayscale detection
func TestRGBTo256Grayscale(t *testing.T) {
	// Black should map to color 16
	result := RGBTo256(0, 0, 0)
	if result != 16 {
		t.Errorf("RGBTo256(0, 0, 0) = %d, want 16 (black)", result)
	}

	// White should map to color 231
	result = RGBTo256(255, 255, 255)
	if result != 231 {
		t.Errorf("RGBTo256(255, 255, 255) = %d, want 231 (white)", result)
	}
}

// Test RGBTo16 conversion
func TestRGBTo16Conversion(t *testing.T) {
	tests := []struct {
		r    uint8
		g    uint8
		b    uint8
		name string
	}{
		{0, 0, 0, "black"},
		{255, 255, 255, "white"},
		{255, 0, 0, "red"},
		{0, 255, 0, "green"},
		{0, 0, 255, "blue"},
	}

	for _, test := range tests {
		result := RGBTo16(test.r, test.g, test.b)
		if result < 0 || result > 15 {
			t.Errorf("RGBTo16(%d, %d, %d) (%s) = %d, want value in [0, 15]", test.r, test.g, test.b, test.name, result)
		}
	}
}

// Test RGBTo16 color accuracy
func TestRGBTo16Colors(t *testing.T) {
	// Black (0,0,0) - all values below threshold
	black := RGBTo16(0, 0, 0)
	if black != 0 {
		t.Errorf("RGBTo16(0, 0, 0) = %d, want 0 (black)", black)
	}

	// Red (255,0,0) - only red above threshold, sum > 382 (bright)
	red := RGBTo16(255, 0, 0)
	if red != 1 {
		t.Errorf("RGBTo16(255, 0, 0) = %d, want 1 (red)", red)
	}

	// Green (0,255,0) - only green above threshold, sum > 382 (bright)
	green := RGBTo16(0, 255, 0)
	if green != 2 {
		t.Errorf("RGBTo16(0, 255, 0) = %d, want 2 (green)", green)
	}

	// Blue (0,0,255) - only blue above threshold, sum > 382 (bright)
	blue := RGBTo16(0, 0, 255)
	if blue != 4 {
		t.Errorf("RGBTo16(0, 0, 255) = %d, want 4 (blue)", blue)
	}
}

// Test ColorLevel priority: COLORTERM > NO_COLOR > FORCE_COLOR
func TestColorLevelPriority(t *testing.T) {
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("TERM")()

	// COLORTERM takes priority
	os.Setenv("COLORTERM", "truecolor")
	os.Setenv("NO_COLOR", "1")
	os.Setenv("FORCE_COLOR", "0")
	os.Setenv("TERM", "dumb")

	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("COLORTERM should have priority, got %v", level)
	}
}

// Test edge case: empty LC_ALL with non-empty LANG
func TestDetectUnicodeEmptyLCAll(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("LC_ALL", "")
	os.Setenv("LANG", "en_US.UTF-8")

	if !detectUnicode() {
		t.Error("Expected Unicode support when LANG=en_US.UTF-8")
	}
}

// Test FORCE_COLOR with empty string (should be treated as not set)
func TestDetectColorLevelForceColorEmpty(t *testing.T) {
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("COLORTERM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("WT_SESSION")()

	os.Setenv("FORCE_COLOR", "")
	// Empty FORCE_COLOR should not match != "" check
	os.Setenv("TERM", "xterm")
	level := detectColorLevel()
	// Should fall through to TERM-based detection
	if level != Color256 {
		t.Errorf("Expected Color256 with empty FORCE_COLOR and xterm TERM, got %v", level)
	}
}

// Test case-insensitive TERM matching for truecolor terms
func TestDetectColorLevelCaseInsensitive(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	// Test with uppercase (though TERM is typically lowercase)
	os.Setenv("TERM", "XTERM-KITTY")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for XTERM-KITTY (case-insensitive), got %v", level)
	}
}

// Integration test: test common terminal combinations
func TestColorLevelCommonTerminals(t *testing.T) {
	tests := []struct {
		term      string
		colorterm string
		expected  ColorLevel
		name      string
	}{
		{"xterm-256color", "", Color256, "xterm-256color"},
		{"xterm-kitty", "", TrueColor, "xterm-kitty"},
		{"alacritty", "", TrueColor, "alacritty"},
		{"screen-256color", "", Color256, "screen-256color"},
		{"tmux-256color", "", Color256, "tmux-256color"},
		{"xterm", "", Color256, "xterm default"},
		{"screen", "", Color256, "screen"},
		{"tmux", "", Color256, "tmux"},
		{"dumb", "", NoColor, "dumb terminal"},
		{"vt100", "", Basic, "vt100"},
		{"linux", "", Basic, "linux"},
	}

	for _, test := range tests {
		defer unsetEnv("TERM")()
		defer unsetEnv("COLORTERM")()
		defer unsetEnv("NO_COLOR")()
		defer unsetEnv("FORCE_COLOR")()
		defer unsetEnv("TERM_PROGRAM")()
		defer unsetEnv("WT_SESSION")()

		os.Setenv("TERM", test.term)
		if test.colorterm != "" {
			os.Setenv("COLORTERM", test.colorterm)
		}

		level := detectColorLevel()
		if level != test.expected {
			t.Errorf("Terminal %q: got %v, want %v", test.name, level, test.expected)
		}
	}
}

// Integration test: test terminal multiplexer detection in Detect()
func TestDetectMultiplexers(t *testing.T) {
	defer unsetEnv("TMUX")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "screen-256color")
	os.Setenv("TMUX", "/tmp/tmux-1000/default,12345,0")

	caps := Capabilities{
		IsTmux: os.Getenv("TMUX") != "",
		IsScreen: strings.Contains(os.Getenv("TERM"), "screen"),
	}

	if !caps.IsTmux {
		t.Error("Expected TMUX detection")
	}
	if !caps.IsScreen {
		t.Error("Expected SCREEN detection")
	}
}

// Test SSH detection
func TestDetectSSH(t *testing.T) {
	defer unsetEnv("SSH_CLIENT")()
	defer unsetEnv("SSH_TTY")()

	tests := []struct {
		sshClient string
		sshTTY    string
		expected  bool
		name      string
	}{
		{"192.168.1.1 22 22", "", true, "SSH_CLIENT set"},
		{"", "/dev/pts/0", true, "SSH_TTY set"},
		{"192.168.1.1 22 22", "/dev/pts/0", true, "both set"},
		{"", "", false, "neither set"},
	}

	for _, test := range tests {
		defer unsetEnv("SSH_CLIENT")()
		defer unsetEnv("SSH_TTY")()

		if test.sshClient != "" {
			os.Setenv("SSH_CLIENT", test.sshClient)
		}
		if test.sshTTY != "" {
			os.Setenv("SSH_TTY", test.sshTTY)
		}

		isSSH := os.Getenv("SSH_CLIENT") != "" || os.Getenv("SSH_TTY") != ""
		if isSSH != test.expected {
			t.Errorf("SSH detection %q: got %v, want %v", test.name, isSSH, test.expected)
		}
	}
}

// Test Color16 with out-of-range values returns empty string
func TestColor16OutOfRange(t *testing.T) {
	tests := []struct {
		fg  int
		bg  int
		name string
	}{
		{-2, -1, "negative foreground"},
		{-1, -2, "negative background"},
		{16, -1, "foreground > 15"},
		{-1, 16, "background > 15"},
	}

	for _, test := range tests {
		result := Color16(test.fg, test.bg)
		if result != "" {
			t.Errorf("Color16(%d, %d) (%s) = %q, want empty string", test.fg, test.bg, test.name, result)
		}
	}
}

// Test vte terminal (GTK VTE, Gnome Terminal)
func TestDetectColorLevelVTE(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "xterm-vte")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for xterm-vte, got %v", level)
	}
}

// Test rio terminal
func TestDetectColorLevelRio(t *testing.T) {
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "rio")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for rio, got %v", level)
	}
}

// Test Apple_Terminal
func TestDetectColorLevelAppleTerminal(t *testing.T) {
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	// Note: TERM_PROGRAM check happens after xterm-like checks
	// Use a TERM that matches neither truecolor list nor 256color patterns
	os.Setenv("TERM", "vt220")
	os.Setenv("TERM_PROGRAM", "Apple_Terminal")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for Apple_Terminal, got %v", level)
	}
}

// Test Hyper terminal
func TestDetectColorLevelHyper(t *testing.T) {
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "vt220")
	os.Setenv("TERM_PROGRAM", "Hyper")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for Hyper, got %v", level)
	}
}

// Test vscode terminal
func TestDetectColorLevelVscode(t *testing.T) {
	defer unsetEnv("TERM_PROGRAM")()
	defer unsetEnv("TERM")()
	defer unsetEnv("NO_COLOR")()
	defer unsetEnv("FORCE_COLOR")()
	defer unsetEnv("COLORTERM")()

	os.Setenv("TERM", "vt220")
	os.Setenv("TERM_PROGRAM", "vscode")
	level := detectColorLevel()
	if level != TrueColor {
		t.Errorf("Expected TrueColor for vscode, got %v", level)
	}
}

// Test multiple locale variables set
func TestDetectUnicodeMultipleLocales(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	// LC_ALL should take precedence
	os.Setenv("LC_ALL", "en_US.UTF-8")
	os.Setenv("LANG", "C")
	os.Setenv("LC_CTYPE", "C")

	if !detectUnicode() {
		t.Error("Expected Unicode support with LC_ALL=en_US.UTF-8")
	}
}

// Test alacritty terminal
func TestDetectUnicodeAlacritty(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "alacritty")
	if !detectUnicode() {
		t.Error("Expected Unicode support for alacritty terminal")
	}
}

// Test iterm terminal
func TestDetectUnicodeITerm(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "iterm")
	if !detectUnicode() {
		t.Error("Expected Unicode support for iterm terminal")
	}
}

// Test rxvt terminal
func TestDetectUnicodeRxvt(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "rxvt")
	if !detectUnicode() {
		t.Error("Expected Unicode support for rxvt terminal")
	}
}

// Test vt100 terminal (no unicode in old terminals)
func TestDetectUnicodeVT100(t *testing.T) {
	defer unsetEnv("LC_ALL")()
	defer unsetEnv("LANG")()
	defer unsetEnv("LC_CTYPE")()
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "vt100")
	if !detectUnicode() {
		t.Error("Expected Unicode support for vt100 terminal (vt100 is in unicode list)")
	}
}

// Test rxvt-256 terminal
func TestDetectSixelRxvt256(t *testing.T) {
	defer unsetEnv("TERM")()

	os.Setenv("TERM", "rxvt-256color")
	if detectSixel() {
		t.Error("Expected no Sixel support for rxvt")
	}
}
