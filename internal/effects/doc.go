// Package effects provides Unicode text transformation effects.
//
// It offers 10+ text transformation effects including upside-down text (Flip/Reverse), mirroring,
// Zalgo (corrupted text), bubble text, and various Unicode styles like bold, italic, small caps,
// fullwidth, script, and Fraktur.
//
// Example usage:
//
//	flipped := effects.Flip("Hello")
//	zalgo := effects.Zalgo("Creepy")
//	bubble := effects.Bubble("Message")
//	bold := effects.Bold("Bold")
//	fullwidth := effects.Fullwidth("Text")
package effects
