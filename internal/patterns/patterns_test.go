package patterns

import (
	"strings"
	"testing"
)

func TestCreateBorderSingle(t *testing.T) {
	result := CreateBorder("Hello", "single", 1)
	if result == "" {
		t.Fatal("CreateBorder() returned empty")
	}
	if !strings.Contains(result, "Hello") {
		t.Error("CreateBorder() should contain the text")
	}
	if !strings.Contains(result, "┌") {
		t.Error("single border should contain ┌")
	}
	if !strings.Contains(result, "┘") {
		t.Error("single border should contain ┘")
	}
}

func TestCreateBorderDouble(t *testing.T) {
	result := CreateBorder("Test", "double", 0)
	if !strings.Contains(result, "╔") {
		t.Error("double border should contain ╔")
	}
}

func TestCreateBorderUnknown(t *testing.T) {
	result := CreateBorder("Test", "nonexistent_style", 1)
	// Should fall back to single
	if result == "" {
		t.Fatal("unknown style should fallback, not return empty")
	}
}

func TestCreateBorderMultiline(t *testing.T) {
	result := CreateBorder("Line1\nLine2\nLine3", "round", 0)
	if !strings.Contains(result, "Line1") {
		t.Error("multiline border should contain all lines")
	}
	if !strings.Contains(result, "Line3") {
		t.Error("multiline border should contain all lines")
	}
}

func TestCreateBorderAllStyles(t *testing.T) {
	for name := range Borders {
		result := CreateBorder("X", name, 0)
		if result == "" {
			t.Errorf("CreateBorder with style %q returned empty", name)
		}
	}
}

func TestCreateDivider(t *testing.T) {
	result := CreateDivider("single", 20)
	if result == "" {
		t.Fatal("CreateDivider() returned empty")
	}
	if len([]rune(result)) < 20 {
		t.Errorf("CreateDivider() length = %d, want >= 20", len([]rune(result)))
	}
}

func TestCreateDividerAllStyles(t *testing.T) {
	for name := range Dividers {
		result := CreateDivider(name, 10)
		if result == "" {
			t.Errorf("CreateDivider with style %q returned empty", name)
		}
	}
}

func TestCreateDividerUnknown(t *testing.T) {
	result := CreateDivider("nonexistent", 10)
	if result == "" {
		t.Fatal("unknown divider style should fallback")
	}
}

func TestListBorders(t *testing.T) {
	borders := ListBorders()
	if len(borders) == 0 {
		t.Fatal("ListBorders() returned empty")
	}
	if len(borders) != len(Borders) {
		t.Errorf("ListBorders() = %d, Borders map has %d", len(borders), len(Borders))
	}
}

func TestListDividers(t *testing.T) {
	dividers := ListDividers()
	if len(dividers) == 0 {
		t.Fatal("ListDividers() returned empty")
	}
	if len(dividers) != len(Dividers) {
		t.Errorf("ListDividers() = %d, Dividers map has %d", len(dividers), len(Dividers))
	}
}

func TestCreateBorderPadding(t *testing.T) {
	noPad := CreateBorder("X", "single", 0)
	withPad := CreateBorder("X", "single", 2)
	// With padding, lines should be wider
	noPadLines := strings.Split(noPad, "\n")
	withPadLines := strings.Split(withPad, "\n")
	if len([]rune(withPadLines[1])) <= len([]rune(noPadLines[1])) {
		t.Error("padding should make border wider")
	}
}
