// Package banner provides FIGlet font integration for generating ASCII text banners.
//
// It manages embedded FIGlet fonts with caching support, allowing users to render text in various
// decorative ASCII art styles. The package provides access to 40+ fonts including 3D, retro/gaming,
// graffiti, and other artistic styles.
//
// Example usage:
//
//	font := banner.GetFont("shadow")
//	result := banner.Generate("HELLO", font)
//	fonts := banner.ListFonts()
//	banner.SetCustomFontPath("/path/to/fonts")
package banner
