package imgproto

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/ddmoney420/moji/internal/terminal"
)

// Protocol represents an image display protocol
type Protocol int

const (
	// ASCII is the default text-based rendering
	ASCII Protocol = iota
	// Sixel is DEC Sixel graphics
	Sixel
	// Kitty is the Kitty graphics protocol
	Kitty
	// ITerm2 is iTerm2 inline images
	ITerm2
)

// String returns the protocol name
func (p Protocol) String() string {
	switch p {
	case Sixel:
		return "sixel"
	case Kitty:
		return "kitty"
	case ITerm2:
		return "iterm2"
	default:
		return "ascii"
	}
}

// ParseProtocol parses a protocol name string
func ParseProtocol(s string) Protocol {
	switch strings.ToLower(s) {
	case "sixel", "six":
		return Sixel
	case "kitty":
		return Kitty
	case "iterm2", "iterm":
		return ITerm2
	case "auto":
		return Detect()
	default:
		return ASCII
	}
}

// Detect auto-detects the best available protocol
func Detect() Protocol {
	caps := terminal.Detect()
	if caps.Kitty {
		return Kitty
	}
	if caps.ITerm2 {
		return ITerm2
	}
	if caps.Sixel {
		return Sixel
	}
	return ASCII
}

// Render renders an image using the specified protocol
func Render(img image.Image, proto Protocol, width int) (string, error) {
	switch proto {
	case Sixel:
		return RenderSixel(img, width)
	case Kitty:
		return RenderKitty(img, width)
	case ITerm2:
		return RenderITerm2(img, width)
	default:
		return "", fmt.Errorf("protocol %s requires ASCII conversion, not image rendering", proto)
	}
}

// RenderKitty renders an image using the Kitty graphics protocol
func RenderKitty(img image.Image, width int) (string, error) {
	// Scale image to fit width
	img = scaleImage(img, width)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("encoding PNG: %w", err)
	}

	// Base64 encode
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Kitty protocol: split into chunks of 4096 bytes
	var result strings.Builder
	chunkSize := 4096

	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[i:end]

		more := 1
		if end >= len(encoded) {
			more = 0
		}

		if i == 0 {
			// First chunk: include format and action
			result.WriteString(fmt.Sprintf("\x1b_Ga=T,f=100,m=%d;%s\x1b\\", more, chunk))
		} else {
			result.WriteString(fmt.Sprintf("\x1b_Gm=%d;%s\x1b\\", more, chunk))
		}
	}

	result.WriteString("\n")
	return result.String(), nil
}

// RenderITerm2 renders an image using iTerm2 inline image protocol
func RenderITerm2(img image.Image, width int) (string, error) {
	// Scale image to fit width
	img = scaleImage(img, width)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("encoding PNG: %w", err)
	}

	// Base64 encode the image
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// iTerm2 inline image protocol: OSC 1337
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	result := fmt.Sprintf("\x1b]1337;File=inline=1;width=%dpx;height=%dpx;preserveAspectRatio=1:%s\x07\n",
		w, h, encoded)

	return result, nil
}

// RenderSixel renders an image using the Sixel graphics protocol
func RenderSixel(img image.Image, width int) (string, error) {
	// Scale image to fit width
	img = scaleImage(img, width)

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Quantize colors to a palette (max 256 for Sixel)
	palette := quantize(img, 256)

	var result strings.Builder

	// Sixel header: DCS P1;P2;P3 q
	// P1=0 (pixel aspect ratio 2:1), P2=1 (background transparent), P3=0 (horizontal grid size)
	result.WriteString("\x1bPq\n")

	// Define color palette
	for i, c := range palette {
		r, g, b, _ := c.RGBA()
		// Sixel uses percentage (0-100) for RGB
		rp := int(float64(r>>8) / 255.0 * 100)
		gp := int(float64(g>>8) / 255.0 * 100)
		bp := int(float64(b>>8) / 255.0 * 100)
		result.WriteString(fmt.Sprintf("#%d;2;%d;%d;%d", i, rp, gp, bp))
	}

	// Create color index map for the image
	colorMap := make([][]int, h)
	for y := 0; y < h; y++ {
		colorMap[y] = make([]int, w)
		for x := 0; x < w; x++ {
			colorMap[y][x] = closestColor(img.At(x+bounds.Min.X, y+bounds.Min.Y), palette)
		}
	}

	// Encode Sixel data (6 rows at a time)
	for band := 0; band < h; band += 6 {
		// For each color in use in this band
		for colorIdx := range palette {
			var line strings.Builder
			hasPixels := false

			for x := 0; x < w; x++ {
				sixelByte := byte(0)
				for row := 0; row < 6; row++ {
					y := band + row
					if y < h && colorMap[y][x] == colorIdx {
						sixelByte |= 1 << uint(row)
						hasPixels = true
					}
				}
				line.WriteByte(sixelByte + 63) // Sixel chars start at '?'
			}

			if hasPixels {
				result.WriteString(fmt.Sprintf("#%d", colorIdx))
				result.WriteString(line.String())
				result.WriteString("$") // Carriage return (stay on same line)
			}
		}
		result.WriteString("-") // Line feed (move to next sixel band)
	}

	// Sixel terminator
	result.WriteString("\x1b\\")

	return result.String(), nil
}

// ListProtocols returns available protocol names
func ListProtocols() []string {
	return []string{"auto", "ascii", "sixel", "kitty", "iterm2"}
}

// WriteToTerminal writes the rendered image directly to the terminal
func WriteToTerminal(img image.Image, proto Protocol, width int) error {
	if proto == ASCII {
		return fmt.Errorf("use ASCII conversion for text-based output")
	}

	output, err := Render(img, proto, width)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, output)
	return err
}

// scaleImage scales an image to fit within the given width (in pixels)
func scaleImage(img image.Image, widthCols int) image.Image {
	bounds := img.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	// Each terminal column is roughly 8 pixels wide
	targetW := widthCols * 8
	if targetW <= 0 || targetW >= origW {
		return img
	}

	scale := float64(targetW) / float64(origW)
	newW := int(float64(origW) * scale)
	newH := int(float64(origH) * scale)

	// Simple nearest-neighbor scaling
	scaled := image.NewRGBA(image.Rect(0, 0, newW, newH))
	for y := 0; y < newH; y++ {
		srcY := int(float64(y) / scale)
		if srcY >= origH {
			srcY = origH - 1
		}
		for x := 0; x < newW; x++ {
			srcX := int(float64(x) / scale)
			if srcX >= origW {
				srcX = origW - 1
			}
			scaled.Set(x, y, img.At(srcX+bounds.Min.X, srcY+bounds.Min.Y))
		}
	}

	return scaled
}

// quantize reduces image colors to a palette of n colors
func quantize(img image.Image, n int) []color.Color {
	bounds := img.Bounds()

	// Simple median-cut quantization approximation
	// Sample colors from the image
	colorCounts := make(map[uint32]int)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 2 {
			r, g, b, _ := img.At(x, y).RGBA()
			// Reduce to 5 bits per channel for grouping
			key := ((r >> 11) << 10) | ((g >> 11) << 5) | (b >> 11)
			colorCounts[key]++
		}
	}

	// Take top N most frequent colors
	type colorFreq struct {
		key   uint32
		count int
	}
	sorted := make([]colorFreq, 0, len(colorCounts))
	for k, v := range colorCounts {
		sorted = append(sorted, colorFreq{k, v})
	}
	// Simple insertion sort for our purposes
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].count > sorted[j-1].count; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	palette := make([]color.Color, 0, n)
	for i := 0; i < len(sorted) && i < n; i++ {
		key := sorted[i].key
		r := uint8(((key >> 10) & 0x1f) << 3)
		g := uint8(((key >> 5) & 0x1f) << 3)
		b := uint8((key & 0x1f) << 3)
		palette = append(palette, color.RGBA{r, g, b, 255})
	}

	// Ensure at least black and white
	if len(palette) == 0 {
		palette = append(palette, color.RGBA{0, 0, 0, 255})
		palette = append(palette, color.RGBA{255, 255, 255, 255})
	}

	return palette
}

// closestColor finds the closest palette color to the given color
func closestColor(c color.Color, palette []color.Color) int {
	r1, g1, b1, _ := c.RGBA()
	bestIdx := 0
	bestDist := uint32(0xFFFFFFFF)

	for i, pc := range palette {
		r2, g2, b2, _ := pc.RGBA()
		dr := int32(r1>>8) - int32(r2>>8)
		dg := int32(g1>>8) - int32(g2>>8)
		db := int32(b1>>8) - int32(b2>>8)
		dist := uint32(dr*dr + dg*dg + db*db)
		if dist < bestDist {
			bestDist = dist
			bestIdx = i
		}
	}

	return bestIdx
}
