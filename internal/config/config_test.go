package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// Helper to create temporary test directory
func setupTestDir(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "moji-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return tmpDir, func() {
		os.RemoveAll(tmpDir)
	}
}

// Helper to create a test config file
func writeConfigFile(t *testing.T, path string, cfg *Config) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}
}

// ============================================================================
// Config Path Tests
// ============================================================================

func TestConfigPath_WithXDGConfigHome(t *testing.T) {
	// Save original env
	original := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", original)

	// Set custom XDG_CONFIG_HOME
	testPath := "/custom/config"
	os.Setenv("XDG_CONFIG_HOME", testPath)

	path := ConfigPath()
	expected := filepath.Join(testPath, "moji", "config.yaml")

	if path != expected {
		t.Errorf("ConfigPath() = %q, want %q", path, expected)
	}
}

func TestConfigPath_WithoutXDGConfigHome(t *testing.T) {
	// Save original env
	original := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", original)

	// Unset XDG_CONFIG_HOME
	os.Setenv("XDG_CONFIG_HOME", "")

	path := ConfigPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "moji", "config.yaml")

	if path != expected {
		t.Errorf("ConfigPath() = %q, want %q", path, expected)
	}
}

func TestLegacyConfigPath(t *testing.T) {
	path := LegacyConfigPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".mojirc")

	if path != expected {
		t.Errorf("LegacyConfigPath() = %q, want %q", path, expected)
	}
}

// ============================================================================
// Default Configuration Tests
// ============================================================================

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Error("DefaultConfig() returned nil")
	}

	// Check defaults are set
	if cfg.Defaults.BannerFont != "standard" {
		t.Errorf("BannerFont = %q, want 'standard'", cfg.Defaults.BannerFont)
	}
	if cfg.Defaults.ConvertWidth != 80 {
		t.Errorf("ConvertWidth = %d, want 80", cfg.Defaults.ConvertWidth)
	}
	if cfg.Defaults.GradientTheme != "rainbow" {
		t.Errorf("GradientTheme = %q, want 'rainbow'", cfg.Defaults.GradientTheme)
	}

	// Check built-in presets exist
	if len(cfg.Presets) < 3 {
		t.Errorf("Expected at least 3 presets, got %d", len(cfg.Presets))
	}

	// Check specific presets
	presets := []string{"neon-banner", "retro", "fire-banner"}
	for _, name := range presets {
		if _, ok := cfg.Presets[name]; !ok {
			t.Errorf("Missing preset: %s", name)
		}
	}
}

func TestDefaultConfig_AllDefaultsSet(t *testing.T) {
	cfg := DefaultConfig()

	// Verify all defaults fields have non-zero values
	tests := []struct {
		name  string
		value interface{}
	}{
		{"BannerFont", cfg.Defaults.BannerFont},
		{"BannerStyle", cfg.Defaults.BannerStyle},
		{"BannerBorder", cfg.Defaults.BannerBorder},
		{"ConvertWidth", cfg.Defaults.ConvertWidth},
		{"ConvertCharset", cfg.Defaults.ConvertCharset},
		{"ConvertDither", cfg.Defaults.ConvertDither},
		{"GradientTheme", cfg.Defaults.GradientTheme},
		{"GradientMode", cfg.Defaults.GradientMode},
		{"BubbleStyle", cfg.Defaults.BubbleStyle},
		{"BubbleWidth", cfg.Defaults.BubbleWidth},
		{"QRCharset", cfg.Defaults.QRCharset},
		{"PatternBorder", cfg.Defaults.PatternBorder},
		{"PatternDivider", cfg.Defaults.PatternDivider},
	}

	for _, test := range tests {
		if test.value == "" || test.value == 0 {
			t.Errorf("Default %s is not set", test.name)
		}
	}
}

// ============================================================================
// Preset Tests
// ============================================================================

func TestGetPreset_Exists(t *testing.T) {
	cfg := DefaultConfig()
	preset, ok := cfg.GetPreset("neon-banner")

	if !ok {
		t.Error("GetPreset returned false for existing preset")
	}
	if preset.Name != "Neon Banner" {
		t.Errorf("Preset name = %q, want 'Neon Banner'", preset.Name)
	}
	if preset.Command != "banner" {
		t.Errorf("Preset command = %q, want 'banner'", preset.Command)
	}
}

func TestGetPreset_NotExists(t *testing.T) {
	cfg := DefaultConfig()
	_, ok := cfg.GetPreset("nonexistent")

	if ok {
		t.Error("GetPreset returned true for non-existent preset")
	}
}

func TestSetPreset_New(t *testing.T) {
	cfg := &Config{Presets: map[string]Preset{}}
	preset := Preset{
		Name:    "Custom",
		Command: "art",
		Args:    map[string]string{"style": "custom"},
	}

	cfg.SetPreset("custom", preset)

	retrieved, ok := cfg.GetPreset("custom")
	if !ok {
		t.Error("SetPreset didn't save preset")
	}
	if retrieved.Name != "Custom" {
		t.Errorf("Retrieved preset name = %q, want 'Custom'", retrieved.Name)
	}
}

func TestSetPreset_Override(t *testing.T) {
	cfg := DefaultConfig()
	newPreset := Preset{
		Name:        "Updated Neon",
		Description: "Updated description",
		Command:     "banner",
		Args:        map[string]string{"gradient": "updated"},
	}

	cfg.SetPreset("neon-banner", newPreset)

	retrieved, _ := cfg.GetPreset("neon-banner")
	if retrieved.Name != "Updated Neon" {
		t.Errorf("Preset not updated, name = %q", retrieved.Name)
	}
}

func TestListPresets(t *testing.T) {
	cfg := DefaultConfig()
	presets := cfg.ListPresets()

	if len(presets) < 3 {
		t.Errorf("ListPresets returned %d presets, want at least 3", len(presets))
	}

	// Check that all expected presets are in the list
	presetMap := make(map[string]bool)
	for _, p := range presets {
		presetMap[p] = true
	}

	expected := []string{"neon-banner", "retro", "fire-banner"}
	for _, name := range expected {
		if !presetMap[name] {
			t.Errorf("Expected preset %q not in list", name)
		}
	}
}

func TestListPresets_NilPresets(t *testing.T) {
	cfg := &Config{}
	presets := cfg.ListPresets()

	// Should handle nil presets gracefully
	if presets != nil && len(presets) > 0 {
		t.Errorf("ListPresets should return empty for nil presets, got %d", len(presets))
	}
}

// ============================================================================
// Alias Tests
// ============================================================================

func TestGetAlias_Exists(t *testing.T) {
	cfg := &Config{
		Aliases: map[string]string{
			"art-demo": "art mountains",
		},
	}

	cmd, ok := cfg.GetAlias("art-demo")
	if !ok {
		t.Error("GetAlias returned false for existing alias")
	}
	if cmd != "art mountains" {
		t.Errorf("GetAlias returned %q, want 'art mountains'", cmd)
	}
}

func TestGetAlias_NotExists(t *testing.T) {
	cfg := &Config{Aliases: map[string]string{}}
	_, ok := cfg.GetAlias("nonexistent")

	if ok {
		t.Error("GetAlias returned true for non-existent alias")
	}
}

func TestSetAlias_New(t *testing.T) {
	cfg := &Config{Aliases: map[string]string{}}

	cfg.SetAlias("test-alias", "test command")

	cmd, ok := cfg.GetAlias("test-alias")
	if !ok {
		t.Error("SetAlias didn't save alias")
	}
	if cmd != "test command" {
		t.Errorf("Retrieved command = %q, want 'test command'", cmd)
	}
}

func TestSetAlias_Override(t *testing.T) {
	cfg := &Config{
		Aliases: map[string]string{
			"test": "old command",
		},
	}

	cfg.SetAlias("test", "new command")

	cmd, _ := cfg.GetAlias("test")
	if cmd != "new command" {
		t.Errorf("Alias not updated, command = %q", cmd)
	}
}

func TestSetAlias_NilAliases(t *testing.T) {
	cfg := &Config{}

	cfg.SetAlias("test", "command")

	cmd, ok := cfg.GetAlias("test")
	if !ok {
		t.Error("SetAlias failed with nil aliases map")
	}
	if cmd != "command" {
		t.Errorf("Retrieved command = %q, want 'command'", cmd)
	}
}

// ============================================================================
// Save and Persistence Tests
// ============================================================================

func TestSave_CreatesDirectory(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Override ConfigPath to use temp directory
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := DefaultConfig()
	err := cfg.Save()

	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Check directory was created
	configDir := filepath.Join(tmpDir, "moji")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Save() didn't create config directory")
	}

	// Check file was created
	configFile := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Save() didn't create config file")
	}
}

func TestSave_WritesValidYAML(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := DefaultConfig()
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify we can read it back
	configFile := filepath.Join(tmpDir, "moji", "config.yaml")
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	if len(data) == 0 {
		t.Error("Saved config file is empty")
	}
}

// ============================================================================
// Load and Round-Trip Tests
// ============================================================================

func TestLoad_DefaultsWhenNoFile(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Mock home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Use non-existent paths
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, "nonexistent"))
	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Error("Load() returned nil config")
	}

	// Should have default values
	if cfg.Defaults.BannerFont != "standard" {
		t.Errorf("Expected default BannerFont, got %q", cfg.Defaults.BannerFont)
	}
}

func TestLoad_FromXDGPath(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create config with custom values
	cfg := DefaultConfig()
	cfg.Defaults.BannerFont = "custom-font"
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load it back
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded.Defaults.BannerFont != "custom-font" {
		t.Errorf("BannerFont = %q, want 'custom-font'", loaded.Defaults.BannerFont)
	}
}

func TestLoad_FromLegacyPath(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("HOME", originalHome)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Set temp as home
	os.Setenv("HOME", tmpDir)
	// Make XDG path non-existent
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, "nonexistent-xdg"))

	// Create legacy config
	cfg := DefaultConfig()
	cfg.Defaults.BannerFont = "legacy-font"

	// Write as YAML to legacy path
	legacyPath := filepath.Join(tmpDir, ".mojirc")
	data, _ := yaml.Marshal(cfg)
	if err := os.WriteFile(legacyPath, data, 0644); err != nil {
		t.Fatalf("failed to write legacy config: %v", err)
	}

	// Load should find legacy config
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded.Defaults.BannerFont != "legacy-font" {
		t.Errorf("BannerFont = %q, want 'legacy-font'", loaded.Defaults.BannerFont)
	}
}

func TestRoundTrip_SaveAndLoad(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create config with custom values
	original := DefaultConfig()
	original.Defaults.BannerFont = "custom"
	original.Defaults.ConvertWidth = 120
	original.SetAlias("custom-alias", "custom command")
	original.SetPreset("custom-preset", Preset{
		Name:    "Custom",
		Command: "art",
		Args:    map[string]string{"style": "test"},
	})

	if err := original.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load it back
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify values match
	if loaded.Defaults.BannerFont != "custom" {
		t.Errorf("BannerFont = %q, want 'custom'", loaded.Defaults.BannerFont)
	}
	if loaded.Defaults.ConvertWidth != 120 {
		t.Errorf("ConvertWidth = %d, want 120", loaded.Defaults.ConvertWidth)
	}

	// Verify alias
	cmd, ok := loaded.GetAlias("custom-alias")
	if !ok || cmd != "custom command" {
		t.Errorf("Alias not preserved in round-trip")
	}

	// Verify preset
	preset, ok := loaded.GetPreset("custom-preset")
	if !ok || preset.Command != "art" {
		t.Errorf("Preset not preserved in round-trip")
	}
}

// ============================================================================
// YAML/JSON Format Tests
// ============================================================================

func TestYAMLParsing_ValidYAML(t *testing.T) {
	yamlData := `
defaults:
  banner_font: "test-font"
  convert_width: 100
presets:
  test-preset:
    name: "Test"
    command: "test"
aliases:
  test-alias: "test cmd"
`
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(yamlData), cfg)

	if err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}

	if cfg.Defaults.BannerFont != "test-font" {
		t.Errorf("BannerFont = %q, want 'test-font'", cfg.Defaults.BannerFont)
	}
	if cfg.Defaults.ConvertWidth != 100 {
		t.Errorf("ConvertWidth = %d, want 100", cfg.Defaults.ConvertWidth)
	}
}

func TestJSONParsing_ValidJSON(t *testing.T) {
	jsonData := `{
		"defaults": {
			"banner_font": "json-font",
			"convert_width": 90
		},
		"presets": {
			"json-preset": {
				"name": "JSON Preset",
				"command": "test"
			}
		}
	}`

	cfg := &Config{}
	err := json.Unmarshal([]byte(jsonData), cfg)

	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if cfg.Defaults.BannerFont != "json-font" {
		t.Errorf("BannerFont = %q, want 'json-font'", cfg.Defaults.BannerFont)
	}
	if cfg.Defaults.ConvertWidth != 90 {
		t.Errorf("ConvertWidth = %d, want 90", cfg.Defaults.ConvertWidth)
	}
}

// ============================================================================
// Init Tests
// ============================================================================

func TestInit_CreatesDefaultConfig(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	err := Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check file was created
	configFile := filepath.Join(tmpDir, "moji", "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Init() didn't create config file")
	}
}

func TestInit_DoesNotOverwrite(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create initial config
	cfg := DefaultConfig()
	cfg.Defaults.BannerFont = "original"
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Init again
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Load and verify original value is preserved
	loaded, _ := Load()
	if loaded.Defaults.BannerFont != "original" {
		t.Error("Init() overwrote existing config")
	}
}

// ============================================================================
// Example Tests
// ============================================================================

func TestExample_ReturnsValidYAML(t *testing.T) {
	example := Example()

	if example == "" {
		t.Error("Example() returned empty string")
	}

	// Verify it contains expected fields
	expectedFields := []string{"defaults", "presets", "aliases"}
	for _, field := range expectedFields {
		if !contains(example, field) {
			t.Errorf("Example() missing field: %s", field)
		}
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestEmptyConfig(t *testing.T) {
	cfg := &Config{}

	// Should handle empty config
	_, ok := cfg.GetPreset("any")
	if ok {
		t.Error("GetPreset on empty config should return false")
	}

	alias, ok := cfg.GetAlias("any")
	if ok {
		t.Error("GetAlias on empty config should return false")
	}

	if alias != "" {
		t.Errorf("GetAlias returned %q, want empty", alias)
	}
}

func TestConfigWithOnlyDefaults(t *testing.T) {
	cfg := &Config{
		Defaults: Defaults{
			BannerFont:     "only-font",
			ConvertWidth:   100,
			GradientTheme:  "test-theme",
		},
		Presets: make(map[string]Preset),
		Aliases: make(map[string]string),
	}

	if cfg.Defaults.BannerFont != "only-font" {
		t.Errorf("BannerFont = %q, want 'only-font'", cfg.Defaults.BannerFont)
	}

	presets := cfg.ListPresets()
	if len(presets) != 0 {
		t.Errorf("Expected no presets, got %d", len(presets))
	}
}

func TestMultipleAliases(t *testing.T) {
	cfg := &Config{Aliases: make(map[string]string)}

	aliases := map[string]string{
		"alias1": "command1",
		"alias2": "command2",
		"alias3": "command3",
	}

	for alias, cmd := range aliases {
		cfg.SetAlias(alias, cmd)
	}

	for alias, expectedCmd := range aliases {
		cmd, ok := cfg.GetAlias(alias)
		if !ok {
			t.Errorf("GetAlias(%q) returned false", alias)
		}
		if cmd != expectedCmd {
			t.Errorf("GetAlias(%q) = %q, want %q", alias, cmd, expectedCmd)
		}
	}
}

func TestMultiplePresets(t *testing.T) {
	cfg := DefaultConfig()

	// Add custom presets
	cfg.SetPreset("custom1", Preset{Name: "Custom 1", Command: "cmd1"})
	cfg.SetPreset("custom2", Preset{Name: "Custom 2", Command: "cmd2"})

	presets := cfg.ListPresets()
	if len(presets) < 5 {
		t.Errorf("Expected at least 5 presets (3 default + 2 custom), got %d", len(presets))
	}
}

// ============================================================================
// Preset Structure Tests
// ============================================================================

func TestPresetStructure_NeonBanner(t *testing.T) {
	cfg := DefaultConfig()
	preset, ok := cfg.GetPreset("neon-banner")

	if !ok {
		t.Fatal("neon-banner preset not found")
	}

	if preset.Name != "Neon Banner" {
		t.Errorf("Name = %q, want 'Neon Banner'", preset.Name)
	}
	if preset.Command != "banner" {
		t.Errorf("Command = %q, want 'banner'", preset.Command)
	}
	if preset.Args["gradient"] != "neon" {
		t.Errorf("gradient arg = %q, want 'neon'", preset.Args["gradient"])
	}
	if preset.Args["border"] != "double" {
		t.Errorf("border arg = %q, want 'double'", preset.Args["border"])
	}
}

func TestPresetStructure_Retro(t *testing.T) {
	cfg := DefaultConfig()
	preset, ok := cfg.GetPreset("retro")

	if !ok {
		t.Fatal("retro preset not found")
	}

	if preset.Name != "Retro Terminal" {
		t.Errorf("Name = %q, want 'Retro Terminal'", preset.Name)
	}
	if preset.Command != "banner" {
		t.Errorf("Command = %q, want 'banner'", preset.Command)
	}
	if preset.Args["gradient"] != "retro" {
		t.Errorf("gradient arg = %q, want 'retro'", preset.Args["gradient"])
	}
}

func TestPresetStructure_FireBanner(t *testing.T) {
	cfg := DefaultConfig()
	preset, ok := cfg.GetPreset("fire-banner")

	if !ok {
		t.Fatal("fire-banner preset not found")
	}

	if preset.Name != "Fire Banner" {
		t.Errorf("Name = %q, want 'Fire Banner'", preset.Name)
	}
	if preset.Command != "banner" {
		t.Errorf("Command = %q, want 'banner'", preset.Command)
	}
	if preset.Args["gradient"] != "fire" {
		t.Errorf("gradient arg = %q, want 'fire'", preset.Args["gradient"])
	}
	if preset.Args["font"] != "block" {
		t.Errorf("font arg = %q, want 'block'", preset.Args["font"])
	}
}

// ============================================================================
// Custom Preset Tests
// ============================================================================

func TestCustomPresetWithArgs(t *testing.T) {
	cfg := &Config{Presets: make(map[string]Preset)}

	customPreset := Preset{
		Name:        "Custom Banner",
		Description: "A custom banner preset",
		Command:     "banner",
		Args: map[string]string{
			"font":   "shadow",
			"border": "bold",
			"style":  "bold",
		},
	}

	cfg.SetPreset("custom-banner", customPreset)

	retrieved, ok := cfg.GetPreset("custom-banner")
	if !ok {
		t.Fatal("Custom preset not found")
	}

	// Verify all args are preserved
	expectedArgs := []string{"font", "border", "style"}
	for _, arg := range expectedArgs {
		if retrieved.Args[arg] == "" {
			t.Errorf("Arg %q not preserved in custom preset", arg)
		}
	}

	// Verify the command is preserved
	if retrieved.Command != "banner" {
		t.Errorf("Command = %q, want 'banner'", retrieved.Command)
	}
}

// ============================================================================
// Path Tests
// ============================================================================

func TestPathPreservation_ArtPaths(t *testing.T) {
	cfg := &Config{
		ArtPaths: []string{"/path/to/art1", "/path/to/art2"},
	}

	if len(cfg.ArtPaths) != 2 {
		t.Errorf("ArtPaths length = %d, want 2", len(cfg.ArtPaths))
	}
	if cfg.ArtPaths[0] != "/path/to/art1" {
		t.Errorf("ArtPaths[0] = %q, want '/path/to/art1'", cfg.ArtPaths[0])
	}
}

func TestPathPreservation_FontPaths(t *testing.T) {
	cfg := &Config{
		FontPaths: []string{"/path/to/fonts"},
	}

	if len(cfg.FontPaths) != 1 {
		t.Errorf("FontPaths length = %d, want 1", len(cfg.FontPaths))
	}
}

func TestPathPreservation_CowfilePaths(t *testing.T) {
	cfg := &Config{
		CowfilePaths: []string{"/path/to/cowfiles1", "/path/to/cowfiles2", "/path/to/cowfiles3"},
	}

	if len(cfg.CowfilePaths) != 3 {
		t.Errorf("CowfilePaths length = %d, want 3", len(cfg.CowfilePaths))
	}
}

// ============================================================================
// Additional Coverage Tests
// ============================================================================

func TestSave_VariousPathStates(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Test saving multiple times
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := DefaultConfig()
	cfg.Defaults.BannerFont = "first"

	if err := cfg.Save(); err != nil {
		t.Fatalf("First Save() failed: %v", err)
	}

	// Modify and save again
	cfg.Defaults.BannerFont = "second"
	if err := cfg.Save(); err != nil {
		t.Fatalf("Second Save() failed: %v", err)
	}

	loaded, _ := Load()
	if loaded.Defaults.BannerFont != "second" {
		t.Errorf("BannerFont = %q, want 'second'", loaded.Defaults.BannerFont)
	}
}

func TestSetPreset_WithNilPresets(t *testing.T) {
	cfg := &Config{}
	// Presets is nil initially

	cfg.SetPreset("new-preset", Preset{
		Name:    "New",
		Command: "test",
	})

	// Should initialize the map
	preset, ok := cfg.GetPreset("new-preset")
	if !ok {
		t.Error("SetPreset with nil presets didn't work")
	}
	if preset.Command != "test" {
		t.Errorf("Preset command = %q, want 'test'", preset.Command)
	}
}

func TestLegacyConfigPath_InvalidHome(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set invalid home that can't be accessed
	os.Setenv("HOME", "")

	path := LegacyConfigPath()
	// When home is not accessible, the function returns ""
	if path != "" {
		t.Errorf("LegacyConfigPath with invalid HOME = %q, want ''", path)
	}
}

func TestConfigPath_InvalidHome(t *testing.T) {
	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("HOME", originalHome)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Unset both
	os.Setenv("HOME", "")
	os.Setenv("XDG_CONFIG_HOME", "")

	path := ConfigPath()
	if path != "" {
		t.Errorf("ConfigPath with no home = %q, want ''", path)
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create malformed YAML file
	configDir := filepath.Join(tmpDir, "moji")
	os.MkdirAll(configDir, 0755)
	configFile := filepath.Join(configDir, "config.yaml")
	os.WriteFile(configFile, []byte("invalid: yaml: content:"), 0644)

	// Load should fail on malformed YAML
	_, err := Load()
	if err == nil {
		t.Error("Load() should fail on malformed YAML")
	}
}

func TestLoad_LegacyJSONFormat(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("HOME", originalHome)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("HOME", tmpDir)
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, "nonexistent-xdg"))

	// Create legacy JSON config
	legacyPath := filepath.Join(tmpDir, ".mojirc")
	jsonData := `{"defaults":{"banner_font":"json-legacy"},"presets":{},"aliases":{}}`
	os.WriteFile(legacyPath, []byte(jsonData), 0644)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() failed on JSON: %v", err)
	}

	if loaded.Defaults.BannerFont != "json-legacy" {
		t.Errorf("BannerFont = %q, want 'json-legacy'", loaded.Defaults.BannerFont)
	}
}

func TestLoad_MalformedLegacyFile(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("HOME", originalHome)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("HOME", tmpDir)
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, "nonexistent-xdg"))

	// Create malformed legacy file
	legacyPath := filepath.Join(tmpDir, ".mojirc")
	os.WriteFile(legacyPath, []byte("invalid: yaml: and: json: {broken"), 0644)

	// Load should fail
	_, err := Load()
	if err == nil {
		t.Error("Load() should fail on malformed legacy file")
	}
}

func TestInit_InvalidPath(t *testing.T) {
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Set to empty to test error handling
	os.Setenv("XDG_CONFIG_HOME", "")

	// When home is not accessible, ConfigPath returns ""
	// Init should handle this gracefully
	err := Init()
	if err != nil {
		// We expect this might fail, but at minimum shouldn't panic
		t.Logf("Init with invalid path returned error: %v", err)
	}
}

func TestAliasEdgeCases(t *testing.T) {
	cfg := &Config{Aliases: make(map[string]string)}

	// Test empty alias name
	cfg.SetAlias("", "command")
	cmd, ok := cfg.GetAlias("")
	if !ok || cmd != "command" {
		t.Error("Empty alias name should work")
	}

	// Test empty command
	cfg.SetAlias("alias", "")
	cmd, ok = cfg.GetAlias("alias")
	if !ok || cmd != "" {
		t.Error("Empty command should be preserved")
	}

	// Test special characters
	cfg.SetAlias("alias-with-dash", "cmd with spaces")
	cmd, ok = cfg.GetAlias("alias-with-dash")
	if !ok || cmd != "cmd with spaces" {
		t.Error("Alias with special chars should work")
	}
}

func TestPresetEdgeCases(t *testing.T) {
	cfg := &Config{Presets: make(map[string]Preset)}

	// Test preset with empty fields
	cfg.SetPreset("empty-preset", Preset{})
	preset, ok := cfg.GetPreset("empty-preset")
	if !ok {
		t.Error("Empty preset should be retrievable")
	}
	if preset.Name != "" {
		t.Errorf("Empty preset Name should be empty, got %q", preset.Name)
	}

	// Test preset with nil args
	cfg.SetPreset("nil-args-preset", Preset{
		Name:    "Nil Args",
		Command: "test",
		Args:    nil,
	})
	preset, _ = cfg.GetPreset("nil-args-preset")
	if preset.Args != nil {
		t.Error("Nil args should remain nil")
	}
}

func TestDefaultConfig_AllPresetsHaveCommand(t *testing.T) {
	cfg := DefaultConfig()

	for name, preset := range cfg.Presets {
		if preset.Command == "" {
			t.Errorf("Preset %q has empty Command", name)
		}
	}
}

func TestDefaultConfig_AllPresetsHaveDescription(t *testing.T) {
	cfg := DefaultConfig()

	for name, preset := range cfg.Presets {
		if preset.Description == "" {
			t.Errorf("Preset %q has empty Description", name)
		}
	}
}

func TestSave_RecreatesExistingDirectory(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create the directory first
	configDir := filepath.Join(tmpDir, "moji")
	os.MkdirAll(configDir, 0755)

	cfg := DefaultConfig()
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file exists
	configFile := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Save() didn't create file in existing directory")
	}
}
