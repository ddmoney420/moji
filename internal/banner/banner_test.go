package banner

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	result, err := Generate("Hi", "standard")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if result == "" {
		t.Fatal("Generate() returned empty")
	}
	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Error("banner should have multiple lines")
	}
}

func TestGenerateAllFonts(t *testing.T) {
	fonts := ListFonts()
	for _, font := range fonts {
		result, err := Generate("A", font.Name)
		if err != nil {
			t.Errorf("Generate with font %q error: %v", font.Name, err)
			continue
		}
		if result == "" {
			t.Errorf("Generate with font %q returned empty", font.Name)
		}
	}
}

func TestGenerateEmpty(t *testing.T) {
	result, err := Generate("", "standard")
	if err != nil {
		t.Fatalf("Generate('') error: %v", err)
	}
	// Empty input may still produce the font's empty output
	_ = result
}

func TestGenerateUnknownFont(t *testing.T) {
	// Unknown fonts fallback to standard
	result, err := Generate("Hi", "nonexistent_font_xyz")
	if err != nil {
		t.Fatalf("Generate with unknown font error: %v", err)
	}
	standard, _ := Generate("Hi", "standard")
	if result != standard {
		t.Error("unknown font should fallback to standard")
	}
}

func TestListFonts(t *testing.T) {
	fonts := ListFonts()
	if len(fonts) == 0 {
		t.Fatal("ListFonts() returned empty")
	}
	for _, f := range fonts {
		if f.Name == "" {
			t.Error("font has empty Name")
		}
	}
}

func TestGenerateMultiWord(t *testing.T) {
	result, err := Generate("Hello World", "standard")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if result == "" {
		t.Fatal("Generate() with multi-word returned empty")
	}
}

func TestGenerateSpecialChars(t *testing.T) {
	result, err := Generate("!@#", "standard")
	if err != nil {
		t.Fatalf("Generate(special chars) error: %v", err)
	}
	if result == "" {
		t.Fatal("Generate(special chars) returned empty")
	}
}
