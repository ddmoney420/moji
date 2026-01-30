// Package speech provides speech bubble text formatting in various styles.
//
// It supports 6 bubble styles (round, square, double, thick, ascii, think) with automatic
// word-wrapping. Bubbles can be combined with ASCII art for speech display.
//
// Example usage:
//
//	bubble := speech.Wrap("Hello!", speech.RoundBubble)
//	bubble := speech.Wrap("Thinking...", speech.ThinkBubble)
//	combined := speech.Combine(asciiArt, bubble)
package speech
