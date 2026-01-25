package calendar

import (
	"fmt"
	"strings"
	"time"
)

// Options for calendar generation
type Options struct {
	Highlight   []time.Time // Days to highlight
	FirstDayMon bool        // Start week on Monday
	ShowWeekNum bool        // Show week numbers
	Compact     bool        // Compact format
}

// Month generates a calendar for a specific month
func Month(year int, month time.Month, opts Options) string {
	var sb strings.Builder

	// Header
	monthName := month.String()
	header := fmt.Sprintf("%s %d", monthName, year)
	padding := (20 - len(header)) / 2
	sb.WriteString(strings.Repeat(" ", padding))
	sb.WriteString("\033[1m")
	sb.WriteString(header)
	sb.WriteString("\033[0m")
	sb.WriteString("\n")

	// Day headers
	if opts.ShowWeekNum {
		sb.WriteString("   ")
	}
	days := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	if opts.FirstDayMon {
		days = []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	}
	for _, d := range days {
		sb.WriteString(d)
		sb.WriteString(" ")
	}
	sb.WriteString("\n")

	// Get first day of month
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Calculate starting position
	startWeekday := int(firstDay.Weekday())
	if opts.FirstDayMon {
		startWeekday = (startWeekday + 6) % 7
	}

	// Week number
	if opts.ShowWeekNum {
		_, week := firstDay.ISOWeek()
		sb.WriteString(fmt.Sprintf("%2d ", week))
	}

	// Padding for first week
	for i := 0; i < startWeekday; i++ {
		sb.WriteString("   ")
	}

	// Days
	today := time.Now()
	currentWeekday := startWeekday

	for day := 1; day <= lastDay.Day(); day++ {
		date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		isToday := date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day()
		isHighlighted := false
		for _, h := range opts.Highlight {
			if h.Year() == date.Year() && h.Month() == date.Month() && h.Day() == date.Day() {
				isHighlighted = true
				break
			}
		}

		if isToday {
			sb.WriteString("\033[7m") // Inverse
		} else if isHighlighted {
			sb.WriteString("\033[1;33m") // Bold yellow
		}

		sb.WriteString(fmt.Sprintf("%2d", day))

		if isToday || isHighlighted {
			sb.WriteString("\033[0m")
		}
		sb.WriteString(" ")

		currentWeekday++
		if currentWeekday == 7 {
			sb.WriteString("\n")
			currentWeekday = 0
			// Week number for next week
			if opts.ShowWeekNum && day < lastDay.Day() {
				nextDate := date.AddDate(0, 0, 1)
				_, week := nextDate.ISOWeek()
				sb.WriteString(fmt.Sprintf("%2d ", week))
			}
		}
	}

	if currentWeekday != 0 {
		sb.WriteString("\n")
	}

	return sb.String()
}

// Year generates a full year calendar
func Year(year int, opts Options) string {
	var sb strings.Builder

	// Year header
	yearHeader := fmt.Sprintf("%d", year)
	sb.WriteString(strings.Repeat(" ", 30))
	sb.WriteString("\033[1m")
	sb.WriteString(yearHeader)
	sb.WriteString("\033[0m")
	sb.WriteString("\n\n")

	// 4 rows of 3 months
	for row := 0; row < 4; row++ {
		months := make([]string, 3)
		for col := 0; col < 3; col++ {
			m := time.Month(row*3 + col + 1)
			months[col] = Month(year, m, opts)
		}

		// Combine side by side
		monthLines := make([][]string, 3)
		maxLines := 0
		for i, m := range months {
			monthLines[i] = strings.Split(m, "\n")
			if len(monthLines[i]) > maxLines {
				maxLines = len(monthLines[i])
			}
		}

		for line := 0; line < maxLines; line++ {
			for i := 0; i < 3; i++ {
				if line < len(monthLines[i]) {
					sb.WriteString(fmt.Sprintf("%-22s", monthLines[i][line]))
				} else {
					sb.WriteString(strings.Repeat(" ", 22))
				}
				if i < 2 {
					sb.WriteString("  ")
				}
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Current generates calendar for current month
func Current(opts Options) string {
	now := time.Now()
	return Month(now.Year(), now.Month(), opts)
}

// Today returns formatted today's date
func Today() string {
	now := time.Now()
	return now.Format("Monday, January 2, 2006")
}

// Countdown returns days until a date
func Countdown(target time.Time) string {
	now := time.Now()
	diff := target.Sub(now)

	if diff < 0 {
		return "Date has passed"
	}

	days := int(diff.Hours() / 24)
	hours := int(diff.Hours()) % 24

	if days == 0 {
		return fmt.Sprintf("%d hours remaining", hours)
	}
	if days == 1 {
		return "1 day remaining"
	}
	return fmt.Sprintf("%d days remaining", days)
}

// ASCIIArt returns the current date as ASCII art
func ASCIIArt() string {
	now := time.Now()
	day := now.Day()
	month := now.Month().String()[:3]
	year := now.Year()

	return fmt.Sprintf(`
╔═══════════════╗
║   %s %4d   ║
║               ║
║      %2d       ║
║               ║
╚═══════════════╝
`, month, year, day)
}

// WeekView returns current week view
func WeekView(opts Options) string {
	var sb strings.Builder
	now := time.Now()

	// Find start of week
	weekday := int(now.Weekday())
	if opts.FirstDayMon {
		weekday = (weekday + 6) % 7
	}
	startOfWeek := now.AddDate(0, 0, -weekday)

	// Header
	_, week := now.ISOWeek()
	sb.WriteString(fmt.Sprintf("\033[1mWeek %d\033[0m\n", week))
	sb.WriteString(strings.Repeat("─", 50))
	sb.WriteString("\n")

	// Days
	for i := 0; i < 7; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		isToday := day.Day() == now.Day() && day.Month() == now.Month()

		if isToday {
			sb.WriteString("\033[7m") // Inverse
		}

		sb.WriteString(fmt.Sprintf("%-10s %s",
			day.Weekday().String(),
			day.Format("Jan 02")))

		if isToday {
			sb.WriteString(" (today)")
			sb.WriteString("\033[0m")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
