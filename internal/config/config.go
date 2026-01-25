package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the moji configuration
type Config struct {
	// Default settings
	Defaults Defaults `json:"defaults" yaml:"defaults"`

	// Presets for quick access
	Presets map[string]Preset `json:"presets" yaml:"presets"`

	// Custom art database paths
	ArtPaths []string `json:"art_paths" yaml:"art_paths"`

	// Custom font paths
	FontPaths []string `json:"font_paths" yaml:"font_paths"`

	// Custom cowfile paths
	CowfilePaths []string `json:"cowfile_paths" yaml:"cowfile_paths"`

	// Aliases for commands
	Aliases map[string]string `json:"aliases" yaml:"aliases"`
}

// Defaults holds default values for commands
type Defaults struct {
	// Banner defaults
	BannerFont   string `json:"banner_font" yaml:"banner_font"`
	BannerStyle  string `json:"banner_style" yaml:"banner_style"`
	BannerBorder string `json:"banner_border" yaml:"banner_border"`

	// Convert defaults
	ConvertWidth   int    `json:"convert_width" yaml:"convert_width"`
	ConvertCharset string `json:"convert_charset" yaml:"convert_charset"`
	ConvertDither  string `json:"convert_dither" yaml:"convert_dither"`

	// Gradient defaults
	GradientTheme string `json:"gradient_theme" yaml:"gradient_theme"`
	GradientMode  string `json:"gradient_mode" yaml:"gradient_mode"`

	// Speech defaults
	BubbleStyle string `json:"bubble_style" yaml:"bubble_style"`
	BubbleWidth int    `json:"bubble_width" yaml:"bubble_width"`

	// QR defaults
	QRCharset string `json:"qr_charset" yaml:"qr_charset"`

	// Pattern defaults
	PatternBorder  string `json:"pattern_border" yaml:"pattern_border"`
	PatternDivider string `json:"pattern_divider" yaml:"pattern_divider"`

	// Output defaults
	CopyToClipboard bool `json:"copy_to_clipboard" yaml:"copy_to_clipboard"`
	JSONOutput      bool `json:"json_output" yaml:"json_output"`
}

// Preset represents a saved configuration preset
type Preset struct {
	Name        string            `json:"name" yaml:"name"`
	Description string            `json:"description" yaml:"description"`
	Command     string            `json:"command" yaml:"command"`
	Args        map[string]string `json:"args" yaml:"args"`
}

// Default configuration
func DefaultConfig() *Config {
	return &Config{
		Defaults: Defaults{
			BannerFont:     "standard",
			BannerStyle:    "none",
			BannerBorder:   "none",
			ConvertWidth:   80,
			ConvertCharset: "standard",
			ConvertDither:  "none",
			GradientTheme:  "rainbow",
			GradientMode:   "horizontal",
			BubbleStyle:    "round",
			BubbleWidth:    40,
			QRCharset:      "blocks",
			PatternBorder:  "single",
			PatternDivider: "single",
		},
		Presets: map[string]Preset{
			"neon-banner": {
				Name:        "Neon Banner",
				Description: "Colorful neon-style banner",
				Command:     "banner",
				Args: map[string]string{
					"gradient": "neon",
					"border":   "double",
				},
			},
			"retro": {
				Name:        "Retro Terminal",
				Description: "Old-school green terminal look",
				Command:     "banner",
				Args: map[string]string{
					"gradient": "retro",
				},
			},
			"fire-banner": {
				Name:        "Fire Banner",
				Description: "Hot fire-colored banner",
				Command:     "banner",
				Args: map[string]string{
					"gradient": "fire",
					"font":     "block",
				},
			},
		},
		ArtPaths:     []string{},
		FontPaths:    []string{},
		CowfilePaths: []string{},
		Aliases:      map[string]string{},
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	// Check XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "moji", "config.yaml")
	}

	// Fall back to ~/.config/moji/config.yaml
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "moji", "config.yaml")
}

// LegacyConfigPath returns the legacy config path (~/.mojirc)
func LegacyConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".mojirc")
}

// Load loads configuration from file
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try modern path first
	configPath := ConfigPath()
	if configPath != "" {
		if data, err := os.ReadFile(configPath); err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
			return cfg, nil
		}
	}

	// Try legacy path
	legacyPath := LegacyConfigPath()
	if legacyPath != "" {
		if data, err := os.ReadFile(legacyPath); err == nil {
			// Try YAML first
			if err := yaml.Unmarshal(data, cfg); err != nil {
				// Try JSON
				if err := json.Unmarshal(data, cfg); err != nil {
					return nil, err
				}
			}
			return cfg, nil
		}
	}

	// Return default config if no file found
	return cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := ConfigPath()
	if configPath == "" {
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetPreset returns a preset by name
func (c *Config) GetPreset(name string) (Preset, bool) {
	p, ok := c.Presets[name]
	return p, ok
}

// SetPreset saves a preset
func (c *Config) SetPreset(name string, preset Preset) {
	if c.Presets == nil {
		c.Presets = make(map[string]Preset)
	}
	c.Presets[name] = preset
}

// ListPresets returns all preset names
func (c *Config) ListPresets() []string {
	var names []string
	for k := range c.Presets {
		names = append(names, k)
	}
	return names
}

// GetAlias returns command for an alias
func (c *Config) GetAlias(alias string) (string, bool) {
	cmd, ok := c.Aliases[alias]
	return cmd, ok
}

// SetAlias creates an alias
func (c *Config) SetAlias(alias, command string) {
	if c.Aliases == nil {
		c.Aliases = make(map[string]string)
	}
	c.Aliases[alias] = command
}

// Init creates a default config file if it doesn't exist
func Init() error {
	configPath := ConfigPath()
	if configPath == "" {
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // Already exists
	}

	// Create default config
	cfg := DefaultConfig()
	return cfg.Save()
}

// Example returns an example config as YAML string
func Example() string {
	cfg := DefaultConfig()
	data, _ := yaml.Marshal(cfg)
	return string(data)
}
