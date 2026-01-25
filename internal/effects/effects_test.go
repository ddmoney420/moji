package effects

import (
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	tests := []struct {
		effect string
		input  string
	}{
		{"flip", "Hello"},
		{"reverse", "Hello"},
		{"mirror", "Hello"},
		{"bubble", "Hello"},
		{"square", "Hello"},
		{"bold", "Hello"},
		{"italic", "Hello"},
		{"strikethrough", "Hello"},
		{"underline", "Hello"},
		{"smallcaps", "Hello"},
		{"fullwidth", "Hello"},
		{"monospace", "Hello"},
		{"script", "Hello"},
		{"fraktur", "Hello"},
		{"double-struck", "Hello"},
		{"sparkle", "Hello"},
	}

	for _, tt := range tests {
		result := Apply(tt.effect, tt.input)
		if result == "" {
			t.Errorf("Apply(%q, %q) returned empty", tt.effect, tt.input)
		}
	}
}

func TestApplyUnknown(t *testing.T) {
	result := Apply("nonexistent_effect", "Hello")
	if result != "Hello" {
		t.Errorf("Apply with unknown effect should return input unchanged, got %q", result)
	}
}

func TestFlip(t *testing.T) {
	result := Flip("Hello")
	if result == "Hello" {
		t.Error("Flip() should transform text")
	}
	if len([]rune(result)) != len([]rune("Hello")) {
		t.Error("Flip() should preserve rune count")
	}
}

func TestReverse(t *testing.T) {
	result := Reverse("Hello")
	if result != "olleH" {
		t.Errorf("Reverse('Hello') = %q, want 'olleH'", result)
	}
}

func TestReverseEmpty(t *testing.T) {
	result := Reverse("")
	if result != "" {
		t.Errorf("Reverse('') = %q, want ''", result)
	}
}

func TestMirror(t *testing.T) {
	result := Mirror("Hi")
	if !strings.Contains(result, "Hi") {
		t.Error("Mirror() should contain original text")
	}
}

func TestBubble(t *testing.T) {
	result := Bubble("abc")
	// Bubble converts to circled letters
	if result == "abc" {
		t.Error("Bubble() should transform characters")
	}
}

func TestSquare(t *testing.T) {
	result := Square("abc")
	if result == "abc" {
		t.Error("Square() should transform characters")
	}
}

func TestFullwidth(t *testing.T) {
	result := Fullwidth("abc")
	if result == "abc" {
		t.Error("Fullwidth() should transform characters")
	}
	// Fullwidth chars are wider
	if len(result) <= len("abc") {
		t.Error("Fullwidth() should produce wider output")
	}
}

func TestZalgo(t *testing.T) {
	result := Zalgo("Hi", 2)
	if len(result) <= len("Hi") {
		t.Error("Zalgo() should add combining characters")
	}
}

func TestZalgoIntensity(t *testing.T) {
	mild := Zalgo("Hello", 1)
	intense := Zalgo("Hello", 6)
	if len(intense) <= len(mild) {
		t.Error("higher intensity zalgo should produce more characters")
	}
}

func TestListEffects(t *testing.T) {
	effects := ListEffects()
	if len(effects) == 0 {
		t.Fatal("ListEffects() returned empty")
	}
	for _, e := range effects {
		if e.Name == "" {
			t.Error("effect has empty Name")
		}
		if e.Desc == "" {
			t.Error("effect has empty Desc")
		}
	}
}

func TestWave(t *testing.T) {
	result := Wave("Hello World")
	if result == "" {
		t.Fatal("Wave() returned empty")
	}
}

func TestSparkle(t *testing.T) {
	result := Sparkle("Hi")
	if !strings.Contains(result, "Hi") {
		t.Error("Sparkle() should contain original text")
	}
	if result == "Hi" {
		t.Error("Sparkle() should add sparkle characters")
	}
}
