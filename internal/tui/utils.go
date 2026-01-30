// utils.go contains utility helper functions used across the TUI package,
// including filename generation, ANSI stripping, and file extension handling.
package tui

import (
	"fmt"
	"path/filepath"
	"time"
)

// makeRange creates a slice of consecutive integers
func makeRange(n int) []int {
	r := make([]int, n)
	for i := range r {
		r[i] = i
	}
	return r
}

// generateFilename generates a timestamped filename
func generateFilename(ext string) string {
	return fmt.Sprintf("moji_%s.%s", time.Now().Format("20060102_150405"), ext)
}

// changeExtension changes the file extension
func changeExtension(filename, newExt string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return filename[:len(filename)-len(ext)] + "." + newExt
	}
	return filename + "." + newExt
}

// stripANSI removes ANSI escape sequences from a string
func stripANSI(s string) string {
	var result []rune
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result = append(result, r)
	}
	return string(result)
}
