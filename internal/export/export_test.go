package export

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		hex  string
		want color.RGBA
	}{
		{"#000000", color.RGBA{0, 0, 0, 255}},
		{"000000", color.RGBA{0, 0, 0, 255}},
		{"#ffffff", color.RGBA{255, 255, 255, 255}},
		{"#ff0000", color.RGBA{255, 0, 0, 255}},
		{"#00ff00", color.RGBA{0, 255, 0, 255}},
		{"#0000ff", color.RGBA{0, 0, 255, 255}},
		{"#0d1117", color.RGBA{13, 17, 23, 255}},
		{"invalid", color.RGBA{0, 0, 0, 255}}, // fallback to black
		{"", color.RGBA{0, 0, 0, 255}},        // fallback to black
		{"#fff", color.RGBA{0, 0, 0, 255}},    // too short, fallback
	}

	for _, tt := range tests {
		got := parseHexColor(tt.hex)
		if got != tt.want {
			t.Errorf("parseHexColor(%q) = %v, want %v", tt.hex, got, tt.want)
		}
	}
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"<tag>", "&lt;tag&gt;"},
		{"a & b", "a &amp; b"},
		{"it's", "it&apos;s"},
		{`say "hi"`, "say &quot;hi&quot;"},
		{"<a href=\"x\">&</a>", "&lt;a href=&quot;x&quot;&gt;&amp;&lt;/a&gt;"},
	}

	for _, tt := range tests {
		got := escapeXML(tt.input)
		if got != tt.want {
			t.Errorf("escapeXML(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestToPNG(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "test.png")

	err := ToPNG("Hello\nWorld", out, "#000000", "#ffffff")
	if err != nil {
		t.Fatalf("ToPNG() error: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Output PNG file is empty")
	}

	// Check PNG magic bytes
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	if len(data) < 8 || string(data[1:4]) != "PNG" {
		t.Error("Output is not a valid PNG file")
	}
}

func TestToPNGInvalidPath(t *testing.T) {
	err := ToPNG("Hello", "/nonexistent/dir/test.png", "#000000", "#ffffff")
	if err == nil {
		t.Error("ToPNG with invalid path should return error")
	}
}

func TestToSVG(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "test.svg")

	err := ToSVG("Hello\nWorld", out, "#000000", "#ffffff")
	if err != nil {
		t.Fatalf("ToSVG() error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "<svg") {
		t.Error("SVG output missing <svg> tag")
	}
	if !strings.Contains(content, "Hello") {
		t.Error("SVG output missing text content")
	}
	if !strings.Contains(content, "World") {
		t.Error("SVG output missing second line")
	}
	if !strings.Contains(content, "monospace") {
		t.Error("SVG output missing monospace font")
	}
}

func TestToSVGEscaping(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "test.svg")

	err := ToSVG("<script>alert('xss')</script>", out, "#000000", "#ffffff")
	if err != nil {
		t.Fatalf("ToSVG() error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "<script>") {
		t.Error("SVG output contains unescaped HTML - XSS vulnerability")
	}
	if !strings.Contains(content, "&lt;script&gt;") {
		t.Error("SVG output should contain escaped HTML entities")
	}
}

func TestToSVGInvalidPath(t *testing.T) {
	err := ToSVG("Hello", "/nonexistent/dir/test.svg", "#000000", "#ffffff")
	if err == nil {
		t.Error("ToSVG with invalid path should return error")
	}
}

func TestToHTML(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "test.html")

	err := ToHTML("Hello World", out, "#000000", "#ffffff", "Test")
	if err != nil {
		t.Fatalf("ToHTML() error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "<!DOCTYPE html>") {
		t.Error("HTML output missing doctype")
	}
	if !strings.Contains(content, "<title>Test</title>") {
		t.Error("HTML output missing title")
	}
	if !strings.Contains(content, "Hello World") {
		t.Error("HTML output missing text content")
	}
	if !strings.Contains(content, "#000000") {
		t.Error("HTML output missing background color")
	}
}

func TestToHTMLInvalidPath(t *testing.T) {
	err := ToHTML("Hello", "/nonexistent/dir/test.html", "#000", "#fff", "Test")
	if err == nil {
		t.Error("ToHTML with invalid path should return error")
	}
}

func TestToPNGMultiline(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "multi.png")

	art := `  /\_/\
 ( o.o )
  > ^ <`

	err := ToPNG(art, out, "#0d1117", "#e6edf3")
	if err != nil {
		t.Fatalf("ToPNG() with multiline art error: %v", err)
	}

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Output PNG file is empty")
	}
}
