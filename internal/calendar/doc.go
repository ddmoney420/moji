// Package calendar provides ASCII calendar generation with multiple view modes.
//
// It generates calendars in various formats including monthly, yearly, current month, and week views.
// The package supports customization options for highlighting specific dates, formatting preferences,
// and ASCII art variations.
//
// Example usage:
//
//	monthStr := calendar.Month(time.Now())
//	yearStr := calendar.Year(2024)
//	currentStr := calendar.Current()
//	weekStr := calendar.WeekView(time.Now())
//	countdown := calendar.Countdown(targetDate)
package calendar
