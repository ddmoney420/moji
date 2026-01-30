package chain

import (
	"fmt"
	"strconv"

	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/speech"
)

// Execute runs the pipeline on the given input
func (pl *Pipeline) Execute(input string) (string, error) {
	result := input

	for _, step := range pl.steps {
		var err error
		result, err = executeStep(result, step)
		if err != nil {
			return "", fmt.Errorf("error executing step %q: %w", step.Command, err)
		}
	}

	return result, nil
}

// executeStep executes a single pipeline step
func executeStep(input string, step *Step) (string, error) {
	switch step.Command {
	case "banner":
		return executeBanner(input, step)
	case "effect":
		return executeEffect(input, step)
	case "gradient":
		return executeGradient(input, step)
	case "border":
		return executeBorder(input, step)
	case "bubble":
		return executeBubble(input, step)
	default:
		return "", fmt.Errorf("unknown command: %q", step.Command)
	}
}

// executeBanner applies banner effect
func executeBanner(input string, step *Step) (string, error) {
	// For banner, we use the text argument if provided, otherwise use input
	text := step.Text
	if text == "" {
		text = input
	}

	// Banner doesn't have variants in the current implementation,
	// so we just apply it to the text
	// This is a placeholder for future banner functionality
	return text, nil
}

// executeEffect applies a text effect (flip, reverse, etc.)
func executeEffect(input string, step *Step) (string, error) {
	effectName := step.Variant
	if effectName == "" {
		// Try to get effect type from args
		if t, ok := step.Args["type"]; ok {
			effectName = t
		} else {
			return "", fmt.Errorf("effect requires either variant (effect:flip) or type argument (type=flip)")
		}
	}

	return effects.Apply(effectName, input), nil
}

// executeGradient applies a gradient effect
func executeGradient(input string, step *Step) (string, error) {
	gradientTheme := step.Variant
	if gradientTheme == "" {
		if t, ok := step.Args["theme"]; ok {
			gradientTheme = t
		} else {
			return "", fmt.Errorf("gradient requires either variant (gradient:rainbow) or theme argument")
		}
	}

	// Validate gradient exists
	if _, ok := gradient.Themes[gradientTheme]; !ok {
		available := make([]string, 0, len(gradient.Themes))
		for k := range gradient.Themes {
			available = append(available, k)
		}
		return "", fmt.Errorf("unknown gradient theme %q, available: %v", gradientTheme, available)
	}

	mode := "horizontal"
	if m, ok := step.Args["mode"]; ok {
		mode = m
	}

	return gradient.Apply(input, gradientTheme, mode), nil
}

// executeBorder applies a border
func executeBorder(input string, step *Step) (string, error) {
	borderStyle := step.Variant
	if borderStyle == "" {
		if s, ok := step.Args["style"]; ok {
			borderStyle = s
		} else {
			return "", fmt.Errorf("border requires either variant (border:double) or style argument")
		}
	}

	padding := 1
	if p, ok := step.Args["padding"]; ok {
		var err error
		padding, err = strconv.Atoi(p)
		if err != nil {
			return "", fmt.Errorf("invalid padding value %q: %w", p, err)
		}
	}

	return patterns.CreateBorder(input, borderStyle, padding), nil
}

// executeBubble applies a speech bubble
func executeBubble(input string, step *Step) (string, error) {
	bubbleStyle := step.Variant
	if bubbleStyle == "" {
		if s, ok := step.Args["style"]; ok {
			bubbleStyle = s
		} else {
			return "", fmt.Errorf("bubble requires either variant (bubble:round) or style argument")
		}
	}

	width := 40
	if w, ok := step.Args["width"]; ok {
		var err error
		width, err = strconv.Atoi(w)
		if err != nil {
			return "", fmt.Errorf("invalid width value %q: %w", w, err)
		}
	}

	return speech.Wrap(input, bubbleStyle, width), nil
}

// ExecuteString parses and executes a DSL string
func ExecuteString(dsl string, input string) (string, error) {
	pipeline, err := Parse(dsl)
	if err != nil {
		return "", err
	}

	return pipeline.Execute(input)
}

// ValidateGradientTheme checks if a gradient theme exists
func ValidateGradientTheme(theme string) error {
	if _, ok := gradient.Themes[theme]; !ok {
		available := make([]string, 0, len(gradient.Themes))
		for k := range gradient.Themes {
			available = append(available, k)
		}
		return fmt.Errorf("unknown gradient theme %q, available: %v", theme, available)
	}
	return nil
}

// GetAvailableGradients returns list of available gradient themes
func GetAvailableGradients() []string {
	themes := make([]string, 0, len(gradient.Themes))
	for k := range gradient.Themes {
		themes = append(themes, k)
	}
	return themes
}

// GetAvailableBorders returns list of available border styles
func GetAvailableBorders() []string {
	borders := make([]string, 0)
	borders = append(borders, "single", "double", "round", "thick", "ascii", "stars", "dots", "blocks", "shadow", "fancy")
	return borders
}

// GetAvailableBubbles returns list of available bubble styles
func GetAvailableBubbles() []string {
	bubbles := make([]string, 0)
	bubbles = append(bubbles, "round", "square", "double", "thick", "ascii", "cloud")
	return bubbles
}

// GetAvailableEffects returns list of available effects
func GetAvailableEffects() []string {
	effects := make([]string, 0)
	effects = append(effects, "flip", "reverse", "zalgo", "small", "full-width", "strike", "superscript", "subscript")
	return effects
}
