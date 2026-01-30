package themes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	// Clear registry for this test
	registry = make(map[string]*Theme)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Should have loaded built-in themes
	if len(registry) == 0 {
		t.Fatal("Init() should register built-in themes")
	}

	// Check a few expected built-in themes exist
	expectedThemes := []string{"rainbow", "neon", "fire", "ice"}
	for _, name := range expectedThemes {
		if _, ok := registry[name]; !ok {
			t.Errorf("Expected theme %q not found in registry", name)
		}
	}
}

func TestRegisterTheme(t *testing.T) {
	// Clear registry for this test
	registry = make(map[string]*Theme)

	theme := &Theme{
		Name:        "test-theme",
		Description: "Test theme",
		Colors:      []string{"#FF0000", "#00FF00", "#0000FF"},
	}

	err := RegisterTheme(theme)
	if err != nil {
		t.Fatalf("RegisterTheme() failed: %v", err)
	}

	retrieved, err := GetTheme("test-theme")
	if err != nil {
		t.Fatalf("GetTheme() failed: %v", err)
	}

	if retrieved.Name != "test-theme" {
		t.Errorf("Expected theme name %q, got %q", "test-theme", retrieved.Name)
	}
}

func TestGetTheme(t *testing.T) {
	// Clear and init registry
	registry = make(map[string]*Theme)
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	theme, err := GetTheme("rainbow")
	if err != nil {
		t.Fatalf("GetTheme(rainbow) failed: %v", err)
	}

	if theme.Name != "rainbow" {
		t.Errorf("Expected theme name %q, got %q", "rainbow", theme.Name)
	}

	if len(theme.Colors) == 0 {
		t.Error("Rainbow theme should have colors")
	}
}

func TestGetThemeNotFound(t *testing.T) {
	// Clear and init registry
	registry = make(map[string]*Theme)
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	_, err := GetTheme("nonexistent-theme")
	if err == nil {
		t.Fatal("GetTheme() should return error for nonexistent theme")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error message should mention 'not found', got: %v", err)
	}
}

func TestListThemes(t *testing.T) {
	// Clear and init registry
	registry = make(map[string]*Theme)
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	themes := ListThemes()
	if len(themes) == 0 {
		t.Fatal("ListThemes() should return theme names")
	}

	// Check for expected built-in themes
	found := make(map[string]bool)
	for _, name := range themes {
		found[name] = true
	}

	expectedThemes := []string{"rainbow", "neon", "fire", "ice", "matrix"}
	for _, expected := range expectedThemes {
		if !found[expected] {
			t.Errorf("Expected theme %q not found in ListThemes()", expected)
		}
	}
}

func TestValidateTheme(t *testing.T) {
	tests := []struct {
		name      string
		theme     *Theme
		shouldErr bool
	}{
		{
			name:      "nil theme",
			theme:     nil,
			shouldErr: true,
		},
		{
			name: "theme with no name",
			theme: &Theme{
				Description: "No name",
				Colors:      []string{"#FF0000"},
			},
			shouldErr: true,
		},
		{
			name: "theme with no colors",
			theme: &Theme{
				Name:        "empty",
				Description: "No colors",
				Colors:      []string{},
			},
			shouldErr: true,
		},
		{
			name: "theme with invalid color format",
			theme: &Theme{
				Name:        "invalid",
				Description: "Invalid color",
				Colors:      []string{"#GGGGGG"},
			},
			shouldErr: true,
		},
		{
			name: "theme with short hex color",
			theme: &Theme{
				Name:        "short",
				Description: "Short hex",
				Colors:      []string{"#FFF"},
			},
			shouldErr: true,
		},
		{
			name: "valid theme",
			theme: &Theme{
				Name:        "valid",
				Description: "Valid theme",
				Colors:      []string{"#FF0000", "#00FF00", "#0000FF"},
			},
			shouldErr: false,
		},
		{
			name: "valid theme with lowercase hex",
			theme: &Theme{
				Name:        "lowercase",
				Description: "Lowercase hex",
				Colors:      []string{"#ff0000", "#00ff00"},
			},
			shouldErr: false,
		},
		{
			name: "theme with metadata",
			theme: &Theme{
				Name:        "with-metadata",
				Description: "Has metadata",
				Colors:      []string{"#FF0000"},
				Metadata: map[string]string{
					"author":  "test",
					"version": "1.0",
				},
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTheme(tt.theme)
			if (err != nil) != tt.shouldErr {
				t.Errorf("validateTheme() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}

func TestLoadTheme(t *testing.T) {
	// Create a temporary YAML file
	tmpDir := t.TempDir()
	themePath := filepath.Join(tmpDir, "test-theme.yaml")

	content := `name: test-custom-theme
description: A custom test theme
colors:
  - "#FF0000"
  - "#00FF00"
  - "#0000FF"
metadata:
  author: test-user
`

	if err := os.WriteFile(themePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test theme file: %v", err)
	}

	theme, err := LoadTheme(themePath)
	if err != nil {
		t.Fatalf("LoadTheme() failed: %v", err)
	}

	if theme.Name != "test-custom-theme" {
		t.Errorf("Expected name %q, got %q", "test-custom-theme", theme.Name)
	}

	if len(theme.Colors) != 3 {
		t.Errorf("Expected 3 colors, got %d", len(theme.Colors))
	}

	if theme.Colors[0] != "#FF0000" {
		t.Errorf("Expected first color #FF0000, got %s", theme.Colors[0])
	}

	if theme.Metadata["author"] != "test-user" {
		t.Errorf("Expected metadata author 'test-user', got %q", theme.Metadata["author"])
	}
}

func TestLoadThemeInvalidFile(t *testing.T) {
	_, err := LoadTheme("/nonexistent/path/theme.yaml")
	if err == nil {
		t.Fatal("LoadTheme() should return error for nonexistent file")
	}
}

func TestLoadThemeInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	themePath := filepath.Join(tmpDir, "invalid.yaml")

	content := `{invalid yaml: [`
	if err := os.WriteFile(themePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := LoadTheme(themePath)
	if err == nil {
		t.Fatal("LoadTheme() should return error for invalid YAML")
	}
}

func TestLoadThemeInvalidColors(t *testing.T) {
	tmpDir := t.TempDir()
	themePath := filepath.Join(tmpDir, "bad-colors.yaml")

	content := `name: bad-theme
description: Bad colors
colors:
  - "#GGGGGG"
`

	if err := os.WriteFile(themePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := LoadTheme(themePath)
	if err == nil {
		t.Fatal("LoadTheme() should return error for invalid color format")
	}

	if !strings.Contains(err.Error(), "invalid color format") {
		t.Errorf("Error should mention invalid color format, got: %v", err)
	}
}

func TestSaveTheme(t *testing.T) {
	tmpDir := t.TempDir()
	themePath := filepath.Join(tmpDir, "saved-theme.yaml")

	theme := &Theme{
		Name:        "saved-test",
		Description: "Saved test theme",
		Colors:      []string{"#FF0000", "#00FF00"},
		Metadata: map[string]string{
			"author": "test",
		},
	}

	err := SaveTheme(theme, themePath)
	if err != nil {
		t.Fatalf("SaveTheme() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(themePath); err != nil {
		t.Fatalf("Saved theme file not found: %v", err)
	}

	// Load it back and verify
	loaded, err := LoadTheme(themePath)
	if err != nil {
		t.Fatalf("Failed to load saved theme: %v", err)
	}

	if loaded.Name != theme.Name {
		t.Errorf("Loaded name %q doesn't match original %q", loaded.Name, theme.Name)
	}

	if len(loaded.Colors) != len(theme.Colors) {
		t.Errorf("Loaded colors count %d doesn't match original %d", len(loaded.Colors), len(theme.Colors))
	}
}

func TestSaveThemeInvalidTheme(t *testing.T) {
	tmpDir := t.TempDir()
	themePath := filepath.Join(tmpDir, "invalid-save.yaml")

	invalidTheme := &Theme{
		Name:   "invalid",
		Colors: []string{}, // Invalid: no colors
	}

	err := SaveTheme(invalidTheme, themePath)
	if err == nil {
		t.Fatal("SaveTheme() should return error for invalid theme")
	}
}

func TestLoadThemesDir(t *testing.T) {
	tmpDir := t.TempDir()
	themesDir := filepath.Join(tmpDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatalf("Failed to create themes directory: %v", err)
	}

	// Create multiple theme files
	theme1 := `name: theme1
description: First theme
colors:
  - "#FF0000"
  - "#00FF00"
`

	theme2 := `name: theme2
description: Second theme
colors:
  - "#0000FF"
`

	if err := os.WriteFile(filepath.Join(themesDir, "theme1.yaml"), []byte(theme1), 0644); err != nil {
		t.Fatalf("Failed to create theme1: %v", err)
	}

	if err := os.WriteFile(filepath.Join(themesDir, "theme2.yml"), []byte(theme2), 0644); err != nil {
		t.Fatalf("Failed to create theme2: %v", err)
	}

	// Clear registry for this test
	registry = make(map[string]*Theme)

	loaded, err := LoadThemesDir(themesDir)
	if err != nil {
		t.Fatalf("LoadThemesDir() failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("Expected 2 themes loaded, got %d", len(loaded))
	}

	// Check registry
	if _, ok := registry["theme1"]; !ok {
		t.Error("theme1 not registered")
	}

	if _, ok := registry["theme2"]; !ok {
		t.Error("theme2 not registered")
	}
}

func TestLoadThemesDirNonexistent(t *testing.T) {
	_, err := LoadThemesDir("/nonexistent/path")
	if err == nil {
		t.Fatal("LoadThemesDir() should return error for nonexistent directory")
	}
}

func TestLoadThemesDirIgnoresNonYAML(t *testing.T) {
	tmpDir := t.TempDir()
	themesDir := filepath.Join(tmpDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create a YAML file and a non-YAML file
	validTheme := `name: valid
description: Valid
colors:
  - "#FF0000"
`

	if err := os.WriteFile(filepath.Join(themesDir, "valid.yaml"), []byte(validTheme), 0644); err != nil {
		t.Fatalf("Failed to create valid theme: %v", err)
	}

	if err := os.WriteFile(filepath.Join(themesDir, "ignored.txt"), []byte("ignored"), 0644); err != nil {
		t.Fatalf("Failed to create txt file: %v", err)
	}

	registry = make(map[string]*Theme)

	loaded, err := LoadThemesDir(themesDir)
	if err != nil {
		t.Fatalf("LoadThemesDir() failed: %v", err)
	}

	if len(loaded) != 1 {
		t.Errorf("Expected 1 theme (txt ignored), got %d", len(loaded))
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex       string
		r, g, b   uint8
		shouldErr bool
	}{
		{"#FF0000", 255, 0, 0, false},
		{"#00FF00", 0, 255, 0, false},
		{"#0000FF", 0, 0, 255, false},
		{"#FFFFFF", 255, 255, 255, false},
		{"#000000", 0, 0, 0, false},
		{"#ff0000", 255, 0, 0, false}, // lowercase
		{"#GGGGGG", 0, 0, 0, true},    // invalid hex
		{"#FFF", 0, 0, 0, true},       // too short
		{"FF0000", 0, 0, 0, true},     // no #
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			r, g, b, err := HexToRGB(tt.hex)
			if (err != nil) != tt.shouldErr {
				t.Errorf("HexToRGB() error = %v, shouldErr %v", err, tt.shouldErr)
				return
			}

			if !tt.shouldErr && (r != tt.r || g != tt.g || b != tt.b) {
				t.Errorf("HexToRGB(%s) = (%d,%d,%d), expected (%d,%d,%d)", tt.hex, r, g, b, tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestRGBToHex(t *testing.T) {
	tests := []struct {
		r, g, b  uint8
		expected string
	}{
		{255, 0, 0, "#FF0000"},
		{0, 255, 0, "#00FF00"},
		{0, 0, 255, "#0000FF"},
		{255, 255, 255, "#FFFFFF"},
		{0, 0, 0, "#000000"},
		{128, 64, 32, "#804020"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := RGBToHex(tt.r, tt.g, tt.b)
			if result != tt.expected {
				t.Errorf("RGBToHex(%d,%d,%d) = %s, expected %s", tt.r, tt.g, tt.b, result, tt.expected)
			}
		})
	}
}

func TestBuiltinThemesValid(t *testing.T) {
	// Clear and init registry
	registry = make(map[string]*Theme)
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	for name, theme := range registry {
		t.Run(name, func(t *testing.T) {
			if err := validateTheme(theme); err != nil {
				t.Errorf("Built-in theme %q failed validation: %v", name, err)
			}

			if theme.Description == "" {
				t.Errorf("Theme %q has empty description", name)
			}
		})
	}
}

// TestRegistryIsolation ensures tests don't interfere with each other
func TestRegistryIsolation(t *testing.T) {
	originalSize := len(registry)

	// Create a new theme
	newTheme := &Theme{
		Name:        "isolation-test",
		Description: "Test isolation",
		Colors:      []string{"#FF0000"},
	}

	if err := RegisterTheme(newTheme); err != nil {
		t.Fatalf("RegisterTheme() failed: %v", err)
	}

	// Check it was added
	if _, err := GetTheme("isolation-test"); err != nil {
		t.Fatal("Theme not registered")
	}

	// Test doesn't cleanup, but that's OK for this test suite
	// In production use, Init() should be called at app startup
	if len(registry) <= originalSize {
		t.Error("Registry size should increase after RegisterTheme")
	}
}
