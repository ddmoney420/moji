package halfblock

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// Half-block characters for 2x vertical resolution
// Each character cell represents 2 vertical pixels
const (
	FullBlock  = '█' // Both pixels on
	UpperHalf  = '▀' // Top pixel on
	LowerHalf  = '▄' // Bottom pixel on
	EmptyBlock = ' ' // Both pixels off
)

// RenderGrayscale converts an image to half-block ASCII with 2x vertical resolution
func RenderGrayscale(img image.Image, width int, threshold uint8) string {
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	// Calculate dimensions
	if width <= 0 {
		width = 80
	}

	// Scale image
	scaleX := float64(imgWidth) / float64(width)
	scaleY := scaleX * 2 // 2x because each char is 2 pixels tall

	height := int(float64(imgHeight) / scaleY)
	if height <= 0 {
		height = 1
	}

	var sb strings.Builder

	// Process 2 rows at a time
	for y := 0; y < height*2; y += 2 {
		for x := 0; x < width; x++ {
			// Sample top and bottom pixels
			srcX := int(float64(x)*scaleX) + bounds.Min.X
			srcY1 := int(float64(y)*scaleY/2) + bounds.Min.Y
			srcY2 := int(float64(y+1)*scaleY/2) + bounds.Min.Y

			// Get brightness values
			top := getBrightness(img, srcX, srcY1, bounds) > threshold
			bottom := getBrightness(img, srcX, srcY2, bounds) > threshold

			// Choose character based on which pixels are "on"
			switch {
			case top && bottom:
				sb.WriteRune(FullBlock)
			case top && !bottom:
				sb.WriteRune(UpperHalf)
			case !top && bottom:
				sb.WriteRune(LowerHalf)
			default:
				sb.WriteRune(EmptyBlock)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderColor converts an image to half-block ASCII with color
func RenderColor(img image.Image, width int) string {
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	if width <= 0 {
		width = 80
	}

	scaleX := float64(imgWidth) / float64(width)
	scaleY := scaleX * 2

	height := int(float64(imgHeight) / scaleY)
	if height <= 0 {
		height = 1
	}

	var sb strings.Builder

	for y := 0; y < height*2; y += 2 {
		for x := 0; x < width; x++ {
			srcX := int(float64(x)*scaleX) + bounds.Min.X
			srcY1 := int(float64(y)*scaleY/2) + bounds.Min.Y
			srcY2 := int(float64(y+1)*scaleY/2) + bounds.Min.Y

			// Get colors
			r1, g1, b1 := getColor(img, srcX, srcY1, bounds)
			r2, g2, b2 := getColor(img, srcX, srcY2, bounds)

			// Use upper half block with foreground=top color, background=bottom color
			sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%c\033[0m",
				r1, g1, b1, r2, g2, b2, UpperHalf))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderWithCharset renders using half-blocks for shape and charset for detail
func RenderWithCharset(img image.Image, width int, charset string) string {
	if charset == "" {
		charset = " ░▒▓█"
	}

	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	if width <= 0 {
		width = 80
	}

	scaleX := float64(imgWidth) / float64(width)
	scaleY := scaleX // Normal aspect ratio for charset mode

	height := int(float64(imgHeight) / scaleY)
	if height <= 0 {
		height = 1
	}

	runes := []rune(charset)

	var sb strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x)*scaleX) + bounds.Min.X
			srcY := int(float64(y)*scaleY) + bounds.Min.Y

			brightness := getBrightness(img, srcX, srcY, bounds)
			charIdx := int(float64(brightness) / 255.0 * float64(len(runes)-1))
			if charIdx >= len(runes) {
				charIdx = len(runes) - 1
			}

			sb.WriteRune(runes[charIdx])
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// getBrightness returns the brightness value at a point
func getBrightness(img image.Image, x, y int, bounds image.Rectangle) uint8 {
	if x < bounds.Min.X {
		x = bounds.Min.X
	}
	if x >= bounds.Max.X {
		x = bounds.Max.X - 1
	}
	if y < bounds.Min.Y {
		y = bounds.Min.Y
	}
	if y >= bounds.Max.Y {
		y = bounds.Max.Y - 1
	}

	c := img.At(x, y)
	r, g, b, _ := c.RGBA()

	// Luminance formula
	return uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
}

// getColor returns RGB values at a point
func getColor(img image.Image, x, y int, bounds image.Rectangle) (uint8, uint8, uint8) {
	if x < bounds.Min.X {
		x = bounds.Min.X
	}
	if x >= bounds.Max.X {
		x = bounds.Max.X - 1
	}
	if y < bounds.Min.Y {
		y = bounds.Min.Y
	}
	if y >= bounds.Max.Y {
		y = bounds.Max.Y - 1
	}

	c := img.At(x, y)
	r, g, b, _ := c.RGBA()

	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
}

// TextToHalfBlock converts text to half-block representation
// Each character becomes a 2x2 pattern
func TextToHalfBlock(text string, filled rune) string {
	if filled == 0 {
		filled = '█'
	}

	lines := strings.Split(text, "\n")
	var sb strings.Builder

	for _, line := range lines {
		// Each line becomes 2 output lines due to half-blocks
		var topLine, bottomLine strings.Builder

		for _, r := range line {
			if r == ' ' || r == '\t' {
				topLine.WriteRune(' ')
				bottomLine.WriteRune(' ')
			} else {
				topLine.WriteRune(filled)
				bottomLine.WriteRune(filled)
			}
		}

		sb.WriteString(topLine.String())
		sb.WriteString("\n")
		sb.WriteString(bottomLine.String())
		sb.WriteString("\n")
	}

	return sb.String()
}

// SmoothScale applies simple smoothing to an image for better half-block rendering
func SmoothScale(img image.Image, targetWidth, targetHeight int) *image.RGBA {
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	result := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	scaleX := float64(srcWidth) / float64(targetWidth)
	scaleY := float64(srcHeight) / float64(targetHeight)

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// Simple bilinear sampling
			srcX := float64(x) * scaleX
			srcY := float64(y) * scaleY

			x0 := int(srcX)
			y0 := int(srcY)
			x1 := x0 + 1
			y1 := y0 + 1

			if x1 >= srcWidth {
				x1 = srcWidth - 1
			}
			if y1 >= srcHeight {
				y1 = srcHeight - 1
			}

			xFrac := srcX - float64(x0)
			yFrac := srcY - float64(y0)

			c00 := img.At(x0+bounds.Min.X, y0+bounds.Min.Y)
			c10 := img.At(x1+bounds.Min.X, y0+bounds.Min.Y)
			c01 := img.At(x0+bounds.Min.X, y1+bounds.Min.Y)
			c11 := img.At(x1+bounds.Min.X, y1+bounds.Min.Y)

			r00, g00, b00, a00 := c00.RGBA()
			r10, g10, b10, a10 := c10.RGBA()
			r01, g01, b01, a01 := c01.RGBA()
			r11, g11, b11, a11 := c11.RGBA()

			r := bilinear(float64(r00), float64(r10), float64(r01), float64(r11), xFrac, yFrac)
			g := bilinear(float64(g00), float64(g10), float64(g01), float64(g11), xFrac, yFrac)
			b := bilinear(float64(b00), float64(b10), float64(b01), float64(b11), xFrac, yFrac)
			a := bilinear(float64(a00), float64(a10), float64(a01), float64(a11), xFrac, yFrac)

			result.Set(x, y, color.RGBA{
				R: uint8(r / 256),
				G: uint8(g / 256),
				B: uint8(b / 256),
				A: uint8(a / 256),
			})
		}
	}

	return result
}

func bilinear(c00, c10, c01, c11, xFrac, yFrac float64) float64 {
	top := c00*(1-xFrac) + c10*xFrac
	bottom := c01*(1-xFrac) + c11*xFrac
	return top*(1-yFrac) + bottom*yFrac
}
