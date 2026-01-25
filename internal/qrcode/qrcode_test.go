package qrcode

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	result, err := Generate("Hello", Options{})
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if result == "" {
		t.Fatal("Generate() returned empty")
	}
	// QR codes have multiple lines
	lines := strings.Split(result, "\n")
	if len(lines) < 5 {
		t.Errorf("QR code should have many lines, got %d", len(lines))
	}
}

func TestGenerateWithCharset(t *testing.T) {
	charsets := ListCharsets()
	for _, cs := range charsets {
		result, err := Generate("Test", Options{Charset: cs})
		if err != nil {
			t.Errorf("Generate with charset %q error: %v", cs, err)
			continue
		}
		if result == "" {
			t.Errorf("Generate with charset %q returned empty", cs)
		}
	}
}

func TestGenerateInvert(t *testing.T) {
	normal, err := Generate("Hi", Options{Invert: false})
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	inverted, err := Generate("Hi", Options{Invert: true})
	if err != nil {
		t.Fatalf("Generate(invert) error: %v", err)
	}
	if normal == inverted {
		t.Error("inverted QR should differ from normal")
	}
}

func TestGenerateURL(t *testing.T) {
	result, err := Generate("https://example.com", Options{})
	if err != nil {
		t.Fatalf("Generate(URL) error: %v", err)
	}
	if result == "" {
		t.Fatal("Generate(URL) returned empty")
	}
}

func TestGenerateEmpty(t *testing.T) {
	_, err := Generate("", Options{})
	if err == nil {
		t.Error("Generate('') should return error")
	}
}

func TestListCharsets(t *testing.T) {
	charsets := ListCharsets()
	if len(charsets) == 0 {
		t.Fatal("ListCharsets() returned empty")
	}
}

func TestGenerateCompact(t *testing.T) {
	result, err := GenerateCompact("Test", false)
	if err != nil {
		t.Fatalf("GenerateCompact() error: %v", err)
	}
	if result == "" {
		t.Fatal("GenerateCompact() returned empty")
	}
}
