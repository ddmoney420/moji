// Package dither provides image dithering algorithms for ASCII conversion.
//
// It implements multiple dithering techniques (Floyd-Steinberg, Bayer, Atkinson, Sierra, etc.)
// that convert images to grayscale while distributing quantization errors for better visual
// appearance in ASCII art.
//
// Example usage:
//
//	grayscale := dither.FloydSteinberg(img)
//	grayscale := dither.Bayer(img)
//	grayscale := dither.Atkinson(img)
package dither
