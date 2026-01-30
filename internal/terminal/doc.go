// Package terminal provides terminal capability detection and information.
//
// It detects terminal capabilities including color level (NoColor/Basic/256/TrueColor),
// terminal size, and support for Unicode, Sixel, and Kitty image protocols. Helper methods
// determine if specific features are available.
//
// Example usage:
//
//	caps := terminal.Detect()
//	if caps.TrueColorSupported {
//		// Use full RGB colors
//	}
//	width, height := caps.Size()
//	hasUnicode := caps.SupportsUnicode()
package terminal
