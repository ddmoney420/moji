// Package config provides configuration management with YAML and JSON support.
//
// It manages moji's configuration including default settings, command presets, aliases, and
// custom paths for art, fonts, and cowfiles. Configurations are stored as YAML or JSON files
// and support both global and user-specific settings.
//
// Example usage:
//
//	cfg := config.Load()
//	preset := cfg.GetPreset("mypreset")
//	cfg.SetPreset("newpreset", presetData)
//	cfg.Save()
package config
