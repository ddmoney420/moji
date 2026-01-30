// Package qrcode provides QR code generation to ASCII art format.
//
// It generates QR codes as ASCII art with support for 9 character sets and ANSI coloring.
// Compact rendering uses half-block characters for higher resolution.
//
// Example usage:
//
//	ascii := qrcode.Generate("https://example.com")
//	compact := qrcode.GenerateCompact("data")
//	colored := qrcode.GenerateColored("text", colorScheme)
package qrcode
