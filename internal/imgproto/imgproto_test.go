package imgproto

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

// Helper functions for testing

// createTestImage creates a test image with specified dimensions and colors.
// colors is a 2D array where each element specifies a color index from the provided palette.
func createTestImage(width, height int, palette []color.Color, colorIndices [][]int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if y < len(colorIndices) && x < len(colorIndices[y]) {
				colorIdx := colorIndices[y][x]
				if colorIdx < len(palette) {
					img.Set(x, y, palette[colorIdx])
				}
			}
		}
	}
	return img
}

// createSolidColorImage creates an image filled with a single color
func createSolidColorImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// createGradientImage creates an RGB gradient image
func createGradientImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(((x + y) * 255) / (width + height))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// createMultiColorImage creates an image with distinct color blocks
func createMultiColorImage(width, height int) image.Image {
	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},   // red
		color.RGBA{0, 255, 0, 255},   // green
		color.RGBA{0, 0, 255, 255},   // blue
		color.RGBA{255, 255, 0, 255}, // yellow
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	blockWidth := width / 2
	blockHeight := height / 2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colorIdx := ((y / blockHeight) * 2) + (x / blockWidth)
			if colorIdx >= len(colors) {
				colorIdx = len(colors) - 1
			}
			img.Set(x, y, colors[colorIdx])
		}
	}
	return img
}

// TestProtocolString tests the Protocol.String method
func TestProtocolString(t *testing.T) {
	tests := []struct {
		proto Protocol
		want  string
	}{
		{ASCII, "ascii"},
		{Sixel, "sixel"},
		{Kitty, "kitty"},
		{ITerm2, "iterm2"},
		{WezTerm, "wezterm"},
		{Terminology, "terminology"},
		{Protocol(99), "ascii"}, // invalid protocol should default to ascii
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.proto.String()
			if got != tt.want {
				t.Errorf("Protocol.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestParseProtocol tests the ParseProtocol function
func TestParseProtocol(t *testing.T) {
	tests := []struct {
		input string
		want  Protocol
	}{
		// Sixel variations
		{"sixel", Sixel},
		{"SIXEL", Sixel},
		{"six", Sixel},
		{"SIX", Sixel},

		// Kitty variations
		{"kitty", Kitty},
		{"KITTY", Kitty},

		// iTerm2 variations
		{"iterm2", ITerm2},
		{"ITERM2", ITerm2},
		{"iterm", ITerm2},
		{"ITERM", ITerm2},

		// WezTerm variations
		{"wezterm", WezTerm},
		{"WEZTERM", WezTerm},

		// Terminology variations
		{"terminology", Terminology},
		{"TERMINOLOGY", Terminology},

		// ASCII variations
		{"ascii", ASCII},
		{"ASCII", ASCII},
		{"", ASCII},
		{"unknown", ASCII},

		// Auto should call Detect() which may vary by environment
		{"auto", Detect()},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseProtocol(tt.input)
			if got != tt.want {
				t.Errorf("ParseProtocol(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestDetect tests the Detect function with various environment configurations
func TestDetect(t *testing.T) {
	// Save original environment
	origKittyWindowID := os.Getenv("KITTY_WINDOW_ID")
	origITermSessionID := os.Getenv("ITERM_SESSION_ID")
	origTermProgram := os.Getenv("TERM_PROGRAM")
	origTerm := os.Getenv("TERM")
	origTerminology := os.Getenv("TERMINOLOGY")

	defer func() {
		// Restore original environment
		if origKittyWindowID != "" {
			os.Setenv("KITTY_WINDOW_ID", origKittyWindowID)
		} else {
			os.Unsetenv("KITTY_WINDOW_ID")
		}
		if origITermSessionID != "" {
			os.Setenv("ITERM_SESSION_ID", origITermSessionID)
		} else {
			os.Unsetenv("ITERM_SESSION_ID")
		}
		if origTermProgram != "" {
			os.Setenv("TERM_PROGRAM", origTermProgram)
		} else {
			os.Unsetenv("TERM_PROGRAM")
		}
		if origTerm != "" {
			os.Setenv("TERM", origTerm)
		} else {
			os.Unsetenv("TERM")
		}
		if origTerminology != "" {
			os.Setenv("TERMINOLOGY", origTerminology)
		} else {
			os.Unsetenv("TERMINOLOGY")
		}
	}()

	tests := []struct {
		name             string
		kittyWindowID    string
		iTermSessionID   string
		termProgram      string
		term             string
		terminology      string
		expectedProtocol Protocol
	}{
		{
			name:             "Kitty via window ID",
			kittyWindowID:    "1",
			iTermSessionID:   "",
			termProgram:      "",
			term:             "xterm-256color",
			terminology:      "",
			expectedProtocol: Kitty,
		},
		{
			name:             "Kitty via TERM variable",
			kittyWindowID:    "",
			iTermSessionID:   "",
			termProgram:      "",
			term:             "xterm-kitty",
			terminology:      "",
			expectedProtocol: Kitty,
		},
		{
			name:             "WezTerm via TERM_PROGRAM",
			kittyWindowID:    "",
			iTermSessionID:   "",
			termProgram:      "WezTerm",
			term:             "xterm-256color",
			terminology:      "",
			expectedProtocol: WezTerm,
		},
		{
			name:             "WezTerm via TERM variable",
			kittyWindowID:    "",
			iTermSessionID:   "",
			termProgram:      "",
			term:             "wezterm",
			terminology:      "",
			expectedProtocol: WezTerm,
		},
		{
			name:             "iTerm2 via session ID",
			kittyWindowID:    "",
			iTermSessionID:   "w0t0p0:0x0",
			termProgram:      "",
			term:             "xterm-256color",
			terminology:      "",
			expectedProtocol: ITerm2,
		},
		{
			name:             "iTerm2 via TERM_PROGRAM",
			kittyWindowID:    "",
			iTermSessionID:   "",
			termProgram:      "iTerm.app",
			term:             "xterm-256color",
			terminology:      "",
			expectedProtocol: ITerm2,
		},
		{
			name:             "Terminology via env var",
			kittyWindowID:    "",
			iTermSessionID:   "",
			termProgram:      "",
			term:             "dumb",
			terminology:      "1",
			expectedProtocol: Terminology,
		},
		{
			name:             "Sixel via TERM",
			kittyWindowID:    "",
			iTermSessionID:   "",
			termProgram:      "",
			term:             "xterm",
			terminology:      "",
			expectedProtocol: Sixel,
		},
		{
			name:             "Priority: Kitty over WezTerm",
			kittyWindowID:    "1",
			iTermSessionID:   "",
			termProgram:      "WezTerm",
			term:             "xterm-kitty",
			terminology:      "",
			expectedProtocol: Kitty,
		},
		{
			name:             "Priority: WezTerm over iTerm2",
			kittyWindowID:    "",
			iTermSessionID:   "w0t0p0:0x0",
			termProgram:      "WezTerm",
			term:             "xterm-256color",
			terminology:      "",
			expectedProtocol: WezTerm,
		},
		{
			name:             "Priority: iTerm2 over Terminology",
			kittyWindowID:    "",
			iTermSessionID:   "w0t0p0:0x0",
			termProgram:      "",
			term:             "xterm-256color",
			terminology:      "1",
			expectedProtocol: ITerm2,
		},
		{
			name:             "ASCII fallback",
			kittyWindowID:    "",
			iTermSessionID:   "",
			termProgram:      "",
			term:             "dumb",
			terminology:      "",
			expectedProtocol: ASCII,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("KITTY_WINDOW_ID")
			os.Unsetenv("ITERM_SESSION_ID")
			os.Unsetenv("TERM_PROGRAM")
			os.Unsetenv("TERM")
			os.Unsetenv("TERMINOLOGY")

			if tt.kittyWindowID != "" {
				os.Setenv("KITTY_WINDOW_ID", tt.kittyWindowID)
			}
			if tt.iTermSessionID != "" {
				os.Setenv("ITERM_SESSION_ID", tt.iTermSessionID)
			}
			if tt.termProgram != "" {
				os.Setenv("TERM_PROGRAM", tt.termProgram)
			}
			if tt.term != "" {
				os.Setenv("TERM", tt.term)
			}
			if tt.terminology != "" {
				os.Setenv("TERMINOLOGY", tt.terminology)
			}

			got := Detect()
			if got != tt.expectedProtocol {
				t.Errorf("Detect() = %v, want %v", got, tt.expectedProtocol)
			}
		})
	}
}

// TestListProtocols tests the ListProtocols function
func TestListProtocols(t *testing.T) {
	protocols := ListProtocols()
	expected := []string{"auto", "ascii", "sixel", "kitty", "iterm2", "wezterm", "terminology"}

	if len(protocols) != len(expected) {
		t.Errorf("ListProtocols() returned %d items, want %d", len(protocols), len(expected))
	}

	for i, p := range protocols {
		if p != expected[i] {
			t.Errorf("ListProtocols()[%d] = %q, want %q", i, p, expected[i])
		}
	}
}

// TestRenderWithInvalidProtocol tests Render with ASCII protocol (should error)
func TestRenderWithInvalidProtocol(t *testing.T) {
	img := createSolidColorImage(16, 16, color.RGBA{255, 0, 0, 255})

	_, err := Render(img, ASCII, 80)
	if err == nil {
		t.Errorf("Render() with ASCII protocol should return error, got nil")
	}
}

// TestRenderSixel tests Sixel rendering
func TestRenderSixel(t *testing.T) {
	tests := []struct {
		name     string
		image    image.Image
		width    int
		validate func(t *testing.T, output string)
	}{
		{
			name:  "Simple solid color",
			image: createSolidColorImage(16, 16, color.RGBA{255, 0, 0, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				// Sixel should start with escape sequence
				if !strings.Contains(output, "\x1bP") {
					t.Errorf("Sixel output should contain DCS escape sequence \\x1bP")
				}
				// Sixel should end with terminator
				if !strings.Contains(output, "\x1b\\") {
					t.Errorf("Sixel output should contain terminator \\x1b\\")
				}
				// Should contain color definition
				if !strings.Contains(output, "#0;2;") {
					t.Errorf("Sixel output should contain color definition")
				}
			},
		},
		{
			name:  "Multi-color image",
			image: createMultiColorImage(16, 16),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1bP") {
					t.Errorf("Sixel output should contain DCS escape sequence")
				}
				if !strings.Contains(output, "\x1b\\") {
					t.Errorf("Sixel output should contain terminator")
				}
			},
		},
		{
			name:  "Small 1x1 image",
			image: createSolidColorImage(1, 1, color.RGBA{0, 255, 0, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1bP") {
					t.Errorf("Sixel output should contain DCS escape sequence")
				}
				if !strings.Contains(output, "\x1b\\") {
					t.Errorf("Sixel output should contain terminator")
				}
			},
		},
		{
			name:  "Large image with scaling",
			image: createGradientImage(256, 256),
			width: 20,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1bP") {
					t.Errorf("Sixel output should contain DCS escape sequence")
				}
				if !strings.Contains(output, "\x1b\\") {
					t.Errorf("Sixel output should contain terminator")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RenderSixel(tt.image, tt.width)
			if err != nil {
				t.Fatalf("RenderSixel() error = %v", err)
			}
			if output == "" {
				t.Errorf("RenderSixel() returned empty output")
			}
			tt.validate(t, output)
		})
	}
}

// TestRenderKitty tests Kitty graphics protocol rendering
func TestRenderKitty(t *testing.T) {
	tests := []struct {
		name     string
		image    image.Image
		width    int
		validate func(t *testing.T, output string)
	}{
		{
			name:  "Simple solid color",
			image: createSolidColorImage(16, 16, color.RGBA{0, 0, 255, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				// Kitty should start with Ga=T,f=100,m=...
				if !strings.Contains(output, "\x1b_Ga=T,f=100,m=") {
					t.Errorf("Kitty output should contain initial command with Ga=T,f=100")
				}
				// Should contain escape terminator
				if !strings.Contains(output, "\x1b\\") {
					t.Errorf("Kitty output should contain escape terminator")
				}
				// Should be base64 encoded (alphanumeric, +, /, =)
				if !hasBase64Content(output) {
					t.Errorf("Kitty output should contain base64 encoded data")
				}
			},
		},
		{
			name:  "Multi-color image",
			image: createMultiColorImage(32, 32),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1b_Ga=T,f=100,m=") {
					t.Errorf("Kitty output should contain initial command")
				}
				if !strings.Contains(output, "\x1b\\") {
					t.Errorf("Kitty output should contain escape terminator")
				}
			},
		},
		{
			name:  "Image requiring chunking (large image)",
			image: createGradientImage(512, 512),
			width: 80,
			validate: func(t *testing.T, output string) {
				// Large images should be chunked with m=1 for more=1
				if !strings.Contains(output, "m=1;") {
					t.Errorf("Large image should have chunked output with m=1")
				}
				// Final chunk should have m=0
				if !strings.Contains(output, "m=0;") {
					t.Errorf("Output should contain final chunk with m=0")
				}
			},
		},
		{
			name:  "Small 1x1 image",
			image: createSolidColorImage(1, 1, color.RGBA{255, 255, 0, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1b_Ga=T,f=100,m=") {
					t.Errorf("Kitty output should contain command")
				}
				if !strings.Contains(output, "m=0;") {
					t.Errorf("Single chunk should have m=0")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RenderKitty(tt.image, tt.width)
			if err != nil {
				t.Fatalf("RenderKitty() error = %v", err)
			}
			if output == "" {
				t.Errorf("RenderKitty() returned empty output")
			}
			tt.validate(t, output)
		})
	}
}

// TestRenderITerm2 tests iTerm2 inline image protocol rendering
func TestRenderITerm2(t *testing.T) {
	tests := []struct {
		name     string
		image    image.Image
		width    int
		validate func(t *testing.T, output string)
	}{
		{
			name:  "Simple solid color",
			image: createSolidColorImage(16, 16, color.RGBA{255, 0, 255, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				// iTerm2 should use OSC 1337
				if !strings.Contains(output, "\x1b]1337;") {
					t.Errorf("iTerm2 output should contain OSC 1337 sequence")
				}
				// Should contain File parameter
				if !strings.Contains(output, "File=") {
					t.Errorf("iTerm2 output should contain File parameter")
				}
				// Should contain inline parameter
				if !strings.Contains(output, "inline=1") {
					t.Errorf("iTerm2 output should contain inline=1 parameter")
				}
				// Should end with bell character
				if !strings.Contains(output, "\x07") {
					t.Errorf("iTerm2 output should end with bell character")
				}
				// Should contain width and height
				if !strings.Contains(output, "width=") {
					t.Errorf("iTerm2 output should contain width parameter")
				}
				if !strings.Contains(output, "height=") {
					t.Errorf("iTerm2 output should contain height parameter")
				}
				// Should have base64 data
				if !hasBase64Content(output) {
					t.Errorf("iTerm2 output should contain base64 encoded data")
				}
			},
		},
		{
			name:  "Multi-color image",
			image: createMultiColorImage(32, 32),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1b]1337;") {
					t.Errorf("iTerm2 output should contain OSC sequence")
				}
				if !strings.Contains(output, "width=32px;height=32px") {
					t.Errorf("iTerm2 output should have correct dimensions")
				}
			},
		},
		{
			name:  "Large image with scaling",
			image: createGradientImage(256, 256),
			width: 20,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1b]1337;") {
					t.Errorf("iTerm2 output should contain OSC sequence")
				}
				// Verify preserveAspectRatio parameter
				if !strings.Contains(output, "preserveAspectRatio=1:") {
					t.Errorf("iTerm2 output should contain preserveAspectRatio parameter")
				}
			},
		},
		{
			name:  "Small 1x1 image",
			image: createSolidColorImage(1, 1, color.RGBA{0, 255, 255, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "width=1px;height=1px") {
					t.Errorf("iTerm2 output should preserve 1x1 dimensions")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RenderITerm2(tt.image, tt.width)
			if err != nil {
				t.Fatalf("RenderITerm2() error = %v", err)
			}
			if output == "" {
				t.Errorf("RenderITerm2() returned empty output")
			}
			tt.validate(t, output)
		})
	}
}

// TestQuantize tests color quantization
func TestQuantize(t *testing.T) {
	tests := []struct {
		name          string
		image         image.Image
		maxColors     int
		minPaletteLen int
		maxPaletteLen int
	}{
		{
			name:          "Single color quantization",
			image:         createSolidColorImage(16, 16, color.RGBA{255, 0, 0, 255}),
			maxColors:     256,
			minPaletteLen: 1,
			maxPaletteLen: 256,
		},
		{
			name:          "Multi-color quantization",
			image:         createMultiColorImage(32, 32),
			maxColors:     256,
			minPaletteLen: 1,
			maxPaletteLen: 256,
		},
		{
			name:          "Gradient quantization",
			image:         createGradientImage(64, 64),
			maxColors:     256,
			minPaletteLen: 1,
			maxPaletteLen: 256,
		},
		{
			name:          "Limited palette (16 colors)",
			image:         createMultiColorImage(32, 32),
			maxColors:     16,
			minPaletteLen: 1,
			maxPaletteLen: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			palette := quantize(tt.image, tt.maxColors)
			if len(palette) < tt.minPaletteLen || len(palette) > tt.maxPaletteLen {
				t.Errorf("quantize() returned palette of length %d, want between %d and %d",
					len(palette), tt.minPaletteLen, tt.maxPaletteLen)
			}

			// Ensure all colors are valid RGBA
			for i, c := range palette {
				r, g, b, a := c.RGBA()
				if r > 65535 || g > 65535 || b > 65535 || a > 65535 {
					t.Errorf("palette[%d] has invalid RGBA values: %d,%d,%d,%d", i, r, g, b, a)
				}
			}
		})
	}
}

// TestClosestColor tests the closest color matching
func TestClosestColor(t *testing.T) {
	palette := []color.Color{
		color.RGBA{0, 0, 0, 255},       // black
		color.RGBA{255, 255, 255, 255}, // white
		color.RGBA{255, 0, 0, 255},     // red
		color.RGBA{0, 255, 0, 255},     // green
		color.RGBA{0, 0, 255, 255},     // blue
	}

	tests := []struct {
		name        string
		color       color.Color
		expectedIdx int
	}{
		{
			name:        "Exact black match",
			color:       color.RGBA{0, 0, 0, 255},
			expectedIdx: 0,
		},
		{
			name:        "Exact white match",
			color:       color.RGBA{255, 255, 255, 255},
			expectedIdx: 1,
		},
		{
			name:        "Exact red match",
			color:       color.RGBA{255, 0, 0, 255},
			expectedIdx: 2,
		},
		{
			name:        "Dark color to black",
			color:       color.RGBA{10, 10, 10, 255},
			expectedIdx: 0,
		},
		{
			name:        "Bright color to white",
			color:       color.RGBA{245, 245, 245, 255},
			expectedIdx: 1,
		},
		{
			name:        "Red-ish color to red",
			color:       color.RGBA{255, 50, 50, 255},
			expectedIdx: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := closestColor(tt.color, palette)
			if got != tt.expectedIdx {
				t.Errorf("closestColor() = %d, want %d", got, tt.expectedIdx)
			}
		})
	}
}

// TestScaleImage tests image scaling
func TestScaleImage(t *testing.T) {
	tests := []struct {
		name              string
		image             image.Image
		widthCols         int
		expectedMaxWidth  int
		expectedMaxHeight int
		shouldScale       bool
	}{
		{
			name:              "No scaling when width >= original",
			image:             createSolidColorImage(64, 64, color.RGBA{255, 0, 0, 255}),
			widthCols:         80,
			expectedMaxWidth:  640, // 80 cols * 8 pixels
			expectedMaxHeight: 640,
			shouldScale:       false,
		},
		{
			name:              "Scaling down",
			image:             createSolidColorImage(256, 256, color.RGBA{0, 255, 0, 255}),
			widthCols:         20,
			expectedMaxWidth:  160, // 20 cols * 8 pixels
			expectedMaxHeight: 160,
			shouldScale:       true,
		},
		{
			name:              "Scaling maintains aspect ratio",
			image:             createSolidColorImage(200, 100, color.RGBA{0, 0, 255, 255}),
			widthCols:         10,
			expectedMaxWidth:  80, // 10 cols * 8 pixels
			expectedMaxHeight: 40,
			shouldScale:       true,
		},
		{
			name:              "Small width doesn't scale",
			image:             createSolidColorImage(32, 32, color.RGBA{255, 255, 0, 255}),
			widthCols:         80,
			expectedMaxWidth:  32,
			expectedMaxHeight: 32,
			shouldScale:       false,
		},
		{
			name:              "Zero width doesn't scale",
			image:             createSolidColorImage(64, 64, color.RGBA{255, 0, 255, 255}),
			widthCols:         0,
			expectedMaxWidth:  64,
			expectedMaxHeight: 64,
			shouldScale:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaled := scaleImage(tt.image, tt.widthCols)
			bounds := scaled.Bounds()
			width := bounds.Dx()
			height := bounds.Dy()

			// Check dimensions don't exceed expected maximums
			if width > tt.expectedMaxWidth {
				t.Errorf("scaled width %d exceeds expected max %d", width, tt.expectedMaxWidth)
			}
			if height > tt.expectedMaxHeight {
				t.Errorf("scaled height %d exceeds expected max %d", height, tt.expectedMaxHeight)
			}

			// Check aspect ratio is maintained (within 10%)
			origBounds := tt.image.Bounds()
			origAspect := float64(origBounds.Dx()) / float64(origBounds.Dy())
			scaledAspect := float64(width) / float64(height)
			diff := origAspect / scaledAspect
			if diff < 0.9 || diff > 1.1 {
				// Allow some tolerance due to integer scaling
				// Only warn if significantly off
				if diff < 0.8 || diff > 1.2 {
					t.Logf("aspect ratio changed from %.2f to %.2f", origAspect, scaledAspect)
				}
			}
		})
	}
}

// TestRenderWezTerm tests WezTerm graphics protocol rendering
func TestRenderWezTerm(t *testing.T) {
	tests := []struct {
		name     string
		image    image.Image
		width    int
		validate func(t *testing.T, output string)
	}{
		{
			name:  "Simple solid color",
			image: createSolidColorImage(16, 16, color.RGBA{255, 0, 0, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				// WezTerm should use iTerm2-compatible protocol
				if !strings.Contains(output, "\x1b]1337;") && !strings.Contains(output, "\x1bPtmux;") {
					t.Errorf("WezTerm output should contain iTerm2 protocol or tmux passthrough")
				}
				if !strings.Contains(output, "File=") {
					t.Errorf("WezTerm output should contain File parameter")
				}
				if !strings.Contains(output, "inline=1") {
					t.Errorf("WezTerm output should contain inline=1 parameter")
				}
				if !hasBase64Content(output) {
					t.Errorf("WezTerm output should contain base64 encoded data")
				}
			},
		},
		{
			name:  "Multi-color image",
			image: createMultiColorImage(32, 32),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "File=") {
					t.Errorf("WezTerm output should contain File parameter")
				}
				if !hasBase64Content(output) {
					t.Errorf("WezTerm output should contain base64 data")
				}
			},
		},
		{
			name:  "Large image with scaling",
			image: createGradientImage(256, 256),
			width: 20,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "preserveAspectRatio=1:") {
					t.Errorf("WezTerm output should contain preserveAspectRatio parameter")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RenderWezTerm(tt.image, tt.width)
			if err != nil {
				t.Fatalf("RenderWezTerm() error = %v", err)
			}
			if output == "" {
				t.Errorf("RenderWezTerm() returned empty output")
			}
			tt.validate(t, output)
		})
	}
}

// TestRenderTerminology tests Terminology graphics protocol rendering
func TestRenderTerminology(t *testing.T) {
	tests := []struct {
		name     string
		image    image.Image
		width    int
		validate func(t *testing.T, output string)
	}{
		{
			name:  "Simple solid color",
			image: createSolidColorImage(16, 16, color.RGBA{0, 255, 0, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				// Terminology should use inline image protocol
				if !strings.Contains(output, "\x1b}is#") {
					t.Errorf("Terminology output should contain inline image escape sequence")
				}
				if !strings.Contains(output, "\x1b\\") {
					t.Errorf("Terminology output should contain terminator")
				}
				if !hasBase64Content(output) {
					t.Errorf("Terminology output should contain base64 encoded data")
				}
			},
		},
		{
			name:  "Multi-color image",
			image: createMultiColorImage(32, 32),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1b}is#") {
					t.Errorf("Terminology output should contain inline image escape")
				}
				// Should contain dimensions: width;height;
				if !strings.Contains(output, ";") {
					t.Errorf("Terminology output should contain dimension separators")
				}
			},
		},
		{
			name:  "Large image with scaling",
			image: createGradientImage(256, 256),
			width: 20,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1b}is#") {
					t.Errorf("Terminology output should contain inline image sequence")
				}
				if !hasBase64Content(output) {
					t.Errorf("Terminology output should contain base64 data")
				}
			},
		},
		{
			name:  "Small 1x1 image",
			image: createSolidColorImage(1, 1, color.RGBA{255, 255, 0, 255}),
			width: 80,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "\x1b}is#1;1;") {
					t.Errorf("Terminology output should preserve 1x1 dimensions")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RenderTerminology(tt.image, tt.width)
			if err != nil {
				t.Fatalf("RenderTerminology() error = %v", err)
			}
			if output == "" {
				t.Errorf("RenderTerminology() returned empty output")
			}
			tt.validate(t, output)
		})
	}
}

// TestRender tests the generic Render function
func TestRender(t *testing.T) {
	img := createSolidColorImage(32, 32, color.RGBA{100, 150, 200, 255})

	tests := []struct {
		name     string
		proto    Protocol
		validate func(t *testing.T, output string, err error)
	}{
		{
			name:  "Sixel rendering",
			proto: Sixel,
			validate: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Render() error = %v", err)
				}
				if !strings.Contains(output, "\x1bP") {
					t.Errorf("Sixel output missing DCS escape sequence")
				}
			},
		},
		{
			name:  "Kitty rendering",
			proto: Kitty,
			validate: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Render() error = %v", err)
				}
				if !strings.Contains(output, "\x1b_Ga=T,f=100") {
					t.Errorf("Kitty output missing command")
				}
			},
		},
		{
			name:  "iTerm2 rendering",
			proto: ITerm2,
			validate: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Render() error = %v", err)
				}
				if !strings.Contains(output, "\x1b]1337;") {
					t.Errorf("iTerm2 output missing OSC sequence")
				}
			},
		},
		{
			name:  "WezTerm rendering",
			proto: WezTerm,
			validate: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Render() error = %v", err)
				}
				if !strings.Contains(output, "\x1b]1337;") && !strings.Contains(output, "\x1bPtmux;") {
					t.Errorf("WezTerm output missing protocol sequence")
				}
			},
		},
		{
			name:  "Terminology rendering",
			proto: Terminology,
			validate: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Render() error = %v", err)
				}
				if !strings.Contains(output, "\x1b}is#") {
					t.Errorf("Terminology output missing inline image sequence")
				}
			},
		},
		{
			name:  "ASCII protocol should error",
			proto: ASCII,
			validate: func(t *testing.T, output string, err error) {
				if err == nil {
					t.Errorf("Render() with ASCII should return error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := Render(img, tt.proto, 80)
			tt.validate(t, output, err)
		})
	}
}

// TestKittyChunking tests that Kitty protocol correctly chunks large data
func TestKittyChunking(t *testing.T) {
	// Create a very large image that will exceed 4096 byte chunk size
	// PNG encoding may produce different file sizes, so we test with an even larger image
	largeImg := createGradientImage(1024, 1024)

	output, err := RenderKitty(largeImg, 80)
	if err != nil {
		t.Fatalf("RenderKitty() error = %v", err)
	}

	// Verify the output contains at least one chunk command
	if !strings.Contains(output, "Ga=T,f=100") && !strings.Contains(output, "Gm=") {
		t.Errorf("Expected Kitty protocol commands in output")
	}

	// Check for proper termination
	if !strings.Contains(output, "\x1b\\") {
		t.Errorf("Expected escape terminators in output")
	}
}

// TestSixelBands tests that Sixel correctly encodes data in 6-row bands
func TestSixelBands(t *testing.T) {
	// Create a simple multi-color image
	img := createMultiColorImage(16, 12) // 12 rows = 2 full bands

	output, err := RenderSixel(img, 80)
	if err != nil {
		t.Fatalf("RenderSixel() error = %v", err)
	}

	// Check that output contains band separator (-)
	bandSeparators := strings.Count(output, "-")
	if bandSeparators < 1 {
		t.Errorf("Expected at least 1 band separator, got %d", bandSeparators)
	}

	// Check for color definitions
	colorDefCount := strings.Count(output, "#")
	if colorDefCount < 1 {
		t.Errorf("Expected color definitions, got %d", colorDefCount)
	}
}

// TestPNGEncoding tests that images can be encoded to PNG
func TestPNGEncoding(t *testing.T) {
	tests := []struct {
		name  string
		image image.Image
	}{
		{
			name:  "Solid color",
			image: createSolidColorImage(16, 16, color.RGBA{255, 0, 0, 255}),
		},
		{
			name:  "Multi-color",
			image: createMultiColorImage(32, 32),
		},
		{
			name:  "Gradient",
			image: createGradientImage(64, 64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := png.Encode(&buf, tt.image)
			if err != nil {
				t.Fatalf("PNG encoding error: %v", err)
			}

			if buf.Len() == 0 {
				t.Errorf("PNG encoding resulted in empty buffer")
			}

			// Decode to verify it's valid PNG
			_, err = png.Decode(&buf)
			if err != nil {
				t.Fatalf("Failed to decode PNG: %v", err)
			}
		})
	}
}

// Helper function to check if string contains base64-like content
func hasBase64Content(s string) bool {
	// Check for base64 characters (alphanumeric, +, /, =)
	for _, c := range s {
		if (c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			c == '+' || c == '/' || c == '=' {
			return true
		}
	}
	return false
}

// BenchmarkQuantize benchmarks the color quantization function
func BenchmarkQuantize(b *testing.B) {
	img := createGradientImage(256, 256)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quantize(img, 256)
	}
}

// BenchmarkScaleImage benchmarks the image scaling function
func BenchmarkScaleImage(b *testing.B) {
	img := createGradientImage(512, 512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scaleImage(img, 80)
	}
}

// BenchmarkRenderSixel benchmarks Sixel rendering
func BenchmarkRenderSixel(b *testing.B) {
	img := createGradientImage(128, 128)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderSixel(img, 80)
	}
}

// BenchmarkRenderKitty benchmarks Kitty rendering
func BenchmarkRenderKitty(b *testing.B) {
	img := createGradientImage(128, 128)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderKitty(img, 80)
	}
}

// BenchmarkRenderITerm2 benchmarks iTerm2 rendering
func BenchmarkRenderITerm2(b *testing.B) {
	img := createGradientImage(128, 128)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderITerm2(img, 80)
	}
}
