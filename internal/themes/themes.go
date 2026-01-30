package themes

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Theme represents a color theme with metadata
type Theme struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Colors      []string          `yaml:"colors"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// registry holds all available themes (built-in + user-defined)
var registry = make(map[string]*Theme)

// colorHexRegex validates hex color format (#RRGGBB)
var colorHexRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// Init initializes the theme system with built-in themes and auto-loads user themes
func Init() error {
	// Register built-in themes
	registerBuiltinThemes()

	// Auto-load user themes from config directory
	configDir, err := getConfigDir()
	if err != nil {
		// Config dir not available, continue with built-in themes only
		return nil
	}

	themesDir := filepath.Join(configDir, "themes")
	if _, err := os.Stat(themesDir); err != nil {
		// Directory doesn't exist yet, that's fine
		return nil
	}

	_, err = LoadThemesDir(themesDir)
	return err
}

// registerBuiltinThemes adds hardcoded themes to the registry
func registerBuiltinThemes() {
	builtinThemes := []Theme{
		{
			Name:        "rainbow",
			Description: "Classic rainbow gradient",
			Colors: []string{
				"#FF0000", "#FF7F00", "#FFFF00", "#00FF00",
				"#0000FF", "#4B0082", "#9400D3",
			},
		},
		{
			Name:        "neon",
			Description: "Bright neon colors",
			Colors: []string{
				"#FF00FF", "#00FFFF", "#FFFF00", "#FF0080", "#00FF80",
			},
		},
		{
			Name:        "fire",
			Description: "Yellow to red fire",
			Colors: []string{
				"#FFFF00", "#FFC800", "#FF8000", "#FF0000", "#800000",
			},
		},
		{
			Name:        "ice",
			Description: "White to blue ice",
			Colors: []string{
				"#FFFFFF", "#C8E6FF", "#64B4FF", "#0064FF", "#003296",
			},
		},
		{
			Name:        "matrix",
			Description: "Green matrix style",
			Colors: []string{
				"#003200", "#006400", "#00C800", "#00FF00", "#96FF96",
			},
		},
		{
			Name:        "sunset",
			Description: "Warm sunset colors",
			Colors: []string{
				"#FF6464", "#FF9632", "#FFC864", "#C86496", "#643296",
			},
		},
		{
			Name:        "ocean",
			Description: "Deep blue ocean",
			Colors: []string{
				"#003264", "#006496", "#0096C8", "#64C8E6", "#C8FFFF",
			},
		},
		{
			Name:        "c64",
			Description: "Commodore 64 palette",
			Colors: []string{
				"#403285", "#665AFF", "#86C741", "#FFF1E0",
			},
		},
		{
			Name:        "dracula",
			Description: "Dracula theme colors",
			Colors: []string{
				"#FF79C6", "#BD93F9", "#8BE9FD", "#50FA7B", "#FFB86C",
			},
		},
		{
			Name:        "vaporwave",
			Description: "80s vaporwave aesthetic",
			Colors: []string{
				"#FF71CE", "#B967FF", "#01CDFE", "#05FFA1", "#FFFB96",
			},
		},
		{
			Name:        "retro",
			Description: "Retro green terminal",
			Colors: []string{
				"#004000", "#008000", "#00C000", "#00FF00", "#80FF80",
			},
		},
		{
			Name:        "pastel",
			Description: "Soft pastel colors",
			Colors: []string{
				"#FFB3BA", "#FFDFBA", "#FFFFBA", "#BAFFC9", "#BAE1FF",
			},
		},
	}

	for i := range builtinThemes {
		registry[builtinThemes[i].Name] = &builtinThemes[i]
	}
}

// LoadTheme loads a theme from a YAML file
func LoadTheme(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	var theme Theme
	if err := yaml.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme YAML: %w", err)
	}

	if err := validateTheme(&theme); err != nil {
		return nil, err
	}

	return &theme, nil
}

// LoadThemesDir loads all themes from a directory
func LoadThemesDir(dir string) ([]Theme, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	var loaded []Theme
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		theme, err := LoadTheme(path)
		if err != nil {
			// Log warning but continue loading other themes
			fmt.Fprintf(os.Stderr, "Warning: failed to load theme %s: %v\n", entry.Name(), err)
			continue
		}

		if err := RegisterTheme(theme); err != nil {
			// Log warning but continue loading other themes
			fmt.Fprintf(os.Stderr, "Warning: failed to register theme %s: %v\n", entry.Name(), err)
			continue
		}

		loaded = append(loaded, *theme)
	}

	return loaded, nil
}

// SaveTheme saves a theme to a YAML file
func SaveTheme(theme *Theme, path string) error {
	if err := validateTheme(theme); err != nil {
		return err
	}

	data, err := yaml.Marshal(theme)
	if err != nil {
		return fmt.Errorf("failed to marshal theme to YAML: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create theme directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	return nil
}

// RegisterTheme adds a theme to the registry
func RegisterTheme(theme *Theme) error {
	if err := validateTheme(theme); err != nil {
		return err
	}
	registry[theme.Name] = theme
	return nil
}

// GetTheme retrieves a theme by name
func GetTheme(name string) (*Theme, error) {
	theme, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("theme %q not found", name)
	}
	return theme, nil
}

// ListThemes returns all available theme names
func ListThemes() []string {
	var names []string
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// validateTheme checks theme structure and color validity
func validateTheme(theme *Theme) error {
	if theme == nil {
		return fmt.Errorf("theme cannot be nil")
	}

	if theme.Name == "" {
		return fmt.Errorf("theme name is required")
	}

	if len(theme.Colors) == 0 {
		return fmt.Errorf("theme must have at least one color")
	}

	for i, colorStr := range theme.Colors {
		if !colorHexRegex.MatchString(colorStr) {
			return fmt.Errorf("invalid color format at index %d: %q (expected #RRGGBB)", i, colorStr)
		}
	}

	return nil
}

// getConfigDir returns the XDG config directory for moji
func getConfigDir() (string, error) {
	// Try XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "moji"), nil
	}

	// Fall back to ~/.config/moji
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "moji"), nil
}

// HexToRGB converts a hex color string to RGB components
func HexToRGB(hexStr string) (uint8, uint8, uint8, error) {
	if !colorHexRegex.MatchString(hexStr) {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %q", hexStr)
	}

	// Remove the '#' prefix
	hexStr = hexStr[1:]

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, 0, 0, err
	}

	return bytes[0], bytes[1], bytes[2], nil
}

// RGBToHex converts RGB components to a hex color string
func RGBToHex(r, g, b uint8) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
