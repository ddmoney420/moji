package chain

import (
	"strings"
	"testing"
)

func TestExecuteEffectStep(t *testing.T) {
	pipeline, err := Parse("effect:flip")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Hello")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Flip should transform the text
	if result == "Hello" {
		t.Error("Effect:flip should transform text")
	}
}

func TestExecuteGradientStep(t *testing.T) {
	pipeline, err := Parse("gradient:rainbow")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Test")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Gradient should add ANSI codes
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Gradient should produce ANSI RGB codes")
	}
}

func TestExecuteGradientWithMode(t *testing.T) {
	pipeline, err := Parse("gradient:fire mode=vertical")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Multi-line input for mode effect
	input := "Line1\nLine2"
	result, err := pipeline.Execute(input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Gradient with mode should produce color codes")
	}
}

func TestExecuteBorderStep(t *testing.T) {
	pipeline, err := Parse("border:single")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Test")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Border should add box characters
	if !strings.Contains(result, "┌") && !strings.Contains(result, "─") {
		t.Error("Border should add border characters")
	}
}

func TestExecuteBorderWithPadding(t *testing.T) {
	pipeline, err := Parse("border:double padding=2")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Text")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "╔") {
		t.Error("Border:double should use double border chars")
	}
}

func TestExecuteBubbleStep(t *testing.T) {
	pipeline, err := Parse("bubble:round")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Hello")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "╭") {
		t.Error("Bubble should add bubble characters")
	}
}

func TestExecuteBubbleWithWidth(t *testing.T) {
	pipeline, err := Parse("bubble:square width=20")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	longText := "This is a long sentence that should wrap"
	result, err := pipeline.Execute(longText)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "┌") {
		t.Error("Bubble:square should use square bubble")
	}
}

func TestExecutePipeline(t *testing.T) {
	pipeline, err := Parse("effect:flip | gradient:neon")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Test")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should have both effects applied
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Result should have gradient colors")
	}
}

func TestExecuteComplexPipeline(t *testing.T) {
	pipeline, err := Parse("gradient:rainbow | border:single padding=1")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Complex")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should have gradient and border
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Should have gradient colors")
	}
	if !strings.Contains(result, "┌") {
		t.Error("Should have border")
	}
}

func TestExecuteInvalidGradient(t *testing.T) {
	pipeline, err := Parse("gradient:nonexistent")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = pipeline.Execute("Test")
	if err == nil {
		t.Error("Should error on invalid gradient theme")
	}
	if !strings.Contains(err.Error(), "unknown gradient") {
		t.Errorf("Error should mention unknown gradient, got: %v", err)
	}
}

func TestExecuteInvalidEffect(t *testing.T) {
	pipeline, err := Parse("effect:nonexistent")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// effects.Apply should handle invalid effect gracefully
	result, err := pipeline.Execute("Test")
	if err != nil {
		t.Logf("Execute returned error (expected): %v", err)
	} else if result == "" {
		t.Error("Result should not be empty")
	}
}

func TestExecuteStringHelper(t *testing.T) {
	result, err := ExecuteString("gradient:fire", "Hello")
	if err != nil {
		t.Fatalf("ExecuteString failed: %v", err)
	}

	if !strings.Contains(result, "\033[38;2;") {
		t.Error("ExecuteString should apply gradient")
	}
}

func TestValidateGradientTheme(t *testing.T) {
	tests := []struct {
		theme    string
		valid    bool
		errorMsg string
	}{
		{"rainbow", true, ""},
		{"fire", true, ""},
		{"neon", true, ""},
		{"invalid", false, "unknown gradient theme"},
	}

	for _, test := range tests {
		err := ValidateGradientTheme(test.theme)
		if test.valid && err != nil {
			t.Errorf("ValidateGradientTheme(%q) should not error: %v", test.theme, err)
		}
		if !test.valid && err == nil {
			t.Errorf("ValidateGradientTheme(%q) should error", test.theme)
		}
		if !test.valid && err != nil && !strings.Contains(err.Error(), test.errorMsg) {
			t.Errorf("Expected error containing %q, got: %v", test.errorMsg, err)
		}
	}
}

func TestGetAvailableGradients(t *testing.T) {
	gradients := GetAvailableGradients()
	if len(gradients) == 0 {
		t.Error("Should return at least one gradient")
	}

	hasRainbow := false
	for _, g := range gradients {
		if g == "rainbow" {
			hasRainbow = true
			break
		}
	}
	if !hasRainbow {
		t.Error("Should include rainbow gradient")
	}
}

func TestGetAvailableBorders(t *testing.T) {
	borders := GetAvailableBorders()
	if len(borders) == 0 {
		t.Error("Should return at least one border style")
	}

	hasSingle := false
	for _, b := range borders {
		if b == "single" {
			hasSingle = true
			break
		}
	}
	if !hasSingle {
		t.Error("Should include single border")
	}
}

func TestGetAvailableBubbles(t *testing.T) {
	bubbles := GetAvailableBubbles()
	if len(bubbles) == 0 {
		t.Error("Should return at least one bubble style")
	}

	hasRound := false
	for _, b := range bubbles {
		if b == "round" {
			hasRound = true
			break
		}
	}
	if !hasRound {
		t.Error("Should include round bubble")
	}
}

func TestGetAvailableEffects(t *testing.T) {
	effects := GetAvailableEffects()
	if len(effects) == 0 {
		t.Error("Should return at least one effect")
	}

	hasFlip := false
	for _, e := range effects {
		if e == "flip" {
			hasFlip = true
			break
		}
	}
	if !hasFlip {
		t.Error("Should include flip effect")
	}
}

func TestExecuteEffectWithoutVariant(t *testing.T) {
	// Effect requires variant or type argument
	pipeline, err := Parse("effect type=flip")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Hello")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result == "Hello" {
		t.Error("Effect should transform text")
	}
}

func TestExecuteGradientWithoutVariant(t *testing.T) {
	// Gradient requires variant or theme argument
	pipeline, err := Parse("gradient theme=neon")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Test")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Gradient should produce colors")
	}
}

func TestExecuteBorderWithoutVariant(t *testing.T) {
	// Border requires variant or style argument
	pipeline, err := Parse("border style=round")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Text")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "╭") {
		t.Error("Border should add border characters")
	}
}

func TestExecuteBubbleWithoutVariant(t *testing.T) {
	// Bubble requires variant or style argument
	pipeline, err := Parse("bubble style=double")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	result, err := pipeline.Execute("Hello")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "╔") {
		t.Error("Bubble:double should use double chars")
	}
}

func TestExecuteInvalidPaddingValue(t *testing.T) {
	pipeline, err := Parse("border:single padding=abc")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = pipeline.Execute("Test")
	if err == nil {
		t.Error("Should error on invalid padding value")
	}
	if !strings.Contains(err.Error(), "invalid padding") {
		t.Errorf("Error should mention invalid padding, got: %v", err)
	}
}

func TestExecuteInvalidWidthValue(t *testing.T) {
	pipeline, err := Parse("bubble:round width=notanumber")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = pipeline.Execute("Test")
	if err == nil {
		t.Error("Should error on invalid width value")
	}
	if !strings.Contains(err.Error(), "invalid width") {
		t.Errorf("Error should mention invalid width, got: %v", err)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	pipeline, err := Parse("unknown_cmd")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = pipeline.Execute("Test")
	if err == nil {
		t.Error("Should error on unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("Error should mention unknown command, got: %v", err)
	}
}
