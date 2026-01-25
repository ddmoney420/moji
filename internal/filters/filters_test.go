package filters

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	for name := range registry {
		f, ok := Get(name)
		if !ok {
			t.Errorf("Get(%q) returned not ok", name)
		}
		if f == nil {
			t.Errorf("Get(%q) returned nil filter", name)
		}
	}
}

func TestGetNotFound(t *testing.T) {
	_, ok := Get("nonexistent_filter_xyz")
	if ok {
		t.Error("Get() should return false for nonexistent filter")
	}
}

func TestGetCaseInsensitive(t *testing.T) {
	_, ok := Get("RAINBOW")
	if !ok {
		t.Error("Get() should be case-insensitive")
	}
	_, ok = Get("Metal")
	if !ok {
		t.Error("Get() should be case-insensitive for Metal")
	}
}

func TestAllFiltersProduceOutput(t *testing.T) {
	input := "Hello World\nSecond Line"
	for name, f := range registry {
		result := f(input)
		if result == "" {
			t.Errorf("filter %q produced empty output for %q", name, input)
		}
	}
}

func TestRainbow(t *testing.T) {
	result := Rainbow("Hello")
	if result == "" {
		t.Fatal("Rainbow() returned empty")
	}
	if len(result) <= len("Hello") {
		t.Error("Rainbow() should add ANSI color codes")
	}
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Rainbow should use RGB color codes")
	}
}

func TestRainbowPreservesSpaces(t *testing.T) {
	result := Rainbow("A B")
	// The output should contain space characters without ANSI wrapping
	if !strings.Contains(result, "A") || !strings.Contains(result, "B") {
		t.Error("Rainbow should preserve characters")
	}
}

func TestMetal(t *testing.T) {
	result := Metal("Hello\nWorld")
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Metal should use RGB color codes")
	}
	if !strings.Contains(result, "H") {
		t.Error("Metal should preserve characters")
	}
}

func TestMetalPreservesSpaces(t *testing.T) {
	result := Metal("A B")
	if !strings.Contains(result, " ") {
		t.Error("Metal should preserve spaces without coloring them")
	}
}

func TestCrop(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"removes leading empty lines", "\n\nHello", "Hello"},
		{"removes trailing empty lines", "Hello\n\n", "Hello"},
		{"removes common indent", "  A\n  B", "A\nB"},
		{"preserves relative indent", "  A\n    B", "A\n  B"},
		{"all empty returns empty", "\n\n\n", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Crop(tt.input)
			if result != tt.expect {
				t.Errorf("Crop(%q) = %q, want %q", tt.input, result, tt.expect)
			}
		})
	}
}

func TestFlip(t *testing.T) {
	result := Flip("Hello\nWorld")
	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Fatalf("Flip() should preserve line count, got %d", len(lines))
	}
	if lines[0] != "World" {
		t.Errorf("Flip() first line = %q, want %q", lines[0], "World")
	}
	if lines[1] != "Hello" {
		t.Errorf("Flip() second line = %q, want %q", lines[1], "Hello")
	}
}

func TestFlop(t *testing.T) {
	result := Flop("AB\nCD")
	lines := strings.Split(result, "\n")
	if lines[0] != "BA" {
		t.Errorf("Flop() first line = %q, want 'BA'", lines[0])
	}
	if lines[1] != "DC" {
		t.Errorf("Flop() second line = %q, want 'DC'", lines[1])
	}
}

func TestRotate180(t *testing.T) {
	result := Rotate180("AB\nCD")
	lines := strings.Split(result, "\n")
	if lines[0] != "DC" {
		t.Errorf("Rotate180 first line = %q, want 'DC'", lines[0])
	}
	if lines[1] != "BA" {
		t.Errorf("Rotate180 second line = %q, want 'BA'", lines[1])
	}
}

func TestBorder(t *testing.T) {
	result := Border("Test")
	if !strings.Contains(result, "Test") {
		t.Error("Border() should contain original text")
	}
	if !strings.Contains(result, "┌") {
		t.Error("Border() should contain top-left corner")
	}
	if !strings.Contains(result, "┘") {
		t.Error("Border() should contain bottom-right corner")
	}
	if !strings.Contains(result, "│ Test │") {
		t.Error("Border() should have padded content line")
	}
}

func TestBorderMultiLine(t *testing.T) {
	result := Border("AB\nCDE")
	if !strings.Contains(result, "│ AB") {
		t.Error("Border should contain first line")
	}
	if !strings.Contains(result, "│ CDE") {
		t.Error("Border should contain second line")
	}
}

func TestShadow(t *testing.T) {
	result := Shadow("Hi")
	if !strings.Contains(result, "Hi") {
		t.Error("Shadow should contain original text")
	}
	if !strings.Contains(result, "\033[90m") {
		t.Error("Shadow should use dark gray")
	}
	if !strings.Contains(result, "░") {
		t.Error("Shadow should use ░ character")
	}
}

func TestShadow3D(t *testing.T) {
	result := Shadow3D("Hi")
	if !strings.Contains(result, "Hi") {
		t.Error("Shadow3D should contain original text")
	}
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Shadow3D should use RGB colors for layers")
	}
}

func TestGlitch(t *testing.T) {
	result := Glitch("Hello World Test String Here")
	if result == "" {
		t.Error("Glitch should produce output")
	}
	// Due to randomness, just check it produces something
}

func TestMatrix(t *testing.T) {
	result := Matrix("Hello")
	if !strings.Contains(result, "\033[38;2;0;") {
		t.Error("Matrix should use green RGB colors")
	}
}

func TestFire(t *testing.T) {
	result := Fire("Hello\nWorld")
	if !strings.Contains(result, "\033[38;2;255;") {
		t.Error("Fire should use red-based RGB colors")
	}
}

func TestIce(t *testing.T) {
	result := Ice("Cold\nIcy")
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Ice should use RGB color codes")
	}
}

func TestNeon(t *testing.T) {
	result := Neon("Glow\nBright")
	if !strings.Contains(result, "\033[1m") {
		t.Error("Neon should include bold")
	}
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("Neon should use RGB colors")
	}
}

func TestRetroGreen(t *testing.T) {
	result := RetroGreen("Code")
	if !strings.Contains(result, "\033[38;2;0;") {
		t.Error("RetroGreen should use green RGB colors")
	}
}

func TestBold(t *testing.T) {
	result := Bold("Test")
	if result != "\033[1mTest\033[0m" {
		t.Errorf("Bold = %q", result)
	}
}

func TestItalic(t *testing.T) {
	result := Italic("Test")
	if result != "\033[3mTest\033[0m" {
		t.Errorf("Italic = %q", result)
	}
}

func TestUnderline(t *testing.T) {
	result := Underline("Test")
	if result != "\033[4mTest\033[0m" {
		t.Errorf("Underline = %q", result)
	}
}

func TestStrikethrough(t *testing.T) {
	result := Strikethrough("Test")
	if result != "\033[9mTest\033[0m" {
		t.Errorf("Strikethrough = %q", result)
	}
}

func TestBlinkFilter(t *testing.T) {
	result := Blink("Test")
	if result != "\033[5mTest\033[0m" {
		t.Errorf("Blink = %q", result)
	}
}

func TestDim(t *testing.T) {
	result := Dim("Test")
	if result != "\033[2mTest\033[0m" {
		t.Errorf("Dim = %q", result)
	}
}

func TestInvert(t *testing.T) {
	result := Invert("Test")
	if result != "\033[7mTest\033[0m" {
		t.Errorf("Invert = %q", result)
	}
}

func TestChain(t *testing.T) {
	result := Chain("Hello", []string{"bold", "underline"})
	if !strings.Contains(result, "\033[1m") {
		t.Error("Chain() should apply bold")
	}
	if !strings.Contains(result, "\033[4m") {
		t.Error("Chain() should apply underline")
	}
}

func TestChainEmpty(t *testing.T) {
	result := Chain("Hello", nil)
	if result != "Hello" {
		t.Errorf("Chain() with no filters should return input, got %q", result)
	}
}

func TestChainUnknownSkipped(t *testing.T) {
	result := Chain("Hello", []string{"nonexistent", "bold"})
	if !strings.Contains(result, "\033[1m") {
		t.Error("Chain should skip unknown and apply known filters")
	}
}

func TestParseChain(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"rainbow", 1},
		{"rainbow,bold,underline", 3},
		{"bold, italic, underline", 3},
		{" bold , ", 1},
	}
	for _, tt := range tests {
		got := ParseChain(tt.input)
		wantNil := tt.expected == 0
		if wantNil && got != nil {
			t.Errorf("ParseChain(%q) = %v, want nil", tt.input, got)
			continue
		}
		if !wantNil && len(got) != tt.expected {
			t.Errorf("ParseChain(%q) = %d items, want %d", tt.input, len(got), tt.expected)
		}
	}
}

func TestList(t *testing.T) {
	names := List()
	if len(names) == 0 {
		t.Fatal("List() returned empty")
	}
	if len(names) != len(registry) {
		t.Errorf("List() returned %d, registry has %d", len(names), len(registry))
	}
}

func TestListFilters(t *testing.T) {
	infos := ListFilters()
	if len(infos) < 20 {
		t.Errorf("ListFilters should return at least 20 entries, got %d", len(infos))
	}
	for _, info := range infos {
		if info.Name == "" || info.Desc == "" {
			t.Error("filter info should not have empty Name or Desc")
		}
	}
}
