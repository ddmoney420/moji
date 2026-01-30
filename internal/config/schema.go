package config

import (
	"encoding/json"
)

// JSONSchemaType represents a JSON Schema type definition
type JSONSchemaType struct {
	Schema      string                 `json:"$schema"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Properties  map[string]interface{} `json:"properties"`
	Required    []string               `json:"required"`
	Examples    []interface{}          `json:"examples,omitempty"`
}

// GenerateSchema returns a JSON Schema for the moji configuration format
func GenerateSchema() string {
	schema := JSONSchemaType{
		Schema:      "http://json-schema.org/draft-07/schema#",
		Title:       "Moji Configuration Schema",
		Description: "Configuration schema for moji CLI tool",
		Type:        "object",
		Required:    []string{"defaults"},
		Properties: map[string]interface{}{
			"defaults": map[string]interface{}{
				"type":        "object",
				"title":       "Default Settings",
				"description": "Default values for moji commands",
				"properties": map[string]interface{}{
					"banner_font": map[string]interface{}{
						"type":        "string",
						"title":       "Banner Font",
						"description": "Default font for banner command (standard, shadow, slant, block, etc.)",
						"default":     "standard",
						"examples":    []string{"standard", "shadow", "slant", "block"},
					},
					"banner_style": map[string]interface{}{
						"type":        "string",
						"title":       "Banner Style",
						"description": "Default style for banner output (none, bold, italic, etc.)",
						"default":     "none",
						"examples":    []string{"none", "bold", "italic"},
					},
					"banner_border": map[string]interface{}{
						"type":        "string",
						"title":       "Banner Border",
						"description": "Default border for banner (none, single, double, round, bold, stars, hash, etc.)",
						"default":     "none",
						"enum":        []string{"none", "single", "double", "round", "bold", "stars", "hash"},
					},
					"convert_width": map[string]interface{}{
						"type":        "integer",
						"title":       "Convert Width",
						"description": "Default width for ASCII art conversion in characters",
						"default":     80,
						"minimum":     10,
						"maximum":     1000,
					},
					"convert_charset": map[string]interface{}{
						"type":        "string",
						"title":       "Convert Charset",
						"description": "Default character set for ASCII conversion (standard, extended, blocks, etc.)",
						"default":     "standard",
						"examples":    []string{"standard", "extended", "blocks"},
					},
					"convert_dither": map[string]interface{}{
						"type":        "string",
						"title":       "Convert Dither",
						"description": "Default dithering algorithm (none, floyd, bayer, etc.)",
						"default":     "none",
						"examples":    []string{"none", "floyd", "bayer"},
					},
					"gradient_theme": map[string]interface{}{
						"type":        "string",
						"title":       "Gradient Theme",
						"description": "Default gradient theme (rainbow, neon, retro, fire, ocean, etc.)",
						"default":     "rainbow",
						"examples":    []string{"rainbow", "neon", "retro", "fire", "ocean"},
					},
					"gradient_mode": map[string]interface{}{
						"type":        "string",
						"title":       "Gradient Mode",
						"description": "Direction of gradient application (horizontal, vertical, diagonal)",
						"default":     "horizontal",
						"enum":        []string{"horizontal", "vertical", "diagonal"},
					},
					"bubble_style": map[string]interface{}{
						"type":        "string",
						"title":       "Bubble Style",
						"description": "Default style for speech bubbles (round, square, cloud, etc.)",
						"default":     "round",
						"examples":    []string{"round", "square", "cloud", "think"},
					},
					"bubble_width": map[string]interface{}{
						"type":        "integer",
						"title":       "Bubble Width",
						"description": "Default width for speech bubbles in characters",
						"default":     40,
						"minimum":     10,
						"maximum":     200,
					},
					"qr_charset": map[string]interface{}{
						"type":        "string",
						"title":       "QR Code Charset",
						"description": "Character set for QR code output (blocks, ascii, etc.)",
						"default":     "blocks",
						"examples":    []string{"blocks", "ascii"},
					},
					"pattern_border": map[string]interface{}{
						"type":        "string",
						"title":       "Pattern Border",
						"description": "Default border style for patterns (single, double, etc.)",
						"default":     "single",
						"examples":    []string{"single", "double"},
					},
					"pattern_divider": map[string]interface{}{
						"type":        "string",
						"title":       "Pattern Divider",
						"description": "Default divider style for patterns (single, double, etc.)",
						"default":     "single",
						"examples":    []string{"single", "double"},
					},
					"copy_to_clipboard": map[string]interface{}{
						"type":        "boolean",
						"title":       "Copy to Clipboard",
						"description": "Whether to automatically copy output to clipboard",
						"default":     false,
					},
					"json_output": map[string]interface{}{
						"type":        "boolean",
						"title":       "JSON Output",
						"description": "Whether to output in JSON format by default",
						"default":     false,
					},
				},
			},
			"presets": map[string]interface{}{
				"type":        "object",
				"title":       "Presets",
				"description": "Named configuration presets for quick access",
				"additionalProperties": map[string]interface{}{
					"type":  "object",
					"title": "Preset",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"title":       "Name",
							"description": "Display name for the preset",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"title":       "Description",
							"description": "Description of what this preset does",
						},
						"command": map[string]interface{}{
							"type":        "string",
							"title":       "Command",
							"description": "The moji command this preset is for (banner, convert, art, etc.)",
						},
						"args": map[string]interface{}{
							"type":        "object",
							"title":       "Arguments",
							"description": "Command arguments as key-value pairs",
							"additionalProperties": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"command"},
				},
			},
			"art_paths": map[string]interface{}{
				"type":        "array",
				"title":       "Custom Art Paths",
				"description": "Additional directories to search for ASCII art files",
				"items": map[string]interface{}{
					"type": "string",
				},
				"examples": [][]string{{"/usr/local/share/moji/art", "/home/user/.moji/art"}},
			},
			"font_paths": map[string]interface{}{
				"type":        "array",
				"title":       "Custom Font Paths",
				"description": "Additional directories to search for font files",
				"items": map[string]interface{}{
					"type": "string",
				},
				"examples": [][]string{{"/usr/local/share/fonts", "/home/user/.moji/fonts"}},
			},
			"cowfile_paths": map[string]interface{}{
				"type":        "array",
				"title":       "Cowfile Paths",
				"description": "Additional directories to search for cowfile art",
				"items": map[string]interface{}{
					"type": "string",
				},
				"examples": [][]string{{"/usr/local/share/cowsay/cows", "/home/user/.moji/cows"}},
			},
			"aliases": map[string]interface{}{
				"type":        "object",
				"title":       "Command Aliases",
				"description": "Custom aliases for frequently used commands",
				"additionalProperties": map[string]interface{}{
					"type":        "string",
					"title":       "Alias Command",
					"description": "The command this alias expands to",
				},
				"examples": []map[string]string{{"my-banner": "banner --font shadow --border bold"}},
			},
		},
		Examples: []interface{}{
			map[string]interface{}{
				"defaults": map[string]interface{}{
					"banner_font":     "shadow",
					"banner_border":   "bold",
					"convert_width":   100,
					"gradient_theme":  "neon",
				},
				"presets": map[string]interface{}{
					"neon-banner": map[string]interface{}{
						"name":        "Neon Banner",
						"description": "Colorful neon-style banner",
						"command":     "banner",
						"args": map[string]string{
							"gradient": "neon",
							"border":   "double",
						},
					},
				},
				"aliases": map[string]string{
					"my-art": "art mountains",
				},
			},
		},
	}

	data, _ := json.MarshalIndent(schema, "", "  ")
	return string(data)
}
