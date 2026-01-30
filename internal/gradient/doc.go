// Package gradient provides color gradient application to text.
//
// It implements 12+ color themes (rainbow, neon, fire, ocean, dracula, vaporwave, etc.)
// that can be applied to text with horizontal, vertical, or diagonal gradient modes.
//
// Example usage:
//
//	result := gradient.Apply("Text", gradient.Rainbow, gradient.Horizontal)
//	result := gradient.Apply("Text", gradient.Neon, gradient.Vertical)
//	themes := gradient.ListThemes()
package gradient
