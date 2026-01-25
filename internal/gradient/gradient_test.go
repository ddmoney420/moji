package gradient

import (
	"testing"
)

func TestApply(t *testing.T) {
	result := Apply("Hello", "rainbow", "horizontal")
	if result == "" {
		t.Fatal("Apply() returned empty")
	}
	if len(result) <= len("Hello") {
		t.Error("Apply() should add ANSI color codes")
	}
}

func TestApplyAllThemes(t *testing.T) {
	themes := ListThemes()
	for _, theme := range themes {
		result := Apply("Test", theme.Name, "horizontal")
		if result == "" {
			t.Errorf("Apply with theme %q returned empty", theme.Name)
		}
	}
}

func TestApplyModes(t *testing.T) {
	modes := []string{"horizontal", "vertical", "diagonal"}
	for _, mode := range modes {
		result := Apply("Hello\nWorld", "rainbow", mode)
		if result == "" {
			t.Errorf("Apply with mode %q returned empty", mode)
		}
	}
}

func TestApplyPerLine(t *testing.T) {
	result := ApplyPerLine("Line1\nLine2\nLine3", "neon")
	if result == "" {
		t.Fatal("ApplyPerLine() returned empty")
	}
}

func TestApplyUnknownTheme(t *testing.T) {
	result := Apply("Hello", "nonexistent_theme", "horizontal")
	// Should still return something (fallback or unchanged)
	if result == "" {
		t.Error("Apply with unknown theme should not return empty")
	}
}

func TestListThemes(t *testing.T) {
	themes := ListThemes()
	if len(themes) == 0 {
		t.Fatal("ListThemes() returned empty")
	}
	for _, theme := range themes {
		if theme.Name == "" {
			t.Error("theme has empty Name")
		}
	}
}
