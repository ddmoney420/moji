package ux

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ddmoney420/moji/internal/terminal"
	"golang.org/x/term"
)

// Check represents a diagnostic check result
type Check struct {
	Name    string
	Status  CheckStatus
	Message string
	Detail  string
}

type CheckStatus int

const (
	CheckOK CheckStatus = iota
	CheckWarn
	CheckFail
	CheckSkip
)

func (s CheckStatus) String() string {
	switch s {
	case CheckOK:
		return "ok"
	case CheckWarn:
		return "warn"
	case CheckFail:
		return "fail"
	case CheckSkip:
		return "skip"
	default:
		return "?"
	}
}

func (s CheckStatus) Icon() string {
	if !ColorEnabled() {
		switch s {
		case CheckOK:
			return "[OK]"
		case CheckWarn:
			return "[WARN]"
		case CheckFail:
			return "[FAIL]"
		case CheckSkip:
			return "[SKIP]"
		default:
			return "[?]"
		}
	}
	switch s {
	case CheckOK:
		return BoldGreen + "✓" + Reset
	case CheckWarn:
		return BoldYellow + "!" + Reset
	case CheckFail:
		return BoldRed + "✗" + Reset
	case CheckSkip:
		return Dim + "-" + Reset
	default:
		return "?"
	}
}

// Doctor runs all diagnostic checks
func Doctor() []Check {
	var checks []Check

	// Version check
	checks = append(checks, Check{
		Name:    "Version",
		Status:  CheckOK,
		Message: fmt.Sprintf("moji %s", Version),
	})

	// Terminal check
	checks = append(checks, checkTerminal())

	// Color support
	checks = append(checks, checkColors())

	// Unicode support
	checks = append(checks, checkUnicode())

	// Clipboard
	checks = append(checks, checkClipboard())

	// Shell
	checks = append(checks, checkShell())

	// OS/Platform
	checks = append(checks, Check{
		Name:    "Platform",
		Status:  CheckOK,
		Message: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	})

	return checks
}

func checkTerminal() Check {
	caps := terminal.Detect()

	if !IsTTY() {
		return Check{
			Name:    "Terminal",
			Status:  CheckWarn,
			Message: "Not running in a TTY",
			Detail:  "Some features may be limited when piping output",
		}
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return Check{
			Name:    "Terminal",
			Status:  CheckWarn,
			Message: caps.Term,
			Detail:  "Could not detect terminal size",
		}
	}

	return Check{
		Name:    "Terminal",
		Status:  CheckOK,
		Message: fmt.Sprintf("%s (%dx%d)", caps.Term, width, height),
	}
}

func checkColors() Check {
	if NoColor {
		return Check{
			Name:    "Colors",
			Status:  CheckSkip,
			Message: "Disabled (NO_COLOR set)",
		}
	}

	caps := terminal.Detect()

	switch caps.ColorLevel {
	case terminal.TrueColor:
		return Check{
			Name:    "Colors",
			Status:  CheckOK,
			Message: "24-bit truecolor supported",
		}
	case terminal.Color256:
		return Check{
			Name:    "Colors",
			Status:  CheckOK,
			Message: "256 colors supported",
		}
	case terminal.Basic:
		return Check{
			Name:    "Colors",
			Status:  CheckWarn,
			Message: "Basic 16 colors only",
			Detail:  "Set COLORTERM=truecolor for full color support",
		}
	default:
		return Check{
			Name:    "Colors",
			Status:  CheckWarn,
			Message: "No color support detected",
		}
	}
}

func checkUnicode() Check {
	// Check LANG/LC_ALL for UTF-8
	lang := os.Getenv("LANG")
	lcAll := os.Getenv("LC_ALL")

	hasUTF8 := strings.Contains(strings.ToLower(lang), "utf-8") ||
		strings.Contains(strings.ToLower(lang), "utf8") ||
		strings.Contains(strings.ToLower(lcAll), "utf-8") ||
		strings.Contains(strings.ToLower(lcAll), "utf8")

	if hasUTF8 {
		return Check{
			Name:    "Unicode",
			Status:  CheckOK,
			Message: "UTF-8 locale configured",
		}
	}

	return Check{
		Name:    "Unicode",
		Status:  CheckWarn,
		Message: "UTF-8 locale not detected",
		Detail:  fmt.Sprintf("LANG=%s (some characters may not display correctly)", lang),
	}
}

func checkClipboard() Check {
	// Try to detect clipboard availability
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "pbcopy"
	case "linux":
		// Check for xclip or xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = "xclip"
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = "xsel"
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = "wl-copy"
		}
	case "windows":
		cmd = "clip"
	}

	if cmd == "" {
		return Check{
			Name:    "Clipboard",
			Status:  CheckWarn,
			Message: "Clipboard tool not found",
			Detail:  "Install xclip, xsel, or wl-copy for clipboard support",
		}
	}

	if _, err := exec.LookPath(cmd); err != nil {
		return Check{
			Name:    "Clipboard",
			Status:  CheckWarn,
			Message: fmt.Sprintf("%s not found", cmd),
		}
	}

	return Check{
		Name:    "Clipboard",
		Status:  CheckOK,
		Message: fmt.Sprintf("Available (%s)", cmd),
	}
}

func checkShell() Check {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "unknown"
	}

	// Extract shell name
	parts := strings.Split(shell, "/")
	shellName := parts[len(parts)-1]

	return Check{
		Name:    "Shell",
		Status:  CheckOK,
		Message: shellName,
		Detail:  fmt.Sprintf("Run 'moji completions %s' for tab completion", shellName),
	}
}

// FormatDoctorReport formats the doctor check results
func FormatDoctorReport(checks []Check) string {
	var sb strings.Builder

	sb.WriteString("\n")
	if ColorEnabled() {
		sb.WriteString(Bold + "moji doctor" + Reset + "\n")
	} else {
		sb.WriteString("moji doctor\n")
	}
	sb.WriteString(strings.Repeat("─", 40) + "\n\n")

	maxNameLen := 0
	for _, c := range checks {
		if len(c.Name) > maxNameLen {
			maxNameLen = len(c.Name)
		}
	}

	hasIssues := false
	for _, c := range checks {
		sb.WriteString(fmt.Sprintf("  %s %-*s  %s\n",
			c.Status.Icon(), maxNameLen, c.Name, c.Message))
		if c.Detail != "" {
			if ColorEnabled() {
				sb.WriteString(fmt.Sprintf("     %s%s%s\n", Dim, c.Detail, Reset))
			} else {
				sb.WriteString(fmt.Sprintf("     %s\n", c.Detail))
			}
		}
		if c.Status == CheckFail || c.Status == CheckWarn {
			hasIssues = true
		}
	}

	sb.WriteString("\n")
	if hasIssues {
		if ColorEnabled() {
			sb.WriteString(Yellow + "Some issues detected. See details above." + Reset + "\n")
		} else {
			sb.WriteString("Some issues detected. See details above.\n")
		}
	} else {
		if ColorEnabled() {
			sb.WriteString(Green + "All checks passed!" + Reset + "\n")
		} else {
			sb.WriteString("All checks passed!\n")
		}
	}

	return sb.String()
}
