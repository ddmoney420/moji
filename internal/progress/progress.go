package progress

import (
	"fmt"
	"strings"
	"time"
)

// Style defines a progress bar style
type Style struct {
	Left      string
	Right     string
	Fill      string
	Empty     string
	Head      string
	ShowPct   bool
	ShowCount bool
}

// Predefined styles
var Styles = map[string]Style{
	"default": {
		Left:    "[",
		Right:   "]",
		Fill:    "=",
		Empty:   " ",
		Head:    ">",
		ShowPct: true,
	},
	"blocks": {
		Left:    "│",
		Right:   "│",
		Fill:    "█",
		Empty:   "░",
		Head:    "█",
		ShowPct: true,
	},
	"smooth": {
		Left:    "│",
		Right:   "│",
		Fill:    "█",
		Empty:   " ",
		Head:    "▓",
		ShowPct: true,
	},
	"dots": {
		Left:    "⟨",
		Right:   "⟩",
		Fill:    "●",
		Empty:   "○",
		Head:    "●",
		ShowPct: true,
	},
	"arrows": {
		Left:    "├",
		Right:   "┤",
		Fill:    "─",
		Empty:   " ",
		Head:    "►",
		ShowPct: true,
	},
	"shades": {
		Left:    "",
		Right:   "",
		Fill:    "▓",
		Empty:   "░",
		Head:    "▓",
		ShowPct: true,
	},
	"braille": {
		Left:    "⣿",
		Right:   "⣿",
		Fill:    "⣿",
		Empty:   "⠀",
		Head:    "⣿",
		ShowPct: true,
	},
	"simple": {
		Left:    "",
		Right:   "",
		Fill:    "#",
		Empty:   "-",
		Head:    "#",
		ShowPct: true,
	},
	"minimal": {
		Left:    "",
		Right:   "",
		Fill:    "━",
		Empty:   "─",
		Head:    "●",
		ShowPct: false,
	},
}

// Bar creates a progress bar
func Bar(current, total int, width int, styleName string) string {
	style, ok := Styles[styleName]
	if !ok {
		style = Styles["default"]
	}

	if total <= 0 {
		total = 1
	}
	if current > total {
		current = total
	}
	if current < 0 {
		current = 0
	}

	pct := float64(current) / float64(total)
	barWidth := width - len(style.Left) - len(style.Right)
	if style.ShowPct {
		barWidth -= 5 // " 100%"
	}

	filled := int(pct * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	var sb strings.Builder
	sb.WriteString(style.Left)

	for i := 0; i < barWidth; i++ {
		if i < filled {
			sb.WriteString(style.Fill)
		} else if i == filled && filled < barWidth {
			sb.WriteString(style.Head)
		} else {
			sb.WriteString(style.Empty)
		}
	}

	sb.WriteString(style.Right)

	if style.ShowPct {
		sb.WriteString(fmt.Sprintf(" %3d%%", int(pct*100)))
	}

	return sb.String()
}

// Spinner creates an animated spinner frame
func Spinner(frame int, styleName string) string {
	spinners := map[string][]string{
		"default": {"|", "/", "-", "\\"},
		"dots":    {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		"line":    {"⎺", "⎻", "⎼", "⎽", "⎼", "⎻"},
		"bounce":  {"⠁", "⠂", "⠄", "⠂"},
		"grow":    {"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"},
		"box":     {"◰", "◳", "◲", "◱"},
		"circle":  {"◐", "◓", "◑", "◒"},
		"arrow":   {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
		"clock":   {"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"},
		"moon":    {"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
		"star":    {"✶", "✷", "✹", "✺"},
	}

	frames, ok := spinners[styleName]
	if !ok {
		frames = spinners["default"]
	}

	return frames[frame%len(frames)]
}

// BarWithLabel creates a labeled progress bar
func BarWithLabel(label string, current, total int, width int, styleName string) string {
	bar := Bar(current, total, width-len(label)-1, styleName)
	return fmt.Sprintf("%s %s", label, bar)
}

// MultiBar creates multiple progress bars
func MultiBar(items []struct {
	Label   string
	Current int
	Total   int
}, width int, styleName string) string {
	var sb strings.Builder

	// Find max label length
	maxLabel := 0
	for _, item := range items {
		if len(item.Label) > maxLabel {
			maxLabel = len(item.Label)
		}
	}

	for i, item := range items {
		label := fmt.Sprintf("%-*s", maxLabel, item.Label)
		bar := Bar(item.Current, item.Total, width-maxLabel-1, styleName)
		sb.WriteString(label)
		sb.WriteString(" ")
		sb.WriteString(bar)
		if i < len(items)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ETA calculates estimated time remaining
func ETA(current, total int, elapsed time.Duration) string {
	if current <= 0 {
		return "calculating..."
	}

	rate := float64(current) / elapsed.Seconds()
	remaining := float64(total-current) / rate

	if remaining < 60 {
		return fmt.Sprintf("%.0fs remaining", remaining)
	} else if remaining < 3600 {
		return fmt.Sprintf("%.0fm remaining", remaining/60)
	} else {
		return fmt.Sprintf("%.1fh remaining", remaining/3600)
	}
}

// Rate calculates items per second
func Rate(count int, elapsed time.Duration) string {
	if elapsed.Seconds() <= 0 {
		return "0/s"
	}
	rate := float64(count) / elapsed.Seconds()
	if rate >= 1000000 {
		return fmt.Sprintf("%.1fM/s", rate/1000000)
	} else if rate >= 1000 {
		return fmt.Sprintf("%.1fK/s", rate/1000)
	} else {
		return fmt.Sprintf("%.1f/s", rate)
	}
}

// ByteRate calculates bytes per second
func ByteRate(bytes int64, elapsed time.Duration) string {
	if elapsed.Seconds() <= 0 {
		return "0 B/s"
	}
	rate := float64(bytes) / elapsed.Seconds()

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case rate >= GB:
		return fmt.Sprintf("%.1f GB/s", rate/GB)
	case rate >= MB:
		return fmt.Sprintf("%.1f MB/s", rate/MB)
	case rate >= KB:
		return fmt.Sprintf("%.1f KB/s", rate/KB)
	default:
		return fmt.Sprintf("%.0f B/s", rate)
	}
}

// ListStyles returns available style names
func ListStyles() []string {
	var names []string
	for k := range Styles {
		names = append(names, k)
	}
	return names
}
