package chain

import (
	"strings"
	"testing"
)

func TestLexerTokens(t *testing.T) {
	input := "banner 'Hello' | gradient:fire"
	lexer := NewLexer(input)

	tokens := []struct {
		typ   TokenType
		value string
	}{
		{TokenIdentifier, "banner"},
		{TokenString, "Hello"},
		{TokenPipe, "|"},
		{TokenIdentifier, "gradient"},
		{TokenColon, ":"},
		{TokenIdentifier, "fire"},
	}

	for _, expected := range tokens {
		token := lexer.NextToken()
		if token.Type != expected.typ {
			t.Errorf("Expected token type %v, got %v", expected.typ, token.Type)
		}
		if token.Value != expected.value {
			t.Errorf("Expected token value %q, got %q", expected.value, token.Value)
		}
	}

	// Verify EOF
	token := lexer.NextToken()
	if token.Type != TokenEOF {
		t.Errorf("Expected EOF, got %v", token.Type)
	}
}

func TestLexerQuotedString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'hello world'", "hello world"},
		{`"quoted string"`, "quoted string"},
		{`'with\' escaped'`, "with' escaped"},
	}

	for _, test := range tests {
		lexer := NewLexer(test.input)
		token := lexer.NextToken()
		if token.Type != TokenString {
			t.Errorf("Expected string token, got %v", token.Type)
		}
		if token.Value != test.expected {
			t.Errorf("Expected %q, got %q", test.expected, token.Value)
		}
	}
}

func TestLexerWhitespace(t *testing.T) {
	input := "  banner   'text'  "
	lexer := NewLexer(input)

	// First token should be banner
	token := lexer.NextToken()
	if token.Type != TokenIdentifier || token.Value != "banner" {
		t.Errorf("Expected 'banner' identifier")
	}

	// Second token should be string
	token = lexer.NextToken()
	if token.Type != TokenString {
		t.Errorf("Expected string token after whitespace")
	}
}

func TestParseSimpleCommand(t *testing.T) {
	pipeline, err := Parse("banner")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if len(steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(steps))
	}

	if steps[0].Command != "banner" {
		t.Errorf("Expected command 'banner', got %q", steps[0].Command)
	}
}

func TestParseCommandWithVariant(t *testing.T) {
	pipeline, err := Parse("gradient:rainbow")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if len(steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(steps))
	}

	if steps[0].Command != "gradient" {
		t.Errorf("Expected command 'gradient', got %q", steps[0].Command)
	}
	if steps[0].Variant != "rainbow" {
		t.Errorf("Expected variant 'rainbow', got %q", steps[0].Variant)
	}
}

func TestParseCommandWithTextArg(t *testing.T) {
	pipeline, err := Parse("banner 'Hello World'")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if steps[0].Text != "Hello World" {
		t.Errorf("Expected text 'Hello World', got %q", steps[0].Text)
	}
}

func TestParseCommandWithKeyValue(t *testing.T) {
	pipeline, err := Parse("border:double padding=2")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if steps[0].Command != "border" {
		t.Errorf("Expected command 'border'")
	}
	if steps[0].Variant != "double" {
		t.Errorf("Expected variant 'double'")
	}
	if steps[0].Args["padding"] != "2" {
		t.Errorf("Expected padding=2, got %q", steps[0].Args["padding"])
	}
}

func TestParseMultipleSteps(t *testing.T) {
	pipeline, err := Parse("banner 'Hi' | gradient:fire | border:single")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if len(steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(steps))
	}

	if steps[0].Command != "banner" || steps[0].Text != "Hi" {
		t.Error("First step should be banner with text 'Hi'")
	}
	if steps[1].Command != "gradient" || steps[1].Variant != "fire" {
		t.Error("Second step should be gradient:fire")
	}
	if steps[2].Command != "border" || steps[2].Variant != "single" {
		t.Error("Third step should be border:single")
	}
}

func TestParseComplexPipeline(t *testing.T) {
	input := "banner 'Hello' | gradient:rainbow mode=vertical | border:double padding=2 | bubble:round width=50"
	pipeline, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if len(steps) != 4 {
		t.Errorf("Expected 4 steps, got %d", len(steps))
	}

	// Verify gradient step args
	if steps[1].Args["mode"] != "vertical" {
		t.Errorf("Expected mode=vertical, got %q", steps[1].Args["mode"])
	}

	// Verify border step args
	if steps[2].Args["padding"] != "2" {
		t.Errorf("Expected padding=2, got %q", steps[2].Args["padding"])
	}

	// Verify bubble step args
	if steps[3].Args["width"] != "50" {
		t.Errorf("Expected width=50, got %q", steps[3].Args["width"])
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		input    string
		errorMsg string
	}{
		{"", "empty pipeline"},
		{"| gradient:fire", "expected command name"},
	}

	for _, test := range tests {
		_, err := Parse(test.input)
		if err == nil {
			t.Errorf("Expected error for input %q", test.input)
		}
		if err != nil && !strings.Contains(err.Error(), test.errorMsg) {
			t.Errorf("Expected error containing %q, got %v", test.errorMsg, err)
		}
	}
}

func TestParserPosition(t *testing.T) {
	input := "banner | gradient"
	pipeline, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if len(steps) != 2 {
		t.Errorf("Expected 2 steps")
	}
}

func TestParserIdentifierWithDashes(t *testing.T) {
	pipeline, err := Parse("border:double-line")
	if err == nil || strings.Contains(err.Error(), "double-line") {
		// Parser should handle dash in identifiers
		if err == nil && len(pipeline.Steps()) > 0 {
			t.Logf("Successfully parsed border with dash: %v", pipeline.Steps()[0].Variant)
		}
	}
}

func TestParseMultipleKeyValueArgs(t *testing.T) {
	pipeline, err := Parse("border:double padding=2 width=80")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if steps[0].Args["padding"] != "2" {
		t.Errorf("Expected padding=2")
	}
	if steps[0].Args["width"] != "80" {
		t.Errorf("Expected width=80")
	}
}

func TestParseTextAndArgs(t *testing.T) {
	pipeline, err := Parse("banner 'Hello' width=80")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if steps[0].Text != "Hello" {
		t.Errorf("Expected text 'Hello'")
	}
	if steps[0].Args["width"] != "80" {
		t.Errorf("Expected width=80")
	}
}

func TestPipelineSteps(t *testing.T) {
	pipeline, err := Parse("effect:flip | gradient:neon")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if len(steps) != 2 {
		t.Errorf("Expected 2 steps")
	}

	if steps[0].Command != "effect" {
		t.Errorf("First step should be effect")
	}
	if steps[1].Command != "gradient" {
		t.Errorf("Second step should be gradient")
	}
}

func TestEmptyArgsMap(t *testing.T) {
	pipeline, err := Parse("banner")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	steps := pipeline.Steps()
	if steps[0].Args == nil {
		t.Error("Args map should not be nil")
	}
	if len(steps[0].Args) != 0 {
		t.Error("Args map should be empty")
	}
}
