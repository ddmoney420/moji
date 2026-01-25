package tree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.MaxDepth != -1 {
		t.Errorf("DefaultOptions MaxDepth = %d, want -1", opts.MaxDepth)
	}
	if opts.ShowHidden {
		t.Error("DefaultOptions ShowHidden should be false")
	}
	if opts.Indent != "│   " {
		t.Errorf("DefaultOptions Indent = %q, want %q", opts.Indent, "│   ")
	}
}

func TestGenerate(t *testing.T) {
	// Create temp directory structure
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "file1.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(dir, "subdir", "nested.go"), []byte("package sub"), 0644)
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("hidden"), 0644)

	opts := DefaultOptions()
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if !entry.IsDir {
		t.Error("Root should be a directory")
	}
	if len(entry.Entries) == 0 {
		t.Error("Root should have entries")
	}

	// Hidden files should be excluded by default
	for _, e := range entry.Entries {
		if strings.HasPrefix(e.Name, ".") {
			t.Errorf("Hidden file %q should be excluded", e.Name)
		}
	}
}

func TestGenerateShowHidden(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("hidden"), 0644)
	os.WriteFile(filepath.Join(dir, "visible"), []byte("visible"), 0644)

	opts := DefaultOptions()
	opts.ShowHidden = true
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	foundHidden := false
	for _, e := range entry.Entries {
		if e.Name == ".hidden" {
			foundHidden = true
		}
	}
	if !foundHidden {
		t.Error("ShowHidden=true should include hidden files")
	}
}

func TestGenerateMaxDepth(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "a", "b", "c"), 0755)
	os.WriteFile(filepath.Join(dir, "a", "b", "c", "deep.txt"), []byte("deep"), 0644)

	opts := DefaultOptions()
	opts.MaxDepth = 1
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Should have 'a' directory
	if len(entry.Entries) == 0 {
		t.Fatal("Should have at least one entry")
	}

	// 'a' should have no sub-entries due to depth limit
	for _, e := range entry.Entries {
		if e.Name == "a" && len(e.Entries) != 0 {
			t.Error("MaxDepth=1 should prevent deeper traversal")
		}
	}
}

func TestGenerateDirsOnly(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello"), 0644)

	opts := DefaultOptions()
	opts.DirsOnly = true
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	for _, e := range entry.Entries {
		if !e.IsDir {
			t.Errorf("DirsOnly should exclude files, found %q", e.Name)
		}
	}
}

func TestGenerateFilesOnly(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello"), 0644)

	opts := DefaultOptions()
	opts.FilesOnly = true
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	for _, e := range entry.Entries {
		if e.IsDir {
			t.Errorf("FilesOnly should exclude dirs, found %q", e.Name)
		}
	}
}

func TestGeneratePattern(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "readme.md"), []byte("# readme"), 0644)
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("package test"), 0644)

	opts := DefaultOptions()
	opts.Pattern = "*.go"
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	for _, e := range entry.Entries {
		if !strings.HasSuffix(e.Name, ".go") {
			t.Errorf("Pattern *.go should exclude %q", e.Name)
		}
	}
}

func TestGenerateMaxItems(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 10; i++ {
		os.WriteFile(filepath.Join(dir, strings.Repeat("a", i+1)+".txt"), []byte("x"), 0644)
	}

	opts := DefaultOptions()
	opts.MaxItems = 3
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if len(entry.Entries) > 3 {
		t.Errorf("MaxItems=3 but got %d entries", len(entry.Entries))
	}
}

func TestFormat(t *testing.T) {
	entry := Entry{
		Name:  "root",
		IsDir: true,
		Entries: []Entry{
			{Name: "subdir", IsDir: true, Entries: []Entry{
				{Name: "file.go", IsDir: false},
			}},
			{Name: "readme.md", IsDir: false},
		},
	}

	output := Format(entry, DefaultOptions())
	if !strings.Contains(output, "root") {
		t.Error("Format should contain root name")
	}
	if !strings.Contains(output, "subdir") {
		t.Error("Format should contain subdir")
	}
	if !strings.Contains(output, "readme.md") {
		t.Error("Format should contain readme.md")
	}
	if !strings.Contains(output, "├──") && !strings.Contains(output, "└──") {
		t.Error("Format should contain tree connectors")
	}
}

func TestGetFileIcon(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"main.go", "🔵"},
		{"script.py", "🐍"},
		{"app.js", "🟨"},
		{"unknown.xyz", "📄"},
	}
	for _, tt := range tests {
		got := getFileIcon(tt.name)
		if got != tt.want {
			t.Errorf("getFileIcon(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestSimple(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte("hello"), 0644)

	output := Simple(dir, 1)
	if output == "" {
		t.Error("Simple() returned empty")
	}
	if !strings.Contains(output, "test.txt") {
		t.Error("Simple() should contain test.txt")
	}
}

func TestGenerateInvalidPath(t *testing.T) {
	_, err := Generate("/nonexistent_path_xyz_123", DefaultOptions())
	if err == nil {
		t.Error("Generate with invalid path should return error")
	}
}

func TestSortBySize(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "small.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "big.txt"), []byte(strings.Repeat("x", 1000)), 0644)

	opts := DefaultOptions()
	opts.SortBySize = true
	entry, err := Generate(dir, opts)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if len(entry.Entries) >= 2 {
		if entry.Entries[0].Size < entry.Entries[1].Size {
			t.Error("SortBySize should sort largest first")
		}
	}
}
