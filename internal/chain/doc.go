// Package chain provides a pipeline for sequencing and chaining text transformations.
//
// It allows composition of multiple effects (text effects, speech bubbles, borders, gradients)
// into a single transformation pipeline, enabling complex visual effects by combining simpler
// operations in sequence.
//
// Example usage:
//
//	chain := chain.New().
//		AddEffect(effects.Flip).
//		AddGradient(gradient.Rainbow).
//		AddBorder(patterns.SolidBorder).
//		Apply(text)
//	chain.ApplyToArt(asciiArt)
package chain
