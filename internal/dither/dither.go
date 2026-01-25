package dither

import (
	"image"
	"image/color"
	"math"
)

// Algorithm represents a dithering algorithm
type Algorithm string

const (
	None           Algorithm = "none"
	FloydSteinberg Algorithm = "floyd-steinberg"
	Bayer2x2       Algorithm = "bayer2x2"
	Bayer4x4       Algorithm = "bayer4x4"
	Bayer8x8       Algorithm = "bayer8x8"
	Atkinson       Algorithm = "atkinson"
	Sierra         Algorithm = "sierra"
	SierraLite     Algorithm = "sierra-lite"
	Stucki         Algorithm = "stucki"
	Burkes         Algorithm = "burkes"
	JarvisJudice   Algorithm = "jarvis"
)

// Bayer matrices for ordered dithering
var bayer2x2 = [][]float64{
	{0.0 / 4.0, 2.0 / 4.0},
	{3.0 / 4.0, 1.0 / 4.0},
}

var bayer4x4 = [][]float64{
	{0.0 / 16.0, 8.0 / 16.0, 2.0 / 16.0, 10.0 / 16.0},
	{12.0 / 16.0, 4.0 / 16.0, 14.0 / 16.0, 6.0 / 16.0},
	{3.0 / 16.0, 11.0 / 16.0, 1.0 / 16.0, 9.0 / 16.0},
	{15.0 / 16.0, 7.0 / 16.0, 13.0 / 16.0, 5.0 / 16.0},
}

var bayer8x8 = [][]float64{
	{0.0 / 64.0, 32.0 / 64.0, 8.0 / 64.0, 40.0 / 64.0, 2.0 / 64.0, 34.0 / 64.0, 10.0 / 64.0, 42.0 / 64.0},
	{48.0 / 64.0, 16.0 / 64.0, 56.0 / 64.0, 24.0 / 64.0, 50.0 / 64.0, 18.0 / 64.0, 58.0 / 64.0, 26.0 / 64.0},
	{12.0 / 64.0, 44.0 / 64.0, 4.0 / 64.0, 36.0 / 64.0, 14.0 / 64.0, 46.0 / 64.0, 6.0 / 64.0, 38.0 / 64.0},
	{60.0 / 64.0, 28.0 / 64.0, 52.0 / 64.0, 20.0 / 64.0, 62.0 / 64.0, 30.0 / 64.0, 54.0 / 64.0, 22.0 / 64.0},
	{3.0 / 64.0, 35.0 / 64.0, 11.0 / 64.0, 43.0 / 64.0, 1.0 / 64.0, 33.0 / 64.0, 9.0 / 64.0, 41.0 / 64.0},
	{51.0 / 64.0, 19.0 / 64.0, 59.0 / 64.0, 27.0 / 64.0, 49.0 / 64.0, 17.0 / 64.0, 57.0 / 64.0, 25.0 / 64.0},
	{15.0 / 64.0, 47.0 / 64.0, 7.0 / 64.0, 39.0 / 64.0, 13.0 / 64.0, 45.0 / 64.0, 5.0 / 64.0, 37.0 / 64.0},
	{63.0 / 64.0, 31.0 / 64.0, 55.0 / 64.0, 23.0 / 64.0, 61.0 / 64.0, 29.0 / 64.0, 53.0 / 64.0, 21.0 / 64.0},
}

// Apply applies dithering algorithm to an image
func Apply(img image.Image, algo Algorithm) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	// Convert to grayscale first
	gray := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x+bounds.Min.X, y+bounds.Min.Y)
			gray.Set(x, y, color.GrayModel.Convert(c).(color.Gray))
		}
	}

	switch algo {
	case FloydSteinberg:
		return floydSteinberg(gray)
	case Bayer2x2:
		return orderedDither(gray, bayer2x2)
	case Bayer4x4:
		return orderedDither(gray, bayer4x4)
	case Bayer8x8:
		return orderedDither(gray, bayer8x8)
	case Atkinson:
		return atkinson(gray)
	case Sierra:
		return sierra(gray)
	case SierraLite:
		return sierraLite(gray)
	case Stucki:
		return stucki(gray)
	case Burkes:
		return burkes(gray)
	case JarvisJudice:
		return jarvisJudice(gray)
	default:
		return gray
	}
}

// floydSteinberg implements Floyd-Steinberg dithering
func floydSteinberg(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	// Copy image
	errors := make([][]float64, height)
	for y := 0; y < height; y++ {
		errors[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			errors[y][x] = float64(img.GrayAt(x, y).Y)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := errors[y][x]
			newPixel := 0.0
			if oldPixel > 127 {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: uint8(newPixel)})

			err := oldPixel - newPixel

			// Distribute error
			if x+1 < width {
				errors[y][x+1] += err * 7.0 / 16.0
			}
			if y+1 < height {
				if x > 0 {
					errors[y+1][x-1] += err * 3.0 / 16.0
				}
				errors[y+1][x] += err * 5.0 / 16.0
				if x+1 < width {
					errors[y+1][x+1] += err * 1.0 / 16.0
				}
			}
		}
	}

	return result
}

// orderedDither applies ordered/Bayer dithering
func orderedDither(img *image.Gray, matrix [][]float64) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	matrixSize := len(matrix)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := float64(img.GrayAt(x, y).Y) / 255.0
			threshold := matrix[y%matrixSize][x%matrixSize]

			newPixel := uint8(0)
			if oldPixel > threshold {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: newPixel})
		}
	}

	return result
}

// atkinson implements Atkinson dithering (used by early Macs)
func atkinson(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	errors := make([][]float64, height)
	for y := 0; y < height; y++ {
		errors[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			errors[y][x] = float64(img.GrayAt(x, y).Y)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := errors[y][x]
			newPixel := 0.0
			if oldPixel > 127 {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: uint8(newPixel)})

			// Atkinson uses 1/8 and only distributes 6/8 of the error
			err := (oldPixel - newPixel) / 8.0

			if x+1 < width {
				errors[y][x+1] += err
			}
			if x+2 < width {
				errors[y][x+2] += err
			}
			if y+1 < height {
				if x > 0 {
					errors[y+1][x-1] += err
				}
				errors[y+1][x] += err
				if x+1 < width {
					errors[y+1][x+1] += err
				}
			}
			if y+2 < height {
				errors[y+2][x] += err
			}
		}
	}

	return result
}

// sierra implements Sierra dithering
func sierra(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	errors := make([][]float64, height)
	for y := 0; y < height; y++ {
		errors[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			errors[y][x] = float64(img.GrayAt(x, y).Y)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := errors[y][x]
			newPixel := 0.0
			if oldPixel > 127 {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: uint8(newPixel)})

			err := oldPixel - newPixel

			// Sierra coefficients (divided by 32)
			if x+1 < width {
				errors[y][x+1] += err * 5.0 / 32.0
			}
			if x+2 < width {
				errors[y][x+2] += err * 3.0 / 32.0
			}
			if y+1 < height {
				if x > 1 {
					errors[y+1][x-2] += err * 2.0 / 32.0
				}
				if x > 0 {
					errors[y+1][x-1] += err * 4.0 / 32.0
				}
				errors[y+1][x] += err * 5.0 / 32.0
				if x+1 < width {
					errors[y+1][x+1] += err * 4.0 / 32.0
				}
				if x+2 < width {
					errors[y+1][x+2] += err * 2.0 / 32.0
				}
			}
			if y+2 < height {
				if x > 0 {
					errors[y+2][x-1] += err * 2.0 / 32.0
				}
				errors[y+2][x] += err * 3.0 / 32.0
				if x+1 < width {
					errors[y+2][x+1] += err * 2.0 / 32.0
				}
			}
		}
	}

	return result
}

// sierraLite implements Sierra Lite (Two-Row Sierra) dithering
func sierraLite(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	errors := make([][]float64, height)
	for y := 0; y < height; y++ {
		errors[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			errors[y][x] = float64(img.GrayAt(x, y).Y)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := errors[y][x]
			newPixel := 0.0
			if oldPixel > 127 {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: uint8(newPixel)})

			err := oldPixel - newPixel

			if x+1 < width {
				errors[y][x+1] += err * 2.0 / 4.0
			}
			if y+1 < height {
				if x > 0 {
					errors[y+1][x-1] += err * 1.0 / 4.0
				}
				errors[y+1][x] += err * 1.0 / 4.0
			}
		}
	}

	return result
}

// stucki implements Stucki dithering
func stucki(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	errors := make([][]float64, height)
	for y := 0; y < height; y++ {
		errors[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			errors[y][x] = float64(img.GrayAt(x, y).Y)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := errors[y][x]
			newPixel := 0.0
			if oldPixel > 127 {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: uint8(newPixel)})

			err := oldPixel - newPixel

			// Stucki coefficients (divided by 42)
			if x+1 < width {
				errors[y][x+1] += err * 8.0 / 42.0
			}
			if x+2 < width {
				errors[y][x+2] += err * 4.0 / 42.0
			}
			if y+1 < height {
				if x > 1 {
					errors[y+1][x-2] += err * 2.0 / 42.0
				}
				if x > 0 {
					errors[y+1][x-1] += err * 4.0 / 42.0
				}
				errors[y+1][x] += err * 8.0 / 42.0
				if x+1 < width {
					errors[y+1][x+1] += err * 4.0 / 42.0
				}
				if x+2 < width {
					errors[y+1][x+2] += err * 2.0 / 42.0
				}
			}
			if y+2 < height {
				if x > 1 {
					errors[y+2][x-2] += err * 1.0 / 42.0
				}
				if x > 0 {
					errors[y+2][x-1] += err * 2.0 / 42.0
				}
				errors[y+2][x] += err * 4.0 / 42.0
				if x+1 < width {
					errors[y+2][x+1] += err * 2.0 / 42.0
				}
				if x+2 < width {
					errors[y+2][x+2] += err * 1.0 / 42.0
				}
			}
		}
	}

	return result
}

// burkes implements Burkes dithering
func burkes(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	errors := make([][]float64, height)
	for y := 0; y < height; y++ {
		errors[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			errors[y][x] = float64(img.GrayAt(x, y).Y)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := errors[y][x]
			newPixel := 0.0
			if oldPixel > 127 {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: uint8(newPixel)})

			err := oldPixel - newPixel

			// Burkes coefficients (divided by 32)
			if x+1 < width {
				errors[y][x+1] += err * 8.0 / 32.0
			}
			if x+2 < width {
				errors[y][x+2] += err * 4.0 / 32.0
			}
			if y+1 < height {
				if x > 1 {
					errors[y+1][x-2] += err * 2.0 / 32.0
				}
				if x > 0 {
					errors[y+1][x-1] += err * 4.0 / 32.0
				}
				errors[y+1][x] += err * 8.0 / 32.0
				if x+1 < width {
					errors[y+1][x+1] += err * 4.0 / 32.0
				}
				if x+2 < width {
					errors[y+1][x+2] += err * 2.0 / 32.0
				}
			}
		}
	}

	return result
}

// jarvisJudice implements Jarvis-Judice-Ninke dithering
func jarvisJudice(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)

	errors := make([][]float64, height)
	for y := 0; y < height; y++ {
		errors[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			errors[y][x] = float64(img.GrayAt(x, y).Y)
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := errors[y][x]
			newPixel := 0.0
			if oldPixel > 127 {
				newPixel = 255
			}
			result.SetGray(x, y, color.Gray{Y: uint8(newPixel)})

			err := oldPixel - newPixel

			// Jarvis-Judice-Ninke coefficients (divided by 48)
			if x+1 < width {
				errors[y][x+1] += err * 7.0 / 48.0
			}
			if x+2 < width {
				errors[y][x+2] += err * 5.0 / 48.0
			}
			if y+1 < height {
				if x > 1 {
					errors[y+1][x-2] += err * 3.0 / 48.0
				}
				if x > 0 {
					errors[y+1][x-1] += err * 5.0 / 48.0
				}
				errors[y+1][x] += err * 7.0 / 48.0
				if x+1 < width {
					errors[y+1][x+1] += err * 5.0 / 48.0
				}
				if x+2 < width {
					errors[y+1][x+2] += err * 3.0 / 48.0
				}
			}
			if y+2 < height {
				if x > 1 {
					errors[y+2][x-2] += err * 1.0 / 48.0
				}
				if x > 0 {
					errors[y+2][x-1] += err * 3.0 / 48.0
				}
				errors[y+2][x] += err * 5.0 / 48.0
				if x+1 < width {
					errors[y+2][x+1] += err * 3.0 / 48.0
				}
				if x+2 < width {
					errors[y+2][x+2] += err * 1.0 / 48.0
				}
			}
		}
	}

	return result
}

// ToGrayscale converts brightness values to grayscale levels
func ToGrayscale(img *image.Gray, levels int) *image.Gray {
	if levels < 2 {
		levels = 2
	}
	if levels > 256 {
		levels = 256
	}

	bounds := img.Bounds()
	result := image.NewGray(bounds)

	step := 256.0 / float64(levels)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldVal := float64(img.GrayAt(x, y).Y)
			level := int(oldVal / step)
			if level >= levels {
				level = levels - 1
			}
			newVal := uint8(float64(level) * step)
			result.SetGray(x, y, color.Gray{Y: newVal})
		}
	}

	return result
}

// ListAlgorithms returns available dithering algorithms
func ListAlgorithms() []string {
	return []string{
		"none",
		"floyd-steinberg",
		"bayer2x2",
		"bayer4x4",
		"bayer8x8",
		"atkinson",
		"sierra",
		"sierra-lite",
		"stucki",
		"burkes",
		"jarvis",
	}
}

// GetAlgorithm returns algorithm from string
func GetAlgorithm(name string) Algorithm {
	switch name {
	case "floyd-steinberg", "fs":
		return FloydSteinberg
	case "bayer2x2", "bayer2":
		return Bayer2x2
	case "bayer4x4", "bayer4", "bayer":
		return Bayer4x4
	case "bayer8x8", "bayer8":
		return Bayer8x8
	case "atkinson":
		return Atkinson
	case "sierra":
		return Sierra
	case "sierra-lite", "sierra2":
		return SierraLite
	case "stucki":
		return Stucki
	case "burkes":
		return Burkes
	case "jarvis", "jjn":
		return JarvisJudice
	default:
		return None
	}
}

// Luminance calculates perceived brightness
func Luminance(r, g, b uint8) float64 {
	// Using standard luminance weights
	return 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
}

// ContrastStretch applies contrast stretching
func ContrastStretch(img *image.Gray, minOut, maxOut uint8) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	// Find min/max in image
	minVal, maxVal := uint8(255), uint8(0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			v := img.GrayAt(x, y).Y
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
	}

	// Avoid division by zero
	if maxVal == minVal {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				result.SetGray(x, y, color.Gray{Y: (minOut + maxOut) / 2})
			}
		}
		return result
	}

	// Stretch
	scale := float64(maxOut-minOut) / float64(maxVal-minVal)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			v := img.GrayAt(x, y).Y
			newV := uint8(math.Round(float64(minOut) + float64(v-minVal)*scale))
			result.SetGray(x, y, color.Gray{Y: newV})
		}
	}

	return result
}
