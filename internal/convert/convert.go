package convert

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// Options for ASCII conversion
type Options struct {
	Width      int    // Target width in characters (0 = auto)
	Height     int    // Target height in characters (0 = auto based on width)
	Charset    string // Characters to use (dark to light)
	Invert     bool   // Invert brightness
	EdgeDetect bool   // Use edge detection
	Color      bool   // Preserve colors (ANSI)
	Dither     bool   // Apply dithering
}

// DefaultOptions returns sensible defaults
func DefaultOptions() Options {
	return Options{
		Width:   80,
		Charset: " .:-=+*#%@",
	}
}

// Character sets for different styles - leveraging full Unicode
var CharSets = map[string]string{
	// Basic sets
	"standard": " .:-=+*#%@",
	"simple":   " .*#",
	"binary":   " █",
	"ascii":    " .,:;i1tfLCG08@",

	// Extended ASCII with more detail
	"detailed": " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$",

	// Block elements (U+2580-259F)
	"blocks": " ░▒▓█",
	"shade":  " ░▒▓█▄▀▌▐",

	// Braille patterns (U+2800-28FF) - 256 patterns for high resolution
	"braille": " ⠁⠂⠃⠄⠅⠆⠇⡀⡁⡂⡃⡄⡅⡆⡇⠈⠉⠊⠋⠌⠍⠎⠏⡈⡉⡊⡋⡌⡍⡎⡏⠐⠑⠒⠓⠔⠕⠖⠗⡐⡑⡒⡓⡔⡕⡖⡗⠘⠙⠚⠛⠜⠝⠞⠟⡘⡙⡚⡛⡜⡝⡞⡟⠠⠡⠢⠣⠤⠥⠦⠧⡠⡡⡢⡣⡤⡥⡦⡧⠨⠩⠪⠫⠬⠭⠮⠯⡨⡩⡪⡫⡬⡭⡮⡯⠰⠱⠲⠳⠴⠵⠶⠷⡰⡱⡲⡳⡴⡵⡶⡷⠸⠹⠺⠻⠼⠽⠾⠿⡸⡹⡺⡻⡼⡽⡾⡿⢀⢁⢂⢃⢄⢅⢆⢇⣀⣁⣂⣃⣄⣅⣆⣇⢈⢉⢊⢋⢌⢍⢎⢏⣈⣉⣊⣋⣌⣍⣎⣏⢐⢑⢒⢓⢔⢕⢖⢗⣐⣑⣒⣓⣔⣕⣖⣗⢘⢙⢚⢛⢜⢝⢞⢟⣘⣙⣚⣛⣜⣝⣞⣟⢠⢡⢢⢣⢤⢥⢦⢧⣠⣡⣢⣣⣤⣥⣦⣧⢨⢩⢪⢫⢬⢭⢮⢯⣨⣩⣪⣫⣬⣭⣮⣯⢰⢱⢲⢳⢴⢵⢶⢷⣰⣱⣲⣳⣴⣵⣶⣷⢸⢹⢺⢻⢼⢽⢾⢿⣸⣹⣺⣻⣼⣽⣾⣿",

	// Box drawing characters (U+2500-257F)
	"box": " ─│┌┐└┘├┤┬┴┼━┃┏┓┗┛┣┫┳┻╋╔╗╚╝╠╣╦╩╬░▒▓█",

	// Geometric shapes (U+25A0-25FF)
	"geometric": " ◦◌○◎●◐◑◒◓◔◕◖◗◘◙◚◛◜◝◞◟◠◡◢◣◤◥◦◧◨◩◪◫◬◭◮◯",

	// Filled geometric with squares and circles
	"shapes": " ·∘○◌◎●◐◑▪▫◻◼◽◾▢▣▤▥▦▧▨▩■□▬▭▮▯",

	// Mathematical operators and symbols
	"math": " ·∙∘∴∵∶∷∸∹∺∻∼∽∾∿≀≁≂≃≄≅≆≇≈≉≊≋≌≍≎≏≐≑≒≓≔≕≖≗≘≙≚≛≜≝≞≟≠≡≢≣≤≥",

	// Stars and asterisks
	"stars": " ·✦✧★☆✡✢✣✤✥✦✧✩✪✫✬✭✮✯✰✱✲✳✴✵✶✷✸✹✺✻✼✽✾✿❀❁❂❃❄❅❆❇❈❉❊❋",

	// Dingbats and symbols
	"dingbats": " ·•‣⁃⁌⁍∎∗※⁂⁎⁑⁕❖❗❘❙❚❛❜❝❞❟❠❡❢❣❤❥❦❧",

	// Arrows
	"arrows": " →↗↑↖←↙↓↘↔↕↯⇐⇑⇒⇓⇔⇕⇖⇗⇘⇙⇚⇛⇜⇝⇞⇟⇠⇡⇢⇣⇤⇥⇦⇧⇨⇩⇪",

	// Musical symbols
	"music": " ♩♪♫♬♭♮♯𝄞𝄢𝄪𝄫",

	// Card suits and misc symbols
	"symbols": " ·‥…※†‡•‣⁂⁃⁎⁑⁕♠♡♢♣♤♥♦♧♨♩♪♫♬☀☁☂☃☄★☆☇☈☉☊☋☌☍☎☏",

	// Dense unicode for maximum detail
	"ultra": " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$░▒▓█▀▄▌▐●○◐◑◒◓■□▪▫◻◼★☆✦✧",

	// Retro/terminal style
	"retro": " ·:;+*%#@█▓▒░",

	// Dots and circles progression
	"dots": " ·∘○◌◎●",

	// Elegant/minimal
	"elegant": " ·‥…∴∵∶∷⁘⁙⁚⁛⁜⁝⁞",

	// CJK-inspired density
	"dense": " 一二三亖卅卌╳╬╫╪╩╨╧╦╥╤╣╢╡╠╟╞╝╜╛╚╙╘╗╖╕╔╓╒║═",

	// Hatching patterns
	"hatch": " /\\|─│┌┐└┘├┤┬┴┼╱╲╳",

	// Stipple/noise pattern
	"stipple": " ⠀⠁⠂⠄⠈⠐⠠⡀⢀⣀⣠⣤⣴⣶⣷⣿",

	// Legacy/compatibility - keeping old name
	"dots_braille": " ⠁⠂⠃⠄⠅⠆⠇⡀⡁⡂⡃⡄⡅⡆⡇⠈⠉⠊⠋⠌⠍⠎⠏⡈⡉⡊⡋⡌⡍⡎⡏⠐⠑⠒⠓⠔⠕⠖⠗⡐⡑⡒⡓⡔⡕⡖⡗⠘⠙⠚⠛⠜⠝⠞⠟⡘⡙⡚⡛⡜⡝⡞⡟⠠⠡⠢⠣⠤⠥⠦⠧⡠⡡⡢⡣⡤⡥⡦⡧⠨⠩⠪⠫⠬⠭⠮⠯⡨⡩⡪⡫⡬⡭⡮⡯⠰⠱⠲⠳⠴⠵⠶⠷⡰⡱⡲⡳⡴⡵⡶⡷⠸⠹⠺⠻⠼⠽⠾⠿⡸⡹⡺⡻⡼⡽⡾⡿⢀⢁⢂⢃⢄⢅⢆⢇⣀⣁⣂⣃⣄⣅⣆⣇⢈⢉⢊⢋⢌⢍⢎⢏⣈⣉⣊⣋⣌⣍⣎⣏⢐⢑⢒⢓⢔⢕⢖⢗⣐⣑⣒⣓⣔⣕⣖⣗⢘⢙⢚⢛⢜⢝⢞⢟⣘⣙⣚⣛⣜⣝⣞⣟⢠⢡⢢⢣⢤⢥⢦⢧⣠⣡⣢⣣⣤⣥⣦⣧⢨⢩⢪⢫⢬⢭⢮⢯⣨⣩⣪⣫⣬⣭⣮⣯⢰⢱⢲⢳⢴⢵⢶⢷⣰⣱⣲⣳⣴⣵⣶⣷⢸⢹⢺⢻⢼⢽⢾⢿⣸⣹⣺⣻⣼⣽⣾⣿",
}

// LoadImageFile loads an image from a file path
func LoadImageFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

// LoadImageURL loads an image from a URL
func LoadImageURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: HTTP %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

func FromFile(path string, opts Options) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	return FromReader(file, opts)
}

// FromURL converts an image from URL to ASCII art
func FromURL(url string, opts Options) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: HTTP %d", resp.StatusCode)
	}

	return FromReader(resp.Body, opts)
}

// FromReader converts an image from io.Reader to ASCII art
func FromReader(r io.Reader, opts Options) (string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	return FromImage(img, opts)
}

// FromImage converts an image.Image to ASCII art
// For large images (above threshold), automatically uses parallel processing
func FromImage(img image.Image, opts Options) (string, error) {
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	// Check if image is large enough to benefit from parallelization
	pixelCount := imgWidth * imgHeight
	if pixelCount >= parallelConfig.threshold {
		// Use parallel processing for large images
		return FromImageParallel(img, opts)
	}

	// Sequential processing for small images
	return fromImageSequential(img, opts)
}

// fromImageSequential performs the actual sequential pixel-by-pixel conversion
// This is the core algorithm used by both FromImage and FromImageParallel
func fromImageSequential(img image.Image, opts Options) (string, error) {
	if opts.Charset == "" {
		opts.Charset = CharSets["standard"]
	}

	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	// Calculate target dimensions
	targetWidth := opts.Width
	if targetWidth <= 0 {
		targetWidth = 80
	}

	// Characters are roughly 2x taller than wide, so adjust
	aspectRatio := float64(imgWidth) / float64(imgHeight)
	targetHeight := opts.Height
	if targetHeight <= 0 {
		targetHeight = int(float64(targetWidth) / aspectRatio / 2.0)
	}

	// Calculate sampling step
	stepX := float64(imgWidth) / float64(targetWidth)
	stepY := float64(imgHeight) / float64(targetHeight)

	// Precompute grayscale if edge detection
	var edges [][]float64
	if opts.EdgeDetect {
		edges = detectEdges(img, targetWidth, targetHeight, stepX, stepY)
	}

	var result strings.Builder
	chars := []rune(opts.Charset)
	numChars := len(chars)

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// Sample the image region
			sampleX := int(float64(x)*stepX) + bounds.Min.X
			sampleY := int(float64(y)*stepY) + bounds.Min.Y

			// Ensure within bounds
			if sampleX >= bounds.Max.X {
				sampleX = bounds.Max.X - 1
			}
			if sampleY >= bounds.Max.Y {
				sampleY = bounds.Max.Y - 1
			}

			var brightness float64
			var r8, g8, b8 uint8

			if opts.EdgeDetect && edges != nil {
				brightness = edges[y][x]
			} else {
				// Get average color of the region
				r8, g8, b8, brightness = sampleRegion(img, sampleX, sampleY, int(stepX), int(stepY))
			}

			if opts.Invert {
				brightness = 1.0 - brightness
			}

			// Map brightness to character
			charIdx := int(brightness * float64(numChars-1))
			if charIdx >= numChars {
				charIdx = numChars - 1
			}
			if charIdx < 0 {
				charIdx = 0
			}

			char := chars[charIdx]

			if opts.Color && !opts.EdgeDetect {
				// ANSI 24-bit color
				result.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", r8, g8, b8, char))
			} else {
				result.WriteRune(char)
			}
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

// sampleRegion samples a region of the image and returns RGB and brightness
func sampleRegion(img image.Image, x, y, width, height int) (uint8, uint8, uint8, float64) {
	bounds := img.Bounds()

	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	var totalR, totalG, totalB float64
	var count float64

	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px := x + dx
			py := y + dy
			if px >= bounds.Max.X || py >= bounds.Max.Y {
				continue
			}

			c := img.At(px, py)
			r, g, b, _ := c.RGBA()

			// Convert from 16-bit to 8-bit
			totalR += float64(r >> 8)
			totalG += float64(g >> 8)
			totalB += float64(b >> 8)
			count++
		}
	}

	if count == 0 {
		return 0, 0, 0, 0
	}

	avgR := totalR / count
	avgG := totalG / count
	avgB := totalB / count

	// Calculate perceived brightness (human eye is more sensitive to green)
	brightness := (0.299*avgR + 0.587*avgG + 0.114*avgB) / 255.0

	return uint8(avgR), uint8(avgG), uint8(avgB), brightness
}

// detectEdges applies Sobel edge detection
func detectEdges(img image.Image, width, height int, stepX, stepY float64) [][]float64 {
	bounds := img.Bounds()

	// First, create grayscale version at target resolution
	gray := make([][]float64, height)
	for y := 0; y < height; y++ {
		gray[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			sampleX := int(float64(x)*stepX) + bounds.Min.X
			sampleY := int(float64(y)*stepY) + bounds.Min.Y

			if sampleX >= bounds.Max.X {
				sampleX = bounds.Max.X - 1
			}
			if sampleY >= bounds.Max.Y {
				sampleY = bounds.Max.Y - 1
			}

			_, _, _, brightness := sampleRegion(img, sampleX, sampleY, int(stepX), int(stepY))
			gray[y][x] = brightness
		}
	}

	// Apply Sobel operator
	edges := make([][]float64, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				edges[y][x] = 0
				continue
			}

			// Sobel kernels
			gx := -gray[y-1][x-1] + gray[y-1][x+1] +
				-2*gray[y][x-1] + 2*gray[y][x+1] +
				-gray[y+1][x-1] + gray[y+1][x+1]

			gy := -gray[y-1][x-1] - 2*gray[y-1][x] - gray[y-1][x+1] +
				gray[y+1][x-1] + 2*gray[y+1][x] + gray[y+1][x+1]

			magnitude := math.Sqrt(gx*gx + gy*gy)

			// Normalize and invert (edges should be dark characters)
			edges[y][x] = 1.0 - math.Min(magnitude*2, 1.0)
		}
	}

	return edges
}

// GetCharset returns a charset by name
func GetCharset(name string) string {
	if cs, ok := CharSets[name]; ok {
		return cs
	}
	return CharSets["standard"]
}

// ListCharsets returns all available charset names
func ListCharsets() []string {
	return []string{
		"standard", "simple", "ascii", "detailed", "binary",
		"blocks", "shade", "braille", "stipple",
		"box", "hatch",
		"geometric", "shapes", "dots",
		"stars", "dingbats", "symbols", "arrows", "music",
		"math", "elegant", "dense",
		"ultra", "retro",
	}
}
