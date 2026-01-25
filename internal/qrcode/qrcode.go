package qrcode

import (
	"strings"

	qr "github.com/skip2/go-qrcode"
)

// CharSets for QR code rendering
var CharSets = map[string]struct {
	Full  string
	Empty string
}{
	"blocks":  {"██", "  "},
	"shaded":  {"▓▓", "░░"},
	"dots":    {"●●", "  "},
	"ascii":   {"##", "  "},
	"braille": {"⣿⣿", "⠀⠀"},
	"compact": {"█", " "},
	"inverse": {"  ", "██"},
	"minimal": {"▀", " "},
	"half":    {"▄▄", "  "},
}

// Options for QR code generation
type Options struct {
	Size    int
	Charset string
	Invert  bool
	ANSI    bool
}

// Generate creates an ASCII QR code from text
func Generate(text string, opts Options) (string, error) {
	// Create QR code
	code, err := qr.New(text, qr.Medium)
	if err != nil {
		return "", err
	}

	// Get bitmap
	bitmap := code.Bitmap()

	// Get charset
	cs, ok := CharSets[opts.Charset]
	if !ok {
		cs = CharSets["blocks"]
	}

	full := cs.Full
	empty := cs.Empty

	if opts.Invert {
		full, empty = empty, full
	}

	// Build ASCII representation
	var sb strings.Builder

	if opts.ANSI {
		// Use ANSI background colors for reliable contrast on any terminal
		darkBg := "\033[40m"  // black background for dark modules
		lightBg := "\033[47m" // white background for light modules
		reset := "\033[0m"

		if opts.Invert {
			darkBg, lightBg = lightBg, darkBg
		}

		for _, row := range bitmap {
			for _, cell := range row {
				if cell {
					sb.WriteString(darkBg + "  " + reset)
				} else {
					sb.WriteString(lightBg + "  " + reset)
				}
			}
			sb.WriteString("\n")
		}
	} else {
		for _, row := range bitmap {
			for _, cell := range row {
				if cell {
					sb.WriteString(full)
				} else {
					sb.WriteString(empty)
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

// GenerateCompact creates a compact QR using half-block characters
func GenerateCompact(text string, invert bool) (string, error) {
	code, err := qr.New(text, qr.Medium)
	if err != nil {
		return "", err
	}

	bitmap := code.Bitmap()
	height := len(bitmap)

	var sb strings.Builder

	// Process two rows at a time using half-blocks
	for y := 0; y < height; y += 2 {
		for x := 0; x < len(bitmap[y]); x++ {
			top := bitmap[y][x]
			bottom := false
			if y+1 < height {
				bottom = bitmap[y+1][x]
			}

			if invert {
				top = !top
				bottom = !bottom
			}

			// Use half-block characters
			switch {
			case top && bottom:
				sb.WriteString("█")
			case top && !bottom:
				sb.WriteString("▀")
			case !top && bottom:
				sb.WriteString("▄")
			default:
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// ListCharsets returns available charset names
func ListCharsets() []string {
	return []string{"blocks", "shaded", "dots", "ascii", "braille", "compact", "inverse", "minimal", "half"}
}
