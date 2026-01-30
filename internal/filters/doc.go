// Package filters provides text filtering effects with 20+ artistic styles.
//
// It includes filters such as rainbow coloring, metal effects, fire, ice, glitch, matrix,
// and many others that can be combined through chaining.
//
// Example usage:
//
//	result := filters.Rainbow("Text")
//	result := filters.Metal("Text")
//	result := filters.Fire("Text")
//	result := filters.Glitch("Text")
//	chain := filters.Chain(text, filters.Rainbow, filters.Glitch)
package filters
