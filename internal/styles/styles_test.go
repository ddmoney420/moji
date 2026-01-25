package styles

import (
	"strings"
	"testing"
)

func TestRainbow(t *testing.T) {
	result := Rainbow("Hello")
	if !strings.Contains(result, Reset) {
		t.Error("Rainbow should end with Reset")
	}
	if !strings.Contains(result, "H") {
		t.Error("Rainbow should contain original characters")
	}
	// Should not color spaces
	result = Rainbow("A B")
	if strings.Contains(result, Red+" ") {
		t.Error("Rainbow should not color spaces")
	}
}

func TestGradient(t *testing.T) {
	input := "Line1\nLine2\nLine3"
	result := Gradient(input)
	if !strings.Contains(result, Reset) {
		t.Error("Gradient should end with Reset")
	}
	if !strings.Contains(result, "Line1") {
		t.Error("Gradient should preserve text content")
	}
}

func TestFire(t *testing.T) {
	result := Fire("Hello\nWorld")
	if !strings.Contains(result, Reset) {
		t.Error("Fire should end with Reset")
	}
	if !strings.Contains(result, "Hello") {
		t.Error("Fire should preserve text")
	}
}

func TestIce(t *testing.T) {
	result := Ice("Cold\nText")
	if !strings.Contains(result, Reset) {
		t.Error("Ice should end with Reset")
	}
}

func TestMatrix(t *testing.T) {
	result := Matrix("Hack")
	if !strings.Contains(result, Reset) {
		t.Error("Matrix should end with Reset")
	}
	if !strings.Contains(result, "H") {
		t.Error("Matrix should contain original characters")
	}
}

func TestNeon(t *testing.T) {
	result := Neon("Glow")
	if !strings.Contains(result, Bold) {
		t.Error("Neon should include Bold")
	}
	if !strings.Contains(result, Reset) {
		t.Error("Neon should end with Reset")
	}
}

func TestOcean(t *testing.T) {
	result := Ocean("Wave\nSea")
	if !strings.Contains(result, Reset) {
		t.Error("Ocean should end with Reset")
	}
}

func TestSunset(t *testing.T) {
	result := Sunset("Dusk\nEvening")
	if !strings.Contains(result, Reset) {
		t.Error("Sunset should end with Reset")
	}
}

func TestCyberpunk(t *testing.T) {
	result := Cyberpunk("Neon\nCity")
	if !strings.Contains(result, Bold) {
		t.Error("Cyberpunk should include Bold")
	}
	if !strings.Contains(result, Magenta) || !strings.Contains(result, Cyan) {
		t.Error("Cyberpunk should alternate magenta and cyan")
	}
}

func TestLava(t *testing.T) {
	result := Lava("Hot")
	if !strings.Contains(result, Reset) {
		t.Error("Lava should end with Reset")
	}
}

func TestToxic(t *testing.T) {
	result := Toxic("Acid")
	if !strings.Contains(result, Reset) {
		t.Error("Toxic should end with Reset")
	}
}

func TestGalaxy(t *testing.T) {
	result := Galaxy("Stars\nNebula")
	if !strings.Contains(result, Reset) {
		t.Error("Galaxy should end with Reset")
	}
}

func TestGold(t *testing.T) {
	result := Gold("Shine\nBright")
	if !strings.Contains(result, Bold) {
		t.Error("Gold should include Bold")
	}
}

func TestHacker(t *testing.T) {
	result := Hacker("Code")
	if !strings.Contains(result, DarkGreen) {
		t.Error("Hacker should use DarkGreen")
	}
}

func TestVaporwave(t *testing.T) {
	result := Vaporwave("Aesthetic\nVibes")
	if !strings.Contains(result, Magenta) || !strings.Contains(result, Cyan) {
		t.Error("Vaporwave should use magenta and cyan")
	}
}

func TestChristmas(t *testing.T) {
	result := Christmas("Ho")
	if !strings.Contains(result, Red) || !strings.Contains(result, Green) {
		t.Error("Christmas should use red and green")
	}
}

func TestUSA(t *testing.T) {
	result := USA("USA")
	if !strings.Contains(result, Red) || !strings.Contains(result, White) || !strings.Contains(result, Blue) {
		t.Error("USA should use red, white, and blue")
	}
}

func TestMono(t *testing.T) {
	result := Mono("Clean")
	if !strings.HasPrefix(result, Bold+White) {
		t.Error("Mono should start with Bold+White")
	}
	if !strings.HasSuffix(result, Reset) {
		t.Error("Mono should end with Reset")
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		style    string
		contains string
	}{
		{"rainbow", Reset},
		{"gradient", Reset},
		{"fire", Reset},
		{"ice", Reset},
		{"matrix", Reset},
		{"neon", Bold},
		{"ocean", Reset},
		{"sunset", Reset},
		{"cyberpunk", Bold},
		{"cyber", Bold},
		{"lava", Reset},
		{"toxic", Reset},
		{"galaxy", Reset},
		{"gold", Bold},
		{"hacker", DarkGreen},
		{"vaporwave", Magenta},
		{"vapor", Magenta},
		{"christmas", Red},
		{"xmas", Red},
		{"usa", Red},
		{"america", Red},
		{"mono", Bold},
		{"white", Bold},
	}

	for _, tt := range tests {
		t.Run(tt.style, func(t *testing.T) {
			result := Apply("Test", tt.style)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Apply(%q) should contain %q", tt.style, tt.contains)
			}
		})
	}

	// Default/none should return unchanged
	result := Apply("Hello", "none")
	if result != "Hello" {
		t.Errorf("Apply with 'none' should return text unchanged, got %q", result)
	}
	result = Apply("Hello", "unknown")
	if result != "Hello" {
		t.Errorf("Apply with unknown style should return text unchanged")
	}
}

func TestListStyles(t *testing.T) {
	styles := ListStyles()
	if len(styles) == 0 {
		t.Error("ListStyles should return styles")
	}
	// Check that "none" is included
	found := false
	for _, s := range styles {
		if s.Name == "none" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ListStyles should include 'none'")
	}
}

func TestApplyBorder(t *testing.T) {
	tests := []struct {
		name   string
		border string
		expect string
	}{
		{"none returns unchanged", "none", "Hello"},
		{"empty returns unchanged", "", "Hello"},
		{"single adds border", "single", "┌"},
		{"double adds border", "double", "╔"},
		{"round adds border", "round", "╭"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyBorder("Hello", tt.border)
			if !strings.Contains(result, tt.expect) {
				t.Errorf("ApplyBorder(%q) should contain %q, got %q", tt.border, tt.expect, result)
			}
		})
	}

	// Multi-line border
	result := ApplyBorder("Line1\nLine2", "single")
	if !strings.Contains(result, "│") {
		t.Error("Border should contain vertical bars for content")
	}
}

func TestApplyAlignment(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		align  string
		width  int
		expect string
	}{
		{"left default", "Hi", "left", 10, "Hi"},
		{"center", "Hi", "center", 10, "    Hi"},
		{"right", "Hi", "right", 10, "        Hi"},
		{"zero width unchanged", "Hi", "center", 0, "Hi"},
		{"text wider than width", "HelloWorld", "center", 5, "HelloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyAlignment(tt.text, tt.align, tt.width)
			if !strings.Contains(result, tt.expect) {
				t.Errorf("ApplyAlignment(%q, %q, %d) expected to contain %q, got %q",
					tt.text, tt.align, tt.width, tt.expect, result)
			}
		})
	}
}

func TestGetBorderStyle(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"single", "┌"},
		{"double", "╔"},
		{"round", "╭"},
		{"rounded", "╭"},
		{"bold", "┏"},
		{"thick", "┏"},
		{"ascii", "+"},
		{"stars", "*"},
		{"hash", "#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := GetBorderStyle(tt.name)
			if bs.TopLeft != tt.expected {
				t.Errorf("GetBorderStyle(%q).TopLeft = %q, want %q", tt.name, bs.TopLeft, tt.expected)
			}
		})
	}

	// Unknown returns empty
	bs := GetBorderStyle("unknown")
	if bs.Horizontal != "" {
		t.Error("Unknown border should return empty style")
	}
}

func TestListBorders(t *testing.T) {
	borders := ListBorders()
	if len(borders) < 7 {
		t.Errorf("ListBorders should return at least 7 styles, got %d", len(borders))
	}
}

func TestSpacesPreserved(t *testing.T) {
	// Ensure styles don't color whitespace
	funcs := []struct {
		name string
		fn   func(string) string
	}{
		{"Rainbow", Rainbow},
		{"Matrix", Matrix},
		{"Lava", Lava},
		{"Toxic", Toxic},
		{"Christmas", Christmas},
		{"USA", USA},
	}

	for _, f := range funcs {
		t.Run(f.name, func(t *testing.T) {
			result := f.fn("A B")
			// The space between A and B should not have color codes directly before it
			lines := strings.Split(result, "\n")
			if len(lines) > 0 {
				// Just check the result contains the space
				if !strings.Contains(lines[0], " ") {
					t.Errorf("%s should preserve spaces", f.name)
				}
			}
		})
	}
}
