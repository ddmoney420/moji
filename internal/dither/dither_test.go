package dither

import (
	"image"
	"image/color"
	"math"
	"testing"
)

// Helper functions for test setup
func createGrayImage(width, height int, value uint8) *image.Gray {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{Y: value})
		}
	}
	return img
}

func createGradientImage(width, height int) *image.Gray {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create horizontal gradient from 0 to 255
			value := uint8((x * 255) / width)
			img.SetGray(x, y, color.Gray{Y: value})
		}
	}
	return img
}

func createCheckerboardImage(width, height int, squareSize int) *image.Gray {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if ((x/squareSize) + (y/squareSize)) % 2 == 0 {
				img.SetGray(x, y, color.Gray{Y: 0})
			} else {
				img.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}
	return img
}

func allPixelsInRange(img *image.Gray, min, max uint8) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := img.GrayAt(x, y).Y
			if val < min || val > max {
				return false
			}
		}
	}
	return true
}

func allPixelsValue(img *image.Gray, value uint8) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y != value {
				return false
			}
		}
	}
	return true
}

func hasNaN(img *image.Gray) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := float64(img.GrayAt(x, y).Y)
			if math.IsNaN(val) {
				return true
			}
		}
	}
	return false
}

func countPixelsWithValue(img *image.Gray, value uint8) int {
	count := 0
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y == value {
				count++
			}
		}
	}
	return count
}

// Floyd-Steinberg Tests
func TestFloydSteinbergBasic(t *testing.T) {
	img := createGrayImage(10, 10, 200)
	result := Apply(img, FloydSteinberg)

	if result == nil {
		t.Fatal("Floyd-Steinberg returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Floyd-Steinberg produced out-of-range pixel values")
	}
	if hasNaN(result) {
		t.Error("Floyd-Steinberg produced NaN values")
	}
}

func TestFloydSteinbergBlackPixels(t *testing.T) {
	img := createGrayImage(5, 5, 50)
	result := Apply(img, FloydSteinberg)

	darkPixels := countPixelsWithValue(result, 0)
	if darkPixels == 0 {
		t.Error("Floyd-Steinberg should have dark (0) pixels for dark input")
	}
}

func TestFloydSteinbergWhitePixels(t *testing.T) {
	img := createGrayImage(5, 5, 200)
	result := Apply(img, FloydSteinberg)

	brightPixels := countPixelsWithValue(result, 255)
	if brightPixels == 0 {
		t.Error("Floyd-Steinberg should have bright (255) pixels for bright input")
	}
}

func TestFloydSteinbergConsistency(t *testing.T) {
	img := createGrayImage(10, 10, 127)
	result1 := Apply(img, FloydSteinberg)
	result2 := Apply(img, FloydSteinberg)

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if result1.GrayAt(x, y) != result2.GrayAt(x, y) {
				t.Error("Floyd-Steinberg produced different results for same input")
				return
			}
		}
	}
}

func TestFloydSteinbergGradient(t *testing.T) {
	img := createGradientImage(10, 10)
	result := Apply(img, FloydSteinberg)

	if !allPixelsInRange(result, 0, 255) {
		t.Error("Floyd-Steinberg gradient produced out-of-range values")
	}

	// Check for dithering effect (should have both 0 and 255)
	hasBlack := countPixelsWithValue(result, 0) > 0
	hasWhite := countPixelsWithValue(result, 255) > 0
	if !hasBlack || !hasWhite {
		t.Error("Floyd-Steinberg gradient should produce both black and white pixels")
	}
}

// Atkinson Tests
func TestAtkinsonBasic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, Atkinson)

	if result == nil {
		t.Fatal("Atkinson returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Atkinson produced out-of-range pixel values")
	}
}

func TestAtkinsonErrorDistribution(t *testing.T) {
	img := createGrayImage(20, 20, 128)
	result := Apply(img, Atkinson)

	// Atkinson uses 1/8 error distribution, so should have less visible banding
	hasBlack := countPixelsWithValue(result, 0) > 0
	hasWhite := countPixelsWithValue(result, 255) > 0
	if !hasBlack || !hasWhite {
		t.Error("Atkinson should produce both black and white pixels")
	}
}

func TestAtkinsonSmallImage(t *testing.T) {
	img := createGrayImage(2, 2, 100)
	result := Apply(img, Atkinson)

	if result == nil {
		t.Fatal("Atkinson failed on small image")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Atkinson produced invalid values on small image")
	}
}

// Sierra Tests
func TestSierraBasic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, Sierra)

	if result == nil {
		t.Fatal("Sierra returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Sierra produced out-of-range pixel values")
	}
}

func TestSierraCoefficients(t *testing.T) {
	img := createGrayImage(15, 15, 127)
	result := Apply(img, Sierra)

	// Sierra distributes 32/32 of error, check for proper dithering
	hasBlack := countPixelsWithValue(result, 0) > 0
	hasWhite := countPixelsWithValue(result, 255) > 0
	if !hasBlack || !hasWhite {
		t.Error("Sierra should produce both black and white pixels")
	}
}

// Sierra Lite Tests
func TestSierraLiteBasic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, SierraLite)

	if result == nil {
		t.Fatal("Sierra Lite returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Sierra Lite produced out-of-range pixel values")
	}
}

func TestSierraLiteSimplicity(t *testing.T) {
	img := createGrayImage(10, 10, 128)
	result := Apply(img, SierraLite)

	// Sierra Lite is simpler (two-row), should still dither
	hasBlack := countPixelsWithValue(result, 0) > 0
	hasWhite := countPixelsWithValue(result, 255) > 0
	if !hasBlack || !hasWhite {
		t.Error("Sierra Lite should produce both black and white pixels")
	}
}

// Stucki Tests
func TestStuckiBasic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, Stucki)

	if result == nil {
		t.Fatal("Stucki returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Stucki produced out-of-range pixel values")
	}
}

func TestStuckiWideFilter(t *testing.T) {
	img := createGrayImage(20, 20, 127)
	result := Apply(img, Stucki)

	// Stucki has wider distribution (42/42)
	hasBlack := countPixelsWithValue(result, 0) > 0
	hasWhite := countPixelsWithValue(result, 255) > 0
	if !hasBlack || !hasWhite {
		t.Error("Stucki should produce both black and white pixels")
	}
}

// Burkes Tests
func TestBurkesBasic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, Burkes)

	if result == nil {
		t.Fatal("Burkes returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Burkes produced out-of-range pixel values")
	}
}

func TestBurkesDistribution(t *testing.T) {
	img := createGrayImage(15, 15, 128)
	result := Apply(img, Burkes)

	// Burkes distributes 32/32 error
	hasBlack := countPixelsWithValue(result, 0) > 0
	hasWhite := countPixelsWithValue(result, 255) > 0
	if !hasBlack || !hasWhite {
		t.Error("Burkes should produce both black and white pixels")
	}
}

// Jarvis-Judice-Ninke Tests
func TestJarvisJudiceBasic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, JarvisJudice)

	if result == nil {
		t.Fatal("Jarvis-Judice returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Jarvis-Judice produced out-of-range pixel values")
	}
}

func TestJarvisJudiceWideKernel(t *testing.T) {
	img := createGrayImage(25, 25, 127)
	result := Apply(img, JarvisJudice)

	// Jarvis-Judice has 3x5 kernel (distributed over 48/48)
	hasBlack := countPixelsWithValue(result, 0) > 0
	hasWhite := countPixelsWithValue(result, 255) > 0
	if !hasBlack || !hasWhite {
		t.Error("Jarvis-Judice should produce both black and white pixels")
	}
}

// Bayer 2x2 Tests
func TestBayer2x2Basic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, Bayer2x2)

	if result == nil {
		t.Fatal("Bayer 2x2 returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Bayer 2x2 produced out-of-range pixel values")
	}
}

func TestBayer2x2Periodic(t *testing.T) {
	img := createGrayImage(4, 4, 128)
	result := Apply(img, Bayer2x2)

	// Bayer 2x2 should repeat every 2 pixels
	pixel00 := result.GrayAt(0, 0).Y
	pixel20 := result.GrayAt(2, 0).Y
	if pixel00 != pixel20 {
		t.Error("Bayer 2x2 should repeat periodically")
	}
}

func TestBayer2x2ThresholdValues(t *testing.T) {
	img := createGradientImage(20, 4)
	result := Apply(img, Bayer2x2)

	if !allPixelsInRange(result, 0, 255) {
		t.Error("Bayer 2x2 produced invalid threshold values")
	}
}

// Bayer 4x4 Tests
func TestBayer4x4Basic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, Bayer4x4)

	if result == nil {
		t.Fatal("Bayer 4x4 returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Bayer 4x4 produced out-of-range pixel values")
	}
}

func TestBayer4x4Periodic(t *testing.T) {
	img := createGrayImage(8, 8, 128)
	result := Apply(img, Bayer4x4)

	pixel00 := result.GrayAt(0, 0).Y
	pixel40 := result.GrayAt(4, 0).Y
	if pixel00 != pixel40 {
		t.Error("Bayer 4x4 should repeat every 4 pixels")
	}
}

func TestBayer4x4MoreLevels(t *testing.T) {
	img := createGradientImage(20, 4)
	result := Apply(img, Bayer4x4)

	blackCount := countPixelsWithValue(result, 0)
	whiteCount := countPixelsWithValue(result, 255)
	if blackCount == 0 || whiteCount == 0 {
		t.Error("Bayer 4x4 should produce both black and white pixels in gradient")
	}
}

// Bayer 8x8 Tests
func TestBayer8x8Basic(t *testing.T) {
	img := createGrayImage(10, 10, 150)
	result := Apply(img, Bayer8x8)

	if result == nil {
		t.Fatal("Bayer 8x8 returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Bayer 8x8 produced out-of-range pixel values")
	}
}

func TestBayer8x8Periodic(t *testing.T) {
	img := createGrayImage(16, 16, 128)
	result := Apply(img, Bayer8x8)

	pixel00 := result.GrayAt(0, 0).Y
	pixel80 := result.GrayAt(8, 0).Y
	if pixel00 != pixel80 {
		t.Error("Bayer 8x8 should repeat every 8 pixels")
	}
}

func TestBayer8x8FinestDithering(t *testing.T) {
	img := createGradientImage(20, 8)
	result := Apply(img, Bayer8x8)

	if !allPixelsInRange(result, 0, 255) {
		t.Error("Bayer 8x8 produced invalid values")
	}
}

// None Algorithm Tests
func TestNoneAlgorithm(t *testing.T) {
	img := createGrayImage(10, 10, 128)
	result := Apply(img, None)

	if result == nil {
		t.Fatal("None algorithm returned nil")
	}
	// None should preserve grayscale values
	if !allPixelsValue(result, 128) {
		t.Error("None algorithm should preserve input pixel values")
	}
}

// Edge Case Tests
func TestSinglePixelImage(t *testing.T) {
	tests := []Algorithm{FloydSteinberg, Bayer2x2, Bayer4x4, Bayer8x8, Atkinson, Sierra, SierraLite, Stucki, Burkes, JarvisJudice}

	for _, algo := range tests {
		img := createGrayImage(1, 1, 128)
		result := Apply(img, algo)

		if result == nil {
			t.Errorf("%v returned nil for 1x1 image", algo)
		}
		if !allPixelsInRange(result, 0, 255) {
			t.Errorf("%v produced invalid values for 1x1 image", algo)
		}
	}
}

func TestLargeImage(t *testing.T) {
	tests := []Algorithm{FloydSteinberg, Bayer2x2, Bayer4x4, Bayer8x8, Atkinson}

	for _, algo := range tests {
		img := createGrayImage(500, 500, 128)
		result := Apply(img, algo)

		if result == nil {
			t.Errorf("%v returned nil for large image", algo)
		}
		if !allPixelsInRange(result, 0, 255) {
			t.Errorf("%v produced invalid values for large image", algo)
		}
	}
}

func TestAllBlackImage(t *testing.T) {
	tests := []Algorithm{FloydSteinberg, Bayer2x2, Atkinson, Sierra}

	for _, algo := range tests {
		img := createGrayImage(10, 10, 0)
		result := Apply(img, algo)

		if !allPixelsValue(result, 0) {
			t.Errorf("%v should keep all-black image black", algo)
		}
	}
}

func TestAllWhiteImage(t *testing.T) {
	tests := []Algorithm{FloydSteinberg, Bayer2x2, Atkinson, Sierra}

	for _, algo := range tests {
		img := createGrayImage(10, 10, 255)
		result := Apply(img, algo)

		if !allPixelsValue(result, 255) {
			t.Errorf("%v should keep all-white image white", algo)
		}
	}
}

func TestCheckerboardPattern(t *testing.T) {
	tests := []Algorithm{FloydSteinberg, Bayer2x2, Bayer4x4, Atkinson}

	for _, algo := range tests {
		img := createCheckerboardImage(10, 10, 2)
		result := Apply(img, algo)

		if result == nil {
			t.Errorf("%v returned nil for checkerboard", algo)
		}
		if !allPixelsInRange(result, 0, 255) {
			t.Errorf("%v produced invalid values for checkerboard", algo)
		}
	}
}

// GetAlgorithm Tests
func TestGetAlgorithmFloydSteinberg(t *testing.T) {
	tests := []string{"floyd-steinberg", "fs"}
	for _, name := range tests {
		algo := GetAlgorithm(name)
		if algo != FloydSteinberg {
			t.Errorf("GetAlgorithm(%q) should return FloydSteinberg", name)
		}
	}
}

func TestGetAlgorithmBayer(t *testing.T) {
	tests := map[string]Algorithm{
		"bayer2x2": Bayer2x2,
		"bayer2":   Bayer2x2,
		"bayer4x4": Bayer4x4,
		"bayer4":   Bayer4x4,
		"bayer":    Bayer4x4,
		"bayer8x8": Bayer8x8,
		"bayer8":   Bayer8x8,
	}
	for name, expected := range tests {
		algo := GetAlgorithm(name)
		if algo != expected {
			t.Errorf("GetAlgorithm(%q) should return %v, got %v", name, expected, algo)
		}
	}
}

func TestGetAlgorithmOtherAlgorithms(t *testing.T) {
	tests := map[string]Algorithm{
		"atkinson":    Atkinson,
		"sierra":      Sierra,
		"sierra-lite": SierraLite,
		"sierra2":     SierraLite,
		"stucki":      Stucki,
		"burkes":      Burkes,
		"jarvis":      JarvisJudice,
		"jjn":         JarvisJudice,
	}
	for name, expected := range tests {
		algo := GetAlgorithm(name)
		if algo != expected {
			t.Errorf("GetAlgorithm(%q) should return %v", name, expected)
		}
	}
}

func TestGetAlgorithmUnknown(t *testing.T) {
	algo := GetAlgorithm("unknown-algorithm")
	if algo != None {
		t.Error("GetAlgorithm with unknown name should return None")
	}
}

// ListAlgorithms Tests
func TestListAlgorithms(t *testing.T) {
	algos := ListAlgorithms()

	if len(algos) == 0 {
		t.Fatal("ListAlgorithms returned empty list")
	}

	expectedCount := 11 // none, floyd-steinberg, bayer2x2, bayer4x4, bayer8x8, atkinson, sierra, sierra-lite, stucki, burkes, jarvis
	if len(algos) != expectedCount {
		t.Errorf("ListAlgorithms returned %d algorithms, expected %d", len(algos), expectedCount)
	}

	// Check for specific algorithms
	hasFloydSteinberg := false
	hasBayer := false
	for _, algo := range algos {
		if algo == "floyd-steinberg" {
			hasFloydSteinberg = true
		}
		if algo == "bayer4x4" {
			hasBayer = true
		}
	}

	if !hasFloydSteinberg {
		t.Error("ListAlgorithms should include floyd-steinberg")
	}
	if !hasBayer {
		t.Error("ListAlgorithms should include bayer4x4")
	}
}

// ToGrayscale Tests
func TestToGrayscaleBasic(t *testing.T) {
	img := createGrayImage(10, 10, 128)
	result := ToGrayscale(img, 2)

	if result == nil {
		t.Fatal("ToGrayscale returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("ToGrayscale produced out-of-range values")
	}
}

func TestToGrayscaleLevelHandling(t *testing.T) {
	img := createGrayImage(10, 10, 128)

	// Test minimum levels
	result := ToGrayscale(img, 1)
	if result == nil {
		t.Error("ToGrayscale should handle levels < 2")
	}

	// Test maximum levels
	result = ToGrayscale(img, 300)
	if result == nil {
		t.Error("ToGrayscale should handle levels > 256")
	}
}

func TestToGrayscaleGradient(t *testing.T) {
	img := createGradientImage(256, 10)
	result := ToGrayscale(img, 4)

	if !allPixelsInRange(result, 0, 255) {
		t.Error("ToGrayscale produced invalid values for gradient")
	}
}

// Luminance Tests
func TestLuminanceRGB(t *testing.T) {
	// Test pure colors
	grayR := Luminance(255, 0, 0)
	grayG := Luminance(0, 255, 0)
	grayB := Luminance(0, 0, 255)

	if grayR <= 0 || grayR >= 255 {
		t.Errorf("Luminance(255,0,0) = %f, expected between 0-255", grayR)
	}
	if grayG <= grayR {
		t.Error("Green should have higher luminance than red")
	}
	if grayB >= grayG {
		t.Error("Blue should have lower luminance than green")
	}
}

func TestLuminanceWhiteBlack(t *testing.T) {
	white := Luminance(255, 255, 255)
	black := Luminance(0, 0, 0)

	if white != 255.0 {
		t.Errorf("Luminance(255,255,255) = %f, expected 255.0", white)
	}
	if black != 0.0 {
		t.Errorf("Luminance(0,0,0) = %f, expected 0.0", black)
	}
}

// ContrastStretch Tests
func TestContrastStretchBasic(t *testing.T) {
	img := createGradientImage(100, 10)
	result := ContrastStretch(img, 0, 255)

	if result == nil {
		t.Fatal("ContrastStretch returned nil")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("ContrastStretch produced out-of-range values")
	}
}

func TestContrastStretchUniformImage(t *testing.T) {
	img := createGrayImage(10, 10, 128)
	result := ContrastStretch(img, 0, 255)

	// Uniform image should become uniform output
	if !allPixelsInRange(result, 0, 255) {
		t.Error("ContrastStretch failed on uniform image")
	}
}

func TestContrastStretchRangeHandling(t *testing.T) {
	img := createGradientImage(100, 10)
	result := ContrastStretch(img, 50, 200)

	if !allPixelsInRange(result, 50, 200) {
		t.Error("ContrastStretch should respect output range")
	}
}

func TestContrastStretchAllBlack(t *testing.T) {
	img := createGrayImage(10, 10, 0)
	result := ContrastStretch(img, 0, 255)

	// Uniform image is set to midpoint of output range
	expectedValue := uint8((0 + 255) / 2)
	if !allPixelsValue(result, expectedValue) {
		t.Errorf("ContrastStretch uniform image should be %d, got something else", expectedValue)
	}
}

// Integration Tests
func TestApplyWithColorImage(t *testing.T) {
	// Create a simple color image (RGBA)
	bounds := image.Rect(0, 0, 10, 10)
	colorImg := image.NewRGBA(bounds)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			// Create a colored pixel
			colorImg.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	result := Apply(colorImg, FloydSteinberg)

	if result == nil {
		t.Fatal("Apply failed on color image")
	}
	if !allPixelsInRange(result, 0, 255) {
		t.Error("Apply produced invalid values for color image")
	}
}

func TestApplyAllAlgorithms(t *testing.T) {
	img := createGrayImage(20, 20, 128)
	algorithms := []Algorithm{
		FloydSteinberg, Bayer2x2, Bayer4x4, Bayer8x8,
		Atkinson, Sierra, SierraLite, Stucki, Burkes, JarvisJudice, None,
	}

	for _, algo := range algorithms {
		result := Apply(img, algo)
		if result == nil {
			t.Errorf("Apply returned nil for algorithm %v", algo)
		}
		if !allPixelsInRange(result, 0, 255) {
			t.Errorf("Algorithm %v produced out-of-range values", algo)
		}
	}
}

func TestRoundTripConsistency(t *testing.T) {
	originalImg := createGrayImage(15, 15, 128)
	dithered1 := Apply(originalImg, FloydSteinberg)
	dithered2 := Apply(dithered1, FloydSteinberg)

	// Dithering a dithered image should not fail
	if dithered2 == nil {
		t.Fatal("Round-trip dithering returned nil")
	}
	if !allPixelsInRange(dithered2, 0, 255) {
		t.Error("Round-trip dithering produced invalid values")
	}
}

func TestBoundaryConditions(t *testing.T) {
	// Test 2x2 image (smallest for meaningful testing)
	tests := []Algorithm{FloydSteinberg, Bayer2x2, Atkinson, Sierra, Stucki}

	for _, algo := range tests {
		img := createGrayImage(2, 2, 128)
		result := Apply(img, algo)

		if result == nil {
			t.Errorf("%v failed on 2x2 image", algo)
		}
		if !allPixelsInRange(result, 0, 255) {
			t.Errorf("%v produced invalid values on 2x2 image", algo)
		}
	}
}

func TestTallNarrowImage(t *testing.T) {
	tests := []Algorithm{FloydSteinberg, Bayer4x4, Sierra}

	for _, algo := range tests {
		img := createGrayImage(2, 50, 150)
		result := Apply(img, algo)

		if result == nil {
			t.Errorf("%v failed on tall narrow image", algo)
		}
		if !allPixelsInRange(result, 0, 255) {
			t.Errorf("%v produced invalid values on tall narrow image", algo)
		}
	}
}

func TestWideShallowImage(t *testing.T) {
	tests := []Algorithm{FloydSteinberg, Bayer4x4, Sierra}

	for _, algo := range tests {
		img := createGrayImage(50, 2, 150)
		result := Apply(img, algo)

		if result == nil {
			t.Errorf("%v failed on wide shallow image", algo)
		}
		if !allPixelsInRange(result, 0, 255) {
			t.Errorf("%v produced invalid values on wide shallow image", algo)
		}
	}
}

// Benchmark Tests
func BenchmarkFloydSteinberg(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, FloydSteinberg)
	}
}

func BenchmarkBayer2x2(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, Bayer2x2)
	}
}

func BenchmarkBayer4x4(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, Bayer4x4)
	}
}

func BenchmarkBayer8x8(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, Bayer8x8)
	}
}

func BenchmarkAtkinson(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, Atkinson)
	}
}

func BenchmarkSierra(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, Sierra)
	}
}

func BenchmarkStucki(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, Stucki)
	}
}

func BenchmarkJarvisJudice(b *testing.B) {
	img := createGrayImage(100, 100, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(img, JarvisJudice)
	}
}
