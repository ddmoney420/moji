package chain

import (
	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/speech"
)

// Options for chaining effects
type Options struct {
	Gradient     string // Gradient theme name (empty = none)
	GradientMode string // horizontal, vertical, diagonal
	Border       string // Border style (empty = none)
	BorderPad    int    // Padding inside border
	Effect       string // Text effect (flip, reverse, etc.)
	Bubble       string // Speech bubble style (empty = none)
	BubbleWidth  int    // Max bubble width
}

// Apply applies a chain of effects to text
func Apply(text string, opts Options) string {
	result := text

	// 1. Apply text effect first (flip, reverse, etc.)
	if opts.Effect != "" {
		result = effects.Apply(opts.Effect, result)
	}

	// 2. Apply speech bubble
	if opts.Bubble != "" {
		width := opts.BubbleWidth
		if width <= 0 {
			width = 40
		}
		result = speech.Wrap(result, opts.Bubble, width)
	}

	// 3. Apply border
	if opts.Border != "" {
		padding := opts.BorderPad
		if padding <= 0 {
			padding = 1
		}
		result = patterns.CreateBorder(result, opts.Border, padding)
	}

	// 4. Apply gradient last (needs to be on final output)
	if opts.Gradient != "" {
		mode := opts.GradientMode
		if mode == "" {
			mode = "horizontal"
		}
		result = gradient.Apply(result, opts.Gradient, mode)
	}

	return result
}

// ApplyToArt applies effects suitable for multi-line ASCII art
func ApplyToArt(art string, opts Options) string {
	result := art

	// 1. Apply border
	if opts.Border != "" {
		padding := opts.BorderPad
		if padding <= 0 {
			padding = 1
		}
		result = patterns.CreateBorder(result, opts.Border, padding)
	}

	// 2. Apply gradient
	if opts.Gradient != "" {
		mode := opts.GradientMode
		if mode == "" {
			mode = "diagonal"
		}
		result = gradient.Apply(result, opts.Gradient, mode)
	}

	return result
}
