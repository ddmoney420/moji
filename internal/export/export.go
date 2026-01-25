package export

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// ToPNG exports ASCII art to a PNG image
func ToPNG(text, filename, bgHex, fgHex string) error {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")

	// Calculate dimensions
	charWidth := 7   // basicfont.Face7x13 character width
	charHeight := 13 // basicfont.Face7x13 character height
	padding := 20

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	imgWidth := maxWidth*charWidth + padding*2
	imgHeight := len(lines)*charHeight + padding*2

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Parse colors
	bgColor := parseHexColor(bgHex)
	fgColor := parseHexColor(fgHex)

	// Fill background
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// Draw text
	face := basicfont.Face7x13
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fgColor),
		Face: face,
	}

	for i, line := range lines {
		x := padding
		y := padding + (i+1)*charHeight - 3 // baseline adjustment
		drawer.Dot = fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		}
		drawer.DrawString(line)
	}

	// Save to file
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// parseHexColor converts a hex color string to color.RGBA
func parseHexColor(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return color.RGBA{0, 0, 0, 255} // default black
	}

	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{r, g, b, 255}
}

// ToSVG exports ASCII art to SVG (for better scalability)
func ToSVG(text, filename, bgHex, fgHex string) error {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")

	charWidth := 8.4   // approximate width for monospace
	charHeight := 14.0 // line height
	padding := 20.0

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	width := float64(maxWidth)*charWidth + padding*2
	height := float64(len(lines))*charHeight + padding*2

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Write SVG header
	fmt.Fprintf(f, `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f">
  <rect width="100%%" height="100%%" fill="%s"/>
  <text font-family="monospace" font-size="12" fill="%s">
`, width, height, bgHex, fgHex)

	// Write text lines
	for i, line := range lines {
		y := padding + float64(i+1)*charHeight
		// Escape XML special characters
		escaped := escapeXML(line)
		fmt.Fprintf(f, `    <tspan x="%.0f" y="%.0f">%s</tspan>
`, padding, y, escaped)
	}

	// Close SVG
	fmt.Fprintf(f, `  </text>
</svg>
`)

	return nil
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// ToHTML exports ASCII art to HTML with styling
func ToHTML(text, filename, bgHex, fgHex, title string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	escaped := escapeXML(text)

	fmt.Fprintf(f, `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>%s</title>
  <style>
    body {
      background-color: %s;
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 100vh;
      margin: 0;
      padding: 20px;
      box-sizing: border-box;
    }
    pre {
      color: %s;
      font-family: 'Courier New', Courier, monospace;
      font-size: 14px;
      line-height: 1.2;
      white-space: pre;
    }
  </style>
</head>
<body>
  <pre>%s</pre>
</body>
</html>
`, title, bgHex, fgHex, escaped)

	return nil
}
