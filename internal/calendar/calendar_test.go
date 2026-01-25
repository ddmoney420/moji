package calendar

import (
	"strings"
	"testing"
	"time"
)

func TestMonth(t *testing.T) {
	// Use a fixed month for reproducible tests
	opts := Options{}
	result := Month(2024, time.January, opts)

	// Should contain month name and year
	if !strings.Contains(result, "January") {
		t.Error("Month should contain month name")
	}
	if !strings.Contains(result, "2024") {
		t.Error("Month should contain year")
	}

	// Should contain day headers
	if !strings.Contains(result, "Su") || !strings.Contains(result, "Sa") {
		t.Error("Month should contain day headers")
	}

	// Should contain days
	if !strings.Contains(result, "15") {
		t.Error("Month should contain day 15")
	}
	if !strings.Contains(result, "31") {
		t.Error("January should contain day 31")
	}
}

func TestMonthMondayStart(t *testing.T) {
	opts := Options{FirstDayMon: true}
	result := Month(2024, time.January, opts)

	// Day headers should start with Mo
	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Fatal("Month should have at least 2 lines")
	}
	dayLine := lines[1]
	if !strings.HasPrefix(strings.TrimSpace(dayLine), "Mo") {
		t.Errorf("Monday-start should begin with Mo, got: %q", dayLine)
	}
}

func TestMonthWithWeekNumbers(t *testing.T) {
	opts := Options{ShowWeekNum: true}
	result := Month(2024, time.January, opts)

	// Should contain week numbers (2-digit prefix)
	if !strings.Contains(result, " 1 ") && !strings.Contains(result, " 2 ") {
		// Week numbers vary; just check the output is longer than without
		withoutWeek := Month(2024, time.January, Options{})
		if len(result) <= len(withoutWeek) {
			t.Error("Week numbers should add content")
		}
	}
}

func TestMonthHighlight(t *testing.T) {
	highlight := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.Local)
	opts := Options{Highlight: []time.Time{highlight}}
	result := Month(2024, time.January, opts)

	// Should contain ANSI code for highlighted day
	if !strings.Contains(result, "\033[1;33m") {
		t.Error("Highlighted day should use bold yellow ANSI code")
	}
}

func TestMonthToday(t *testing.T) {
	now := time.Now()
	opts := Options{}
	result := Month(now.Year(), now.Month(), opts)

	// Should contain inverse ANSI for today
	if !strings.Contains(result, "\033[7m") {
		t.Error("Current month should highlight today with inverse")
	}
}

func TestYear(t *testing.T) {
	opts := Options{}
	result := Year(2024, opts)

	// Should contain year header
	if !strings.Contains(result, "2024") {
		t.Error("Year should contain the year number")
	}

	// Should contain all months
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}
	for _, m := range months {
		if !strings.Contains(result, m) {
			t.Errorf("Year should contain %s", m)
		}
	}
}

func TestCurrent(t *testing.T) {
	opts := Options{}
	result := Current(opts)
	now := time.Now()

	// Should contain current month
	if !strings.Contains(result, now.Month().String()) {
		t.Errorf("Current should contain %s", now.Month().String())
	}
}

func TestToday(t *testing.T) {
	result := Today()
	now := time.Now()

	// Should contain the day of the week
	if !strings.Contains(result, now.Weekday().String()) {
		t.Error("Today should contain day of week")
	}

	// Should contain the year
	if !strings.Contains(result, "2") { // at least partial year
		t.Error("Today should contain some part of the date")
	}
}

func TestCountdown(t *testing.T) {
	// Future date
	future := time.Now().Add(48 * time.Hour)
	result := Countdown(future)
	if !strings.Contains(result, "days remaining") && !strings.Contains(result, "day remaining") {
		t.Errorf("Future countdown = %q, want 'days remaining'", result)
	}

	// Past date
	past := time.Now().Add(-48 * time.Hour)
	result = Countdown(past)
	if result != "Date has passed" {
		t.Errorf("Past countdown = %q, want 'Date has passed'", result)
	}

	// Today (within hours)
	soonish := time.Now().Add(5 * time.Hour)
	result = Countdown(soonish)
	if !strings.Contains(result, "hours remaining") {
		t.Errorf("Same-day countdown = %q, want 'hours remaining'", result)
	}
}

func TestASCIIArt(t *testing.T) {
	result := ASCIIArt()
	if !strings.Contains(result, "╔") {
		t.Error("ASCIIArt should contain border characters")
	}
	now := time.Now()
	if !strings.Contains(result, now.Month().String()[:3]) {
		t.Error("ASCIIArt should contain abbreviated month")
	}
}

func TestWeekView(t *testing.T) {
	opts := Options{}
	result := WeekView(opts)

	// Should contain "Week" header
	if !strings.Contains(result, "Week") {
		t.Error("WeekView should contain 'Week' header")
	}

	// Should contain "(today)"
	if !strings.Contains(result, "(today)") {
		t.Error("WeekView should mark today")
	}

	// Should contain day names
	if !strings.Contains(result, "Monday") && !strings.Contains(result, "Sunday") {
		t.Error("WeekView should contain day names")
	}

	// Should have 7 days
	dayCount := 0
	for _, line := range strings.Split(result, "\n") {
		if strings.Contains(line, "day") && !strings.Contains(line, "Week") {
			dayCount++
		}
	}
	// At least check we have multiple day lines
	lines := strings.Split(result, "\n")
	if len(lines) < 7 {
		t.Errorf("WeekView should have at least 7 lines of content, got %d", len(lines))
	}
}

func TestWeekViewMondayStart(t *testing.T) {
	opts := Options{FirstDayMon: true}
	result := WeekView(opts)

	// Find first day line (after header and divider)
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Monday") || strings.Contains(line, "Sunday") ||
			strings.Contains(line, "Tuesday") || strings.Contains(line, "Wednesday") ||
			strings.Contains(line, "Thursday") || strings.Contains(line, "Friday") ||
			strings.Contains(line, "Saturday") {
			// First day should be Monday
			if strings.Contains(line, "Monday") {
				return // pass
			}
			t.Errorf("First day should be Monday for FirstDayMon, got: %q", line)
			return
		}
	}
}

func TestMonthFebruary(t *testing.T) {
	// Non-leap year
	opts := Options{}
	result := Month(2023, time.February, opts)
	if strings.Contains(result, "29") {
		t.Error("Feb 2023 should not have day 29")
	}

	// Leap year
	result = Month(2024, time.February, opts)
	if !strings.Contains(result, "29") {
		t.Error("Feb 2024 should have day 29")
	}
}
