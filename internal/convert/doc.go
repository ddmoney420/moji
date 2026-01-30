// Package convert provides image-to-ASCII conversion with multiple charsets and advanced features.
//
// It converts images to ASCII art with support for various character sets, dithering algorithms,
// edge detection, color preservation, and parallel processing. The package handles files, URLs,
// and in-memory image data.
//
// Example usage:
//
//	result := convert.FromFile("image.jpg", opts)
//	result := convert.FromURL("https://example.com/image.png", opts)
//	result := convert.FromImage(img, opts)
//	result.Print()
package convert
