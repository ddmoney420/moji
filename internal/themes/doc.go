// Package themes provides a theme system for color palettes.
//
// It manages 12+ built-in themes (rainbow, neon, dracula, etc.) and supports custom theme
// registration from YAML files. The package maintains a theme registry for consistent color
// application across the application.
//
// Example usage:
//
//	theme := themes.LoadTheme("rainbow")
//	theme := themes.GetTheme("neon")
//	themes.RegisterTheme("custom", customTheme)
//	themes.LoadFromFile("custom.yaml")
package themes
