// Package halfblock provides half-block ASCII rendering for higher vertical resolution.
//
// It uses half-block characters (▀▄) to achieve 2x vertical resolution compared to standard
// ASCII, with support for grayscale and color rendering using bilinear sampling.
//
// Example usage:
//
//	result := halfblock.RenderGrayscale(img)
//	result := halfblock.RenderColor(img)
//	result := halfblock.RenderWithCharset(img, charset)
package halfblock
