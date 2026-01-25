package chain

import (
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	// No options - returns unchanged
	result := Apply("Hello", Options{})
	if result != "Hello" {
		t.Errorf("Apply with empty opts should return text unchanged, got %q", result)
	}
}

func TestApplyEffect(t *testing.T) {
	result := Apply("Hello", Options{Effect: "reverse"})
	if result == "Hello" {
		t.Error("Apply with Effect should transform text")
	}
}

func TestApplyBubble(t *testing.T) {
	result := Apply("Hello", Options{Bubble: "round"})
	if !strings.Contains(result, "╭") {
		t.Error("Apply with Bubble should add speech bubble")
	}
	if !strings.Contains(result, "Hello") {
		t.Error("Bubble should contain original text")
	}
}

func TestApplyBubbleWidth(t *testing.T) {
	longText := "This is a very long sentence that should be wrapped according to BubbleWidth"
	result := Apply(longText, Options{Bubble: "round", BubbleWidth: 20})
	// Should wrap into multiple lines
	lines := strings.Split(result, "\n")
	if len(lines) < 4 {
		t.Error("BubbleWidth should cause text wrapping")
	}
}

func TestApplyBubbleDefaultWidth(t *testing.T) {
	result := Apply("Short", Options{Bubble: "square", BubbleWidth: 0})
	// Should use default width of 40
	if !strings.Contains(result, "┌") {
		t.Error("Should apply square bubble with default width")
	}
}

func TestApplyBorder(t *testing.T) {
	result := Apply("Hello", Options{Border: "double"})
	if !strings.Contains(result, "═") {
		t.Error("Apply with Border should add double border")
	}
}

func TestApplyBorderPad(t *testing.T) {
	result := Apply("Hi", Options{Border: "single", BorderPad: 2})
	// Should have padding inside border
	if !strings.Contains(result, "Hi") {
		t.Error("Border should contain text")
	}
}

func TestApplyGradient(t *testing.T) {
	result := Apply("Hello\nWorld", Options{Gradient: "rainbow"})
	// Gradient produces ANSI color codes
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Apply with Gradient should produce ANSI RGB codes")
	}
}

func TestApplyGradientMode(t *testing.T) {
	result := Apply("Hello\nWorld", Options{Gradient: "fire", GradientMode: "vertical"})
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Gradient with vertical mode should produce color codes")
	}
}

func TestApplyGradientDefaultMode(t *testing.T) {
	result := Apply("Test", Options{Gradient: "neon", GradientMode: ""})
	// Should default to "horizontal"
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Gradient with default mode should produce color codes")
	}
}

func TestApplyChainOrder(t *testing.T) {
	// Effect is applied first, then bubble, then border, then gradient
	result := Apply("Hello", Options{
		Effect:  "reverse",
		Bubble:  "round",
		Border:  "single",
		Gradient: "rainbow",
	})

	// Should have gradient colors (applied last)
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Chained result should have gradient colors")
	}
	// Should have border
	if !strings.Contains(result, "─") {
		t.Error("Chained result should have border")
	}
	// Should have bubble
	if !strings.Contains(result, "╭") || !strings.Contains(result, "╰") {
		t.Error("Chained result should have bubble")
	}
}

func TestApplyToArt(t *testing.T) {
	art := " /\\_/\\\n( o.o )\n > ^ <"
	result := ApplyToArt(art, Options{})
	if result != art {
		t.Error("ApplyToArt with empty opts should return unchanged")
	}
}

func TestApplyToArtBorder(t *testing.T) {
	art := "Cat\nArt"
	result := ApplyToArt(art, Options{Border: "round"})
	if !strings.Contains(result, "Cat") {
		t.Error("ApplyToArt should preserve art content")
	}
}

func TestApplyToArtGradient(t *testing.T) {
	art := "Hello\nWorld\nTest"
	result := ApplyToArt(art, Options{Gradient: "fire"})
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("ApplyToArt with gradient should produce color codes")
	}
}

func TestApplyToArtGradientDefaultMode(t *testing.T) {
	art := "Art\nPiece"
	result := ApplyToArt(art, Options{Gradient: "rainbow", GradientMode: ""})
	// Default mode for art is "diagonal"
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("ApplyToArt gradient should default to diagonal mode")
	}
}

func TestApplyToArtBorderAndGradient(t *testing.T) {
	art := "Star"
	result := ApplyToArt(art, Options{
		Border:   "single",
		Gradient: "neon",
	})
	// Should have both border and gradient
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Should have gradient colors")
	}
}
