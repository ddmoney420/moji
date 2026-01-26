package ux

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Version information (set via ldflags or defaults)
var (
	Version   = "1.0.1"
	BuildDate = "2026-01-25"
	GitCommit = "dev"
)

// Global settings
var (
	Quiet   bool
	Verbose bool
	NoColor bool
)

// Color codes
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Dim        = "\033[2m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldCyan   = "\033[1;36m"
)

func init() {
	// Check NO_COLOR environment variable (https://no-color.org/)
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		NoColor = true
	}
	// Also check TERM=dumb
	if os.Getenv("TERM") == "dumb" {
		NoColor = true
	}
}

// IsTTY returns true if stdout is a terminal
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// IsStderrTTY returns true if stderr is a terminal
func IsStderrTTY() bool {
	return term.IsTerminal(int(os.Stderr.Fd()))
}

// ColorEnabled returns true if colors should be used
func ColorEnabled() bool {
	return !NoColor && IsTTY()
}

// C returns the color code if colors are enabled, empty string otherwise
func C(code string) string {
	if ColorEnabled() {
		return code
	}
	return ""
}

// VersionString returns formatted version info
func VersionString() string {
	return fmt.Sprintf("moji %s (built %s, %s)", Version, BuildDate, GitCommit)
}

// VersionLong returns detailed version info
func VersionLong() string {
	var sb strings.Builder
	sb.WriteString(C(BoldCyan))
	sb.WriteString("moji")
	sb.WriteString(C(Reset))
	sb.WriteString(" ")
	sb.WriteString(C(Bold))
	sb.WriteString(Version)
	sb.WriteString(C(Reset))
	sb.WriteString("\n")
	sb.WriteString(C(Dim))
	sb.WriteString(fmt.Sprintf("Build date: %s\n", BuildDate))
	sb.WriteString(fmt.Sprintf("Git commit: %s\n", GitCommit))
	sb.WriteString(C(Reset))
	return sb.String()
}

// Error prints an error message to stderr
func Error(format string, args ...interface{}) {
	if Quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if IsStderrTTY() && !NoColor {
		fmt.Fprintf(os.Stderr, "%s%s%s %s\n", BoldRed, "Error:", Reset, msg)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
}

// ErrorWithSuggestion prints an error with a suggestion
func ErrorWithSuggestion(err string, suggestion string) {
	if Quiet {
		return
	}
	if IsStderrTTY() && !NoColor {
		fmt.Fprintf(os.Stderr, "%s%s%s %s\n", BoldRed, "Error:", Reset, err)
		if suggestion != "" {
			fmt.Fprintf(os.Stderr, "%s%s%s %s\n", Yellow, "Hint:", Reset, suggestion)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		if suggestion != "" {
			fmt.Fprintf(os.Stderr, "Hint: %s\n", suggestion)
		}
	}
}

// ErrorWithCommand prints an error with a command suggestion
func ErrorWithCommand(err string, cmd string) {
	if Quiet {
		return
	}
	if IsStderrTTY() && !NoColor {
		fmt.Fprintf(os.Stderr, "%s%s%s %s\n", BoldRed, "Error:", Reset, err)
		fmt.Fprintf(os.Stderr, "%s%s%s %s%s%s\n", Dim, "Try:", Reset, Cyan, cmd, Reset)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		fmt.Fprintf(os.Stderr, "Try: %s\n", cmd)
	}
}

// Warn prints a warning message to stderr
func Warn(format string, args ...interface{}) {
	if Quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if IsStderrTTY() && !NoColor {
		fmt.Fprintf(os.Stderr, "%s%s%s %s\n", BoldYellow, "Warning:", Reset, msg)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", msg)
	}
}

// Info prints an info message to stderr (only in verbose mode)
func Info(format string, args ...interface{}) {
	if Quiet || !Verbose {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if IsStderrTTY() && !NoColor {
		fmt.Fprintf(os.Stderr, "%s%s%s %s\n", BoldBlue, "Info:", Reset, msg)
	} else {
		fmt.Fprintf(os.Stderr, "Info: %s\n", msg)
	}
}

// Debug prints a debug message to stderr (only in verbose mode)
func Debug(format string, args ...interface{}) {
	if !Verbose {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if IsStderrTTY() && !NoColor {
		fmt.Fprintf(os.Stderr, "%s[debug]%s %s\n", Dim, Reset, msg)
	} else {
		fmt.Fprintf(os.Stderr, "[debug] %s\n", msg)
	}
}

// Success prints a success message to stderr
func Success(format string, args ...interface{}) {
	if Quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if IsStderrTTY() && !NoColor {
		fmt.Fprintf(os.Stderr, "%s%s%s %s\n", BoldGreen, "Done:", Reset, msg)
	} else {
		fmt.Fprintf(os.Stderr, "Done: %s\n", msg)
	}
}

// Status prints a status message (non-quiet mode)
func Status(format string, args ...interface{}) {
	if Quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

// DidYouMean formats a "did you mean" suggestion
func DidYouMean(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}
	if ColorEnabled() {
		formatted := make([]string, len(suggestions))
		for i, s := range suggestions {
			formatted[i] = Cyan + s + Reset
		}
		return "Did you mean: " + strings.Join(formatted, ", ") + "?"
	}
	return "Did you mean: " + strings.Join(suggestions, ", ") + "?"
}

// Header prints a section header
func Header(title string) {
	if Quiet {
		return
	}
	if ColorEnabled() {
		fmt.Printf("%s%s%s\n", Bold, title, Reset)
	} else {
		fmt.Println(title)
	}
}

// SubHeader prints a sub-section header
func SubHeader(title string) {
	if Quiet {
		return
	}
	if ColorEnabled() {
		fmt.Printf("\n%s%s%s\n", BoldCyan, title, Reset)
	} else {
		fmt.Printf("\n%s\n", title)
	}
}

// ListItem prints a list item with optional description
func ListItem(name, desc string, nameWidth int) {
	if ColorEnabled() {
		fmt.Printf("  %s%-*s%s  %s%s%s\n", Cyan, nameWidth, name, Reset, Dim, desc, Reset)
	} else {
		fmt.Printf("  %-*s  %s\n", nameWidth, name, desc)
	}
}

// Divider prints a horizontal divider
func Divider(char string, width int) {
	fmt.Println(strings.Repeat(char, width))
}

// Box wraps text in a simple box
func Box(text string) string {
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	var sb strings.Builder
	border := "+" + strings.Repeat("-", maxLen+2) + "+"
	sb.WriteString(border + "\n")
	for _, line := range lines {
		sb.WriteString(fmt.Sprintf("| %-*s |\n", maxLen, line))
	}
	sb.WriteString(border)
	return sb.String()
}

// Highlight returns highlighted text
func Highlight(text string) string {
	if ColorEnabled() {
		return Bold + text + Reset
	}
	return text
}

// Code returns text formatted as code
func Code(text string) string {
	if ColorEnabled() {
		return Cyan + text + Reset
	}
	return "`" + text + "`"
}

// FormatExample returns a formatted example
func FormatExample(cmd, desc string) string {
	if ColorEnabled() {
		return fmt.Sprintf("  %s%s%s\n    %s%s%s", Cyan, cmd, Reset, Dim, desc, Reset)
	}
	return fmt.Sprintf("  %s\n    %s", cmd, desc)
}
